package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type OwnerSnftCollectionCatch struct {
	NftInfo []SnftCollectInfo
	Total   uint64
}

func (nft NftDb) QueryOwnerSnftCollection(owner, Categories, startIndex, count string) ([]SnftCollectInfo, int64, error) {
	spendT := time.Now()
	queryCatchSql := owner + Categories + startIndex + count
	nftCatch := OwnerSnftCollectionCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryOwnerSnftCollection", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryOwnerSnftCollection() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, int64(nftCatch.Total), nil
	}
	stageCollection := []SnftCollectInfo{}
	//sql := `SELECT  a.createaddr, a.name, a.img, a.desc, count(*) as chipcount  FROM collects as a ` +
	//	   `JOIN (SELECT collectcreator, collections FROM nfts whereownaddr and deleted_at  is null ) as b ON a.Createaddr = b.collectcreator AND a.name = b.collections `
	sql := `SELECT a.createaddr, a.name, a.img, a.desc, b.chipcount  FROM collects as a  ` +
		`JOIN (SELECT collectcreator, collections, count(*) as chipcount FROM nfts whereownaddr and deleted_at is null group by collectcreator, collections  ) as b ` +
		`ON a.Createaddr = b.collectcreator AND a.name = b.collections `
	sqlcount := `SELECT a.createaddr, a.name FROM collects as a  ` +
		`JOIN (SELECT collectcreator, collections FROM nfts whereownaddr and deleted_at is null group by collectcreator, collections  ) as b ` +
		`ON a.Createaddr = b.collectcreator AND a.name = b.collections `
	if Categories != "*" {
		sql = sql + " where categories = " + "\"" + Categories + "\""
		sqlcount = sqlcount + " where categories = " + "\"" + Categories + "\""
	}
	sql = strings.Replace(sql, "whereownaddr", "where ownaddr =  "+"\""+owner+"\"  ", -1)
	sqlcount = strings.Replace(sqlcount, "whereownaddr", "where ownaddr =  "+"\""+owner+"\"  ", -1)
	//sql = sql + " GROUP BY a.createaddr, a.name, a.img, a.desc "
	sqlcout := "select count(c.Createaddr) from " + "( " + sqlcount + " ) as c "
	var recount int64
	err := nft.db.Raw(sqlcout).Scan(&recount)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println("QueryOwnerSnftCollection() Scan(&recount) err=", err)
		return nil, 0, ErrDataBase
	}
	fmt.Printf("QueryOwnerSnftCollection() Scan(&recount)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	sql = sql + " limit " + startIndex + ", " + count
	err = nft.db.Raw(sql).Scan(&stageCollection)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println("QueryOwnerSnftCollection() Scan(&stageCollection) err=", err)
		return nil, 0, ErrDataBase
	}
	GetRedisCatch().CatchQueryData("QueryOwnerSnftCollection", queryCatchSql, &OwnerSnftCollectionCatch{stageCollection, uint64(recount)})
	fmt.Printf("QueryOwnerSnftCollection() Scan(&stageCollection) spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return stageCollection, recount, nil
}
