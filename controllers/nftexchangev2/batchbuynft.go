package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"regexp"
	"time"
)

//Buy nft works: post
func (nft *NftExchangeControllerV2) BatchBuyNft() {
	fmt.Println("BatchBuyNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("BuyNft() connect database err = %s\n", err)
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
		inputDataErr := nft.verifyInputData_BatchBuyNft(data, token)
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
					inputDatarr = nd.BatchMakeOffer(data["user_addr"], data["offer_list"])
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
	fmt.Println("BatchBuyNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_BatchBuyNft(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	//regNumber, _ := regexp.Compile(PattenNumber)
	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["nft_token_id"] != "" {
	//	match := regString.MatchString(data["nft_token_id"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["vote_stage"] != "" {
	//	match := regString.MatchString(data["vote_stage"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["pay_channel"] != "" {
	//	match := regString.MatchString(data["pay_channel"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["currency_type"] != "" {
	//	match := regString.MatchString(data["currency_type"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["price"] != "" {
	//	match := regNumber.MatchString(data["price"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["trade_sig"] != "" {
	//	match := regString.MatchString(data["trade_sig"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["dead_time"] != "" {
	//	match := regNumber.MatchString(data["dead_time"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	if data["sig"] != "" {
		match := regString.MatchString(data["sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["nft_contract_addr"] != "" {
	//	match := regString.MatchString(data["nft_contract_addr"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	if getToken != token {
		return ERRTOKEN
	}

	return nil
}
