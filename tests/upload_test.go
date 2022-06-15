package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"testing"
)

func TestUpload(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestUpload() err=", err)
		return
	}
	fmt.Println("TestUpload() login end.")
	fmt.Println("start Test Upload.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			//err := NewCollect("collect_test_" + strconv.Itoa(i), userAddr, tokens[i], tKey[i].WorkKey)
			meta := "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm"
			err := Upload(userAddr, "collect_test_" + strconv.Itoa(i), "upload_test_" + strconv.Itoa(i), tokens[i], meta, tKey[i].WorkKey)
			if err != nil {
				fmt.Println("upload() err=", err)
			}
		}(i)
	}
	wd.Wait()
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestUpload() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestUpload() login err=", err)
		return
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	meta := "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm"
	err = Upload(useraddr, "collect_test_batch", "upload_test_batch", batchToken, meta, batchkey)
	if err != nil {
		fmt.Println("upload() err=", err)
	}
	fmt.Println("end test upload().")
}

func TestUploadOne(t *testing.T) {
	//batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestNewCollections() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestNewCollections() login err=", err)
		return
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	err = NewCollect("collect_test_batch", useraddr, batchToken, batchkey)
	if err != nil {
		fmt.Println("TestNewCollections() err=", err)
	}
	meta := "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm"
	err = Upload(useraddr, "collect_test_batch", "upload_test_batch", batchToken, meta, batchkey)
	if err != nil {
		fmt.Println("upload() err=", err)
	}
	fmt.Println("end test upload().")
}



func TestUploadOneNftImage(t *testing.T) {
	//batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestNewCollections() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestNewCollections() login err=", err)
		return
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	err = NewCollect("collect_test_batch", useraddr, batchToken, batchkey)
	if err != nil {
		fmt.Println("TestNewCollections() err=", err)
	}
	meta := "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm"
	err = UploadImage(useraddr, "collect_test_batch", "upload_test_batch", batchToken, meta, batchkey)
	if err != nil {
		fmt.Println("upload() err=", err)
	}
	fmt.Println("end test upload().")
}

func TestUserUpload(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestUpload() err=", err)
		return
	}
	fmt.Println("TestUpload() login end.")
	fmt.Println("start Test Upload.")
	userAddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	//err := NewCollect("collect_test_" + strconv.Itoa(i), userAddr, tokens[i], tKey[i].WorkKey)
	meta := "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm"
	err = Upload(userAddr, "collect_test_1", "upload_test_1", tokens[1], meta, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("upload() err=", err)
	}
}