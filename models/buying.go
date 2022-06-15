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

const (
	WaitTransTime = 60
)

type BuyingParams struct {
	UserAddr     string `json:"user_addr"`
	BuyerAddr    string `json:"buyer_addr"`
	ContractAddr string `json:"contract_addr"`
	TokenId      string `json:"token_id"`
	Price        string `json:"price"`
	BuyerSig     string `json:"buyer_sig"`
	VoteStage    string `json:"vote_stage"`
	SellerSig    string `json:"seller_sig"`
}

type SellParams struct {
	UserAddr     string `json:"user_addr"`
	ContractAddr string `json:"contract_addr"`
	TokenId      string `json:"token_id"`
	Price1       string `json:"price1"`
	Price2       string `json:"price2"`
	Day          string `json:"day"`
	SellType     string `json:"sell_type"`
	PayChannel   string `json:"pay_channel"`
	Currency     string `json:"currency"`
	Hide         string `json:"hide"`
	Sig          string `json:"sig"`
	VoteStage    string `json:"vote_stage"`
	TradeSig     string `json:"trade_sig"`
}

type CancelSellParams struct {
	UserAddr     string `json:"user_addr"`
	ContractAddr string `json:"contract_addr"`
	TokenId      string `json:"token_id"`
	Sig          string `json:"sig"`
}

func WormTrans(mintState, sellerSig, buyerSig string) error {
	if mintState == Minted.String() {
		buyer := contracts.Buyer{}
		err := json.Unmarshal([]byte(buyerSig), &buyer)
		if err != nil {
			fmt.Println("WormTrans() Minted Buyer Unmarshal() err=", err)
			return err
		}
		err = contracts.ExchangeTrans(buyer, contracts.SuperAdminAddr)
		if err != nil {
			fmt.Println("WormTrans() ExchangeTrans() err=", err)
			return err
		}
	} else {
		seller := contracts.Seller2{}
		err := json.Unmarshal([]byte(sellerSig), &seller)
		if err != nil {
			fmt.Println("WormTrans() unminted Seller2 Unmarshal() err=", err)
			return err
		}
		buyer := contracts.Buyer1{}
		fmt.Println("WormTrans() Buyer buyerSig=", buyerSig)
		err = json.Unmarshal([]byte(buyerSig), &buyer)
		if err != nil {
			fmt.Println("WormTrans() unminted Buyer Unmarshal() err=", err)
			return err
		}
		err = contracts.ExchangerMint(seller, buyer, contracts.SuperAdminAddr)
		if err != nil {
			fmt.Println("WormTrans() ExchangerMint() err=", err)
			return err
		}
	}
	return nil
}

func AuthWormTrans(mintState, sellerSig, buyerSig, authSign string) (string, error) {
	var txhash string
	if mintState == Minted.String() {
		buyer := contracts.Buyer{}
		err := json.Unmarshal([]byte(buyerSig), &buyer)
		if err != nil {
			log.Println("AuthWormTrans() Minted Buyer Unmarshal() err=", err)
			return "", err
		}
		txhash, err = contracts.AuthExchangeTrans(buyer, authSign, contracts.SuperAdminAddr)
		if err != nil {
			log.Println("AuthWormTrans() ExchangeTrans() err=", err)
			return "", err
		}
	} else {
		seller := contracts.Seller2{}
		err := json.Unmarshal([]byte(sellerSig), &seller)
		if err != nil {
			log.Println("AuthWormTrans() unminted Seller2 Unmarshal() err=", err)
			return "", err
		}
		buyer := contracts.Buyer1{}
		log.Println("AuthWormTrans() Buyer buyerSig=", buyerSig)
		err = json.Unmarshal([]byte(buyerSig), &buyer)
		if err != nil {
			log.Println("AuthWormTrans() unminted Buyer Unmarshal() err=", err)
			return "", err
		}

		txhash, err = contracts.AuthExchangerMint(seller, buyer, authSign, contracts.SuperAdminAddr)
		if err != nil {
			log.Println("AuthWormTrans() ExchangerMint() err=", err)
			return "", err
		}
	}
	return txhash, nil
}

