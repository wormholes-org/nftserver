package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/ethhelper"
	"github.com/nftexchange/nftserver/ethhelper/database"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
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

//const sqlsvrLcT = "admin:user123456@tcp(192.168.1.235:3306)/"
const sqlsvrLcT = "admin:user123456@tcp(192.168.1.235:3307)/"


//const sqlsvrLcT = "admin:user123456@tcp(192.168.1.237:3306)/"

//const sqlsvrLcT = "admin:user123456@tcp(192.168.56.128:3306)/"
//
//const sqlsvrLcT = "demo:123456@tcp(192.168.56.129:3306)/"

//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
//var SqlSvrT = "admin:user123456@tcp(192.168.1.238:3306)/"
//const dbNameT = "mynftdb"
//const dbNameT = "tnftdb"
//const dbNameT = "nftdb"
//const dbNameT = "snftdb8012"
//const dbNameT = "tttt"
const dbNameT = "c0xbc8ac1fe086809fdaab2568dd3e8025218a62bb5"

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

func TestGetAdminAddr(t *testing.T) {
	addrs, err := ethhelper.AdminList()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addrs)
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
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	snftChip, count, err := nd.QuerySnftChip("0x01842a2cf56400a245a56955dc407c2c4137321e", "7679889549168", "0", "10")
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

