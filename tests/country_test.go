package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestQueryCountrys(t *testing.T) {
	fmt.Println("start Test TestQueryCountrys.")
	err := QueryCountry()
	if err != nil {
		fmt.Println("TestQueryCountrys() err=", err)
	}
	fmt.Println("start Test TestKycAuditSysParams.")
	privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestQueryCountrys() HexToECDSA() err=", err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("TestQueryCountrys() err=", err)
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(fromAddress.String())
	err = ModdifyCountrys("China", "中国", "CN", "86", privateKey)
	if err != nil {
		t.Fatal("TestQueryCountrys() err=", err)
	}
	err = ModdifyCountrys("United States of America", "美国", "US", "1", privateKey)
	if err != nil {
		t.Fatal("TestQueryCountrys() err=", err)
	}
	err = QueryCountry()
	if err != nil {
		fmt.Println("TestQueryCountrys() err=", err)
	}
	fmt.Println("end test TestQueryCountrys().")
}
