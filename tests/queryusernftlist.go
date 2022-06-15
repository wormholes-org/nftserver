package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)

type UserNft struct {
	UserAddr        string `json:"user_addr"`
	CreatorAddr     string `json:"creator_addr"`
	OwnerAddr       string `json:"owner_addr"`
	Md5             string `json:"md5"`
	Name            string `json:"name"`
	Desc            string `json:"desc"`
	Meta            string `json:"meta"`
	SourceUrl       string `json:"source_url"`
	NftContractAddr string `json:"nft_contract_addr"`
	NftTokenId      string `json:"nft_token_id"`
	Categories      string `json:"categories"`
	Collections     string `json:"collections"`
	//AssetSample     string `json:"asset_sample"`
	Hide            string `json:"hide"`
}

type ResponseNftList struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data	[]UserNft 	`json:"data"`
	Total_count int		`json:"total_count"`
}

func QueryUserNFTList(userAddr string, count string, index string, token string, workKey *ecdsa.PrivateKey) ([]UserNft, error) {
	url := SrcUrl + "queryUserNFTList"
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
	var revData ResponseNftList
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
