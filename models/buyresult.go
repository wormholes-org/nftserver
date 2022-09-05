package models

import (
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) BuyResult(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty, txhash string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	contractAddr = strings.ToLower(contractAddr)

	fmt.Println("BuyResult() price = ", price)
	if IsPriceValid(price) != true {
		fmt.Println("BuyResult() price err")
		return ErrPrice
	}
	fmt.Println(time.Now().String()[:25], "BuyResult() Begin", "from=", from, "to=", to, "price=", price,
		"contractAddr=", contractAddr, "tokenId=", tokenId, "txhash=", txhash,
		"royalty=", royalty /*, "sig=", sig, "trade_sig=", trade_sig*/)
	fmt.Println("BuyResult()++q++++++++++++++++++")
	if royalty != "" {
		fmt.Println("BuyResult() royalty!=Null mint royalty=", royalty)
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
		if err.Error != nil {
			fmt.Println("BuyResult() royalty err =", ErrNftNotExist)
			return ErrNftNotExist
		}
		trans := Trans{}
		trans.Contract = contractAddr
		trans.Fromaddr = ""
		trans.Toaddr = to
		trans.Signdata = sig
		trans.Tradesig = trade_sig
		trans.Tokenid = tokenId
		trans.Price, _ = strconv.ParseUint(price, 10, 64)
		trans.Transtime = time.Now().Unix()
		trans.Selltype = SellTypeMintNft.String()
		trans.Name = nftRec.Name
		trans.Meta = nftRec.Meta
		trans.Desc = nftRec.Desc
		trans.Txhash = txhash
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("BuyResult() royalty create trans err=", err.Error)
				return err.Error
			}
			nftrecord := Nfts{}
			nftrecord.Signdata = sig

			nftrecord.Royalty, _ = strconv.Atoi(royalty)
			//nftrecord.Royalty = nftrecord.Royalty / 100
			nftrecord.Mintstate = Minted.String()
			err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
				contractAddr, tokenId).Updates(&nftrecord)
			if err.Error != nil {
				fmt.Println("BuyResult() royalty update nfts record err=", err.Error)
				return err.Error
			}
			fmt.Println("BuyResult() royalty!=Null Ok")
			return nil
		})
	}
	fmt.Println("BuyResult()-------------------")
	if from != "" && to != "" {
		fmt.Println("BuyResult() from != Null && to != Null")
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
		if err.Error != nil {
			fmt.Println("BuyResult() auction not find err=", err.Error)
			return ErrNftNotExist
		}
		if price == "" {
			fmt.Println("BuyResult() price == null")
			return nft.db.Transaction(func(tx *gorm.DB) error {
				var auctionRec Auction
				err = tx.Set("gorm:query_option", "FOR UPDATE").Where("contract = ? AND tokenid = ? AND ownaddr =?",
					contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
				if err.Error != nil {
					fmt.Println("BuyResult() auction not find err=", err.Error)
					return err.Error
				}
				trans := Trans{}
				trans.Auctionid = auctionRec.ID
				trans.Contract = auctionRec.Contract
				trans.Createaddr = nftRec.Createaddr
				trans.Fromaddr = from
				trans.Toaddr = to
				trans.Signdata = sig
				trans.Tradesig = trade_sig
				trans.Tokenid = auctionRec.Tokenid
				trans.Nftid = auctionRec.Nftid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Price = 0
				trans.Transtime = time.Now().Unix()
				trans.Selltype = SellTypeAsset.String()
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResult() create trans record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResult() price == null OK")
				return nil
			})
		} else {
			fmt.Println("BuyResult() price != null")
			return nft.db.Transaction(func(tx *gorm.DB) error {
				var auctionRec Auction
				err = tx.Where("contract = ? AND tokenid = ? AND ownaddr =?",
					contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
				if err.Error != nil {
					fmt.Println("BuyResult() auction not find err=", err.Error)
					return err.Error
				}
				trans := Trans{}
				trans.Auctionid = auctionRec.ID
				trans.Contract = auctionRec.Contract
				trans.Createaddr = nftRec.Createaddr
				trans.Fromaddr = from
				trans.Toaddr = to
				trans.Signdata = sig
				trans.Tradesig = trade_sig
				trans.Nftid = auctionRec.Nftid
				trans.Tokenid = auctionRec.Tokenid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Txhash = txhash
				trans.Name = nftRec.Name
				trans.Meta = nftRec.Meta
				trans.Desc = nftRec.Desc
				trans.Price, _ = strconv.ParseUint(price, 10, 64)
				trans.Transtime = time.Now().Unix()
				/*if auctionRec.Selltype == SellTypeWaitSale.String() {
					trans.Selltype = SellTypeHighestBid.String()
				}else {
					trans.Selltype = auctionRec.Selltype
				}*/
				trans.Selltype = auctionRec.Selltype
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResult() create trans record err=", err.Error)
					return err.Error
				}
				var collectRec Collects
				err = nft.db.Where("createaddr = ? AND  name=?",
					nftRec.Collectcreator, nftRec.Collections).First(&collectRec)
				if err.Error == nil {
					transCnt := collectRec.Transcnt + 1
					transAmt := collectRec.Transamt + trans.Price
					collectRec = Collects{}
					collectRec.Transcnt = transCnt
					collectRec.Transamt = transAmt
					err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
						nftRec.Collections, nftRec.Collectcreator).Updates(&collectRec)
					if err.Error != nil {
						fmt.Println("BuyResult() update collectRec err=", err.Error)
						return err.Error
					}
				}
				nftrecord := Nfts{}
				nftrecord.Ownaddr = to
				nftrecord.Selltype = SellTypeNotSale.String()
				nftrecord.Paychan = auctionRec.Paychan
				nftrecord.TransCur = auctionRec.Currency
				nftrecord.Transprice = trans.Price
				nftrecord.Transamt = nftRec.Transamt + trans.Price
				nftrecord.Transcnt = nftRec.Transcnt + 1
				nftrecord.Transtime = time.Now().Unix()
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
				if err.Error != nil {
					fmt.Println("BuyResult() update record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
				if err.Error != nil {
					fmt.Println("BuyResult() delete auction record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
				if err.Error != nil {
					fmt.Println("BuyResult() delete bid record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResult() from != Null && to != Null --> price != Null OK")
				return nil
			})
		}
	}
	fmt.Println("BuyResult() End.")
	return ErrFromToAddrZero
}

func (nft NftDb) BuyResultWithAmount(from, to, contractAddr, tokenId, amount, price, royalty, txhash, txtime string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	contractAddr = strings.ToLower(contractAddr)

	fmt.Println("BuyResultWithAmount() price = ", price)
	if IsUint64DataValid(price) != true {
		fmt.Println("BuyResultWithAmount() price err")
		return ErrPrice
	}
	fmt.Println(time.Now().String()[:25], "BuyResultWithAmount() Begin", "from=", from, "to=", to, "price=", price,
		"contractAddr=", contractAddr, "tokenId=", tokenId, "txhash=", txhash,
		"royalty=", royalty /*, "sig=", sig, "trade_sig=", trade_sig*/)
	trans := Trans{}
	err := nft.db.Select("id").Where("contract = ? AND tokenid = ? AND txhash = ? AND (selltype = ? or selltype = ? or selltype = ?)",
		contractAddr, tokenId, txhash, SellTypeFixPrice.String(), SellTypeBidPrice.String(), SellTypeHighestBid.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultWithAmount() trans exist.")
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultWithAmount() err =", err.Error)
		return err.Error
	}
	fmt.Println("BuyResultWithAmount()-------------------")
	if from != "" && to != "" && price != "" {
		fmt.Println("BuyResultWithAmount() from != Null && to != Null")
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			contractAddr, tokenId, from).First(&nftRec)
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				fmt.Println("BuyResultWithAmount() nft not find ")
				return nil
			}
			fmt.Println("BuyResultWithAmount() nft find err=", err.Error)
			return ErrNftNotExist
		}
		fmt.Println("BuyResultWithAmount() price != null")
		aucFlag := false
		var auctionRec Auction
		count, _ := strconv.Atoi(amount)
		if count > nftRec.Count {
			fmt.Println("BuyResultWithAmount() nft.count < amount.")
			return ErrNftAmount
		}
		err = nft.db.Where("contract = ? AND tokenid = ? AND ownaddr =? AND count = ?",
			contractAddr, tokenId, from, count).First(&auctionRec)
		if err.Error != nil {
			if err.Error != gorm.ErrRecordNotFound {
				fmt.Println("BuyResultWithAmount() auction not find err=", err.Error)
				return err.Error
			}
		} else {
			aucFlag = true
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			trans := Trans{}
			trans.Contract = contractAddr
			trans.Createaddr = nftRec.Createaddr
			trans.Count = count
			trans.Fromaddr = from
			trans.Toaddr = to
			trans.Tokenid = tokenId
			if aucFlag {
				trans.Auctionid = auctionRec.ID
				trans.Nftid = auctionRec.Nftid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Selltype = auctionRec.Selltype
			} else {
				trans.Selltype = SellTypeFixPrice.String()
			}
			trans.Txhash = txhash
			trans.Name = nftRec.Name
			trans.Meta = nftRec.Meta
			trans.Desc = nftRec.Desc
			trans.Price, _ = strconv.ParseUint(price, 10, 64)
			//trans.Transtime = time.Now().Unix()
			trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
			/*if auctionRec.Selltype == SellTypeWaitSale.String() {
				trans.Selltype = SellTypeHighestBid.String()
			}else {
				trans.Selltype = auctionRec.Selltype
			}*/
			err := tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("BuyResultWithAmount() create trans record err=", err.Error)
				return err.Error
			}
			var collectRec Collects
			err = nft.db.Where("createaddr = ? AND  name=?",
				nftRec.Collectcreator, nftRec.Collections).First(&collectRec)
			if err.Error == nil {
				transCnt := collectRec.Transcnt + 1
				transAmt := collectRec.Transamt + trans.Price
				collectRec = Collects{}
				collectRec.Transcnt = transCnt
				collectRec.Transamt = transAmt
				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
					nftRec.Collections, nftRec.Collectcreator).Updates(&collectRec)
				if err.Error != nil {
					fmt.Println("BuyResultWithAmount() update collectRec err=", err.Error)
					return err.Error
				}
			}
			var toRec Nfts
			err = nft.db.Where("contract = ? AND tokenid = ? AND ownaddr = ?",
				contractAddr, tokenId, to).First(&toRec)
			if err.Error != nil {
				if err.Error == gorm.ErrRecordNotFound {
					if nftRec.Count == count {
						nftrecord := Nfts{}
						nftrecord.Ownaddr = to
						nftrecord.Selltype = SellTypeNotSale.String()
						nftrecord.Transprice = trans.Price
						nftrecord.Transamt = nftRec.Transamt + trans.Price
						nftrecord.Transcnt = nftRec.Transcnt + 1
						nftrecord.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
						//nftrecord.Transtime = time.Now().Unix()
						//trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
						err = tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Updates(&nftrecord)
						if err.Error != nil {
							fmt.Println("BuyResultWithAmount() update record err=", err.Error)
							return err.Error
						}
					} else {
						nftrecord := Nfts{}
						nftrecord.Selltype = SellTypeNotSale.String()
						nftrecord.Transprice = trans.Price
						nftrecord.Transamt = nftRec.Transamt + trans.Price
						nftrecord.Transcnt = nftRec.Transcnt + 1
						nftrecord.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
						nftrecord.Count = nftRec.Count - count
						//nftrecord.Transtime = time.Now().Unix()
						//trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
						err = tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Updates(&nftrecord)
						if err.Error != nil {
							fmt.Println("BuyResultWithAmount() update record err=", err.Error)
							return err.Error
						}
						nftrecord = Nfts{}
						nftrecord = nftRec
						nftrecord.ID = 0
						nftrecord.Ownaddr = to
						nftrecord.Selltype = SellTypeNotSale.String()
						nftrecord.Transprice = 0
						nftrecord.Transamt = 0
						nftrecord.Transcnt = 0
						nftrecord.Count = count
						nftrecord.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
						//nftrecord.Transtime = time.Now().Unix()
						//trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
						err := tx.Model(&nftrecord).Create(&nftrecord)
						if err.Error != nil {
							fmt.Println("BuyResultWithAmount() create new nft err=", err.Error)
							return err.Error
						}
					}
				} else {
					fmt.Println("BuyResultWithAmount() dbase err=", err.Error)
					return ErrNftNotExist
				}
			} else {
				if nftRec.Count == count {
					nftrecord := Nfts{}
					nftrecord = toRec
					nftrecord.Selltype = SellTypeNotSale.String()
					nftrecord.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
					nftrecord.Transprice = trans.Price
					nftrecord.Transamt = nftRec.Transamt + trans.Price + toRec.Transamt
					nftrecord.Transcnt = nftRec.Transcnt + 1 + toRec.Transcnt
					nftrecord.Count = toRec.Count + count
					//nftrecord.Transtime = time.Now().Unix()
					//trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
					err = tx.Model(&Nfts{}).Where("id = ?", toRec.ID).Updates(&nftrecord)
					if err.Error != nil {
						fmt.Println("BuyResultWithAmount() update record err=", err.Error)
						return err.Error
					}
					err = tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Delete(&nftRec)
					if err.Error != nil {
						fmt.Println("BuyResultWithAmount() delete record err=", err.Error)
						return err.Error
					}
				} else {
					nftrecord := Nfts{}
					nftrecord.Selltype = SellTypeNotSale.String()
					nftrecord.Transprice = trans.Price
					nftrecord.Transamt = nftRec.Transamt + trans.Price
					nftrecord.Transcnt = nftRec.Transcnt + 1
					nftrecord.Count = nftRec.Count - count
					//nftrecord.Transtime = time.Now().Unix()
					nftrecord.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
					//trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
					err = tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Updates(&nftrecord)
					if err.Error != nil {
						fmt.Println("BuyResultWithAmount() update record err=", err.Error)
						return err.Error
					}
					nftrecord = Nfts{}
					nftrecord.Selltype = SellTypeNotSale.String()
					nftrecord.Transprice = trans.Price
					nftrecord.Count = toRec.Count + count
					//nftrecord.Transtime = time.Now().Unix()
					//nftrecord.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
					//trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
					err = tx.Model(&Nfts{}).Where("id = ?", toRec.ID).Updates(&nftrecord)
					if err.Error != nil {
						fmt.Println("BuyResultWithAmount() update record err=", err.Error)
						return err.Error
					}
				}
			}
			if aucFlag {
				err = tx.Model(&Auction{}).Where("id = ?", auctionRec.ID).Delete(&Auction{})
				if err.Error != nil {
					fmt.Println("BuyResultWithAmount() delete auction record err=", err.Error)
					return err.Error
				}
				err = nft.db.Model(&Bidding{}).Where("Auctionid = ?", auctionRec.ID).Delete(&Bidding{})
				if err.Error != nil {
					fmt.Println("BuyResultWithAmount() delete bid record err=", err.Error)
					return err.Error
				}
			}
			fmt.Println("BuyResultWithAmount() from != Null && to != Null --> price != Null OK")
			return nil
		})
	}
	fmt.Println("BuyResultWithAmount() End.")
	return nil
}

