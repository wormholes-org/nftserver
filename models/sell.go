package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) MakeOffer(userAddr,
	contractAddr,
	tokenId string,
	PayChannel string,
	CurrencyType string,
	price uint64,
	TradeSig string,
	dead_time int64,
	voteStage string,
	sigdata string,
	buyauthsig string) error {
	userAddr = strings.ToLower(userAddr)
	contractAddr = strings.ToLower(contractAddr)

	fmt.Println("MakeOffer() userAddr=", userAddr, "      time=", time.Now().String())
	UserSync.Lock(userAddr)
	defer UserSync.UnLock(userAddr)
	if price <= LowPrice {
		fmt.Println("MakeOffer() price <= 0.")
		return ErrBidOutRange
	}
	if ExchangerAuth == "" {
		fmt.Println("MakeOffer() Unauthorized exchange error.")
		return ErrUnauthExchange
	}
	if !nft.UserKYCAduit(userAddr) {
		return ErrUserNotVerify
	}
	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid =?", contractAddr, tokenId).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("MakeOffer() bidprice not find nft err= ", err.Error)
		return ErrNftNotExist
	}
	if nftrecord.Pledgestate == Pledge.String() {
		log.Println("MakeOffer() snft Has pledged.", err.Error)
		return ErrSnftPledge
	}
	if nftrecord.Ownaddr == userAddr {
		fmt.Println("MakeOffer() don't buy your own nft.")
		return ErrBuyOwn
	}
	valid, errmsg, cerr := WormsAmountValid(price, userAddr)
	if cerr != nil {
		return ErrGetBalance
	}
	if !valid {
		return errors.New(ErrBlockchain.Error() + errmsg)
	}
	if buyauthsig == "" {
		rerr := MakeofferSigVerify(TradeSig, userAddr)
		if rerr != nil {
			log.Println("MakeOffer() SigVerify buyerSig err=", rerr)
			return rerr
		}
	}

	var auctionRec Auction
	err = nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&auctionRec)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			fmt.Println("MakeOffer() RecordNotFound")
			auctionRec = Auction{}
			auctionRec.Selltype = SellTypeBidPrice.String()
			auctionRec.Paychan = PayChannel
			auctionRec.Ownaddr = nftrecord.Ownaddr
			auctionRec.Nftaddr = nftrecord.Nftaddr
			auctionRec.Nftid = nftrecord.ID
			auctionRec.Contract = contractAddr
			auctionRec.Tokenid = tokenId
			auctionRec.Count = 1
			auctionRec.Currency = CurrencyType
			//auctionRec.Startprice = price
			//auctionRec.Endprice = price
			auctionRec.Startdate = time.Now().Unix()
			auctionRec.Enddate = dead_time
			auctionRec.Signdata = sigdata
			auctionRec.Tradesig = TradeSig
			//auctionRec.Authsig = authsig
			//auctionHistory := AuctionHistory{}
			//auctionHistory.AuctionRecord = auctionRec.AuctionRecord
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err = tx.Model(&auctionRec).Create(&auctionRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create auctionRec record err=", err.Error)
					return ErrDataBase
				}
				//err = tx.Model(&AuctionHistory{}).Create(&auctionHistory)
				//if err.Error != nil {
				//	fmt.Println("MakeOffer() create auctionHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftrecord = Nfts{}
				nftrecord.Selltype = auctionRec.Selltype
				nftrecord.Sellprice = auctionRec.Startprice
				nftrecord.Offernum = 1
				nftrecord.Maxbidprice = price
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
				if err.Error != nil {
					fmt.Println("MakeOffer() update record err=", err.Error)
					return ErrDataBase
				}
				bidRec := Bidding{}
				bidRec.Bidaddr = userAddr
				bidRec.Auctionid = auctionRec.ID
				bidRec.Contract = contractAddr
				bidRec.Tokenid = tokenId
				bidRec.Count = 1
				bidRec.Price = price
				bidRec.Currency = CurrencyType
				bidRec.Paychan = PayChannel
				bidRec.Tradesig = TradeSig
				bidRec.Buyauthsig = buyauthsig
				bidRec.Bidtime = time.Now().Unix()
				bidRec.Signdata = sigdata
				bidRec.Deadtime = dead_time
				bidRec.Nftid = auctionRec.Nftid
				bidRec.VoteStage = voteStage
				//bidRecHistory := BiddingHistory{}
				//bidRecHistory.BidRecord = bidRec.BidRecord
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
					return ErrDataBase
				}
				//err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				//if err.Error != nil {
				//	fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				//NftCatch.SetFlushFlag()
				GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
				fmt.Println("MakeOffer() RecordNotFound OK")
				return nil
			})
		}
		return ErrNftNotSell
	}
	//if time.Now().Unix() < auctionRec.Startdate {
	//	return ErrAuctionNotBegan
	//}
	//NftCatch.SetFlushFlag()
	GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
	if auctionRec.Selltype == SellTypeHighestBid.String() {
		//addrs, err := ethhelper.BalanceOfWeth()
		fmt.Println("MakeOffer() Selltype == SellTypeHighestBid")
		if time.Now().Unix() >= auctionRec.Enddate {
			fmt.Println("MakeOffer() time.Now().Unix() >= auctionRec.Enddate")
			return ErrAuctionEnd
		}
		if auctionRec.Startprice > price {
			fmt.Println("MakeOffer() auctionRec.Startprice > price")
			return ErrBidOutRange
		}
		var bidRec Bidding
		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
		if err.Error == nil {
			fmt.Println("MakeOffer() first bidding.")
			bidRec = Bidding{}
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Buyauthsig = buyauthsig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.VoteStage = voteStage
			bidRec.Deadtime = dead_time
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
					return ErrDataBase
				}
				//bidRecHistory := BiddingHistory(bidRec)
				//err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
				//if err.Error != nil {
				//	log.Println("MakeOffer() update bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftTemp := Nfts{}
				if nftrecord.Maxbidprice < price {
					nftTemp.Maxbidprice = price
					err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
						auctionRec.Contract, auctionRec.Tokenid).Updates(&nftTemp)
					if err.Error != nil {
						log.Println("MakeOffer() update record err=", err.Error)
						return ErrDataBase
					}
				}
				fmt.Println("MakeOffer() first bidding OK.")
				return nil
			})
		} else {
			bidRec = Bidding{}
			bidRec.Bidaddr = userAddr
			bidRec.Auctionid = auctionRec.ID
			bidRec.Nftid = auctionRec.Nftid
			bidRec.Contract = contractAddr
			bidRec.Tokenid = tokenId
			bidRec.Count = 1
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Deadtime = dead_time
			bidRec.Tradesig = TradeSig
			bidRec.Buyauthsig = buyauthsig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.VoteStage = voteStage
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create record err=", err.Error)
					return ErrDataBase
				}
				//bidRecHistory := BiddingHistory{}
				//bidRecHistory.BidRecord = bidRec.BidRecord
				//err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				//if err.Error != nil {
				//	fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftTemp := Nfts{}
				if nftrecord.Maxbidprice < price {
					nftTemp.Maxbidprice = price
				}
				nftTemp.Offernum = nftrecord.Offernum + 1
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftTemp)
				if err.Error != nil {
					log.Println("MakeOffer() update record err=", err.Error)
					return ErrDataBase
				}
				fmt.Println("MakeOffer() change bidding OK.")
				return nil
			})
		}
	}
	if auctionRec.Selltype == SellTypeBidPrice.String() {
		fmt.Println("MakeOffer() Selltype == SellTypeBidPrice")
		var bidRec Bidding
		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
		if err.Error == nil {
			bidRec = Bidding{}
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Buyauthsig = buyauthsig
			bidRec.VoteStage = voteStage
			bidRec.Bidtime = time.Now().Unix()
			bidRec.Deadtime = dead_time
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
					return ErrDataBase
				}
				//bidRecHistory := BiddingHistory(bidRec)
				//err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
				//if err.Error != nil {
				//	fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftTemp := Nfts{}
				if nftrecord.Maxbidprice < price {
					nftTemp.Maxbidprice = price
					err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
						auctionRec.Contract, auctionRec.Tokenid).Updates(&nftTemp)
					if err.Error != nil {
						log.Println("MakeOffer() update record err=", err.Error)
						return ErrDataBase
					}
				}
				fmt.Println("MakeOffer() change bidding OK.")
				return nil
			})
		} else {
			return nft.db.Transaction(func(tx *gorm.DB) error {
				bidRec := Bidding{}
				bidRec.Bidaddr = userAddr
				bidRec.Auctionid = auctionRec.ID
				bidRec.Nftid = auctionRec.Nftid
				bidRec.Contract = contractAddr
				bidRec.Tokenid = tokenId
				bidRec.Count = 1
				bidRec.Price = price
				bidRec.Currency = CurrencyType
				bidRec.Paychan = PayChannel
				bidRec.Tradesig = TradeSig
				bidRec.Buyauthsig = buyauthsig
				bidRec.Bidtime = time.Now().Unix()
				bidRec.Deadtime = dead_time
				bidRec.Signdata = sigdata
				bidRec.VoteStage = voteStage
				//bidRecHistory := BiddingHistory{}
				//bidRecHistory.BidRecord = bidRec.BidRecord
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
					return ErrDataBase
				}
				//err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				//if err.Error != nil {
				//	fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftTemp := Nfts{}
				if nftrecord.Maxbidprice < price {
					nftTemp.Maxbidprice = price
				}
				nftTemp.Offernum = nftrecord.Offernum + 1
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftTemp)
				if err.Error != nil {
					log.Println("MakeOffer() update record err=", err.Error)
					return ErrDataBase
				}
				fmt.Println("MakeOffer() first bidding OK.")
				return nil
			})
		}
	}
	if auctionRec.Selltype == SellTypeFixPrice.String() {
		fmt.Println("MakeOffer() Selltype == SellTypeFixPrice")
		var bidRec Bidding
		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
		if err.Error == nil {
			bidRec = Bidding{}
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Buyauthsig = buyauthsig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.Deadtime = dead_time
			bidRec.VoteStage = voteStage
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
					return ErrDataBase
				}
				//bidRecHistory := BiddingHistory(bidRec)
				//err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
				//if err.Error != nil {
				//	fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftTemp := Nfts{}
				if nftrecord.Maxbidprice < price {
					nftTemp.Maxbidprice = price
					err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
						auctionRec.Contract, auctionRec.Tokenid).Updates(&nftTemp)
					if err.Error != nil {
						log.Println("MakeOffer() update record err=", err.Error)
						return ErrDataBase
					}
				}
				fmt.Println("MakeOffer() change bidding OK.")
				return nil
			})
		} else {
			return nft.db.Transaction(func(tx *gorm.DB) error {
				bidRec := Bidding{}
				bidRec.Bidaddr = userAddr
				bidRec.Auctionid = auctionRec.ID
				bidRec.Nftid = auctionRec.Nftid
				bidRec.Contract = contractAddr
				bidRec.Tokenid = tokenId
				bidRec.Count = 1
				bidRec.Price = price
				bidRec.Currency = CurrencyType
				bidRec.Paychan = PayChannel
				bidRec.Tradesig = TradeSig
				bidRec.Buyauthsig = buyauthsig
				bidRec.Bidtime = time.Now().Unix()
				bidRec.VoteStage = voteStage
				bidRec.Deadtime = dead_time
				bidRec.Signdata = sigdata
				//bidRecHistory := BiddingHistory{}
				//bidRecHistory.BidRecord = bidRec.BidRecord
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					log.Println("MakeOffer() create bidRec record err=", err.Error)
					return ErrDataBase
				}
				//err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				//if err.Error != nil {
				//	log.Println("MakeOffer() create bidRecHistory record err=", err.Error)
				//	return ErrDataBase
				//}
				nftTemp := Nfts{}
				if nftrecord.Maxbidprice < price {
					nftTemp.Maxbidprice = price
				}
				nftTemp.Offernum = nftrecord.Offernum + 1
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftTemp)
				if err.Error != nil {
					log.Println("MakeOffer() update record err=", err.Error)
					return ErrDataBase
				}
				fmt.Println("MakeOffer() first bidding OK.")
				return nil
			})
		}
	}
	return ErrNftNotSell
}

