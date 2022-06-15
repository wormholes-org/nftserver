package models

import (
	"log"
	"strconv"
)

func (nft NftDb) QueryStageList(start_Index, count string) ([]string, int64, error) {
	stageList := []string{}
	var recCount int64
	countSql := `select count(*) from (select snftstage from nfts where snftstage != "" GROUP BY snftstage) count`
	err := nft.db.Raw(countSql).Scan(&recCount)
	if err.Error != nil {
		log.Println("QueryStageList() Count(&recCount) err=", err.Error)
		return nil, 0, err.Error
	}
	startIndex, _ := strconv.Atoi(start_Index)
	nftCount, _ := strconv.Atoi(count)
	err = nft.db.Model(Nfts{}).Select("snftstage").Where("snftstage != \"\"").Group("snftstage").Offset(startIndex).Limit(nftCount).Find(&stageList)
	if err.Error != nil {
		log.Println("QueryStageList() Find(&stageList) err=", err.Error)
		return nil, 0, err.Error
	}
	return stageList, recCount, nil
}

