package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"testing"
)

func TestNewCollections(t *testing.T) {
	testCount := 200
	tKey, tokens, err := Slogin(testCount)
	if err != nil {
		fmt.Println("TestNewCollections err=", err)
		return
	}
	fmt.Println("login end.")
	fmt.Println("start Test NewCollections.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			err := NewCollect("collect_test_" + strconv.Itoa(i), userAddr, tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestNewCollections() err=", err)
			}
		}(i)
	}
	wd.Wait()
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
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
	fmt.Println("end test NewCollections.")
}