func (nft NftDb) BuyResultWithWAmount(nftTx *contracts.NftTx) error {
	from := strings.ToLower(nftTx.From)
	to := strings.ToLower(nftTx.To)
	contractAddr := strings.ToLower(nftTx.Contract)
	tokenId := strings.ToLower(nftTx.TokenId)
	nftaddr := strings.ToLower(nftTx.NftAddr)
	txhash := strings.ToLower(nftTx.TxHash)
	txtime, _ := strconv.ParseInt(nftTx.Ts, 10, 64)
	if contractAddr != ExchangeOwer && nftaddr[:3] != "0x8" {
		return nil
	}

	var price string
	if len(nftTx.Price) >= 9 {
		price = nftTx.Price[:len(nftTx.Price)-9]
	} else {
		price = "0"
	}
	//fmt.Println("BuyResultWithWAmount() price = ", nftTx.Price)
	//if IsUint64DataValid(nftTx.Price) != true {
	//	fmt.Println("BuyResultWithWAmount() price err")
	//	return ErrPrice
	//}
	fmt.Println(time.Now().String()[:25], "BuyResultWithWAmount() Begin", " from=", from, " to=", to, " price=", nftTx.Price,
		" contractAddr=", contractAddr, " tokenId=", tokenId, " txhash=", txhash, " nftaddr=", nftaddr)
	trans := Trans{}
	err := nft.db.Select("id").Where("nftaddr = ? AND txhash = ? AND (selltype = ? or selltype = ? or selltype = ?)",
		nftaddr, txhash, SellTypeFixPrice.String(), SellTypeBidPrice.String(), SellTypeHighestBid.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultWithWAmount() trans exist.")
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultWithWAmount() err =", err.Error)
		return err.Error
	}
	fmt.Println("BuyResultWithWAmount()-------------------")
	if from != "" && to != "" /* && nftTx.Price != ""*/ {
		fmt.Println("BuyResultWithWAmount() from != Null && to != Null")
		var nftRec Nfts
		err := nft.db.Where("nftaddr = ? AND ownaddr = ?",
			nftaddr, from).First(&nftRec)
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				fmt.Println("BuyResultWithWAmount() nft not find ")
				return nil
			}
			fmt.Println("BuyResultWithWAmount() nft find err=", err.Error)
			return ErrNftNotExist
		}
		aucFlag := false
		aucSellType := ""
		var auctionRec Auction
		count, _ := strconv.Atoi(nftTx.Value)
		if count > nftRec.Count {
			fmt.Println("BuyResultWithWAmount() nft.count < amount.")
			return ErrNftAmount
		}
		err = nft.db.Where("tokenid = ? AND ownaddr =?",
			nftRec.Tokenid, from).First(&auctionRec)
		if err.Error != nil {
			if err.Error != gorm.ErrRecordNotFound {
				fmt.Println("BuyResultWithWAmount() auction not find err=", err.Error)
				return err.Error
			}
		} else {
			aucSellType = auctionRec.Selltype
			aucFlag = true
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			trans := Trans{}
			trans.Contract = contractAddr
			trans.Createaddr = nftRec.Createaddr
			trans.Url = nftRec.Url
			trans.Count = count
			trans.Fromaddr = from
			trans.Toaddr = to
			trans.Tokenid = nftRec.Tokenid
			trans.Nftaddr = nftaddr
			if aucFlag {
				trans.Auctionid = auctionRec.ID
				trans.Nftid = auctionRec.Nftid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Selltype = auctionRec.Selltype
			} else {
				trans.Selltype = SellTypeFixPrice.String()
			}
			if !nftTx.Status {
				trans.Selltype = SellTypeError.String()
			}
			trans.Txhash = txhash
			trans.Name = nftRec.Name
			trans.Meta = nftRec.Meta
			trans.Desc = nftRec.Desc
			trans.Price, _ = strconv.ParseUint(price, 10, 64)
			fmt.Println("BuyResultWithWAmount() trans.Price=", trans.Price)
			trans.Transtime, _ = strconv.ParseInt(nftTx.Ts, 10, 64)
			if contractAddr == ExchangeOwer {
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResultWithWAmount() create trans record err=", err.Error)
					return err.Error
				}
			}
			if nftTx.Status && contractAddr == ExchangeOwer {
				var collectRec Collects
				err = nft.db.Where("createaddr = ? AND  name=?",
					nftRec.Collectcreator, nftRec.Collections).First(&collectRec)
				if err.Error == nil {
					transCnt := collectRec.Transcnt + 1
					transAmt := collectRec.Transamt + trans.Price
					collectRec = Collects{}
					collectRec.Transcnt = transCnt
					collectRec.Transamt = transAmt
					err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
						nftRec.Collections, nftRec.Collectcreator).Updates(&collectRec)
					if err.Error != nil {
						fmt.Println("BuyResultWithWAmount() update collectRec err=", err.Error)
						return err.Error
					}
				}
			}
			tonftRec := Nfts{}
			tonftRec.Selltype = SellTypeNotSale.String()
			tonftRec.Transtime = txtime
			if nftTx.Status {
				tonftRec.Ownaddr = to
				if contractAddr == ExchangeOwer {
					tonftRec.Transprice = trans.Price
					tonftRec.Transamt = nftRec.Transamt + trans.Price
					tonftRec.Transcnt = nftRec.Transcnt + 1
				}
			}
			err = tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Updates(&tonftRec)
			if err.Error != nil {
				fmt.Println("BuyResultWithWAmount() update record err=", err.Error)
				return err.Error
			}
			if aucFlag {
				if aucSellType == SellTypeHighestBid.String() {
					if nftTx.Status {
						bidRec := Bidding{}
						err = nft.db.Model(&Bidding{}).Where("bidaddr = ? and Auctionid = ?", to, auctionRec.ID).First(&bidRec)
						if err.Error != nil {
							fmt.Println("BuyResultWithWAmount() get bid record err=", err.Error)
							return err.Error
						}
						if bidRec.VoteStage != "" {
							snftP := SnftPhase{}
							err = nft.db.Model(&SnftPhase{}).Where("tokenid = ?", bidRec.VoteStage).First(&snftP)
							if err.Error != nil {
								fmt.Println("BuyResultWithWAmount() get bidRec.VoteStage record err=", err.Error)
								return err.Error
							}
							vote := snftP.Vote
							snftP = SnftPhase{}
							snftP.Vote = vote + 1
							err = tx.Model(&SnftPhase{}).Where("tokenid = ?", bidRec.VoteStage).Updates(&snftP)
							if err.Error != nil {
								fmt.Println("BuyResultWithWAmount() get bidRec.VoteStage record err=", err.Error)
								return err.Error
							}
						}
					}
				} else {
					if nftTx.Status {
						fmt.Println("BuyResultWithWAmount() auctionRec.VoteStage=", auctionRec.VoteStage)
						if auctionRec.VoteStage != "" {
							snftP := SnftPhase{}
							err = nft.db.Model(&SnftPhase{}).Where("tokenid = ?", auctionRec.VoteStage).First(&snftP)
							if err.Error != nil {
								fmt.Println("BuyResultWithWAmount() get bidRec.VoteStage record err=", err.Error)
								return err.Error
							}
							vote := snftP.Vote
							fmt.Println("BuyResultWithWAmount() snftP.Vote=", snftP.Vote)
							snftP = SnftPhase{}
							snftP.Vote = vote + 1
							err = tx.Model(&SnftPhase{}).Where("tokenid = ?", auctionRec.VoteStage).Updates(&snftP)
							if err.Error != nil {
								fmt.Println("BuyResultWithWAmount() get bidRec.VoteStage record err=", err.Error)
								return err.Error
							}
						}
					}
				}
				err = tx.Model(&Auction{}).Where("id = ?", auctionRec.ID).Delete(&Auction{})
				if err.Error != nil {
					fmt.Println("BuyResultWithWAmount() delete auction record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Bidding{}).Where("Auctionid = ?", auctionRec.ID).Delete(&Bidding{})
				if err.Error != nil {
					fmt.Println("BuyResultWithWAmount() delete bid record err=", err.Error)
					return err.Error
				}
			}
			fmt.Println("BuyResultWithWAmount() from != Null && to != Null --> price != Null OK")
			return nil
		})
	}
	fmt.Println("BuyResultWithWAmount() End.")
	return nil
}

