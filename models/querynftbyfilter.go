package models

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NftInfo struct {
	Ownaddr string `json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft owner address'"`
	//Md5				string		`json:"md5" gorm:"type:longtext NOT NULL;comment:'Picture md5 value'"`
	Name     string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft classification'"`
	Desc     string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'Review description: Failed review description'"`
	Meta     string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'meta information'"`
	Url      string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc raw data hold address'"`
	Contract string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Tokenid  string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'nft token id'"`
	Snft     string `json:"snft" gorm:"type:char(42) ;comment:'wormholes chain snft'"`
	Nftaddr  string `json:"nft_address" gorm:"type:char(42) ;comment:'wormholes chain nft address'"`
	Count    int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	//Approve			string		`json:"approve" gorm:"type:longtext NOT NULL;comment:'Authorize'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft classification'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'Collection creator address'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT collection name'"`
	//Image			string		`json:"asset_sample" gorm:"type:longtext NOT NULL;comment:'Thumbnail binary data'"`
	Hide string `json:"hide" gorm:"type:char(20) NOT NULL;comment:'Whether to let others see'"`
	//Signdata		string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'Signature data, generated when created'"`
	Createaddr string `json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'Create nft address'"`
	//Verifyaddr		string		`json:"vrf_addr" gorm:"type:char(42) NOT NULL;comment:'Validator address'"`
	Currency string `json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'Transaction currency'"`
	Price    uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Price at creation time'"`
	Royalty  int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	//Paychan    		string 		`json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'trading channel'"`
	//TransCur    	string 		`json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'Transaction currency'"`
	Transprice   uint64 `json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'transaction price'"`
	Transtime    int64  `json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'Last trading time'"`
	Createdate   int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft creation time'"`
	Favorited    int    `json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'Follow count'"`
	Transcnt     int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt     uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	Verified     string `json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'Whether the nft work has passed the review'"`
	Verifiedtime int64  `json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'Review time'"`
	Selltype     string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft transaction type'"`
	Sellprice    uint64 `json:"sellprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'selling price'"`
	Mintstate    string `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'minting status'"`
	Pledgestate  string `json:"pledgestate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'Pledgestate status'"`
	//Extend			string		`json:"extend" gorm:"type:longtext NOT NULL;comment:'extend data'"`
	Offernum    uint64 `json:"offernum" gorm:"type:bigint unsigned DEFAULT NULL;comment:'number of bids'"`
	Maxbidprice uint64 `json:"maxbidprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Highest bid price'"`
}

type SnftAddr struct {
	Snft string `json:"snft" gorm:"type:char(42) ;comment:'wormholes chain snft'"`
}

type NftFilter struct {
	NftInfos []NftInfo
	Total    uint64
}

type NftFilterCatch struct {
	Mux          sync.Mutex
	NftFlush     bool
	AuctionFlush bool
	TimeTag      time.Time
	NftInfo      map[string]*NftFilter
}

func (n *NftFilterCatch) GetByHash(hash string, nftType NftFlushType) *NftFilter {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	fmt.Println("GetByHash() GetByHash n.NftInfo catch len=", len(n.NftInfo))
	if len(n.NftInfo) == 0 {
		n.NftInfo = make(map[string]*NftFilter)
	}
	if nftinfo := n.NftInfo[hash]; nftinfo != nil {
		/*if nftType == NftFlushTypeNewNft {
			if n.NftFlush {
				n.NftFlush = false
				return nil
			}
		}
		if nftType == NftFlushTypeAuction {
			if n.AuctionFlush {
				n.AuctionFlush = false
				return nil
			}
		}*/
		fmt.Println("GetByHash() NftFilterCatch hash=", hash)
		return nftinfo
	}
	return nil
}

func (n *NftFilterCatch) SetByHash(hash string, nftinfo *NftFilter) *NftFilter {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if len(n.NftInfo) == 0 {
		fmt.Println("SetByHash() NftFilterCatch len ==0 ")
		n.NftInfo = make(map[string]*NftFilter)
	}
	n.NftInfo[hash] = nftinfo
	n.TimeTag = time.Now()
	fmt.Println("SetByHash() NftFilterCatch", "len=", len(n.NftInfo), " hash=", hash)
	return nil
}

func (n *NftFilterCatch) SetFlushFlag( /*flag NftFlushType*/ ) {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	/*switch flag {
	case NftFlushTypeNewNft:
		n.NftFlush = true
	case NftFlushTypeAuction:
		n.AuctionFlush = true
	}*/
	fmt.Println("SetFlushFlag() clear catch  hash")
	n.NftInfo = make(map[string]*NftFilter)
}

func (n NftFilterCatch) NftCatchHash(data string) string {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return hexutil.Encode(hash)
}

var NftCatch NftFilterCatch

type NftFlushType int

const (
	NftFlushTypeNewNft NftFlushType = iota
	NftFlushTypeAuction
)

func (this NftFlushType) String() string {
	switch this {
	case NftFlushTypeNewNft:
		return "newnft"
	case NftFlushTypeAuction:
		return "newAuction"
	default:
		return "Unknow"
	}
}

func (nft NftDb) QueryNftByFilter(filter []StQueryField, sort []StSortField,
	startIndex string, count string) ([]NftInfo, uint64, error) {
	var queryWhere string
	var orderBy string
	var totalCount int64
	nftInfo := []NftInfo{}

	sql := "select * from (" + "SELECT nfts.*, auctionstemp.startprice AS sellprice, bidcount.offernum, bidcount.maxbidprice " +
		"FROM nfts LEFT JOIN (select * from auctions WHERE deleted_at IS NULL ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid  LEFT JOIN " +
		"(SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid) bidcount " +
		"ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid " + ") a "
	countSql := "select count(*) from (" + "SELECT nfts.*, auctionstemp.startprice AS sellprice, bidcount.offernum, bidcount.maxbidprice " +
		"FROM nfts LEFT JOIN (select * from auctions WHERE deleted_at IS NULL ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid  LEFT JOIN " +
		"(SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid) bidcount " +
		"ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid " + ") a "
	snftSql := `select min(nftaddr) from nfts where snft != "" group by snft `
	snftSqlct := `select sum(total)/256 from (select count(snft) as total from nfts where snft != "" group by snft) as b `

	snftAddrs := []string{}
	snftInfo := []NftInfo{}
	snftTotalCount := 0
	if len(filter) > 0 {
		queryWhere = nft.joinFilters(filter)
		if len(queryWhere) > 0 {
			sql = sql + " where deleted_at is null and " + queryWhere
			countSql = countSql + " where deleted_at is null and " + queryWhere
		} else {
			sql = sql + " where deleted_at is null "
			countSql = countSql + " where deleted_at is null "
		}
	} else {
		sql = sql + " where deleted_at is null " + "and snft =\"\""
		countSql = countSql + " where deleted_at is null " + "and snft =\"\""
		scount, _ := strconv.Atoi(count)
		scount = scount / 2
		count = strconv.Itoa(scount)
		if len(startIndex) > 0 {
			snftSql = snftSql + " limit " + startIndex + ", " + count
		}
		err := nft.db.Raw(snftSql).Scan(&snftAddrs)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Raw(snftSql).Scan(&snftAddrs) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		snftT := 0.1
		err = nft.db.Raw(snftSqlct).Scan(&snftT)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Raw(snftSqlct).Scan(&snftAddrs) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		snftTotalCount = int(snftT)
		fmt.Println("QueryNftByFilter() snftTotalCount=", snftTotalCount)
		err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
			return nil, uint64(0), ErrDataBase
		}
	}
	if len(sort) > 0 {
		for k, v := range sort {
			if k > 0 {
				orderBy = orderBy + ", "
			}
			orderBy += v.By + " " + v.Order
		}
	}
	if len(orderBy) > 0 {
		orderBy = orderBy + ", id desc"
	} else {
		orderBy = "createdate desc, id desc"
	}
	sql = sql + " order by " + orderBy
	countSql = countSql + " order by " + orderBy

	if len(startIndex) > 0 {
		sql = sql + " limit " + startIndex + ", " + count
	}
	fmt.Println("QueryNftByFilter() sql=", sql)
	fmt.Println("QueryNftByFilter() countSql=", countSql)
	err := nft.db.Raw(sql).Scan(&nftInfo)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println("QueryNftByFilter() Scan(&nftInfo) err=", err)
		return nil, uint64(0), ErrDataBase
	}
	err = nft.db.Raw(countSql).Scan(&totalCount)
	if err.Error != nil {
		log.Println("QueryNftByFilter() Scan(&totalCount) err=", err)
		return nil, uint64(0), ErrDataBase
	}
	/*for k, _ := range nftInfo {
		nftInfo[k].Image = ""
		nftInfo[k].Snft = ""
	}*/
	nftInfo = append(nftInfo, snftInfo...)
	return nftInfo, uint64(totalCount) + uint64(snftTotalCount), nil
}

