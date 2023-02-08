package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LevelSnfts struct {
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft owner address'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft name'"`
	Desc           string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'nft description'"`
	Url            string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc raw data hold address'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Tokenid        string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'Uniquely identifies the nft flag'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) ;comment:'Chain of wormholes uniquely identifies the nft flag'"`
	Count          int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft classification'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'Collection creator address'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT collection name'"`
	Hide           string `json:"hide" gorm:"type:char(20) NOT NULL;comment:'Whether to let others see'"`
	Createaddr     string `json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'Create nft address'"`
	Price          uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Price at creation time'"`
	Royalty        int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	Createdate     int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft creation time'"`
	Transcnt       int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt       uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	Pledgestate    string `json:"pledgestate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'Pledgestate status'"`
	Selltype       string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft transaction type'"`
	Sellprice      uint64 `json:"sellprice" gorm:"type:bigint unsigned DEFAULT 0;comment:'selling price'"`
	Maxbidprice    uint64 `json:"maxbidprice" gorm:"type:bigint unsigned DEFAULT 0;comment:'Highest bid price'"`
	Mergetype      uint8
	Mergelevel     uint8
	Offernum       uint64 `json:"offernum" gorm:"type:bigint unsigned DEFAULT 0;comment:'number of bids'"`
	Exchangecnt    int
}

type OwnerLevelSnfts struct {
	SnftInfo []LevelSnfts
	Total    uint64
}

func (nft NftDb) QueryOwnerLevelSnfts(owner, sellType, snftLevel, Index, count string) ([]LevelSnfts, int64, error) {
	owner = strings.ToLower(owner)
	slevel, _ := strconv.Atoi(snftLevel)
	startIndex, _ := strconv.Atoi(Index)
	nftCount, _ := strconv.Atoi(count)
	spendT := time.Now()
	queryCatchSql := owner + sellType + snftLevel + Index + count
	fmt.Println("QueryOwnerLevelSnfts() queryCatchSql=", queryCatchSql)
	nftCatch := OwnerLevelSnfts{}
	cerr := GetRedisCatch().GetCatchData("QueryOwnerLevelSnfts", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryOwnerLevelSnfts() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.SnftInfo, int64(nftCatch.Total), nil
	}
	var lsnfts []LevelSnfts
	var recCount int64
	var sellwg sync.WaitGroup

	go func() {
		sellwg.Add(1)
		defer sellwg.Done()
		switch sellType {
		case "BidPrice":
			spendT := time.Now()
			err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and offernum != 0 and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel).Count(&recCount)
			if err.Error != nil {
				log.Println("QueryOwnerLevelSnfts() recCount err=", err)
			}
			fmt.Printf("QueryOwnerLevelSnfts() BidPrice count spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		case "NotSale", "FixPrice":
			spendT := time.Now()
			err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and selltype = ? and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel, sellType).Count(&recCount)
			if err.Error != nil {
				log.Println("QueryOwnerLevelSnfts() recCount err=", err)
			}
			fmt.Printf("QueryOwnerLevelSnfts() count spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		case "All":
			err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel).Count(&recCount)
			if err.Error != nil {
				log.Println("QueryOwnerLevelSnfts() recCount err=", err)
			}
		}
	}()

	switch sellType {
	case "BidPrice":
		spendT := time.Now()
		//err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and offernum != 0 and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel).Count(&recCount)
		//if err.Error != nil {
		//	log.Println("QueryOwnerLevelSnfts() recCount err=", err)
		//	return nil, 0, ErrNftNotExist
		//}
		//fmt.Printf("QueryOwnerLevelSnfts() BidPrice count spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		spendT = time.Now()

		err := nft.db.Model(Nfts{}).Where("ownaddr = ? and mergetype = ? and offernum != 0 and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel).Limit(nftCount).Offset(startIndex).Scan(&lsnfts)
		if err.Error != nil {
			log.Println("QueryOwnerLevelSnfts() recCount err=", err)
			return nil, 0, ErrNftNotExist
		}
		//if int64(startIndex) >= recCount || recCount == 0 {
		//	return nil, 0, ErrNotMore
		//}
		fmt.Printf("QueryOwnerLevelSnfts() BidPrice get snft spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	case "NotSale", "FixPrice":
		spendT := time.Now()
		//err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and selltype = ? and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel, sellType).Count(&recCount)
		//if err.Error != nil {
		//	log.Println("QueryOwnerLevelSnfts() recCount err=", err)
		//	return nil, 0, ErrNftNotExist
		//}
		//fmt.Printf("QueryOwnerLevelSnfts() count spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())

		spendT = time.Now()
		//if int64(startIndex) >= recCount || recCount == 0 {
		//	return nil, 0, ErrNotMore
		//}
		err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and selltype = ? and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel, sellType).Limit(nftCount).Offset(startIndex).Scan(&lsnfts)
		if err.Error != nil {
			log.Println("QueryOwnerLevelSnfts() get snft err=", err)
			return nil, 0, ErrNftNotExist
		}

		fmt.Printf("QueryOwnerLevelSnfts() get snft spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	case "All":
		//err := nft.db.Model(Nfts{}).Debug().Where("ownaddr = ? and mergetype = ? and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel).Count(&recCount)
		//if err.Error != nil {
		//	log.Println("QueryOwnerLevelSnfts() recCount err=", err)
		//	return nil, 0, ErrNftNotExist
		//}
		//if int64(startIndex) >= recCount || recCount == 0 {
		//	return nil, 0, ErrNotMore
		//}
		err := nft.db.Model(Nfts{}).Where("ownaddr = ? and mergetype = ? and Pledgestate = \"NoPledge\" and exchange = 0 and ( mergetype = mergelevel)", owner, slevel).Limit(nftCount).Offset(startIndex).Scan(&lsnfts)
		if err.Error != nil {
			log.Println("QueryOwnerLevelSnfts() recCount err=", err)
			return nil, 0, ErrNftNotExist
		}
	default:
		return nil, 0, ErrNftNotExist
	}
	sellwg.Wait()
	if int64(startIndex) >= recCount || recCount == 0 {
		return nil, 0, ErrNotMore
	}
	GetRedisCatch().CatchQueryData("QueryOwnerLevelSnfts", queryCatchSql, &OwnerLevelSnfts{lsnfts, uint64(recCount)})
	fmt.Printf("QueryOwnerLevelSnfts() total spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return lsnfts, recCount, nil
}
