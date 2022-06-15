package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"testing"
)

func TestLogin(t *testing.T) {
	logKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("failed GenerateKey with.", err)
	}
	workKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("failed GenerateKey with.", err)
	}
	url := SrcUrl + "/login"
	datam := make(map[string]string)
	datam["user_addr"] = crypto.PubkeyToAddress(logKey.PublicKey).String()
	datam["approve_addr"] = crypto.PubkeyToAddress(workKey.PublicKey).String()
	datam["result"] = ""
	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("login() err=", err)
	}
	fmt.Println(string(b))
	b = DelDataItem(b)
	var revData ResponseLogin
	err = json.Unmarshal(b, &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return
	}
	i := GetMd5Index(revData.Data.Secret, revData.Data.Hash)
	fmt.Println("i=", i)
	datam = make(map[string]string)
	datam["user_addr"] = crypto.PubkeyToAddress(logKey.PublicKey).String()
	datam["approve_addr"] = crypto.PubkeyToAddress(workKey.PublicKey).String()
	datam["result"] = strconv.Itoa(i)
	datam["time_stamp"] = strconv.FormatInt(revData.Data.TimeStamp, 10)


	datas, _ = json.Marshal(&datam)
	data, err := HttpSendSign(datas, logKey)
	if err != nil {
		fmt.Println("sign err=", err)
	}
	b, err = HttpSendRev(url, data, "")
	if err != nil {
		fmt.Println("login() err=", err)
	}
	fmt.Println(string(b))
	b = DelDataItem(b)
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return
	}

	fmt.Println("i=", i)
	datam = make(map[string]string)

	datam["portrait"] = ""
	datam["user_addr"] = crypto.PubkeyToAddress(logKey.PublicKey).String()
	datam["user_info"] = "demo"
	datam["user_mail"] = "test@test.com"
	datam["user_name"] = "test"


	datas, _ = json.Marshal(&datam)
	data, err = HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("sign err=", err)
	}

	b, err = HttpSendRev("http://127.0.0.1:8081/v2/modifyUserInfo", data, revData.Data.Token)
	if err != nil {
		fmt.Println("modifyUserInfo() err=", err)
	}
	b = DelDataItem(b)
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return
	}
}

func ModifyUser(url string, token string, logKey, workKey *ecdsa.PrivateKey) error {
	datam := make(map[string]string)

	datam["portrait"] = ""
	datam["user_addr"] = crypto.PubkeyToAddress(logKey.PublicKey).String()
	datam["user_info"] = "demo"
	datam["user_mail"] = "test@test.com"
	datam["user_name"] = "test"


	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("modifyUserInfo() err=", err)
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

func TestMlogin(t *testing.T) {
	type testKey struct {
		logKey *ecdsa.PrivateKey
		workKey *ecdsa.PrivateKey
	}
	testCount := 10
	tKey, err := GetUserAddr(testCount)
	if err != nil {
		fmt.Println("TestMlogin() err=", err)
		return
	}
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			loginAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			fmt.Println("i=", i, "loginAddr=", loginAddr)
			_, err := Login(tKey[i].LogKey, tKey[i].WorkKey)
			if err != nil {
				fmt.Println("login err=", err)
			}
		}(i)
	}
	wd.Wait()
	fmt.Println("login test end.")
}

func TestLoginOne(t *testing.T) {
	batchkey, err := crypto.HexToECDSA("ca03d6e46914b47d7dbd4fa02cc28f8e209ae903443016218d5345b1a2d5fb21")
	if err != nil {
		fmt.Println("TestLoginOne() key err=", err)
		return
	}
	_, err = Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestLoginOne() login err=", err)
		return
	}
}