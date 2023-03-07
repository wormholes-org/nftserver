package models

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"log"
	"strings"
	"time"
)

type BuyList struct {
	Snft []struct {
		ContractAddr string `json:"contractAddr"`
		TokenId      string `json:"tokenId"`
	}
	BuyAuthSig string `json:"authSig"`
}

func (nft NftDb) BatchForceBuyingNft(userAddr, buylist string) error {
	userAddr = strings.ToLower(userAddr)
	fmt.Println("BatchForceBuyingNft() userAddr=", userAddr, "time=", time.Now().String())

	if UserSync.LockTran(userAddr) {
		return ErrUserTrading
	} else {
		defer UserSync.UnLockTran(userAddr)
	}
	UserSync.Lock(userAddr)
	defer UserSync.UnLock(userAddr)
	buyList := BuyList{}
	if buylist != "" {
		uerr := json.Unmarshal([]byte(buylist), &buyList)
		if uerr != nil {
			log.Println("BatchForceBuyingNft() input buylist err = ", uerr)
			return ErrDataFormat
		}
	}
	txhashs := []string{}
	blocknumber := contracts.GetCurrentBlockNumber()
	for _, s := range buyList.Snft {
		var nftrecord Nfts
		dberr := nft.db.Where("contract = ? AND tokenid =? and mergelevel = 0", s.ContractAddr, s.TokenId).First(&nftrecord)
		if dberr.Error != nil {
			log.Println("BatchForceBuyingNft() bidprice not find nft err= ", dberr.Error)
			continue
		}
		if nftrecord.Mergetype != nftrecord.Mergelevel {
			log.Println("BatchForceBuyingNft() snft has been merged")
			continue
		}
		buyer := contracts.Buyer{}
		buyer.Nftaddress = nftrecord.Nftaddr
		buyer.Exchanger = ExchangeOwer
		buyer.Seller = nftrecord.Ownaddr
		txhash, err := contracts.ForceBuyingAuthExchangeTrans(buyer, buyList.BuyAuthSig, ExchangerAuth, contracts.SuperAdminAddr)
		if err != nil {
			log.Println("BatchForceBuyingNft() ForceBuyingAuthExchangeTrans err=", err, "txhash=", txhash, "nftaddr=", nftrecord.Nftaddr)
			continue
		}
		txhashs = append(txhashs, txhash)
	}
	if len(txhashs) != 0 {
		GetRedisCatch().SetDirtyFlag(TradingDirtyName)
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
				log.Println("BatchForceBuyingNft() trans ok ")
				return nil
			}
			time.Sleep(time.Second)
		}
	}
	return ErrWaitingClose
}
