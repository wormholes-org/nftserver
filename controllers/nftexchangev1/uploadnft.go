package nftexchangev1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"net/http"
)

//Upload nft works: post
func (nft *NftExchangeControllerV1) UploadNft() {

	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		//Call nftipfs to upload the image to the ipfs file server
		cid, err := nft.AddFileToIpfs(data["asset_sample"])
		if err != nil {
			httpResponseData.Code = "400"
			httpResponseData.Msg = err.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			fmt.Printf(">>>>>>>cid=%s\n", cid)
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
			}
		}

	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "Incorrect user information entered"
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
}

func (nft *NftExchangeControllerV1) AddFileToIpfs(content string) (string, error) {
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
			return "", errors.New("nftipfs server return data format error！")
		}
	}
}

//Send POST request
//url: request address data: data submitted by POST request contentType: request body format, such as: application/json
//content: The content returned by the request
func (nft *NftExchangeControllerV1) SendPost(url string, data interface{}, contentType string) (respData controllers.HttpResponseData, err error) {
	jsonStr, _ := json.Marshal(data)
	fmt.Println("SendPost,url=", url, "jsonStr=", string(jsonStr))
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

func (nft *NftExchangeControllerV1) IpfsTest() {
	content := "test content"
	s, e := nft.AddFileToIpfs(content)
	if e != nil {
		nft.Ctx.WriteString(e.Error())
	} else {
		nft.Ctx.WriteString(s)
	}
}
