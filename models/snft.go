package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type SnftCollectRec struct {
	Createaddr   string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Contract     string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Contracttype string `json:"contracttype" gorm:"type:char(20) CHARACTER SET utf8mb4 NOT NULL;comment:'contract type'"`
	Name         string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'collection name'"`
	Desc         string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'collection discription'"`
	Categories   string `json:"categories" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'collection category'"`
	Totalcount   int    `json:"total_count" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'Total number of nfts in the collection'"`
	Transcnt     int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt     uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	SigData      string `json:"sig" gorm:"type:longtext NOT NULL;comment:'sign'"`
	Img          string `json:"img" gorm:"type:longtext NOT NULL;comment:'logo picture'"`
	Tokenid      string `json:"tokenid" gorm:"type:char(42) ;comment:'nft token id'"`
	Snft         string `json:"snft" gorm:"type:longtext ;comment:'16 snft'"`
	//Period       string `json:"period" gorm:"type:char(42) ;comment:'stage'"`
	Local     string `json:"local" gorm:"type:char(42) ;comment:'stage position'"`
	Exchanger string `json:"exchanger" gorm:"type:char(42) ;comment:'The address of the exchange to which snft belongs'"`
	Extend    string `json:"extend" gorm:"type:longtext ;comment:'extend data'"`
}

type SnftCollect struct {
	gorm.Model
	SnftCollectRec
}

func (v SnftCollect) TableName() string {
	return "snftcollect"
}

type SnftPhaseRec struct {
	Createaddr string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'creator's address'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'name'"`
	Desc       string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'description'"`
	Vote       int    `json:"vote" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'votes'"`
	Accedvote  string `json:"accedvote" gorm:"type:char(20) ;comment:'whether to vote'"`
	Accedeth   string `json:"accedeth" gorm:"type:char(20) ;comment:'Whether to select snft'"`
	Categories string `json:"categories" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'category'"`
	Tokenid    string `json:"tokenid" gorm:"type:char(42) ;comment:'snft token id'"`
	Collect    string `json:"collect" gorm:"type:longtext ;comment:'collection'"`
	Meta       string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta information'"`
	Extend     string `json:"extend" gorm:"type:longtext ;comment:'extend data'"`
}

type SnftPhase struct {
	gorm.Model
	SnftPhaseRec
}

func (v SnftPhase) TableName() string {
	return "snftphase"
}

type SnftRec struct {
	Ownaddr    string `json:"ownaddr" gorm:"type:char(42) ;comment:'nft owner address'"`
	Md5        string `json:"md5" gorm:"type:longtext ;comment:'Picture md5 value'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft name'"`
	Desc       string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft description'"`
	Meta       string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta information'"`
	Nftmeta    string `json:"nftmeta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta、tokenid'"`
	Url        string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc raw data hold address'"`
	Contract   string `json:"nft_contract_addr" gorm:"type:char(42) ;comment:'contract address'"`
	Tokenid    string `json:"nft_token_id" gorm:"type:char(42) ;comment:'snft token id'"`
	Nftaddr    string `json:"nft_address" gorm:"type:char(42) ;comment:'wormholes chain address'"`
	Count      int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	Approve    string `json:"approve" gorm:"type:longtext ;comment:'Authorize'"`
	Categories string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'snft category'"`
	//Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) ;comment:'Collection creator address'"`
	//Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'NFT collection name'"`
	Image        string `json:"asset_sample" gorm:"type:longtext ;comment:'Thumbnail binary data'"`
	Hide         string `json:"hide" gorm:"type:char(20) ;comment:'Whether to let others see'"`
	Signdata     string `json:"sig" gorm:"type:longtext ;comment:'Signature data, generated when created'"`
	Createaddr   string `json:"user_addr" gorm:"type:char(42) ;comment:'Create nft address'"`
	Verifyaddr   string `json:"vrf_addr" gorm:"type:char(42) ;comment:'Validator address'"`
	Currency     string `json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'Transaction currency'"`
	Price        uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Price at creation time'"`
	Royalty      int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	Paychan      string `json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'trading channel'"`
	TransCur     string `json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'Transaction currency'"`
	Transprice   uint64 `json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'transaction price'"`
	Transtime    int64  `json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'Last trading time'"`
	Createdate   int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft creation time'"`
	Favorited    int    `json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'Follow count'"`
	Transcnt     int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt     uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	Verified     string `json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'Whether the nft work has passed the review'"`
	Verifieddesc string `json:"Verifieddesc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'Review description: Failed review description'"`
	Verifiedtime int64  `json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'Review time'"`
	Selltype     string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft transaction type'"`
	Mintstate    string `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'minting status'"`
	Collection   string `json:"collection" gorm:"type:char(42) ;comment:'collection'"`
	Local        string `json:"local" gorm:"type:char(42) ;comment:'position in collection'"`
	Extend       string `json:"extend" gorm:"type:longtext ;comment:'extend data'"`
}

