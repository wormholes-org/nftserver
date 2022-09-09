package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"time"
)

func (nft *NftExchangeControllerV2) SetAgreement() {
	fmt.Println("SetAgreement()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetAgreement() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	data, sigData := nft.GetAgreementData()
	inputDataErr := nft.verifyInputData_SetAdmins(data)
	if inputDataErr != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = inputDataErr.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		_, inputDatarr := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, sigData, data["sig"])
		if inputDatarr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDatarr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			f, h, err := nft.GetFile("myfile")
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				f.Close()
				err := nd.SetAgreement(data["param"], f, h.Filename)
				if err == nil {
					httpResponseData.Code = "200"
					httpResponseData.Data = []interface{}{}
				} else {
					httpResponseData.Code = "500"
					httpResponseData.Msg = err.Error()
					httpResponseData.Data = []interface{}{}
				}
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetAgreement()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) GetAgreementData() (map[string]string, string) {
	var data = make(map[string]string)
	var sigData string
	data["param"] = nft.GetString("param", "")
	data["sig"] = nft.GetString("sig", "")
	sigData = data["param"]
	return data, sigData
}
