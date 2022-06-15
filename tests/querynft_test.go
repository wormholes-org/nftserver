package main

import (
	"fmt"
	"testing"
)

func TestQueryNFT(t *testing.T) {
	_, tokens, err := Mlogin(1)
	if err != nil {
		fmt.Println("TestQueryUserNFTList() err=", err)
		return
	}
	fmt.Println("TestQueryUserNFTList() login end.")
	fmt.Println("start Test QueryUserNFTList.")
	_, err = QueryNFT("0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "2339705851328", tokens[0])
	if err != nil {
		fmt.Println("QueryUserNFTList() err=", err)
	}
	fmt.Println("end test QueryUserNFTList().")
}