func (nft NftDb) BuyResultAsset(from, to, contractAddr, tokenId, amount, price, royalty, txhash, txtime string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	contractAddr = strings.ToLower(contractAddr)

	fmt.Println("BuyResultAsset() price = ", price)
	if IsUint64DataValid(price) != true {
		fmt.Println("BuyResultAsset() price err")
		return ErrPrice
	}
	fmt.Println(time.Now().String()[:25], "BuyResultWithAmount() Begin", "from=", from, "to=", to, "price=", price,
		"contractAddr=", contractAddr, "tokenId=", tokenId, "txhash=", txhash,
		"royalty=", royalty /*, "sig=", sig, "trade_sig=", trade_sig*/)
	trans := Trans{}
	err := nft.db.Select("id").Where("contract = ? AND tokenid = ? AND txhash = ? AND selltype = ?",
		contractAddr, tokenId, txhash, SellTypeAsset.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultAsset() trans exist.")
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultAsset() err =", err.Error)
		return err.Error
	}
	fmt.Println("BuyResultAsset()-------------------")
	if from != "" && to != "" && price == "" {
		fmt.Println("BuyResultAsset() from != Null && to != Null")
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			contractAddr, tokenId, from).First(&nftRec)
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				fmt.Println("BuyResultAsset() nft not find ")
				return nil
			}
			fmt.Println("BuyResultAsset() nft find err=", err.Error)
			return ErrNftNotExist
		}
		fmt.Println("BuyResultAsset() price == null")
		return nft.db.Transaction(func(tx *gorm.DB) error {
			var auctionRec Auction
			err = tx.Where("contract = ? AND tokenid = ? AND ownaddr =?",
				contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
			if err.Error != nil {
				fmt.Println("BuyResultAsset() auction not find err=", err.Error)
				return err.Error
			}
			trans := Trans{}
			trans.Auctionid = auctionRec.ID
			trans.Contract = auctionRec.Contract
			trans.Createaddr = nftRec.Createaddr
			trans.Fromaddr = from
			trans.Toaddr = to
			trans.Tokenid = auctionRec.Tokenid
			trans.Nftid = auctionRec.Nftid
			trans.Paychan = auctionRec.Paychan
			trans.Currency = auctionRec.Currency
			trans.Price = 0
			//trans.Transtime = time.Now().Unix()
			trans.Transtime, _ = strconv.ParseInt(txtime, 10, 64)
			trans.Selltype = SellTypeAsset.String()
			err := tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("BuyResultAsset() create trans record err=", err.Error)
				return err.Error
			}
			fmt.Println("BuyResultAsset() price == null OK")
			return nil
		})
	}
	fmt.Println("BuyResultAsset() End.")
	return nil
}

