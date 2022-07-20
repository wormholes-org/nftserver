package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

//Create period:post
func (nft *NftExchangeControllerV2) CreatsSnftphase() {
	fmt.Println("CreatsSnftphase()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_NewPeriod(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			_, err := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, rawData, data["sig"])
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
					err = nd.NewSnftPhase(data["user_addr"], data["name"],
						data["desc"])
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
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("CreatsSnftphase()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//Modify period
func (nft *NftExchangeControllerV2) SetSnftPeriod() {
	fmt.Println("SetSnftPeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	period := models.Period{}
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_NewPeriod(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			_, err := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, rawData, data["sig"])

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
					err = json.Unmarshal(bytes, &period)
					if err != nil {
						nft.Ctx.ResponseWriter.Write([]byte("Failed to update data！"))
					} else {
						err = nd.SetSnftPeriod(period)
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
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetSnftPeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//获取period
func (nft *NftExchangeControllerV2) GetSnftPeriod() {
	fmt.Println("GetSnftPeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	homePageData, err := nd.GetSnftPeriod()
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
		httpResponseData.TotalCount = 0
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetSnftPeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//Get all voting periods
func (nft *NftExchangeControllerV2) GetAllVotePeriod() {
	fmt.Println("GetAllVotePeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	homePageData, err := nd.GetAllVotePeriod()
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
		httpResponseData.TotalCount = 0
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetAllVotePeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//Get three voting periods
func (nft *NftExchangeControllerV2) GetAccedVotePeriod() {
	fmt.Println("GetAccedVotePeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	homePageData, err := nd.GetAccedVotePeriod()
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
		httpResponseData.TotalCount = 0
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetAccedVotePeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

//delete period
func (nft *NftExchangeControllerV2) DelSnftPeriod() {
	fmt.Println("DelSnftPeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("DelAnnounces() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
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
				inputDatarr := nd.DelSnftPeriod(data["id"])
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
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("DelSnftPeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) GetVotePeriod() {
	fmt.Println("GetVotePeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
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
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	} else {

		inputDataErr := nft.verifyInputData_QueryNftCollectionList(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			nftCollections, totalCount, err := nd.GetVoteSnftPeriod(data["start_index"], data["count"])
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
				httpResponseData.Data = nftCollections
				httpResponseData.TotalCount = uint64(totalCount)
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetVotePeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) SetVotePeriod() {
	fmt.Println("SetVotePeriod()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	_ = http.Client{Timeout: time.Hour}

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_NewPeriod(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			_, err := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, rawData, data["sig"])

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
					err = nd.SetVoteSnftPeriod(data["param"])
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
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetVotePeriod()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) SetPeriodEth() {
	fmt.Println("SetPeriodEth()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("CreateCollection() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_Tokenid(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			_, err := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminEdit, rawData, data["sig"])

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
					err = nd.SetPeriodEth(data["param"])
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
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetPeriodEth()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func sayHelloTimeout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	fmt.Println(w, "hello world")
	time.Sleep(5 * time.Second)

}

func (nft *NftExchangeControllerV2) verifyInputData_NewPeriod(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["name"] != "" {
	//	match := regString.MatchString(data["name"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["desc"] != "" {
	//	match := regString.MatchString(data["desc"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	return nil
}

func (nft *NftExchangeControllerV2) verifyInputData_Tokenid(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["param"] != "" {
		match := regString.MatchString(data["param"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["name"] != "" {
	//	match := regString.MatchString(data["name"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["desc"] != "" {
	//	match := regString.MatchString(data["desc"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	return nil
}
