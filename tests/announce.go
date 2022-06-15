package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nftexchange/nftserver/models"
)


type ResponseAnnouce struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data 	[]models.Announces	`json:"data"`
	Total_count int		`json:"total_count"`
}

func QueryAnnounce(start_index, count string) error {
	url := SrcUrl + "queryAnnounce"
	datam := make(map[string]string)
	datam["start_index"] = start_index
	datam["count"] = count

	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("QueryAnnounce() err=", err)
		return err
	}
	//b = DelDataItem(b)
	var revData ResponseAnnouce
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("QueryAnnounce() get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

func ModdifyAnnounce(title, content string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "modifyAnnounce"
	datam := make(map[string]string)
	datam["title"] = title
	datam["content"] = content

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

func DelAnnounce(workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "delAnnounces"
	datam := make(map[string]string)
	delannouces := []int{1, 2, 3}
	str, _ := json.Marshal(&delannouces)
	datam["del_announces"] = string(str)

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