func QueryWhereSplit(query string) map[string]string {
	spitStr := make(map[string]string)
	for {
		s := strings.Index(query, "(")
		if s == -1 {
			break
		}
		e := strings.Index(query, ")")
		if s == -1 {
			break
		}
		str := query[s : e+1]
		query = query[e+1:]
		if strings.Contains(str, "selltype") {
			spitStr["selltype"] = str
		}
		if strings.Contains(str, "offernum") {
			spitStr["offernum"] = str
		}
		if strings.Contains(str, "sellprice") {
			str = strings.Replace(str, "sellprice", "startprice", -1)
			spitStr["sellprice"] = str
		}
		if strings.Contains(str, "createdate") {
			spitStr["createdate"] = str
			i1 := strings.Index(str, "=")
			i2 := strings.Index(str, ")")
			date := str[i1+1 : i2]
			date = date[:len(date)-4] + "0000"
			spitStr["createdate"] = "(createdate>=" + date + ")"
		}
		if strings.Contains(str, "collectcreator") {
			spitStr["collectcreator"] = str
			str = strings.Replace(str, "collectcreator", "createaddr", -1)
			spitStr["createaddr"] = str

		}
		if strings.Contains(str, "collections") {
			spitStr["collections"] = str
			str = strings.Replace(str, "collections", "name", -1)
			spitStr["name"] = str
		}
		if strings.Contains(str, "categories") {
			spitStr["categories"] = str
		}
	}
	return spitStr
}

