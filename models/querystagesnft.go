package models

import (
	"log"
)

func (nft NftDb) QueryStageSnft(stage, collect string) ([]SnftChipInfo, error) {
	snftInfo := []SnftChipInfo{}
	snftAddrs := []string{}
	err := nft.db.Model(&Nfts{}).Select("min(nftaddr)").Where("Collections = ? and Snftstage = ?", collect, stage).Group("snft").Find(&snftAddrs)
	if err.Error != nil {
		log.Println("QueryStageSnft() Select(snft) err=", err.Error)
		return nil, err.Error
	}
	err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
	if err.Error != nil {
		log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
		return nil, err.Error
	}
	return snftInfo, nil
}

