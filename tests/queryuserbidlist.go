package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)

type UserBid struct {
	NftContractAddr string `json:"nft_contract_addr"`
	NftTokenId      string `json:"nft_token_id"`
	Name      	 	string `json:"name"`
	Price           uint64 `json:"price"`
	Count           uint64 `json:"count"`
	Date            int64  `json:"date"`
}

type ResponseBidList struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data	[]UserBid 	`json:"data"`
	Total_count int		`json:"total_count"`
}

func QueryUserBidList(userAddr string, count string, index string, token string, workKey *ecdsa.PrivateKey) ([]UserBid, error) {
	url := SrcUrl + "queryUserBidList"
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
	var revData ResponseBidList
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
