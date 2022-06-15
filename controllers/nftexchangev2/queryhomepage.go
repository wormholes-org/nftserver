package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"time"
)

//查询首页数据
func (nft *NftExchangeControllerV2) QueryHomePage() {
	var spendT = time.Now()
	fmt.Println("QueryHomePage()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", spendT)
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryHomePage() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	homePageData, err := nd.QueryHomePage(false)
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
		httpResponseData.Data = homePageData
		httpResponseData.TotalCount = 1
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryHomePage() Spend time",time.Now().Sub(spendT), "time.now=", time.Now())
	fmt.Println("QueryHomePage() <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<",time.Now())
}