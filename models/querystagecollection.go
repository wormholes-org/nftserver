package models

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type SnftCollectInfo struct {
	Createaddr string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'collection name'"`
	Img        string `json:"img" gorm:"type:longtext ;comment:'logo image'"`
	Desc       string `json:"desc" gorm:"type:longtext NOT NULL;comment:'Collection description'"`
	Chipcount  int64  `json:"chipcount" gorm:"type:bigint ;comment:'Number of slices'"`
}

type StageCollectionCatch struct {
	NftInfo []SnftCollectInfo
}

func (nft NftDb) QueryStageCollection(stage string) ([]SnftCollectInfo, error) {
	spendT := time.Now()
	queryCatchSql := stage
	nftCatch := StageCollectionCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryStageCollection", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryStageCollection() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, nil
	}
	stageCollection := []SnftCollectInfo{}
	err := nft.db.Model(&Collects{}).Select([]string{"createaddr", "name", "img"}).Where("snftstage = ?", stage).Find(&stageCollection)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryStageCollection() Select(snft) err=", err.Error)
			return nil, ErrDataBase
		}
		return []SnftCollectInfo{}, nil
	}
	GetRedisCatch().CatchQueryData("QueryStageCollection", queryCatchSql, &StageCollectionCatch{stageCollection})
	log.Printf("QueryStageCollection() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return stageCollection, nil
}
