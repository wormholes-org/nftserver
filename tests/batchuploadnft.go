package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"time"
)

//func main()  {
//	batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
//	if err != nil {
//		fmt.Println("TestNewCollections() key err=", err)
//		return
//	}
//	batchToken, err := Login(batchkey, batchkey)
//	if err != nil {
//		fmt.Println("TestNewCollections() login err=", err)
//		return
//	}
//	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
//	//err = NewCollect("collect_test_batch", useraddr, batchToken, batchkey)
//	//if err != nil {
//	//	fmt.Println("TestNewCollections() err=", err)
//	//}
//	meta := "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm"
//	st := time.Now()
//	fmt.Printf("UpLoadNft start time=%s \n", time.Now())
//	for i := 0; i < 100000000; i++ {
//		t := time.Now()
//		err = Upload(useraddr, "collect_test_batch", "upload_test_batch", batchToken, meta, batchkey)
//		if err != nil {
//			fmt.Println("upload() err=", err)
//		}
//		fmt.Printf("UpLoadNft i=%-9d spend time=%s time.now=%s\n", i, time.Now().Sub(t), time.Now())
//	}
//	fmt.Printf("UpLoadNft spend time=%s end time =%s\n", time.Now().Sub(st), time.Now())
//	fmt.Println("end test upload().")
//}

func main()  {
	testCount := 1
	upLoadCount := 100000000
	//upLoadCount := 100
	//SrcUrl = "http://192.168.1.235:9006/c0x70eb3d3f80b577e9c3954d04b787c40b763a369b/v2/"
	//SrcUrl = "http://192.168.1.235:9006/c0xbc8ac1fe086809fdaab2568dd3e8025218a62bb5/v2/"
	SrcUrl = "http://192.168.1.235:10582/v2/"
	tKey, tokens, err := Slogin(testCount)
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
			err := NewCollect("collect_test_" + strconv.Itoa(i), userAddr, tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("NewCollect() err=", err)
				//return
			}
			meta := "/ipfs/QmaNavMMmsgCyxtHMUhNM3pNfHqiiaJX4mrutXM6RF8JPy"
			for n := 0; n < upLoadCount/testCount; n++ {
				st := time.Now()
				err := Upload(userAddr, "collect_test_" + strconv.Itoa(i), "upload_test_" + strconv.Itoa(i), tokens[i], meta, tKey[i].WorkKey)
				if err != nil {
					fmt.Println("upload(i)", i, " err=", err)
					return
				}
				fmt.Printf("userAddr=%s UpLoadNft spend time=%s end time =%s\n", userAddr, time.Now().Sub(st), time.Now())
			}

		}(i)
	}
	wd.Wait()
}

func mainnew()  {
	testCount := 10
	upLoadCount := 100000000
	//upLoadCount := 10
	//upLoadCount := 100
	//SrcUrl = "http://192.168.1.235:9006/c0xbc8ac1fe086809fdaab2568dd3e8025218a62bb5/v2/"
	//SrcUrl = "http://192.168.1.235:9006/c0x70eb3d3f80b577e9c3954d04b787c40b763a369b/v2/"
	SrcUrl = "http://192.168.1.235:10582/v2/"
	tKey, tokens, err := Slogin(testCount)
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
			err := NewCollect("collect_test_" + strconv.Itoa(i), userAddr, tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("NewCollect() err=", err)
				//return
			}
			nftmeta := GetTestNfts(i)
			for n := 0; n < upLoadCount/testCount; n++ {
				st := time.Now()
				err := UploadWithImage(userAddr, "collect_test_" + strconv.Itoa(i), "upload_test_" + strconv.Itoa(i), tokens[i], nftmeta.Meta, nftmeta.Image, tKey[i].WorkKey)
				if err != nil {
					fmt.Println("upload(i)", i, " err=", err)
					return
				}
				fmt.Printf("userAddr=%s UpLoadNft spend time=%s end time =%s\n", userAddr, time.Now().Sub(st), time.Now())
			}

		}(i)
	}
	wd.Wait()
}
