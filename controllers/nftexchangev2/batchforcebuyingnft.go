package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

// @Title BatchForceBuyingNft
// @Description To force buy snft works, the transaction is initiated by the exchange: post
// @Param Token header string true "token"
// @Param user_addr body string  true "user addr"
// @Param sig body string true "data sig"
// @Param buy_list body string true "BuyList"
// @Success 200
// @Failure 500
// @router /v2/batchForceBuyingNft [post]
func (nft *NftExchangeControllerV2) BatchForceBuyingNft() {
	fmt.Println("BatchForceBuyingNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("BatchForceBuyingNft() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		token := nft.Ctx.Request.Header.Get("Token")
		log.Println("token is : ", token)
		log.Println("data is : ", data)
		inputDataErr := nft.verifyInputData_BatchForceBuyingNft(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
			_, inputDatarr := signature.IsValidAddr(rawData, data["sig"], approveAddr)
			if inputDatarr != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = inputDatarr.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				inputDatarr = nd.InsertSigData(data["sig"], rawData)
				if inputDatarr != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = inputDatarr.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					inputDatarr = nd.BatchForceBuyingNft(data["user_addr"], data["buy_list"])
					if inputDatarr == nil {
						httpResponseData.Code = "200"
						httpResponseData.Data = []interface{}{}
					} else {
						httpResponseData.Code = "500"
						httpResponseData.Msg = inputDatarr.Error()
						httpResponseData.Data = []interface{}{}
					}
				}
			}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
	fmt.Println("BatchForceBuyingNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_BatchForceBuyingNft(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	//regNumber, _ := regexp.Compile(PattenNumber)
	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["sig"] != "" {
		match := regString.MatchString(data["sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//TODO token
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	log.Printf("gentoken is : %v,and token is : %v", getToken, token)
	if getToken != token {
		return ERRTOKEN
	}
	return nil
}
