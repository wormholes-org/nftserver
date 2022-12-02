package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"regexp"
	"time"
)

//Query snft in the collection according to the creator, collection name with deleted
func (nft *NftExchangeControllerV2) QueryCollectionAllSnft() {
	fmt.Println("QuerySnftByCollection()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QuerySnftByCollection() connect database err = %s\n", err)
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
		inputDataErr := nft.verifyInputData_QueryCollectionAllSnft(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			snfts, err := nd.QueryCollectionAllSnft(data["createaddr"], data["name"])
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
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QuerySnftByCollection()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryCollectionAllSnft(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["start_index"] != "" {
		match := regNumber.MatchString(data["start_index"])
		if !match {
			fmt.Println("verifyInputData_QuerySnftByCollection createaddr error", data["start_index"])
			return ERRINPUTINVALID
		}
	}
	if data["count"] != "" {
		match := regNumber.MatchString(data["count"])
		if !match {
			fmt.Println("verifyInputData_QuerySnftByCollection createaddr error", data["count"])
			return ERRINPUTINVALID
		}
	}
	if data["createaddr"] != "" {
		match := regString.MatchString(data["createaddr"])
		if !match {
			fmt.Println("verifyInputData_QuerySnftByCollection createaddr error", data["createaddr"])
			return ERRINPUTINVALID
		}
	}
	/*if data["name"] != "" {
		match := regString.MatchString(data["name"])
		if !match {
			fmt.Println("verifyInputData_QuerySnftByCollection name error", data["name"])
			return ERRINPUTINVALID
		}
	}*/
	return nil
}
