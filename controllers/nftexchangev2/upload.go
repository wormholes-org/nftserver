package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"regexp"
	"strings"
	"time"
)

func (nft *NftExchangeControllerV2) verifyInputData_Upload(data map[string]string) error {
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

func (nft *NftExchangeControllerV2) verifySign_Upload(user, contract, tokenid, sig string) error {
	return nil
	msg := user + contract + tokenid
	msghash := crypto.Keccak256([]byte(msg))
	hexsig, err := hexutil.Decode("0x" + sig)
	if err != nil {
		fmt.Println("verifySign_Upload() Decode() err=", err)
		return err
	}
	pub, err := crypto.SigToPub(msghash, hexsig)
	if err != nil {
		fmt.Println("verifySign_Upload() TransactionSender() err=", err)
		return err
	}
	toaddr := crypto.PubkeyToAddress(*pub)
	if toaddr.String() != user {
		fmt.Println("verifySign_Upload() PubkeyToAddress() buyer address error.")
		return err
	}
	return nil
}

func (nft *NftExchangeControllerV2) UpLoadFile() {
	fmt.Println("UpLoad()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	//nd, err := models.NewNftDb(models.Sqldsndb)
	//if err != nil {
	//	fmt.Printf("UpLoad() connect database err = %s\n", err)
	//	return
	//}
	//defer nd.Close()

	contract := nft.GetString("contract", "")
	tokenid := nft.GetString("tokenid", "")
	user := nft.GetString("user", "")
	sig := nft.GetString("sig", "")
	keyname := nft.GetString("keyname", "")
	fmt.Println("keyname=", keyname)
	if err := nft.verifySign_Upload(user, contract, tokenid, sig); err == nil {
		f, h, err := nft.GetFile("myfile") //获取上传的文件
		if err != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = err.Error()
			httpResponseData.Data = []interface{}{}
			fmt.Println("SaveToFile err=", err)
		} else {
			path := models.ImageDir + h.Filename
			f.Close()
			err = nft.SaveToFile("myfile", path)
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
				fmt.Println("SaveToFile err=", err)
			} else {
				httpResponseData.Code = "200"
				httpResponseData.Data = []interface{}{}
			}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "Incorrect user signature entered"
		httpResponseData.Data = []interface{}{}
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("UpLoad()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
