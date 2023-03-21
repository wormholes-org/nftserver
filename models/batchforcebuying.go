package models

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"log"
	"strings"
	"sync"
	"time"
)

type SnftSyncMapList struct {
	Mux   sync.Mutex
	Snfts map[string]struct{}
}

func (u *SnftSyncMapList) LockSnft(snftAddr string) bool {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	if len(u.Snfts) == 0 {
		u.Snfts = make(map[string]struct{})
	}
	_, ok := u.Snfts[snftAddr]
	if !ok {
		u.Snfts[snftAddr] = struct{}{}
		return false
	} else {
		return true
	}
}

func (u *SnftSyncMapList) UnLockSnft(snftAddr string) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	_, ok := u.Snfts[snftAddr]
	if ok {
		delete(u.Snfts, snftAddr)
	}
}

type BuyList struct {
	Snft []struct {
		ContractAddr string `json:"contractAddr"`
		TokenId      string `json:"tokenId"`
	}
	BuyAuthSig string `json:"authSig"`
}

var snftLock SnftSyncMapList

func (nft NftDb) BatchForceBuyingNft(userAddr, buylist string) error {
	userAddr = strings.ToLower(userAddr)
	fmt.Println("BatchForceBuyingNft() userAddr=", userAddr, "time=", time.Now().String())

	buyList := BuyList{}
	if buylist != "" {
		uerr := json.Unmarshal([]byte(buylist), &buyList)
		if uerr != nil {
			log.Println("BatchForceBuyingNft() input buylist err = ", uerr)
			return ErrDataFormat
		}
	} else {
		log.Println("BatchForceBuyingNft() buyList.Snft = 0 ")
		return ErrDataFormat
	}
	var snft string
	var nftrecord Nfts
	if len(buyList.Snft) > 0 {
		dberr := nft.db.Where("contract = ? AND tokenid =? and mergelevel = 0", buyList.Snft[0].ContractAddr, buyList.Snft[0].TokenId).First(&nftrecord)
		if dberr.Error != nil {
			log.Println("BatchForceBuyingNft() bidprice not find nft err= ", dberr.Error)
			return ErrDataBase
		}
		snft = nftrecord.Snft
	} else {
		log.Println("BatchForceBuyingNft() buyList.Snft = 0 ")
		return ErrDataFormat
	}
	if snftLock.LockSnft(snft) {
		return ErrUserTrading
	} else {
		defer snftLock.UnLockSnft(snft)
	}

	//txhashs := []string{}
	//blocknumber := contracts.GetCurrentBlockNumber()
	//for _, s := range buyList.Snft {
	//	var nftrecord Nfts
	//	dberr := nft.db.Where("contract = ? AND tokenid =? and mergelevel = 0", s.ContractAddr, s.TokenId).First(&nftrecord)
	//	if dberr.Error != nil {
	//		log.Println("BatchForceBuyingNft() bidprice not find nft err= ", dberr.Error)
	//		continue
	//	}
	//	if nftrecord.Mergetype != nftrecord.Mergelevel {
	//		log.Println("BatchForceBuyingNft() snft has been merged")
	//		continue
	//	}
	//	buyer := contracts.Buyer{}
	//	buyer.Nftaddress = nftrecord.Nftaddr
	//	buyer.Exchanger = ExchangeOwer
	//	buyer.Seller = nftrecord.Ownaddr
	//	txhash, blockn, err := contracts.ForceBuyingAuthExchangeTrans(buyer, buyList.BuyAuthSig, ExchangerAuth, contracts.SuperAdminAddr)
	//	if err != nil {
	//		log.Println("BatchForceBuyingNft() ForceBuyingAuthExchangeTrans err=", err, "txhash=", txhash, "nftaddr=", nftrecord.Nftaddr, "blocknumber=", blockn)
	//		continue
	//	}
	//	txhashs = append(txhashs, txhash)
	//}
	if ExchangerAuth == "" {
		fmt.Println("BatchForceBuyingNft() Unauthorized exchange error.")
		return ErrUnauthExchange
	}
	if nftrecord.Mergetype != nftrecord.Mergelevel {
		log.Println("BatchForceBuyingNft() snft has been merged")
		return ErrDataFormat
	}
	buyer := contracts.Buyer{}
	buyer.Nftaddress = nftrecord.Snft
	buyer.Exchanger = ExchangeOwer
	buyer.Seller = ExchangeOwer
	txhash, blockn, err := contracts.ForceBuyingAuthExchangeTrans(buyer, buyList.BuyAuthSig, ExchangerAuth, contracts.SuperAdminAddr)
	if err != nil {
		log.Println("BatchForceBuyingNft() ForceBuyingAuthExchangeTrans err=", err, "txhash=", txhash, "nftaddr=", nftrecord.Nftaddr, "blocknumber=", blockn)
		return err
	}
	blocknumber := contracts.GetCurrentBlockNumber()
	for i := blocknumber + 10; blocknumber < i; {
		blocknumber = contracts.GetCurrentBlockNumber()
		txStatus, err := contracts.GetTransStatus(txhash)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if txStatus {
			GetRedisCatch().SetDirtyFlag(TradingDirtyName)
			log.Println("BatchForceBuyingNft() trans ok ")
			return nil
		}
		time.Sleep(time.Second)
	}
	return ErrWaitingClose
}
