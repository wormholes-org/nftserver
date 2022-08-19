package models

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
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

type SnftChipCatch struct {
	NftInfo []SnftChipInfo
	Total   uint64
}

func (nft NftDb) QuerySnftChip(contract, tokenid, start_Index, count string) ([]SnftChipInfo, uint64, error) {
	spendT := time.Now()
	queryCatchSql := contract + tokenid + start_Index + count
	nftCatch := SnftChipCatch{}
	cerr := GetRedisCatch().GetCatchData("QuerySnftChip", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QuerySnftChip() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, nftCatch.Total, nil
	}
	var nftRec Nfts
	err := nft.db.Model(&Nfts{}).Select("snft").Where("contract = ? AND tokenid = ?", contract, tokenid).First(&nftRec)
	if err.Error != nil {
		log.Println("QuerySnftChip() Select(snft) err=", err.Error)
		return nil, 0, ErrNotFound
	}
	var recCount int64
	err = nft.db.Model(Nfts{}).Where("snft = ?", nftRec.Snft).Count(&recCount)
	if err.Error != nil {
		log.Println("QuerySnftChip() Count(&recCount) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	startIndex, _ := strconv.Atoi(start_Index)
	nftCount, _ := strconv.Atoi(count)
	nftInfo := []SnftChipInfo{}
	err = nft.db.Model(Nfts{}).Where("snft = ?", nftRec.Snft).Offset(startIndex).Limit(nftCount).Scan(&nftInfo)
	if err.Error != nil {
		log.Println("QuerySnftChip() Find(&nftInfo) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	GetRedisCatch().CatchQueryData("QuerySnftChip", queryCatchSql, &SnftChipCatch{nftInfo, uint64(recCount)})
	log.Printf("QuerySnftChip() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return nftInfo, uint64(recCount), nil
}

func (nft NftDb) QueryOwnerSnftChip(owner, start_Index, count string) ([]SnftChipInfo, uint64, error) {
	spendT := time.Now()
	queryCatchSql := owner + start_Index + count
	nftCatch := SnftChipCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryOwnerSnftChip", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryOwnerSnftChip() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, nftCatch.Total, nil
	}
	var recCount int64
	err := nft.db.Model(Nfts{}).Where("ownaddr = ?", owner).Count(&recCount)
	if err.Error != nil {
		log.Println("QueryOwnerSnftChip() Count(&recCount) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	startIndex, _ := strconv.Atoi(start_Index)
	nftCount, _ := strconv.Atoi(count)

	if int64(startIndex) >= recCount || recCount == 0 {
		return nil, 0, ErrNotMore
	}

	nftInfo := []SnftChipInfo{}
	err = nft.db.Model(Nfts{}).Where("ownaddr = ?", owner).Offset(startIndex).Limit(nftCount).Scan(&nftInfo)
	if err.Error != nil {
		log.Println("QueryOwnerSnftChip() Find(&nftInfo) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	GetRedisCatch().CatchQueryData("QueryOwnerSnftChip", queryCatchSql, &SnftChipCatch{nftInfo, uint64(recCount)})
	log.Printf("QueryOwnerSnftChip() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return nftInfo, uint64(recCount), nil
}

func (nft NftDb) QueryArraySnft(array string) ([]SnftChipInfo, error) {
	snftarray := []string{}
	uerr := json.Unmarshal([]byte(array), &snftarray)
	if uerr != nil {
		log.Println("input data err =", uerr)
		return nil, ErrData
	}
	nftInfo := []SnftChipInfo{}
	err := nft.db.Model(Nfts{}).Where("nftaddr  in ?", snftarray).Scan(&nftInfo)
	if err.Error != nil {
		log.Println("QueryArraySnft() Find(&nftInfo) err=", err.Error)
		return nil, ErrDataBase
	}
	return nftInfo, nil
}
