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

// @Title QueryRecommendSnfts
// @Description Query Homepage  recommend snft:post
// @Param user_addr body string true "user addr"  example()
// @Success 200 {object} models.RecommendBuyingSell.Buying
// @Failure 500
// @router /v2/queryRecommendSnfts [post]
func (nft *NftExchangeControllerV2) QueryRecommendSnfts() {
	fmt.Println("QueryRecommendSnfts()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	//defer nft.Ctx.Request.Body.Close()
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("QueryRecommendSnfts() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]interface{}
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	json.Unmarshal(bytes, &data)
	s, ok := data["user_addr"].(string)
	fmt.Printf("QueryRecommendSnfts() user_addr =%s\n", s)
	if ok {
		inputDataErr := nft.verifyInputData_QueryRecommendSnfts(s)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			snfts, err := nd.QueryRecommendSnfts(s)
			if err == nil {
				httpResponseData.Code = "200"
				httpResponseData.Data = snfts
			} else {
				if err == gorm.ErrRecordNotFound || err == models.ErrNftNotExist {
					httpResponseData.Code = "200"
				} else {
					httpResponseData.Code = "500"
				}
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Data = []interface{}{}
		httpResponseData.Msg = ERRINPUT.Error()
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
	fmt.Println("QueryRecommendSnfts()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryRecommendSnfts(user string) error {
	regString, _ := regexp.Compile(PattenString)
	if user != "" {
		match := regString.MatchString(user)
		if !match {
			return ERRINPUTINVALID
		}
	}
	return nil
}
