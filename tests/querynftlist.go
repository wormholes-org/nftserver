package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nftexchange/nftserver/models"
)

type NftInfo struct {
	Ownaddr			string		`json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft拥有者地址'"`
	Md5				string		`json:"md5" gorm:"type:longtext NOT NULL;comment:'图片md5值'"`
	Name			string 		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft分类'"`
	Desc			string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'审核描述：未通过审核描述'"`
	Meta			string		`json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'元信息'"`
	Url				string		`json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc原始数据保持地址'"`
	Contract		string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Tokenid			string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
	Count	     	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Approve			string		`json:"approve" gorm:"type:longtext NOT NULL;comment:'授权'"`
	Categories		string 		`json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'nft分类'"`
	Collectcreator	string		`json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'合集创建者地址'"`
	Collections 	string  	`json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT合集名'"`
	Image			string		`json:"asset_sample" gorm:"type:longtext NOT NULL;comment:'缩略图二进制数据'"`
	Hide			string		`json:"hide" gorm:"type:char(20) NOT NULL;comment:'是否让其他人看到'"`
	Signdata		string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据，创建时产生'"`
	Createaddr		string		`json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'创建nft地址'"`
	Verifyaddr		string		`json:"vrf_addr" gorm:"type:char(42) NOT NULL;comment:'验证人地址'"`
	Currency    	string 		`json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Price			uint64		`json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'创建时定的价格'"`
	Royalty     	int 		`json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Paychan    		string 		`json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'交易通道'"`
	TransCur    	string 		`json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Transprice		uint64		`json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'交易成交价格'"`
	Transtime		int64		`json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'最后交易时间'"`
	Createdate		int64		`json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft创建时间'"`
	Favorited		int			`json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'被关注计数'"`
	Transcnt		int			`json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt		uint64		`json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	Verified		string 		`json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'nft作品是否通过审核'"`
	Verifiedtime	int64		`json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'审核时间'"`
	Selltype    	string      `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft交易类型'"`
	//Sellprice		uint64		`json:"sellingprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'正在销售价格'"`
	Mintstate   	string      `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'铸币状态'"`
	Extend			string		`json:"extend" gorm:"type:longtext NOT NULL;comment:'扩展字段'"`
	Sellprice		uint64		`json:"sellprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'正在销售的价格'"`
	Offernum		uint64		`json:"offernum" gorm:"type:bigint unsigned DEFAULT NULL;comment:'出价个数'"`
	Maxbidprice		uint64		`json:"maxbidprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'最高出价价格'"`
}

type ResponseNftFilter struct {
	Code string			`json:"code"`
	Msg string			`json:"msg"`
	//Data interface{}	`json:"data"`
	Data []NftInfo	`json:"data"`
	TotalCount uint64	`json:"total_count"`
}

type Filter struct {
	Field     string `json:"field"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type Sort struct {
	By    string `json:"by"`
	Order string `json:"order"`
}

type HttpRequestFilter struct {
	Match string `json:"match"`
	Filter []models.StQueryField `json:"filter"`
	Sort []models.StSortField `json:"sort"`
	StartIndex string `json:"start_index"`
	Count string `json:"count"`
}

func QueryNFTList(filter []models.StQueryField, sort []models.StSortField, start_index, count string) (*[]NftInfo, error) {
	url := SrcUrl + "queryNFTList"
	//datam := make(map[string]string)
	//f, _ := json.Marshal(&filter)
	//datam["filter"] = string(f)
	//s, _ := json.Marshal(&sort)
	//datam["sort"] = string(s)
	//datam["match"] = ""
	//datam["start_index"] = start_index
	//datam["count"] = count
	var datam HttpRequestFilter
	datam.Filter = filter
	datam.Sort = sort
	datam.Match = ""
	datam.StartIndex = "0"
	datam.Count = "20"
	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("QueryNFT() err=", err)
		return nil, err
	}
	var revData ResponseNftFilter
	err = json.Unmarshal([]byte(b), &revData)
	if err !=nil {
		fmt.Println("QueryNFT() Unmarshal err=", err)
		return nil, err
	}
	if revData.Code != "200" {
		return nil, errors.New(revData.Msg)
	}
	return &revData.Data, nil
}
