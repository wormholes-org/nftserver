package models

import (
	"gorm.io/gorm"
	"log"
)

type SnftCollectInfo struct {
	Createaddr   string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'创建者者地址'"`
	Name         string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'集合名称'"`
	Img          string `json:"img" gorm:"type:longtext ;comment:'logo图片'"`
	Desc         string	`json:"desc" gorm:"type:longtext NOT NULL;comment:'集合描述'"`
	Chipcount    int64  `json:"chipcount" gorm:"type:bigint ;comment:'切片数'"`
}

func (nft NftDb) QueryStageCollection(stage string) ([]SnftCollectInfo, error) {
	stageCollection := []SnftCollectInfo{}
	err := nft.db.Model(&Collects{}).Select([]string{"createaddr", "name", "img"}).Where("snftstage = ?", stage).Find(&stageCollection)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryStageCollection() Select(snft) err=", err.Error)
			return nil, err.Error
		}
		return []SnftCollectInfo{}, nil
	}
	return stageCollection, nil
}

