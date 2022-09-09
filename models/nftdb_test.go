package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const sqlsvrLcT = "admin:user123456@tcp(192.168.1.235:3306)/"

//
//const sqlsvrLcT = "admin:user123456@tcp(192.168.1.237:3306)/"

//
//const sqlsvrLcT = "admin:user123456@tcp(192.168.56.122:3306)/"

//
//const sqlsvrLcT = "demo:123456@tcp(192.168.56.129:3306)/"

//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
//var SqlSvrT = "admin:user123456@tcp(192.168.1.238:3306)/"
//const dbNameT = "mynftdb"
//const dbNameT = "tnftdb"
const dbNameT = "c0x7cfbfeab24f646ab96b31da34ac689bc71faf655"

//const dbNameT = "snftdb8012"

//const dbNameT = "nftdb8011"
//const dbNameT = "tttt"
//const dbNameT = "c0x5051580802283c7b053d234d124b199045ead750"

//const dbNameT = "c0x655b762461ba84b9c55de7346729594ce6361ac7"

//const dbNameT = "c0x544d5471284271f0cc0b48669d553c72a0877070"

const localtimeT = "?parseTime=true&loc=Local"

//const localtimeT = "?charset=utf8mb4&parseTime=True&loc=Local"

const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

func TestCreateDb(t *testing.T) {
	nd := new(NftDb)
	err := nd.InitDb(sqlsvrLcT, dbNameT)
	if err != nil {
		fmt.Printf("InitDb() err=%s\n", err)
	}
}

func TestDbMaxConnect(t *testing.T) {
	for i := 0; i < 2000; i++ {
		_, err := NewNftDb(sqldsnT)
		if err != nil {
			fmt.Printf("connet count=%d err=%s\n", i, err)
			break
		}
	}
	fmt.Println("end.")
}

