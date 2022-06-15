package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestQueryAnnounce(t *testing.T) {
	fmt.Println("start Test TestQueryAnnounce.")
	err := QueryAnnounce("0", "2")
	if err != nil {
		fmt.Println("TestQueryAnnounce() err=", err)
	}
	fmt.Println("start Test TestQueryAnnounce.")
	privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestQueryAnnounce() HexToECDSA() err=", err)
		return
	}
	err = DelAnnounce(privateKey)
	if err != nil {
		t.Fatal("TestQueryAnnounce() err=", err)
	}
	err = ModdifyAnnounce("nft", "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900", privateKey)
	if err != nil {
		t.Fatal("TestQueryAnnounce() err=", err)
	}
	err = ModdifyAnnounce("kyc", "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900", privateKey)
	if err != nil {
		t.Fatal("TestQueryAnnounce() err=", err)
	}
	err = ModdifyAnnounce("admin", "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900", privateKey)
	if err != nil {
		t.Fatal("TestQueryAnnounce() err=", err)
	}
	err = ModdifyAnnounce("kyc", "0x0000000000000000000000000000000000000", privateKey)
	if err != nil {
		t.Fatal("TestQueryAnnounce() err=", err)
	}
	err = QueryAnnounce("2", "10")
	if err != nil {
		fmt.Println("TestQueryAnnounce() err=", err)
	}
	fmt.Println("end test TestQueryAnnounce().")
}
