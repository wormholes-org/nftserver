package main

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/models"
)

func QueryNftInfo() error {
	//url := SrcUrl + "queryHomePage"
	//url := "http://api.wormholestest.com:8666/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/00"
	url := "http://api.wormholestest.com:8666/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/02"
	//url = "http://192.168.1.237:11002/nft_src/044bdc18-bdfb-46b9-9743-c969a6264b19.png"
	//datam := make(map[string]string)
	//
	//datas, _ := json.Marshal(&datam)
	b, err := HttpGetSendRev(url, "", "")
	if err != nil {
		fmt.Println("QueryNftInfo() err=", err)
		return err
	}
	//b = DelDataItem(b)
	ImageDir := "./"
	imagerr := models.SaveNftImage(ImageDir, "test", "01", string(b))
	if imagerr != nil {
		fmt.Println("err=", err)
	}
	var snft models.SnftInfo
	err = json.Unmarshal([]byte(b), &snft)
	if err != nil {
		fmt.Println("QueryNftInfo() get resp failed, err", err)
		return err
	}
	return nil
}

