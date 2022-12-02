package models

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type SnftPledgeCollectInfo struct {
	Createaddr  string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Name        string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'collection name'"`
	Img         string `json:"img" gorm:"type:longtext ;comment:'logo image'"`
	Desc        string `json:"desc" gorm:"type:longtext NOT NULL;comment:'Collection description'"`
	Chipcount   int64  `json:"chipcount" gorm:"type:bigint ;comment:'Number of slices'"`
	Totalcount  int64  `json:"totalcount" gorm:"type:bigint ;comment:'Number of slices'"`
	MergeLevel  uint8
	MergeNumber uint32
	BlockNumber uint64
}

type OwnerSnftCollectionCatch struct {
	NftInfo []SnftPledgeCollectInfo
	Total   uint64
}

type SnftPledgeInfo struct {
	Snft       string `json:"snft" gorm:"type:char(42) ;comment:'wormholes chain snft'"`
	Createaddr string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'collection name'"`
	Img        string `json:"img" gorm:"type:longtext ;comment:'logo image'"`
	Desc       string `json:"desc" gorm:"type:longtext NOT NULL;comment:'Collection description'"`
	Chipcount  int64  `json:"chipcount" gorm:"type:bigint ;comment:'Number of slices'"`
	Totalcount int64  `json:"totalcount" gorm:"type:bigint ;comment:'Number of slices'"`
}

type SnftExist struct {
	Createaddr string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'collection name'"`
}

