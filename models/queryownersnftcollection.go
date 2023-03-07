package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
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
		sql = `SELECT a.createaddr, a.snftcollection as snft, a.name, a.img, a.desc, b.chipcount, a.totalcount FROM collects as a  
			JOIN (select collections,collectcreator,count(*) as chipcount from nfts where collections in 
			(select collections from nfts where exchange = 0 and ownaddr =?  and mergetype = mergelevel group by collections) 
			and ownaddr = ? and  mergetype =0 and exchange = 0 group  by collections,collectcreator) as b 
			ON a.Createaddr = b.collectcreator AND a.name = b.collections `
		sqlcount = `select collections as chipcount from nfts where collections in 
		(select collections from nfts where exchange = 0 and ownaddr =? and mergetype = mergelevel group by collections) and 
		ownaddr = ? and  mergetype =0 and exchange = 0 group  by collections,collectcreator`

		var recountAddr []string
		err := nft.db.Raw(sqlcount, owner, owner).Scan(&recountAddr)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryOwnerSnftCollection() Scan(&recount) err=", err)
			return nil, 0, ErrDataBase
		}
		recount = int64(len(recountAddr))
		fmt.Printf("QueryOwnerSnftCollection() Scan(&recount)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		spendT = time.Now()
		sql = sql + " limit " + startIndex + ", " + count
		err = nft.db.Raw(sql, owner, owner).Scan(&snftInf)
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

			//addr := common.HexToAddress(info.Snft + "00")
			//accountInfo, err := contracts.GetAccountInfo(addr, nil)
			//if err != nil {
			//	log.Println("QueryOwnerSnftCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
			//	return nil, 0, ErrBlockchain
			//}
			//if accountInfo.Owner.String() == ZeroAddr {
			//	addr = common.HexToAddress(info.Snft[:len(info.Snft)-1] + "000")
			//	accountInfo, err = contracts.GetAccountInfo(addr, nil)
			//	if err != nil {
			//		log.Println("QueryOwnerSnftCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
			//		return nil, 0, ErrBlockchain
			//	}
			//}
			addr := info.Snft + "mm"
			var accountInfo Nfts
			err := nft.db.Model(&Nfts{}).Where("nftaddr = ?", addr).First(&accountInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryOwnerSnftCollection() First nfts err=", err)
				return nil, 0, ErrDataBase
			}
			tmp.MergeLevel = accountInfo.Mergelevel
			tmp.MergeNumber = accountInfo.Mergenumber
			//tmp.BlockNumber = accountInfo.NFTPledgedBlockNumber
			stageCollection = append(stageCollection, tmp)
			//if strings.ToLower(accountInfo.Ownaddr) == strings.ToLower(owner) || strings.ToLower(accountInfo.Ownaddr) == strings.ToLower(ZeroAddr) {
			//	if accountInfo.Exchange != 0 {
			//		continue
			//	}
			//	tmp.MergeLevel = accountInfo.Mergelevel
			//	tmp.MergeNumber = accountInfo.Mergenumber
			//	//tmp.BlockNumber = accountInfo.NFTPledgedBlockNumber
			//	stageCollection = append(stageCollection, tmp)
			//}

		}
	}
	GetRedisCatch().CatchQueryData("QueryOwnerSnftCollection", queryCatchSql, &OwnerSnftCollectionCatch{stageCollection, uint64(recount)})
	fmt.Printf("QueryOwnerSnftCollection() Scan(&stageCollection) spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return stageCollection, recount, nil
}

type OwnerPahseChip struct {
	OwnerCollections []OwnerCollections
}

type OwnerCollections struct {
	Collect   Nfts
	OwnerSnft []OwnerSnft
	status    int
}
type OwnerSnft struct {
	Snft      Nfts
	OwnerChip []OwnerChip
	status    int
}
type OwnerChip struct {
	Chip   Nfts
	status int
}

type OwnerSnftChipData struct {
	Createaddr     string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'collection name'"`
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft owner address'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Tokenid        string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'Uniquely identifies the nft flag'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) ;comment:'Chain of wormholes uniquely identifies the nft flag'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'Collection creator address'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT collection name'"`
	Chipcount      int64  `json:"chipcount" gorm:"type:bigint ;comment:'Number of slices'"`
	Totalcount     int64  `json:"totalcount" gorm:"type:bigint ;comment:'Number of slices'"`
}

