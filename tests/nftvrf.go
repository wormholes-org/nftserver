package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)


/*
"vrf_addr": "0x8d904263c5383a9c22e80e7b45d3f11f2df63d0c",
		"owner": n["ownaddr"],
		"nft_contract_addr": n["nft_contract_addr"],
		"nft_token_id": n["nft_token_id"],
		"desc": "none",
		"vrf_res": "Passed"
*/
func NftVrf(workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "vrfNFT"
	datam := make(map[string]string)
	datam["vrf_addr"] = "0xBAaeeab54cDFF708a8dCc51F56f4e2A4CE7c2ABc"
	datam["owner"] = "0x4d5bde96fe35a42bac7d2aba227207347e3cf66e"
	datam["nft_contract_addr"] = "0xa1e67a33e090afe696d7317e05c506d7687bb2e5"
	datam["nft_token_id"] = "2339705851328"
	datam["desc"] = "110000000"
	datam["vrf_res"] = "Passed"

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
