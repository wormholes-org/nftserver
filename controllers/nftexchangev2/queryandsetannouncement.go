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
	"time"
)

func (nft *NftExchangeControllerV2) verifyInputData_QueryAnnounce(data map[string]string) error {
	regNumber, _ := regexp.Compile(PattenNumber)
	if data["start_index"] != "" {
		match := regNumber.MatchString(data["start_index"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["count"] != "" {
		match := regNumber.MatchString(data["count"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftExchangeControllerV2) QueryAnnounce() {
	fmt.Println("QueryAnnounce()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("GetAdmins() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	defer nft.Ctx.Request.Body.Close()
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		inputDataErr := nft.verifyInputData_QueryAnnounce(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			count, annouce, err := nd.QueryAnnouncements(data["start_index"], data["count"])
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
				httpResponseData.Data = annouce
				httpResponseData.TotalCount = uint64(count)
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryAnnounce()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_SetAnnounce(data map[string]string) error {
	//regString, _ := regexp.Compile(PattenString)
	//regNumber, _ := regexp.Compile(PattenNumber)
	//
	//if data["title"] != "" {
	//	match := regString.MatchString(data["title"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["content"] != "" {
	//	match := regNumber.MatchString(data["content"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	return nil
}

func (nft *NftExchangeControllerV2) SetAnnounce() {
	fmt.Println("SetAnnounce()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetAdmins() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_SetAnnounce(data)
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
					inputDatarr := nd.SetAnnouncement(data["title"], data["content"])
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
	fmt.Println("SetAnnounce()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_DelAnnounce(data map[string]string) error {
	//regString, _ := regexp.Compile(PattenString)
	//regNumber, _ := regexp.Compile(PattenNumber)
	//
	//if data["title"] != "" {
	//	match := regString.MatchString(data["title"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["content"] != "" {
	//	match := regNumber.MatchString(data["content"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	return nil
}

func (nft *NftExchangeControllerV2) DelAnnounces() {
	fmt.Println("SetAnnounce()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("DelAnnounces() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_DelAnnounce(data)
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
					inputDatarr := nd.DelAnnounce(data["del_announces"])
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
	fmt.Println("DelAnnounces()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