func IsSellTypeValid(selltype string) error {
	switch selltype {
	case SellTypeFixPrice.String(), SellTypeBidPrice.String(), SellTypeHighestBid.String():
		return nil
	default:
		fmt.Println("IsSellTypeValid() sell type error.")
		return ErrSellType
	}
}

func (nft NftDb) Sell(ownAddr,
	PrivAddr string,
	contractAddr,
	tokenId string,
	sellType string,
	payChan string,
	days int,
	startPrice,
	endPrice uint64,
	royalty string,
	currency string,
	hide string,
	sigData string,
	voteStage string,
	tradeSig, sellauthsig string) error {

	ownAddr = strings.ToLower(ownAddr)
	PrivAddr = strings.ToLower(PrivAddr)
	contractAddr = strings.ToLower(contractAddr)
	if err := IsSellTypeValid(sellType); err != nil {
		fmt.Println("Sell() sell type err")
		return ErrSellType
	}
	if ExchangerAuth == "" {
		fmt.Println("Sell() Unauthorized exchange error.")
		return ErrUnauthExchange
	}
	fmt.Println("Sell() ownAddr=", ownAddr, "      time=", time.Now().String())

	if UserSync.LockTran(ownAddr) {
		return ErrUserTrading
	} else {
		defer UserSync.UnLockTran(ownAddr)
	}
	//UserSync.Lock(ownAddr)
	//defer UserSync.UnLock(ownAddr)
	if !nft.UserKYCAduit(ownAddr) {
		return ErrUserNotVerify
	}
	if sellauthsig == "" {
		rerr := SellSigVerify(tradeSig, ownAddr)
		if rerr != nil {
			log.Println("Sell() SigVerify sigData err=", rerr)
			return rerr
		}
	}
	fmt.Println(time.Now().String()[:22], "Sell() Start.",
		"tokenId=", tokenId,
		"SellType=", sellType,
		"startPrice=", startPrice,
		"endPrice=", endPrice)
	defer fmt.Println(time.Now().String()[:22], "Sell() end.")

	if startPrice <= LowPrice {
		fmt.Println("Sell() startPrice <= 0.")
		return ErrPrice
	}
	if days > ToolongAuciton {
		fmt.Println("Sell() Auction date too long.")
		return ErrAuctionDate
	}
	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid =? AND ownaddr = ?", contractAddr, tokenId, ownAddr).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("Sell() err= ", err.Error)
		return ErrNftNotExist
	}
	if nftrecord.Mergetype != nftrecord.Mergelevel {
		fmt.Println("Sell() snft has been merged")
		return ErrNftMerged
	}
	if nftrecord.Verified != Passed.String() {
		return ErrNotVerify
	}
	if nftrecord.Pledgestate == Pledge.String() {
		log.Println("MakeOffer() snft Has pledged.", err.Error)
		return ErrSnftPledge
	}
	/*if nftrecord.Mintstate != Minted.String() {
		return ErrNftNotMinted
	}*/
	//if startDate.After(endDate) {
	//	return ErrAuctionStartAfterEnd
	//}
	//if startDate.Before(time.Now()) {
	//	startDate = time.Now()
	//	//return ErrAuctionStartBeforeNow
	//}
	var auctionRec Auction
	err = nft.db.Where("contract = ? AND nftid = ? AND ownaddr = ?",
		nftrecord.Contract, nftrecord.ID, ownAddr).First(&auctionRec)
	if err.Error == nil {
		if auctionRec.Selltype != SellTypeBidPrice.String() {
			log.Println("Sell() err=", err.Error, ErrNftSelling)
			return ErrNftSelling
		} else {
			err := nft.db.Transaction(func(tx *gorm.DB) error {
				err = tx.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
				if err.Error != nil {
					log.Println("Sell() delete bid record err=", err.Error)
					return ErrDataBase
				}
				err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
				if err.Error != nil {
					log.Println("Sell() delete bidprice auction record err=", err.Error)
					return ErrDataBase
				}
				nfttab := map[string]interface{}{
					"Selltype":    SellTypeNotSale.String(),
					"Sellprice":   0,
					"Offernum":    0,
					"Maxbidprice": 0,
				}
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nfttab)
				if err.Error != nil {
					fmt.Println("Sell() update record err=", err.Error)
					return ErrDataBase
				}
				return nil
			})
			if err != nil {
				log.Println("Sell() delete bidprice err=", err)
				return err
			}
		}
	}
	auctionRec = Auction{}
	auctionRec.Selltype = sellType
	auctionRec.Paychan = payChan
	auctionRec.Ownaddr = ownAddr
	auctionRec.Nftid = nftrecord.ID
	auctionRec.Contract = contractAddr
	auctionRec.Tokenid = tokenId
	auctionRec.Nftaddr = nftrecord.Nftaddr
	auctionRec.Count = 1
	auctionRec.Currency = currency
	auctionRec.Startprice = startPrice
	auctionRec.Endprice = endPrice
	auctionRec.Privaddr = PrivAddr
	auctionRec.Sellauthsig = sellauthsig
	auctionRec.Startdate = time.Now().Unix()
	auctionRec.Enddate = time.Now().AddDate(0, 0, days).Unix()
	//auctionRec.Enddate = time.Now().Add(3 * time.Minute).Unix()
	auctionRec.Signdata = sigData
	auctionRec.VoteStage = voteStage
	auctionRec.Tradesig = tradeSig
	auctionRec.SellState = SellStateStart.String()

	if sellType == SellTypeFixPrice.String() {
		auctionRec.Startprice = startPrice
		auctionRec.Endprice = startPrice
	}
	//auctionHistory := AuctionHistory{}
	//auctionHistory.AuctionRecord = auctionRec.AuctionRecord
	sysInfoRec := SysInfos{}
	err = nft.db.Last(&sysInfoRec)
	if err.Error != nil {
		log.Println("BuyResultWithWAmount() Last(&sysInfoRec) err=", err.Error)
		return ErrDataBase
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&auctionRec).Create(&auctionRec)
		if err.Error != nil {
			log.Println("Sell() create auctionRec record err=", err.Error)
			return ErrDataBase
		}
		//err = tx.Model(&AuctionHistory{}).Create(&auctionHistory)
		//if err.Error != nil {
		//	log.Println("Sell() create auctionHistory record err=", err.Error)
		//	return ErrDataBase
		//}
		nftrecord = Nfts{}
		nftrecord.Hide = hide
		nftrecord.Selltype = sellType
		nftrecord.Sellprice = auctionRec.Startprice

		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
			auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
		if err.Error != nil {
			log.Println("Sell() update record err=", err.Error)
			return ErrDataBase
		}
		switch sellType {
		case SellTypeFixPrice.String():
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfoRec.ID).Update("Fixpricecnt", gorm.Expr("Fixpricecnt + ?", 1))
			if err.Error != nil {
				log.Println("Sell() Fixpricecnt  err= ", err.Error)
				return ErrDataBase
			}
		case SellTypeHighestBid.String():
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfoRec.ID).Update("Highestbidcnt", gorm.Expr("Highestbidcnt + ?", 1))
			if err.Error != nil {
				log.Println("Sell() Highestbidcnt  err= ", err.Error)
				return ErrDataBase
			}
		}
		/*nftrecord = Nfts{}
		nftrecord.Royalty, _ = strconv.Atoi(royalty)
		nftrecord.Royalty = nftrecord.Royalty / 100
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =? AND royalty = ?",
			auctionRec.Contract, auctionRec.Tokenid, 0).Updates(&nftrecord)
		if err.Error != nil {
			fmt.Println("Sell() update record err=", err.Error)
			return err.Error
		}*/
		//NftCatch.SetFlushFlag()
		GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
		return nil
	})
}

