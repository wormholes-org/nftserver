package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)


func AuditKYC(workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "userKYC"
	datam := make(map[string]string)
	datam["vrf_addr"] = "0xBAaeeab54cDFF708a8dCc51F56f4e2A4CE7c2ABc"
	datam["user_addr"] = "0x2d81524c4e6443b8795c123c8ac2c2b64ce12c75"
	datam["desc"] = "test."
	datam["kyc_res"] = "110000000"

	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("AuditKYC() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, "")
	if err != nil {
		fmt.Println("AuditKYC() err=", err)
		return err
	}
	b = DelDataItem(b)
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
