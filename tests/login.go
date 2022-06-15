package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
)

type ResponseLogin struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data struct {
		Hash 		string 	`json:"hash"`
		Secret 		string 	`json:"secret"`
		TimeStamp 	int64 	`json:"time_stamp"`
		Token 		string 	`json:"token"`
	} `json:"data"`
	Total_count int	`json:"total_count"`
}

func Login(logKey, workKey *ecdsa.PrivateKey) (string, error) {
	url := SrcUrl + "login"
	datam := make(map[string]string)
	datam["user_addr"] = crypto.PubkeyToAddress(logKey.PublicKey).String()
	fmt.Println("login() user_prv=", hex.EncodeToString(crypto.FromECDSA(logKey)))
	fmt.Println("login() user_addr=", crypto.PubkeyToAddress(logKey.PublicKey).String())
	datam["approve_addr"] = crypto.PubkeyToAddress(workKey.PublicKey).String()
	datam["result"] = ""
	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), "")
	if err != nil {
		fmt.Println("login() err=", err)
	}
	b = DelDataItem(b)
	var revData ResponseLogin
	err = json.Unmarshal(b, &revData)
	if err != nil {
		fmt.Println("login() err=", err)
		return "", err
	}
	if revData.Code != "200" {
		return "", errors.New(revData.Msg)
	}
	i := GetMd5Index(revData.Data.Secret, revData.Data.Hash)
	datam = make(map[string]string)
	datam["user_addr"] = crypto.PubkeyToAddress(logKey.PublicKey).String()
	datam["approve_addr"] = crypto.PubkeyToAddress(workKey.PublicKey).String()
	datam["result"] = strconv.Itoa(i)
	datam["time_stamp"] = strconv.FormatInt(revData.Data.TimeStamp, 10)

	datas, _ = json.Marshal(&datam)
	data, err := HttpSendSign(datas, logKey)
	if err != nil {
		fmt.Println("HttpSendSign() err=", err)
		return "", err
	}
	b, err = HttpSendRev(url, data, "")
	if err != nil {
		fmt.Println("HttpSendRev() err=", err)
		return "", err
	}
	b = DelDataItem(b)
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return "", err
	}
	if revData.Code != "200" {
		return "", errors.New(revData.Msg)
	}
	return revData.Data.Token, nil
}

func Mlogin(c int) ([]UserKeys, []string, error) {
	tKey, err := GetUserAddr(c)
	if err != nil {
		return nil, nil, err
	}
	wd := sync.WaitGroup{}
	tokens := make([]string, c)
	for i := 0; i < c; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			loginAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			fmt.Println( "loginAddr=", loginAddr)
			token, err := Login(tKey[i].LogKey, tKey[i].WorkKey)
			if err != nil {
				fmt.Println("login err=", err)
			}
			tokens[i] = token

		}(i)
	}
	wd.Wait()
	return tKey, tokens, nil
}

func Slogin(c int) ([]UserKeys, []string, error) {
	tKey, err := GetUserAddr(c)
	if err != nil {
		return nil, nil, err
	}
	tokens := make([]string, c)
	for i := 0; i < c; i++ {
		loginAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
		fmt.Println( "loginAddr=", loginAddr)
		token, err := Login(tKey[i].LogKey, tKey[i].WorkKey)
		if err != nil {
			fmt.Println("login err=", err)
		}
		tokens[i] = token

	}
	return tKey, tokens, nil
}