func TestQueryNftByFilterNew(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

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
		{
			"collectcreator",
			"=",
			"0x7a149f02e5e4571c42d5cf69b4ccb5772fa1b275",
		},
		{
			"collections",
			"=",
			"collect_test_0",
		},
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
		/*{
			"createdate",
			">=",
			"1654226774",
		},*/
	}
	//sorts := []StSortField{{By: "createdate", Order: "desc"}}
	//nfilters = []StQueryField{}
	sorts := []StSortField{
		/*{
			"sellprice",
			"asc",
		},*/
		{
			"verifiedtime",
			"asc",
		},
	}
	nftByFilter, count, err := nd.QueryNftByFilterNftSnft(nfilters, sorts, "nft", "0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	nftByFilter, count, err = nd.QueryNftByFilterNew(nfilters, sorts, "nftsnft", "5", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	t.Logf("nft = %v %v\n", nftByFilter, count)
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

func TestBalanceOfWeth(t *testing.T) {
	c, err := ethhelper.BalanceOfWeth("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c")
	fmt.Println(c, err)
	fmt.Println(c > 1000, err)
}

func TestOwnOf(t *testing.T) {
	ct, err := ethhelper.IsErc721("0xa1e67a33e090afe696d7317e05c506d7687bb2e5")
	if err != nil {
		fmt.Println(err)
	}
	if ct == 1 {
		isNft721, err := ethhelper.IsOwnerOfNFT721("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "9767110817076")
		if err != nil {
			fmt.Println(err)
		}
		if isNft721 {
			fmt.Println("isNft721")
		}
		approve, err := ethhelper.IsApprovedNFT721("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "9767110817076")
		if err != nil {
			fmt.Println(err)
		}
		if approve {
			fmt.Println("Nft721 is approve")
		}
	} else {
		sNft1155, err := ethhelper.IsOwnerOfNFT1155("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "9767110817076")
		if err != nil {
			fmt.Println(err)
		}
		if sNft1155 {
			fmt.Println("isNft721")
		}
		approve, err := ethhelper.IsApprovedNFT1155("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5")
		if err != nil {
			fmt.Println(err)
		}
		if approve {
			fmt.Println("Nft1155 is approve")
		}
	}
}

func TestAllowanceOfWeth(t *testing.T) {
	c, err := ethhelper.AllowanceOfWeth("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c")
	fmt.Println(c, err)
	c = c[:len(c)-9]
	fmt.Println(c, err)
	wei := new(big.Int)
	wei.SetString(c, 10)
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
		"sigedata",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	err = nd.ModifyCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"test", "img", "contract_type", "contract_addr",
		"test desc.", "art", "sig string")
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
		"sigedata",
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

func TestSyncProc(t *testing.T) {
	//9532550
	InitSyncBlockTs(sqldsnT)
	//syncFlag := make(chan struct{})
	//SyncProc("", syncFlag)
	//<-syncFlag
}

func TestSyncNftFromChain(t *testing.T) {
	buyResultCh := make(chan []*database.NftTx)
	wethTransferCh := make(chan *ethhelper.WethTransfer)
	wethApproveCh := make(chan *ethhelper.WethTransfer)
	var BlockTxs []*database.NftTx
	var wethTransfers []*ethhelper.WethTransfer
	var wethApproves []*ethhelper.WethTransfer
	endCh := make(chan bool)
	go ethhelper.SyncNftFromChain(strconv.Itoa(9651405 /*9570987*/), true, buyResultCh, wethTransferCh, wethApproveCh, endCh)
	isOver := false
	for {
		select {
		case buyResult := <-buyResultCh:
			fmt.Println(buyResult)
			BlockTxs = append(BlockTxs, buyResult...)
		case wethTransfer := <-wethTransferCh:
			fmt.Println(wethTransfers)
			wethTransfers = append(wethTransfers, wethTransfer)
		case wethApprove := <-wethApproveCh:
			fmt.Println(wethTransfers)
			wethApproves = append(wethApproves, wethApprove)
		case <-endCh:
			isOver = true
			break
		default:
		}
		if isOver {
			break
		}
	}
	fmt.Println("end")
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

func TestGetBlockTxs(t *testing.T) {
	//txs := GetBlockTxs(9508909)
	txs, _, _ := GetBlockTxs(9508910)
	fmt.Println(len(txs))
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

func TestVerifySign(t *testing.T) {
	sig := "0x29381026df3b9cb57d67eaa620c4b4ace3886b62e344586aaef09adaf484941d35e1b824be66e0ea10fbf6dd63d6a9822fbb6c855a8f9bd3922c134d3a9a0f871b"
	msg := `{"def_language":"en_us"}`
	_, err := IsValidAddr(msg, sig, "")
	if err != nil {
		t.Fatal("err=", err)
	}

}
func TestNftDb_DefaultCountry(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	//approve :="{\"msg\": \"0xBf5Fabb29D464B41Eaf88096654dd813Ad7bcF580x8488277e4221ccd61a854c6bc9de027a3dd7fafa0x989680\",\"sig\":\"0x89a7c2e30d7bc84f8a86e78c712c235febc16bdc1e03a949cbd042d9d9e6835619ef99ce17bcc86c933c955408b434e1821a3941ced3fd714530ca43528102a81c\"}"
	approve := "{\"exchanger_owner\":\"0xBf5Fabb29D464B41Eaf88096654dd813Ad7bcF58\",\"to\":\"0x8488277e4221ccd61a854c6bc9de027a3dd7fafa\",\"block_number\":\"0x989680\",\"sig\":\"0x89a7c2e30d7bc84f8a86e78c712c235febc16bdc1e03a949cbd042d9d9e6835619ef99ce17bcc86c933c955408b434e1821a3941ced3fd714530ca43528102a81c\"}"
	approve = strings.ToLower(approve)
	fmt.Println(approve)
	if err != nil {
		fmt.Println(err)
	}

	var Exchangerauth map[string]string
	authSign := "{\"block_number\":\"0x989680\",\"exchanger_owner\":\"0xbf5fabb29d464b41eaf88096654dd813ad7bcf58\",\"sig\":\"0xf9fb29af97785eb0bc267f5bb5169b68a41e381f06595a80b48f4d71f4b0e60c62a8532431fdea72adcf3d21057cad4c7b4a73bde545c39c966333ba26c7543d1c\",\"to\":\"0x8488277e4221ccd61a854c6bc9de027a3dd7fafa\"}"
	authSign = "{\"block_number\":\"0x989680\",\"exchanger_owner\":\"0xbf5fabb29d464b41eaf88096654dd813ad7bcf58\",\"sig\":\"0xf9fb29af97785eb0bc267f5bb5169b68a41e381f06595a80b48f4d71f4b0e60c62a8532431fdea72adcf3d21057cad4c7b4a73bde545c39c966333ba26c7543d1c\",\"to\":\"0x8488277e4221ccd61a854c6bc9de027a3dd7fafa\"}"
	authSign = "{\"block_number\":\"0xcd09\",\"exchanger_owner\":\"0xc9176a2386d78aa7ea1e8b9e602eb6b5f6dd37b0\",\"sig\":\"0xc9fc6d7f11be5c3120c3013a391d443c2818e64dbe94f0426f40e8461603899969059fcc75c1f6043ac7e7b610118d54fd4203a5f39e5f560e7c6f4eba4d8c3a1b\"}"
	err = json.Unmarshal([]byte(authSign), &Exchangerauth)
	fmt.Println(Exchangerauth)
	if err != nil {
		log.Println("AuthExchangerMint()  Unmarshal() err=", err)

	}
	//key, err := crypto.GenerateKey()

	//updateP:=key.D.String()
	//fmt.Println(updateP)
	//err= nd.DefaultCountry()
	//if err != nil {
	//	t.Errorf("GetSysParams() err=%s", err)
	//}
}

func TestGetRecove(t *testing.T) {
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

	data := map[string]string{"block_number": "0x936e", "exchanger_owner": "0x48cae23c1e43ce233952d2b15b6461dba83767d8", "sig": "0xf85c89b71602fad9fe9fc6e710df024f71708f8c3cadab7b92cfdc3f70e2ef9a774d11250068fb5e07ec2b0ab799fc93c0a9d510d47c2f081bf4d49a2567cee61c", "to": "0x9147e89e031d7466a79aced31e2f8ab2e80ab7da"}
	data2 := "{\"block_number\":\"0x936e\",\"exchanger_owner\":\"0x48cae23c1e43ce233952d2b15b6461dba83767d8\",\"to\":\"0x9147e89e031d7466a79aced31e2f8ab2e80ab7da\"}"

	sigData, _ := hexutil.Decode(data["sig"])
	sigData[64] -= 27
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len([]byte(data2)), data2)
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	hash := hasher.Sum(nil)
	fmt.Println("sigdebug hash=", hexutil.Encode(hash))
	rpk, err := crypto.SigToPub(hash, sigData)
	fmt.Println(crypto.PubkeyToAddress(*rpk))

	if err != nil {
	}
}

func TestNftDb_SigRecover(t *testing.T) {
	//ss := fmt.Sprintf("%02x", "321")
	//fmt.Println(ss)
	CreatedAt := time.Time{}
	fmt.Println(CreatedAt)
	m := make(map[string]string)
	m["name"] = "1"
	_, ok := m["name"]
	fmt.Println(ok)
	if !ok {
		m["name"] = "name"
	} else {
		fmt.Println("data err")
	}
	for i := 0; i < 256; i++ {
		ss := fmt.Sprintf("%02x", i)
		fmt.Printf("%v,", ss)
	}

	rawData := "{\"exchanger_owner\":\"0xe39e6081f1da6a2a40b5dd0758dac3a52746a8e3\",\"to\":\"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900\",\"block_number\":\"0x52c8c338e\",\"sig\":\"\"}"
	rawData = "0xe39e6081f1da6a2a40b5dd0758dac3a52746a8e30x7fbc8ad616177c6519228fca4a7d9ec7d18049000x52c8c338e"
	sig := "e98fe81d45bccd5941029d2a3c3dfd40da00df647bb52d9d5c0e13317ae4580921354ba533349b133c538eeddb76103556dcf59f4f7ca28a3b459b1590814fd11b"
	addr := "0xe39e6081f1da6a2a40b5dd0758dac3a52746a8e3"
	verificationAddr, err := GetEthAddr(rawData, sig)
	if err != nil {
		fmt.Println(err)
	}
	verificationAddrS := verificationAddr.String()
	verificationAddrS = strings.ToLower(verificationAddrS)

	addr = strings.ToLower(addr)
	fmt.Printf("sigdebug verificationAddrS = [%s], approveAddr's addr = [%s]\n", verificationAddrS, addr)
	if verificationAddrS == addr {
		fmt.Println("sigdebug verify [Y]")
	}
	fmt.Println("sigdebug verify [N]")
}
func TestSnftSearch(t *testing.T) {
	nd, nerr := NewNftDb(sqldsnT)
	if nerr != nil {
		fmt.Printf("connect database err = %s\n", nerr)
	}
	defer nd.Close()

	collectperiod := SnftCollectPeriod{}
	err := nd.db.Model(&SnftCollectPeriod{}).Where("period=? and collect =? and local=?", "1", "1", "1").First(&collectperiod)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			snftperiod := &SnftCollectPeriod{}
			snftperiod.Period = "1"
			snftperiod.Collect = "1"
			snftperiod.Local = "1"
			//err = tx.Model(&SnftCollectPeriod{}).Create(&snftperiod)
			err = nd.db.Model(&SnftCollectPeriod{}).Create(SnftCollectPeriodRec{})
			if err.Error != nil {
				fmt.Printf("SetSnftPeriod() create SnftCollectPeriod err=%s", err.Error)
				//return err.Error
			}
			//continue
		}
		fmt.Println("SnftCollectPeriod err =", err.Error)
		//return err.Error
	}
	//admins, err := nd.CollectSearch("collect,art", "collect")
	//admins, err := nd.GetSnftCollection("1")
	//err = nd.SetCollectSnft("[{\"collect\":\"1\",\"local\":\"1\"}]", "1")
	//if err != nil {
	//	t.Errorf("GetSysParams() err=%s", err)
	//}
	//fmt.Println(admins)
}

func TestSnftPeriod(t *testing.T) {
	nd, nerr := NewNftDb(sqldsnT)
	if nerr != nil {
		fmt.Printf("connect database err = %s\n", nerr)
	}
	defer nd.Close()
	existnft := Nfts{}
	existsnft := Snfts{}

	//serr := nd.SetCollectSnft("[{\"collect\":\"4687775155334\",\"local\":\"1\"}]", "1")
	//if serr != nil {
	//	fmt.Printf("input nft err=%v", serr)
	//	//return s.New("input nft err")
	//}
	collect := Collects{}
	str := "("
	str += "qqqqqqqqqqqqq,demo01"
	str += ")"
	err := nd.db.Model(&Collects{}).Limit(18).Offset(1).Find(&collect)
	if err.Error != nil {
		fmt.Printf("input nft err=%s", err.Error)
		//return s.New("input nft err")
	}
	fmt.Println(collect)
	err = nd.db.Model(&Collects{}).Where("(length( name ) - length( REPLACE ( name,',','' )) )= ?", 4).Find(&collect)
	if err.Error != nil {
		fmt.Printf("input nft err=%s", err.Error)
		//return s.New("input nft err")
	}
	fmt.Println(collect)
	err = nd.db.Model(&Collects{}).Where("name not in ? ", []string{"qqqqqqqqqqqqq", "demo01"}).First(&collect)
	if err.Error != nil {
		fmt.Printf("input nft err=%s", err.Error)
		//return s.New("input nft err")
	}
	fmt.Println(collect)

	err = nd.db.Model(&Collects{}).Where("name = ? ", "qqqqqqqqqqqqq").First(&collect)
	err = nd.db.Model(&Collects{}).Unscoped().Where("name = ?", "qqqqqqqqqqqqq").Delete(&Collects{})
	if err.Error != nil {
		fmt.Printf("input nft err=%s", err.Error)
		//return s.New("input nft err")
	}
	fmt.Println()

	err = nd.db.Model(&Snfts{}).Where("tokenid = ? ", "4687775155334").First(&existsnft)
	err = nd.db.Model(&Snfts{}).Where("tokenid =?", "4687775155334").Update("meta", "1")

	//snftphase := []SnftPhase{}
	//snftphasecollect := []QueryPeriod{}
	////var snftcollectrec SnftPhase
	//db := nft.db.Model(&SnftPhase{}).Where("accedvote = ", "true").Find(&snftphase)
	//
	if err.Error != nil {
		fmt.Printf("input nft err=%s", err.Error)
		//return s.New("input nft err")
	}

	fmt.Println()
	insnft := Snfts{}
	insnft.Name = existnft.Name
	insnft.Desc = existnft.Desc
	insnft.Ownaddr = existnft.Ownaddr
	insnft.Image = existnft.Image
	insnft.Md5 = existnft.Md5
	insnft.Meta = existnft.Meta
	insnft.Nftmeta = existnft.Nftmeta
	insnft.Url = existnft.Url
	insnft.Contract = existnft.Contract
	insnft.Tokenid = existnft.Tokenid
	insnft.Nftaddr = existnft.Nftaddr
	insnft.Count = 1
	insnft.Approve = existnft.Approve
	insnft.Categories = existnft.Categories
	insnft.Hide = existnft.Hide
	insnft.Signdata = existnft.Signdata
	insnft.Createaddr = existnft.Createaddr
	insnft.Verifyaddr = existnft.Verifyaddr
	insnft.Currency = existnft.Currency
	insnft.Price = existnft.Price
	insnft.Royalty = existnft.Royalty
	insnft.Collection = "2"
	insnft.Local = "1"
	err = nd.db.Model(&Snfts{}).Create(&insnft)
	if err != nil {
		fmt.Printf("SetCollectSnt() create  snft err=%v", err.Error)
		//return errors.New("SetCollectSnt() create  snft err")
	}
	//admins, err := nd.GetSnftPeriod()
	//admins, err := nd.GetSnftCollection("1")
	//err = nd.SetCollectSnft("[{\"collect\":\"1\",\"local\":\"1\"}]", "1")
	//if err != nil {
	//	t.Errorf("GetSysParams() err=%s", err)
	//}
	//fmt.Println(admins)
}
func TestIpfs(t *testing.T) {
	//url := "http://192.168.1.237:5001"
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", rand.Int63())
		fmt.Println("UploadNft() s=", s, ",len =", len(s))
		//s = s[len(s)-13:]
		//NewTokenid := s
		//fmt.Println("UploadNft() NewTokenid=", NewTokenid)
	}

	nft := []NftRecord{}
	nftlist := "[{\"ownaddr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"md5\":\"699412f4e1d2a0b65e144baa5679ab44\",\"name\":\"demo\",\"desc\":\"demo\",\"meta\":\"/ipfs/QmUJ8uXxZaNsN3uspzkjexrTWqQbMG6r4kFwCS6erqfBDq\",\"source_url\":\"/ipfs/QmdqM4HQrAWdpihPfy4Rgq4exGn81af26zXsUQeAvSLCv1\",\"nft_contract_addr\":\"0x0000000000000000000000000000000000000000\",\"nft_token_id\":\"9257838269658\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x53e4a985c7167b8f607a5c4647d5f768959d3429e318397cf6090bdfc667e8ac79f419a50679b70b50cc9bc9abf7991382f9997731d470a84932d0b34603351c1c\",\"user_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":200,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651734870,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651734870,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"md5\":\"699412f4e1d2a0b65e144baa5679ab44\",\"name\":\"demo\",\"desc\":\"demo\",\"meta\":\"/ipfs/QmdhfJ5yFdZQ9DPESbPH9D4kBCjn1fEzWA94qHBnkMrhsi\",\"source_url\":\"/ipfs/QmdqM4HQrAWdpihPfy4Rgq4exGn81af26zXsUQeAvSLCv1\",\"nft_contract_addr\":\"0x0000000000000000000000000000000000000000\",\"nft_token_id\":\"2019716535149\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x53e4a985c7167b8f607a5c4647d5f768959d3429e318397cf6090bdfc667e8ac79f419a50679b70b50cc9bc9abf7991382f9997731d470a84932d0b34603351c1c\",\"user_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":200,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651731071,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651731071,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"md5\":\"b698c829f9a7c9b3b763213067d3d580\",\"name\":\"demo\",\"desc\":\"demo\",\"meta\":\"/ipfs/QmVi7cPaXFfJppLEsWmxPhWpAE845S2J6MwP7aXGkcWMM4\",\"source_url\":\"/ipfs/QmbFaundABKbfUN39NkwZGgJWSbTDf9X7u45TwqC1hUv7y\",\"nft_contract_addr\":\"0x0000000000000000000000000000000000000000\",\"nft_token_id\":\"5116009548960\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x7c2100075eef0216b188d7a26b82ccb8c6b50998315c7ec1511e522c22ee70416df4bf38041e0e159bb007af4b7a524598f0ccedec3703aa2ff9e69f0e8d658c1c\",\"user_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":200,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651728874,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651728874,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"27977f709b648b41efa81bdafa02325f\",\"name\":\"0001\",\"desc\":\"1111111111\",\"meta\":\"/ipfs/QmSzSvBAPYexd8sTrtinHWHWVnhqqgRUkG7EDRbaiiu8k5\",\"source_url\":\"/ipfs/QmcPZf28bXGyJpocZqZz3t1h26p7CfbopJyuTFQ412kqpN\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"4993190995868\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"collections\":\"0001q\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0xd30a8430d2b09d91a462aa479bc55ddaae839016b973d8451a32b45f62b6eeda300a58216d82e31ac7ee7c1995f260e0c0d5ce0988092f39eaf0fe78a7ac33381b\",\"user_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651726806,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651726806,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"27977f709b648b41efa81bdafa02325f\",\"name\":\"0001\",\"desc\":\"1111111111\",\"meta\":\"/ipfs/QmSzSvBAPYexd8sTrtinHWHWVnhqqgRUkG7EDRbaiiu8k5\",\"source_url\":\"/ipfs/QmcPZf28bXGyJpocZqZz3t1h26p7CfbopJyuTFQ412kqpN\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"6139224986592\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"collections\":\"0001q\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0xd30a8430d2b09d91a462aa479bc55ddaae839016b973d8451a32b45f62b6eeda300a58216d82e31ac7ee7c1995f260e0c0d5ce0988092f39eaf0fe78a7ac33381b\",\"user_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651720801,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651720801,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"27977f709b648b41efa81bdafa02325f\",\"name\":\"0001\",\"desc\":\"1111111111\",\"meta\":\"/ipfs/QmSzSvBAPYexd8sTrtinHWHWVnhqqgRUkG7EDRbaiiu8k5\",\"source_url\":\"/ipfs/QmcPZf28bXGyJpocZqZz3t1h26p7CfbopJyuTFQ412kqpN\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"3397519365732\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"collections\":\"111111111111\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x93ca5a21dc7e8b380d7655a47ae62f83f86a9a6c619d4fd008728f02be7ccdf9386b685d1164a397e3744f07ab7709e9dcfedaaf27c2f1ae317b3170d91b1ca71b\",\"user_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651715796,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651715796,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"a304e1245655cc7ac0b2392fd2bc8571\",\"name\":\"0001\",\"desc\":\"33333\",\"meta\":\"/ipfs/QmcvHaL6ZPcSdinziAejFQ6akTMQpToqXRzFirrYm8NaLP\",\"source_url\":\"/ipfs/QmXQDiihuJRzyPMN1nuh45ysbzEeY7XtGQeNtjnChvQMRK\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"6023705250027\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"collections\":\"111111111111\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x39ac33771fbc707f2feb3e3570e5d55ff09da4ad614cb7987ea25aaa2cd2039c21a0011f5136bf66de4a08297cd0762a9fba36422960747f3233e88c42b011da1c\",\"user_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651715713,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651715713,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"a304e1245655cc7ac0b2392fd2bc8571\",\"name\":\"0001\",\"desc\":\"33333\",\"meta\":\"/ipfs/QmcvHaL6ZPcSdinziAejFQ6akTMQpToqXRzFirrYm8NaLP\",\"source_url\":\"/ipfs/QmXQDiihuJRzyPMN1nuh45ysbzEeY7XtGQeNtjnChvQMRK\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9263934753983\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"collections\":\"111111111111\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x39ac33771fbc707f2feb3e3570e5d55ff09da4ad614cb7987ea25aaa2cd2039c21a0011f5136bf66de4a08297cd0762a9fba36422960747f3233e88c42b011da1c\",\"user_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651715470,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651715470,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"3a3ae0f7378df518acb4f56f816318ff\",\"name\":\"11111111111\",\"desc\":\"1111111111111\",\"meta\":\"/ipfs/QmR4TtmLDovKeic4hGLdUUiLoJz2B3Y6GSujuf4NYCqP2H\",\"source_url\":\"/ipfs/QmUNGcUwVjDfTJP4bvmtbYPukXkqHeH78ktQBfwWdvFMZS\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"8387211489329\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"collections\":\"111111111111\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0xbb4fec86122b72d78f9562293e714fb090ec0b982dc7f9fa778a9e1a0324252c692ce228ece46e8688e6cc90952c545b8a5c890fda2435ed0979d8dfdbd5c3621b\",\"user_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651714106,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651714106,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x0a0b071c3b7c7ee0e973993359833526b30c165d\",\"md5\":\"a304e1245655cc7ac0b2392fd2bc8571\",\"name\":\"0001\",\"desc\":\"33333\",\"meta\":\"/ipfs/QmcvHaL6ZPcSdinziAejFQ6akTMQpToqXRzFirrYm8NaLP\",\"source_url\":\"/ipfs/QmXQDiihuJRzyPMN1nuh45ysbzEeY7XtGQeNtjnChvQMRK\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9958266705353\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x0a0b071c3b7c7ee0e973993359833526b30c165d\",\"collections\":\"0001\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0xcd979bfa092f15d602c31f132ffb2f318ea383313f78457f80f059e9b1927337486d051901469fa4c23989f7f5b8bd17767fc35f9f8f368bc9621da6086db2d41b\",\"user_addr\":\"0x0a0b071c3b7c7ee0e973993359833526b30c165d\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651666295,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651666295,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x88c41ce51023b2891dc7b6ae4c87c1a67163a46f\",\"md5\":\"b698c829f9a7c9b3b763213067d3d580\",\"name\":\"test0\",\"desc\":\"test0\",\"meta\":\"/ipfs/QmUXtQJUndSvt9NEfnZk4SDRYsfpkZ4nedZ2hsUGJBZSoL\",\"source_url\":\"/ipfs/QmbFaundABKbfUN39NkwZGgJWSbTDf9X7u45TwqC1hUv7y\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"2988025786196\",\"snft\":\"\",\"nft_address\":\"0x000000000000000000000000000000000000001d\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x8fbf399d77bc8c14399afb0f6d32dbe22189e169\",\"collections\":\"test0\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x436031e8cb722aa2e9d5d34fe1d07be128a7a052f2ddab2e43be86267cd30aca72bc71e997ab1121c5622e763ad578ca67e6ad46ae36ac249ed8e5b0f1d32a321b\",\"user_addr\":\"0x8fbf399d77bc8c14399afb0f6d32dbe22189e169\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":200,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":1000000000,\"last_trans_time\":1651151799,\"createdate\":1651147739,\"favorited\":0,\"transcnt\":1,\"transamt\":1000000000,\"verified\":\"Passed\",\"vrf_time\":1651147739,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x88c41ce51023b2891dc7b6ae4c87c1a67163a46f\",\"md5\":\"4fc9dc60dbf9bee3fd9fbe78ccf3b942\",\"name\":\"1231312\",\"desc\":\"123123131\",\"meta\":\"/ipfs/QmUKngpGK1nozxybRZQ1sasqzAYtPuB7gZtAnNcS7yHe7N\",\"source_url\":\"/ipfs/QmS2duiZ7ZUWfw7wMBEnyDB219EnDvLqHUQfFoZLmiJer6\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"7929082544413\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x88c41ce51023b2891dc7b6ae4c87c1a67163a46f\",\"collections\":\"1231231312\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x744c39d018fff6f4db65aa1264bf6ef8fc72be17cdccc627efa32cbf4e19458471c0b8e2f4ea4cdf83ff646cd4cf0d33191136a0a19d00038833cba260c185911b\",\"user_addr\":\"0x88c41ce51023b2891dc7b6ae4c87c1a67163a46f\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651147271,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651147271,\"selltype\":\"FixPrice\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"md5\":\"b698c829f9a7c9b3b763213067d3d580\",\"name\":\"demo\",\"desc\":\"demo\",\"meta\":\"/ipfs/QmZAR7QbjCoTLvxd2PihbveQpP2cARjf5BM4CA8uQMKwoD\",\"source_url\":\"/ipfs/QmbFaundABKbfUN39NkwZGgJWSbTDf9X7u45TwqC1hUv7y\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9044355746691\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"collections\":\"1232131231231231\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0x7c1603350ee8204e82c24b0cd3ad2b2fe24b0dedeb02ec556bc02cf467c25c206067f453f847eeca405603c38f606298b0d545cd0d6a7d87da2050d06929c6da1b\",\"user_addr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":200,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651126345,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651126345,\"selltype\":\"FixPrice\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x8b07aff2327a3b7e2876d899cafac99f7ae16b10\",\"md5\":\"b698c829f9a7c9b3b763213067d3d580\",\"name\":\"demo\",\"desc\":\"demo\",\"meta\":\"/ipfs/QmZAR7QbjCoTLvxd2PihbveQpP2cARjf5BM4CA8uQMKwoD\",\"source_url\":\"/ipfs/QmbFaundABKbfUN39NkwZGgJWSbTDf9X7u45TwqC1hUv7y\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"4794394950601\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x8b07aff2327a3b7e2876d899cafac99f7ae16b10\",\"collections\":\"demo\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0xfe4645d2d8c598886339ce7663a85013ee290150edc6e5dc7ad2c3d62840d6d552efd09b821406beecc40945b22fa71aabe632c039ad489b26cd6d8a66b869d21b\",\"user_addr\":\"0x8b07aff2327a3b7e2876d899cafac99f7ae16b10\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":200,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651102724,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651102724,\"selltype\":\"FixPrice\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"md5\":\"4fc9dc60dbf9bee3fd9fbe78ccf3b942\",\"name\":\"demo02\",\"desc\":\"1232131\",\"meta\":\"/ipfs/QmZH29cozsR2TPCY8yfLeaeLuKGHuHae7v7pG9S83GZZeK\",\"source_url\":\"/ipfs/QmS2duiZ7ZUWfw7wMBEnyDB219EnDvLqHUQfFoZLmiJer6\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9234950143618\",\"snft\":\"\",\"nft_address\":\"\",\"count\":1,\"approve\":\"\",\"categories\":\"Art\",\"collection_creator_addr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"collections\":\"1232131231231231\",\"asset_sample\":\"\",\"hide\":\"true\",\"sig\":\"0xb440606ca508df4e080b58788abac6331b087ff39eb6ce66b8f10c27bfaa5aa87c4ad45dfe508fbc79796e10dc7d02db26a54f9ce1b6314230660984d71d06691b\",\"user_addr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":100,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1651023540,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1651023540,\"selltype\":\"FixPrice\",\"sellprice\":0,\"mintstate\":\"NoMinted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf50cbaffa72cc902de3f4f1e61132d858f3361d9\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000010000\",\"source_url\":\"/ipfs/QmWmrjHBoEC1jgd2e9dCJLuZPoyGtHCMcrY5oDXy3YCjcL\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"8401083909235\",\"snft\":\"\",\"nft_address\":\"0x8000000000000000000000000000000000010000\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000010.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976482,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976482,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xd8861d235134ef573894529b577af28ae0e3449c\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000FFfF\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9080540620005\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000ffff\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976482,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976482,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xfff531a2da46d051fde4c47f042ee6322407df3f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000fffE\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"4949653114920\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fffe\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976481,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976481,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xedfc22e9cfb4e24815c3a12e81bf10cab9ce4d26\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000ffFd\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"3476174801171\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fffd\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976481,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976481,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x033eecd45d8c8ec84516359f39b11c260a56719e\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000FFFC\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"3297330970609\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fffc\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976481,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976481,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x83c43f6f7bb4d8e429b21ff303a16b4c99a59b05\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000FfFB\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"4325144237707\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fffb\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976480,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976480,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x492fe79f6f162aea18c88e40228feeb9c6bca2aa\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000FfFA\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9764091788416\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fffa\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976480,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976480,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf6e55448acd5bc3fa5e2666344a89eada16ce65a\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000ffF9\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"8246957956815\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fff9\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976480,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976480,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf8d4d274558cbb9a8cdc813b5d118ded10353a0f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000fFF8\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"5813710339255\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fff8\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976479,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976479,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x1b27137606881995d8249bc687325fe080de0377\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x800000000000000000000000000000000000FFf7\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"6700831544522\",\"snft\":\"\",\"nft_address\":\"0x800000000000000000000000000000000000fff7\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x800000000000000000000000000000000000f.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650976479,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650976479,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x1b27137606881995d8249bc687325fe080de0377\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000000\",\"source_url\":\"/ipfs/QmWmrjHBoEC1jgd2e9dCJLuZPoyGtHCMcrY5oDXy3YCjcL\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"9163593943075\",\"snft\":\"0x80000000000000000000000000000000000000\",\"nft_address\":\"0x8000000000000000000000000000000000000000\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650945938,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650945938,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x492fe79f6f162aea18c88e40228feeb9c6bca2aa\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000100\",\"source_url\":\"/ipfs/QmcdAU4DG7Az66ja82YxSgDJqQWrogoq3y1B947kQ2wrR3\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"4762786620361\",\"snft\":\"0x80000000000000000000000000000000000001\",\"nft_address\":\"0x8000000000000000000000000000000000000100\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946059,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946059,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x93f24e8a3162b45611ab17a62dd0c95999cda60f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000200\",\"source_url\":\"/ipfs/QmTKrSVUUuy3bLt2mJrPYYbH3Qaf6X11vzb5PeG5MAgkTw\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"1195494277638\",\"snft\":\"0x80000000000000000000000000000000000002\",\"nft_address\":\"0x8000000000000000000000000000000000000200\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946172,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946172,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf50cbaffa72cc902de3f4f1e61132d858f3361d9\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000300\",\"source_url\":\"/ipfs/QmcJFxVQMKR5EVZwN8eKmdoxKpx8MCbkk2uYQAscEbLXVV\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"5626364554966\",\"snft\":\"0x80000000000000000000000000000000000003\",\"nft_address\":\"0x8000000000000000000000000000000000000300\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946262,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946262,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf8d4d274558cbb9a8cdc813b5d118ded10353a0f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000400\",\"source_url\":\"/ipfs/Qmcufskz8gCPL9ohadjZkcAfWWxm3KfZygh9vbd7pXvt5S\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"6899378292727\",\"snft\":\"0x80000000000000000000000000000000000004\",\"nft_address\":\"0x8000000000000000000000000000000000000400\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946375,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946375,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x9d196915f63dbdb97dea552648123655109d98a5\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000500\",\"source_url\":\"/ipfs/QmQf858Mq6qcecssXZdvVcy4NDHQ2pG5zxyk6uaxavLop2\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"3011775454896\",\"snft\":\"0x80000000000000000000000000000000000005\",\"nft_address\":\"0x8000000000000000000000000000000000000500\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946460,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946460,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000600\",\"source_url\":\"/ipfs/QmfPQaKA9quurhKa4F7yP9RPR5Z4gftA7eaVAJd5qCQ3Q9\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"8606937505391\",\"snft\":\"0x80000000000000000000000000000000000006\",\"nft_address\":\"0x8000000000000000000000000000000000000600\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946531,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946531,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x93f24e8a3162b45611ab17a62dd0c95999cda60f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000700\",\"source_url\":\"/ipfs/QmWRojhwrJeRZdXdrz6mLkj5CWeKC1FJg6aM6nv173BjwK\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"1486702902730\",\"snft\":\"0x80000000000000000000000000000000000007\",\"nft_address\":\"0x8000000000000000000000000000000000000700\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946582,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946582,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf6e55448acd5bc3fa5e2666344a89eada16ce65a\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000800\",\"source_url\":\"/ipfs/QmUZuREi9Nr3eHapbCFLvMb82MJmMSzUVK2TQU8aetzSED\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"4041688175359\",\"snft\":\"0x80000000000000000000000000000000000008\",\"nft_address\":\"0x8000000000000000000000000000000000000800\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946630,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946630,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xd8861d235134ef573894529b577af28ae0e3449c\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000900\",\"source_url\":\"/ipfs/QmbSmDX9iWbqjKqCs3DzvLmzCHJxJELwSzPzNi4DGD6vUd\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"6546006756844\",\"snft\":\"0x80000000000000000000000000000000000009\",\"nft_address\":\"0x8000000000000000000000000000000000000900\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946675,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946675,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xd8861d235134ef573894529b577af28ae0e3449c\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000a00\",\"source_url\":\"/ipfs/QmX3cAiUMv1MQsCGUA6xSskNBEhumiitM6q1hf6VdDcuCH\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"6946198360002\",\"snft\":\"0x8000000000000000000000000000000000000a\",\"nft_address\":\"0x8000000000000000000000000000000000000a00\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946719,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946719,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x82ce2100a8abdd8862746b15dde245115042476f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000b00\",\"source_url\":\"/ipfs/QmW8CaCKCarLzYjtbqrc8oSxADhiKqpYK5qm57aonLMWn2\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"5461931460995\",\"snft\":\"0x8000000000000000000000000000000000000b\",\"nft_address\":\"0x8000000000000000000000000000000000000b00\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946764,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946764,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xc9ab0174cacb94209ba0d4fd36ab01767c47ac5a\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000C00\",\"source_url\":\"/ipfs/QmXxFyv9sUV9yodaZavhP5h4nF4mmEC7SXSU1fpGykrJWP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"5817269610026\",\"snft\":\"0x8000000000000000000000000000000000000c\",\"nft_address\":\"0x8000000000000000000000000000000000000c00\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946810,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946810,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xd8861d235134ef573894529b577af28ae0e3449c\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000D00\",\"source_url\":\"/ipfs/QmYcnvEGo4yLAjuQEUFY3VajgRzNm9UvAp5h9sVwshVoGm\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"2557496972651\",\"snft\":\"0x8000000000000000000000000000000000000d\",\"nft_address\":\"0x8000000000000000000000000000000000000d00\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946855,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946855,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xd8861d235134ef573894529b577af28ae0e3449c\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000E00\",\"source_url\":\"/ipfs/QmctEG5CZ9uRoJA63wMxGtucM3z6ddbL3fHT69Kmd2NTja\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"1648145636364\",\"snft\":\"0x8000000000000000000000000000000000000e\",\"nft_address\":\"0x8000000000000000000000000000000000000e00\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946918,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946918,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xa81b2900a01b7282421e98fbfb9356dc0e684313\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000000F00\",\"source_url\":\"/ipfs/QmQU5EcS6uZkZmmgyo5PWzFR57RsVDKz7idmhTWs37maNP\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"8066415176867\",\"snft\":\"0x8000000000000000000000000000000000000f\",\"nft_address\":\"0x8000000000000000000000000000000000000f00\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000000.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650946973,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650946973,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x93f24e8a3162b45611ab17a62dd0c95999cda60f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001000\",\"source_url\":\"/ipfs/QmWmrjHBoEC1jgd2e9dCJLuZPoyGtHCMcrY5oDXy3YCjcL\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"2825910380042\",\"snft\":\"0x80000000000000000000000000000000000010\",\"nft_address\":\"0x8000000000000000000000000000000000001000\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947022,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947022,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xfff531a2da46d051fde4c47f042ee6322407df3f\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001100\",\"source_url\":\"/ipfs/QmcdAU4DG7Az66ja82YxSgDJqQWrogoq3y1B947kQ2wrR3\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"5320654705092\",\"snft\":\"0x80000000000000000000000000000000000011\",\"nft_address\":\"0x8000000000000000000000000000000000001100\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947079,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947079,\"selltype\":\"BidPrice\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x9d196915f63dbdb97dea552648123655109d98a5\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001200\",\"source_url\":\"/ipfs/QmTKrSVUUuy3bLt2mJrPYYbH3Qaf6X11vzb5PeG5MAgkTw\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"8294146895760\",\"snft\":\"0x80000000000000000000000000000000000012\",\"nft_address\":\"0x8000000000000000000000000000000000001200\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947145,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947145,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf6e55448acd5bc3fa5e2666344a89eada16ce65a\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001300\",\"source_url\":\"/ipfs/QmcJFxVQMKR5EVZwN8eKmdoxKpx8MCbkk2uYQAscEbLXVV\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"5816576433698\",\"snft\":\"0x80000000000000000000000000000000000013\",\"nft_address\":\"0x8000000000000000000000000000000000001300\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947209,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947209,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x9d196915f63dbdb97dea552648123655109d98a5\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001400\",\"source_url\":\"/ipfs/Qmcufskz8gCPL9ohadjZkcAfWWxm3KfZygh9vbd7pXvt5S\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"7703552488431\",\"snft\":\"0x80000000000000000000000000000000000014\",\"nft_address\":\"0x8000000000000000000000000000000000001400\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947285,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947285,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x033eecd45d8c8ec84516359f39b11c260a56719e\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001500\",\"source_url\":\"/ipfs/QmQf858Mq6qcecssXZdvVcy4NDHQ2pG5zxyk6uaxavLop2\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"2687813192166\",\"snft\":\"0x80000000000000000000000000000000000015\",\"nft_address\":\"0x8000000000000000000000000000000000001500\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947351,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947351,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0x1b27137606881995d8249bc687325fe080de0377\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001600\",\"source_url\":\"/ipfs/QmfPQaKA9quurhKa4F7yP9RPR5Z4gftA7eaVAJd5qCQ3Q9\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"2374905239320\",\"snft\":\"0x80000000000000000000000000000000000016\",\"nft_address\":\"0x8000000000000000000000000000000000001600\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947422,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947422,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xc9ab0174cacb94209ba0d4fd36ab01767c47ac5a\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001700\",\"source_url\":\"/ipfs/QmWRojhwrJeRZdXdrz6mLkj5CWeKC1FJg6aM6nv173BjwK\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"3158726984294\",\"snft\":\"0x80000000000000000000000000000000000017\",\"nft_address\":\"0x8000000000000000000000000000000000001700\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947489,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947489,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0},{\"ownaddr\":\"0xf50cbaffa72cc902de3f4f1e61132d858f3361d9\",\"md5\":\"774f15fe95dfb25ed10110d631bdac94\",\"name\":\"DemoSNFT\",\"desc\":\"nft desc\",\"meta\":\"/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7/0x8000000000000000000000000000000000001800\",\"source_url\":\"/ipfs/QmUZuREi9Nr3eHapbCFLvMb82MJmMSzUVK2TQU8aetzSED\",\"nft_contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"nft_token_id\":\"1153184310239\",\"snft\":\"0x80000000000000000000000000000000000018\",\"nft_address\":\"0x8000000000000000000000000000000000001800\",\"count\":1,\"approve\":\"\",\"categories\":\"Virtual World\",\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"collections\":\"0x8000000000000000000000000000000000001.合集\",\"asset_sample\":\"\",\"hide\":\"\",\"sig\":\"\",\"user_addr\":\"0x085abc35ed85d26c2795b64c6ffb89b68ab1c479\",\"vrf_addr\":\"\",\"currency\":\"\",\"price\":0,\"royalty\":250,\"paychan\":\"\",\"trans_cur\":\"\",\"transprice\":0,\"last_trans_time\":0,\"createdate\":1650947563,\"favorited\":0,\"transcnt\":0,\"transamt\":0,\"verified\":\"Passed\",\"vrf_time\":1650947563,\"selltype\":\"NotSale\",\"sellprice\":0,\"mintstate\":\"Minted\",\"extend\":\"\",\"offernum\":0,\"maxbidprice\":0}]"
	err := json.Unmarshal([]byte(nftlist), &nft)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(nft)

	nftcollect := []CollectRec{}
	collect := "[{\"collection_creator_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"name\":\"0x8000000000000000000000000000000000000.合集\",\"img\":\"\",\"contract_addr\":\"0x0000000000000000000000000000000000000000\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0x8fbf399d77bc8c14399afb0f6d32dbe22189e169\",\"name\":\"test0\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"name\":\"0001q\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0xb31b41e5ef219fb0cc9935ad914158cf8970db44\",\"name\":\"111111111111\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0x0a0b071c3b7c7ee0e973993359833526b30c165d\",\"name\":\"0001\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0x88c41ce51023b2891dc7b6ae4c87c1a67163a46f\",\"name\":\"1231231312\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0x8b07aff2327a3b7e2876d899cafac99f7ae16b10\",\"name\":\"demo\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0},{\"collection_creator_addr\":\"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c\",\"name\":\"1232131231231231\",\"img\":\"\",\"contract_addr\":\"0x01842a2cf56400a245a56955dc407c2c4137321e\",\"desc\":\"\",\"contracttype\":\"\",\"categories\":\"\",\"total_count\":0}]"
	err = json.Unmarshal([]byte(collect), &nftcollect)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(nftcollect)

	rand.Seed(time.Now().UnixNano())
	//randomNum := rand.Intn(50)
	//for i := 0; i < 10; i++ {
	//	fmt.Printf("%v\t", rand.Intn(50))
	//	fmt.Println(nft[rand.Intn(50)])
	//}
	nftnumber := 0
	collectnumber := 0
	namenumer := 0
	nfts := []nftInfo{}
	for i := 0; i < 16; i++ {
		//fmt.Printf("%v\t", rand.Intn(50))
		collectnumber = rand.Intn(8)
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 16; i++ {
			//fmt.Printf("%v\t", rand.Intn(8))
			nftnumber = rand.Intn(50)
			ss := fmt.Sprintf("%02x", namenumer)
			var nftMeta nftInfo
			nftMeta.CreatorAddr = nft[nftnumber].Createaddr
			nftMeta.Contract = nft[nftnumber].Contract
			nftMeta.Name = nft[nftnumber].Name
			nftMeta.Desc = nft[nftnumber].Desc
			nftMeta.Category = nft[nftnumber].Categories
			nftMeta.Royalty = strconv.Itoa(nft[nftnumber].Royalty)
			nftMeta.SourceImageName = ss
			nftMeta.FileType = ""
			nftMeta.SourceUrl = nft[nftnumber].Url
			nftMeta.Md5 = nft[nftnumber].Md5
			nftMeta.CollectionsName = nftcollect[collectnumber].Name
			nftMeta.CollectionsCreator = nftcollect[collectnumber].Createaddr
			nftMeta.CollectionsExchanger = nftcollect[collectnumber].Contract
			nftMeta.CollectionsCategory = nftcollect[collectnumber].Categories
			nftMeta.CollectionsImgUrl = nft[collectnumber].Url
			//nftMeta.Attributes = attributes
			namenumer++
			nfts = append(nfts, nftMeta)

		}
	}
	metaStr, err := json.Marshal(&nfts)
	meta, err := SaveToIpfs(string(metaStr))
	if err != nil {
		t.Errorf("GetSysParams() err=%s", err)
	}
	fmt.Println(meta)
	fmt.Println(namenumer, collectnumber, nftnumber)
	//for i := 0; i < 255; i++ {
	//	fmt.Printf("%v\t", rand.Intn(50))
	//	fmt.Println(nft[rand.Intn(50)])
	//	ss := fmt.Sprintf("%02x", i)
	//	var nftMeta nftInfo
	//	nftMeta.CreatorAddr = "0xb31b41e5ef219fb0cc9935ad914158cf8970db44"
	//	nftMeta.Contract = ""
	//	nftMeta.Name = "testsnft"
	//	nftMeta.Desc = "testdesc"
	//	nftMeta.Category = "Art"
	//	nftMeta.Royalty = "100"
	//	nftMeta.SourceImageName = ss
	//	nftMeta.FileType = ""
	//	nftMeta.SourceUrl = "/ipfs/" + "QmSzSvBAPYexd8sTrtinHWHWVnhqqgRUkG7EDRbaiiu8k5"
	//	nftMeta.Md5 = "27977f709b648b41efa81bdafa02325f"
	//	nftMeta.CollectionsName = "0001q"
	//	nftMeta.CollectionsCreator = "0xb31b41e5ef219fb0cc9935ad914158cf8970db44"
	//	nftMeta.CollectionsExchanger = ""
	//	nftMeta.CollectionsCategory = "Art"
	//	nftMeta.CollectionsImgUrl = ""
	//	//nftMeta.Attributes = attributes
	//	metaStr, err := json.Marshal(&nftMeta)
	//	meta, err := SaveToIpfs(string(metaStr))
	//	if err != nil {
	//		t.Errorf("GetSysParams() err=%s", err)
	//	}
	//	fmt.Println(meta)
	//}

}

func TestDirs(t *testing.T) {

	param := strings.Split("724761865667,5118325964426,3935361596550,2258522108353,1979835067310,5958807637084,8394128951493,4619720007408,9361396615789,8916548021829,4094151285352,6930117693541,4660839740217,1907681606499,7564023841144,6965060794590", ",")

	if len(param) == 16 || param[0] == "1" {
		fmt.Println(len(param))

	} else {
		fmt.Println("SaveToIpfs() err=")
		fmt.Println(len(param))
	}
	fmt.Println(param)
	if param != nil {
		fmt.Println("SaveToIpfs() err=")
		fmt.Println(len(param))
	}
	os.Mkdir("./snft", 0777)
	for i := 0; i < 10; i++ {
		file := "./snft/" + fmt.Sprintf("%02x", i)
		file6, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0766)
		if err != nil {
			fmt.Println("creat file error")
		}
		//defer file6.Close()

		//metaStr, merr := json.Marshal(&nftMeta)
		//if merr != nil {
		//	fmt.Println("SetVoteSnftPeriod() save nftmeta info err=", merr)
		//}
		metastr := "asasasas"
		file6.WriteString(metastr)
	}

	file := "./snft/" + fmt.Sprintf("%02x", "50")
	file6, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("creat file error")
	}
	metastr := "asasasas"
	file6.WriteString(metastr)
	url := "192.168.1.237:5001"
	spendT := time.Now()
	s := shell.NewShell(url)
	s.SetTimeout(5 * time.Second)
	mhash, err := s.AddDir("./snft")
	if err != nil {
		log.Println("SaveToIpfs() err=", err)
	}
	fmt.Println(mhash)
	fmt.Printf("SaveToIpfs  Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())

	//url := "http://192.168.1.237:5001"
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", rand.Int63())
		fmt.Println("UploadNft() s=", s, ",len =", len(s))
		//s = s[len(s)-13:]
		//NewTokenid := s
		//fmt.Println("UploadNft() NewTokenid=", NewTokenid)
	}

}

