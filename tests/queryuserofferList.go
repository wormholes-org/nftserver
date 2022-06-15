package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type UserOffer struct {
	Contract	 string `json:"nft_contract_addr"`
	Tokenid      string `json:"nft_token_id"`
	Name      	 string `json:"name"`
	Price        uint64 `json:"price"`
	Count        uint64 `json:"count"`
	Bidtime      int64  `json:"date"`
}

type ResponseOfferList struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data	[]UserOffer 	`json:"data"`
	Total_count int		`json:"total_count"`
}

func QueryUserOfferList(userAddr string, count string, index string, token string) ([]UserOffer, error) {
	url := SrcUrl + "queryUserOfferList"
	datam := make(map[string]string)
	datam["count"] = count
	datam["start_index"] = index
	datam["user_addr"] = userAddr

	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), token)
	if err != nil {
		fmt.Println("QueryUserNFTList() err=", err)
		return nil, err
	}
	var revData ResponseOfferList
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("QueryUserNFTList() Unmarshal err=", err)
		return nil, err
	}
	if revData.Code != "200" {
		return nil, errors.New(revData.Msg)
	}
	return revData.Data, nil
}
