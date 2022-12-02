package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

//Query single NFT information
func (nft *NftExchangeControllerV2) QueryTokenIdByNftaddr() {
	fmt.Println("QueryTokenIdByNftaddr()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryTokenIdByNftaddr() connect database err = %s\n", err)
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
		inputDataErr := nft.verifyInputData_QueryTokenIdByNftaddr(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			tokenId, err := nd.QueryTokenIdByNftaddr(data["nft_addr"])
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				httpResponseData.Code = "200"
				httpResponseData.Data = tokenId
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryTokenIdByNftaddr()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryTokenIdByNftaddr(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["nft_addr"] != "" {
		match := regString.MatchString(data["nft_addr"])
		if !match {
			log.Println("verifyInputData_QueryTokenIdByNftaddr() error nft_addr=", data["nft_addr"])
			return ERRINPUTINVALID
		}
	}
	return nil
}