func (nft NftDb) BuyingNft(userAddr,
	buyerAddr,
	contractAddr,
	tokenId string,
	price uint64,
	buyerSig string,
	voteStage string,
	sellerSig string,
) error {
	userAddr = strings.ToLower(userAddr)
	buyerAddr = strings.ToLower(buyerAddr)
	contractAddr = strings.ToLower(contractAddr)
	fmt.Println("BuyingNft() userAddr=", userAddr, "time=", time.Now().String())
	fmt.Println("BuyingNft() buyerAddr=", buyerAddr)
	fmt.Println("BuyingNft() contractAddr=", contractAddr)
	fmt.Println("BuyingNft() tokenId=", tokenId)
	fmt.Println("BuyingNft() buyerSig=", buyerSig)
	fmt.Println("BuyingNft() sellerSig=", sellerSig)
	fmt.Println("BuyingNft() price=", price)
	fmt.Println("BuyingNft() voteStage=", voteStage)
	UserSync.Lock(userAddr)
	defer UserSync.UnLock(userAddr)

	var auctionRec Auction
	err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&auctionRec)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("BuyingNft() RecordNotFound")
			return err.Error
		}
		return ErrNftNotSell
	}
	if auctionRec.SellState == SellStateWait.String() {
		fmt.Println("BuyingNft() on sale.")
		return ErrNftSelling
	}
	if price <= LowPrice {
		fmt.Println("BuyingNft() price <= 0.")
		return ErrBidOutRange
	}
	var nftrecord Nfts
	err = nft.db.Where("contract = ? AND tokenid =?", contractAddr, tokenId).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("BuyingNft() bidprice not find nft err= ", err.Error)
		return ErrNftNotExist
	}
	if userAddr == buyerAddr {
		//err := WormTrans(nftrecord.Mintstate, auctionRec.Tradesig, buyerSig)
		//ExchangerAuth = `{"exchanger_owner":"0x01842a2cf56400a245a56955dc407c2c4137321e","to":"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900","block_number":"0x2540be400","sig":"0x7f1ca96714208959c5a75bdbf4770893b76b13c0bca26da2086c3365e537d57444f79b31498301c5c1d55400eec4b469c83a88a527159112f27ff934c222e4191b"}`
		auctRec := Auction{}
		auctRec.SellState = SellStateWait.String()
		auctRec.Price = price
		auctRec.VoteStage = voteStage
		dberr := nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
		if dberr.Error != nil {
			log.Println("BuyingNft() update auction record err=", dberr.Error)
			return dberr.Error
		}
		txhash, err := AuthWormTrans(nftrecord.Mintstate, auctionRec.Tradesig, buyerSig, ExchangerAuth)
		if err != nil {
			log.Println("BuyingNft() AuthWormTrans err=", err)
			auctRec = Auction{}
			auctRec.SellState = SellStateStart.String()
			auctRec.Price = price
			dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
				auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
			if dberr.Error != nil {
				log.Println("BuyingNft() update auction record err=", dberr.Error)
				return dberr.Error
			}
			return errors.New("BuyingNft() AuthWormTrans." + err.Error())
		}

		for i := 0; i < WaitTransTime; i++ {
			trans := Trans{}
			result := nft.db.Model(&Trans{}).Where("txhash = ?", txhash).First(&trans)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return result.Error
			}
			if result.Error == gorm.ErrRecordNotFound {
				time.Sleep(time.Second)
				continue
			}
			//time.Sleep(time.Second)
			log.Println("BuyingNft() trans ok ")
			return nil
		}
		auctRec = Auction{}
		auctRec.SellState = SellStateStart.String()
		auctRec.Price = price

		dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
		if dberr.Error != nil {
			log.Println("BuyingNft() update auction record err=", dberr.Error)
			return dberr.Error
		}
		return errors.New("Waiting for the deal to close!")
	} else {
		var bidRec Bidding
		dberr := nft.db.Where("contract = ? AND tokenid = ? AND bidaddr = ?", contractAddr, tokenId, buyerAddr).First(&bidRec)
		if dberr.Error != nil {
			if dberr.Error != gorm.ErrRecordNotFound {
				log.Println("BuyingNft() RecordNotFound")
				return dberr.Error
			}
			return ErrNftNotSell
		}
		//err := WormTrans(nftrecord.Mintstate, sellerSig, bidRec.Tradesig)
		//ExchangerAuth = `{"exchanger_owner":"0x01842a2cf56400a245a56955dc407c2c4137321e","to":"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900","block_number":"0x2540be400","sig":"0x7f1ca96714208959c5a75bdbf4770893b76b13c0bca26da2086c3365e537d57444f79b31498301c5c1d55400eec4b469c83a88a527159112f27ff934c222e4191b"}`
		auctRec := Auction{}
		auctRec.SellState = SellStateWait.String()
		auctRec.Price = price
		auctRec.VoteStage = bidRec.VoteStage
		dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
		if dberr.Error != nil {
			log.Println("BuyingNft() update auction record err=", dberr.Error)
			return dberr.Error
		}
		txhash, err := AuthWormTrans(nftrecord.Mintstate, sellerSig, bidRec.Tradesig, ExchangerAuth)
		if err != nil {
			fmt.Println("BuyingNft() WormTrans err=", err)
			auctRec = Auction{}
			auctRec.SellState = SellStateStart.String()
			auctRec.Price = price
			dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
				auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
			if dberr.Error != nil {
				log.Println("BuyingNft() update auction record err=", dberr.Error)
				return dberr.Error
			}
			return errors.New("BuyingNft() AuthWormTrans." + err.Error())
		}
		for i := 0; i < WaitTransTime; i++ {
			trans := Trans{}
			result := nft.db.Model(&Trans{}).Where("txhash = ?", txhash).First(&trans)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return result.Error
			}
			if result.Error == gorm.ErrRecordNotFound {
				time.Sleep(time.Second)
				continue
			}
			//time.Sleep(time.Second)
			log.Println("BuyingNft() trans ok ")
			return nil
		}
		auctRec = Auction{}
		auctRec.SellState = SellStateStart.String()
		auctRec.Price = price
		dberr = nft.db.Model(&Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
			auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
		if dberr.Error != nil {
			log.Println("BuyingNft() update auction record err=", dberr.Error)
			return dberr.Error
		}
		return errors.New("Waiting for the deal to close!")
	}
	return nil
}

func (nft NftDb) GroupBuyingNft(userAddr, params string) error {

	if params == "" {
		fmt.Println("input param nil")
		return errors.New("input param nil")
	}
	var Buying []BuyingParams
	err := json.Unmarshal([]byte(params), &Buying)
	if err != nil {
		fmt.Println("Unmarshal input err=", err)
		return err
	}
	fmt.Println("buying:   ", Buying)
	for _, j := range Buying {
		price, _ := strconv.ParseUint(j.Price, 10, 64)
		err = nft.BuyingNft(j.UserAddr, j.BuyerSig, j.ContractAddr, j.TokenId, price, j.BuyerAddr, j.VoteStage, j.SellerSig)
		if err != nil {
			fmt.Println("BuyingNft err=", err)
			return err
		}
	}
	return nil
}