type Snfts struct {
	gorm.Model
	SnftRec
}

func (v Snfts) TableName() string {
	return "snfts"
}

type SnftCollectPeriodRec struct {
	Period  string `json:"period" gorm:"type:char(42) ;comment:'stage'"`
	Collect string `json:"collect" gorm:"type:varchar(42) ;comment:'colletion'"`
	Local   string `json:"local" gorm:"type:char(42) ;comment:'Location'"`
}

type SnftCollectPeriod struct {
	gorm.Model
	SnftCollectPeriodRec
}

func (v SnftCollectPeriod) TableName() string {
	return "snftcollectperiod"
}

type Period struct {
	ID         string `json:"id"`
	Createaddr string `json:"createaddr"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Vote       int    `json:"vote"`
	Accedvote  string `json:"accedvote"`
	Categories string `json:"categories"`
	Collect    string `json:"collect"`
	TokenID    string `json:"tokenid"`
	Extend     string `json:"extend"`
}
type SnftCollection struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Accedvote  string `json:"accedvote"`
	Categories string `json:"categories"`
	TokenID    string `json:"tokenid"`
	Snft       string `json:"snft"`
	Extend     string `json:"extend"`
}

func (nft NftDb) SnftSearch(categories, param string) ([]Nfts, error) {
	snftsearch := []Nfts{}
	//cerr := GetRedisCatch().GetCatchData("SnftSearch", categories+param, &snftsearch)
	//if cerr == nil {
	//	log.Printf("SnftSearch() cache default  time.now=%s\n", time.Now())
	//	return snftsearch, nil
	//}
	if param == "" && categories == "" {
		err := nft.db.Model(&Nfts{}).Where("snft = ? ", "").Find(&snftsearch)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			fmt.Printf("search nft err=%s", err.Error)
			return nil, ErrDataBase
		}
		for i, _ := range snftsearch {
			snftsearch[i].Image = ""
		}
		//GetRedisCatch().CatchQueryData("SnftSearch", categories+param, &snftsearch)

		return snftsearch, nil
	}
	if categories == "" {
		err := nft.db.Model(&Nfts{}).Where("name like ? and snft = ?", "%"+param+"%", "").Find(&snftsearch)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			fmt.Printf("search nft err=%s", err.Error)
			return nil, ErrDataBase
		}
		for i, _ := range snftsearch {
			snftsearch[i].Image = ""
		}
		//GetRedisCatch().CatchQueryData("SnftSearch", categories+param, &snftsearch)

		return snftsearch, nil
	} else {
		catestr := strings.Split(categories, ",")
		catesql := "select * from nfts where deleted_at IS NULL and ( "
		for i, str := range catestr {
			if i < len(catestr)-1 {
				catesql = catesql + " categories = " + "'" + str + "'" + " or "
			}
			if i == len(catestr)-1 {
				catesql = catesql + " categories =  " + "'" + str + "'"
			}
		}
		catesql += " ) and name like ?"
		err := nft.db.Raw(catesql, "%"+param+"%").Scan(&snftsearch)
		if err.Error != nil {
			fmt.Println("SnftSearch()  err=", err)
			return nil, ErrDataBase
		}
		for i, _ := range snftsearch {
			snftsearch[i].Image = ""
		}
		//GetRedisCatch().CatchQueryData("SnftSearch", categories+param, &snftsearch)

		return snftsearch, nil
	}
}
