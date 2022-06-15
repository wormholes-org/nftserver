package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestUpLoadFile(t *testing.T) {
	fmt.Println("start Test TestUpLoadFile.")
	fmt.Println("start Test TestUpLoadFile.")
	privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestUpLoadFile() HexToECDSA() err=", err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("TestUpLoadFile() err=", err)
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(fromAddress.String())
	err = UpLoadFile("test.txt", privateKey)
	if err != nil {
		t.Fatal("TestUpLoadFile() err=", err)
	}
	fmt.Println("end test TestUpLoadFile().")
}
