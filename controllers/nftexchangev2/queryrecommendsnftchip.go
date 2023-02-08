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
func (nft *NftExchangeControllerV2) QueryRecommendSnftChip() {
	fmt.Println("QueryRecommendSnftChip()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryRecommendSnftChip() connect database err = %s\n", err)
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
		inputDataErr := nft.verifyInputData_QueryRecommednSnftChip(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			snfts, count, err := nd.QueryRecommendSnftChip(data["user_addr"], data["nft_contract_addr"], data["nft_token_id"])
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
				httpResponseData.Data = snfts
				httpResponseData.TotalCount = count
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryRecommendSnftChip()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryRecommednSnftChip(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	//regNumber, _ := regexp.Compile(PattenNumber)

	//if data["start_index"] != "" {
	//	match := regNumber.MatchString(data["start_index"])
	//	if !match {
	//		log.Println("start_index input error")
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["count"] != "" {
	//	match := regNumber.MatchString(data["count"])
	//	if !match {
	//		log.Println("count input error")
	//		return ERRINPUTINVALID
	//	}
	//}
	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nft_contract_addr"] != "" {
		match := regString.MatchString(data["nft_contract_addr"])
		if !match {
			log.Println("nft_contract_addr input error")
			return ERRINPUTINVALID
		}
	}
	if data["nft_token_id"] != "" {
		match := regString.MatchString(data["nft_token_id"])
		if !match {
			log.Println("nft_token_id input error")
			return ERRINPUTINVALID
		}
	}
	return nil
}