func (nft NftDb) BuyResultRoyalty(from, to, contractAddr, tokenId, price, royalty, txhash, txtime string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	contractAddr = strings.ToLower(contractAddr)
	transTime, _ := strconv.ParseInt(txtime, 10, 64)
	fmt.Println(time.Now().String()[:25], "BuyResultRoyalty() Begin", "from=", from, "to=", to, "price=", price,
		"contractAddr=", contractAddr, "tokenId=", tokenId, "txhash=", txhash,
		"royalty=", royalty /*, "sig=", sig, "trade_sig=", trade_sig*/)
	trans := Trans{}
	err := nft.db.Select("id").Where("contract = ? AND tokenid = ? AND txhash = ? AND selltype = ?",
		contractAddr, tokenId, txhash, SellTypeMintNft.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultRoyalty() err =", ErrTransExist)
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultRoyalty() err =", err.Error)
		return err.Error
	}
	if royalty != "" {
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ? AND createaddr = ?",
			contractAddr, tokenId, to).First(&nftRec)
		if err.Error != nil {
			if err.Error != gorm.ErrRecordNotFound {
				fmt.Println("BuyResultRoyalty() royalty err =", ErrNftNotExist)
				return ErrNftNotExist
			}
			return nil
		}
		trans := Trans{}
		trans.Contract = contractAddr
		trans.Fromaddr = ""
		trans.Toaddr = to
		trans.Tokenid = tokenId
		//trans.Price, _ = strconv.ParseUint(price, 10, 64)
		//trans.Transtime = time.Now().Unix()
		trans.Transtime = transTime
		trans.Selltype = SellTypeMintNft.String()
		trans.Name = nftRec.Name
		trans.Meta = nftRec.Meta
		trans.Desc = nftRec.Desc
		trans.Txhash = txhash
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("BuyResultRoyalty() royalty create trans err=", err.Error)
				return err.Error
			}
			nftrecord := Nfts{}
			nftrecord.Royalty, _ = strconv.Atoi(royalty)
			//nftrecord.Royalty = nftrecord.Royalty / 100
			nftrecord.Mintstate = Minted.String()
			err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =? AND createaddr = ?",
				contractAddr, tokenId, to).Updates(&nftrecord)
			if err.Error != nil {
				fmt.Println("BuyResultRoyalty() royalty update nfts record err=", err.Error)
				return err.Error
			}
			fmt.Println("BuyResultRoyalty() royalty!=Null Ok")
			return nil
		})
	}
	return nil
}

