package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
)

//Query the homepage data
func (nft *NftExchangeControllerV2) GetOverview() {
	var spendT = time.Now()
	fmt.Println("GetOverview()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", spendT)
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("GetOverview() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	OverviewData, err := nd.GetOverviewExcel()
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
		httpResponseData.Data = OverviewData
		httpResponseData.TotalCount = 1
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetOverview() <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) GetSnftPeriodNum() {
	fmt.Println("GetSnftPeriodNum()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("ModifyCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		resCount, err := nd.GetSnftPeriodNum(data["addr"])
		if err == nil {
			httpResponseData.Code = "200"
			httpResponseData.Data = resCount
		} else {
			httpResponseData.Code = "500"
			httpResponseData.Msg = err.Error()
			httpResponseData.Data = []interface{}{}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetSnftPeriodNum()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) GetSnftPledge() {
	fmt.Println("GetSnftPledge()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("ModifyCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	resCount, pledge, err := nd.GetSnftPeledge()
	if err == nil {
		httpResponseData.Code = "200"
		httpResponseData.Data = resCount
		httpResponseData.TotalCount = uint64(pledge)
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetSnftPledge()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
