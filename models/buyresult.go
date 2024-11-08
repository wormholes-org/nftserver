package models

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"math/big"
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

func SnftMerge(nftaddr, toAddr string, accountInfo *contracts.Account, tx *gorm.DB, db *gorm.DB) error {
	fmt.Println("SnftMerge start ...")
	switch len(nftaddr) {
	case SnftExchangeStage:
		//nfttab := Nfts{}
		//nfttab.Ownaddr = toAddr
		//fmt.Println("SnftMerge() blocknumber=", NftTx.BlockNumber, " nftaddr=", NftTx.NftAddr, " to=", nfttab.Ownaddr)
		//err := tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Updates(&nfttab)
		//if err.Error != nil {
		//	fmt.Println("SnftMerge() create nfts record err=", err.Error)
		//	return err.Error
		//}
		//nfttab := Nfts{}
		//nfttab.Ownaddr = toAddr
		//nfttab.Mergelevel = accountInfo.MergeLevel
		//nfttab.Mergenumber = accountInfo.MergeNumber
		//err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").Updates(&nfttab)
		//if err.Error != nil {
		//	log.Println("SnftMerge() update nfts record err=", err.Error)
		//	return err.Error
		//}
		return nil
	case SnftExchangeColletion:
		nftaddr = nftaddr + "00"
		if accountInfo.MergeLevel == 3 {
			mnft := Nfts{}
			err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").First(&mnft)
			if err.Error != nil {
				log.Println("SnftMerge() update nfts record err=", err.Error)
				return err.Error
			}
			//if mnft.Mergelevel < 3 {
			if true {
				nfttab := Nfts{}
				nfttab.Ownaddr = toAddr
				nfttab.Mergelevel = accountInfo.MergeLevel
				nfttab.Mergenumber = accountInfo.MergeNumber
				err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").Updates(&nfttab)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				nfttab = Nfts{}
				nfttab.Mergelevel = accountInfo.MergeLevel
				err = tx.Model(&Nfts{}).Where("snftstage = ? and (mergetype = 1 or mergetype = 2)", nftaddr[:len(nftaddr)-3]+"m").Updates(&nfttab)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				nfttab = Nfts{}
				nfttab.Nftaddr = nftaddr[:len(nftaddr)-3] + "mmm"
				nerr := ClearAuction(db, &nfttab)
				if nerr != nil {
					log.Println("SnftMerge() ClearAuction err=", err.Error)
					return err.Error
				}
			}

		}
		return nil
	case SnftExchangeSnft:
		if nftaddr[:3] == "0x8" {
			nftaddr = nftaddr + "0"
			switch accountInfo.MergeLevel {
			case 2:
				mnft := Nfts{}
				err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-2]+"mm").First(&mnft)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				//if mnft.Mergelevel < 2 {
				if true {
					nfttab := Nfts{}
					nfttab.Ownaddr = toAddr
					nfttab.Mergelevel = accountInfo.MergeLevel
					nfttab.Mergenumber = accountInfo.MergeNumber
					err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-2]+"mm").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("Snftcollection = ? and mergetype = 1", nftaddr[:len(nftaddr)-2]+"m").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Nftaddr = nftaddr[:len(nftaddr)-2] + "mm"
					nerr := ClearAuction(db, &nfttab)
					if nerr != nil {
						log.Println("SnftMerge() ClearAuction err=", err.Error)
						return err.Error
					}
				}
			case 3:
				mnft := Nfts{}
				err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").First(&mnft)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				//if mnft.Mergelevel < 3 {
				if true {
					nfttab := Nfts{}
					nfttab.Ownaddr = toAddr
					nfttab.Mergelevel = accountInfo.MergeLevel
					nfttab.Mergenumber = accountInfo.MergeNumber
					err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("snftstage = ? and (mergetype = 1 or mergetype = 2)", nftaddr[:len(nftaddr)-3]+"m").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Nftaddr = nftaddr[:len(nftaddr)-3] + "mmm"
					nerr := ClearAuction(db, &nfttab)
					if nerr != nil {
						log.Println("SnftMerge() ClearAuction err=", err.Error)
						return err.Error
					}
				}
			}
		}
		return nil
	case SnftExchangeChip:
		if nftaddr[:3] == "0x8" && accountInfo.MergeLevel != 0 {
			switch accountInfo.MergeLevel {
			case 1:
				mnft := Nfts{}
				err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-1]+"m").First(&mnft)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				if mnft.Mergelevel < 1 {
					err = tx.Model(&Users{}).Where("useraddr = ?", toAddr).Update("Rewards", gorm.Expr("Rewards + ?", 1))
					if err.Error != nil {
						log.Println("SnftMerge() update users record err=", err.Error)
						return err.Error
					}
				}

				//if mnft.Mergelevel < 1 {
				if true {
					nfttab := Nfts{}
					nfttab.Ownaddr = toAddr
					nfttab.Mergelevel = accountInfo.MergeLevel
					nfttab.Mergenumber = accountInfo.MergeNumber
					err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-1]+"m").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("snft = ?", nftaddr[:len(nftaddr)-1]).Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Nftaddr = nftaddr[:len(nftaddr)-1] + "m"
					nerr := ClearAuction(db, &nfttab)
					if nerr != nil {
						log.Println("SnftMerge() ClearAuction err=", err.Error)
						return err.Error
					}
				}
			case 2:
				mnft := Nfts{}
				err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-2]+"mm").First(&mnft)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				//if mnft.Mergelevel < 2 {
				if true {
					nfttab := Nfts{}
					nfttab.Ownaddr = toAddr
					nfttab.Mergelevel = accountInfo.MergeLevel
					nfttab.Mergenumber = accountInfo.MergeNumber
					err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-2]+"mm").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr[:len(nftaddr)-2]).Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("Snftcollection = ? and mergetype = 1", nftaddr[:len(nftaddr)-2]+"m").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Nftaddr = nftaddr[:len(nftaddr)-2] + "mm"
					nerr := ClearAuction(db, &nfttab)
					if nerr != nil {
						log.Println("SnftMerge() ClearAuction err=", err.Error)
						return err.Error
					}
				}
			case 3:
				mnft := Nfts{}
				err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").First(&mnft)
				if err.Error != nil {
					log.Println("SnftMerge() update nfts record err=", err.Error)
					return err.Error
				}
				//if mnft.Mergelevel < 3 {
				if true {
					nfttab := Nfts{}
					nfttab.Ownaddr = toAddr
					nfttab.Mergelevel = accountInfo.MergeLevel
					nfttab.Mergenumber = accountInfo.MergeNumber
					err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("snftstage = ? and (mergetype = 1 or mergetype = 2)", nftaddr[:len(nftaddr)-3]+"m").Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Mergelevel = accountInfo.MergeLevel
					err = tx.Model(&Nfts{}).Where("snftstage = ?", nftaddr[:len(nftaddr)-3]).Updates(&nfttab)
					if err.Error != nil {
						log.Println("SnftMerge() update nfts record err=", err.Error)
						return err.Error
					}
					nfttab = Nfts{}
					nfttab.Nftaddr = nftaddr[:len(nftaddr)-3] + "mmm"
					nerr := ClearAuction(db, &nfttab)
					if nerr != nil {
						log.Println("SnftMerge() ClearAuction err=", err.Error)
						return err.Error
					}
				}
			}
		}
		return nil
	}

	return nil
}

