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

func (nft *NftExchangeControllerV2) GetSysParams() {
	sysParams := &models.SysParamsInfo{}
	fmt.Println("GetSysParams()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("GetSysParams() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	defer nft.Ctx.Request.Body.Close()
	sysParams, err = nd.QuerySysParams()
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
		httpResponseData.Data = sysParams
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetSysParams()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//Return system parameters according to params
func (nft *NftExchangeControllerV2) GetSysParamByParams() {
	fmt.Println("GetSysParamByParams()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	//defer nft.Ctx.Request.Body.Close()
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("GetSysParamByParams() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]interface{}
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	json.Unmarshal(bytes, &data)
	s, ok := data["params"].(string)
	fmt.Printf(">>>>>>>>s=%s\n", s)
	if ok {
		userInfo, err := nd.GetSysParam(s)
		if err == nil {
			httpResponseData.Code = "200"
			httpResponseData.Data = userInfo
		} else {
			if err == gorm.ErrRecordNotFound || err == models.ErrNftNotExist {
				httpResponseData.Code = "200"
			} else {
				httpResponseData.Code = "500"
			}
			httpResponseData.Msg = err.Error()
			httpResponseData.Data = []interface{}{}
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
	fmt.Println("GetSysParamByParams()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_SetSysParams(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["nft1155addr"] != "" {
		match := regString.MatchString(data["nft1155addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["adminaddr"] != "" {
		match := regString.MatchString(data["adminaddr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["lowprice"] != "" {
		match := regNumber.MatchString(data["lowprice"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["blocknumber"] != "" {
		match := regNumber.MatchString(data["blocknumber"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["scannumber"] != "" {
		match := regNumber.MatchString(data["scannumber"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["royaltylimit"] != "" {
		match := regNumber.MatchString(data["royaltylimit"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nftaudit"] != "" {
		s := strings.ToLower(data["nftaudit"])
		if s != "true" && s != "false" {
			return ERRINPUTINVALID
		}
	}
	if data["userkyc"] != "" {
		s := strings.ToLower(data["userkyc"])
		if s != "true" && s != "false" {
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftExchangeControllerV2) SetSysParams() {
	fmt.Println("SetSysParams()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	sysParams := models.SysParamsInfo{}
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetSysParams() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_SetSysParams(data)
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
					err = json.Unmarshal(bytes, &sysParams)
					if err != nil {
						nft.Ctx.ResponseWriter.Write([]byte("Failed to update data！"))
					} else {
						inputDatarr := nd.SetSysParams(sysParams)
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
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetSysParams()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) SetAnnouncementParams() {
	fmt.Println("SetAnnouncementParams()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetAnnouncementParams() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_SetAnnouncementParams(data, token)
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
					inputDatarr := nd.SetAnnouncementParam(data["params"])
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
	fmt.Println("SetAnnouncementParams()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) SetExchangeSig() {
	fmt.Println("SetExchangeSig()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetExchangeSig() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_SetExchangeSig(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			//approveAddr, _ := approveAddrsMap.GetApproveAddr(data["exchanger_owner"] + data["to"] + data["block_number"])
			//_, err := signature.IsValidAddr(rawData, data["sig"], approveAddr)
			//rawData := signature.RemoveSignData(string(bytes))
			_, err := signature.IsValidAddr(data["exchanger_owner"]+data["to"]+data["block_number"], data["sig"], data["exchanger_owner"])
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
					exchange, err := json.Marshal(data)
					if err != nil {
						nft.Ctx.ResponseWriter.Write([]byte("Failed to update data！"))
					} else {
						inputDatarr := nd.SetExchageSig(string(exchange))
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
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetExchangeSig()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) GetExchangeSig() {
	fmt.Println("GetExchangeSig()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("GetExchangeSig() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	input, err := nd.GetExchageSig()
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
		httpResponseData.Data = input
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetExchangeSig()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_SetAnnouncementParams(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)

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
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	if getToken != token {
		return ERRTOKEN
	}
	return nil
}

func (nft *NftExchangeControllerV2) verifyInputData_SetExchangeSig(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["exchanger_owner"] != "" {
		match := regString.MatchString(data["exchanger_owner"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["to"] != "" {
		match := regString.MatchString(data["to"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["block_number"] != "" {
		match := regString.MatchString(data["block_number"])
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

	return nil
}
