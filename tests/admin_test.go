package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

type deladmins struct {
	Admins [][]string `json:"del_admins"`
}

func TestQueryAdmin(t *testing.T) {
	fmt.Println("start Test TestQueryAdmin.")
	err := QueryAdmin("nft")
	if err != nil {
		fmt.Println("TestQueryAdmin() err=", err)
	}
	fmt.Println("start Test TestKycAuditSysParams.")
	privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestQueryAdmin() HexToECDSA() err=", err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("TestQueryAdmin() err=", err)
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(fromAddress.String())
	err = ModdifyAdmins("0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"nft", "6", privateKey)
	if err != nil {
		t.Fatal("ModdifyAdmins() err=", err)
	}
	err = ModdifyAdmins("0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"kyc", "6", privateKey)
	if err != nil {
		t.Fatal("ModdifyAdmins() err=", err)
	}
	err = ModdifyAdmins("0xbf5fabb29d464b41eaf88096654dd813ad7bcf58",
		"admin", "6", privateKey)
	if err != nil {
		t.Fatal("ModdifyAdmins() err=", err)
	}

	deladmins := [][]string{
		{"0x572bcacb7ae32db658c8dee49e156d455ad59e15", "nft"},
		{"0x6e40b6deb1671b48b8b7efecac58b21f4f96468a", "admin"},
	}
	delstr, err := json.Marshal(&deladmins)
	err = DelAdmins(string(delstr), privateKey)
	if err != nil {
		t.Fatal("DelAdmins() err=", err)
	}
	fmt.Println("end test TestQueryAdmin().")
}