func (nft NftDb) QueryNftByFilterNew(filter []StQueryField, sort []StSortField, nftType,
	startIndex string, count string) ([]NftInfo, uint64, error) {
	var queryWhere string
	var orderBy string
	var totalCount int64
	nftInfo := []NftInfo{}
	spendT := time.Now()
	snftAddrs := []string{}
	snftTotalCount := 0

	/*sql := "select * from (" + "SELECT nfts.*, auctionstemp.startprice AS sellprice, bidcount.offernum, bidcount.maxbidprice " +
		"FROM nfts LEFT JOIN (select * from auctions WHERE deleted_at IS NULL ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid  LEFT JOIN " +
		"(SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid) bidcount " +
		"ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid " + ") a "
	countSql := "select count(*) from (" + "SELECT nfts.*, auctionstemp.startprice AS sellprice, bidcount.offernum, bidcount.maxbidprice " +
		"FROM nfts LEFT JOIN (select * from auctions WHERE deleted_at IS NULL ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid  LEFT JOIN " +
		"(SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid) bidcount " +
		"ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid " + ") a "*/
	if len(sort) > 0 {
		for k, v := range sort {
			if k > 0 {
				orderBy = orderBy + ", "
			}
			orderBy += v.By + " " + v.Order
		}
	}
	if len(orderBy) > 0 {
		orderBy = orderBy + ", id desc"
	} else {
		orderBy = "createdate desc, id desc"
	}
	if len(filter) > 0 {
		queryWhere = nft.joinFilters(filter)
		nftCatchHash := NftCatch.NftCatchHash(queryWhere + orderBy + startIndex + count)
		nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeAuction)
		if nftCatch != nil {
			fmt.Printf("QueryNftByFilter() filter spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}
		querySplit := QueryWhereSplit(queryWhere)
		if querySplit["collections"] != "" {
			if len(querySplit) == 4 {
				collectRec := CollectRec{}
				name := querySplit["name"]
				createaddr := querySplit["createaddr"]
				collectSql := "select createaddr, name, snftcollection, Contracttype from collects where " + name + " AND " + createaddr + " "
				err := nft.db.Raw(collectSql).Scan(&collectRec)
				if err.Error != nil {
					if err.Error != gorm.ErrRecordNotFound {
						log.Println("QueryNftByFilter() Select(snft) err=", err.Error)
						return nil, 0, ErrDataBase
					}
					return []NftInfo{}, 0, nil
				}
				if collectRec.Contracttype == "snft" {
					err := nft.db.Model(&Nfts{}).Select("min(nftaddr)").Where("Collectcreator = ? and Collections = ? ", collectRec.Createaddr, collectRec.Name).Group("snft").Find(&snftAddrs)
					if err.Error != nil {
						log.Println("QueryNftByFilter() Select(snft) err=", err.Error)
						return nil, 0, ErrDataBase
					}
					err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&nftInfo)
					if err.Error != nil {
						log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
						return nil, 0, ErrDataBase
					}
					fmt.Printf("QueryNftByFilter() collect snft spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
					NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(len(nftInfo))})
					return nftInfo, uint64(len(nftInfo)), nil
				} else {
					offset := 0
					if len(startIndex) > 0 {
						offset, _ = strconv.Atoi(startIndex)
					}
					limit := 0
					if len(count) > 0 {
						limit, _ = strconv.Atoi(count)
					}
					err := nft.db.Model(&Nfts{}).Where("Collectcreator = ? and Collections = ? ", collectRec.Createaddr, collectRec.Name).
						Order(orderBy).Offset(offset).Limit(limit).Scan(&nftInfo)
					if err.Error != nil {
						log.Println("QueryNftByFilter() Select(snft) err=", err.Error)
						return nil, 0, ErrDataBase
					}
					err = nft.db.Model(&Nfts{}).Where("Collectcreator = ? and Collections = ? ", collectRec.Createaddr, collectRec.Name).Count(&totalCount)
					if err.Error != nil {
						log.Println("QueryNftByFilter() Select(snft) err=", err.Error)
						return nil, 0, ErrDataBase
					}
					fmt.Printf("QueryNftByFilter() collect nft spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
					NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
					return nftInfo, uint64(totalCount), nil
				}
			}
		}
		if len(querySplit) == 1 && querySplit["createdate"] != "" {
			spendStart := time.Now()
			nftCatchHash := NftCatch.NftCatchHash(startIndex + count)
			nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeNewNft)
			if nftCatch != nil {
				fmt.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
				return nftCatch.NftInfos, nftCatch.Total, nil
			}
			switch nftType {
			case "nft":
				spendStart = time.Now()
				var nftCount int64
				countSql := `select count(id) from nfts where snft = "" and deleted_at is null `
				//countSql = `select count(id) from nfts where snft is NULL `
				fmt.Printf("QueryNftByFilter() countSql sql = %s \n", countSql)
				err := nft.db.Raw(countSql).Scan(&nftCount)
				if err.Error != nil {
					log.Println("QueryNftByFilter() Raw(countSql).Scan(&recCount) err=", err.Error)
					return nil, uint64(0), ErrDataBase
				}
				totalCount = nftCount
				fmt.Printf("QueryNftByFilter() nftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
				spendStart = time.Now()
				nftItem := `Ownaddr, nfts.Name, nfts.Desc, Contract, Tokenid, Nftaddr, Count, Categories, Collectcreator, Collections, Hide, Createaddr, Price, Royalty, Createdate, Favorited, Transcnt, Transamt, Verified, Selltype, Mintstate `
				nftSql := `select ` + nftItem + ` from nfts where snft = "" and deleted_at is null `
				nftSql = nftSql + " order by " + orderBy + " limit " + startIndex + "," + count
				fmt.Printf("QueryNftByFilter() nftInfo sql = %s \n", nftSql)
				err = nft.db.Raw(nftSql).Scan(&nftInfo)
				if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
					log.Println("QueryNftByFilter() Scan(&nftInfo) err=", err)
					return nil, uint64(0), ErrDataBase
				}
				fmt.Printf("QueryNftByFilter() nftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			case "snft":
				countSql := `select count(a.snft) from (select snft from nfts where snft != "" GROUP BY snft) as a`
				var snftCount int64
				err := nft.db.Raw(countSql).Scan(&snftCount)
				if err.Error != nil {
					log.Println("QueryNftByFilter() Raw(countSql).Scan(&snftCount) err=", err.Error)
					return nil, uint64(0), ErrDataBase
				}
				fmt.Printf("QueryNftByFilter() snftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
				spendStart = time.Now()
				totalCount = snftCount
				snftSql := `select min(nftaddr) from nfts where snft != "" group by snft `
				snftInfo := []NftInfo{}
				snftSql = snftSql + " limit " + startIndex + ", " + count
				err = nft.db.Raw(snftSql).Scan(&snftAddrs)
				if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
					log.Println("QueryNftByFilter() Raw(snftSql).Scan(&snftAddrs) err=", err.Error)
					return nil, uint64(0), ErrDataBase
				}
				fmt.Printf("QueryNftByFilter() snftAddrs spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
				spendStart = time.Now()
				fmt.Println("QueryNftByFilter() snftTotalCount=", snftTotalCount)
				err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
				if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
					log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err.Error)
					return nil, uint64(0), ErrDataBase
				}
				nftInfo = append(nftInfo, snftInfo...)
				fmt.Printf("QueryNftByFilter() snftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			//case "nftsnft", "":
			default:
				return nil, uint64(0), errors.New("QueryNftByFilter no params error.")
			}
			fmt.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
			return nftInfo, uint64(totalCount), nil
		}
		snftSql := `SELECT nfts.*, auctionstemp.startprice AS sellprice, offernum, maxbidprice FROM nfts JOIN (select * from auctions WHERE deleted_at IS NULL selltype_condition price_condition ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid left Join (SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid ) bidcount ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid   `
		snftCountSql := `SELECT count(nfts.id) FROM nfts JOIN (select * from auctions WHERE deleted_at IS NULL selltype_condition price_condition ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid left Join (SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid ) bidcount ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid   `
		if querySplit["selltype"] != "" {
			snftSql = strings.Replace(snftSql, "selltype_condition", "and "+querySplit["selltype"], -1)
			snftCountSql = strings.Replace(snftCountSql, "selltype_condition", "and "+querySplit["selltype"], -1)
		} else {
			snftSql = strings.Replace(snftSql, "selltype_condition", " ", -1)
			snftCountSql = strings.Replace(snftCountSql, "selltype_condition", " ", -1)
		}
		if querySplit["sellprice"] != "" {
			snftSql = strings.Replace(snftSql, "price_condition", "and "+querySplit["sellprice"], -1)
			snftCountSql = strings.Replace(snftCountSql, "price_condition", "and "+querySplit["sellprice"], -1)
		} else {
			snftSql = strings.Replace(snftSql, "price_condition", " ", -1)
			snftCountSql = strings.Replace(snftCountSql, "price_condition", " ", -1)
		}
		whereFlag := false
		if querySplit["createdate"] != "" {
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + querySplit["createdate"]
				snftCountSql = snftCountSql + " where " + querySplit["createdate"]
			} else {
				snftSql = snftSql + " and " + querySplit["createdate"]
				snftCountSql = snftCountSql + " and " + querySplit["createdate"]
			}
		}

		if querySplit["offernum"] != "" {
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + querySplit["offernum"]
				snftCountSql = snftCountSql + " where " + querySplit["offernum"]
			} else {
				snftSql = snftSql + " and " + querySplit["offernum"]
				snftCountSql = snftCountSql + " and " + querySplit["offernum"]
			}
		}

		if querySplit["collectcreator"] != "" {
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + querySplit["collectcreator"]
				snftCountSql = snftCountSql + " where " + querySplit["collectcreator"]
			} else {
				snftSql = snftSql + " and " + querySplit["collectcreator"]
				snftCountSql = snftCountSql + " and " + querySplit["collectcreator"]
			}
		}
		if querySplit["collections"] != "" {
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + querySplit["collections"]
				snftCountSql = snftCountSql + " where " + querySplit["collections"]
			} else {
				snftSql = snftSql + " and " + querySplit["collections"]
				snftCountSql = snftCountSql + " and " + querySplit["collections"]
			}
		}
		if querySplit["categories"] != "" {
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + querySplit["categories"]
				snftCountSql = snftCountSql + " where " + querySplit["categories"]
			} else {
				snftSql = snftSql + " and " + querySplit["categories"]
				snftCountSql = snftCountSql + " and " + querySplit["categories"]
			}
		}
		if whereFlag == false {
			whereFlag = true
			snftSql = snftSql + " where " + " (nfts.deleted_at is null) "
			snftCountSql = snftCountSql + " where " + " (nfts.deleted_at is null) "
		} else {
			snftSql = snftSql + " and " + " (nfts.deleted_at is null) "
			snftCountSql = snftCountSql + " and " + " (nfts.deleted_at is null) "
		}
		snftSql = snftSql + " order by " + orderBy
		if len(startIndex) > 0 && len(count) > 0 {
			snftSql = snftSql + " limit " + startIndex + ", " + count
		} else {
			snftSql = snftSql + " limit " + "0" + ", " + "1"
		}
		fmt.Println("QueryNftByFilter() snftSql=", snftSql)
		err := nft.db.Raw(snftSql).Scan(&nftInfo)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Raw(snftSql).Scan(&nftInfo) err=", err.Error)
			return nil, uint64(0), ErrDataBase
		}
		err = nft.db.Raw(snftCountSql).Scan(&totalCount)
		if err.Error != nil {
			log.Println("QueryNftByFilter() Scan(&totalCount) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		fmt.Printf("QueryNftByFilter() normal spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
		return nftInfo, uint64(totalCount), nil
	} else {
		spendStart := time.Now()
		nftCatchHash := NftCatch.NftCatchHash(nftType + startIndex + count)
		nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeNewNft)
		if nftCatch != nil {
			fmt.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}
		switch nftType {
		case "nft":
			spendStart = time.Now()
			var nftCount int64
			countSql := `select count(id) from nfts where snft = "" and deleted_at  is null `
			err := nft.db.Raw(countSql).Scan(&nftCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Raw(countSql).Scan(&recCount) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() nftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			nftItem := `Ownaddr, nfts.Name, nfts.Desc, Contract, Tokenid, Nftaddr, Count, Categories, Collectcreator, Collections, Hide, Createaddr, Price, Transprice, Royalty, Createdate, Favorited, Transcnt, Transamt, Verified, Selltype, Mintstate `
			nftSql := `select ` + nftItem + ` from nfts where snft = "" and deleted_at  is null `
			nftSql = nftSql + " order by " + orderBy + " limit " + startIndex + "," + count
			fmt.Printf("QueryNftByFilter() nftInfo sql = %s\n", nftSql)
			err = nft.db.Raw(nftSql).Scan(&nftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&nftInfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			totalCount = nftCount

			fmt.Printf("QueryNftByFilter() nftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		case "snft":
			spendStart = time.Now()
			countSql := `select count(a.snft) from (select snft from nfts where snft != "" GROUP BY snft) as a`
			var snftCount int64
			err := nft.db.Raw(countSql).Scan(&snftCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Raw(countSql).Scan(&snftCount) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() snftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			totalCount = snftCount
			snftInfo := []NftInfo{}
			spendStart = time.Now()
			snftSql := `select min(nftaddr) from nfts where snft != "" group by snft `
			snftSql = snftSql + " limit " + startIndex + ", " + count
			err = nft.db.Raw(snftSql).Scan(&snftAddrs)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Raw(snftSql).Scan(&snftAddrs) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() snftAddrs spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			fmt.Println("QueryNftByFilter() snftTotalCount=", snftTotalCount)
			err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			nftInfo = append(nftInfo, snftInfo...)
			fmt.Printf("QueryNftByFilter() snftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		case "nftsnft", "":
			spendStart = time.Now()
			countSql := `select count(a.snft) from (select snft from nfts where snft != "" GROUP BY snft) as a`
			var snftCount int64
			err := nft.db.Raw(countSql).Scan(&snftCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Raw(countSql).Scan(&snftCount) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() snftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			var nftCount int64
			countSql = `select count(id) from nfts where snft = "" and deleted_at  is null  `
			//countSql = `select count(id) from nfts where snft is NULL `
			err = nft.db.Raw(countSql).Scan(&nftCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Raw(countSql).Scan(&recCount) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() nftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			totalCount = snftCount + nftCount
			snftSql := `select min(nftaddr) from nfts where snft != "" group by snft `
			snftInfo := []NftInfo{}
			if len(startIndex) >= 0 && len(count) > 0 {
				scount, _ := strconv.Atoi(count)
				scount = scount / 2
				count = strconv.Itoa(scount)
				index, _ := strconv.Atoi(startIndex)
				index = index / 2
				startIndex = strconv.Itoa(index)
				snftSql = snftSql + " limit " + startIndex + ", " + count
			} else {
				count = "1"
				startIndex = "0"
				snftSql = snftSql + " limit " + startIndex + ", " + count
			}
			fmt.Printf("QueryNftByFilter() snftSql sql = %s \n", snftSql)
			err = nft.db.Raw(snftSql).Scan(&snftAddrs)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Raw(snftSql).Scan(&snftAddrs) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() snftAddrs spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			fmt.Println("QueryNftByFilter() snftTotalCount=", snftTotalCount)

			err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() snftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			nftItem := `Ownaddr, nfts.Name, nfts.Desc, Contract, Tokenid, Nftaddr, Count, Categories, Collectcreator, Collections, Hide, Createaddr, Price, Transprice, Royalty, Createdate, Favorited, Transcnt, Transamt, Verified, Selltype, Mintstate `
			nftSql := `select ` + nftItem + ` from nfts where snft = ""  and deleted_at  is null `
			nftSql = nftSql + " order by " + orderBy + " limit " + startIndex + "," + count
			fmt.Printf("QueryNftByFilter() nftInfo sql = %s\n", nftSql)
			err = nft.db.Raw(nftSql).Scan(&nftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&nftInfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() nftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			fmt.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			nftInfo = append(nftInfo, snftInfo...)
		}
		NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
		return nftInfo, uint64(totalCount), nil
		/*countSql := `select count(a.snft) from (select snft from nfts where snft != "" GROUP BY snft) as a`
		var recCount int64
		err := nft.db.Raw(countSql).Scan(&recCount)
		if err.Error != nil {
			log.Println("QueryNftByFilter() Raw(countSql).Scan(&recCount) err=", err.Error)
			return nil, uint64(0), ErrDataBase
		}
		totalCount = recCount
		countSql = `select count(id) from nfts where snft = "" `
		err = nft.db.Raw(countSql).Scan(&recCount)
		if err.Error != nil {
			log.Println("QueryNftByFilter() Raw(countSql).Scan(&recCount) err=", err.Error)
			return nil, uint64(0), ErrDataBase
		}
		totalCount = totalCount + recCount
		snftSql := `SELECT nfts.* FROM nfts JOIN (select contract, tokenid from nfts ) nfts1 ON nfts.contract = nfts1.contract AND nfts.tokenid = nfts1.tokenid left JOIN (select min(nftaddr) as minnftaddr from nfts where snft != "" GROUP BY snft ) nfts2 ON nfts.nftaddr = nfts2.minnftaddr `
		snftSql = snftSql + " order by " + orderBy
		if len(startIndex) > 0 && len(count) > 0 {
			snftSql = snftSql + " limit " + startIndex + ", " + count
		} else {
			snftSql = snftSql + " limit " + "0" + ", " + "1"
		}
		err = nft.db.Raw(snftSql).Scan(&nftInfo)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound{
			log.Println("QueryNftByFilter() Raw(snftSql).Scan(&nftInfo) err=", err.Error)
			return nil, uint64(0), ErrDataBase
		}
		fmt.Printf("QueryNftByFilter() filter=nil spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftInfo, uint64(totalCount) + uint64(snftTotalCount), nil*/
	}
	return nftInfo, uint64(totalCount) + uint64(snftTotalCount), nil
}

func (nft NftDb) NftFilterProc(filter []StQueryField, sort []StSortField, startIndex string, count string) ([]NftInfo, uint64, error) {
	var queryWhere string
	var orderBy string
	var totalCount int64
	nftInfo := []NftInfo{}
	spendT := time.Now()

	/*sql := "select * from (" + "SELECT nfts.*, auctionstemp.startprice AS sellprice, bidcount.offernum, bidcount.maxbidprice " +
		"FROM nfts LEFT JOIN (select * from auctions WHERE deleted_at IS NULL ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid  LEFT JOIN " +
		"(SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid) bidcount " +
		"ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid " + ") a "
	countSql := "select count(*) from (" + "SELECT nfts.*, auctionstemp.startprice AS sellprice, bidcount.offernum, bidcount.maxbidprice " +
		"FROM nfts LEFT JOIN (select * from auctions WHERE deleted_at IS NULL ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid  LEFT JOIN " +
		"(SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid) bidcount " +
		"ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid " + ") a "*/
	if len(sort) > 0 {
		for k, v := range sort {
			if k > 0 {
				orderBy = orderBy + ", "
			}
			orderBy += v.By + " " + v.Order
		}
	}
	if len(orderBy) > 0 {
		orderBy = orderBy + ", id desc"
	} else {
		orderBy = "createdate desc, id desc"
	}
	if len(filter) > 0 {
		queryWhere = nft.joinFilters(filter)
		/*nftCatchHash := NftCatch.NftCatchHash(queryWhere + orderBy + startIndex + count)
		nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeAuction)
		if nftCatch != nil {
			log.Printf("QueryNftByFilter() filter spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}*/
		queryCatchSql := queryWhere + orderBy + startIndex + count
		nftCatch := NftFilter{}
		cerr := GetRedisCatch().GetCatchData("QueryNftByFilterNftSnft", queryCatchSql, &nftCatch)
		if cerr == nil {
			log.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}
		querySplit := QueryWhereSplit(queryWhere)
		if querySplit["selltype"] == "" && querySplit["offernum"] == "" && querySplit["sellprice"] == "" {
			spendStart := time.Now()
			whereFlag := false
			nftItem := `Ownaddr, nfts.Name, nfts.Desc, Contract, Tokenid, Nftaddr, Count, Categories, Collectcreator, Collections, Hide, Createaddr, Price, Transprice, Royalty, Createdate, Favorited, Transcnt, Transamt, Verified, Selltype, Mintstate `
			nftSql := `select ` + nftItem + ` ,transprice as sellprice from nfts  `
			nftCountSql := `select count(id) from nfts `
			if querySplit["createdate"] != "" {
				if whereFlag == false {
					whereFlag = true
					nftSql = nftSql + " where " + querySplit["createdate"]
					nftCountSql = nftCountSql + " where " + querySplit["createdate"]
				} else {
					nftSql = nftSql + " and " + querySplit["createdate"]
					nftCountSql = nftCountSql + " and " + querySplit["createdate"]
				}
			}
			if querySplit["categories"] != "" {
				if whereFlag == false {
					whereFlag = true
					nftSql = nftSql + " where " + querySplit["categories"]
					nftCountSql = nftCountSql + " where " + querySplit["categories"]
				} else {
					nftSql = nftSql + " and " + querySplit["categories"]
					nftCountSql = nftCountSql + " and " + querySplit["categories"]
				}
			}
			if querySplit["collectcreator"] != "" {
				if whereFlag == false {
					whereFlag = true
					nftSql = nftSql + " where " + querySplit["collectcreator"]
					nftCountSql = nftCountSql + " where " + querySplit["collectcreator"]
				} else {
					nftSql = nftSql + " and " + querySplit["collectcreator"]
					nftCountSql = nftCountSql + " and " + querySplit["collectcreator"]
				}
			}
			if querySplit["collections"] != "" {
				if whereFlag == false {
					whereFlag = true
					nftSql = nftSql + " where " + querySplit["collections"]
					nftCountSql = nftCountSql + " where " + querySplit["collections"]
				} else {
					nftSql = nftSql + " and " + querySplit["collections"]
					nftCountSql = nftCountSql + " and " + querySplit["collections"]
				}
			}
			if whereFlag == false {
				whereFlag = true
				nftSql = nftSql + " where " + " (nfts.deleted_at is null) "
				nftCountSql = nftCountSql + " where " + " (nfts.deleted_at is null) "
			} else {
				nftSql = nftSql + " and " + " (nfts.deleted_at is null) "
				nftCountSql = nftCountSql + " and " + " (nfts.deleted_at is null) "
			}
			//snftSql = snftSql + " group by snft "
			//nftCountSql = nftCountSql + " ) as a"
			nftSql = nftSql + " order by " + orderBy
			if len(startIndex) > 0 && len(count) > 0 {
				nftSql = nftSql + " limit " + startIndex + ", " + count
			} else {
				nftSql = nftSql + " limit " + "0" + ", " + "1"
			}
			//nftInfo := []NftInfo{}
			log.Printf("QueryNftByFilter() nftSql=%s\n", nftSql)
			err := nft.db.Raw(nftSql).Scan(&nftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			//nftInfo = append(nftInfo, snftInfo...)
			log.Printf("QueryNftByFilter() nftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendT = time.Now()
			log.Printf("QueryNftByFilter() nftCountSql=%s\n", nftCountSql)
			var nftCount int64
			if querySplit["collectcreator"] != "" && querySplit["collections"] != "" {
				nftCountSql := `select totalcount from collects where `
				createaddr := querySplit["collectcreator"]
				createaddr = strings.Replace(createaddr, "collectcreator", "createaddr", -1)
				name := querySplit["collections"]
				name = strings.Replace(name, "collections", "name", -1)
				nftCountSql = nftCountSql + createaddr + " and " + name + " and " + " deleted_at IS NULL "
				err = nft.db.Raw(nftCountSql).Scan(&nftCount)
				if err.Error != nil {
					log.Println("QueryNftByFilter() Raw(countSql).Scan(&nftCount) err=", err.Error)
					return []NftInfo{}, uint64(0), err.Error
				}
			} else {
				err = nft.db.Raw(nftCountSql).Scan(&nftCount)
				if err.Error != nil {
					log.Println("QueryNftByFilter() Raw(countSql).Scan(&nftCount) err=", err.Error)
					return []NftInfo{}, uint64(0), err.Error
				}
			}
			nftTokenIds := []string{}
			for _, info := range nftInfo {
				nftTokenIds = append(nftTokenIds, info.Tokenid)
			}
			sellPrices := []struct {
				Tokenid string
				Price   uint64
			}{}
			err = nft.db.Model(&Auction{}).Select("tokenid, startprice as price").Where("tokenid in ?", nftTokenIds).Scan(&sellPrices)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&sellPrices) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			offerPrices := []struct {
				Tokenid     string
				Offernum    uint64
				Maxbidprice uint64
			}{}
			err = nft.db.Model(&Bidding{}).Select("tokenid, count(id) as offernum, max(price) as maxbidprice").Where("tokenid in ?", nftTokenIds).Group("tokenid").Scan(&offerPrices)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Select(snft) err=", err.Error)
				return nil, 0, ErrDataBase
			}

			for i, info := range nftInfo {
				for _, price := range sellPrices {
					if info.Tokenid == price.Tokenid {
						nftInfo[i].Sellprice = price.Price
					}
				}
				for _, price := range offerPrices {
					if info.Tokenid == price.Tokenid {
						nftInfo[i].Offernum = price.Offernum
						nftInfo[i].Maxbidprice = price.Maxbidprice
					}
				}
			}
			totalCount = nftCount
			log.Printf("QueryNftByFilter() nftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			//NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
			GetRedisCatch().CatchQueryData("QueryNftByFilterNftSnft", queryCatchSql, &NftFilter{nftInfo, uint64(totalCount)})
			return nftInfo, uint64(totalCount), nil
		} else {
			snftSql := `SELECT nfts.*, auctionstemp.startprice AS sellprice, offernum, maxbidprice FROM nfts JOIN (select * from auctions WHERE deleted_at IS NULL selltype_condition price_condition ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid left Join (SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid ) bidcount ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid   `
			snftCountSql := `SELECT count(nfts.id) FROM nfts JOIN (select * from auctions WHERE deleted_at IS NULL selltype_condition price_condition ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid left Join (SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid ) bidcount ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid   `
			if querySplit["selltype"] != "" {
				snftSql = strings.Replace(snftSql, "selltype_condition", "and "+querySplit["selltype"], -1)
				snftCountSql = strings.Replace(snftCountSql, "selltype_condition", "and "+querySplit["selltype"], -1)
			} else {
				snftSql = strings.Replace(snftSql, "selltype_condition", " ", -1)
				snftCountSql = strings.Replace(snftCountSql, "selltype_condition", " ", -1)
			}
			if querySplit["sellprice"] != "" {
				snftSql = strings.Replace(snftSql, "price_condition", "and "+querySplit["sellprice"], -1)
				snftCountSql = strings.Replace(snftCountSql, "price_condition", "and "+querySplit["sellprice"], -1)
			} else {
				snftSql = strings.Replace(snftSql, "price_condition", " ", -1)
				snftCountSql = strings.Replace(snftCountSql, "price_condition", " ", -1)
			}
			whereFlag := false
			if querySplit["createdate"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["createdate"]
					snftCountSql = snftCountSql + " where " + querySplit["createdate"]
				} else {
					snftSql = snftSql + " and " + querySplit["createdate"]
					snftCountSql = snftCountSql + " and " + querySplit["createdate"]
				}
			}
			if querySplit["offernum"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["offernum"]
					snftCountSql = snftCountSql + " where " + querySplit["offernum"]
				} else {
					snftSql = snftSql + " and " + querySplit["offernum"]
					snftCountSql = snftCountSql + " and " + querySplit["offernum"]
				}
			}
			if querySplit["collectcreator"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["collectcreator"]
					snftCountSql = snftCountSql + " where " + querySplit["collectcreator"]
				} else {
					snftSql = snftSql + " and " + querySplit["collectcreator"]
					snftCountSql = snftCountSql + " and " + querySplit["collectcreator"]
				}
			}
			if querySplit["collections"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["collections"]
					snftCountSql = snftCountSql + " where " + querySplit["collections"]
				} else {
					snftSql = snftSql + " and " + querySplit["collections"]
					snftCountSql = snftCountSql + " and " + querySplit["collections"]
				}
			}
			if querySplit["categories"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["categories"]
					snftCountSql = snftCountSql + " where " + querySplit["categories"]
				} else {
					snftSql = snftSql + " and " + querySplit["categories"]
					snftCountSql = snftCountSql + " and " + querySplit["categories"]
				}
			}
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + " (nfts.deleted_at is null) "
				snftCountSql = snftCountSql + " where " + " (nfts.deleted_at is null) "
			} else {
				snftSql = snftSql + " and " + " (nfts.deleted_at is null) "
				snftCountSql = snftCountSql + " and " + " (nfts.deleted_at is null) "
			}
			snftSql = snftSql + " order by " + orderBy
			if len(startIndex) > 0 && len(count) > 0 {
				snftSql = snftSql + " limit " + startIndex + ", " + count
			} else {
				snftSql = snftSql + " limit " + "0" + ", " + "1"
			}
			log.Println("QueryNftByFilter() snftSql=", snftSql)
			err := nft.db.Raw(snftSql).Scan(&nftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Raw(snftSql).Scan(&nftInfo) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			log.Println("QueryNftByFilter() snftCountSql=", snftCountSql)
			err = nft.db.Raw(snftCountSql).Scan(&totalCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Scan(&totalCount) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			log.Printf("QueryNftByFilter() normal spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			//NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
			GetRedisCatch().CatchQueryData("QueryNftByFilterNftSnft", queryCatchSql, &NftFilter{nftInfo, uint64(totalCount)})
			return nftInfo, uint64(totalCount), nil
		}
	} else {
		spendStart := time.Now()
		/*nftCatchHash := NftCatch.NftCatchHash(startIndex + count)
		nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeNewNft)
		if nftCatch != nil {
			log.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}*/
		queryCatchSql := startIndex + count
		nftCatch := NftFilter{}
		cerr := GetRedisCatch().GetCatchData("QueryNftByFilterNftSnft", queryCatchSql, &nftCatch)
		if cerr == nil {
			log.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}
		spendStart = time.Now()
		var nftCount int64
		//countSql := `select max(id) from nfts where snft = "" and deleted_at  is null `
		//err := nft.db.Raw(countSql).Scan(&nftCount)
		err := nft.db.Model(&SysInfos{}).Select("nfttotal").Last(&nftCount)
		if err.Error != nil {
			if !strings.Contains(err.Error.Error(), "converting NULL to int64 is unsupported") {
				log.Println("QueryNftByFilter() Raw(countSql).Scan(&recCount) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			return nftInfo, 0, nil
		}
		log.Printf("QueryNftByFilter() nftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		spendStart = time.Now()
		nftItem := `Ownaddr, nfts.Name, nfts.Desc, Contract, Tokenid, Nftaddr, Count, Categories, Collectcreator, Collections, Hide, Createaddr, Price, Transprice, Royalty, Createdate, Favorited, Transcnt, Transamt, Verified, Selltype, Mintstate `
		nftSql := `select ` + nftItem + ` from nfts where snft = "" and deleted_at  is null `
		nftSql = nftSql + " order by " + orderBy + " limit " + startIndex + "," + count
		log.Printf("QueryNftByFilter() nftInfo sql = %s\n", nftSql)
		err = nft.db.Raw(nftSql).Scan(&nftInfo)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Scan(&nftInfo) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		nftTokenIds := []string{}
		for _, info := range nftInfo {
			nftTokenIds = append(nftTokenIds, info.Tokenid)
		}
		sellPrices := []struct {
			Tokenid string
			Price   uint64
		}{}
		err = nft.db.Model(&Auction{}).Select("tokenid, startprice as price").Where("tokenid in ?", nftTokenIds).Scan(&sellPrices)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Scan(&sellPrices) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		offerPrices := []struct {
			Tokenid  string
			Offernum uint64
			Price    uint64
		}{}
		err = nft.db.Model(&Bidding{}).Select("tokenid, count(id) as offernum, max(price) as price").Where("tokenid in ?", nftTokenIds).Group("tokenid").Scan(&offerPrices)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Select(snft) err=", err.Error)
			return nil, 0, err.Error
		}

		for i, info := range nftInfo {
			for _, price := range sellPrices {
				if info.Tokenid == price.Tokenid {
					nftInfo[i].Sellprice = price.Price
				}
			}
			for _, price := range offerPrices {
				if info.Tokenid == price.Tokenid {
					nftInfo[i].Offernum = price.Offernum
					nftInfo[i].Maxbidprice = price.Price
				}
			}
		}
		totalCount = nftCount
		log.Printf("QueryNftByFilter() nftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		//NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
		GetRedisCatch().CatchQueryData("QueryNftByFilterNftSnft", queryCatchSql, &NftFilter{nftInfo, uint64(totalCount)})
		return nftInfo, uint64(totalCount), nil
	}
	return nil, 0, nil
}

type SnftFilters struct {
	Minnftaddr    string `json:"minnftaddr"`
	Transamt0     int    `json:"transamt"`
	Sellprice     int    `json:"transprice"`
	Transtime0    int    `json:"transtime"`
	Verifiedtime0 int    `json:"verifiedtime"`
}

func (nft NftDb) SnftFilterProc(filter []StQueryField, sort []StSortField, startIndex string, count string) ([]NftInfo, uint64, error) {
	var queryWhere string
	var orderBy string
	var totalCount int64
	nftInfo := []NftInfo{}
	spendT := time.Now()
	snftAddrs := []string{}
	snftTotalCount := 0

	if len(sort) > 0 {
		for k, v := range sort {
			if k > 0 {
				orderBy = orderBy + ", "
			}
			orderBy += v.By + " " + v.Order
		}
	}
	if len(orderBy) > 0 {
		//orderBy = orderBy + ", id desc"
		orderBy = orderBy
	} else {
		//orderBy = "createdate desc, id desc"
		orderBy = "createdate desc"
	}
	if len(filter) > 0 {
		queryWhere = nft.joinFilters(filter)
		/*nftCatchHash := NftCatch.NftCatchHash(queryWhere + orderBy + startIndex + count)
		nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeAuction)
		if nftCatch != nil {
			fmt.Printf("QueryNftByFilter() filter spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}*/
		queryCatchSql := queryWhere + orderBy + startIndex + count
		nftCatch := NftFilter{}
		cerr := GetRedisCatch().GetCatchData("QueryNftByFilterNftSnft", queryCatchSql, &nftCatch)
		if cerr == nil {
			log.Printf("QueryNftByFilter() snft filter spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}
		querySplit := QueryWhereSplit(queryWhere)
		if querySplit["selltype"] == "" && querySplit["offernum"] == "" && querySplit["sellprice"] == "" {
			spendStart := time.Now()
			whereFlag := false
			snftSql := `select min(nftaddr) as minnftaddr, sum(transamt) as transamt, sum(transprice) as sellprice,min(transtime) as transtime, min(verifiedtime) as verifiedtime, min(createdate) as createdate from nfts  `
			snftCountSql := `select count(minnftaddr) from ( select min(nftaddr) as minnftaddr from nfts `
			if querySplit["createdate"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["createdate"]
					snftCountSql = snftCountSql + " where " + querySplit["createdate"]
				} else {
					snftSql = snftSql + " and " + querySplit["createdate"]
					snftCountSql = snftCountSql + " and " + querySplit["createdate"]
				}
			}
			if querySplit["categories"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["categories"]
					snftCountSql = snftCountSql + " where " + querySplit["categories"]
				} else {
					snftSql = snftSql + " and " + querySplit["categories"]
					snftCountSql = snftCountSql + " and " + querySplit["categories"]
				}
			}
			if querySplit["collectcreator"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["collectcreator"]
					snftCountSql = snftCountSql + " where " + querySplit["collectcreator"]
				} else {
					snftSql = snftSql + " and " + querySplit["collectcreator"]
					snftCountSql = snftCountSql + " and " + querySplit["collectcreator"]
				}
			}
			if querySplit["collections"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["collections"]
					snftCountSql = snftCountSql + " where " + querySplit["collections"]
				} else {
					snftSql = snftSql + " and " + querySplit["collections"]
					snftCountSql = snftCountSql + " and " + querySplit["collections"]
				}
			}
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
				snftCountSql = snftCountSql + " where " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
			} else {
				snftSql = snftSql + " and " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
				snftCountSql = snftCountSql + " and " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
			}
			snftSql = snftSql + " group by snft "
			snftCountSql = snftCountSql + " group by snft ) as a"
			snftSql = snftSql + " order by " + orderBy
			if len(startIndex) > 0 && len(count) > 0 {
				snftSql = snftSql + " limit " + startIndex + ", " + count
			} else {
				snftSql = snftSql + " limit " + "0" + ", " + "1"
			}
			fmt.Printf("QueryNftByFilter() snftSql=%s \n", snftSql)
			var snftf []SnftFilters
			err := nft.db.Raw(snftSql).Scan(&snftf)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Raw(snftSql).Scan(&snftinfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() snftAddrs spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendStart = time.Now()
			fmt.Println("QueryNftByFilter() snftTotalCount=", snftTotalCount)
			var snftAddrs []string
			for _, snft := range snftf {
				snftAddrs = append(snftAddrs, snft.Minnftaddr)
			}
			snftInfo := []NftInfo{}
			err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			nftInfo = append(nftInfo, snftInfo...)
			fmt.Printf("QueryNftByFilter() snftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			spendT = time.Now()
			var snftCount int64
			err = nft.db.Raw(snftCountSql).Scan(&snftCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Raw(countSql).Scan(&snftCount) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			totalCount = snftCount
			fmt.Printf("QueryNftByFilter() snftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
			//NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
			GetRedisCatch().CatchQueryData("QueryNftByFilterNftSnft", queryCatchSql, &NftFilter{nftInfo, uint64(totalCount)})
			return nftInfo, uint64(totalCount), nil
		} else {
			snftSql := `SELECT nfts.*, auctionstemp.startprice AS sellprice, offernum, maxbidprice FROM nfts JOIN (select * from auctions WHERE deleted_at IS NULL selltype_condition price_condition ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid left Join (SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid ) bidcount ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid   `
			snftCountSql := `SELECT count(nfts.id) FROM nfts JOIN (select * from auctions WHERE deleted_at IS NULL selltype_condition price_condition ) auctionstemp ON nfts.contract = auctionstemp.contract AND nfts.tokenid = auctionstemp.tokenid left Join (SELECT contract, tokenid, COUNT(*) AS offernum, MAX(price) AS maxbidprice FROM biddings WHERE deleted_at IS NULL GROUP BY contract, tokenid ) bidcount ON nfts.contract = bidcount.contract AND nfts.tokenid = bidcount.tokenid   `
			if querySplit["selltype"] != "" {
				snftSql = strings.Replace(snftSql, "selltype_condition", "and "+querySplit["selltype"], -1)
				snftCountSql = strings.Replace(snftCountSql, "selltype_condition", "and "+querySplit["selltype"], -1)
			} else {
				snftSql = strings.Replace(snftSql, "selltype_condition", " ", -1)
				snftCountSql = strings.Replace(snftCountSql, "selltype_condition", " ", -1)
			}
			if querySplit["sellprice"] != "" {
				snftSql = strings.Replace(snftSql, "price_condition", "and "+querySplit["sellprice"], -1)
				snftCountSql = strings.Replace(snftCountSql, "price_condition", "and "+querySplit["sellprice"], -1)
			} else {
				snftSql = strings.Replace(snftSql, "price_condition", " ", -1)
				snftCountSql = strings.Replace(snftCountSql, "price_condition", " ", -1)
			}
			whereFlag := false
			if querySplit["createdate"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["createdate"]
					snftCountSql = snftCountSql + " where " + querySplit["createdate"]
				} else {
					snftSql = snftSql + " and " + querySplit["createdate"]
					snftCountSql = snftCountSql + " and " + querySplit["createdate"]
				}
			}
			if querySplit["offernum"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["offernum"]
					snftCountSql = snftCountSql + " where " + querySplit["offernum"]
				} else {
					snftSql = snftSql + " and " + querySplit["offernum"]
					snftCountSql = snftCountSql + " and " + querySplit["offernum"]
				}
			}
			if querySplit["collectcreator"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["collectcreator"]
					snftCountSql = snftCountSql + " where " + querySplit["collectcreator"]
				} else {
					snftSql = snftSql + " and " + querySplit["collectcreator"]
					snftCountSql = snftCountSql + " and " + querySplit["collectcreator"]
				}
			}
			if querySplit["collections"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["collections"]
					snftCountSql = snftCountSql + " where " + querySplit["collections"]
				} else {
					snftSql = snftSql + " and " + querySplit["collections"]
					snftCountSql = snftCountSql + " and " + querySplit["collections"]
				}
			}
			if querySplit["categories"] != "" {
				if whereFlag == false {
					whereFlag = true
					snftSql = snftSql + " where " + querySplit["categories"]
					snftCountSql = snftCountSql + " where " + querySplit["categories"]
				} else {
					snftSql = snftSql + " and " + querySplit["categories"]
					snftCountSql = snftCountSql + " and " + querySplit["categories"]
				}
			}
			if whereFlag == false {
				whereFlag = true
				snftSql = snftSql + " where " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
				snftCountSql = snftCountSql + " where " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
			} else {
				snftSql = snftSql + " and " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
				snftCountSql = snftCountSql + " and " + " Pledgestate != " + ` "Pledge" ` + " and " + " deleted_at is null "
			}
			snftSql = snftSql + " order by " + orderBy
			if len(startIndex) > 0 && len(count) > 0 {
				snftSql = snftSql + " limit " + startIndex + ", " + count
			} else {
				snftSql = snftSql + " limit " + "0" + ", " + "1"
			}
			fmt.Println("QueryNftByFilter() snftSql=", snftSql)
			err := nft.db.Raw(snftSql).Scan(&nftInfo)
			if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
				log.Println("QueryNftByFilter() Raw(snftSql).Scan(&nftInfo) err=", err.Error)
				return nil, uint64(0), ErrDataBase
			}
			err = nft.db.Raw(snftCountSql).Scan(&totalCount)
			if err.Error != nil {
				log.Println("QueryNftByFilter() Scan(&totalCount) err=", err)
				return nil, uint64(0), ErrDataBase
			}
			fmt.Printf("QueryNftByFilter() normal spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			//NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
			GetRedisCatch().CatchQueryData("QueryNftByFilterNftSnft", queryCatchSql, &NftFilter{nftInfo, uint64(totalCount)})
			return nftInfo, uint64(totalCount), nil
		}
	} else {
		spendStart := time.Now()
		countSql := `select count(a.snft) from (select snft from nfts where snft != "" GROUP BY snft) as a`
		/*nftCatchHash := NftCatch.NftCatchHash(countSql + startIndex + count)
		nftCatch := NftCatch.GetByHash(nftCatchHash, NftFlushTypeNewNft)
		if nftCatch != nil {
			fmt.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}*/
		queryCatchSql := countSql + startIndex + count
		nftCatch := NftFilter{}
		cerr := GetRedisCatch().GetCatchData("QueryNftByFilterNftSnft", queryCatchSql, &nftCatch)
		if cerr == nil {
			log.Printf("QueryNftByFilter() default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			return nftCatch.NftInfos, nftCatch.Total, nil
		}
		var snftCount int64
		//err := nft.db.Raw(countSql).Scan(&snftCount)
		//if err.Error != nil {
		//	log.Println("QueryNftByFilter() Raw(countSql).Scan(&snftCount) err=", err.Error)
		//	return nil, uint64(0), ErrDataBase
		//}
		//
		err := nft.db.Model(&SysInfos{}).Select("snfttotal").Last(&snftCount)
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				log.Println("SnftFilterProc() select nfttotal err=", err)
				return nil, uint64(0), ErrNotFound
			}
			log.Println("SnftFilterProc() select nfttotal err=", err.Error)
			return nil, uint64(0), ErrDataBase
		}
		fmt.Printf("QueryNftByFilter() snftCount spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		totalCount = snftCount
		snftInfo := []NftInfo{}
		spendStart = time.Now()
		//snftSql := ` select min(nftaddr) from nfts where snft != "" group by snft `
		snftSql := `select min(nftaddr) from nfts where snft != "" and  Pledgestate != "Pledge" and deleted_at is null group by snft `
		snftSql = snftSql + " limit " + startIndex + ", " + count
		err = nft.db.Raw(snftSql).Scan(&snftAddrs)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Raw(snftSql).Scan(&snftAddrs) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		fmt.Printf("QueryNftByFilter() snftAddrs spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		spendStart = time.Now()
		fmt.Println("QueryNftByFilter() snftTotalCount=", snftTotalCount)
		err = nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("QueryNftByFilter() Scan(&snftInfo) err=", err)
			return nil, uint64(0), ErrDataBase
		}
		nftInfo = append(nftInfo, snftInfo...)
		fmt.Printf("QueryNftByFilter() snftInfo spend time=%s time.now=%s\n", time.Now().Sub(spendStart), time.Now())
		//NftCatch.SetByHash(nftCatchHash, &NftFilter{nftInfo, uint64(totalCount)})
		GetRedisCatch().CatchQueryData("QueryNftByFilterNftSnft", queryCatchSql, &NftFilter{nftInfo, uint64(totalCount)})
		return nftInfo, uint64(totalCount), nil
	}
	return nftInfo, uint64(totalCount) + uint64(snftTotalCount), nil
}

func (nft NftDb) QueryNftByFilterNftSnft(filter []StQueryField, sort []StSortField, nftType,
	startIndex string, count string) ([]NftInfo, uint64, error) {
	switch nftType {
	case "nft":
		return nft.NftFilterProc(filter, sort, startIndex, count)
	case "snft":
		return nft.SnftFilterProc(filter, sort, startIndex, count)
	default:
		return nil, 0, nil
	}
}
