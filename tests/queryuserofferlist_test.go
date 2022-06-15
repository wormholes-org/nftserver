package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"sync"
	"testing"
)

func TestQueryUserOfferList(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestQueryUserOfferList() err=", err)
		return
	}
	fmt.Println("TestQueryUserOfferList() login end.")
	fmt.Println("start Test TestQueryUserOfferList.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			_, err := QueryUserOfferList(userAddr, strconv.Itoa(testCount), "0", tokens[i])
			if err != nil {
				fmt.Println("TestQueryUserOfferList() err=", err, "userAddr=", userAddr)
				return
			}
			fmt.Println("TestQueryUserOfferList() Ok.", "userAddr=", userAddr)
		}(i)
	}
	wd.Wait()
	fmt.Println("end test TestQueryUserOfferList().")
}

func TestQueryUserOfferListSingle(t *testing.T) {
	testCount := 10
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestQueryUserOfferListSingle() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestQueryUserOfferListSingle() login err=", err)
	}
	fmt.Println("TestQueryUserOfferListSingle() login end.")
	fmt.Println("start Test TestQueryUserOfferListSingle.")
	userAddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	_, err = QueryUserOfferList(userAddr, strconv.Itoa(testCount), "0", batchToken)
	if err != nil {
		fmt.Println("TestQueryUserOfferListSingle() err=", err, "userAddr=", userAddr)
		return
	}
	fmt.Println("TestQueryUserOfferListSingle() Ok.", "userAddr=", userAddr)
	fmt.Println("end test TestQueryUserOfferListSingle().")
}