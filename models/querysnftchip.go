package models

import (
	"log"
	"strconv"
)

type SnftChipInfo struct {
	Ownaddr			string		`json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft拥有者地址'"`
	Name			string 		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft名称'"`
	Desc			string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'nft描述'"`
	Url				string		`json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc原始数据保持地址'"`
	Contract		string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Tokenid			string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
	Nftaddr			string 		`json:"nft_address" gorm:"type:char(42) ;comment:'wormholes链唯一标识nft标志'"`
	Count	     	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Categories		string 		`json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft分类'"`
	Collectcreator	string		`json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'合集创建者地址'"`
	Collections 	string  	`json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT合集名'"`
	Hide			string		`json:"hide" gorm:"type:char(20) NOT NULL;comment:'是否让其他人看到'"`
	Createaddr		string		`json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'创建nft地址'"`
	Price			uint64		`json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'创建时定的价格'"`
	Royalty     	int 		`json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Createdate		int64		`json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft创建时间'"`
	Transcnt		int			`json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt		uint64		`json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	Chipcount		int
	//Sellprice		uint64		`json:"sellprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'正在销售的价格'"`
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