func (nft NftDb) GroupSell(useraddr, contract, params, alljudge, level, status, authsig string) error {
	if alljudge == "true" {
		var allauction []string
		switch status {
		case SellTypeBidPrice.String():
			sql := `SELECT bd.tokenid FROM biddings  bd  join nfts nf on nf.tokenid =bd.tokenid where nf.ownaddr = ? 
		and bd.deleted_at is null and nf.deleted_at is null and nf.mergelevel = ? limit ?`
			db := nft.db.Raw(sql, useraddr, level, 500).Find(&allauction)
			if db.Error != nil {
				log.Println("GroupSell get all auction err=", db.Error)
				return ErrDataFormat
			}
		case SellTypeNotSale.String():
			sql := `select tokenid from nfts where tokenid  not in  (select tokenid from biddings where deleted_at is null ) 
		and ownaddr =? and  tokenid  not in (select tokenid from auctions where deleted_at is null ) and 
		mergetype = mergelevel and exchange =0 and mergelevel = ? and deleted_at is null limit ?`
			db := nft.db.Raw(sql, useraddr, level, 500).Find(&allauction)
			if db.Error != nil {
				log.Println("GroupSell get all auction err=", db.Error)
				return ErrDataFormat
			}
		}
		tokenidmap := map[string]string{}
		for _, value := range allauction {
			tokenidmap[value] = value
		}
		var allSell []SellParams
		var Sell SellParams
		err := json.Unmarshal([]byte(params), &allSell)
		if err != nil {
			fmt.Println("Unmarshal input err=", err)
			return ErrDataFormat
		}
		if len(allSell) < 1 {
			log.Println("GroupSell input params data unmarshal is empty data")
			return ErrDataFormat
		}
		Sell = allSell[0]
		for _, j := range tokenidmap {
			price1, _ := strconv.ParseUint(Sell.Price1, 10, 64)
			price2, _ := strconv.ParseUint(Sell.Price2, 10, 64)
			days, _ := strconv.Atoi(strings.TrimSpace(Sell.Day))
			err = nft.Sell(useraddr, "", contract, j, Sell.SellType, Sell.PayChannel, days, price1, price2,
				"", Sell.Currency, Sell.Hide, Sell.Sig, Sell.VoteStage, Sell.TradeSig, authsig)
			if err != nil {
				log.Println("GroupSell err=", err)
				return ErrDataBase
			}
		}
	} else {
		var Sell []SellParams
		err := json.Unmarshal([]byte(params), &Sell)
		if err != nil {
			fmt.Println("Unmarshal input err=", err)
			return ErrDataFormat
		}
		for _, j := range Sell {
			price1, _ := strconv.ParseUint(j.Price1, 10, 64)
			price2, _ := strconv.ParseUint(j.Price2, 10, 64)
			days, _ := strconv.Atoi(strings.TrimSpace(j.Day))
			err = nft.Sell(j.UserAddr, "", j.ContractAddr, j.TokenId, j.SellType, j.PayChannel, days, price1, price2,
				"", j.Currency, j.Hide, j.Sig, j.VoteStage, j.TradeSig, authsig)
			if err != nil {
				fmt.Println("GroupCancelSell  GroupSell err=", err)
				return ErrDataBase
			}
		}
	}

	return nil
}

