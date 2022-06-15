package nftexchangev2

import (
	"fmt"
	"github.com/nftexchange/nftserver/models"
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
