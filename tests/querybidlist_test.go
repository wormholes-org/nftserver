package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"testing"
)

func TestQueryUserBidList(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestQueryUserBidList() err=", err)
		return
	}
	fmt.Println("TestQueryUserBidList() login end.")
	fmt.Println("start Test TestQueryUserBidList.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			_, err := QueryUserBidList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestQueryUserBidList() err=", err, "userAddr=", userAddr)
				return
			}
			fmt.Println("TestQueryUserBidList() Ok.", "userAddr=", userAddr)
		}(i)
	}
	wd.Wait()
	fmt.Println("end test TestQueryUserBidList().")
}