func TestDirss(t *testing.T) {

	nd, nerr := NewNftDb(sqldsnT)
	if nerr != nil {
		fmt.Printf("connect database err = %s\n", nerr)
	}
	defer nd.Close()
	meta := "/ipfs/QmVyVJTMQVbHRz8dr8RHrW4c1pgnspcM3Ee1pj9vae2oo8"
	err := contracts.SendSnftTrans(meta, "")
	if err != nil {
		fmt.Println("SetVoteSnftPeriod() save nftmeta info err=", err)
	}
	//img := "data:image/jpeg;base64,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"
	//_, serr := saveIpfsjpgImage("222", img)
	//if serr != nil {
	//	fmt.Println("SaveToIpfs() save collection image err=", serr)
	//}
}

func TestSlice(t *testing.T) {
	bb := "0x8000000000000000000000000000000000001234"
	bb = bb[len(bb)-2:]
}

func TestQueryOwnerSnftCollection(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	//SELECT a.Createaddr, a.name, a.img, count(a.name) FROM collects as a
	//		   JOIN (SELECT collectcreator, collections FROM nfts where ownaddr = "0x83c43f6f7bb4d8e429b21ff303a16b4c99a59b05" and deleted_at  is null ) as b ON a.Createaddr = b.collectcreator AND a.name = b.collections
	//			 GROUP BY a.createaddr, a.name, a.img
	//
	page, recount, err := nd.QueryOwnerSnftCollection("0x83c43f6f7bb4d8e429b21ff303a16b4c99a59b05", "Art", "0", "16")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(recount)
	t.Log(page)
}

