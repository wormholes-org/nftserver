package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type RecommendSnftChipInfo struct {
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft owner address'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft name'"`
	Desc           string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'nft description'"`
	Url            string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc raw data hold address'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Tokenid        string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'Uniquely identifies the nft flag'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) ;comment:'Chain of wormholes uniquely identifies the nft flag'"`
	Count          int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft classification'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'Collection creator address'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT collection name'"`
	Hide           string `json:"hide" gorm:"type:char(20) NOT NULL;comment:'Whether to let others see'"`
	Createaddr     string `json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'Create nft address'"`
	Price          uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Price at creation time'"`
	Royalty        int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	Createdate     int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft creation time'"`
	Transcnt       int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt       uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	Pledgestate    string `json:"pledgestate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'Pledgestate status'"`
	Selltype       string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft transaction type'"`
	Sellprice      uint64 `json:"sellprice" gorm:"type:bigint unsigned DEFAULT 0;comment:'selling price'"`
	Mergetype      uint8
	Mergelevel     uint8
	Bidding        bool
	Maxbidprice    uint64 `json:"maxbidprice" gorm:"type:bigint unsigned DEFAULT 0;comment:'Highest bid price'"`
}

func (nft NftDb) QueryRecommendSnftChip(userAddr, contract, tokenid string) ([]*RecommendSnftChipInfo, uint64, error) {
	spendT := time.Now()
	userAddr = strings.ToLower(userAddr)
	contract = strings.ToLower(contract)
	var nftRec Nfts
	err := nft.db.Model(&Nfts{}).Select("snft").Where("contract = ? AND tokenid = ?", contract, tokenid).First(&nftRec)
	if err.Error != nil {
		log.Println("QuerySnftChip() Select(snft) err=", err.Error)
		return nil, 0, ErrNftNotExist
	}
	if nftRec.Snft[len(nftRec.Snft)-1:len(nftRec.Snft)] == "m" {
		//nftRec.Snft = nftRec.Snft[:len(nftRec.Snft)-1]
		return nil, 0, ErrNftNotExist
	}
	var recCount int64
	err = nft.db.Model(Nfts{}).Where("snft = ? and exchange = 0 and pledgestate = 0", nftRec.Snft).Count(&recCount)
	if err.Error != nil {
		log.Println("QueryRecommendSnftChip() Count(&recCount) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	nftInfo := []*RecommendSnftChipInfo{}
	err = nft.db.Model(Nfts{}).Where("snft = ? and exchange = 0 and pledgestate = 0", nftRec.Snft).Scan(&nftInfo)
	if err.Error != nil {
		log.Println("QueryRecommendSnftChip() Find(&nftInfo) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	bidRecs := []Bidding{}
	dberr := nft.db.Model(Bidding{}).Where("bidaddr = ?", userAddr).Scan(&bidRecs)
	if dberr.Error != nil {
		if dberr.Error != gorm.ErrRecordNotFound {
			log.Println("BuyingNft() RecordNotFound")
			return nil, 0, ErrDataBase
		}
	}
	if len(bidRecs) != 0 {
		for _, info := range nftInfo {
			for _, rec := range bidRecs {
				if info.Contract == rec.Contract && info.Tokenid == rec.Tokenid {
					info.Bidding = true
				}
			}
		}
	}

	log.Printf("QueryRecommendSnftChip() spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return nftInfo, uint64(recCount), nil
}

type SnftOtherChipInfo struct {
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft owner address'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft name'"`
	Desc           string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'nft description'"`
	Url            string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc raw data hold address'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Tokenid        string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'Uniquely identifies the nft flag'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) ;comment:'Chain of wormholes uniquely identifies the nft flag'"`
	Count          int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft classification'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'Collection creator address'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT collection name'"`
	Hide           string `json:"hide" gorm:"type:char(20) NOT NULL;comment:'Whether to let others see'"`
	Createaddr     string `json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'Create nft address'"`
	Price          uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Price at creation time'"`
	Royalty        int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	Createdate     int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft creation time'"`
	Transcnt       int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt       uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	Pledgestate    string `json:"pledgestate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'Pledgestate status'"`
	Selltype       string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft transaction type'"`
	Sellprice      uint64 `json:"sellprice" gorm:"type:bigint unsigned DEFAULT 0;comment:'selling price'"`
	Mergetype      uint8
	Mergelevel     uint8
	Bidding        bool
	Maxbidprice    uint64 `json:"maxbidprice" gorm:"type:bigint unsigned DEFAULT 0;comment:'Highest bid price'"`
	Bidowner       string `json:"bidowner"`
	Bigtime        string `json:"bigtime"`
}

func (nft NftDb) QuerySnftOtherChip(userAddr, contract, tokenid string) ([]*SnftOtherChipInfo, uint64, error) {
	spendT := time.Now()
	userAddr = strings.ToLower(userAddr)
	contract = strings.ToLower(contract)
	var nftRec Nfts
	err := nft.db.Model(&Nfts{}).Select("snft").Where("contract = ? AND tokenid = ?", contract, tokenid).First(&nftRec)
	if err.Error != nil {
		log.Println("QuerySnftOtherChip() Select(snft) err=", err.Error)
		return nil, 0, ErrNftNotExist
	}
	if nftRec.Snft[len(nftRec.Snft)-1:len(nftRec.Snft)] == "m" {
		//nftRec.Snft = nftRec.Snft[:len(nftRec.Snft)-1]
		return nil, 0, ErrNftNotExist
	}
	var recCount int64
	err = nft.db.Model(Nfts{}).Where("snft = ? and exchange = 0 and pledgestate = 0", nftRec.Snft).Count(&recCount)
	if err.Error != nil {
		log.Println("QuerySnftOtherChip() Count(&recCount) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	nftInfo := []*SnftOtherChipInfo{}
	err = nft.db.Model(Nfts{}).Where("snft = ? and exchange = 0 and pledgestate = 0", nftRec.Snft).Scan(&nftInfo)
	if err.Error != nil {
		log.Println("QuerySnftOtherChip() Find(&nftInfo) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	bidRecs := Bidding{}
	for _, info := range nftInfo {
		if info.Maxbidprice == 0 {
			continue
		}
		bidRecs = Bidding{}
		dberr := nft.db.Model(Bidding{}).Where("tokenid = ? and price=?", info.Tokenid, info.Maxbidprice).First(&bidRecs)
		if dberr.Error != nil {
			log.Println("QuerySnftOtherChip() first bidding RecordNotFound")
			continue
		}
		info.Bidowner = bidRecs.Bidaddr
		info.Bigtime = fmt.Sprint(bidRecs.UpdatedAt.Unix())

	}

	log.Printf("QuerySnftOtherChip() spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return nftInfo, uint64(recCount), nil
}