type NftMintInfo struct {
	md5        string
	name       string
	desc       string
	meta       string
	source_url string
	//nft_contract_addr string "0xfffffffffffffffffff"
	nft_token_id   string
	categories     string
	collections    string
	Collectcreator string
	asset_sample   string
	hide           string
	count          string
}

func GetNftMintInfo(creatorAddr, toAddr, blockNum, creatorNonce string) (*NftMintInfo, error) {
	return &NftMintInfo{}, nil
}

type NftSysMintInfo struct {
	Createaddr     string `json:"user_addr" gorm:"type:char(42) ;comment:'create nft address'"`
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) ;comment:'nft owner address'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft name'"`
	Desc           string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft description'"`
	Meta           string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta information'"`
	Url            string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nft raw data hold address'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) ;comment:'contract address'"`
	Count          int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) ;comment:'chain of wormholes uniquely identifies the nft flag'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft classification'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) ;comment:'collection creator address'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'NFT collection name'"`
	Image          string `json:"asset_sample" gorm:"type:longtext ;comment:'thumbnail binary data'"`
	Royalty        string `json:"royalty"`
	BlockNumber    string `json:"block_number"`
}

type SnftInfo struct {
	CreatorAddr string  `json:"creator_addr"`
	Ownaddr     string  `json:"ownaddr"`
	Contract    string  `json:"nft_contract_addr"`
	Nftaddr     string  `json:"nft_address"`
	Name        string  `json:"name"`
	Desc        string  `json:"desc"`
	Meta        string  `json:"meta"`
	Category    string  `json:"category"`
	Royalty     float64 `json:"royalty"`
	//Royalty              string `json:"royalty"`
	SourceUrl string `json:"source_url"`
	Md5       string `json:"md5"`
	//Collections          string  `json:"collections"`
	CollectionsName      string `json:"collections_name"`
	CollectionsCreator   string `json:"collections_creator"`
	CollectionsExchanger string `json:"collections_exchanger"`
	CollectionsCategory  string `json:"collections_category"`
	CollectionsImgUrl    string `json:"collections_img_url"`
	CollectionsDesc      string `json:"collections_desc"`
}

