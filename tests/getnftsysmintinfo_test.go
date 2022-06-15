package main

import (
	"fmt"
	"testing"
)

func TestGetNftSysMintInfo(t *testing.T) {
	fmt.Println("start Test Testgetnftinfo.")
	//privateKey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	//if err != nil {
	//	fmt.Println("Testgetnftinfo() HexToECDSA() err=", err)
	//	return
	//}
	err := GetNftSysMintInfo(10)
	if err != nil {
		fmt.Println("Testgetnftinfo() err=", err)
	}
	fmt.Println("end test Testgetnftinfo().")
}
