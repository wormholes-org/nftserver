package models

import (
	"testing"
	"fmt"
)

func TestBatchUploadNft(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	RoyaltyLimit = 10000
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			NFT1155Addr = "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
		} else {
			NFT1155Addr = "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F"
		}
		err = nd.UploadNft(
			"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
			"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
			"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
			"md5 string",
			"name string",
			"desc string",
			"meta string",
			"source_url string",
			"",
			"",
			"Art",
			"test",
			Default_image,
			"true",
			"2",
			"1",
			"sig string")
		if err != nil {
			fmt.Println("UploadNft err=", err)
		}
	}
}