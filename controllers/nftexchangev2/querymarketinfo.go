package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"time"
)

//获取市场数据
func (nft *NftExchangeControllerV2) QueryMarketInfo() {
	fmt.Println("QueryMarketInfo()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryMarketInfo() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	marketInfo, err := nd.QueryMarketInfo()
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
		httpResponseData.Data = marketInfo
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryMarketInfo()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//获取nft总量数据
func (nft *NftExchangeControllerV2) GetNftMarketInfo() {
	fmt.Println("QueryMarketInfo()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryMarketInfo() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	marketInfo, err := nd.GetNftMarketInfo()
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
		httpResponseData.Data = marketInfo
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryMarketInfo()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
