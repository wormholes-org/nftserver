package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)

func UpLoadFile(file string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "modifyAdmins"
	datam := make(map[string]string)
	datam["adminaddr"] = file


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

