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

//获取用户nft列表
func (nft *NftExchangeControllerV2) QueryUserNftList() {
	fmt.Println("QueryUserNftList()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryUserNftList() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	defer nft.Ctx.Request.Body.Close()
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	} else {

		inputDataErr := nft.verifyInputData_QueryUserNftList(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			userNfts, totalCount, err := nd.QueryUserNFTList(data["user_addr"], data["start_index"],
				data["count"])
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
				httpResponseData.Data = userNfts
				httpResponseData.TotalCount = uint64(totalCount)
			}
		}
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryUserNftList()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryUserNftList(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["start_index"] !=  "" {
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

	return nil
}
