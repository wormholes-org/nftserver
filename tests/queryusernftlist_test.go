package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"testing"
)

func TestQueryUserNFTList(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestQueryUserNFTList() err=", err)
		return
	}
	fmt.Println("TestQueryUserNFTList() login end.")
	fmt.Println("start Test QueryUserNFTList.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			_, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("QueryUserNFTList() err=", err)
			}
		}(i)
	}
	wd.Wait()
	fmt.Println("end test QueryUserNFTList().")
}