type SellerVerify struct {
	Price         string `json:"price"`
	Nftaddress    string `json:"nft_address"`
	Royalty       string `json:"royalty"`
	Metaurl       string `json:"meta_url"`
	Exclusiveflag string `json:"exclusive_flag"`
	Exchanger     string `json:"exchanger"`
	Blocknumber   string `json:"block_number"`
	Sig           string `json:"sig"`
}

func SellSigVerify(sigstr, buyerAddr string) error {
	buyer := SellerVerify{}
	err := json.Unmarshal([]byte(sigstr), &buyer)
	if err != nil {
		log.Println("SigVerify Unmarshal err=", err)
		return errors.New(ErrData.Error() + "sig data err")
	}
	msg := buyer.Price + buyer.Nftaddress + buyer.Royalty + buyer.Metaurl + buyer.Exclusiveflag + buyer.Exchanger + buyer.Blocknumber
	toaddr, rerr := contracts.RecoverAddress(msg, buyer.Sig)
	fmt.Println("toaddr =", toaddr.String(), "  buyaddr =", buyerAddr)
	if rerr != nil {
		log.Println("SigVerify() recoverAddress() err=", err)
		return errors.New(ErrData.Error() + "buyer sig recover err")
	}
	if strings.ToLower(toaddr.String()) != strings.ToLower(buyerAddr) {
		log.Println("SigVerify()   address error.")
		return errors.New(ErrData.Error() + " address error.")
	}
	return nil
}

func MakeofferSigVerify(sigstr, buyerAddr string) error {
	buyer := contracts.Buyer{}
	err := json.Unmarshal([]byte(sigstr), &buyer)
	if err != nil {
		log.Println("SigVerify Unmarshal err=", err)
		return errors.New(ErrData.Error() + "sig data err")
	}
	msg := buyer.Price + buyer.Nftaddress + buyer.Exchanger + buyer.Blocknumber + buyer.Seller
	toaddr, rerr := contracts.RecoverAddress(msg, buyer.Sig)
	fmt.Println("toaddr =", toaddr.String(), "  buyaddr =", buyerAddr)
	if rerr != nil {
		log.Println("SigVerify() recoverAddress() err=", err)
		return errors.New(ErrData.Error() + "Seller sig recover err")
	}
	if strings.ToLower(toaddr.String()) != strings.ToLower(buyerAddr) {
		log.Println("SigVerify()  Seller address error.")
		return errors.New(ErrData.Error() + " address error.")
	}
	return nil
}
