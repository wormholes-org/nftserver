package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type OfferList struct {
	Snft []struct {
		ContractAddr string `json:"contractAddr"`
		TokenId      string `json:"tokenId"`
	}
	AuthSig      string `json:"authSig"`
	Price        string `json:"price"`
	PayChannel   string `json:"payChannel"`
	CurrencyType string `json:"currencyType"`
	DeadTime     string `json:"deadTime"`
}

type SellList struct {
	Snft []struct {
		ContractAddr string `json:"contractAddr"`
		TokenId      string `json:"tokenId"`
	}
	AuthSig string `json:"authSig"`
}

func (nft NftDb) BatchBuyingNft(userAddr, offerlist string, selllist string) error {
	userAddr = strings.ToLower(userAddr)
	fmt.Println("BatchBuyingNft() userAddr=", userAddr, "time=", time.Now().String())

	if UserSync.LockTran(userAddr) {
		return ErrUserTrading
	} else {
		defer UserSync.UnLockTran(userAddr)
	}
	//UserSync.Lock(userAddr)
	//defer UserSync.UnLock(userAddr)
	offerList := OfferList{}
	var uerr error
	if offerlist != "" {
		uerr := json.Unmarshal([]byte(offerlist), &offerList)
		if uerr != nil {
			log.Println("BatchBuyingNft() input offerList err = ", uerr)
			return ErrDataFormat
		}
	}
	for _, s := range offerList.Snft {
		price, _ := strconv.ParseUint(offerList.Price, 10, 64)
		deadTime, _ := strconv.ParseInt(offerList.DeadTime, 10, 64)
		uerr = nft.MakeOffer(userAddr, s.ContractAddr, s.TokenId,
			offerList.PayChannel, offerList.CurrencyType, price, "", deadTime, "", "", offerList.AuthSig)
		if uerr != nil {
			log.Println("BatchBuyingNft() makeoffer err = ", uerr)
			return uerr
		}
	}
	sellList := SellList{}
	if selllist != "" {
		uerr = json.Unmarshal([]byte(selllist), &sellList)
		if uerr != nil {
			log.Println("BatchBuyingNft() input buylist err = ", uerr)
			return ErrDataFormat
		}
	} else {
		log.Println("BatchBuyingNft() SellList{} is nil ")
		return nil
	}
	txhashs := []string{}
	blocknumber := contracts.GetCurrentBlockNumber()
	for _, s := range sellList.Snft {
		var auctionRec Auction
		dberr := nft.db.Where("contract = ? AND tokenid = ?", s.ContractAddr, s.TokenId).First(&auctionRec)
		if dberr.Error != nil {
			if dberr.Error != gorm.ErrRecordNotFound {
				fmt.Println("BatchBuyingNft() RecordNotFound")
				return errors.New(ErrDataBase.Error() + dberr.Error.Error())
			}
			continue
		}
		if auctionRec.SellState == SellStateWait.String() {
			log.Println("BatchBuyingNft() on sale.")
			continue
		}
		if auctionRec.Selltype != SellTypeFixPrice.String() {
			log.Println("BatchBuyingNft() selltype is not FixPrice.")
			continue
		}
		var nftrecord Nfts
		dberr = nft.db.Where("contract = ? AND tokenid =?", s.ContractAddr, s.TokenId).First(&nftrecord)
		if dberr.Error != nil {
			fmt.Println("BatchBuyingNft() bidprice not find nft err= ", dberr.Error)
			continue
		}
		if nftrecord.Mergetype != nftrecord.Mergelevel {
			fmt.Println("BatchBuyingNft() snft has been merged")
			continue
		}
		auctRec := Auction{}
		auctRec.SellState = SellStateWait.String()
		auctRec.Price = auctionRec.Startprice
		dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
		if dberr.Error != nil {
			log.Println("BatchBuyingNft() update auction record err=", dberr.Error)
			return errors.New(ErrDataBase.Error() + dberr.Error.Error())
		}
		price := big.NewInt(0).SetUint64(auctionRec.Startprice)
		price = price.Mul(price, big.NewInt(1000000000))
		sell := contracts.Seller1{}
		if auctionRec.Sellauthsig == "" {
			err := json.Unmarshal([]byte(auctionRec.Tradesig), &sell)
			if err != nil {
				log.Println("BatchBuyingNft() seller Unmarshal() err=", err)
				continue
			}
		} else {
			sell.Price = hexutil.EncodeBig(price)
			sell.Exchanger = ExchangeOwer
			sell.Nftaddress = auctionRec.Nftaddr
		}
		buyer := contracts.Buyer{}
		buyer.Nftaddress = auctionRec.Nftaddr
		buyer.Exchanger = ExchangeOwer
		buyer.Price = hexutil.EncodeBig(price)
		buyer.Seller = auctionRec.Ownaddr
		txhash, err := contracts.BatchAuthExchangeTrans(sell, buyer, auctionRec.Sellauthsig, sellList.AuthSig, ExchangerAuth, contracts.SuperAdminAddr)
		if err != nil {
			log.Println("BatchBuyingNft() AuthWormTrans err=", err)
			auctRec = Auction{}
			auctRec.SellState = SellStateStart.String()
			dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
				auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
			if dberr.Error != nil {
				log.Println("BatchBuyingNft() update auction record err=", dberr.Error)
				return errors.New(ErrDataBase.Error() + dberr.Error.Error())
			}
			continue
		}
		txhashs = append(txhashs, txhash)
	}
	if len(txhashs) != 0 {
		for i := blocknumber + 10; blocknumber < i; {
			blocknumber = contracts.GetCurrentBlockNumber()
			txflag := false
			for _, txhash := range txhashs {
				txStatus, err := contracts.GetTransStatus(txhash)
				if err != nil {
					txflag = false
					break
				}
				if txStatus {
					txflag = true
				} else {
					txflag = false
					break
				}
			}
			if txflag {
				log.Println("BatchBuyingNft() trans ok ")
				return nil
			}
			time.Sleep(time.Second)
		}
	}
	GetRedisCatch().SetDirtyFlag(TradingDirtyName)
	return ErrWaitingClose
}
