package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nftexchange/nftserver/models"
)

type NftInfo struct {
	Ownaddr			string		`json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:''"`
	Md5				string		`json:"md5" gorm:"type:longtext NOT NULL;comment:''"`
	Name			string 		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:''"`
	Desc			string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:'：'"`
	Meta			string		`json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  NOT NULL;comment:''"`
	Url				string		`json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:''"`
	Contract		string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:''"`
	Tokenid			string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:''"`
	Count	     	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:''"`
	Approve			string		`json:"approve" gorm:"type:longtext NOT NULL;comment:''"`
	Categories		string 		`json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:''"`
	Collectcreator	string		`json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:''"`
	Collections 	string  	`json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:''"`
	Image			string		`json:"asset_sample" gorm:"type:longtext NOT NULL;comment:''"`
	Hide			string		`json:"hide" gorm:"type:char(20) NOT NULL;comment:''"`
	Signdata		string		`json:"sig" gorm:"type:longtext NOT NULL;comment:''"`
	Createaddr		string		`json:"user_addr" gorm:"type:char(42) NOT NULL;comment:''"`
	Verifyaddr		string		`json:"vrf_addr" gorm:"type:char(42) NOT NULL;comment:''"`
	Currency    	string 		`json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:''"`
	Price			uint64		`json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:''"`
	Royalty     	int 		`json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:''"`
	Paychan    		string 		`json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:''"`
	TransCur    	string 		`json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:''"`
	Transprice		uint64		`json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:''"`
	Transtime		int64		`json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:''"`
	Createdate		int64		`json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:''"`
	Favorited		int			`json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:''"`
	Transcnt		int			`json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:''"`
	Transamt		uint64		`json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:''"`
	Verified		string 		`json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:''"`
	Verifiedtime	int64		`json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:''"`
	Selltype    	string      `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:''"`
	//Sellprice		uint64		`json:"sellingprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:''"`
	Mintstate   	string      `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:''"`
	Extend			string		`json:"extend" gorm:"type:longtext NOT NULL;comment:''"`
	Sellprice		uint64		`json:"sellprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:''"`
	Offernum		uint64		`json:"offernum" gorm:"type:bigint unsigned DEFAULT NULL;comment:''"`
	Maxbidprice		uint64		`json:"maxbidprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:''"`
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
