package models

import (
	"log"
	"time"
)

func (nft NftDb) QueryStageSnft(stage, collect string) ([]SnftChipInfo, error) {
	spendT := time.Now()
	queryCatchSql := stage + collect
	nftCatch := SnftChipCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryStageSnft", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryStageSnft() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, nil
	}
	snftInfo := []SnftChipInfo{}
	snftAddrs := []string{}
	err := nft.db.Model(&Nfts{}).Select("min(nftaddr)").Where("Collections = ? and Snftstage = ?", collect, stage).Group("snft").Find(&snftAddrs)
	if err.Error != nil {
		log.Println("QueryStageSnft() Select(snft) err=", err.Error)
		return nil, ErrDataBase
	}
	err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
	if err.Error != nil {
		log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
		return nil, ErrDataBase
	}
	GetRedisCatch().CatchQueryData("QueryStageSnft", queryCatchSql, &SnftChipCatch{snftInfo, uint64(0)})
	log.Printf("QueryStageSnft() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return snftInfo, nil
}
