package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nftexchange/nftserver/models"
)

type ResponseC struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data struct {
		Countrys []models.CountryRec
	} `json:"data"`
	Total_count int	`json:"total_count"`
}

func QueryCountry() error {
	url := SrcUrl + "queryCountrys"
	datam := make(map[string]string)

	datas, _ := json.Marshal(&datam)
	b, err := HttpGetSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("QueryCountry() err=", err)
		return err
	}
	//b = DelDataItem(b)
	var revData ResponseC
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("QueryCountry() get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

func ModdifyCountrys(regionen, regioncn, domain, telecode string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "modifyCountrys"
	datam := make(map[string]string)
	datam["regionen"] = regionen
	datam["regioncn"] = regioncn
	datam["domain"] = domain
	datam["telecode"] = telecode


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
	var revData ResponseC
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

