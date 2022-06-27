package nftexchangev2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

//Upload nft works:post
func (nft *NftExchangeControllerV2) UploadNft() {
	fmt.Println("UploadNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	spendT := time.Now()
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("UploadNft() connect database err = %s\n", err)
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
		inputDataErr := nft.verifyInputData_UploadNft(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			//cid, err := nft.AddFileToIpfs(data["asset_sample"])
			//if err != nil {
			//	httpResponseData.Code = "400"
			//	httpResponseData.Msg = err.Error()
			//	httpResponseData.Data = []interface{}{}
			//} else {
			//	fmt.Printf(">>>>>>>cid=%s\n", cid)
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
					totalt := time.Now()
					err = nd.UploadNft(data["user_addr"], data["creator_addr"], data["owner_addr"],
						data["md5"], data["name"], data["desc"],
						data["meta"], data["source_url"],
						data["nft_contract_addr"], data["nft_token_id"],
						data["categories"], data["collections"],
						data["asset_sample"], data["hide"], data["royalty"], data["count"], data["sig"])
					if err == nil {
						httpResponseData.Code = "200"
						httpResponseData.Data = []interface{}{}
					} else {
						httpResponseData.Code = "500"
						httpResponseData.Msg = err.Error()
						httpResponseData.Data = []interface{}{}
					}
					fmt.Printf(" nd.UploadNft() total Spend time=%s time.now=%s\n", time.Now().Sub(totalt), time.Now())
				}
			}
			//}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "Incorrect user information entered"
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()

	fmt.Printf("UploadNft() Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	fmt.Println("UploadNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) AddFileToIpfs(content string) (string, error) {
	serverIP, _ := beego.AppConfig.String("nftIpfsServerIP")
	serverPort, _ := beego.AppConfig.String("nftIpfsServerPort")
	url := "http://" + serverIP + ":" + serverPort + "/v1/ipfsadd"
	fmt.Println("NftExchangeController.AddFileToIpfs(), url=", url)
	var mapImage map[string]string
	mapImage = make(map[string]string, 0)
	mapImage["asset"] = content
	respData, err := nft.SendPost(url, mapImage, "application/json")
	if err != nil {
		return "", err
	} else {
		cid, ok := respData.Data.(string)
		if ok {
			return cid, nil
		} else {
			return "", errors.New("nftipfs server 返回数据格式错误！")
		}
	}
}

func (nft *NftExchangeControllerV2) SendPost(url string, data interface{}, contentType string) (respData controllers.HttpResponseData, err error) {
	jsonStr, _ := json.Marshal(data)
	//fmt.Println("SendPost,url=", url, "jsonStr=", string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonStr))
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Add("content-type", contentType)
	if err != nil {
		return
	}
	defer req.Body.Close()
	//client := &http.Client{Timeout: 5 * time.Second}
	client := &http.Client{}
	fmt.Println("SendPost url=", url)
	resp, error := client.Do(req)
	if error != nil {
		return
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(result, &respData)
	return
}

func (nft *NftExchangeControllerV2) IpfsTest() {
	content := "test content"
	s, e := nft.AddFileToIpfs(content)
	if e != nil {
		nft.Ctx.WriteString(e.Error())
	} else {
		nft.Ctx.WriteString(s)
	}
}

func (nft *NftExchangeControllerV2) DelNft() {
	fmt.Println("DelNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		rawData := signature.RemoveSignData(string(bytes))
		token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_DelNft(data, token)
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
				inputDatarr := nd.DelNft(data["user_addr"], data["contract"], data["tokenid"])
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
		httpResponseData.Msg = "Incorrect user information entered"
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("DelNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) SetNft() {
	fmt.Println("SetNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	spendT := time.Now()
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("UploadNftImage() connect database err = %s\n", err)
		return
	}
	defer nd.Close()

	data, sigData := nft.GetData()
	token := nft.Ctx.Request.Header.Get("Token")
	inputDataErr := nft.verifyInputData_UploadNftImage(data, token)
	if inputDataErr != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = inputDataErr.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
		_, err := signature.IsValidAddr(sigData, data["sig"], approveAddr)
		if err != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = err.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			err = nd.InsertSigData(data["sig"], sigData)
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				f, h, err := nft.GetFile("myfile")
				imageHash := ""
				filename := ""
				if err == nil {
					f.Close()
					imageHash, err = nft.SaveToIpfs("myfile")
					filename = h.Filename
				} else {
					err = nil
				}
				if err != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = err.Error()
					httpResponseData.Data = []interface{}{}
					fmt.Println("SaveToFile err=", err)
				} else {
					totalt := time.Now()
					nftimage, err := nd.SetNft(data["user_addr"],
						data["md5"], data["name"], data["desc"],
						imageHash, filename,
						data["nft_token_id"],
						data["categories"], data["collections"],
						data["asset_sample"], data["hide"], data["royalty"], data["count"], data["attributes"], data["sig"])
					if err == nil {
						httpResponseData.Code = "200"
						httpResponseData.Data = nftimage
					} else {
						httpResponseData.Code = "500"
						httpResponseData.Msg = err.Error()
						httpResponseData.Data = []interface{}{}
					}
					fmt.Printf(" nd.SetNft() total Spend time=%s time.now=%s\n", time.Now().Sub(totalt), time.Now())
				}
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)

	fmt.Printf("SetNft() Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	fmt.Println("SetNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_UploadNft(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)
	regImage, _ := regexp.Compile(PattenImageBase64)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["creator_addr"] != "" {
		match := regString.MatchString(data["creator_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["owner_addr"] != "" {
		match := regString.MatchString(data["owner_addr"])
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
	if data["md5"] != "" {
		match := regString.MatchString(data["md5"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["categories"] != "" {
		match := regString.MatchString(data["categories"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["collections"] != "" {
	//	match := regString.MatchString(data["collections"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	match := regImage.MatchString(data["asset_sample"])
	if !match {
		return ERRINPUTINVALID
	}
	if data["hide"] != "" {
		match := regString.MatchString(data["hide"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["royalty"] != "" {
		match := regNumber.MatchString(data["royalty"])
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

func (nft *NftExchangeControllerV2) verifyInputData_DelNft(data map[string]string, token string) error {
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