func (nft NftDb) QueryOwnerSnftCollection(owner, Categories, startIndex, count, status string) ([]SnftPledgeCollectInfo, int64, error) {
	spendT := time.Now()
	queryCatchSql := owner + Categories + startIndex + count + status
	nftCatch := OwnerSnftCollectionCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryOwnerSnftCollection", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryOwnerSnftCollection() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, int64(nftCatch.Total), nil
	}
	stageCollection := []SnftPledgeCollectInfo{}
	snftInf := []SnftPledgeInfo{}
	//sql := `SELECT  a.createaddr, a.name, a.img, a.desc, count(*) as chipcount  FROM collects as a ` +
	//	   `JOIN (SELECT collectcreator, collections FROM nfts whereownaddr and deleted_at  is null ) as b ON a.Createaddr = b.collectcreator AND a.name = b.collections `
	var sql, sqlcount string
	var recount int64
	if status == "" {
		sql = `SELECT a.createaddr, a.snftcollection as snft, a.name, a.img, a.desc, b.chipcount, a.totalcount FROM collects as a  ` +
			`JOIN (SELECT collectcreator, collections, count(*) as chipcount FROM nfts whereownaddr and pledgestate= "NoPledge" and exchange = 0 and mergetype = 0 and deleted_at is null group by collectcreator, collections) as b ` +
			`ON a.Createaddr = b.collectcreator AND a.name = b.collections `
		sqlcount = `SELECT a.createaddr, a.name FROM collects as a  ` +
			`JOIN (SELECT collectcreator, collections FROM nfts whereownaddr and pledgestate= "NoPledge" and exchange = 0 and mergetype = 0 and deleted_at is null group by collectcreator, collections  ) as b ` +
			`ON a.Createaddr = b.collectcreator AND a.name = b.collections `
		if Categories != "*" {
			sql = sql + " where categories = " + "\"" + Categories + "\""
			sqlcount = sqlcount + " where categories = " + "\"" + Categories + "\""
		}
		sql = strings.Replace(sql, "whereownaddr", "where ownaddr =  "+"\""+owner+"\"  ", -1)
		sqlcount = strings.Replace(sqlcount, "whereownaddr", "where ownaddr =  "+"\""+owner+"\"  ", -1)
		//sql = sql + " GROUP BY a.createaddr, a.name, a.img, a.desc "
		sqlcout := "select count(c.Createaddr) from " + "( " + sqlcount + " ) as c "
		err := nft.db.Raw(sqlcout).Scan(&recount)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryOwnerSnftCollection() Scan(&recount) err=", err)
			return nil, 0, ErrDataBase
		}
		fmt.Printf("QueryOwnerSnftCollection() Scan(&recount)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		spendT = time.Now()
		sql = sql + " limit " + startIndex + ", " + count
		err = nft.db.Raw(sql).Scan(&snftInf)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryOwnerSnftCollection() Scan(&stageCollection) err=", err)
			return nil, 0, ErrDataBase
		}
	} else {
		sql = `SELECT b.snft, a.snftcollection, a.createaddr, a.name, a.img, a.desc, b.chipcount, a.totalcount FROM collects as a  ` +
			`JOIN (SELECT snft, collectcreator, collections, count(*) as chipcount FROM nfts whereownaddr and exchange = 0 and mergetype =0 and deleted_at is null group by snft, collectcreator, collections  ) as b ` +
			`ON a.Createaddr = b.collectcreator AND a.name = b.collections `
		sqlcount = `SELECT a.createaddr, a.name, b.chipcount, a.snftcollection FROM collects as a  ` +
			`JOIN (SELECT collectcreator, collections, count(*) as chipcount FROM nfts whereownaddr and exchange = 0 and mergetype = 0 and deleted_at is null group by snft, collectcreator, collections  ) as b ` +
			`ON a.Createaddr = b.collectcreator AND a.name = b.collections `
		if Categories != "*" {
			sql = sql + " where categories = " + "\"" + Categories + "\"" + " and b.chipcount = 16 "
			sqlcount = sqlcount + " where categories = " + "\"" + Categories + "\"" + " and b.chipcount = 16 "
		} else {
			sql = sql + " where b.chipcount = 16 "
			sqlcount = sqlcount + " where b.chipcount = 16 "
		}
		sql = strings.Replace(sql, "whereownaddr", "where ownaddr =  "+"\""+owner+"\"  and pledgestate= "+"\""+status+"\" ", -1)
		sqlcount = strings.Replace(sqlcount, "whereownaddr", "where ownaddr =  "+"\""+owner+"\"  and pledgestate= "+"\""+status+"\" ", -1)
		//sqlcout := "select count(c.Createaddr) from " + "( " + sqlcount + " ) as c "

		sqlcout := "select createaddr,name from ( " + sqlcount + ") as c group by createaddr ,name, snftcollection"
		sqlcout = "select count(*) from (" + sqlcout + ") as d"
		err := nft.db.Raw(sqlcout).Scan(&recount)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryOwnerSnftCollection() Scan(&recount) err=", err)
			return nil, 0, ErrDataBase
		}
		fmt.Printf("QueryOwnerSnftCollection() Scan(&recount)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		spendT = time.Now()
		sql = "select createaddr ,name ,img, snftcollection as snft, sum(chipcount) as chipcount, totalcount from ( " + sql + " ) as c group by createaddr, name, img, chipcount, totalcount, snftcollection"
		sql = sql + " limit " + startIndex + ", " + count

		err = nft.db.Raw(sql).Scan(&snftInf)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryOwnerSnftCollection() Scan(&stageCollection) err=", err)
			return nil, 0, ErrDataBase
		}
	}
	exist := make(map[SnftExist]struct{})
	for _, info := range snftInf {
		m := SnftExist{info.Createaddr, info.Name}
		if _, ok := exist[m]; !ok {
			exist[m] = struct{}{}
			tmp := SnftPledgeCollectInfo{}
			tmp.Desc = info.Desc
			tmp.Img = info.Img
			tmp.Createaddr = info.Createaddr
			tmp.Chipcount = info.Chipcount
			tmp.Totalcount = info.Totalcount
			tmp.Name = info.Name

			addr := common.HexToAddress(info.Snft + "00")
			accountInfo, err := contracts.GetAccountInfo(addr, nil)
			if err != nil {
				log.Println("QueryOwnerSnftCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
				return nil, 0, ErrBlockchain
			}
			if accountInfo.Owner.String() == ZeroAddr {
				addr = common.HexToAddress(info.Snft[:len(info.Snft)-1] + "000")
				accountInfo, err = contracts.GetAccountInfo(addr, nil)
				if err != nil {
					log.Println("QueryOwnerSnftCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
					return nil, 0, ErrBlockchain
				}
			}
			tmp.MergeLevel = accountInfo.MergeLevel
			tmp.MergeNumber = accountInfo.MergeNumber
			tmp.BlockNumber = accountInfo.NFTPledgedBlockNumber
			stageCollection = append(stageCollection, tmp)
		}
	}
	GetRedisCatch().CatchQueryData("QueryOwnerSnftCollection", queryCatchSql, &OwnerSnftCollectionCatch{stageCollection, uint64(recount)})
	fmt.Printf("QueryOwnerSnftCollection() Scan(&stageCollection) spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return stageCollection, recount, nil
}
