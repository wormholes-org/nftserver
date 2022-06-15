package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)


func (nft NftDb) QueryOwnerSnftCollection(owner, Categories, startIndex, count string) ([]SnftCollectInfo, int64, error) {
	spendT := time.Now()
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
	sql = strings.Replace(sql, "whereownaddr", "where ownaddr =  " + "\"" + owner + "\"  ", -1 )
	sqlcount = strings.Replace(sqlcount, "whereownaddr", "where ownaddr =  " + "\"" + owner + "\"  ", -1 )
	//sql = sql + " GROUP BY a.createaddr, a.name, a.img, a.desc "
	sqlcout := "select count(c.Createaddr) from " + "( " + sqlcount + " ) as c "
	var recount int64
	err := nft.db.Raw(sqlcout).Scan(&recount)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound{
		log.Println("QueryOwnerSnftCollection() Scan(&recount) err=", err)
		return nil, 0, err.Error
	}
	fmt.Printf("QueryOwnerSnftCollection() Scan(&recount)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	sql = sql + " limit " + startIndex + ", " + count
	err = nft.db.Raw(sql).Scan(&stageCollection)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound{
		log.Println("QueryOwnerSnftCollection() Scan(&stageCollection) err=", err)
		return nil, 0, err.Error
	}
	fmt.Printf("QueryOwnerSnftCollection() Scan(&stageCollection) spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return stageCollection, recount, nil
}

