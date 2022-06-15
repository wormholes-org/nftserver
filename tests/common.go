package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//const SrcUrl = "http://192.168.32.128:8081/v2/"
//const SrcUrl = "http://192.168.56.129:8081/v2/"
//var SrcUrl = "http://192.168.1.237:8091/v2/"
var SrcUrl = "http://192.168.1.235:10582/v2/"
//const SrcUrl = "http://192.168.1.8:8081/v2/"
//const SrcUrl = "http://192.168.4.237:8081/v2/"
//const SrcUrl = "https://192.168.1.237:9002/v2/"

type UserKeys struct {
	LogKey *ecdsa.PrivateKey 	//`json:"log_key"`
	WorkKey *ecdsa.PrivateKey	//`json:"work_key"`
}

type SaveKeys struct {
	LogKey []byte
	WorkKey []byte
}


type Response struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data struct {
		Hash 		string 	`json:"hash"`
		Secret 		string 	`json:"secret"`
		TimeStamp 	int64 	`json:"time_stamp"`
		Token 		string 	`json:"token"`
	} `json:"data"`
	Total_count int	`json:"total_count"`
}

func GetMd5Index(Secret, Hash string) int {
	var i int
	for {
		rawData := Secret + strconv.Itoa(i)
		sum := md5.Sum([]byte(rawData))
		hashString := hex.EncodeToString(sum[:])
		if hashString == Hash {
			break
		}
		i ++
	}
	return i
}

func Sign(data []byte, prv *ecdsa.PrivateKey) (string, error) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	sig, err := crypto.Sign(crypto.Keccak256([]byte(msg)), prv)
	if err != nil {
		fmt.Println("signature error: ", err)
		return "", err
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	return sigstr, err
}

func DelDataItem(data []byte) []byte {
	datastr := strings.ReplaceAll(string(data), "\"data\":[],", "")
	datastr = strings.ReplaceAll(datastr, "[", "")
	datastr = strings.ReplaceAll(datastr, "]", "")
	return []byte(datastr)
}

func HttpSendRev(url string, data string, token string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	if strings.Index(url, "https") != -1 {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("token", token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HttpGetSendRev(url string, data string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	if strings.Index(url, "https") != -1 {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("token", token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}


func HttpGetParamsSendRev(remoteUrl string, queryValues url.Values) ([]byte, error) {
	//params := url.Values{}
	parseURL, err := url.Parse(remoteUrl)
	if err != nil {
		log.Println("err")
	}
	parseURL.RawQuery = queryValues.Encode()
	urlPathWithParams := parseURL.String()
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}
	fmt.Println(string(b))
	return b, nil
}

func HttpSendSign(msg []byte, prv *ecdsa.PrivateKey) (string, error) {
	sig, err := Sign(msg, prv)
	if err != nil {
		fmt.Println("sign err=", err)
		return "", err
	}
	Index := bytes.LastIndex(msg, []byte("}"))
	msg[Index] = ','
	msg = append(msg, []byte("\"sig\":\"")...)
	msg = append(msg, sig...)
	msg = append(msg, []byte("\"}")...)
	return string(msg), nil
}

func GetUserAddr(c int) ([]UserKeys, error) {
	tKey, err := GetUserKeys("./key")
	if err != nil {
		fmt.Println("GetUserKeys() err= ", err)
		return nil, err
	}
	return tKey[:c], nil
}


func GenUserKeys(filePath string, c int) error {
	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	file.Truncate(0)
	var tKey []SaveKeys
	for i := 0; i < c; i++ {
		logKey, err := crypto.GenerateKey()
		if err != nil {
			fmt.Println("failed GenerateKey with.", err)
		}
		workKey, err := crypto.GenerateKey()
		if err != nil {
			fmt.Println("failed GenerateKey with.", err)
		}
		tKey = append(tKey, SaveKeys{crypto.FromECDSA(logKey), crypto.FromECDSA(workKey)})
	}
	data, err := json.Marshal(tKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = file.WriteAt(data, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetUserKeys(filePath string) ([]UserKeys, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var tKey []SaveKeys
	err = json.Unmarshal(buffer, &tKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var userKey []UserKeys
	for _, key := range tKey {
		logkey, err := crypto.ToECDSA(key.LogKey)
		if err != nil {
			return nil, err
		}
		workKey, err := crypto.ToECDSA(key.WorkKey)
		if err != nil {
			return nil, err
		}
		userKey = append(userKey, UserKeys{logkey, workKey})
	}
	return userKey, nil
}

