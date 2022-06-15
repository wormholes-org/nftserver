package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestQuerySysParams(t *testing.T) {
	_, _, err := Mlogin(1)
	if err != nil {
		fmt.Println("TestQuerySysParams() err=", err)
		return
	}
	fmt.Println("TestQuerySysParams() login end.")
	fmt.Println("start Test TestQuerySysParams.")
	_, err = QuerySysParams()
	if err != nil {
		fmt.Println("TestQuerySysParams() err=", err)
	}
	fmt.Println("end test TestQuerySysParams().")
}


func TestModifySysParams(t *testing.T) {
	_, _, err := Mlogin(1)
	if err != nil {
		fmt.Println("TestQuerySysParams() err=", err)
		return
	}
	fmt.Println("TestQuerySysParams() login end.")
	fmt.Println("start Test TestQuerySysParams.")
	privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("SendTrans() err=", err)
		return
	}
	err = ModdifySysParams(privateKey)
	if err != nil {
		fmt.Println("TestQuerySysParams() err=", err)
	}
	fmt.Println("end test TestQuerySysParams().")
}
