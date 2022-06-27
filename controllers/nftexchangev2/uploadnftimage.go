package nftexchangev2

import (
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"log"
	"regexp"
	"time"
)

//const UpLoadSize = 100000000

func (nft *NftExchangeControllerV2) GetData() (map[string]string, string) {
	var data = make(map[string]string)
	var sigData string
	data["user_addr"] = nft.GetString("user_addr", "")
	data["creator_addr"] = nft.GetString("creator_addr", "")
	data["owner_addr"] = nft.GetString("owner_addr", "")
	data["md5"] = nft.GetString("md5", "")
	data["name"] = nft.GetString("name", "")
	data["desc"] = nft.GetString("desc", "")
	data["source_url"] = nft.GetString("source_url", "")
	data["nft_contract_addr"] = nft.GetString("nft_contract_addr", "")
	data["nft_token_id"] = nft.GetString("nft_token_id", "")
	data["categories"] = nft.GetString("categories", "")
	data["collections"] = nft.GetString("collections", "")
	data["asset_sample"] = nft.GetString("asset_sample", "")
	data["hide"] = nft.GetString("hide", "")
	data["royalty"] = nft.GetString("royalty", "")
	data["count"] = nft.GetString("count", "")
	data["attributes"] = nft.GetString("attributes", "")
	data["sig"] = nft.GetString("sig", "")
	sigData = data["user_addr"] /*+ data["creator_addr"] + data["owner_addr"]*/ + data["md5"] + data["name"] +
		data["desc"] + /*data["meta"] +*/ /*data["source_url"] +*/ /*data["nft_contract_addr"] + data["nft_token_id"] +*/
		data["categories"] + data["collections"] + /*data["asset_sample"] +*/ data["hide"] + data["royalty"] + data["count"]
	return data, sigData
}

func (nft *NftExchangeControllerV2) SaveToIpfs(fromfile string) (string, error) {
	spendT := time.Now()
	file, fileh, err := nft.Ctx.Request.FormFile(fromfile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if fileh.Size > int64(models.UploadSize) {
		log.Println("SaveToIpfs() upload size too big!")
		return "", errors.New("upload size too big!")
	}
	url := models.NftIpfsServerIP + ":" + models.NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(500 * time.Second)
	mhash, err := s.Add(file)
	if err != nil {
		log.Println("SaveToIpfs() err=", err)
		return "", err
	}
	fmt.Printf("SaveToIpfs  Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return mhash, nil
}

//Upload nft works:post
func (nft *NftExchangeControllerV2) UploadNftImage() {
	fmt.Println("UploadNftImage()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
				if err != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = err.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					f.Close()
					imageHash, err := nft.SaveToIpfs("myfile")
					if err != nil {
						httpResponseData.Code = "500"
						httpResponseData.Msg = err.Error()
						httpResponseData.Data = []interface{}{}
					} else {
						totalt := time.Now()

						if err != nil {
							httpResponseData.Code = "500"
							httpResponseData.Msg = err.Error()
							httpResponseData.Data = []interface{}{}
							fmt.Println("SaveToFile err=", err)
						} else {
							nftimage, err := nd.UploadNftImage(data["user_addr"], data["creator_addr"], data["owner_addr"],
								data["md5"], data["name"], data["desc"],
								imageHash, h.Filename,
								data["nft_contract_addr"], data["nft_token_id"],
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
						}
						fmt.Printf(" nd.UploadNftImage() total Spend time=%s time.now=%s\n", time.Now().Sub(totalt), time.Now())
					}
				}
			}
		}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()

	fmt.Printf("UploadNftImage() Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	fmt.Println("UploadNftImage()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) AddImageFileToIpfs(content string) (string, error) {
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

func (nft *NftExchangeControllerV2) verifyInputData_UploadNftImage(data map[string]string, token string) error {
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