func GetNftSysMintInfo(blockNum string) ([]NftSysMintInfo, error) {
	nft := make([]NftSysMintInfo, 0, 20)
	return nft, nil
}

func (nft NftDb) BuyResultWRoyalty(mintTx *contracts.NftTx) error {
	from := strings.ToLower(mintTx.From)
	to := strings.ToLower(mintTx.To)
	contractAddr := strings.ToLower(mintTx.Contract)
	transTime, _ := strconv.ParseInt(mintTx.Ts, 10, 64)
	tokenId := strings.ToLower(mintTx.TokenId)
	txhash := strings.ToLower(mintTx.TxHash)
	nftaddr := strings.ToLower(mintTx.NftAddr)

	fmt.Println(time.Now().String()[:25], "BuyResultWRoyalty() Begin", "from=", from, "to=", to, "price=", mintTx.Price,
		"contractAddr=", contractAddr, "tokenId=", tokenId, "txhash=", mintTx.TxHash,
		"royalty=", mintTx.Ratio /*, "sig=", sig, "trade_sig=", trade_sig*/)
	trans := Trans{}
	err := nft.db.Select("id").Where("contract = ? AND tokenid = ? AND txhash = ? AND selltype = ?",
		contractAddr, tokenId, txhash, SellTypeMintNft.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultWRoyalty() err =", ErrTransExist)
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultWRoyalty() err =", err.Error)
		return err.Error
	} else {
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ? AND ownaddr = ? AND mintstate != ?",
			mintTx.Contract, mintTx.TokenId, to, Minted.String()).First(&nftRec)
		if err.Error != nil {
			if err.Error != gorm.ErrRecordNotFound {
				fmt.Println("BuyResultWRoyalty() database err =", err.Error)
				return err.Error
			}
		} else {
			nfttab := Nfts{}
			nfttab.Nftaddr = nftaddr
			nfttab.Mintstate = Minted.String()
			trans := Trans{}
			trans.Contract = contractAddr
			trans.Fromaddr = ""
			trans.Toaddr = to
			trans.Tokenid = tokenId
			trans.Nftaddr = nftaddr
			trans.Transtime = transTime
			if !mintTx.Status {
				return nil
			}
			trans.Selltype = SellTypeMintNft.String()
			trans.Txhash = txhash
			trans.Count = nftRec.Count
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResultWRoyalty() royalty create trans err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
					contractAddr, tokenId, to).Updates(&nfttab)
				if err.Error != nil {
					fmt.Println("BuyResultWRoyalty() royalty create nfts record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResultWRoyalty() royalty Ok", " to=", to, " nftaddr=", nftaddr)
				return nil
			})
		}
	}
	return nil
}

