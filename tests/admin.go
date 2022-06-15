package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)


func QueryAdmin(adminType string) error {
	url := SrcUrl + "queryAdmins"
	datam := make(map[string]string)
	datam["admintype"] = adminType
	datam["start_index"] = "1"
	datam["count"] = "10"

	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("AuditKYC() err=", err)
		return err
	}
	//b = DelDataItem(b)
	var revData ResponseLogin
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

func ModdifyAdmins(adminAddr, adminType, adminAuth string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "modifyAdmins"
	datam := make(map[string]string)
	datam["adminaddr"] = adminAddr
	datam["admintype"] = adminType
	datam["adminauth"] = adminAuth


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

func DelAdmins(adminAddrs string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "delAdmins"
	datam := make(map[string]string)
	datam["del_admins"] = adminAddrs


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