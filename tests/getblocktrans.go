package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type NftTx struct {
	Operator         string
	From             string
	To               string
	Contract         string
	TokenId          string
	Value            string
	Price            string
	Ratio            string
	TxHash           string
	Ts               string
	BlockNumber      string
	TransactionIndex string
	MetaUrl          string
	NftAddr          string
	Nonce            string
	Status           bool
	TransType        int
}

type ResponseGetBlockTrans struct {
	Code       string  `json:"code"`
	Msg        string  `json:"msg"`
	Data       []NftTx `json:"data"`
	TotalCount uint64  `json:"total_count"`
}

func GetBlockTrans(blocknumber string) (*[]NftTx, error) {
	//url := "http://192.168.56.1:8081/" + "/v1/getBlockTrans"
	url := "http://192.168.1.237:8089/" + "/v1/getBlockTrans"
	datam := make(map[string]string)
	datam["blocknumber"] = blocknumber
	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("GetBlockTrans() err=", err)
		return nil, err
	}
	var revData ResponseGetBlockTrans
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("QueryNFT() Unmarshal err=", err)
		return nil, err
	}
	if revData.Code != "200" {
		return nil, errors.New(revData.Msg)
	}
	return &revData.Data, nil
}
