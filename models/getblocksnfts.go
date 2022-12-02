package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*type SnftInfo struct {
	CreatorAddr          string  `json:"creator_addr"`
	Ownaddr              string  `json:"ownaddr"`
	Contract             string  `json:"nft_contract_addr"`
	Nftaddr              string  `json:"nft_address"`
	Name                 string  `json:"name"`
	Desc                 string  `json:"desc"`
	Meta                 string  `json:"meta"`
	Category             string  `json:"category"`
	Royalty              float64 `json:"royalty"`
	SourceUrl            string  `json:"source_url"`
	Md5                  string  `json:"md5"`
	CollectionsName      string  `json:"collections_name"`
	CollectionsCreator   string  `json:"collections_creator"`
	CollectionsExchanger string  `json:"collections_exchanger"`
	CollectionsCategory  string  `json:"collections_category"`
	CollectionsImgUrl    string  `json:"collections_img_url"`
	CollectionsDesc      string  `json:"collections_desc"`
}
*/

type ResponseGetBlockSnfts struct {
	Code       string     `json:"code"`
	Msg        string     `json:"msg"`
	Data       []SnftInfo `json:"data"`
	TotalCount uint64     `json:"total_count"`
}

func GetBlockSnfts(blocknumber string) ([]SnftInfo, error) {
	url := NftScanServer + "/v1/getBlockSnfts"
	datam := make(map[string]string)
	datam["blocknumber"] = blocknumber
	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("GetBlockTrans() err=", err)
		return nil, err
	}
	var revData ResponseGetBlockSnfts
	err = json.Unmarshal([]byte(b), &revData)
	if err !=nil {
		fmt.Println("QueryNFT() Unmarshal err=", err)
		return nil, err
	}
	if revData.Code != "200" {
		return nil, errors.New(revData.Msg)
	}
	return revData.Data, nil
}