func (nft NftDb) BuyResultWTransfer(mintTx *contracts.NftTx) error {
	to := strings.ToLower(mintTx.To)
	//contractAddr := strings.ToLower(mintTx.Contract)
	transTime, _ := strconv.ParseInt(mintTx.Ts, 10, 64)
	//tokenId := strings.ToLower(mintTx.TokenId)
	txhash := strings.ToLower(mintTx.TxHash)
	nftaddr := strings.ToLower(mintTx.NftAddr)
	fmt.Println("BuyResultWTransfer() to=", to)
	fmt.Println("BuyResultWTransfer() nftaddr=", nftaddr)
	fmt.Println("BuyResultWTransfer() txhash=", txhash)
	if nftaddr == "" {
		fmt.Println("BuyResultWTransfer() error nftaddr equal null.")
		return nil
	}

	fmt.Println(time.Now().String()[:25], "BuyResultWTransfer() Begin", "to=", to, "nftaddr=", nftaddr, "block=", mintTx.BlockNumber)
	/*trans := Trans{}
	err := nft.db.Select("id").Where("txhash = ? AND selltype = ?", txhash, SellTypeTransfer.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultWRoyalty() err =", ErrTransExist)
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultWTransfer() err =", err.Error)
		return err.Error
	}*/
	var nftRec Nfts
	err := nft.db.Where("nftaddr = ?", nftaddr).First(&nftRec)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("BuyResultWTransfer() database err =", err.Error)
			return err.Error
		}
	} else {
		nfttab := Nfts{}
		nfttab.Ownaddr = to
		trans := Trans{}
		trans.Contract = nftRec.Contract
		trans.Fromaddr = ZeroAddr
		trans.Toaddr = to
		trans.Tokenid = nftRec.Tokenid
		trans.Nftaddr = nftaddr
		trans.Transtime = transTime
		trans.Txhash = txhash
		trans.Count = nftRec.Count
		trans.Selltype = SellTypeTransfer.String()
		trans.Txhash = txhash
		return nft.db.Transaction(func(tx *gorm.DB) error {
			/*err := tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("BuyResultWTransfer() royalty create trans err=", err.Error)
				return err.Error
			}*/
			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Updates(&nfttab)
			if err.Error != nil {
				fmt.Println("BuyResultWTransfer() create nfts record err=", err.Error)
				return err.Error
			}
			fmt.Println("BuyResultWTransfer() Ok blocknumber=", mintTx.BlockNumber, " nftaddr=", mintTx.NftAddr, " to=", nfttab.Ownaddr)
			return nil
		})
	}
	return nil
}

