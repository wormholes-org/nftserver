package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestNftVrf(t *testing.T) {
	fmt.Println("start Test TestNftVrf.")
	privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestNftVrf() HexToECDSA() err=", err)
		return
	}
	err = NftVrf(privateKey)
	if err != nil {
		fmt.Println("TestNftVrf() err=", err)
	}
	fmt.Println("end test TestNftVrf().")
}
