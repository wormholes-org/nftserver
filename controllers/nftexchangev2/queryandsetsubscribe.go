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

func (nft *NftExchangeControllerV2) verifyInputData_QuerySubscribes(data map[string]string) error {
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

func (nft *NftExchangeControllerV2) QuerySubscribeEmails() {
	fmt.Println("QuerySubscribes()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		inputDataErr := nft.verifyInputData_QuerySubscribes(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			count, annouce, err := nd.QuerySubscribeEmails(data["start_index"], data["count"])
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
	fmt.Println("QuerySubscribes()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_SetSubscribeEmail(data map[string]string) error {
	regEmail, _ := regexp.Compile(PattenEmail)
	regString, _ := regexp.Compile(PattenString)
	if data["user_mail"] != "" {
		match := regEmail.MatchString(data["user_mail"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftExchangeControllerV2) SetSubscribeEmail() {
	fmt.Println("SetSubscribeEmail()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("SetSubscribeEmail() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_SetSubscribeEmail(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			token := nft.Ctx.Request.Header.Get("Token")
			inputDataErr := nft.verifyInputData_BuyingNft(data, token)
			if inputDataErr != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = inputDataErr.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				rawData := signature.RemoveSignData(string(bytes))
				approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
				_, inputDatarr := signature.IsValidAddr(rawData, data["sig"], approveAddr)
				inputDatarr = nd.InsertSigData(data["sig"], rawData)
				if inputDatarr != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = inputDatarr.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					inputDatarr := nd.SetSubscribeEmail(data["user_addr"], data["user_mail"])
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
	fmt.Println("SetSubscribeEmail()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_DelSubscribeEmail(data map[string]string) error {
	regEmail, _ := regexp.Compile(PattenEmail)
	regString, _ := regexp.Compile(PattenString)
	if data["user_mail"] != "" {
		match := regEmail.MatchString(data["user_mail"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftExchangeControllerV2) DelSubscribeEmail() {
	fmt.Println("DelSubscribeEmail()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("DelSubscribeEmail() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_DelSubscribeEmail(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			token := nft.Ctx.Request.Header.Get("Token")
			inputDataErr := nft.verifyInputData_BuyingNft(data, token)
			if inputDataErr != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = inputDataErr.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
				_, inputDatarr := signature.IsValidAddr(rawData, data["sig"], approveAddr)
				inputDatarr = nd.InsertSigData(data["sig"], rawData)
				if inputDatarr != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = inputDatarr.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					inputDatarr := nd.DelSubscribeEmail(data["user_addr"], data["user_mail"])
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
	fmt.Println("DelSubscribeEmail()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
