package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

//Querying information about a single SNFT fragment
func (nft *NftExchangeControllerV2) QueryOwnerLevelSnfts() {
	fmt.Println("QueryOwnerLevelSnfts()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryOwnerLevelSnfts() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	defer nft.Ctx.Request.Body.Close()
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		inputDataErr := nft.verifyInputData_QueryOwnerLevelSnfts(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			collections, recount, err := nd.QueryOwnerLevelSnfts(data["owner_addr"], data["sell_type"], data["snft_level"], data["start_index"], data["count"])
			if err != nil {
				if err == gorm.ErrRecordNotFound || err == models.ErrNftNotExist {
					httpResponseData.Code = "200"
				} else {
					httpResponseData.Code = "500"
				}
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				httpResponseData.Code = "200"
				httpResponseData.TotalCount = uint64(recount)
				httpResponseData.Data = collections
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryOwnerLevelSnfts()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryOwnerLevelSnfts(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["start_index"] != "" {
		match := regNumber.MatchString(data["start_index"])
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
	if data["owner_addr"] != "" {
		match := regString.MatchString(data["owner_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["snft_level"] != "" {
		if data["snft_level"] != "*" {
			match := regNumber.MatchString(data["snft_level"])
			log.Println("verifyInputData_QueryOwnerLevelSnfts() data[\"snft_level\"]")
			if !match {
				return ERRINPUTINVALID
			}
		}
	}
	if data["sell_type"] != "" {
		if data["sell_type"] != "*" {
			match := regString.MatchString(data["sell_type"])
			if !match {
				log.Println("verifyInputData_QueryOwnerLevelSnfts() data[\"sell_type\"]")
				return ERRINPUTINVALID
			}
		}
	}
	return nil
}