func (nft NftDb) BuyResultExchange(exchangeTx *contracts.NftTx) error {
	nftaddr := strings.ToLower(exchangeTx.NftAddr)
	nftaddress := ""
	recCount := int64(0)
	switch len(nftaddr) {
	case SnftExchangeStage:
		nftaddress = nftaddr + "000"
	case SnftExchangeColletion:
		nftaddress = nftaddr + "00"
	case SnftExchangeSnft:
		nftaddress = nftaddr + "0"
	case SnftExchangeChip:
		nftaddress = nftaddr
		snft := nftaddress[:len(nftaddress)-1]
		err := nft.db.Model(Nfts{}).Where("snft = ? and ownaddr = ?", snft, ZeroAddr).Count(&recCount)
		if err.Error != nil {
			log.Println("BuyResultExchange() recCount err=", err)
			return err.Error
		}
	}
	fmt.Println("BuyResultExchange() beging. nftaddr=", nftaddr)
	nftRec := Nfts{}
	//err := nft.db.Select([]string{"collectcreator", "Collections"}).Where("nftaddr = ?", nftaddress).First(&nftRec)
	err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddress).First(&nftRec)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println("BuyResultExchange() dbase error.")
		return ErrDataBase
	}
	if err.Error == gorm.ErrRecordNotFound {
		log.Println("BuyResultExchange() snft not find error.")
		return nil
	}
	if nftRec.Ownaddr == ZeroAddr {
		log.Println("BuyResultExchange() snft already exchange.")
		return nil
	}
	collectRec := Collects{}
	err = nft.db.Where("createaddr = ? AND  name=?",
		nftRec.Collectcreator, nftRec.Collections).First(&collectRec)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println("BuyResultExchange() database err=", err.Error)
		return ErrDataBase
	}
	if err.Error == gorm.ErrRecordNotFound {
		log.Println("BuyResultExchange() snft not find error.")
		return nil
	}
	sysInfo := SysInfos{}
	err = nft.db.Model(&SysInfos{}).Last(&sysInfo)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			log.Println("UploadWNft() SysInfos err=", err)
			return ErrCollectionNotExist
		}
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		switch len(nftaddr) {
		case SnftExchangeStage:
			fmt.Println("BuyResultExchange() exchange 38 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			err := tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("ownaddr", ZeroAddr)
			if err.Error != nil {
				log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Collects{}).Where("Snftstage = ?", nftaddr).Delete(&Collects{})
			if err.Error != nil {
				log.Println("BuyResultExchange() 38 exchange deleted collect recorde err= ", err.Error)
				return err.Error
			}
			if sysInfo.Snfttotal >= 256 {
				sysInfo.Snfttotal -= 256
				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
				if err.Error != nil {
					log.Println("BuyResultExchange() 38 exchange sub SysInfos snfttotal err=", err.Error)
					return err.Error
				}
			}
			//NftCatch.SetFlushFlag()
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(SnftExchange)
		case SnftExchangeColletion:
			fmt.Println("BuyResultExchange() exchange 40 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			err := tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Update("ownaddr", ZeroAddr)
			if err.Error != nil {
				log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
				nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
			if err.Error != nil {
				log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
				return err.Error
			}
			if sysInfo.Snfttotal >= 16 {
				sysInfo.Snfttotal -= 16
				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
				if err.Error != nil {
					log.Println("BuyResultExchange() 39 exchange sub SysInfos snfttotal err=", err.Error)
					return err.Error
				}
			}
			//NftCatch.SetFlushFlag()
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(SnftExchange)
		case SnftExchangeSnft:
			fmt.Println("BuyResultExchange() exchange 41 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			err := tx.Model(&Nfts{}).Where("Snft = ?", nftaddr).Update("ownaddr", ZeroAddr)
			if err.Error != nil {
				log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
				return err.Error
			}
			if collectRec.Totalcount >= 16 {
				collectRec.Totalcount -= 16
				err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
					nftRec.Collectcreator, nftRec.Collections).Update("totalcount", collectRec.Totalcount)
				if err.Error != nil {
					log.Println("BuyResultExchange() 40 exchange update collect recorde err= ", err.Error)
					return err.Error
				}
				/*if collectRec.Totalcount == 0 {
					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
						nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
					if err.Error != nil {
						log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
						return err.Error
					}
				}*/
			}
			if sysInfo.Snfttotal > 0 {
				sysInfo.Snfttotal -= 1
				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
				if err.Error != nil {
					log.Println("BuyResultExchange() 40 exchange sub  SysInfos snfttotal err=", err.Error)
					return err.Error
				}
			}
			//NftCatch.SetFlushFlag()
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(SnftExchange)
		case SnftExchangeChip:
			fmt.Println("BuyResultExchange() exchange 42 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress).Update("ownaddr", ZeroAddr)
			if err.Error != nil {
				log.Println("BuyResultExchange() err=", err.Error)
				return err.Error
			}
			if collectRec.Totalcount >= 1 {
				collectRec.Totalcount -= 1
				err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
					nftRec.Collectcreator, nftRec.Collections).Update("totalcount", collectRec.Totalcount)
				if err.Error != nil {
					fmt.Println("BuyResultExchange() add collectins totalcount err= ", err.Error)
					return err.Error
				}
				/*if collectRec.Totalcount == 0 {
					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
						nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
					if err.Error != nil {
						log.Println("BuyResultExchange() deleted collect recorde err= ", err.Error)
						return err.Error
					}
				}*/
				if recCount+1 == 16 {
					if sysInfo.Snfttotal > 0 {
						sysInfo.Snfttotal -= 1
						err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
						if err.Error != nil {
							log.Println("BuyResultExchange() add  SysInfos snfttotal err=", err.Error)
							return err.Error
						}
					}
				}
				//NftCatch.SetFlushFlag()
				GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
				GetRedisCatch().SetDirtyFlag(SnftExchange)
			}
		}
		return nil
	})
	return nil
}

func (nft NftDb) BuyResultWPledge(Tx *contracts.NftTx) error {
	nftaddr := strings.ToLower(Tx.NftAddr)
	fmt.Println("BuyResultWPledge() nftaddr=", nftaddr, " transType=", Tx.TransType)
	fmt.Println("BuyResultWPledge() transType=", Tx.TransType)
	switch len(nftaddr) {
	case SnftExchangeSnft:
		if Tx.TransType == contracts.WormHolesPledge {
			err := nft.db.Model(&Nfts{}).Where("Snft = ?", nftaddr).Update("Pledgestate", Pledge.String())
			if err.Error != nil {
				log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
				return err.Error
			}
		} else {
			err := nft.db.Model(&Nfts{}).Where("Snft = ?", nftaddr).Update("Pledgestate", NoPledge.String())
			if err.Error != nil {
				log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
				return err.Error
			}
		}
	case snftCollectionOffset:
		if Tx.TransType == contracts.WormHolesPledge {
			err := nft.db.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Update("Pledgestate", Pledge.String())
			if err.Error != nil {
				log.Println("BuyResultWPledge() snftCollectionOffset err=", err.Error)
				return err.Error
			}
		} else {
			err := nft.db.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Update("Pledgestate", NoPledge.String())
			if err.Error != nil {
				log.Println("BuyResultWPledge() snftCollectionOffset err=", err.Error)
				return err.Error
			}
		}
	case SnftExchangeStage:
		if Tx.TransType == contracts.WormHolesPledge {
			err := nft.db.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("Pledgestate", Pledge.String())
			if err.Error != nil {
				log.Println("BuyResultWPledge() SnftExchangeStage err=", err.Error)
				return err.Error
			}
		} else {
			err := nft.db.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("Pledgestate", NoPledge.String())
			if err.Error != nil {
				log.Println("BuyResultWPledge() SnftExchangeStage err=", err.Error)
				return err.Error
			}
		}
	}
	return nil
}
