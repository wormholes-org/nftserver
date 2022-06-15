package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResponseHomePage struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	//Data        string `json:"data"`
	Total_count int    `json:"total_count"`
}

func QueryHomePage() error {
	url := SrcUrl + "queryHomePage"
	//datam := make(map[string]string)
	//
	//datas, _ := json.Marshal(&datam)
	b, err := HttpGetSendRev(url, "", "")
	if err != nil {
		fmt.Println("QueryHomePage() err=", err)
		return err
	}
	//b = DelDataItem(b)
	var revData ResponseHomePage
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("QueryHomePage() get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