func TestLoginNew(t *testing.T) {
	wd := sync.WaitGroup{}
	wd.Add(1)
	go func() {
		nd, err := NewNftDb(sqldsnT)
		if err != nil {
			fmt.Printf("connect database err = %s\n", err)
		}
		defer nd.Close()
		err = nd.LoginNew("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "sigdata")
		if err != nil {
			fmt.Printf("login err.\n")
		}
		wd.Done()
	}()
	wd.Add(1)
	go func() {
		nd, err := NewNftDb(sqldsnT)
		if err != nil {
			fmt.Printf("connect database err = %s\n", err)
		}
		defer nd.Close()
		err = nd.LoginNew("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "sigdata")
		if err != nil {
			fmt.Printf("login err.\n")
		}
		wd.Done()
	}()
	wd.Wait()
	fmt.Println("login test end.")
}

var testsync UserSyncMapList

func TestLogin(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.Login("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "sigdata")
	if err != nil {
		fmt.Printf("login err.\n")
	}
	err = nd.Login("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "sigdata")
	if err != nil {
		fmt.Printf("login err.\n")
	}
	err = nd.Login("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e166", "sigdata")
	if err != nil {
		fmt.Printf("login err.\n")
	}
	fmt.Println("login test end.")
}

func TestModifyUserInfo(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, err = nd.QueryUserInfo("bbbbbbbbbbbbbbbbbbbbb")
	if err != nil {
		fmt.Println("QueryUserInfo() err=", err)
	}
	//err = nd.ModifyUserInfo("bbbbbbbbbbbbbbbbbbbbb", "renameuser",
	//	"portrait", "my bio.", "test@test.com", "sigdata")
	//if err != nil {
	//	fmt.Println("ModifyUserInfo() err=", err)
	//}
	_, err = nd.QueryUserInfo("bbbbbbbbbbbbbbbbbbbbb")
	if err != nil {
		fmt.Println("QueryUserInfo() err=", err)
	}
	//err = nd.ModifyUserInfo("bbbbbbbbbbbbbbbbbcbbbb", "renameuser",
	//	"portrait", "my bio.", "test@test.com", "sigdata")
	//if err != nil {
	//	fmt.Println("ModifyUserInfo() err=", err)
	//}
}

func TestQueryUserInfo(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, err = nd.QueryUserInfo("0x572bcacb7ae32db658c8dee49e156d455ad59ec8")
	if err != nil {
		fmt.Println("QueryUserInfo() err=", err)
	}
}

/*func TestUpload(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.UploadNft("0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
		"md5",
		"url",
		"1000",
		"signdata",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"1631679689395",
		"0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
		"image",
		"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.UploadNft("0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
	"md5",
	"url",
	"2000",
	"signdata",
	"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
	"1631679689395",
	"0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
	"image",
	"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.UploadNft("useraddr",
		"md5",
		"url",
		"3000",
		"signdata",
		"contract22",
		"tokenid22",
		"ownaddr22",
		"image",
		"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.UploadNft("useraddr",
		"md5",
		"url",
		"5000",
		"signdata",
		"contract55",
		"tokenid55",
		"ownaddr55",
		"image",
		"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
}
*/

func TestBuyNft(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.BuyNft("mynft", "tradeSig", "sigdata", "contract11", "TokenId11")
	if err != nil {
		fmt.Printf("buyNft err=%s.\n", err)
	}
	err = nd.BuyNft("mynft", "tradeSig", "sigdata", "contract22", "TokenId22")
	if err != nil {
		fmt.Printf("buyNft err=%s.\n", err)
	}
}

func TestQueryNft(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	_, err = nd.QueryNft()
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
}

func TestNftbyUser(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	_, err = nd.QueryNftbyUser("mynft")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
}

func TestRenameTab(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	nd.db.Migrator().RenameTable("users", "user_infos")
}

//func TestTimePro(t *testing.T) {
//	TimeProc(sqldsnT)
//}

func TestMash(t *testing.T) {
	type test struct {
		Num int64 `json:"num"`
	}
	price, _ := strconv.ParseUint("", 10, 64)
	fmt.Println(price)
	tt := test{98708097097987098}
	marshal, _ := json.Marshal(tt)
	t.Log(string(marshal))
}

func TestFavorited(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.Like("useraddr", "contract", "tokenid", "sig")
	if err != nil {
		fmt.Printf("AddFavor err = %s\n", err)
	}
	err = nd.Like("useraddr", "contract11", "tokenid11", "sig")
	if err != nil {
		fmt.Printf("AddFavor err = %s\n", err)
	}
	_, err = nd.QueryNftFavorited("useraddr")
	if err != nil {
		fmt.Printf("QueryFavorited err = %s\n", err)
	}
	err = nd.DelNftFavor("useraddr", "contract11", "tokenid11")
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	_, err = nd.QueryNftFavorited("useraddr")
	if err != nil {
		fmt.Printf("QueryFavorited err = %s\n", err)
	}
}

func TestUserFavorited(t *testing.T) {

}

func TestSell(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "HighestBid", "paychan",
		1, 1001, 2000, "royalty", "美元", "false", "sigdate", "0569376186306", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "1", "1", 1200, "Tradesig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "1", "1", 1500, "TradeSig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	//test2
	err = nd.Sell("ownAddr11", "", "contract11", "tokenid11",
		"FixPrice", "paychan",
		1, 2001, 5000, "royalty", "use", "false", "sigdata", "0569376186306", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer1", "contract11", "tokenid11", "1", "1", 2100, "Tradesig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer2", "contract11", "tokenid11", "1", "1", 2200, "Tradesig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer3", "contract11", "tokenid11", "1", "1", 6300,
		"Tradesig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	//test3
	err = nd.Sell("ownAddr22", "", "contract22", "tokenid22", "HighestBid", "paychan",
		1, 6000, 8000, "royalty", "use", "false", "sigdata", "0569376186306", "tradeSig")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer1", "contract22", "tokenid22", "1", "1", 6100, "tradesig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer2", "contract22", "tokenid22", "1", "1", 6200, "TradeSig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer3", "contract22", "tokenid22", "1", "1", 6300, "tradesig", 0, "0569376186306", "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	nd.Close()
}

func TestMakeOffer(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	ExchangerAuth = "0x01842a2cf56400a245a56955dc407c2c4137321e"
	contracts.EthNode = "http://api.wormholestest.com:8561"
	//err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
	//	"",
	//	"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
	//	"0569376186306", "HighestBid", "paychan",
	//	1, 1001, 2000, "royalty", "美元", "false", "sigdate", "tradedate")
	//if err != nil {
	//	fmt.Printf("Sell() err = %s\n", err)
	//}
	err = nd.MakeOffer("0x0109cc44df1c9ae44bac132ed96f146da9a26b88", "0x01842a2cf56400a245a56955dc407c2c4137321e",
		"7401585102779", "1", "1", 11000000000, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	nd.Close()
}
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func TestGetSign(t *testing.T) {
	var message []byte = []byte("Hello World!")
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed GenerateKey with %s.", err)
	}
	//带有0x的私钥
	fmt.Println("private key have 0x   \n", hexutil.Encode(crypto.FromECDSA(key)))
	fmt.Println("public key have 0x   \n", hexutil.Encode(crypto.FromECDSAPub(&key.PublicKey)))
	fmt.Println("addr   \n", crypto.PubkeyToAddress(key.PublicKey).String())
	//不含0x的私钥
	fmt.Println("private key no 0x \n", hex.EncodeToString(crypto.FromECDSA(key)))
	sig, err := crypto.Sign(signHash(message), key)
	if err != nil {
		t.Errorf("signature error: %s", err)
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	addr, err := NftDb{}.GetEthAddr("Hello World!", sigstr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("addr=%x\n", addr)
}

func TestSignAppconf(t *testing.T) {
	file, err := os.OpenFile("D:\\temp\\app.conf", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bytesread)
	msg := string(buffer)
	tm := fmt.Sprintf(time.Now().String())
	msg = msg + "[time]\n" + "date = " + tm + "\n\n"

	var message []byte = []byte(msg)
	key, err := crypto.HexToECDSA("8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed")
	if err != nil {
		fmt.Println(err)
	}
	sig, err := crypto.Sign(signHash(message), key)
	if err != nil {
		t.Errorf("signature error: %s", err)
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	msg = msg + "#签名数据\n" + "[sig]\n" + "app.conf.sig = " + sigstr
	_, err = file.WriteAt([]byte(msg), 0)
	if err != nil {
		fmt.Println(err)
	}
}

func TestExchangePrv(t *testing.T) {
	SuperAdminPrv = DefSuperAdminPrv
	privateKey, err := crypto.HexToECDSA(SuperAdminPrv)
	if err != nil {
		fmt.Printf("InitSysParams() AdminListPrv err = %s\n", err)
		return
	}
	hexprv := crypto.FromECDSA(privateKey)
	hexprvstr := hexutil.Encode(hexprv)
	hexprvstr = hexprvstr[2:]
	fmt.Println(hexprvstr)
}

func TestVerifyAppconf(t *testing.T) {
	file, err := os.OpenFile("D:\\temp\\app.conf", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bytesread)
	msg := string(buffer)
	index := strings.Index(msg, "app.conf.sig = ")
	sig := msg[index+len("app.conf.sig = "):]
	var message []byte = []byte(msg[:strings.Index(msg, "#签名数据")])
	addr, err := NftDb{}.GetEthAddr(string(message), sig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addr)
}

func TestGetEthAddr(t *testing.T) {
	/*{
	  “msg": "Hello World!"
	  "address": "0x0109cc44df1c9ae44bac132ed96f146da9a26b88",
	  "msg": "0x48656c6c6f20576f726c6421",
	  "sig": "23ad293d6976499c11905c2c811502af9c47c2a0388bec4acb7cf2005554f39226a74d6aec36cdca868dd7ecf62fdd92888e2f9f45939f7f4450362eea1cb5ad1c",
	  "version": "3",
	  "signer": "MEW"
	}*/
	nd := new(NftDb)
	addr, err := nd.GetEthAddr("Hello World!", "0x23ad293d6976499c11905c2c811502af9c47c2a0388bec4acb7cf2005554f39226a74d6aec36cdca868dd7ecf62fdd92888e2f9f45939f7f4450362eea1cb5ad1c")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addr)
}

func TestQueryNftCurTransInfo(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nftTranInfo, err := nd.QuerySingleNft("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "9985")
	if err != nil {
		fmt.Println(err)
	}
	marshal, _ := json.Marshal(nftTranInfo)
	fmt.Printf("%s\n", string(marshal))
}

func TestDbPing(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	db, err := nd.db.DB()
	if db.Ping() != nil {

	}
}

func TestQueryMarketInfo(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	minfo, _ := nd.QueryMarketInfo()
	data, _ := json.Marshal(minfo)
	fmt.Println(data)
}

func TestInsertNft(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("TestInsertNft() connect database err = %s\n", err)
	}
	nftrecord := Nfts{}
	dberr := nd.db.Where("id = ?", 77).First(&nftrecord)
	if dberr.Error != nil {
		fmt.Println("TestInsertNft() bidprice not find nft err= ", dberr.Error)
		return
	}

	tokenid, _ := strconv.ParseInt(nftrecord.Tokenid, 10, 64)
	for i := 0; i < 100; i++ {
		nftcreade := Nfts{}
		nftcreade.NftRecord = nftrecord.NftRecord
		tokenid = tokenid + int64(i)
		nftcreade.Tokenid = strconv.FormatInt(tokenid, 10)
		/*dberr := nd.db.Model(&Nfts{}).Create(&nftcreade)
		if dberr.Error != nil {
			fmt.Println("TestInsertNft() bidprice not find nft err= ", dberr.Error)
			return
		}*/
		imagerr := SaveNftImage("d:/temp/image/", nftcreade.Contract, nftcreade.Tokenid, nftcreade.Image)
		if imagerr != nil {
			fmt.Println("TestInsertNft() SaveNftImage err= ", imagerr)
			return
		}
	}
	defer nd.Close()

}

func TestGetBalance(t *testing.T) {

}

func TestBuyResult(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nd1 := new(NftDb)
	err1 := nd1.ConnectDB(sqldsnT)
	if err1 != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd1.Close()
	/*
		 auctionRec.Ownaddr= 0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC
		5873 bidRecs.Bidaddr= 0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169
		5874 auctionRec.Contract= 0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F
		5875 auctionRec.Tokenid= 1631681392629
		5876 price= 50000000000000000
	*/
	if true {
		err = nd.BuyResult("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
			"0x8fbf399d77bc8c14399afb0f6d32dbe22189e169",
			"0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5",
			"1062183305419",
			"tradesig",
			"200000000", "sigData", "", "txhash")
		if err != nil {
			fmt.Println(err)
		}
		err = nd.BuyResult("",
			"0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c",
			"0xA1e67a33e090Afe696D7317e05c506d",
			"9161528579394",
			"tradesig",
			"", "sigData", "200", "txhash")
		if err != nil {
			fmt.Println(err)
		}
	} else {

		go nd.BuyResult("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
			"1631753648255",
			"tradesig",
			"20000000000", "sigData", "", "txhash")
		go nd1.BuyResult("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
			"1631753648255",
			"tradesig",
			"", "sigData", "", "txhash")
		select {}
	}
}

func TestQueryNftByFilter(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	/*
		{
			"filter": [],
			"match": "",
			"sort": [],
		}
	*/
	//filters := []StQueryField{
	//	/*{"collectcreator", "=", "0x572bcacb7ae32db658c8dee49e156d455ad59ec8"},
	//	{"collections", "=", "Buyer"},*/
	//}
	filters := []StQueryField{
		{
			"collectcreator",
			"=",
			"0x01842a2cf56400a245a56955dc407c2c4137321e",
		},
		{
			"collections",
			"=",
			"0000000.合集",
		},
	}
	sorts := []StSortField{}
	nftByFilter, count, err := nd.QueryNftByFilter(filters, sorts, "0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	filters = []StQueryField{
		{
			"selltype",
			"=",
			"FixPrice",
		},
		{
			"selltype",
			"=",
			"HighestBid",
		},
		{
			"offernum",
			">",
			"0",
		},
		{
			"createdate",
			">=",
			"1650355481",
		},
		{
			"collectcreator",
			"=",
			"0x01842a2cf56400a245a56955dc407c2c4137321e",
		},
		{
			"collections",
			"=",
			"0000000.合集",
		},
	}
	//sorts := []StSortField{{By: "createdate", Order: "desc"}}
	sorts = []StSortField{}
	nftByFilter, count, err = nd.QueryNftByFilter(filters, sorts, "0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	t.Logf("nft = %v %v\n", nftByFilter, count)
}

func TestQuerySnfChip(t *testing.T) {
	err := NewQueryCatch("192.168.1.235:6379", "user123456")
	fmt.Println(err)
	nd := new(NftDb)
	err = nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	snftChip, count, err := nd.QuerySnftChip("0x0109cc44df1c9ae44bac132ed96f146da9a26b88", "5640175101587", "0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	snftChip, count, err = nd.QuerySnftChip("0x01842a2cf56400a245a56955dc407c2c4137321e", "7607070612728", "10", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	snftChip, count, err = nd.QuerySnftChip("0x01842a2cf56400a245a56955dc407c2c4137321e", "7607070612728", "20", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	t.Logf("nft = %v %v\n", snftChip, count)
}

func hashMsg(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func recoverAddress(msg string, sigStr string) (*common.Address, error) {
	sigData, err := hexutil.Decode(sigStr)
	if err != nil {
		log.Println("recoverAddress() err=", err)
		return nil, err
	}
	if len(sigData) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := hashMsg([]byte(msg))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(*rpk)
	return &addr, nil
}

func TestDecodeMint(t *testing.T) {
	Buyer := contracts.Buyer1{}
	err := json.Unmarshal([]byte("{\"price\":\"0xd3c21bcecceda0000000\",\"exchanger\":\"0x62e0c8032fb51bc401558b58b1e7733276c1ec8a\",\"block_number\":\"0x8639\",\"sig\":\"0x4d07c38394bdf14590a7b4d8b0d432f493c40c08369cac0f1d920580ad0cd72372bd58714e029398f5909cde8fc116d5a763333fae42db451e492766a47cbf261c\"}"), &Buyer)
	if err != nil {
		log.Println("AuthExchangerMint()  Unmarshal() err=", err)
		return
	}
	msg := Buyer.Price + Buyer.Exchanger + Buyer.Blocknumber
	toaddr, err := recoverAddress(msg, Buyer.Sig)
	if err != nil {
		log.Println("GetBlockTxs() recoverAddress() err=", err)
		return
	}
	fmt.Println(toaddr)
	Seller := contracts.Seller2{}
	err = json.Unmarshal([]byte("{\"price\":\"0x38d7ea4c68000\",\"royalty\":\"0x64\",\"meta_url\":\"7b226d657461223a222f697066732f516d63666e4c647174704b4c55794c325155716f36546a4438695857667161545448716d6a503470636f78656d47222c22746f6b656e5f6964223a2231323234303239323432313134227d\",\"exclusive_flag\":\"1\",\"exchanger\":\"0x62e0c8032fb51bc401558b58b1e7733276c1ec8a\",\"block_number\":\"0x168146307200100\",\"sig\":\"0xf267eaa1edd2d6a4f438e11a65bc0836e6296cd9d1a731737846eb27d34247a45e85505e4eabf7e2cbd0f9ea30b6609a9b7ac7beb1fcb53b98a6fadaa713ad741b\"}"), &Seller)
	if err != nil {
		log.Println("AuthExchangerMint()  Unmarshal() err=", err)
		return
	}
	msg = Seller.Price + Seller.Royalty + Seller.Metaurl + Seller.Exclusiveflag +
		Seller.Exchanger + Seller.Blocknumber
	fromAddr, err := recoverAddress(msg, Seller.Sig)
	if err != nil {
		log.Println("GetBlockTxs() recoverAddress() err=", err)
		return
	}
	fmt.Println(fromAddr)
	var nftmeta contracts.NftMeta
	metabyte, jsonErr := hex.DecodeString("7b226d657461223a222f697066732f516d63666e4c647174704b4c55794c325155716f36546a4438695857667161545448716d6a503470636f78656d47222c22746f6b656e5f6964223a2231323234303239323432313134227d")
	if jsonErr != nil {
		log.Println("GetBlockTxs() hex.DecodeString err=", jsonErr)
	}
	jsonErr = json.Unmarshal(metabyte, &nftmeta)
	if jsonErr != nil {
		log.Println("GetBlockTxs() NftMeta unmarshal type err=", jsonErr)
	}
}

func TestQueryNftByFilterNew(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	NewQueryCatch("192.168.1.235:6379", "user123456")
	GetRedisCatch().SetDirtyFlag(AllDirty)
	nfilters := []StQueryField{
		/*{
			"selltype",
			"=",
			"FixPrice",
		},*/
		{
			"selltype",
			"=",
			"HighestBid",
		},
		/*{
			"offernum",
			">",
			"0",
		},*/
		/*{
			"createdate",
			">=",
			"1654149924",
		},
		{
			"collectcreator",
			"=",
			"0x01842a2cf56400a245a56955dc407c2c4137321e",
		},
		{
			"collections",
			"=",
			"0x8000000000000000000000000000000000000.合集",
		},*/
		/*{
			"sellprice",
			">=",
			"1000000000",
		},
		{
			"sellprice",
			"<=",
			"6000000000",
		},*/
	}
	nfilters = []StQueryField{
		{
			"collectcreator",
			"=",
			"0xbe8c75133a7e4f29b7cdc15d4a45f7593a4f8898",
		},
		{
			"collections",
			"=",
			"测试1",
		},
	}
	/*nfilters = []StQueryField{
		{
			"selltype",
			"=",
			"FixPrice",
		},
		{
			"selltype",
			"=",
			"HighestBid",
		},
		{
			"offernum",
			">",
			"0",
		},
		{
			"createdate",
			">=",
			"1650609876",
		},
		{
			"sellprice",
			">=",
			"10000000000",
		},
		{
			"sellprice",
			"<=",
			"20000000000",
		},
	}*/

	nfilters = []StQueryField{
		//{
		//	"collectcreator",
		//	"=",
		//	"0x7a149f02e5e4571c42d5cf69b4ccb5772fa1b275",
		//},
		//{
		//	"collections",
		//	"=",
		//	"collect_test_0",
		//},
		//{
		//	"selltype",
		//	"=",
		//	"FixPrice",
		//},
		//{
		//	"selltype",
		//	"=",
		//	"HighestBid",
		//},
		//{
		//	"categories",
		//	"=",
		//	"Music",
		//},
		/*{
			"sellprice",
			">=",
			"10000000000",
		},
		{
			"sellprice",
			"<=",
			"20000000000",
		},*/
		{
			"createdate",
			">=",
			"1662009658",
		},
		/*{
			"offernum",
			">",
			"0",
		},*/
	}
	//sorts := []StSortField{{By: "createdate", Order: "desc"}}
	//nfilters = []StQueryField{}
	sorts := []StSortField{
		{
			"sellprice",
			"asc",
		},
		/*{
			"verifiedtime",
			"asc",
		},*/
	}
	sorts = []StSortField{}
	nftByFilter, count, err := nd.QueryNftByFilterNftSnft(nfilters, sorts, "snft", "0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	nftByFilter, count, err = nd.QueryNftByFilterNew(nfilters, sorts, "nftsnft", "5", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	t.Logf("nft = %v %v\n", nftByFilter, count)
}

func TestNftDb_BuyingNft(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	contracts.EthNode = "http://43.129.181.130:8561"
	contracts.SuperAdminAddr = "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
	ExchangerAuth = "{\"exchanger_owner\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"to\":\"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900\",\"block_number\":\"0x2540be400\",\"sig\":\"0x7f1ca96714208959c5a75bdbf4770893b76b13c0bca26da2086c3365e537d57444f79b31498301c5c1d55400eec4b469c83a88a527159112f27ff934c222e4191b\"}"
	price, _ := strconv.ParseUint("11000000000", 10, 64)
	buyerSig := `"buyer_sig":"{\"price\":\"0x98a7d9b8314c0000\",\"exchanger\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"block_number\":\"0x603f87\",\"sig\":\"0x53f3c8018682637d73414b416fe62c825af97e6b2cbe1390ed4d34bc33c136162285b117087fdad81352cf9da7a2edc11d5be093255cb3e847862e95853ad0c21c\"}"`
	sellerSig := `"{\"price\":\"0x98a7d9b8314c0000\",\"royalty\":\"0xc8\",\"meta_url\":\"7b226d657461223a222f697066732f516d6445736b6968627a7a53626f6a3843693845395147707a5833507378316363694e6b47505739505051716d66222c22746f6b656e5f6964223a2239353734323632373430333136227d\",\"exclusive_flag\":\"1\",\"exchanger\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"block_number\":\"0x603f3a\",\"sig\":\"0x57e167bffd120f88a3563c19a809553de7454f3775da1fcb4c09ff99ee3397df554722cdbffc27d649d744d333edb9427e02f6199d6bd046b2787c443a7b5ba81b\"}"`
	//_ = nd.BuyingNft("0x0109CC44df1C9ae44Bac132eD96f146Da9A26B88",
	//	"0x0109CC44df1C9ae44Bac132eD96f146Da9A26B88",
	//	"0x01842a2cf56400a245a56955dc407c2c4137321e",
	//	"9574262740316", price,
	//	buyerSig, "", sellerSig)
	buyerSig = "{\"price\":\"0x98a7d9b8314c0000\",\"nft_address\":\"0x0000000000000000000000000000000000000001\",\"exchanger\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"block_number\":\"0x6040d0\",\"seller\":\"0x0109cc44df1c9ae44bac132ed96f146da9a26b88\",\"sig\":\"0xe0864c8fd9a6929639d4b41e22f1f8296e420a8e162f61f197402c4f346ccf2d02a47ed9ce757bcd5c1a77b9f1051b31641447c79d473bd8fbd9f0495af6070e1b\"}"
	sellerSig = "{\"price\":\"0x98a7d9b8314c0000\",\"nft_address\":\"0x0000000000000000000000000000000000000001\",\"exchanger\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"block_number\":\"0x60470c\",\"sig\":\"0x3b372442331a620ff98055d3c647ec23df2872cb0e08067a655c13fa5d5d4c5e373600af52398442940e6f51300b77c9921534fff8bb0d636e9b5f20a71afddd1b\"}"

	_ = nd.BuyingNft("0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
		"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
		"0x01842a2cf56400a245a56955dc407c2c4137321e",
		"9574262740316", price,
		buyerSig, "", sellerSig)

}

func TestQueryStageList(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	stageList, count, err := nd.QueryStageList("0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	fmt.Println(stageList, count)
}

type HttpResponseData struct {
	Code       string      `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
	TotalCount uint64      `json:"total_count"`
}

func TestQueryStageCollection(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var httpResponseData HttpResponseData
	collections, err := nd.QueryStageCollection("0x800000000000000000000000000000000000")
	if err != nil {
		httpResponseData.Data = []interface{}{}
	} else {
		httpResponseData.Code = "200"
		httpResponseData.Data = collections
	}
	responseData, _ := json.Marshal(httpResponseData)
	fmt.Println(responseData)
	collections, err = nd.QueryStageCollection("0x800000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	fmt.Println(collections)
	collections, err = nd.QueryStageCollection("0x800000000000000000000000000000000001")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	fmt.Println(collections)
}

func TestQueryStageSnft(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	spendT := time.Now()
	stageList, err := nd.QueryStageSnft("0x800000000000000000000000000000000024", "测试6")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	fmt.Println("spend time =", time.Now().Sub(spendT))
	spendT = time.Now()
	stageList, err = nd.QueryStageSnft("0x800000000000000000000000000000000000", "0x8000000000000000000000000000000000001.合集")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	fmt.Println("spend time =", time.Now().Sub(spendT))
	fmt.Println(stageList)
}

func TestTimeStamp(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var nftRecs []Nfts
	//errr := nd.db.Where("createaddr = ?", "useraddr").Distinct("createaddr").Find(&nftRecs)
	errr := nd.db.Where("createaddr = ?", "useraddr").Find(&nftRecs)
	//errr := nd.db.Model(&Nfts{}).Find(&nftRecs)

	if errr.Error != nil {
		fmt.Println(err.Error())
	}
	/*fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))*/
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		//fmt.Println(rand.Int63())
		s := fmt.Sprintf("%d", rand.Int63())
		if len(s) > 16 {
			continue
		}
		s1 := s[len(s)-13:]
		fmt.Println(s1, "=", len(s))
		//fmt.Println(rand.New(rand.NewSource(time.Now().UnixNano())).Int63())
	}
}

func TestSysParams(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	data, err := nd.QuerySysParams()
	if err != nil {
		fmt.Println(err)
	}
	data.Icon = ""
	fmt.Println(data)

}

func TestInitSysParams(t *testing.T) {
	InitSysParams(sqldsnT)
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	RoyaltyLimit = 10000
	ImageDir = "D:\\home\\user1\\chengdu"
	err = nd.UploadNft(
		"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
		"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
		"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"0000000.合集",
		Default_image,
		"true",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.SetSysParams(SysParamsInfo{Lowprice: "1000000"})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	nd.db.Migrator().DropColumn(&SysParams{}, "exchangeaddr")

	err = nd.SetSysParams(SysParamsInfo{NFT1155addr: "0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC", Adminaddr: "", Lowprice: "1000000"})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	err = nd.SetSysParams(SysParamsInfo{NFT1155addr: "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F", Adminaddr: "", Lowprice: "100000000"})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	err = nd.SetSysParams(SysParamsInfo{NFT1155addr: "", Adminaddr: "", Lowprice: "100000000"})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	nd.Close()
}

func TestCollections(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.NewCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"test",
		"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAK8ArwDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD+/igAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAo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gAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoA//9k=",
		"",
		"",
		"test.",
		"art",
		"sigedata", "",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	err = nd.ModifyCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"test", "img", "contract_type", "contract_addr",
		"test desc.", "art", "sig string", "")
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	err = nd.ModifyCollectionsImage("test",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAK8ArwDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD+/igAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAo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gAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoA//9k=",
		"sig string")
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
}

func TestUploadNftNew(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	RoyaltyLimit = 10000
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			NFT1155Addr = "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
		} else {
			NFT1155Addr = "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F"
		}
		err = nd.UploadNft(
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"md5 string",
			"name string",
			"desc string",
			"meta string",
			"source_url string",
			"",
			"",
			"categories string",
			"",
			"asset_sample string",
			"true",
			"2",
			"1",
			"sig string")
		if err != nil {
			fmt.Println("UploadNft err=", err)
		}
	}
}

func TestForeignContract(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.NewCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"foreign-contract-test",
		"",
		"",
		"0x9e2576747C2525062a77667E4E88A97b6034C461",
		"foreign-test.",
		"art",
		"sigedata", "",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	RoyaltyLimit = 10000
	err = nd.UploadNft(
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"foreign-contract-test",
		"asset_sample string",
		"false",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Println("UploadNft err=", err)
	}
}

func TestDbQuery(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var count int64
	dberr := nd.db.Model(Nfts{}).Where("contract = ? ", "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169").Count(&count)
	dberr = nd.db.Model(Nfts{}).Where("contract = ? ", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F").Count(&count)

	var nfttab []Nfts
	dberr = nd.db.Model(Nfts{}).Where("contract = ? ", "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169").Limit(2).Offset(2).Find(&nfttab)
	dberr = nd.db.Model(Nfts{}).Where("contract = ? ", "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169").Limit(2).Offset(5).Find(&nfttab)
	dberr = nd.db.Where("contract = ? ", "").First(&nfttab)
	if dberr.Error == nil {
		fmt.Println("UploadNft() err=nft already exist.")
	}
}

func TestQueryUserCollectionList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, err = nd.QueryUserCollectionList("0x8fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"0", "5")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestQueryUserFavorited(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, err = nd.QueryUserFavoriteList("0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"0", "5")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestModifyCollectionsImage(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nn := time.Now().Unix()
	fmt.Println(nn)
	err = nd.ModifyCollectionsImage("test", "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"modify", "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestQueryUserTradingHistory(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, err = nd.QueryUserTradingHistory("0x572bcAcB7ae32Db658C8dEe49e156d455Ad59eC8",
		"0", "5")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestLike(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.Like("0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F", "1632799124069", "sig")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestSearch(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	//testCond1 := "name"
	testCond2 := "0x"

	//searchData1, _ := nd.Search(testCond1)
	searchdata2, _ := nd.Search(testCond2)
	//for _, data := range searchData1 {
	//	for _, data1 := range data.CollectsRecords {
	//		t.Log(data1)
	//	}
	//	for _, data1 := range data.UserAddrs {
	//		t.Log(data1)
	//	}
	//	for _, data1 := range data.NftsRecords {
	//		t.Log(data1)
	//	}
	//}
	for _, data := range searchdata2 {
		for _, data1 := range data.CollectsRecords {
			t.Log(data1)
		}
		for _, data1 := range data.UserAddrs {
			t.Log(data1)
		}
		for _, data1 := range data.NftsRecords {
			t.Log(data1)
		}
	}
}

func TestBidPriceWithBuy(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	/*err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "HighestBid", "paychan",
		1, 1001, 2000, "royalty","美元", "false", "sigdate", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}*/
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e160", "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5",
		"7070595686952", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e161", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.BuyResult("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c",
		"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162",
		"0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350",
		"tradesig",
		"200000000", "sigData", "", "txhash")
	if err != nil {
		fmt.Println(err)
	}
}

func TestSignal(t *testing.T) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}

func TestBidPriceWithSell(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e160", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e161", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162",
		"",
		"0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "HighestBid", "paychan",
		1, 1001, 2000, "royalty", "美元", "false", "sigdate", "0569376186306", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
}

func TestBidPriceWithTime(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e160", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", time.Now().Unix()+10, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e161", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", time.Now().Unix()+1000, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "0569376186306", "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162",
		"",
		"0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "HighestBid", "paychan",
		1, 1001, 2000, "royalty", "美元", "false", "sigdate", "0569376186306", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
}

func TestQueryMarketTradingHistory(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	sorts := []StSortField{{By: "transtime", Order: "desc"}}

	history, i, err := nd.QueryMarketTradingHistory(nil, sorts, "0", "10")

	t.Log(history)
	t.Log(i)
	t.Log(err)
}

func TestAnnouncements(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.SetAnnouncement("title one", "content one")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title two", "content two")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title three", "content three")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title four", "content three")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title five", "content five")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title six", "content six")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title seven", "content seven")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title eight", "content eight")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	_, err = nd.QueryAnnouncement()
	if err != nil {
		fmt.Println("insert announcement.")
	}
}

func TestSearchSql(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var useroffer []UserOffer
	sql := "SELECT biddings.contract as Contract, biddings.tokenid as Tokenid, biddings.price as Price, " +
		"biddings.count as Count, biddings.bidtime as Bidtime FROM biddings LEFT JOIN nfts ON biddings.contract = nfts.contract AND biddings.tokenid = nfts.tokenid WHERE ownaddr = ? AND biddings.deleted_at is null"
	sql = sql + " limit 1, 2"
	db := nd.db.Raw(sql, "0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5").Scan(&useroffer)
	if db.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
	}
	var count int64
	sql = "SELECT biddings.contract as Contract, biddings.tokenid as Tokenid, biddings.price as Price, " +
		"biddings.count as Count, biddings.bidtime as Bidtime FROM biddings LEFT JOIN nfts ON biddings.contract = nfts.contract AND biddings.tokenid = nfts.tokenid WHERE ownaddr = ? AND biddings.deleted_at is null"
	db = nd.db.Raw(sql, "0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5").Count(&count)
	if db.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
	}
}

func TestIsValidCategory(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	var category1 = "virtual_worlds"
	validCategory1 := nd.IsValidCategory(category1)
	var category2 = "virtual"
	validCategory2 := nd.IsValidCategory(category2)
	t.Log("validCategory1=", validCategory1, "validCategory2=", validCategory2)
}

func TestQueryCollectionInfo(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	info, _ := nd.QueryCollectionInfo("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"实用合集")
	t.Log(info)
}

func TestQueryHomePage(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	spendT := time.Now()

	fmt.Println("spend time = ", time.Now().Sub(spendT))
	page, err := nd.QueryHomePage(true)
	page, err = nd.QueryHomePage(false)

	t.Log(page)
	t.Log(err)
}

func TestMultiQueryHomePage(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	testCount := 10
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			nd, err := NewNftDb(sqldsnT)
			if err != nil {
				fmt.Printf("connect database err = %s\n", err)
			}
			defer nd.Close()
			_, err = nd.QueryHomePage(true)
			if err != nil {
				fmt.Println("TestMultiQueryHomePage err=", err)
			}
			fmt.Println("TestMultiQueryHomePage i=", i)
		}(i)
	}
	wd.Wait()
}

func TestConvert(t *testing.T) {
	m := -19
	mstr := strconv.Itoa(m)
	u64, err := strconv.ParseUint(mstr, 10, 64)
	fmt.Println(u64, err)
	mstr = "ffffabdcdef"
	u64, err = strconv.ParseUint(mstr, 16, 64)
	fmt.Println(u64, err)

	mstr = ""
	u64, err = strconv.ParseUint(mstr, 10, 64)
	fmt.Println(u64, err)
	data, err := strconv.Atoi(mstr)
	fmt.Println(data, err)
	mstr = "ffffabdcdef"
	u64, err = strconv.ParseUint(mstr, 16, 64)
	fmt.Println(u64, err)
}

func TestConvertValid(t *testing.T) {
	err := IsIntDataValid("")
	if err != true {
		fmt.Println("datat err")
	}
	err = IsPriceValid("")
	if err != true {
		fmt.Println("datat err")
	}
}

func TestName(t *testing.T) {
	valid, errmsg, err := AmountValid(100000, "0xc9a9caa0147adc101138920ac7905ca6b62e9a2a")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(valid, errmsg)
}
func TestQueryUserOfferList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, _ = nd.QueryUserOfferList("0x572bcacb7ae32db658c8dee49e156d455ad59ec8",
		"0", "10")
}

func TestQueryNftCollectionList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	u1, _, _ := nd.QueryNFTCollectionList("0", "25")
	u2, _, _ := nd.QueryNFTCollectionList("0", "10")
	u3, _, _ := nd.QueryNFTCollectionList("10", "10")
	u4, _, _ := nd.QueryNFTCollectionList("20", "10")
	fmt.Println(u1, u2, u3, u4)
}

func TestUpdateBlockNumber(t *testing.T) {
	UpdateBlockNumber(sqldsnT)
}

func TestLog(t *testing.T) {
	//log.SetPrefix("TRACE: ")
	//log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	// Println writes to the standard logger.
	log.Println("message")

	// Fatalln is Println() followed by a call to os.Exit(1).
	log.Fatalln("fatal message")

	// Panicln is Println() followed by a call to panic().
	log.Panicln("panic message")
}

func TestImgTools(t *testing.T) {
	_, _, err := ParseBase64Type(Testimage)
	SaveNftImage("./test", "0xaaaaaaaaaaaa", "9999", Testimage)
	fmt.Println(err)
	SavePortrait("./test", "0x1AbCDERFG", "")
	SaveCollectionsImage("./test", "0x1AbCDERFG", "test1115", Default_image)
	SaveNftImage("./test", "0xaaaaaaaaaaaa", "9999", Default_image)
}

func TestAdminModify(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	err = nd.ModifyAdmin("0x572bcacb7ae32db658c8dee49e156d455ad59ec8", "nft", "1")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	err = nd.ModifyAdmin("0x572bcacb7ae32db658c8dee49e156d455ad59ec9", "nft", "1")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	err = nd.ModifyAdmin("0x572bcacb7ae32db658c8dee49e156d455ad59e10", "admin", "7")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	err = nd.ModifyAdmin("0x572bcacb7ae32db658c8dee49e156d455ad59e11", "admin", "6")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	err = nd.ModifyAdmin("0x572bcacb7ae32db658c8dee49e156d455ad59e12", "kyc", "5")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	err = nd.ModifyAdmin("0x572bcacb7ae32db658c8dee49e156d455ad59e15", "kyc", "2")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	count, admins, err := nd.QueryAdmins("nft", "0", "10")
	fmt.Println(count, admins)
	count, admins, err = nd.QueryAdmins("kyc", "0", "10")
	fmt.Println(count, admins)
	count, admins, err = nd.QueryAdmins("admin", "0", "10")
	fmt.Println(count, admins)
	var dellst DelAdmiList
	dellst.DelAdmins = append(dellst.DelAdmins, "0x572bcacb7ae32db658c8dee49e156d455ad59e15")
	dellst.DelAdmins = append(dellst.DelAdmins, "0x572bcacb7ae32db658c8dee49e156d455ad59e15")
	dellst.DelAdmins = append(dellst.DelAdmins, "0x572bcacb7ae32db658c8dee49e156d455ad59e15")
	dellst.DelAdmins = append(dellst.DelAdmins, "0x572bcacb7ae32db658c8dee49e156d455ad59e15")
	dellst.DelAdmins = append(dellst.DelAdmins, "0x572bcacb7ae32db658c8dee49e156d455ad59e15")
	dellst.DelAdmins = append(dellst.DelAdmins, "0x572bcacb7ae32db658c8dee49e156d455ad59e15")
	kd, err := json.Marshal(dellst)
	fmt.Println(kd)
	//{"del_admins":[["0x572bcacb7ae32db658c8dee49e156d455ad59e15","nft"],["0x572bcacb7ae32db658c8dee49e156d455ad59e15","admin"]]}
	deladmins := "[[\"0x572bcacb7ae32db658c8dee49e156d455ad59e15\",\"nft\"],[\"0x6e40b6deb1671b48b8b7efecac58b21f4f96468a\",\"admin\"]]"
	err = nd.DelAdmins(deladmins)
	if err != nil {
		t.Errorf("DelAdmins() err=%s", err)
	}
}

func TestQueryAdminByAddr(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	admins, err := nd.QueryAdminByAddr("0x572bcacb7ae32db658c8dee49e156d455ad59ec8")
	if err != nil {
		t.Errorf("ModifyAdmin() err=%s", err)
	}
	admins, err = nd.QueryAdminByAddr("0xa1e67a33e090afe696d7317e05c506d7687bb2e5")
	if err != nil {
		t.Errorf("QueryAdminByAddr() err=%s", err)
	}
	admins, err = nd.QueryAdminByAddr("0x7fbc8ad616177c6519228fca4a7d9ec7d1804900")
	if err != nil {
		t.Errorf("QueryAdminByAddr() err=%s", err)
	}
	fmt.Println(admins)
}
func TestGetSysParams(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	admins, err := nd.GetSysParam("adminaddr")
	if err != nil {
		t.Errorf("GetSysParams() err=%s", err)
	}
	fmt.Println(admins)
}
func TestPeriodAccedEth(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	uerr := nd.db.Model(&SnftPhase{}).Where("meta =? ", "/ipfs/QmY3K7qboPzZayGU8dfwVZw8HfEdriaxuvKwQENagrNMmp").Update("accedeth", "true")
	if uerr.Error != nil {
		fmt.Println(uerr.Error)
	}
	fmt.Println(uerr.RowsAffected)
	uerr = nd.db.Model(&SnftPhase{}).Where("accedeth =? and meta <>?", "false", "/ipfs/QmY3K7qboPzZayGU8dfwVZw8HfEdriaxuvKwQENagrNMmp").Update("accedeth", "")
	if uerr.Error != nil {
		fmt.Println(uerr.Error)
	}
}
func TestGetNftMarketInfo(t *testing.T) {

	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	fmt.Println(strings.ToLower("0xCFB25DF91cC483cAAFF30d5301fFEede1c59562F"))
	//nerr := nd.db.Last(&SysParams{}).Updates(&SysParamsRec{Blocknumber: 1, Scannumber: 1})

	nerr := nd.db.Last(&SysParams{}).Updates(map[string]interface{}{"blocknumber": 1, "scannumber": 1})
	if nerr.Error != nil {
		fmt.Println("TranSnft() update Blocknumber err= ", nerr.Error)
	}
	data, gerr := nd.GetNftMarketInfo()
	if gerr != nil {
		fmt.Println(gerr)
	}
	fmt.Println(data)
}

func TestGetBuyingParams(t *testing.T) {
	fmt.Println(len("4306d6a4d99b989f347784dd1039039746173f851555e006f9e3e6087ba9142850947206a3482becae4658c3ba6ff59ccb58b743d4d76ad6717ed24707b16a6e1b"))
	fmt.Println(strings.ToLower("0x70Eb3d3f80b577e9C3954d04b787c40b763a369B"))
	batchkey, serr := crypto.HexToECDSA("59ff5f189705c99b9034359ae8222aec959bfd7f0d6640ed41980319f049447f")
	if serr != nil {
		fmt.Println(serr)
	}
	fmt.Println(batchkey)
	//datasss := "{\"user_addr\":\"0x9aa8fef730ebf39660cc444919e9f314645e155b\",\"buying_param\":\"[{\\\"user_addr\\\":\\\"0x9aa8fef730ebf39660cc444919e9f314645e155b\\\",\\\"buyer_addr\\\":\\\"0x9aa8fef730ebf39660cc444919e9f314645e155b\\\",\\\"contract_addr\\\":\\\"0x01842a2cf56400a245a56955dc407c2c4137321e\\\",\\\"token_id\\\":\\\"4046121528691\\\",\\\"price\\\":100,\\\"buyer_sig\\\":\\\"0x0b30062094b32372435b4ac062f7bf911074fe5e6ce563ba7366477db9b426bd1ee6b19b1c1896cf8328ec0ed48a0bfdb5587eff7c0a2ebd4afc7324bfd7ec3c1b\\\",\\\"vote_stage\\\":\\\"1\\\",\\\"seller_sig\\\":\\\"0x5d51c4953f11ad77e62be5d856527bda0c29ba3b4dbec9edf67f2ef12ec3de683c64fe79af986fa91ca7ba0e5034c7ac4f33b1697024b7a229640130ad57164e1b\\\"}]\"}"
	dm := map[string]string{"user_addr": "0x9aa8fef730ebf39660cc444919e9f314645e155b", "buying_param": "[{\\\"user_addr\\\":\\\"0x9aa8fef730ebf39660cc444919e9f314645e155b\\\",\\\"buyer_addr\\\":\\\"0x9aa8fef730ebf39660cc444919e9f314645e155b\\\",\\\"contract_addr\\\":\\\"0x01842a2cf56400a245a56955dc407c2c4137321e\\\",\\\"token_id\\\":\\\"4046121528691\\\",\\\"price\\\":100,\\\"buyer_sig\\\":\\\"0x0b30062094b32372435b4ac062f7bf911074fe5e6ce563ba7366477db9b426bd1ee6b19b1c1896cf8328ec0ed48a0bfdb5587eff7c0a2ebd4afc7324bfd7ec3c1b\\\",\\\"vote_stage\\\":\\\"1\\\",\\\"seller_sig\\\":\\\"0x5d51c4953f11ad77e62be5d856527bda0c29ba3b4dbec9edf67f2ef12ec3de683c64fe79af986fa91ca7ba0e5034c7ac4f33b1697024b7a229640130ad57164e1b\\\"}]"}
	//sig, err := Sign([]byte(sigData), workKey)
	dj, _ := json.Marshal(dm)
	bd := []byte(string(dj))
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(bd), bd)
	sigs, serr := crypto.Sign(crypto.Keccak256([]byte(msg)), batchkey)
	if serr != nil {
		fmt.Println("signature error: ", serr)
	}
	sigs[64] += 27
	sigstr := hexutil.Encode(sigs)
	fmt.Println(sigstr)

	if serr != nil {
		fmt.Println(serr)
	}
	var Buying []BuyingParams
	addr := []string{
		"0x9d9a7dbd7e1e731e36af665786f9a0578ce6aebe", "0x7fbc8ad616177c6519228fca4a7d9ec7d1804900", "0x9b1168d0ba701448d6b9ec344e986d4a2f16971e",
		"0x80800580228ee590940fb75e86c8997611f079b3", "0xae407d12944877a969f38d5f2e313fe739271736", "0xb8b2c20bfdace578e716b594a33367c3036a7f0f",
	}
	token := []string{
		"8251060341162", "8233886113499", "7041514922311", "2662474950581", "7741799044198",
		"2091858617664", "3431090137985", "4548576859231", "8114837095171", "6024977233011",
		"1153977244697", "1894851936255", "6365531960120", "1520322497495", "9407391814828",
		"2236998921270",
	}
	sig := []string{
		"0x48f880bb5a4e6c9cf6e8170a407c06081c1f431664530f697aafc6ad7e63d92409a785bc03bb2c0cb67fca1bbd4ff6b08e50d1aa5ac73d5b4c5cf74c17a57bd01c",
		"0xf91ac6637539d0bf2797ab7bbe65234db94e64b363e938365a4ed21dd21700264dd8b9211f729d39b08562ba1e63d038b6d4104a0a41cbd3e4f7e9136f8cfa621b",
		"0x0b30062094b32372435b4ac062f7bf911074fe5e6ce563ba7366477db9b426bd1ee6b19b1c1896cf8328ec0ed48a0bfdb5587eff7c0a2ebd4afc7324bfd7ec3c1b",
		"0xce10f7adaea1fba71374584aef9c697a9cd2829963d27fedeb08d1934a96d36872f6c4c02d931eac1136d1b3e0777406e5fa04d3877413c34cd65a7ab28492c91c",
		"0x5d51c4953f11ad77e62be5d856527bda0c29ba3b4dbec9edf67f2ef12ec3de683c64fe79af986fa91ca7ba0e5034c7ac4f33b1697024b7a229640130ad57164e1b",
		"0xaac3852895ba0d1b98fe0ca6422dfdeb83b9a30352d47bd09c6bfcf1d668844a5c4c3f77d0f19961351e6233b43df823f5a7b4236fc3562fc98260af3b1eda871c",
		"0xa89d5bcd540d2e397deacb95bed4c9ef93a2c2e0b06c4659d9c258dc77003b6c0af4909eb1d663105982bd72bcfd77fb8ac0200d6acee41714868979ec1572db1c",
		"0x6beb749a3b77e43849cd1726f96c9e45abbb4f3081c17f66a04c7853471637981c2f844910b00e55e74e74c5a2392b605542a3f5903325d94f74b22d8a2967bb1c",
		"0x0edbaf0a57c8c92ab8e1c4129da3fb83f0adf8f2418b0e6cc61a7bfd2b25e41e10e92fd1b9e94c048fd6d4d306d7562cbb5d50a047b9b2cfde2b3c14b1799a971c",
		"0x54e52aa1a735280e3cb553fe7eadfe8c9c375193fca4ba8146f36555e20b3e872f654edf89accede87e2e2b7dd151eb8d02c29d57448d8bc49c1a21d18902bfc1c",
		"0xe067bd64c0e812d2757670f898d646b340479d398c3226a18e5af0976396f55305d860e6e394b9d37f6309d86bcb5f4dda544da2ef0b347b869e5a9108c9a4891c",
		"0x511b84e7fdb2c4c7fb83cdd16155f04f61111cb5f3b0effd2668e7ac25020eed353fa3d0f73e12ff2d2e080177a0a18b8abb4cfc2c1ac6b33388de74c47dec4d1c",
		"0xe719e460d8a34ff80ea882d143142fc788911174df1facd2600c37462691a1c5015adda7ba0c915ebe53f0e9009774274093fe11091e19757d69707fc93302c71c",
		"0xb1f86c8d9851d2db3a74b6a70742590148d7012a441cad653846c107ad4d2cee7896ac0e5876ef5148ba4c6436154f46145034508ad08fd791e4f867833214a01c",
	}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		var buyingparam BuyingParams
		buyingparam.UserAddr = addr[rand.Intn(6)]
		buyingparam.BuyerAddr = addr[rand.Intn(6)]
		buyingparam.ContractAddr = addr[rand.Intn(6)]
		buyingparam.TokenId = token[rand.Intn(16)]
		buyingparam.Price = "100"
		buyingparam.BuyerSig = sig[rand.Intn(14)]
		buyingparam.VoteStage = "1"
		buyingparam.SellerSig = sig[rand.Intn(14)]
		Buying = append(Buying, buyingparam)
	}
	//fmt.Println(Buying)
	datas, _ := json.Marshal(&Buying)
	fmt.Println("buying:   ", string(datas))

	//rand.Seed(time.Now().UnixNano())
	var Sell []SellParams
	for i := 0; i < 5; i++ {
		var buyingparam SellParams
		buyingparam.UserAddr = addr[rand.Intn(6)]
		buyingparam.ContractAddr = addr[rand.Intn(6)]
		buyingparam.TokenId = token[rand.Intn(16)]
		buyingparam.Price1 = "100"
		buyingparam.Price2 = "50"
		buyingparam.Day = "1"
		buyingparam.SellType = "FixPrice"
		buyingparam.PayChannel = "eth"
		buyingparam.VoteStage = "1"
		buyingparam.TradeSig = sig[rand.Intn(14)]
		buyingparam.Sig = sig[rand.Intn(14)]
		buyingparam.Hide = ""
		buyingparam.Currency = "eth"
		Sell = append(Sell, buyingparam)
	}
	//fmt.Println(Buying)
	datas, _ = json.Marshal(&Sell)
	fmt.Println("sell:   ", string(datas))

	var Cancel []CancelSellParams
	for i := 0; i < 5; i++ {
		var buyingparam CancelSellParams
		buyingparam.UserAddr = addr[rand.Intn(6)]
		buyingparam.ContractAddr = addr[rand.Intn(6)]
		buyingparam.TokenId = token[rand.Intn(16)]
		buyingparam.Sig = sig[rand.Intn(14)]
		Cancel = append(Cancel, buyingparam)
	}
	//fmt.Println(Buying)
	datas, _ = json.Marshal(&Cancel)
	fmt.Println("cancel:   ", string(datas))

	mar := "[{\"user_addr\":\"0xae407d12944877a969f38d5f2e313fe739271736\",\"contract_addr\":\"0xb8b2c20bfdace578e716b594a33367c3036a7f0f\",\"token_id\":\"1520322497495\",\"sig\":\"0xe719e460d8a34ff80ea882d143142fc788911174df1facd2600c37462691a1c5015adda7ba0c915ebe53f0e9009774274093fe11091e19757d69707fc93302c71c\"},{\"user_addr\":\"0xb8b2c20bfdace578e716b594a33367c3036a7f0f\",\"contract_addr\":\"0x9b1168d0ba701448d6b9ec344e986d4a2f16971e\",\"token_id\":\"1894851936255\",\"sig\":\"0x0edbaf0a57c8c92ab8e1c4129da3fb83f0adf8f2418b0e6cc61a7bfd2b25e41e10e92fd1b9e94c048fd6d4d306d7562cbb5d50a047b9b2cfde2b3c14b1799a971c\"},{\"user_addr\":\"0x80800580228ee590940fb75e86c8997611f079b3\",\"contract_addr\":\"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900\",\"token_id\":\"2091858617664\",\"sig\":\"0xe067bd64c0e812d2757670f898d646b340479d398c3226a18e5af0976396f55305d860e6e394b9d37f6309d86bcb5f4dda544da2ef0b347b869e5a9108c9a4891c\"},{\"user_addr\":\"0x80800580228ee590940fb75e86c8997611f079b3\",\"contract_addr\":\"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900\",\"token_id\":\"1894851936255\",\"sig\":\"0x6beb749a3b77e43849cd1726f96c9e45abbb4f3081c17f66a04c7853471637981c2f844910b00e55e74e74c5a2392b605542a3f5903325d94f74b22d8a2967bb1c\"},{\"user_addr\":\"0x9b1168d0ba701448d6b9ec344e986d4a2f16971e\",\"contract_addr\":\"0xb8b2c20bfdace578e716b594a33367c3036a7f0f\",\"token_id\":\"1520322497495\",\"sig\":\"0x48f880bb5a4e6c9cf6e8170a407c06081c1f431664530f697aafc6ad7e63d92409a785bc03bb2c0cb67fca1bbd4ff6b08e50d1aa5ac73d5b4c5cf74c17a57bd01c\"}]"
	var MarCancel []CancelSellParams
	err := json.Unmarshal([]byte(mar), &MarCancel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("mar:   ", MarCancel)

}

func TestGetIsVliaddr(t *testing.T) {
	fmt.Println(strings.ToLower("0x571CbB911fE99118B230585BA0cC7c5054324F85"))
	//tm := time.Unix(1661146660, 0)
	//tm := time.Unix(1663825060, 0).AddDate(0, 1, 0)
	tm := time.Unix(1662103355, 0)
	fmt.Println(tm)
	tm = time.Unix(1662016955, 0)
	fmt.Println(tm)
	f := tm.String()
	end := tm.Unix()
	fmt.Println("end=", end)
	fmt.Println(f)
	fmt.Println(f)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
	tm = tm.AddDate(0, 1, 0)
	fmt.Printf("%T", tm)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
	tn := time.Now().Unix()
	fmt.Println(tn)
	fmt.Println(time.Now().AddDate(0, 1, 0).Unix())

}
func TestQueryUserBidList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	userOffers, totalCount, gerr := nd.QueryUserBidList("0x9aa8fef730ebf39660cc444919e9f314645e155b", "0", "10")
	//data, gerr := nd.GetNftMarketInfo()
	if gerr != nil {
		fmt.Println(gerr)
	}
	fmt.Println(userOffers, totalCount)
}

type HttpRequestFilter struct {
	Match      string         `json:"match"`
	Filter     []StQueryField `json:"filter"`
	Sort       []StSortField  `json:"sort"`
	Nfttype    string         `json:"nfttype"`
	StartIndex string         `json:"start_index"`
	Count      string         `json:"count"`
}

func TestQuerNFTList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	//serr := nd.db.Model(&SysInfos{}).Select("nfttotal").Last(&nftCount)
	//if serr.Error != nil {
	//	fmt.Println(serr)
	//}
	var data HttpRequestFilter
	data.Filter = []StQueryField{}
	data.Filter = append(data.Filter, StQueryField{
		Field:     "collectcreator",
		Operation: "=",
		Value:     "0x400ed949861be04a4e0a6f0d1464fc61d89cc4f2",
	})
	data.Filter = append(data.Filter, StQueryField{
		Field:     "collections",
		Operation: "=",
		Value:     "democollect",
	})
	data.Filter = []StQueryField{}

	data.Match = ""
	data.StartIndex = "0"
	data.Nfttype = "nft"
	data.Count = "50"
	data.Sort = []StSortField{}

	collectData := Collects{}
	result := nd.db.Model(&Collects{}).Select([]string{"categories", "createaddr", "desc", "name", "contract", "contracttype", "img", "totalcount"}).Where("createaddr = ? and name = ?", "0xbe8c75133a7e4f29b7cdc15d4a45f7593a4f8898", "测试10").
		First(&collectData)
	fmt.Println(result)

	datah, gerr := nd.QueryHomePage(false)
	//data, gerr := nd.GetNftMarketInfo()
	if gerr != nil {
		fmt.Println(gerr)
	}
	fmt.Println(datah)
	userOffers, totalCount, gerr := nd.QueryNftByFilterNftSnft(data.Filter, data.Sort, data.Nfttype, data.StartIndex, data.Count)
	//data, gerr := nd.GetNftMarketInfo()
	if gerr != nil {
		fmt.Println(gerr)
	}
	fmt.Println(userOffers, totalCount)
}

func TestPendingNFT(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	userOffers, totalCount, gerr := nd.QueryUnverifiedNfts("0", "2", "")
	//data, gerr := nd.GetNftMarketInfo()
	if gerr != nil {
		fmt.Println(gerr)
	}
	fmt.Println(userOffers, totalCount)

	s1, s2, gerr := nd.QueryPendingKYCList("0", "2", "")
	//data, gerr := nd.GetNftMarketInfo()
	if gerr != nil {
		fmt.Println(gerr)
	}
	fmt.Println(s1, s2)
}

func TestSetSuccessETH(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	type Nftdel struct {
		Creat string `json:"creat"`
		Total string `json:"total"`
	}
	setnfts := Nfts{}

	nerr := nd.db.Model(&Nfts{}).Where("tokenid = ? ", 0).First(&setnfts)
	if nerr.Error != nil {
		log.Println("SetNft() err = Nft not exist.")
	}
	nn := []Nftdel{}
	nftsql := `select creat,count(*) as total from (select DATE_FORMAT(created_at,"%Y-%m-%d") as creat 
	from nfts where deleted_at is null ) as ss GROUP BY creat`
	//var total int64
	nerr = nd.db.Raw(nftsql).Scan(&nn)
	if nerr.Error != nil {
		fmt.Println(nerr)
	}
	fmt.Println(nn)

	nfts := Nfts{}
	nerr = nd.db.Model(&Nfts{}).Last(&nfts)
	if nerr.Error != nil {
		fmt.Println(nerr)
	}
	fmt.Println(nfts.Meta)
	slen := strings.LastIndex(nfts.Meta, "/")
	fmt.Println(slen)

	fmt.Println(nfts.Meta)

	fmt.Println(nfts.Meta[:strings.LastIndex(nfts.Meta, "/")])
}

func TestSaveDirtoipfs(t *testing.T) {
	//nfts := []Nfts{}
	//err := json.Unmarshal([]byte("[{\"contract\":\"1\",\"tokenid\":\"\"}]"), &nfts)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//HomePageCatchs.NftLoop = nfts
	//fmt.Println(nfts)

	//imagexint := 220
	//imgxint, err := strconv.Atoi("222")
	//if err != nil {
	//	fmt.Println("imgxint transfer int err =", err)
	//
	//}
	//if imagexint <= imgxint+10 && imagexint >= imgxint-10 {
	//
	//	fmt.Println("captcha ok")
	//
	//} else {
	//	fmt.Println("captcha auth error")
	//
	//}
	//fmt.Println(imgxint)
	//fmt.Println(ImageDir)
	////CaptchaDefault()
	//data, err := GetNftInfoFromIPFSWithShell("/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/mask.png")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(data)
	//url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	url := "http://192.168.1.235:9006"
	//url := "https://www.wormholestest.com"
	//url := "http://43.129.181.130:8561"
	//url := "http://192.168.1.237:9006"
	//url := "https://snft.wormholestest.com"

	spendT := time.Now()
	s := shell.NewShell(url)
	s.SetTimeout(5 * time.Second)
	mhash, err := s.AddDir("D:\\workdir\\go\\code\\captcha\\captcha\\bg\\worm")
	//mhash, err := s.AddDir("D:\\picture\\2022_6_3")
	if err != nil {
		fmt.Println("SaveToIpfs() err=", err)

	}
	fmt.Printf("SaveToIpfs  Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	fmt.Println(mhash)
	fmt.Println(mhash)

	//v, err := os.Open("D:\\workdir\\go\\code\\captcha\\captcha\\bg\\1.jpeg")
	//if err != nil {
	//	fmt.Printf("Http get [%v] failed! %v", v, err)
	//	return
	//}
	//content, err := ioutil.ReadAll(v)
	//if err != nil {
	//	fmt.Printf("Read http response failed! %v", err)
	//	return
	//}
	//collectImageUrl, serr := SaveToIpfss(string(content))
	//if serr != nil {
	//	fmt.Println("SaveToIpfs() save image err=", serr)
	//	return
	//}
	//fmt.Println(collectImageUrl)

}

func TestSaveToIpfss(t *testing.T) {
	url := "https://www.wormholestest.com"
	//url = "https://www.wormholestest.com/c0x9ac8846ec59116e2e63b54c81670ff15f1d00f1a/#/"
	//url = "https://www.wormholestest.com:443"
	sjson := "[{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/1.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/2.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/3.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/4.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/5.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/6.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/7.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/8.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/9.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/10.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/11.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/12.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/13.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/14.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/15.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/16.jpeg\"}]"
	//url = "http://192.168.1.235:9006"

	sjson = "[{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/1.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/2.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/3.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/4.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/5.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/6.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/7.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/8.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/9.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/10.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/11.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/12.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/13.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/14.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/15.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/16.jpeg\"}]"
	sjson = "[{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/1.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/10.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/11.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/12.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/13.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/14.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/15.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/16.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/2.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/3.png\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/4.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/5.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/6.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/7.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/8.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/9.jpeg\"}]"
	//sjson = "[{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/1.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/10.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/11.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/12.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/13.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/14.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/15.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/16.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/2.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/3.png\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/4.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/5.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/6.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/7.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/8.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmSgeszvZ278adFRvWxfsKHbvHwXdYVecyQDt3BbeT2NS2/9.jpeg\"},{\"url\":\"http://192.168.1.235:9006/ipfs/QmRbPUov8yYA2H4qSq1h2nsrCQiqamZBqUByevPtj4AQw4/mask.png\"}]"
	sjson = "[{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/1.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/10.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/11.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/12.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/13.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/14.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/15.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/16.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/2.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/3.png\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/4.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/5.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/6.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/7.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/8.jpeg\"},{\"url\":\"https://www.wormholestest.com/ipfs/QmcJQiyCuUtqboNX18ymt1UAEvihRz3DF9QGoFGf7XjovF/9.jpeg\"}]"
	spendT := time.Now()
	s := shell.NewShell(url)
	s.SetTimeout(500 * time.Second)
	mhash, err := s.Add(bytes.NewBufferString(sjson))
	if err != nil {
		log.Println("SaveToIpfs() err=", err)
		//return "", err
	}
	fmt.Printf("SaveToIpfs  Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	fmt.Println(mhash)

	//return mhash, nil
}

func TestNewRedisCatch(t *testing.T) {
	err := NewQueryCatch("192.168.1.235:6379", "user123456")
	fmt.Println(err)
	qCatch := GetRedisCatch()
	err = qCatch.HsetRedisData("user1", "name", []byte("name1"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	//err = qCatch.HdelRedisData("user1", "name")
	name, err := qCatch.HgetRedisData("user1", "name")
	if err != nil {
		t.Fatalf(err.Error())
	}
	qCatch.SetDirtyFlag([]string{"user1"})
	err = qCatch.HsetRedisData("user1", "age", []byte("20"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = qCatch.HsetRedisData("user2", "name", []byte("name2"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	qCatch.SetDirtyFlag([]string{"user1"})
	err = qCatch.HsetRedisData("user2", "age", []byte("50"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	name, err = qCatch.HgetRedisData("user1", "name")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(name))
	age, err := qCatch.HgetRedisData("user1", "age")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(age))
	name, err = qCatch.HgetRedisData("user2", "name")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(name))
	age, err = qCatch.HgetRedisData("user2", "age")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(age))
	//err = qCatch.SetRedisData("dskfjdkf", []byte("12345678"))
	//nn, err := qCatch.GetRedisData("dskfjdkf")
	//fmt.Println(nn)
	//err = qCatch.SetRedisData("dskfjdkf", []byte("555555"))
	//nn, err = qCatch.GetRedisData("dskfjdkf")
	//fmt.Println(nn)
	type test struct {
		Name string
		Age  int
	}
	ExchangeOwer = "0x01842a2cf56400a245a56955dc407c2c4137321e"
	err = qCatch.CatchQueryData("1", "2", test{"1", 2})
	err = qCatch.CatchQueryData("1", "2", test{"2", 200})
	err = qCatch.CatchQueryData("1", "1", test{"1", 1})
	err = qCatch.CatchQueryData("1", "3", test{"1", 1})
	err = qCatch.CatchQueryData("1", "5", test{"1", 1})
	err = qCatch.CatchQueryData("2", "3", test{"2", 3})
	err = qCatch.CatchQueryData("2", "1", test{"2", 2})
	err = qCatch.CatchQueryData("5", "2", test{"5", 2})
	//var tt = test{name: "name", age: 1,}
	tt := test{Name: "name", Age: 1}
	//qCatch.SetDirtyFlag([]string{"1", "2"})
	nftCatch := NftFilter{}

	err = qCatch.GetCatchData("QueryNftByFilterNftSnft", "050", &nftCatch)
	err = qCatch.GetCatchData("1", "2", &tt)
	err = qCatch.GetCatchData("1", "2", &tt)
	err = qCatch.GetCatchData("1", "2", &tt)
	err = qCatch.GetCatchData("1", "2", &tt)
	err = qCatch.GetCatchData("2", "1", &tt)
	err = qCatch.GetCatchData("2", "3", &tt)
	err = qCatch.GetCatchData("5", "2", &tt)
	//qCatch.SetDirtyFlag([]string{"1", "2"})
}

func TestDelCatch(t *testing.T) {
	err := NewQueryCatch("192.168.1.235:6379", "user123456")
	fmt.Println(err)
	qCatch := GetRedisCatch()
	ExchangeOwer = "0x671a9f50d3f1a1aed7310ebb67cc7fe810a06998"
	var paraminfo SysParamsInfo
	cerr := qCatch.GetCatchData("QuerySysParams", "QuerySysParams", &paraminfo)
	if cerr != nil {
		log.Printf("QueryPendingKYCList() default  time.now=%s\n", time.Now())
		return
	}
	//GetRedisCatch().SetDirtyFlag(SysParamsDirtyName)
	cerr = qCatch.GetCatchData("QuerySysParams", "QuerySysParams", &paraminfo)
	if cerr == nil {
		log.Printf("QueryPendingKYCList() default  time.now=%s\n", time.Now())
		return
	}
}

func TestAddDirtyQuery(t *testing.T) {
	err := NewQueryCatch("192.168.56.128:6379", "")
	err = NewQueryMainCatch("192.168.1.235:6379", "user123456")
	fmt.Println(err)
	ExchangeOwer = "0x671a9f50d3f1a1aed7310ebb67cc7fe810a06998"
	ExchangeOwer = "0x57ed0c503c40308e802414405ce3d399fe3a42c6"
	mCatch := GetRedisMainCatch()
	qCatch := GetRedisCatch()
	mDirtyQuerys, _ := mCatch.GetDirtyQuerys()
	fmt.Println(mDirtyQuerys)
	DirtyQuerys, _ := qCatch.GetDirtyQuerys()
	fmt.Println(DirtyQuerys)

	cerr := mCatch.SaveDirtyQuerys([]string{"QuerySysParams", "QuerySysParams0"})
	if cerr != nil {
		log.Printf("TestAddDirtyQuery() default  time.now=%s\n", time.Now())
		return
	}
	cerr = mCatch.SaveDirtyQuerys([]string{"QuerySysParams0", "QuerySysParams1"})
	if cerr != nil {
		log.Printf("TestAddDirtyQuery() default  time.now=%s\n", time.Now())
		return
	}
	querydirtys, _ := mCatch.GetDirtyQuerys()
	fmt.Println(querydirtys)
	qCatch.scanDirtyQuery(mCatch)
	cerr = mCatch.SaveDirtyQuerys([]string{"QuerySysParams0", "QuerySysParams1"})
	if cerr != nil {
		log.Printf("TestAddDirtyQuery() default  time.now=%s\n", time.Now())
		return
	}
	qCatch.scanDirtyQuery(mCatch)
}

func TestSig_recover(t *testing.T) {

	sigdata := "{\"user_addr\":\"0x36a6bfa3f64a6cd668645a979edc6436257f9ea2\",\"buyer_Addr\":\"0x36a6bfa3f64a6cd668645a979edc6436257f9ea2\",\"nft_contract_addr\":\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\",\"nft_token_id\":\"3556780204294\",\"price\":\"1000000000\",\"buyer_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x67dd81\\\",\\\"sig\\\":\\\"0xf6fd493cfbff6639f9f0855752a4af09bf35bc4e0fab7674c041ab2161c5d3c27540d63ac46885754bba808637e16583d1fc8136d16a38ff62b20c9e8f8eb70c1b\\\"}\",\"seller_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"royalty\\\":\\\"0x64\\\",\\\"meta_url\\\":\\\"7b226d657461223a222f697066732f516d6538586b6a665155706167755a64624e577035644368366d704c6b66586152625a3637437367796769553151222c22746f6b656e5f6964223a2233353536373830323034323934227d\\\",\\\"exclusive_flag\\\":\\\"1\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x6748af\\\",\\\"sig\\\":\\\"0x2f1680c2a35527a1daf33d2c18af7e9fbe892ea1756d0535b2d7ed51b5457c29382e69f234e69aea08a2c78ce78cfb5d24fce15399be4c1f10a3230dedfc21ce1c\\\"}\",\"vote_stage\":\"\"}"
	sigdata = "{\"user_addr\":\"0xf2edd6638a965e1b540d9c081df092e67dd36026\",\"nft_contract_addr\":\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\",\"nft_token_id\":\"3556780204294\",\"pay_channel\":\"ERB\",\"currency_type\":\"ERB\",\"price\":\"1000000000\",\"trade_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x8ae25\\\",\\\"sig\\\":\\\"0xa1e2591ee7c6e182019429d7b3e564753174e95dc9a7d3f3343c6cbc9aeeee6d6eb5b2cb9f5d445bfb56e8acfb4eeae43ddde649f1e95ec0f0322659206efe4c1b\\\"}\",\"dead_time\":\"1658566687\",\"vote_stage\":\"\"}"
	sigdata = "{\"user_addr\":\"0xf2edd6638a965e1b540d9c081df092e67dd36026\",\"buyer_Addr\":\"0xf2edd6638a965e1b540d9c081df092e67dd36026\",\"nft_contract_addr\":\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\",\"nft_token_id\":\"3556780204294\",\"price\":\"1000000000\",\"buyer_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x67de77\\\",\\\"sig\\\":\\\"0x86b2ea329ee83526cdd163c5a60f78573b09b9c9be44ae7d654c4ee447d1b8be6d963df09888a61c58188346b148fc9f6e3cb63a67470f012d5035496a5a12601c\\\"}\",\"seller_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"royalty\\\":\\\"0x64\\\",\\\"meta_url\\\":\\\"7b226d657461223a222f697066732f516d6538586b6a665155706167755a64624e577035644368366d704c6b66586152625a3637437367796769553151222c22746f6b656e5f6964223a2233353536373830323034323934227d\\\",\\\"exclusive_flag\\\":\\\"1\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x6748af\\\",\\\"sig\\\":\\\"0x2f1680c2a35527a1daf33d2c18af7e9fbe892ea1756d0535b2d7ed51b5457c29382e69f234e69aea08a2c78ce78cfb5d24fce15399be4c1f10a3230dedfc21ce1c\\\"}\",\"vote_stage\":\"\"}"
	sigdata = `{"user_addr":"0xf2edd6638a965e1b540d9c081df092e67dd36026","buyer_Addr":"0xf2edd6638a965e1b540d9c081df092e67dd36026","nft_contract_addr":"0xbacd84846337bf50532cd3a32fe48bbeded102b3","nft_token_id":"3556780204294","price":"1000000000","buyer_sig":"{\"price\":\"0xde0b6b3a7640000\",\"exchanger\":\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\",\"block_number\":\"0x67de77\",\"sig\":\"0x86b2ea329ee83526cdd163c5a60f78573b09b9c9be44ae7d654c4ee447d1b8be6d963df09888a61c58188346b148fc9f6e3cb63a67470f012d5035496a5a12601c\"}","seller_sig":"{\"price\":\"0xde0b6b3a7640000\",\"royalty\":\"0x64\",\"meta_url\":\"7b226d657461223a222f697066732f516d6538586b6a665155706167755a64624e577035644368366d704c6b66586152625a3637437367796769553151222c22746f6b656e5f6964223a2233353536373830323034323934227d\",\"exclusive_flag\":\"1\",\"exchanger\":\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\",\"block_number\":\"0x6748af\",\"sig\":\"0x2f1680c2a35527a1daf33d2c18af7e9fbe892ea1756d0535b2d7ed51b5457c29382e69f234e69aea08a2c78ce78cfb5d24fce15399be4c1f10a3230dedfc21ce1c\"}","vote_stage":""}`
	sigdata = "{\"user_addr\":\"0x389a82b54d1c8d3c803b56973fe75f20865eb8dc\",\"buyer_Addr\":\"0xf2edd6638a965e1b540d9c081df092e67dd36026\",\"nft_contract_addr\":\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\",\"nft_token_id\":\"3556780204294\",\"price\":\"1000000000\",\"buyer_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x67de77\\\",\\\"sig\\\":\\\"0x86b2ea329ee83526cdd163c5a60f78573b09b9c9be44ae7d654c4ee447d1b8be6d963df09888a61c58188346b148fc9f6e3cb63a67470f012d5035496a5a12601c\\\"}\",\"seller_sig\":\"{\\\"price\\\":\\\"0xde0b6b3a7640000\\\",\\\"royalty\\\":\\\"0x64\\\",\\\"meta_url\\\":\\\"7b226d657461223a222f697066732f516d6538586b6a665155706167755a64624e577035644368366d704c6b66586152625a3637437367796769553151222c22746f6b656e5f6964223a2233353536373830323034323934227d\\\",\\\"exclusive_flag\\\":\\\"1\\\",\\\"exchanger\\\":\\\"0xbacd84846337bf50532cd3a32fe48bbeded102b3\\\",\\\"block_number\\\":\\\"0x6748af\\\",\\\"sig\\\":\\\"0x2f1680c2a35527a1daf33d2c18af7e9fbe892ea1756d0535b2d7ed51b5457c29382e69f234e69aea08a2c78ce78cfb5d24fce15399be4c1f10a3230dedfc21ce1c\\\"}\",\"vote_stage\":\"\"}"
	sig := "0xfc143e7844c53170b003cd2edd2f55f149ec5bd22ab09dbab7288b0c1c0c7aca0aed47e5c06a0ce0a15efb3c50a3f28c98fa5f67dbcef210e2ba97b8bd6c03b91c"
	sig = "0xaa1135121983e991f38fb01f7b9efabe0c942e33a3384578cb8b3f2c994ea9a43b30e7c29bcf2c765317ad7ed5873dcb2769c560d382abd1f101fe5cc61c6db01c"
	sig = "0x317fb3d50c8a7bb9895fa1279439c043879ec8325ff78b1fd9c46514cac31cd5522cc0c2078821eb9cca62276d57bbfba18dd225f550b236b9726954cdbec9a31b"
	sig = "0xa7949dd6c3bea22ff18030b5b8b07ed50f1d84b7d2956d2f9906c451a9efe5446040b36be59d02fd6668924564e30a1af100541a4b0388c6b7e252c694cbbce91c"
	sig = "0x175c54e06d44d32140d67a4d93310d814845197e2c627f0a136d2c758d4fa32f371e832f13ef54e39049be6d3950fc708b7c2076fa04c31535ae89aa87c0ae001b"
	sig = "0xcc4fe96e9155fa84cf86b1e34c88096955b6d5e4a1fe663565926d64502e59d5447567b65624c64751cb1518d0ace6ceeb6f2d6d658e1b04bdea577435d825b31c"
	sig = "0x5981a44574412756ac8ffec8a24a712060d1a97f44c277b4aede8e5cc2c678722e08225df725608f1347639f783352e166a85d568cc2a3461143a5c75d66bcfa1c"
	addr := strings.ToLower("0xf2edd6638a965e1b540d9c081df092e67dd36026")
	verificationAddr, err := GetEthAddr(sigdata, sig)
	if err != nil {
	}

	verificationAddrS := verificationAddr.String()
	verificationAddrS = strings.ToLower(verificationAddrS)

	addr = strings.ToLower(addr)
	fmt.Printf("sigdebug verificationAddrS = [%s], approveAddr's addr = [%s]\n", verificationAddrS, addr)
	if verificationAddrS == addr {
		fmt.Println("sigdebug verify [Y]")
	}
	fmt.Println("sigdebug verify [N]")
	//return true, nil

}

func TestNftDb_UserKYC(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	uerr := nd.UserKYC("0x9aa8fef730ebf39660cc444919e9f314645e155b", "0x400ed949861be04a4e0a6f0d1464fc61d89cc4f2,0xf13ea37c72860a79521526b5c09584b0cab06db5",
		"1", "Passed", "")
	uerr = nd.VerifyNft("0x9aa8fef730ebf39660cc444919e9f314645e155b",
		"5736644112628,6494497382315",
		"1", "Passed")
	if uerr != nil {
		fmt.Println(uerr)
	}
	fmt.Println("ok")
}

func TestNftDb_Exchange(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	//uerr := nd.DelPartnerLogo("6")
	//data, count, uerr := nd.QueryExpireNft("0", "10", "30")
	snftdata, count, uerr := nd.SNftCollectionSearch("Art", "de", "0", "30")
	if uerr != nil {
		fmt.Println(uerr)
	}
	fmt.Println(snftdata)
	fmt.Println(count)
	if uerr != nil {
		fmt.Println(uerr)
	}
	fmt.Println("ok")
}

func TestNftDb_ExpiredNft(t *testing.T) {

	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	//uerr := nd.DelPartnerLogo("6")
	//data, count, uerr := nd.QueryExpireNft("0", "10", "30")
	snftdata, count, uerr := nd.QueryExpireNft("0", "1000", "0")
	if uerr != nil {
		fmt.Println(uerr)
	}
	//uerr = nd.DelExpiredNft("0")
	fmt.Println(snftdata)
	fmt.Println(count)
	if uerr != nil {
		fmt.Println(uerr)
	}
	fmt.Println("ok")
}

func TestSyncChain(t *testing.T) {
	contracts.EthNode = "http://43.129.181.130:8561"
	NftScanServer = "http://192.168.1.237:8089"
	//NewQueryCatch("192.168.56.120:6379", "")
	NewQueryCatch("192.168.1.237:6379", "user123456")
	TransferSNFT = true
	RoyaltyLimit = 10000
	SyncBlock(sqldsnT)
}

func TestBigData(t *testing.T) {
	b := big.NewInt(10000000)
	g := big.NewInt(1000000000000000000)
	b = b.Mul(b, g)
	fmt.Println(b.String())
	fmt.Println(b.Bytes())
	h := hexutil.EncodeBig(b)
	fmt.Println(h)
	b, _ = b.SetString(h[2:], 16)
	fmt.Println(b.String())
}
func TestBigDatass(t *testing.T) {
	adminKey, aderr := crypto.GenerateKey()
	if aderr != nil {
		log.Println("InitSysParams create Agent adminkey err=", aderr)

	}
	addr := crypto.PubkeyToAddress(adminKey.PublicKey)
	fmt.Println(addr)
	addprv := crypto.FromECDSA(adminKey)
	fmt.Println(hexutil.Encode(addprv))
	//"0xbdc996a45840937b92c3bea53ab0d76664e4ecc3437958f68e86c6b99757bf78
}

func TestTestmode(t *testing.T) {
	m := 256 % 16
	fmt.Println(m)
	m = 253 % 16
	fmt.Println(m)
	m = 25 % 16
	fmt.Println(m)
	m = 2 % 16
	fmt.Println(m)
	m = 11 % 16
	fmt.Println(m)
	m = 16 % 16
	fmt.Println(m)
}

func TestHomePageRenew(t *testing.T) {
	//nd, err := NewNftDb(sqldsnT)
	//if err != nil {
	//	fmt.Printf("connect database err = %s\n", err)
	//}
	//defer nd.Close()
	//for i := 0; i < 100 && i < 10; {
	//	fmt.Println(i)
	//	i++
	//}
	rand.Seed(time.Now().UnixNano())
	scaned := make(map[int64]bool)
	for i := 0; i < 5 && int64(i) < 1; {
		index := rand.Int63()%1 + 1
		log.Println("HomePageRenew() rand.Int63() index= ", index)
		/*if index == 0 {
			index = 1
		}*/
		flag := scaned[index]
		if flag {
			//i = i - 1
			log.Println("HomePageRenew() scaned[index] index= ", index)
			time.Sleep(time.Second)
			continue
		}
		scaned[index] = true
		if i != 2 {
			continue
		}
		i++
	}
	//err := HomePageRenew()
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func TestAlice(t *testing.T) {
	txs := []string{"1", "2", "3"}
	tx := []string{}
	for i, j := range txs {
		fmt.Println(j, i)
		if j == "3" {
		} else {
			tx = append(tx, txs[i])

		}
		//txs = append(txs[:i], txs[i+1:]...)
	}
	fmt.Println(tx)
}

func TestRecover(t *testing.T) {
	sig := "{\"price\":\"0x6046f37e5945c0000\",\"nft_address\":\"0x0000000000000000000000000000000000000001\",\"exchanger\":\"0xd4947f643ef5f1420519e6c6b71d443914efa743\",\"block_number\":\"0x8b88\",\"sig\":\"0x753bb26209021610ad003b3f42656a056c94b4e5e0ed70340c760f9282f6f1e94372a6be920082d7176d44dd20e23ed4693331125a6193de2e5430dffa04b77f1c\"}"
	sig = "{\"price\":\"0x6046f37e5945c0000\",\"nft_address\":\"0x0000000000000000000000000000000000000001\",\"exchanger\":\"0xd4947f643ef5f1420519e6c6b71d443914efa743\",\"block_number\":\"0x8b88\",\"sig\":\"0x753bb26209021610ad003b3f42656a056c94b4e5e0ed70340c760f9282f6f1e94372a6be920082d7176d44dd20e23ed4693331125a6193de2e5430dffa04b77f1c\"}"
	sig = "{\"price\":\"0x56bc75e2d63100000\",\"nft_address\":\"0x0000000000000000000000000000000000000001\",\"exchanger\":\"0xd4947f643ef5f1420519e6c6b71d443914efa743\",\"block_number\":\"0x603fb9\",\"seller\":\"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900\",\"sig\":\"0x694c643f462931444abc52752b24af61f58cfdc37ab891632734a384b8108f6f69093339d0760313f88ade72e22a5e630c7c767a316ad02c1abd6d865684a0fd1c\"}"
	buyer := contracts.Buyer{}
	err := json.Unmarshal([]byte(sig), &buyer)
	if err != nil {
		log.Println("SigVerify Unmarshal err=", err)
	}
	msg := buyer.Price + buyer.Nftaddress + buyer.Exchanger + buyer.Blocknumber + buyer.Seller
	msg = "0x8ac7230489e800000x647b226d657461223a222f697066732f516d5174556351744168423941533169427a6a7354675742343967776d6655323743344133724a387a5132556b71222c22746f6b656e5f6964223a2238343236373938323735333537227d10xb146bf8b0a069630ac5d151c5bf87fad77a479f70x12016307200100"
	sig = "0x7417a8ae80c7dd4e61077cfd6ee9c007a2c703e18308de64544fd3e174a6e9922ed2cbfed5e28fe044a907c68f18cb3f67f1ea96fe8a28c4c39cda878b37744b1c"
	sig = "0x382d215867c00889117a16b22cb37c90a817264491e48eb479b1fa098577f76152b2df9f1042d690e17f1d3fbd37297a882cad4148d29ae58d6827d382d1aa4e1b"
	msg = "0xde0b6b3a76400000x647b226d657461223a222f697066732f516d5757326f447a4d486638335664704e6d426450524a754a3151536a54596a5651315350365358744b326f3842222c22746f6b656e5f6964223a2231363731333633333535393439227d10x01842a2cf56400a245a56955dc407c2c4137321e0x64f626"
	sig = "0x3a1055b1c1a0194e2b5ab08b6c5e8caf0a54edfd4456854089f58475aaf89a936e0455edabe18f99b09eeaa7f7bac1b319d1bc82529003f92a3d7bc0defb65ca1c"
	msg = "user"
	toaddr, rerr := contracts.RecoverAddress(msg, sig)
	fmt.Println("toaddr =", toaddr.String(), "  buyaddr =")
	if rerr != nil {
		log.Println("SigVerify() recoverAddress() err=", err)
	}

}

func TestDelFavor(t *testing.T) {
	nd, nerr := NewNftDb(sqldsnT)
	if nerr != nil {
		fmt.Printf("connect database err = %s\n", nerr)
	}
	defer nd.Close()
	var nftdata []string
	err := nd.db.Model(&Nfts{}).Select("tokenid").Where("collectcreator =? and collections=? and mintstate =? ", "0x7cfbfeab24f646ab96b31da34ac689bc71faf655", "wormholes", "NoMinted").Find(&nftdata)
	if err.Error != nil {
		fmt.Println("DelCollection() delete collection under nfts err= ", err.Error)
	}
	fmt.Println(nftdata)

	var favorited []NftFavorited
	err = nd.db.Model(&NftFavorited{}).Where("tokenid in ? ", nftdata).Find(&favorited)
	if err.Error != nil {
		fmt.Println("DelCollection() delete collection under nfts err= ", err.Error)
	}
	fmt.Println(favorited)

}
