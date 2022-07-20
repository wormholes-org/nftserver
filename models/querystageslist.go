package models

import (
	"log"
	"strconv"
	"time"
)

type StageListCatch struct {
	Stage []string
	Total int64
}

func (nft NftDb) QueryStageList(start_Index, count string) ([]string, int64, error) {
	spendT := time.Now()
	queryCatchSql := start_Index + count
	stagelist := StageListCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryStageList", queryCatchSql, &stagelist)
	if cerr == nil {
		log.Printf("QueryStageList() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return stagelist.Stage, stagelist.Total, nil
	}
	stageList := []string{}
	var recCount int64
	countSql := `select count(*) from (select snftstage from nfts where snftstage != "" GROUP BY snftstage) count`
	err := nft.db.Raw(countSql).Scan(&recCount)
	if err.Error != nil {
		log.Println("QueryStageList() Count(&recCount) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	startIndex, _ := strconv.Atoi(start_Index)
	nftCount, _ := strconv.Atoi(count)
	err = nft.db.Model(Nfts{}).Select("snftstage").Where("snftstage != \"\"").Group("snftstage").Offset(startIndex).Limit(nftCount).Find(&stageList)
	if err.Error != nil {
		log.Println("QueryStageList() Find(&stageList) err=", err.Error)
		return nil, 0, ErrDataBase
	}
	GetRedisCatch().CatchQueryData("QueryStageList", queryCatchSql, &StageListCatch{stageList, recCount})
	log.Printf("QueryStageList() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return stageList, recCount, nil
}
