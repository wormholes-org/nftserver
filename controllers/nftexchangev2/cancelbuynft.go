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

//Cancel the sale (off the shelf)
func (nft *NftExchangeControllerV2) CancelBuyNft() {
	fmt.Println("CancelBuyNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CancelBuyNft() connect database err = %s\n", err)
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
		token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_CancelBuyNft(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			rawData := signature.RemoveSignData(string(bytes))
			approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
			_, err := signature.IsValidAddr(rawData, data["sig"], approveAddr)
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				err = nd.InsertSigData(data["sig"], rawData)
				if err != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = err.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					err := nd.CancelBuy(data["user_addr"], data["nft_contract_addr"],
						data["nft_token_id"], data["trade_sig"], data["sig"])
					if err != nil {
						httpResponseData.Code = "500"
						httpResponseData.Msg = err.Error()
						httpResponseData.Data = []interface{}{}
					} else {
						httpResponseData.Code = "200"
						httpResponseData.Data = []interface{}{}
					}
				}
			}
		}
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("CancelBuyNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_CancelBuyNft(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nft_contract_addr"] != "" {
		match := regString.MatchString(data["nft_contract_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nft_token_id"] != "" {
		match := regString.MatchString(data["nft_token_id"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	/*if data["trade_sig"] != "" {
		match := regString.MatchString(data["trade_sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}*/
	if data["sig"] != "" {
		match := regString.MatchString(data["sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	if getToken != token {
		return ERRTOKEN
	}

	return nil
}
