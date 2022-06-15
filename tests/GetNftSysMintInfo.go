package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type NftSysMintInfo struct {
	User_addr         string `json:"user_addr"`
	Creator_addr      string `json:"creator_addr"`
	Owner_addr        string `json:"owner_addr"`
	Md5               string `json:"md5"`
	Name              string `json:"name"`
	Desc              string `json:"desc"`
	Meta              string `json:"meta"`
	Source_url        string `json:"source_url"`
	Nft_contract_addr string `json:"nft_contract_addr"`
	Nft_token_id      string `json:"nft_token_id"`
	Categories        string `json:"categories"`
	Collections       string `json:"collections"`
	Asset_sample      string `json:"asset_sample"`
	Hide              string `json:"hide"`
	Royalty           string `json:"royalty"`
	Count             string `json:"count"`
}

type ResponseNftInfo struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data struct {
		NftInfo []NftSysMintInfo	`json:"nft_info"`
	} `json:"data"`
	Total_count int	`json:"total_count"`
}
func GetNftSysMintInfo(blockNumber uint64) error {
	//url := SrcUrl + "v2/querymetaurl"
	srcurl := "http://192.168.56.1:8080/" + "v2/querymetaurl/"
	bn  := strconv.FormatUint(blockNumber, 10)
	srcurl = srcurl + "blocknumber=" + bn
	b, err := HttpGetSendRev(srcurl, "", "")
	if err != nil {
		fmt.Println("GetNftSysMintInfo() err=", err)
		return err
	}
	//b = DelDataItem(b)
	var revData ResponseNftInfo
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("AuditKYC() get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}
