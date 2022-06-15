package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type NftAuction struct {
	Selltype        string `json:"selltype"`
	Ownaddr         string `json:"ownaddr"`
	NftTokenId      string `json:"nft_token_id"`
	NftContractAddr string `json:"nft_contract_addr"`
	Paychan         string `json:"paychan"`
	Currency        string `json:"currency"`
	Startprice      uint64 `json:"startprice"`
	Endprice        uint64 `json:"endprice"`
	Startdate       int64  `json:"startdate"`
	Enddate         int64  `json:"enddate"`
	Tradesig       	string `json:"tradesig"`
}

type NftTran struct {
	NftContractAddr string `json:"nft_contract_addr"`
	Fromaddr        string `json:"fromaddr"`
	Toaddr          string `json:"toaddr"`
	NftTokenId      string `json:"nft_token_id"`
	Transtime       int64  `json:"transtime"`
	Paychan         string `json:"paychan"`
	Currency        string `json:"currency"`
	Price           uint64 `json:"price"`
	Txhash			string `json:"trade_hash"`
	Selltype        string `json:"selltype"`
}

type NftBid struct {
	Bidaddr         string `json:"bidaddr"`
	NftTokenId      string `json:"nft_token_id"`
	NftContractAddr string `json:"nft_contract_addr"`
	Paychan         string `json:"paychan"`
	Currency        string `json:"currency"`
	Price           uint64 `json:"price"`
	Bidtime         int64  `json:"bidtime"`
	Tradesig       	string `json:"tradesig"`
}

type NftSingleInfo struct {
	Name            string 			`json:"name"`
	CreatorAddr     string 			`json:"creator_addr"`
	//CreatorPortrait string 			`json:"creator_portrait"`
	OwnerAddr       string 			`json:"owner_addr"`
	//OwnerPortrait   string 			`json:"owner_portrait"`
	Md5             string 			`json:"md5"`
	//AssetSample     string 			`json:"asset_sample"`
	Desc            string 			`json:"desc"`
	Collectiondesc  string 			`json:"collection_desc"`
	Meta            string 			`json:"meta"`
	SourceUrl       string 			`json:"source_url"`
	NftContractAddr string 			`json:"nft_contract_addr"`
	NftTokenId      string 			`json:"nft_token_id"`
	Categories      string 			`json:"categories"`
	CollectionCreatorAddr string    `json:"collection_creator_addr"`
	Collections     string 			`json:"collections"`
	//Img             string 			`json:"img"`
	Approve         string 			`json:"approve"`
	Royalty         int 			`json:"royalty"`
	Verified        string 			`json:"verified"`
	Selltype        string 			`json:"selltype"`
	Mintstate       string	 		`json:"mintstate"`
	Likes	        int 			`json:"likes"`

	Auction 		NftAuction		`json:"auction"`
	Trans   		[]NftTran		`json:"trans"`
	Bids    		[]NftBid	 	`json:"bids"`
}

type ResponseNft struct {
	Code string			`json:"code"`
	Msg string			`json:"msg"`
	//Data interface{}	`json:"data"`
	Data NftSingleInfo	`json:"data"`
	TotalCount uint64	`json:"total_count"`
}

func QueryNFT(contract, tokenId string, token string) (*NftSingleInfo, error) {
	url := SrcUrl + "queryNFT"
	datam := make(map[string]string)
	datam["nft_contract_addr"] = contract
	datam["nft_token_id"] = tokenId

	datas, _ := json.Marshal(&datam)
	b, err := HttpSendRev(url, string(datas), token)
	if err != nil {
		fmt.Println("QueryNFT() err=", err)
		return nil, err
	}
	var revData ResponseNft
	err = json.Unmarshal([]byte(b), &revData)
	if err !=nil {
		fmt.Println("QueryNFT() Unmarshal err=", err)
		return nil, err
	}
	if revData.Code != "200" {
		return nil, errors.New(revData.Msg)
	}
	return &revData.Data, nil
}