func (nft NftDb) QueryOwnerPhaseChip(owner, start_index string) ([]OwnerCollections, int64, error) {
	var ownerdata OwnerPahseChip
	sql := `select left(nftaddr,39) from nfts where ownaddr = ? and mergetype=mergelevel  group by left(nftaddr,39)`
	var phase []string
	db := nft.db.Raw(sql, owner).Scan(&phase)
	if db.Error != nil {
		log.Println("QueryOwnerPhaseChip find phase err=", db.Error)
		return nil, 0, db.Error
	}
	recCount := int64(len(phase))
	startIndex, _ := strconv.Atoi(start_index)
	if int64(startIndex) >= recCount || recCount == 0 {
		return nil, 0, ErrNotMore
	} else {
		ownerphase := phase[startIndex]
		sql = `select left(nftaddr,40) from nfts where ownaddr = ? and nftaddr like ? and mergetype=mergelevel group by left(nftaddr,40)`
		var collect []string
		db = nft.db.Raw(sql, owner, ownerphase+"%").Scan(&collect)
		if db.Error != nil {
			log.Println("QueryOwnerPhaseChip find collect err=", db.Error)
			return nil, 0, db.Error
		}
		var collectnft Nfts
		var ownercollectdata OwnerCollections
		var nfts []Nfts
		for _, singe := range collect {
			collectnft = Nfts{}
			nfts = []Nfts{}
			ownercollectdata = OwnerCollections{}
			if strings.Index(singe, "m") >= 0 {
				continue
			}
			db = nft.db.Model(&Nfts{}).Where("nftaddr = ?", singe+"mm").First(&collectnft)
			if db.Error != nil {
				log.Println("QueryOwnerPhaseChip first singe collect err=", db.Error)
				return nil, 0, db.Error
			}
			ownercollectdata.Collect = collectnft
			if collectnft.Exchange == 0 {
				if collectnft.Mergelevel != 0 {
					ownercollectdata.status = 2
				} else {
					ownercollectdata.status = 0
				}
			} else {
				ownercollectdata.status = 1
			}
			db = nft.db.Model(&Nfts{}).Where("nftaddr like ?", singe+"%").Find(&nfts)
			if db.Error != nil {
				log.Println("QueryOwnerPhaseChip find singe nfts err=", db.Error)
				return nil, 0, db.Error
			}
			snftmap := make(map[string][]Nfts)
			var ownersnftdata OwnerSnft
			for _, singenft := range nfts {
				switch strings.Index(singenft.Nftaddr, "m") {
				case 41:
					ownersnftdata.Snft = singenft
					if singenft.Exchange == 0 {
						if singenft.Mergelevel != 0 {
							ownersnftdata.status = 2
						} else {
							ownersnftdata.status = 0
						}
					} else {
						ownersnftdata.status = 1
					}
					ownercollectdata.OwnerSnft = append(ownercollectdata.OwnerSnft, ownersnftdata)
				default:
					snftmap[singenft.Nftaddr[:41]] = append(snftmap[singenft.Nftaddr[:41]], singenft)
				}

			}
			var ownerchipdata []OwnerChip
			var singeownerchipdata OwnerChip
			for i, singesnft := range ownercollectdata.OwnerSnft {
				ownerchipdata = []OwnerChip{}
				singeownerchipdata = OwnerChip{}
				for _, singeownersnft := range snftmap[singesnft.Snft.Nftaddr[:41]] {
					singeownerchipdata.Chip = singeownersnft
					if singeownersnft.Exchange == 0 {
						if singeownersnft.Mergelevel != 0 {
							singeownerchipdata.status = 2
						} else {
							singeownerchipdata.status = 0
						}
					} else {
						singeownerchipdata.status = 1
					}
					ownerchipdata = append(ownerchipdata, singeownerchipdata)
				}
				ownercollectdata.OwnerSnft[i].OwnerChip = ownerchipdata
			}
			ownerdata.OwnerCollections = append(ownerdata.OwnerCollections, ownercollectdata)

		}

	}
	return ownerdata.OwnerCollections, recCount, nil
}
