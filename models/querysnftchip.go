package models

import (
	"log"
	"strconv"
)

type SnftChipInfo struct {
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
	Chipcount      int
	//Sellprice		uint64		`json:"sellprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'price being sold'"`
}

func (nft NftDb) QuerySnftChip(contract, tokenid, start_Index, count string) ([]SnftChipInfo, uint64, error) {
	var nftRec Nfts
	err := nft.db.Model(&Nfts{}).Select("snft").Where("contract = ? AND tokenid = ?", contract, tokenid).First(&nftRec)
	if err.Error != nil {
		log.Println("QuerySnftChip() Select(snft) err=", err.Error)
		return nil, 0, err.Error
	}
	var recCount int64
	err = nft.db.Model(Nfts{}).Where("snft = ?", nftRec.Snft).Count(&recCount)
	if err.Error != nil {
		log.Println("QuerySnftChip() Count(&recCount) err=", err.Error)
		return nil, 0, err.Error
	}
	startIndex, _ := strconv.Atoi(start_Index)
	nftCount, _ := strconv.Atoi(count)
	nftInfo := []SnftChipInfo{}
	err = nft.db.Model(Nfts{}).Where("snft = ?", nftRec.Snft).Offset(startIndex).Limit(nftCount).Scan(&nftInfo)
	if err.Error != nil {
		log.Println("QuerySnftChip()Find(&nftInfo) err=", err.Error)
		return nil, 0, err.Error
	}
	return nftInfo, uint64(recCount), nil
}