func TestQueryOwnerSnftEth(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	//SELECT a.Createaddr, a.name, a.img, count(a.name) FROM collects as a
	//		   JOIN (SELECT collectcreator, collections FROM nfts where ownaddr = "0x83c43f6f7bb4d8e429b21ff303a16b4c99a59b05" and deleted_at  is null ) as b ON a.Createaddr = b.collectcreator AND a.name = b.collections
	//			 GROUP BY a.createaddr, a.name, a.img
	//
	serr := nd.SetPeriodEth("2633215263482")
	if serr != nil {
		fmt.Println(serr)
	}
}

func TestQueryOwnerSnftChipAmount(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	count, err := nd.QueryOwnerSnftChipAmount("0x83c43f6f7bb4d8e429b21ff303a16b4c99a59b05", "collect2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestQuerySnftByCollection(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	page, err := nd.QuerySnftByCollection("", "0xbe8c75133a7e4f29b7cdc15d4a45f7593a4f8898", "4-测试4", "0", "16")
	if err != nil {
		t.Fatal(err)
	}
	page, err = nd.QuerySnftByCollection("0x085ABc35ed85d26C2795b64C6fFb89B68aB1c479", "0xbe8c75133a7e4f29b7cdc15d4a45f7593a4f8898", "4-测试4", "0", "16")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(page)
}

func TestSubscribeEmails(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	err = nd.SetSubscribeEmail("user0", "www@ttest.com")
	if err != nil {
		t.Fatal(err)
	}
	err = nd.SetSubscribeEmail("user1", "www1@ttest.com")
	if err != nil {
		t.Fatal(err)
	}
	err = nd.SetSubscribeEmail("user1", "www2@ttest.com")
	if err != nil {
		t.Fatal(err)
	}
	err = nd.SetSubscribeEmail("user2", "www2@ttest.com")
	if err != nil {
		t.Fatal(err)
	}
	count, emails, err := nd.QuerySubscribeEmails("0", "10")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(count, emails)
	err = nd.DelSubscribeEmail("user2", "www2@ttest.com")
	if err != nil {
		t.Fatal(err)
	}
	count, emails, err = nd.QuerySubscribeEmails("0", "10")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(count, emails)
}

func TestCreateIndexs(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	err = nd.CreateIndexs()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTables(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	err = nd.CreateTables()
	if err != nil {
		t.Fatal(err)
	}
	sysinfo := SysInfos{}
	db := nd.db.Model(&SysInfos{}).Create(&sysinfo)
	if db.Error != nil {
		fmt.Println("TestCreateTables()->create() err=", db.Error)
		return
	}
}

func TestFmttest(t *testing.T) {
	str := fmt.Sprintf("%04x", 0)
	str = fmt.Sprintf("%04x", 16)
	str = fmt.Sprintf("%04x", 256)
	str = fmt.Sprintf("%04x", 4096)
	str = fmt.Sprintf("%04x", 65535)
	fmt.Println(str)
}

func TestUndeleted(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	sysinfo := SysInfos{}
	db := nd.db.Model(&SysInfos{}).Create(&sysinfo)
	if db.Error != nil {
		fmt.Println("TestCreateTables()->create() err=", db.Error)
		return
	}
	sysinfo = SysInfos{}
	db = nd.db.Model(&SysInfos{}).Create(&sysinfo)
	if db.Error != nil {
		fmt.Println("TestCreateTables()->create() err=", db.Error)
		return
	}
	db = nd.db.Model(&SysInfos{}).Where("id = ?", sysinfo.ID).Delete(&sysinfo)
	if db.Error != nil {
		fmt.Println("TestCreateTables()->create() err=", db.Error)
		return
	}
	sysinfo.DeletedAt = gorm.DeletedAt{}
	sysinfo.Snfttotal = 100
	db = nd.db.Unscoped().Model(&sysinfo).Where("id = ?", sysinfo.ID).Updates(&sysinfo )
	if db.Error != nil {
		fmt.Println("TestCreateTables()->create() err=", db.Error)
		return
	}
}