func (nft NftDb) BuyResultWithWAmount(nftTx *contracts.NftTx) error {
	from := strings.ToLower(nftTx.From)
	to := strings.ToLower(nftTx.To)
	contractAddr := strings.ToLower(nftTx.Contract)
	tokenId := strings.ToLower(nftTx.TokenId)
	nftaddr := strings.ToLower(nftTx.NftAddr)
	OldNftaddr := nftaddr
	txhash := strings.ToLower(nftTx.TxHash)
	txtime, _ := strconv.ParseInt(nftTx.Ts, 10, 64)
	//if contractAddr != ExchangeOwer && nftaddr[:3] != "0x8" {
	//	return nil
	//}

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
	switch len(nftaddr) {
	case SnftExchangeStage:
		nftaddr = nftaddr + "mmm"
	case SnftExchangeColletion:
		nftaddr = nftaddr + "mm"
	case SnftExchangeSnft:
		nftaddr = nftaddr + "m"
	}
	fmt.Println(time.Now().String()[:25], "BuyResultWithWAmount() Begin", " from=", from, " to=", to, " price=", nftTx.Price,
		" contractAddr=", contractAddr, " tokenId=", tokenId, " txhash=", txhash, " nftaddr=", nftaddr)
	trans := Trans{}
	err := nft.db.Select("id").Where("nftaddr = ? AND txhash = ? AND (selltype = ? or selltype = ? or selltype = ? or selltype = ?)",
		nftaddr, txhash, SellTypeFixPrice.String(), SellTypeBidPrice.String(), SellTypeHighestBid.String(), SellTypeForceBuy.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("BuyResultWithWAmount() trans exist.")
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("BuyResultWithWAmount() err =", err.Error)
		return err.Error
	}
	fmt.Println("BuyResultWithWAmount()-------------------")
	var count int
	if from != "" && to != "" /* && nftTx.Price != ""*/ {
		fmt.Println("BuyResultWithWAmount() from != Null && to != Null")
		var nftRec Nfts
		if nftTx.TransType != contracts.WormHolesExForceBuyingAuthTransfer {
			err := nft.db.Where("nftaddr = ? AND ownaddr = ?", nftaddr, from).First(&nftRec)
			if err.Error != nil {
				if err.Error == gorm.ErrRecordNotFound {
					log.Println("BuyResultWithWAmount() nft not find ")
					return nil
				}
				fmt.Println("BuyResultWithWAmount() nft find err=", err.Error)
				return ErrNftNotExist
			}
			count, _ := strconv.Atoi(nftTx.Value)
			if count > nftRec.Count {
				log.Println("BuyResultWithWAmount() nft.count < amount.")
				return ErrNftAmount
			}
		} else {
			err := nft.db.Where("nftaddr = ?", nftaddr).First(&nftRec)
			if err.Error != nil {
				if err.Error == gorm.ErrRecordNotFound {
					log.Println("BuyResultWithWAmount() nft not find ")
					return nil
				}
				fmt.Println("BuyResultWithWAmount() nft find err=", err.Error)
				return ErrNftNotExist
			}
			count, _ := strconv.Atoi(nftTx.Value)
			if count > nftRec.Count {
				log.Println("BuyResultWithWAmount() nft.count < amount.")
				return ErrNftAmount
			}
		}

		aucFlag := false
		//aucSellType := ""
		var auctionRec Auction
		if nftTx.TransType != contracts.WormHolesExForceBuyingAuthTransfer {
			err = nft.db.Where("tokenid = ? AND ownaddr =?",
				nftRec.Tokenid, from).First(&auctionRec)
			if err.Error != nil {
				if err.Error != gorm.ErrRecordNotFound {
					log.Println("BuyResultWithWAmount() auction not find err=", err.Error)
					return err.Error
				}
			} else {
				//aucSellType = auctionRec.Selltype
				aucFlag = true
			}
		}

		//sysNft := Sysnfts{}
		//if nftaddr[:3] == "0x8" {
		//	if nftRec.Mergetype != 0 {
		//		if nftRec.Mergelevel != 3 && nftRec.Mergelevel != 2 {
		//			err = nft.db.Where("snft = ?", nftRec.Snft[:len(nftRec.Snft)-1]).First(&sysNft)
		//			if err.Error != nil {
		//				log.Println("BuyResultWithWAmount() database err=", err.Error)
		//				return ErrNftNotExist
		//			}
		//		}
		//	} else {
		//		err = nft.db.Where("snft = ?", nftRec.Snft).First(&sysNft)
		//		if err.Error != nil {
		//			log.Println("BuyResultWithWAmount() database err=", err.Error)
		//			return ErrNftNotExist
		//		}
		//	}
		//
		//}
		var accountInfo *contracts.Account
		if nftaddr[:3] == "0x8" {
			var gerr error
			if nftTx.TransType == contracts.WormHolesExForceBuyingAuthTransfer {
				accountInfo, gerr = GetMergeLevel(OldNftaddr+"0", nftTx.BlockNumber)
			} else {
				accountInfo, gerr = GetMergeLevel(OldNftaddr, nftTx.BlockNumber)
			}

			if gerr != nil {
				log.Println("BuyResultWithWAmount() ger mergelevel error.")
				return ErrBlockchain
			}
		}
		sysInfoRec := SysInfos{}
		err = nft.db.Last(&sysInfoRec)
		if err.Error != nil {
			log.Println("BuyResultWithWAmount() Last(&sysInfoRec) err=", err.Error)
			return ErrDataBase
		}
		collectionRec := Collects{}
		err := nft.db.Select("id").Where("Createaddr = ? and Name = ?", nftRec.Collectcreator, nftRec.Collections).First(&collectionRec)
		if err.Error != nil {
			log.Println("BuyResultWithWAmount() Collects err=", err.Error)
			return ErrDataBase
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
			trans.Collectionid = collectionRec.ID
			if nftaddr[:3] == "0x8" {
				trans.Nfttype = "snft"
			} else {
				trans.Nfttype = "nft"
			}
			if aucFlag {
				trans.Auctionid = auctionRec.ID
				trans.Nftid = auctionRec.Nftid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Selltype = auctionRec.Selltype
			} else {
				trans.Selltype = SellTypeFixPrice.String()
			}
			if nftTx.TransType == contracts.WormHolesExForceBuyingAuthTransfer {
				trans.Selltype = SellTypeForceBuy.String()
				trans.Nftid = nftRec.ID
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
					log.Println("BuyResultWithWAmount() create trans record err=", err.Error)
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
			nfttab := map[string]interface{}{
				"Selltype":    SellTypeNotSale.String(),
				"Transtime":   txtime,
				"Sellprice":   0,
				"Offernum":    0,
				"Maxbidprice": 0,
			}
			if nftTx.Status {
				nfttab["Ownaddr"] = to
				if contractAddr == ExchangeOwer {
					nfttab["Transprice"] = trans.Price
					nfttab["Transamt"] = nftRec.Transamt + trans.Price
					nfttab["Transcnt"] = nftRec.Transcnt + 1
				}
			}
			if nftTx.TransType == contracts.WormHolesExForceBuyingAuthTransfer {
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Updates(&nfttab)
				if err.Error != nil {
					fmt.Println("BuyResultWithWAmount() update record err=", err.Error)
					return err.Error
				}
			} else {
				err = tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Updates(&nfttab)
				if err.Error != nil {
					fmt.Println("BuyResultWithWAmount() update record err=", err.Error)
					return err.Error
				}
			}

			//if nftaddr[:3] == "0x8" {
			//	snt := Sysnfts{}
			//	snt.Transcnt = sysNft.Transcnt + 1
			//	snt.Transamt = sysNft.Transamt + trans.Price
			//	snt.Transprice = (snt.Transamt / uint64(snt.Transcnt))
			//	err = tx.Model(&Sysnfts{}).Where("id = ?", sysNft.ID).Updates(&snt)
			//	if err.Error != nil {
			//		fmt.Println("BuyResultWithWAmount() update record err=", err.Error)
			//		return err.Error
			//	}
			//}
			if aucFlag {
				/* vote cancel
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
					}*/
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
				switch auctionRec.Selltype {
				case SellTypeFixPrice.String():
					if sysInfoRec.Fixpricecnt > 1 {
						err = tx.Model(&SysInfos{}).Where("id = ?", sysInfoRec.ID).Update("Fixpricecnt", gorm.Expr("Fixpricecnt - ?", 1))
						if err.Error != nil {
							log.Println("Sell() Fixpricecnt  err= ", err.Error)
							return ErrDataBase
						}
					}
				case SellTypeHighestBid.String():
					if sysInfoRec.Highestbidcnt > 1 {
						err = tx.Model(&SysInfos{}).Where("id = ?", sysInfoRec.ID).Update("Highestbidcnt", gorm.Expr("Highestbidcnt - ?", 1))
						if err.Error != nil {
							log.Println("Sell() Highestbidcnt  err= ", err.Error)
							return ErrDataBase
						}
					}
				}
			}
			if nftaddr[:3] == "0x8" {
				if nftTx.TransType == contracts.WormHolesExForceBuyingAuthTransfer {
					OldNftaddr = OldNftaddr + "0"
				}
				nerr := SnftMerge(OldNftaddr, to, accountInfo, tx, nft.db)
				if nerr != nil {
					fmt.Println("BuyResultWithWAmount() SnftMerge err=", nerr)
					return nerr
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
	CreatorAddr          string  `json:"creator_addr"`
	Ownaddr              string  `json:"ownaddr"`
	Contract             string  `json:"nft_contract_addr"`
	Nftaddr              string  `json:"nft_address"`
	Name                 string  `json:"name"`
	Desc                 string  `json:"desc"`
	Meta                 string  `json:"meta"`
	Category             string  `json:"category"`
	Royalty              float64 `json:"royalty"`
	SourceUrl            string  `json:"source_url"`
	Md5                  string  `json:"md5"`
	CollectionsName      string  `json:"collections_name"`
	CollectionsCreator   string  `json:"collections_creator"`
	CollectionsExchanger string  `json:"collections_exchanger"`
	CollectionsCategory  string  `json:"collections_category"`
	CollectionsImgUrl    string  `json:"collections_img_url"`
	CollectionsDesc      string  `json:"collections_desc"`
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

func ClearAuction(db *gorm.DB, nftRec *Nfts) error {
	//var auctionRec Auction
	//err := tx.Where("tokenid = ? AND ownaddr =?", nftRec.Tokenid, nftRec.Ownaddr).First(&auctionRec)
	//if err.Error != nil {
	// if err.Error != gorm.ErrRecordNotFound {
	//    log.Println("ClearAuction() auction not find err=", err.Error)
	//    return err
	// } else {
	//    return nil
	// }
	//}
	//fmt.Println("ClearAuction start ...  nftrec=", nftRec.Nftaddr)
	//return nil
	addr := strings.ReplaceAll(nftRec.Nftaddr, "m", "")
	aucNfts := []Auction{}
	err := db.Model(&Auction{}).Select([]string{"id", "nftaddr"}).Where("nftaddr like ?", addr+"%").Find(&aucNfts)
	if err.Error != nil {
		log.Println("ClearAuction() find record err=", err.Error)
		return err.Error
	}
	if len(aucNfts) > 0 {
		go func() {
			for _, aucnft := range aucNfts {
				db.Transaction(func(tx *gorm.DB) error {
					nfttab := map[string]interface{}{
						"Selltype":    SellTypeNotSale.String(),
						"Sellprice":   0,
						"Offernum":    0,
						"Maxbidprice": 0,
					}
					fmt.Println("ClearAuction() update record nftaddr=", aucnft.Nftaddr)
					err := tx.Model(&Nfts{}).Where("nftaddr = ?", aucnft.Nftaddr).Updates(&nfttab)
					if err.Error != nil {
						log.Println("ClearAuction() update record err=", err.Error)
						return err.Error
					}
					err = tx.Model(&Bidding{}).Where("Auctionid = ?", aucnft.ID).Delete(&Bidding{})
					if err.Error != nil {
						log.Println("ClearAuction() delete bid record err=", err.Error)
						return err.Error
					}
					err = db.Model(&Auction{}).Where("id = ?", aucnft.ID).Delete(&Auction{})
					if err.Error != nil {
						log.Println("ClearAuction() delete auction record err=", err.Error)
						return err.Error
					}
					return nil
				})
			}
			fmt.Println("ClearAuction() clear end nftaddr=", nftRec.Nftaddr)
		}()
	}
	return nil
}

//func (nft NftDb) BuyResultWTransferOld(mintTx *contracts.NftTx) error {
//	to := strings.ToLower(mintTx.To)
//	//contractAddr := strings.ToLower(mintTx.Contract)
//	//transTime, _ := strconv.ParseInt(mintTx.Ts, 10, 64)
//	//tokenId := strings.ToLower(mintTx.TokenId)
//	txhash := strings.ToLower(mintTx.TxHash)
//	nftaddr := strings.ToLower(mintTx.NftAddr)
//	fmt.Println("BuyResultWTransfer() to=", to)
//	fmt.Println("BuyResultWTransfer() nftaddr=", nftaddr)
//	fmt.Println("BuyResultWTransfer() txhash=", txhash)
//	if nftaddr == "" {
//		fmt.Println("BuyResultWTransfer() error nftaddr equal null.")
//		return nil
//	}
//
//	fmt.Println(time.Now().String()[:25], "BuyResultWTransfer() Begin", "to=", to, "nftaddr=", nftaddr, "block=", mintTx.BlockNumber)
//	/*trans := Trans{}
//	err := nft.db.Select("id").Where("txhash = ? AND selltype = ?", txhash, SellTypeTransfer.String()).First(&trans)
//	if err.Error == nil {
//		fmt.Println("BuyResultWRoyalty() err =", ErrTransExist)
//		return nil
//	}
//	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//		fmt.Println("BuyResultWTransfer() err =", err.Error)
//		return err.Error
//	}*/
//	var accountInfo *contracts.Account
//	if nftaddr[:3] == "0x8" {
//		var gerr error
//		accountInfo, gerr = GetMergeLevel(nftaddr, mintTx.BlockNumber)
//		if gerr != nil {
//			log.Println("BuyResultWTransfer() ger mergelevel error.")
//			return ErrBlockchain
//		}
//	}
//	switch len(nftaddr) {
//	case SnftExchangeStage:
//		if nftaddr[:3] != "0x8" {
//			return nil
//		}
//		var nftRec Nfts
//		err := nft.db.Where("nftaddr = ?", nftaddr+"mmm").First(&nftRec)
//		if err.Error != nil {
//			if err.Error != gorm.ErrRecordNotFound {
//				fmt.Println("BuyResultWTransfer() database err =", err.Error)
//				return err.Error
//			}
//		} else {
//			nfttab := Nfts{}
//			nfttab.Ownaddr = to
//			fmt.Println("BuyResultWTransfer() blocknumber=", mintTx.BlockNumber, " nftaddr=", mintTx.NftAddr, " to=", nfttab.Ownaddr)
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err = ClearAuction(tx, &nftRec)
//				if err != nil {
//					log.Println("BuyResultWTransfer() ClearAuction err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Updates(&nfttab)
//				if err.Error != nil {
//					fmt.Println("BuyResultWTransfer() create nfts record err=", err.Error)
//					return err.Error
//				}
//				nfttab = Nfts{}
//				nfttab.Ownaddr = to
//				nfttab.Mergelevel = accountInfo.MergeLevel
//				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").Updates(&nfttab)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//					return err.Error
//				}
//				return nil
//			})
//		}
//	case SnftExchangeColletion:
//		if nftaddr[:3] != "0x8" {
//			return nil
//		}
//		var nftRec Nfts
//		err := nft.db.Where("nftaddr = ?", nftaddr+"mm").First(&nftRec)
//		if err.Error != nil {
//			if err.Error != gorm.ErrRecordNotFound {
//				fmt.Println("BuyResultWTransfer() database err =", err.Error)
//				return err.Error
//			}
//		} else {
//			var mnftRec Nfts
//			if accountInfo.MergeLevel == 3 {
//				err := nft.db.Where("nftaddr = ?", nftaddr+"mm").First(&mnftRec)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() database err =", err.Error)
//					return err.Error
//				}
//			}
//			nfttab := Nfts{}
//			nfttab.Ownaddr = to
//			fmt.Println("BuyResultWTransfer() blocknumber=", mintTx.BlockNumber, " nftaddr=", mintTx.NftAddr, " to=", nfttab.Ownaddr)
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err = ClearAuction(tx, &nftRec)
//				if err != nil {
//					log.Println("BuyResultWTransfer() ClearAuction err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Updates(&nfttab)
//				if err.Error != nil {
//					fmt.Println("BuyResultWTransfer() create nfts record err=", err.Error)
//					return err.Error
//				}
//				nfttab = Nfts{}
//				nfttab.Ownaddr = to
//				nfttab.Mergelevel = accountInfo.MergeLevel
//				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mm").Updates(&nfttab)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//					return err.Error
//				}
//				if accountInfo.MergeLevel == 3 {
//					nfttab = Nfts{}
//					nfttab.Ownaddr = to
//					nfttab.Mergelevel = accountInfo.MergeLevel
//					err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-1]+"mmm").Updates(&nfttab)
//					if err.Error != nil {
//						log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//						return err.Error
//					}
//					nfttab = Nfts{}
//					nfttab.Ownaddr = to
//					nfttab.Mergelevel = accountInfo.MergeLevel
//					err = tx.Model(&Nfts{}).Where("snftstage = ?", mnftRec.Snftstage).Updates(&nfttab)
//					if err.Error != nil {
//						log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//						return err.Error
//					}
//				}
//				return nil
//			})
//		}
//	case SnftExchangeSnft:
//		if nftaddr[:3] != "0x8" {
//			return nil
//		}
//		var nftRec Nfts
//		err := nft.db.Where("nftaddr = ?", nftaddr+"m").First(&nftRec)
//		if err.Error != nil {
//			if err.Error != gorm.ErrRecordNotFound {
//				fmt.Println("BuyResultWTransfer() database err =", err.Error)
//				return err.Error
//			}
//		} else {
//			var mnftRec Nfts
//			if accountInfo.MergeLevel > 1 {
//				err := nft.db.Where("nftaddr = ?", nftaddr+"m").First(&mnftRec)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() database err =", err.Error)
//					return err.Error
//				}
//			}
//			nfttab := Nfts{}
//			nfttab.Ownaddr = to
//			fmt.Println("BuyResultWTransfer() blocknumber=", mintTx.BlockNumber, " nftaddr=", mintTx.NftAddr, " to=", nfttab.Ownaddr)
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err = ClearAuction(tx, &nftRec)
//				if err != nil {
//					log.Println("BuyResultWTransfer() ClearAuction err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&Nfts{}).Where("snft = ?", nftaddr).Updates(&nfttab)
//				if err.Error != nil {
//					fmt.Println("BuyResultWTransfer() create nfts record err=", err.Error)
//					return err.Error
//				}
//				nfttab = Nfts{}
//				nfttab.Ownaddr = to
//				nfttab.Mergelevel = accountInfo.MergeLevel
//				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"m").Updates(&nfttab)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//					return err.Error
//				}
//				if accountInfo.MergeLevel > 1 {
//					switch accountInfo.MergeLevel {
//					case 2:
//						nfttab = Nfts{}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-1]+"mm").Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("Snftcollection = ?", mnftRec.Snftcollection).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//					case 3:
//						nfttab = Nfts{}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-2]+"mmm").Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("snftstage = ?", mnftRec.Snftstage).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//					}
//				}
//				return nil
//			})
//		}
//	case SnftExchangeChip:
//		var nftRec Nfts
//		err := nft.db.Where("nftaddr = ?", nftaddr).First(&nftRec)
//		if err.Error != nil {
//			if err.Error != gorm.ErrRecordNotFound {
//				fmt.Println("BuyResultWTransfer() database err =", err.Error)
//				return err.Error
//			}
//		} else {
//			var mnftRec Nfts
//			if nftaddr[:3] == "0x8" && accountInfo.MergeLevel != 0 {
//				err := nft.db.Where("nftaddr = ?", nftaddr[:len(nftaddr)-1]+"m").First(&mnftRec)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() database err =", err.Error)
//					return err.Error
//				}
//			}
//			nfttab := Nfts{}
//			nfttab.Ownaddr = to
//			fmt.Println("BuyResultWTransfer() blocknumber=", mintTx.BlockNumber, " nftaddr=", mintTx.NftAddr, " to=", nfttab.Ownaddr)
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err = ClearAuction(tx, &nftRec)
//				if err != nil {
//					log.Println("BuyResultWTransfer() ClearAuction err=", err.Error)
//					return err.Error
//				}
//				if nftaddr[:3] == "0x8" && accountInfo.MergeLevel != 0 {
//					nfttab.Mergelevel = 1
//				}
//				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Updates(&nfttab)
//				if err.Error != nil {
//					log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//					return err.Error
//				}
//				if nftaddr[:3] == "0x8" && accountInfo.MergeLevel != 0 {
//					nfttab = Nfts{}
//					if accountInfo.MergeLevel != 0 {
//						nfttab.Mergelevel = 1
//						err = tx.Model(&Nfts{}).Where("snft = ?", nftaddr[:len(nftaddr)-1]).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//					}
//					switch accountInfo.MergeLevel {
//					case 1:
//						nfttab := Nfts{}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("nftaddr = ?", mnftRec.Nftaddr).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//					case 2:
//						nfttab = Nfts{}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-2]+"mm").Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//						nfttab = Nfts{}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("Snftcollection = ?", mnftRec.Snftcollection).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//					case 3:
//						nfttab = Nfts{}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr[:len(nftaddr)-3]+"mmm").Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//						nfttab.Ownaddr = to
//						nfttab.Mergelevel = accountInfo.MergeLevel
//						err = tx.Model(&Nfts{}).Where("snftstage = ?", mnftRec.Snftstage).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
//							return err.Error
//						}
//					}
//				}
//				return nil
//			})
//		}
//	}
//	return nil
//}

func (nft NftDb) BuyResultWTransfer(mintTx *contracts.NftTx) error {
	to := strings.ToLower(mintTx.To)
	//contractAddr := strings.ToLower(mintTx.Contract)
	//transTime, _ := strconv.ParseInt(mintTx.Ts, 10, 64)
	//tokenId := strings.ToLower(mintTx.TokenId)
	txhash := strings.ToLower(mintTx.TxHash)
	nftaddr := strings.ToLower(mintTx.NftAddr)
	OldNftaddr := nftaddr
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
	var accountInfo *contracts.Account
	if nftaddr[:3] == "0x8" {
		var gerr error
		accountInfo, gerr = GetMergeLevel(OldNftaddr, mintTx.BlockNumber)
		if gerr != nil {
			log.Println("BuyResultWTransfer() ger mergelevel error.")
			return ErrBlockchain
		}
	}
	switch len(nftaddr) {
	case SnftExchangeStage:
		nftaddr = nftaddr + "mmm"
	case SnftExchangeColletion:
		nftaddr = nftaddr + "mm"
	case SnftExchangeSnft:
		nftaddr = nftaddr + "m"
	}
	var nftRec Nfts
	err := nft.db.Where("nftaddr = ?", nftaddr).First(&nftRec)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			log.Println("BuyResultWTransfer() nft not find ")
			return nil
		}
		log.Println("BuyResultWTransfer() nft find err=", err.Error)
		return ErrNftNotExist
	}
	auctRec := Auction{}
	aucFlag := false
	err = nft.db.Where("Contract = ? and tokenid = ?", nftRec.Contract, nftRec.Tokenid).First(&auctRec)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("BuyResultWTransfer() auction not find err=", err.Error)
			return err.Error
		}
	} else {
		aucFlag = true
	}
	fmt.Println("BuyResultWTransfer() blocknumber=", mintTx.BlockNumber, " nftaddr=", mintTx.NftAddr, " to=", mintTx.To)
	return nft.db.Transaction(func(tx *gorm.DB) error {
		if aucFlag {
			nfttab := map[string]interface{}{
				"ownaddr":     to,
				"Selltype":    SellTypeNotSale.String(),
				"Sellprice":   0,
				"Offernum":    0,
				"Maxbidprice": 0,
			}
			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Updates(&nfttab)
			if err.Error != nil {
				fmt.Println("BuyResultWithWAmount() update record err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Auction{}).Where("id = ?", auctRec.ID).Delete(&Auction{})
			if err.Error != nil {
				fmt.Println("BuyResultWithWAmount() delete auction record err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Bidding{}).Where("Auctionid = ?", auctRec.ID).Delete(&Bidding{})
			if err.Error != nil {
				fmt.Println("BuyResultWithWAmount() delete bid record err=", err.Error)
				return err.Error
			}
		} else {
			nfttab := Nfts{}
			nfttab.Ownaddr = to
			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Updates(&nfttab)
			if err.Error != nil {
				log.Println("BuyResultWTransfer() update nfts record err=", err.Error)
				return err.Error
			}
		}
		if nftaddr[:3] == "0x8" {
			nerr := SnftMerge(OldNftaddr, to, accountInfo, tx, nft.db)
			if nerr != nil {
				fmt.Println("BuyResultWTransfer() SnftMerge err=", nerr)
				return nerr
			}
		}
		return nil
	})
}

func GetMergeLevel(OldNftaddr string, blocknumber string) (*contracts.Account, error) {
	switch len(OldNftaddr) {
	case SnftExchangeChip:
		addr := common.HexToAddress(OldNftaddr[:41] + "0")
		nb, _ := big.NewInt(0).SetString(blocknumber, 10)
		accountInfo, err := contracts.GetAccountInfo(addr, nb)
		if err != nil {
			log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
			return nil, err
		}
		oldInfo := accountInfo
		if accountInfo.MergeLevel == 1 {
			addr := common.HexToAddress(OldNftaddr[:40] + "00")
			accountInfo, err = contracts.GetAccountInfo(addr, nb)
			if err != nil {
				log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
				return nil, err
			}
			if accountInfo.MergeLevel == 2 {
				oldInfo = accountInfo
				addr := common.HexToAddress(OldNftaddr[:39] + "000")
				accountInfo, err = contracts.GetAccountInfo(addr, nb)
				if err != nil {
					log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
					return nil, err
				}
				if accountInfo.MergeLevel == 3 {
					oldInfo = accountInfo
				}
			}
		}
		return oldInfo, nil
	case SnftExchangeSnft:
		addr := common.HexToAddress(OldNftaddr[:40] + "00")
		nb, _ := big.NewInt(0).SetString(blocknumber, 10)
		accountInfo, err := contracts.GetAccountInfo(addr, nb)
		if err != nil {
			log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
			return nil, err
		}
		oldInfo := accountInfo
		if accountInfo.MergeLevel == 2 {
			addr := common.HexToAddress(OldNftaddr[:39] + "000")
			accountInfo, err = contracts.GetAccountInfo(addr, nb)
			if err != nil {
				log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
				return nil, err
			}
			if accountInfo.MergeLevel == 3 {
				oldInfo = accountInfo
			}
		}
		return oldInfo, nil
	case SnftExchangeColletion:
		addr := common.HexToAddress(OldNftaddr[:39] + "000")
		nb, _ := big.NewInt(0).SetString(blocknumber, 10)
		accountInfo, err := contracts.GetAccountInfo(addr, nb)
		if err != nil {
			log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
			return nil, err
		}
		oldInfo := accountInfo
		if accountInfo.MergeLevel == 3 {
			oldInfo = accountInfo
		}
		return oldInfo, nil
	case SnftExchangeStage:
		addr := common.HexToAddress(OldNftaddr + "000")
		nb, _ := big.NewInt(0).SetString(blocknumber, 10)
		accountInfo, err := contracts.GetAccountInfo(addr, nb)
		if err != nil {
			log.Println("BuyResultWTransfer() GetAccountInfo err =", err, "NftAddress= ", addr)
			return nil, err
		}
		return accountInfo, nil
	default:
		return nil, errors.New("nftaddr err.")
	}
	return nil, errors.New("nftaddr err.")
}

//func (nft NftDb) BuyResultExchangeOld(exchangeTx *contracts.NftTx) error {
//	nftaddr := strings.ToLower(exchangeTx.NftAddr)
//	nftaddress := ""
//	OldNftaddr := nftaddr
//	nftRec := Nfts{}
//	switch len(nftaddr) {
//	case SnftExchangeStage:
//		nftaddress = nftaddr + "000"
//		err := nft.db.Model(&Nfts{}).Where("snftstage = ?", nftaddr).First(&nftRec)
//		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() dbase error.")
//			return ErrDataBase
//		}
//		if err.Error == gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() snft not find error.")
//			return nil
//		}
//	case SnftExchangeColletion:
//		nftaddress = nftaddr + "00"
//		err := nft.db.Model(&Nfts{}).Where("snftcollection = ?", nftaddr).First(&nftRec)
//		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() dbase error.")
//			return ErrDataBase
//		}
//		if err.Error == gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() snft not find error.")
//			return nil
//		}
//	case SnftExchangeSnft:
//		nftaddress = nftaddr + "0"
//		err := nft.db.Model(&Nfts{}).Where("Snft = ?", nftaddr).First(&nftRec)
//		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() dbase error.")
//			return ErrDataBase
//		}
//		if err.Error == gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() snft not find error.")
//			return nil
//		}
//	case SnftExchangeChip:
//		nftaddress = nftaddr
//		err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).First(&nftRec)
//		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() dbase error.")
//			return ErrDataBase
//		}
//		if err.Error == gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() snft not find error.")
//			return nil
//		}
//	}
//	fmt.Println("BuyResultExchange() begin. nftaddr=", nftaddr)
//	collectRec := Collects{}
//	err := nft.db.Where("createaddr = ? AND  name=?",
//		nftRec.Collectcreator, nftRec.Collections).First(&collectRec)
//	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//		log.Println("BuyResultExchange() database err=", err.Error)
//		return ErrDataBase
//	}
//	if err.Error == gorm.ErrRecordNotFound {
//		log.Println("BuyResultExchange() snft not find error.")
//		return nil
//	}
//	MnftRec := Nfts{}
//	var to string
//	var accountInfo *contracts.Account
//	if nftaddr[:3] == "0x8" {
//		var gerr error
//		accountInfo, gerr = GetMergeLevel(OldNftaddr, exchangeTx.BlockNumber)
//		if gerr != nil {
//			log.Println("BuyResultExchange() ger mergelevel error.")
//			return ErrBlockchain
//		}
//		switch len(OldNftaddr) {
//		case SnftExchangeChip:
//			if accountInfo.MergeLevel > 0 {
//				err := nft.db.Where("nftaddr = ?", OldNftaddr[:len(OldNftaddr)-1]+"m").First(&MnftRec)
//				if err.Error != nil {
//					log.Println("BuyResultExchange() database err =", err.Error)
//					return err.Error
//				}
//			}
//		case SnftExchangeSnft:
//			if accountInfo.MergeLevel > 1 {
//				err := nft.db.Where("nftaddr = ?", OldNftaddr+"m").First(&MnftRec)
//				if err.Error != nil {
//					log.Println("BuyResultExchange() database err =", err.Error)
//					return err.Error
//				}
//			}
//		case SnftExchangeColletion:
//			if accountInfo.MergeLevel > 2 {
//				err := nft.db.Where("nftaddr = ?", OldNftaddr+"mm").First(&MnftRec)
//				if err.Error != nil {
//					log.Println("BuyResultExchange() database err =", err.Error)
//					return err.Error
//				}
//			}
//			/*	case SnftExchangeStage:
//				if accountInfo.MergeLevel > 3 {
//					err := nft.db.Where("nftaddr = ?", nftaddr+"m").First(&MnftRec)
//					if err.Error != nil {
//						log.Println("BuyResultWTransfer() database err =", err.Error)
//						return err.Error
//					}
//				}*/
//		}
//
//	}
//	return nft.db.Transaction(func(tx *gorm.DB) error {
//		sysInfo := SysInfos{}
//		err = tx.Model(&SysInfos{}).Last(&sysInfo)
//		if err.Error != nil {
//			if err.Error != gorm.ErrRecordNotFound {
//				log.Println("BuyResultExchange() SysInfos err=", err)
//				return ErrCollectionNotExist
//			}
//		}
//		sysNft := Sysnfts{}
//		snft := nftaddress[:len(nftaddress)-1]
//		err := tx.Where("snft = ?", snft).First(&sysNft)
//		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
//			log.Println("BuyResultExchange() database err=", err.Error)
//			return ErrCollectionNotExist
//		}
//		switch len(nftaddr) {
//		case SnftExchangeStage:
//			fmt.Println("BuyResultExchange() exchange 38 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
//			//err := tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Delete(&Nfts{})
//			//if err.Error != nil {
//			//	log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//			//	return err.Error
//			//}
//			err := tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 39 Snftstage err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr+"m").Update("exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 39 Snftstage err=", err.Error)
//				return err.Error
//			}
//			//err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").Update("exchange", 1)
//			//if err.Error != nil {
//			//	log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//			//	return err.Error
//			//}
//			err = tx.Model(&Sysnfts{}).Where("Snftstage = ?", nftaddr).Delete(&Sysnfts{})
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//				return err.Error
//			}
//			nftcount := err.RowsAffected
//			err = tx.Model(&Collects{}).Where("Snftstage = ?", nftaddr).Delete(&Collects{})
//			if err.Error != nil {
//				log.Println("BuyResultExchange() 38 exchange deleted collect recorde err= ", err.Error)
//				return err.Error
//			}
//			if sysInfo.Snfttotal >= 256 {
//				sysInfo.Snfttotal -= uint64(nftcount)
//				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
//				if err.Error != nil {
//					log.Println("BuyResultExchange() 38 exchange sub SysInfos snfttotal err=", err.Error)
//					return err.Error
//				}
//			}
//			//NftCatch.SetFlushFlag()
//			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
//			GetRedisCatch().SetDirtyFlag(SnftExchange)
//		case SnftExchangeColletion:
//			fmt.Println("BuyResultExchange() exchange 40 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
//			/*	err := tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Delete(&Nfts{})
//				if err.Error != nil {
//					log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
//					return err.Error
//				}*/
//			err := tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Update("Exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Nfts{}).Where("Snftcollection = ? and (mergetype = 1 or mergetype = 2)", nftaddr+"m").Update("Exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
//				return err.Error
//			}
//			oldNfts := Nfts{}
//			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress[:len(nftaddress)-2]+"mm").First(&oldNfts)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress[:len(nftaddress)-3]+"mmm").Update("Exchangecnt", gorm.Expr("Exchangecnt + ?", oldNfts.Exchangecnt))
//			if err.Error != nil {
//				log.Println("BuyResultExchange() err=", err.Error)
//				return err.Error
//			}
//			//err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mm").Update("Exchange", 1)
//			//if err.Error != nil {
//			//	log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
//			//	return err.Error
//			//}
//			nftcount := uint64(err.RowsAffected)
//			err = tx.Model(&Sysnfts{}).Where("Snftcollection = ?", nftaddr).Delete(&Sysnfts{})
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
//				nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
//			if err.Error != nil {
//				log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
//				return err.Error
//			}
//			if sysInfo.Snfttotal >= 16 {
//				sysInfo.Snfttotal -= nftcount / 16
//				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
//				if err.Error != nil {
//					log.Println("BuyResultExchange() 39 exchange sub SysInfos snfttotal err=", err.Error)
//					return err.Error
//				}
//			}
//			//NftCatch.SetFlushFlag()
//			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
//			GetRedisCatch().SetDirtyFlag(SnftExchange)
//		case SnftExchangeSnft:
//			fmt.Println("BuyResultExchange() exchange 41 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
//			err := tx.Model(&Sysnfts{}).Where("Snft = ?", nftaddr).Delete(&Sysnfts{})
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//				return err.Error
//			}
//			/*	err = tx.Model(&Nfts{}).Where("Snft = ?", nftaddr).Delete(&Nfts{})
//				if err.Error != nil {
//					log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
//					return err.Error
//				}*/
//			err = tx.Model(&Nfts{}).Where("Snft = ?", nftaddr).Update("Exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Nfts{}).Where("nftaddr = ? and mergetype = 1", nftaddr+"m").Update("Exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
//				return err.Error
//			}
//			oldNfts := Nfts{}
//			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress[:len(nftaddress)-1]+"m").First(&oldNfts)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress[:len(nftaddress)-2]+"mm").Update("Exchangecnt", gorm.Expr("Exchangecnt + ?", oldNfts.Exchangecnt))
//			if err.Error != nil {
//				log.Println("BuyResultExchange() err=", err.Error)
//				return err.Error
//			}
//			nftcount := int(err.RowsAffected)
//			if collectRec.Totalcount >= 16 {
//				collectRec.Totalcount -= nftcount
//				if collectRec.Totalcount != 0 {
//					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
//						nftRec.Collectcreator, nftRec.Collections).Update("totalcount", collectRec.Totalcount)
//					if err.Error != nil {
//						log.Println("BuyResultExchange() 40 exchange update collect recorde err= ", err.Error)
//						return err.Error
//					}
//				} else {
//					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
//						nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
//					if err.Error != nil {
//						log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
//						return err.Error
//					}
//				}
//			}
//			if sysInfo.Snfttotal > 0 {
//				sysInfo.Snfttotal -= 1
//				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
//				if err.Error != nil {
//					log.Println("BuyResultExchange() 40 exchange sub  SysInfos snfttotal err=", err.Error)
//					return err.Error
//				}
//			}
//			//NftCatch.SetFlushFlag()
//			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
//			GetRedisCatch().SetDirtyFlag(SnftExchange)
//		case SnftExchangeChip:
//			fmt.Println("BuyResultExchange() exchange 42 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
//			//err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress).Delete(&Nfts{})
//			//if err.Error != nil {
//			//	log.Println("BuyResultExchange() err=", err.Error)
//			//	return err.Error
//			//}
//			err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress).Update("exchange", 1)
//			if err.Error != nil {
//				log.Println("BuyResultExchange() err=", err.Error)
//				return err.Error
//			}
//			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress[:len(nftaddress)-1]+"m").Update("Exchangecnt", gorm.Expr("Exchangecnt + ?", 1))
//			if err.Error != nil {
//				log.Println("BuyResultExchange() err=", err.Error)
//				return err.Error
//			}
//
//			if collectRec.Totalcount >= 1 {
//				collectRec.Totalcount -= 1
//				if collectRec.Totalcount != 0 {
//					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
//						nftRec.Collectcreator, nftRec.Collections).Update("totalcount", collectRec.Totalcount)
//					if err.Error != nil {
//						log.Println("BuyResultExchange() 40 exchange update collect recorde err= ", err.Error)
//						return err.Error
//					}
//				} else {
//					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
//						nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
//					if err.Error != nil {
//						log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
//						return err.Error
//					}
//				}
//				if sysNft.Chipcount-1 != 0 {
//					err = tx.Model(&Sysnfts{}).Where("Snft = ?", snft).Update("chipcount", sysNft.Chipcount-1)
//					if err.Error != nil {
//						log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//						return err.Error
//					}
//				} else {
//					err = tx.Model(&Sysnfts{}).Where("Snft = ?", snft).Delete(&sysNft)
//					if err.Error != nil {
//						log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
//						return err.Error
//					}
//					sysInfo.Snfttotal -= 1
//					err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal)
//					if err.Error != nil {
//						log.Println("BuyResultExchange() add  SysInfos snfttotal err=", err.Error)
//						return err.Error
//					}
//				}
//				//NftCatch.SetFlushFlag()
//				GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
//				GetRedisCatch().SetDirtyFlag(SnftExchange)
//			}
//		}
//		if nftaddr[:3] == "0x8" {
//			nerr := SnftMerge(OldNftaddr, to, accountInfo, tx)
//			if nerr != nil {
//				fmt.Println("BuyResultExchange() SnftMerge err=", nerr)
//				return nerr
//			}
//		}
//		return nil
//	})
//	return nil
//}

func (nft NftDb) BuyResultExchange(exchangeTx *contracts.NftTx) error {
	nftaddr := strings.ToLower(exchangeTx.NftAddr)
	OldNftaddr := nftaddr

	switch len(nftaddr) {
	case SnftExchangeStage:
		nftaddr = nftaddr + "mmm"
	case SnftExchangeColletion:
		nftaddr = nftaddr + "mm"
	case SnftExchangeSnft:
		nftaddr = nftaddr + "m"
	}
	var accountInfo *contracts.Account
	if nftaddr[:3] == "0x8" {
		var gerr error
		accountInfo, gerr = GetMergeLevel(OldNftaddr, exchangeTx.BlockNumber)
		if gerr != nil {
			log.Println("BuyResultExchange() ger mergelevel error.")
			return ErrBlockchain
		}
	}
	var nftRec Nfts
	err := nft.db.Where("nftaddr = ?", nftaddr).First(&nftRec)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			log.Println("BuyResultExchange() nft not find ")
			return nil
		}
		log.Println("BuyResultExchange() nft find err=", err.Error)
		return ErrNftNotExist
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
	auctRec := Auction{}
	aucFlag := false
	err = nft.db.Where("Contract = ? and tokenid = ?", nftRec.Contract, nftRec.Tokenid).First(&auctRec)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("BuyResultWTransfer() auction not find err=", err.Error)
			return err.Error
		}
	} else {
		aucFlag = true
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		if aucFlag {
			err = tx.Model(&Auction{}).Where("id = ?", auctRec.ID).Delete(&Auction{})
			if err.Error != nil {
				fmt.Println("BuyResultWithWAmount() delete auction record err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Bidding{}).Where("Auctionid = ?", auctRec.ID).Delete(&Bidding{})
			if err.Error != nil {
				fmt.Println("BuyResultWithWAmount() delete bid record err=", err.Error)
				return err.Error
			}
		}
		sysInfo := SysInfos{}
		err = tx.Model(&SysInfos{}).Last(&sysInfo)
		if err.Error != nil {
			if err.Error != gorm.ErrRecordNotFound {
				log.Println("BuyResultExchange() SysInfos err=", err)
				return ErrCollectionNotExist
			}
		}
		switch len(OldNftaddr) {
		case SnftExchangeStage:
			fmt.Println("BuyResultExchange() exchange 38 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			//err := tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Delete(&Nfts{})
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
			//	return err.Error
			//}
			//err := tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("exchange", 1)
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() exchange 39 Snftstage err=", err.Error)
			//	return err.Error
			//}
			err = tx.Model(&Nfts{}).Where("Snftstage = ? ", OldNftaddr+"m").Update("exchange", 1)
			if err.Error != nil {
				log.Println("BuyResultExchange() exchange 39 Snftstage err=", err.Error)
				return err.Error
			}
			//err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").Update("exchange", 1)
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() exchange 38 Snftstage err=", err.Error)
			//	return err.Error
			//}
			err = tx.Model(&Collects{}).Where("Snftstage = ?", OldNftaddr).Delete(&Collects{})
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
			//GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			//GetRedisCatch().SetDirtyFlag(SnftExchange)
		case SnftExchangeColletion:
			fmt.Println("BuyResultExchange() exchange 40 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			/*	err := tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Delete(&Nfts{})
				if err.Error != nil {
					log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
					return err.Error
				}*/
			//err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Update("Exchange", 1)
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
			//	return err.Error
			//}
			err = tx.Model(&Nfts{}).Where("Snftcollection = ? and (mergetype = 1 or mergetype = 2)", OldNftaddr+"m").Update("Exchange", 1)
			if err.Error != nil {
				log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
				return err.Error
			}
			//err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mm").Update("Exchange", 1)
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() exchange 40 Snftstage err=", err.Error)
			//	return err.Error
			//}
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
			//GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			//GetRedisCatch().SetDirtyFlag(SnftExchange)
		case SnftExchangeSnft:
			fmt.Println("BuyResultExchange() exchange 41 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			/*	err = tx.Model(&Nfts{}).Where("Snft = ?", nftaddr).Delete(&Nfts{})
				if err.Error != nil {
					log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
					return err.Error
				}*/
			err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Update("Exchange", 1)
			if err.Error != nil {
				log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
				return err.Error
			}
			//snftAddr := strings.Replace(nftaddr, "m", "", -1)
			//err = tx.Model(&Nfts{}).Where("Snft = ?", snftAddr).Update("Exchange", 1)
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() exchange 41 Snftstage err=", err.Error)
			//	return err.Error
			//}
			if collectRec.Totalcount >= 16 {
				collectRec.Totalcount -= 16
				if collectRec.Totalcount != 0 {
					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
						nftRec.Collectcreator, nftRec.Collections).Update("totalcount", collectRec.Totalcount)
					if err.Error != nil {
						log.Println("BuyResultExchange() 40 exchange update collect recorde err= ", err.Error)
						return err.Error
					}
				} else {
					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
						nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
					if err.Error != nil {
						log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
						return err.Error
					}
				}
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
			//GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			//GetRedisCatch().SetDirtyFlag(SnftExchange)
		case SnftExchangeChip:
			fmt.Println("BuyResultExchange() exchange 42 nftaddr=", nftaddr, " blocknumber=", exchangeTx.BlockNumber)
			//err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddress).Delete(&Nfts{})
			//if err.Error != nil {
			//	log.Println("BuyResultExchange() err=", err.Error)
			//	return err.Error
			//}
			err := tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Update("exchange", 1)
			if err.Error != nil {
				log.Println("BuyResultExchange() err=", err.Error)
				return err.Error
			}
			if collectRec.Totalcount >= 1 {
				collectRec.Totalcount -= 1
				if collectRec.Totalcount != 0 {
					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
						nftRec.Collectcreator, nftRec.Collections).Update("totalcount", collectRec.Totalcount)
					if err.Error != nil {
						log.Println("BuyResultExchange() 40 exchange update collect recorde err= ", err.Error)
						return err.Error
					}
				} else {
					err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
						nftRec.Collectcreator, nftRec.Collections).Delete(&Collects{})
					if err.Error != nil {
						log.Println("BuyResultExchange() 40 exchange deleted collect recorde err= ", err.Error)
						return err.Error
					}
				}
				//NftCatch.SetFlushFlag()
				//GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
				//GetRedisCatch().SetDirtyFlag(SnftExchange)
			}
		}
		if nftaddr[:3] == "0x8" {
			to := strings.ToLower(accountInfo.Owner.String())
			nerr := SnftMerge(OldNftaddr, to, accountInfo, tx, nft.db)
			if nerr != nil {
				fmt.Println("BuyResultExchange() SnftMerge err=", nerr)
				return nerr
			}
		}
		return nil
	})
	return nil
}

func (nft NftDb) BuyResultWPledge(Tx *contracts.NftTx) error {
	nftaddr := strings.ToLower(Tx.NftAddr)
	if nftaddr[:3] != "0x8" {
		fmt.Println("BuyResultWPledge() nftaddr=", nftaddr, " err=not snft")
		return nil
	}
	fmt.Println("BuyResultWPledge() nftaddr=", nftaddr, " transType=", Tx.TransType)
	fmt.Println("BuyResultWPledge() transType=", Tx.TransType)
	switch len(nftaddr) {
	//case SnftExchangeChip:
	//	if Tx.TransType == contracts.WormHolesPledge {
	//		err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Update("Pledgestate", Pledge.String())
	//		if err.Error != nil {
	//			log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
	//			return err.Error
	//		}
	//	} else {
	//		err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddr).Update("Pledgestate", NoPledge.String())
	//		if err.Error != nil {
	//			log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
	//			return err.Error
	//		}
	//	}
	case SnftExchangeSnft:
		nftRec := Nfts{}
		err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"m").First(&nftRec)
		if err.Error != nil {
			log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
			return err.Error
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			if Tx.TransType == contracts.WormHolesPledge {
				//err = ClearAuction(tx, &nftRec)
				//if err != nil {
				//	log.Println("BuyResultWPledge() ClearAuction err=", err.Error)
				//	return err.Error
				//}
				err = tx.Model(&Nfts{}).Where("Snft = ?", nftaddr).Update("Pledgestate", Pledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"m").Update("Pledgestate", Pledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
					return err.Error
				}
			} else {
				err = tx.Model(&Nfts{}).Where("Snft = ?", nftaddr).Update("Pledgestate", NoPledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"m").Update("Pledgestate", NoPledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
					return err.Error
				}
			}
			return nil
		})
	case snftCollectionOffset:
		nftRec := Nfts{}
		err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mm").First(&nftRec)
		if err.Error != nil {
			log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
			return err.Error
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			if Tx.TransType == contracts.WormHolesPledge {
				//err = ClearAuction(tx, &nftRec)
				//if err != nil {
				//	log.Println("BuyResultWPledge() ClearAuction err=", err.Error)
				//	return err.Error
				//}
				err := tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Update("Pledgestate", Pledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() snftCollectionOffset err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mm").Update("Pledgestate", Pledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() snftCollectionOffset err=", err.Error)
					return err.Error
				}
			} else {
				err := tx.Model(&Nfts{}).Where("Snftcollection = ?", nftaddr).Update("Pledgestate", NoPledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() snftCollectionOffset err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mm").Update("Pledgestate", NoPledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() snftCollectionOffset err=", err.Error)
					return err.Error
				}
			}
			return nil
		})
	case SnftExchangeStage:
		nftRec := Nfts{}
		err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").First(&nftRec)
		if err.Error != nil {
			log.Println("BuyResultWPledge() SnftExchangeSnft err=", err.Error)
			return err.Error
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			if Tx.TransType == contracts.WormHolesPledge {
				//err = ClearAuction(tx, &nftRec)
				//if err != nil {
				//	log.Println("BuyResultWPledge() ClearAuction err=", err.Error)
				//	return err.Error
				//}
				err = tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("Pledgestate", Pledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeStage err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").Update("Pledgestate", Pledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeStage err=", err.Error)
					return err.Error
				}
			} else {
				err = tx.Model(&Nfts{}).Where("Snftstage = ?", nftaddr).Update("Pledgestate", NoPledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeStage err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Nfts{}).Where("nftaddr = ?", nftaddr+"mmm").Update("Pledgestate", NoPledge.String())
				if err.Error != nil {
					log.Println("BuyResultWPledge() SnftExchangeStage err=", err.Error)
					return err.Error
				}
			}
			return nil
		})
	}
	return nil
}
