package main

import (
	"fmt"
	"github.com/nftexchange/nftserver/models"
	"testing"
)

func TestQueryNFTList(t *testing.T) {
	//_, tokens, err := Mlogin(1)
	//if err != nil {
	//	fmt.Println("TestQueryNFTList() err=", err)
	//	return
	//}
	fmt.Println("TestQueryNFTList() login end.")
	fmt.Println("start Test TestQueryNFTList.")
	_, err := QueryNFTList([]models.StQueryField{{Field: "selltype", Operation: "=", Value: "HighestBid"},},
	[]models.StSortField{{By: "verifiedtime", Order: "desc"}}, "0", "20")
	if err != nil {
		fmt.Println("TestQueryNFTList() err=", err)
	}
	fmt.Println("end test TestQueryNFTList().")
}