package nftexchangev2

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/big"
	rands "math/rand"
	"regexp"
	"strconv"
	"time"
)

type Captcha struct {
	Imageid   string `json:"imageid"`
	Imagey    int    `json:"imagey"`
	Image     string `json:"image"`
	Imagecrop string `json:"imagecrop"`
}

var CaptchaMap = map[string]string{}

//Get bot captcha
func (nft *NftExchangeControllerV2) GetCaptcha() {
	fmt.Println("GetCaptcha()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData

	defer nft.Ctx.Request.Body.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_Captcha(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			singleNft, err := CreateNetCode()
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				httpResponseData.Code = "200"
				httpResponseData.Data = singleNft
			}
		}

	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetCaptcha()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_Captcha(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
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

func (nft *NftExchangeControllerV2) AuthCaptcha() {
	fmt.Println("AuthCaptcha()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData

	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err := json.Unmarshal(bytes, &data)
	if err == nil {
		singleNft, err := AuthCode(data["tokenid"], data["imagex"])
		if err != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = err.Error()
			httpResponseData.Data = singleNft
		} else {
			httpResponseData.Code = "200"
			httpResponseData.Data = singleNft
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "Incorrect user information entered"
		httpResponseData.Data = []interface{}{}
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("AuthCaptcha()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func GetRandInt(max int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-1)))
	return int(num.Int64())
}

func AuthCode(tokenid, imagex string) (bool, error) {
	imgx, ok := CaptchaMap[tokenid]
	if !ok {
		fmt.Println("captcha nil")
		return false, errors.New("captcha nil")
	}
	imgxint, err := strconv.Atoi(imgx)
	if err != nil {
		fmt.Println("imgxint transfer int err =", err)
		return false, err
	}
	imagexint, err := strconv.Atoi(imagex)
	if err != nil {
		fmt.Println("imagexint transfer int err =", err)
		return false, err
	}
	if imagexint <= imgxint+10 && imagexint >= imgxint-10 {
		delete(CaptchaMap, tokenid)
		fmt.Println("captcha ok")
		return true, nil
	} else {
		fmt.Println("captcha auth error")
		return false, errors.New("captcha auth err")
	}

}

func CreateNetCode() (Captcha, error) {

	captcha := Captcha{}
	nums := GetRandInt(models.DefaultCaptchaNum)
	var imageId string
	rands.Seed(time.Now().UnixNano())
	var i int
	for i = 0; i < 20; i++ {
		s := fmt.Sprintf("%d", rands.Int63())
		if len(s) < 15 {
			continue
		}
		s = s[len(s)-13:]
		imageId = s
		if s[0] == '0' {
			continue
		}
		fmt.Println("imageId() NewTokenid=", imageId)
		_, ok := CaptchaMap[imageId]
		if !ok {
			break
		}
	}
	fi := fmt.Sprintf("%05x", nums)
	f, err := ioutil.ReadFile(models.ImageDir + "/captcha/" + fi)
	if err != nil {
		fmt.Println("CreateNetCode() readfile err=", err)
		return Captcha{}, err
	}
	originalimg, _, err := image.Decode(bytes.NewReader(f))
	if err != nil {
		fmt.Printf("image Decode  failed! %v", err)
		return Captcha{}, err
	}
	imageRandX := GetRandInt(400)
	if imageRandX < 60 {
		imageRandX += 60
	}

	imageRandY := GetRandInt(300)
	if imageRandY < 60 {
		imageRandY += 60
	}

	maxPotion := image.Point{
		X: imageRandX,
		Y: imageRandY,
	}
	minPotion := image.Point{
		X: imageRandX - 60,
		Y: imageRandY - 60,
	}
	subimg := image.Rectangle{
		Max: maxPotion,
		Min: minPotion,
	}

	overdata := imaging.Crop(originalimg, subimg)

	f, err = ioutil.ReadFile(models.ImageDir + "/captcha/" + "mask")
	if err != nil {
		fmt.Println("CreateNetCode() maskimg readfile err=", err)
		return Captcha{}, err
	}
	maskimg, _, err := image.Decode(bytes.NewReader(f))
	if err != nil {
		fmt.Printf("maskimg image Decode  failed! %v", err)
		return Captcha{}, err
	}

	bgdata := imaging.Overlay(originalimg, maskimg, minPotion, 1.0)
	//f, err = os.Create("./captcha/code/" + imageId + ".jpeg")
	//defer f.Close()
	//jpeg.Encode(f, data2, nil)
	fmt.Println("captcha:", imageId, ",", imageRandX, ",", imageRandY)

	img := image.NewNRGBA(image.Rect(0, 0, 60, 60))
	fmt.Println(maskimg.Bounds().Dy())
	trans := maskimg.At(0, 0)
	fmt.Println(maskimg.At(0, 0) == trans)
	for x := 0; x < 60; x++ {
		for y := 0; y < 60; y++ {
			if maskimg.At(x, y) == trans {
				img.Set(x, y, color.Transparent)
			} else {
				img.Set(x, y, overdata.At(x, y))
			}
		}
	}

	CaptchaMap[imageId] = strconv.Itoa(imageRandX)
	emptyBuff := bytes.NewBuffer(nil)
	jpeg.Encode(emptyBuff, bgdata, nil)
	dist := make([]byte, 5000000)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())

	emptyBuff = bytes.NewBuffer(nil)
	png.Encode(emptyBuff, img)
	maskdist := make([]byte, 5000000)
	base64.StdEncoding.Encode(maskdist, emptyBuff.Bytes())

	captcha.Image = "data:image/jpeg;base64," + string(dist)
	captcha.Imageid = imageId
	captcha.Imagey = imageRandY
	captcha.Imagecrop = "data:image/png;base64," + string(maskdist)
	//fmt.Println(captcha)
	return captcha, nil
}
