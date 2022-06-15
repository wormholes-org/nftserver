package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type SnftCollectRec struct {
	Createaddr   string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'创建者地址'"`
	Contract     string `json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Contracttype string `json:"contracttype" gorm:"type:char(20) CHARACTER SET utf8mb4 NOT NULL;comment:'合约类型'"`
	Name         string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'集合名称'"`
	Desc         string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'集合描述'"`
	Categories   string `json:"categories" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'集合分类'"`
	Totalcount   int    `json:"total_count" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'集合中nft总数'"`
	Transcnt     int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt     uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	SigData      string `json:"sig" gorm:"type:longtext NOT NULL;comment:'签名'"`
	Img          string `json:"img" gorm:"type:longtext NOT NULL;comment:'logo图片'"`
	Tokenid      string `json:"tokenid" gorm:"type:char(42) ;comment:'唯一标识标志'"`
	Snft         string `json:"snft" gorm:"type:longtext ;comment:'16个snft'"`
	//Period       string `json:"period" gorm:"type:char(42) ;comment:'期号'"`
	Local     string `json:"local" gorm:"type:char(42) ;comment:'在期中位置'"`
	Exchanger string `json:"exchanger" gorm:"type:char(42) ;comment:'snft所属交易所地址'"`
	Extend    string `json:"extend" gorm:"type:longtext ;comment:'扩展字段'"`
}

type SnftCollect struct {
	gorm.Model
	SnftCollectRec
}

func (v SnftCollect) TableName() string {
	return "snftcollect"
}

type SnftPhaseRec struct {
	Createaddr string `json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'创建者地址'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'名称'"`
	Desc       string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'描述'"`
	Vote       int    `json:"vote" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'投票数'"`
	Accedvote  string `json:"accedvote" gorm:"type:char(20) ;comment:'是否参与投票'"`
	Accedeth   string `json:"accedeth" gorm:"type:char(20) ;comment:'是否选中snft'"`
	Categories string `json:"categories" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'分类'"`
	Tokenid    string `json:"tokenid" gorm:"type:char(42) ;comment:'唯一标识标志'"`
	Collect    string `json:"collect" gorm:"type:longtext ;comment:'16个collection'"`
	Meta       string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'元信息'"`
	Extend     string `json:"extend" gorm:"type:longtext ;comment:'扩展字段'"`
}

type SnftPhase struct {
	gorm.Model
	SnftPhaseRec
}

func (v SnftPhase) TableName() string {
	return "snftphase"
}

type SnftRec struct {
	Ownaddr    string `json:"ownaddr" gorm:"type:char(42) ;comment:'nft拥有者地址'"`
	Md5        string `json:"md5" gorm:"type:longtext ;comment:'图片md5值'"`
	Name       string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft名称'"`
	Desc       string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft描述'"`
	Meta       string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'元信息'"`
	Nftmeta    string `json:"nftmeta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'元信息、tokenid'"`
	Url        string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc原始数据保持地址'"`
	Contract   string `json:"nft_contract_addr" gorm:"type:char(42) ;comment:'合约地址'"`
	Tokenid    string `json:"nft_token_id" gorm:"type:char(42) ;comment:'唯一标识nft标志'"`
	Nftaddr    string `json:"nft_address" gorm:"type:char(42) ;comment:'wormholes链唯一标识nft标志'"`
	Count      int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Approve    string `json:"approve" gorm:"type:longtext ;comment:'授权'"`
	Categories string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft分类'"`
	//Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) ;comment:'合集创建者地址'"`
	//Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'NFT合集名'"`
	Image        string `json:"asset_sample" gorm:"type:longtext ;comment:'缩略图二进制数据'"`
	Hide         string `json:"hide" gorm:"type:char(20) ;comment:'是否让其他人看到'"`
	Signdata     string `json:"sig" gorm:"type:longtext ;comment:'签名数据，创建时产生'"`
	Createaddr   string `json:"user_addr" gorm:"type:char(42) ;comment:'创建nft地址'"`
	Verifyaddr   string `json:"vrf_addr" gorm:"type:char(42) ;comment:'验证人地址'"`
	Currency     string `json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Price        uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'创建时定的价格'"`
	Royalty      int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Paychan      string `json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'交易通道'"`
	TransCur     string `json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Transprice   uint64 `json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'交易成交价格'"`
	Transtime    int64  `json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'最后交易时间'"`
	Createdate   int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft创建时间'"`
	Favorited    int    `json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'被关注计数'"`
	Transcnt     int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt     uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	Verified     string `json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'nft作品是否通过审核'"`
	Verifieddesc string `json:"Verifieddesc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'审核描述：未通过审核描述'"`
	Verifiedtime int64  `json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'审核时间'"`
	Selltype     string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft交易类型'"`
	Mintstate    string `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'铸币状态'"`
	Collection   string `json:"collection" gorm:"type:char(42) ;comment:'合集'"`
	Local        string `json:"local" gorm:"type:char(42) ;comment:'在合集中位置'"`
	Extend       string `json:"extend" gorm:"type:longtext ;comment:'扩展字段'"`
}

type Snfts struct {
	gorm.Model
	SnftRec
}

func (v Snfts) TableName() string {
	return "snfts"
}

type SnftCollectPeriodRec struct {
	Period  string `json:"period" gorm:"type:char(42) ;comment:'期号'"`
	Collect string `json:"collect" gorm:"type:varchar(42) ;comment:'合集号'"`
	Local   string `json:"local" gorm:"type:char(42) ;comment:'位置'"`
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
	if param == "" && categories == "" {
		err := nft.db.Model(&Nfts{}).Where("snft = ? ", "").Find(&snftsearch)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			fmt.Printf("search nft err=%s", err.Error)
			return nil, err.Error
		}
		for i, _ := range snftsearch {
			snftsearch[i].Image = ""
		}
		return snftsearch, nil
	}
	if categories == "" {
		err := nft.db.Model(&Nfts{}).Where("name like ? and snft = ?", "%"+param+"%", "").Find(&snftsearch)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			fmt.Printf("search nft err=%s", err.Error)
			return nil, err.Error
		}
		for i, _ := range snftsearch {
			snftsearch[i].Image = ""
		}
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
			return nil, err.Error
		}
		for i, _ := range snftsearch {
			snftsearch[i].Image = ""
		}
		return snftsearch, nil
	}
}
