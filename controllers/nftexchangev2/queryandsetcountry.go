package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

func (nft *NftExchangeControllerV2) GetCountrys() {
	fmt.Println("GetCountrys()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("GetCountrys() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	defer nft.Ctx.Request.Body.Close()
	countrys, err := nd.QueryCountrys()
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
		httpResponseData.Data = countrys
		httpResponseData.TotalCount = uint64(len(countrys))
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetCountrys()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_SetCountrys(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["adminaddr"] != "" {
		match := regString.MatchString(data["adminaddr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["AdminAuth"] != "" {
		match := regNumber.MatchString(data["AdminAuth"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["Admintype"] != "" {
		s := strings.ToLower(data["Admintype"])
		if s != "nfc" && s != "kyc" && s != "admin" {
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftExchangeControllerV2) verifyInputData_SetCountry(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)

	//if data["regionen"] != "" {
	//	match := regString.MatchString(data["regionen"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["regioncn"] != "" {
	//	match := regString.MatchString(data["regioncn"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	if data["domain"] != "" {
		match := regString.MatchString(data["domain"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["telecode"] != "" {
		match := regString.MatchString(data["telecode"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftExchangeControllerV2) SetCountrys() {
	fmt.Println("SetCountry()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetCountry() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_SetCountry(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			_, inputDatarr := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, rawData, data["sig"])
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
					inputDatarr := nd.ModifyCountry(data["regionen"], data["regioncn"], data["domain"], data["telecode"])
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
	fmt.Println("SetCountry()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
