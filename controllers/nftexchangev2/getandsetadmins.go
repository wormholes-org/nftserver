package nftexchangev2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*func (nft *NftExchangeControllerV2) IsValidSuperAddr(rawData string, sig string) (bool, error) {
	var addrList []string
	addrList = append(addrList, models.SuperAdminAddr)
	verificationAddr, err := signature.GetEthAddr(rawData, sig)
	if err != nil {
		return false, err
	}
	verificationAddrS := verificationAddr.String()
	verificationAddrS = strings.ToLower(verificationAddrS)

	if verificationAddrS == models.SuperAdminAddr {
		fmt.Println("sigdebug verify [Y]")
		return true, nil
	}
	fmt.Println("sigdebug verify [N]")

	return false, errors.New("verification address is invalid")
}*/

func (nft *NftExchangeControllerV2) IsValidVerifyAddr(nd *models.NftDb, adminType models.AdminType, adminAuth models.AdminAuthType, rawData, sig string) (bool, error) {
	verificationAddr, err := signature.GetEthAddr(rawData, sig)
	if err != nil {
		log.Println("IsValidVerifyAddr() GetEthAddr() err=", err)
		return false, err
	}
	verificationAddrS := verificationAddr.String()
	verificationAddrS = strings.ToLower(verificationAddrS)
	fmt.Println("IsValidVerifyAddr() verificationAddrS=", verificationAddrS)
	admin := models.Admins{}
	db := nd.GetDB().Model(&models.Admins{}).Where("adminaddr = ? AND admintype = ?", verificationAddrS, adminType.String()).First(&admin)
	if db.Error != nil {
		fmt.Println("IsValidVerifyAddr() err=", err)
		return false, errors.New("address error ")
	}
	fmt.Println("IsValidVerifyAddr() verificationAddrS=", verificationAddrS, "admin.AdminAuth=", admin.AdminAuth)

	switch adminAuth {
	case models.AdminBrowse:
		break
	case models.AdminEdit:
		auth, _ := strconv.Atoi(admin.AdminAuth)
		if models.AdminAuthType(auth) == models.AdminEdit || models.AdminAuthType(auth) == models.AdminBrowseEditAudit {
			return true, nil
		} else {
			return false, errors.New("address not permission.")
		}
	case models.AdminAudit:
		auth, _ := strconv.Atoi(admin.AdminAuth)
		if models.AdminAuthType(auth) == models.AdminAudit || models.AdminAuthType(auth) == models.AdminBrowseEditAudit {
			return true, nil
		} else {
			return false, errors.New("address not permission.")
		}
	case models.AdminBrowseEditAudit:
		auth, _ := strconv.Atoi(admin.AdminAuth)
		if models.AdminAuthType(auth) == models.AdminBrowseEditAudit {
			return true, nil
		} else {
			return false, errors.New("address not permission.")
		}
	}
	return false, errors.New("verification address is invalid")
}

func (nft *NftExchangeControllerV2) GetAdmins() {
	fmt.Println("GetAdmins()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		inputDataErr := nft.verifyInputData_GetAdmins(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			count, admins, err := nd.QueryAdmins(data["admintype"], data["start_index"], data["count"])
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
				httpResponseData.Data = admins
				httpResponseData.TotalCount = uint64(count)
			}
		}
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetAdmins()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_GetAdmins(data map[string]string) error {
	regNumber, _ := regexp.Compile(PattenNumber)
	if data["Admintype"] != "" {
		s := strings.ToLower(data["Admintype"])
		if s != "nfc" && s != "kyc" && s != "admin" {
			return ERRINPUTINVALID
		}
	}
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

func (nft *NftExchangeControllerV2) verifyInputData_SetAdmins(data map[string]string) error {
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

func (nft *NftExchangeControllerV2) SetAdmins() {
	fmt.Println("SetSysParams()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		inputDataErr := nft.verifyInputData_SetAdmins(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			rawData := signature.RemoveSignData(string(bytes))
			_, inputDatarr := nft.IsValidVerifyAddr(nd, models.AdminTypeAdmin, models.AdminBrowseEditAudit, rawData, data["sig"])
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
					inputDatarr := nd.ModifyAdmin(data["adminaddr"], data["admintype"], data["adminauth"])
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
		httpResponseData.Msg = "输入的用户信息错误"
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("SetAdmins()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_DelAdmins(data map[string]string) error {
	//regString, _ := regexp.Compile(PattenString)

	//if data["del_admins"] != "" {
	//	match := regString.MatchString(data["del_admins"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	return nil
}

func (nft *NftExchangeControllerV2) DelAdmins() {
	fmt.Println("DelAdmins()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd, err := models.NewNftDb(models.Sqldsndb)
	if err != nil {
		fmt.Printf("DelAdmins() connect database err = %s\n", err)
		return
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_DelAdmins(data)
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
					inputDatarr := nd.DelAdmins(data["del_admins"])
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
		httpResponseData.Msg = "输入的用户信息错误"
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("DelAdmins()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
