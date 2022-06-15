package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	geth "github.com/ethereum/go-ethereum/mobile"
	"github.com/ethereum/go-ethereum/params"
	"github.com/nftexchange/nftserver/common/contracts"
	"testing"
	"time"
)

const (
	TestCount = 10
	rechargeAddr = "0xBAaeeab54cDFF708a8dCc51F56f4e2A4CE7c2ABc"
	rechargePrv = "2604f045d7b600e440541113415a434fc489381053c36887cc0f1f7133abda64"
	TradeContract = "0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"

)
var DepositPrice *geth.BigInt
var ChargePrice *geth.BigInt
func init()  {
	DepositPrice = geth.NewBigInt(0.005 * params.Ether)
	ChargePrice = geth.NewBigInt(0.01 * params.Ether)
}

func TestGenUserKeys(t *testing.T) {
	//err := GenUserKeys("./key", 1000)
	//if err != nil {
	//	fmt.Println("GenUserKeys() err= ", err)
	//}
}

func TestGetUserKeys(t *testing.T) {
	userKeys, err := GetUserKeys("./key")
	if err != nil {
		fmt.Println("GenUserKeys() err= ", err)
	}
	fmt.Println(hexutil.Encode(crypto.FromECDSA(userKeys[0].LogKey)), hexutil.Encode(crypto.FromECDSAPub(&userKeys[0].LogKey.PublicKey)), crypto.PubkeyToAddress(userKeys[0].LogKey.PublicKey).String())
}

func TestInitUserAddr(t *testing.T) {
	userKeys, err := GetUserKeys("./key")
	if err != nil {
		fmt.Println("GenUserKeys() err= ", err)
	}
	for i := 0; i < TestCount; i++ {
		publicKey := userKeys[i].LogKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			fmt.Println("TestInitKeys() err=", err)
			return
		}
		to := crypto.PubkeyToAddress(*publicKeyECDSA)
		fmt.Println("TestInitKeys() rev Addr=", to)
		price := ChargePrice
		//fmt.Println(price.String())
		privateKey, err := crypto.HexToECDSA(rechargePrv)
		if err != nil {
			fmt.Println("SendTrans() err=", err)
			return
		}
		err = contracts.SendTrans(to.String(), price.String(), privateKey)
		if err != nil {
			fmt.Println("TestInitKeys() SendTrans() err=", err)
			return
		}
		time.Sleep(time.Second)
	}
}

func TestDeposit(t *testing.T) {
	userKeys, err := GetUserKeys("./key")
	if err != nil {
		fmt.Println("GenUserKeys() err= ", err)
		return
	}
	for i := 0; i < TestCount; i++ {
		publicKey := userKeys[i].LogKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			fmt.Println("TestInitKeys() err=", err)
			return
		}
		to := crypto.PubkeyToAddress(*publicKeyECDSA)
		fmt.Println("TestInitKeys() rev Addr=", to)
		//price := geth.NewBigInt(0.005 * params.Ether)
		price := DepositPrice
		err = contracts.Deposit(to.String(), price.String(), userKeys[i].LogKey)
		if err != nil {
			fmt.Println("TestInitKeys() SendTrans() err=", err)
			return
		}
		time.Sleep(time.Second)
	}
}

func TestSetApprove(t *testing.T) {
	userKeys, err := GetUserKeys("./key")
	if err != nil {
		fmt.Println("GenUserKeys() err= ", err)
	}
	for i := 0; i < TestCount; i++ {
		tx, err := contracts.SetApprove(TradeContract, userKeys[i].LogKey)
		if err != nil {
			fmt.Println("TestInitKeys() SendTrans() err=", err)
			return
		}
		fmt.Println("TestSetApprove() tx.hash=", tx.Hash())
	}
}

func TestApproveValue(t *testing.T) {
	userKeys, err := GetUserKeys("./key")
	if err != nil {
		fmt.Println("GenUserKeys() err= ", err)
	}
	price := geth.NewBigInt(0.011 * params.Ether)
	for i := 0; i < TestCount; i++ {
		err := contracts.Approve("", price.String(), userKeys[i].LogKey)
		if err != nil {
			fmt.Println("TestApproveValue() SendTrans() err=", err)
			return
		}
	}
}