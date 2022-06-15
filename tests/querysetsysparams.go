package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)

type SysParamsInfo struct {
	NFT1155addr            string `json:"nft1155addr"`
	Adminaddr              string `json:"adminaddr"`
	Lowprice               string `json:"lowprice"`
	Blocknumber            string `json:"blocknumber"`
	Scannumber             string `json:"scannumber"`
	Royaltylimit           string `json:"royaltylimit"`
	Homepage               string `json:"homepage"`
	ExchangerInfo          string `json:"exchangerInfo"`
	Icon	               string `json:"icon"`
	Data	               string `json:"data"`
	Categories             string `json:"categories"`
	Nftaudit			   string `json:"nftaudit"`
	Sig					   string `json:"sig"`
}

type ResponseParams struct {
	Code string			`json:"code"`
	Msg string			`json:"msg"`
	//Data interface{}	`json:"data"`
	Data SysParamsInfo	`json:"data"`
	TotalCount uint64	`json:"total_count"`
}

func QuerySysParams() (*SysParamsInfo, error) {
	url := SrcUrl + "querySysParams"
	datam := make(map[string]string)
	//datam["nft_contract_addr"] = contract
	//datam["nft_token_id"] = tokenId

	datas, _ := json.Marshal(&datam)
	b, err := HttpGetSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("QueryNFT() err=", err)
		return nil, err
	}
	var revData ResponseParams
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

func ModdifySysParams(workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "modifySysParams"
	datam := make(map[string]string)
	//datam["nft1155addr"] = "0xBAaeeab54cDFF708a8dCc51F56f4e2A4CE7c2ABc"
	//datam["adminaddr"] = ""
	//datam["lowprice"] = "1000"
	//datam["scannumber"] = "110000000"
	//datam["nftaudit"] = "true"
	datam["exchangerauth"] = `{"block_number":"0x936e","exchanger_owner":"0x48cae23c1e43ce233952d2b15b6461dba83767d8","sig":"0xf85c89b71602fad9fe9fc6e710df024f71708f8c3cadab7b92cfdc3f70e2ef9a774d11250068fb5e07ec2b0ab799fc93c0a9d510d47c2f081bf4d49a2567cee61c","to":"0x9147e89e031d7466a79aced31e2f8ab2e80ab7da"}`
	//datam["data"] = "Data"
	//datam["exchangerinfo"] = "交易所NO'1"
	//datam["homepage"] = models.HomePages


	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("NewCollect() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, "")
	if err != nil {
		fmt.Println("NewCollect() err=", err)
		return err
	}
	b = DelDataItem(b)
	var revData ResponseLogin
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}
