package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/models"
	"regexp"
	"testing"
)

//const sqlsvrLcT = "admin:user123456@tcp(192.168.1.237:3306)/"
const sqlsvrLcT = "admin:user123456@tcp(192.168.32.128:3306)/"

//
//const sqlsvrLcT = "demo:123456@tcp(192.168.56.129:3306)/"

//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
//var SqlSvrT = "admin:user123456@tcp(192.168.1.238:3306)/"
//const dbNameT = "nftdbdemo"
const dbNameT = "nftdb"
const localtimeT = "?parseTime=true&loc=Local"

//const localtimeT = "?charset=utf8mb4&parseTime=True&loc=Local"

const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

func TestVerify(t *testing.T) {
	nft := NftExchangeControllerV2{}
	nd, err := models.NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	rawData := `{"def_language":"en_us"}`
	sig := `0x29381026df3b9cb57d67eaa620c4b4ace3886b62e344586aaef09adaf484941d35e1b824be66e0ea10fbf6dd63d6a9822fbb6c855a8f9bd3922c134d3a9a0f871b`
	_, inputDatarr := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, rawData, sig)
	fmt.Println(inputDatarr)
}

func verifyInputData_UploadNft(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)
	regImage, _ := regexp.Compile(PattenImageBase64)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["creator_addr"] != "" {
		match := regString.MatchString(data["creator_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["owner_addr"] != "" {
		match := regString.MatchString(data["owner_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["name"] != "" {
	//	match := regString.MatchString(data["name"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["desc"] != "" {
	//	match := regString.MatchString(data["desc"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	if data["nft_contract_addr"] != "" {
		match := regString.MatchString(data["nft_contract_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nft_token_id"] != "" {
		match := regString.MatchString(data["nft_token_id"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["md5"] != "" {
		match := regString.MatchString(data["md5"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["categories"] != "" {
		match := regString.MatchString(data["categories"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["collections"] != "" {
	//	match := regString.MatchString(data["collections"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	match := regImage.MatchString(data["asset_sample"])
	if !match {
		return ERRINPUTINVALID
	}
	if data["hide"] != "" {
		match := regString.MatchString(data["hide"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["royalty"] != "" {
		match := regNumber.MatchString(data["royalty"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["count"] != "" {
		match := regNumber.MatchString(data["count"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["sig"] != "" {
		match := regString.MatchString(data["sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	if getToken != token {
		return ERRTOKEN
	}

	return nil
}

func TestUploadVerifyData(t *testing.T) {
	data := `{
    "creator_addr": "0xde253dbc82978f0890b5bfa6e75ecd39985c3efe",
    "nft_contract_addr": "0x82132502557b8a00d2ac2e8eb4670498ac5e32e8",
    "name": "abcee",
    "desc": "444444",
    "category": "Music",
    "royalty": "100",
    "source_image_name": "5-6练习册听力.mp3",
    "fileType": "mp3",
    "source_url": "/ipfs/QmQWLoxaL39pB7AXYTUcuAPUqN51RxWUjZjKUwR16fSxPn",
    "md5": "8232309c27b26beee23a0a570b2953fc",
    "collections_name": "测试合集",
    "collections_creator": "0xde253dbc82978f0890b5bfa6e75ecd39985c3efe",
    "collections_exchanger": "0x82132502557b8a00d2ac2e8eb4670498ac5e32e8",
    "collections_category": "Utility",
    "collections_img_url": "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH",
    "collections_desc": "",
    "attributes": "[]"
}`
	vdata := make(map[string]string)
	err := json.Unmarshal([]byte(data), &vdata)
	if err != nil {

	}
	inputDataErr := verifyInputData_UploadNft(vdata, "")
	if inputDataErr != nil {
	}
}
