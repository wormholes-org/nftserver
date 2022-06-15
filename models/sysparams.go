package models

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	initNFT1155Addr = "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5"
	initAdminAddr   = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"

	initNFT1155 = "0x53d76f1988B50674089e489B5ad1217AaC08CC85"
	initTrade   = "0x3dE836C28a578da26D846f27353640582761909f"

	initLowprice      = 1000000000
	initRoyaltylimit  = 50 * 100
	SysRoyaltylimit   = 50 * 100
	ZeroAddr          = "0x0000000000000000000000000000000000000000"
	genTokenIdRetry   = 20
	initCategories    = "art,music,domain_names,virtual_worlds,trading_cards,collectibles,sports,utility"
	LenName           = 60
	LenEmail          = 60
	LenPriceStr       = 9
	LowPrice          = 0
	ToolongAuciton    = 365
	HomePages         = "{\"announcement\":[\"m1\",\"m2\",\"m3\",\"m4\",\"m5\"],\"nft_loop\":[{\"contract\":\"\",\"tokenid\":\"\"}],\"collections\":[{\"creator\":\"\",\"name\":\"\"}],\"nfts\":[{\"contract\":\"\",\"tokenid\":\"\"}]}"
	DefAutoFlag       = "true"
	DefAutoSnft       = "false"
	DefUploadSize     = 100000000
	DefCatchTime      = 15
	DefNftloopcount   = 5
	DefSNftStartBlock = 1
	DefNftloopflush   = 15
	DefCollectcount   = 5
	DefCollectflush   = 15
	DefNftcount       = 5
	DefNftlush        = 15
	Deflanguage       = "中文"
	FlushTime         = time.Second * 10
	AesKey            = "CISzdrmfuTQvFJXpLySugjzTqorIMKSZ"
)

var (
	TradeAddr              string
	NFT1155Addr            string
	Weth9Addr              string
	AdminAddr              string
	BrowseNode             string
	EthersNode             string
	NftIpfsServerIP        string
	NftstIpfsServerPort    string
	EthersWsNode           string
	ImageDir               string
	AdminListPrv           string
	SuperAdminPrv          string
	SuperAdminAddr         string
	TradeAuthAddrPrv       string
	AdminMintPrv           string
	Lowprice               uint64
	RoyaltyLimit           int
	NFTUploadAuditRequired bool
	KYCUploadAuditRequired bool
	Authorize              string
	ExchangeOwer           string
	ExchangeName           string
	ExchangeBlocknumber    uint64
	AnnouncementRequired   bool
	TransferNFT            bool
	AutocommitSnft         bool
	ExchangerPrv           *ecdsa.PrivateKey
	ExchangerAddr          string
	ExchangerAuth          string
	DebugPort              string
	DebugAllowNft          string
	AllowNft               bool
	AllowUserMinit         bool
	UploadSize             uint64
)

type ExchangerAuthrize struct {
	Type          string `json:"type"`
	Version       int    `json:"version"`
	ExchangeOwner string `json:"exchanger_owner"`
	ExchangeName  string `json:"exchange_name"`
	To            string `json:"to"`
	BlockNumber   int    `json:"block_number"`
	Sig           string `json:"sig"`
}

type SysParamsRec struct {
	NFT1155addr    string `json:"nft1155addr" gorm:"type:char(42) ;comment:'nft1155合约地址'"`
	Adminaddr      string `json:"adminaddr" gorm:"type:char(42) ;comment:'管理员合约地址'"`
	Lowprice       uint64 `json:"lowprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'底价'"`
	Blocknumber    uint64 `json:"blocknumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'区块高度'"`
	Scannumber     uint64 `json:"scannumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'已扫描区块高度'"`
	Scansnftnumber uint64 `json:"scansnftnumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'已扫描snft区块高度'"`
	Royaltylimit   int    `json:"royaltylimit" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Signdata       string `json:"sig" gorm:"type:longtext ;comment:'签名数据'"`
	Homepage       string `json:"homepage" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'homepage数据'"`
	Exchangerinfo  string `json:"exchangerInfo" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'交易所信息数据'"`
	Icon           string `json:"icon" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'图形数据'"`
	Data           string `json:"data" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'轮播图、看板等数据'"`
	Categories     string `json:"categories" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'nft分类'"`
	Extend         string `json:"extend" gorm:"type:longtext ;comment:'扩展'"`
	Nftaudit       string `json:"nftaudit" gorm:"type:varchar(10) ;comment:'nft上传是否需要审核'"`
	Userkyc        string `json:"userkyc" gorm:"type:varchar(10) ;comment:'kyc是否需要审核'"`
	Deflanguage    string `json:"def_language" gorm:"type:varchar(50) CHARACTER SET utf8mb4 ;comment:'交易所默认语言'"`
	Restrictcode   string `json:"restrictcode" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'交易所国家代码限制'"`
	Autoflag       string `json:"autoflag" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'轮播图、看板等数据自动提取标志'"`
	Catchtime      int    `json:"catchtime" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'HomePage提取时间间隔'"`
	Nftloopcount   int    `json:"nftloopcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'轮播图选取的nft数'"`
	//Nftloopflush           int    `json:"nftloopflush" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'轮播图提取间隔'"`
	Collectcount int `json:"collectcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'热门合集选取数'"`
	//Collectflush           int    `json:"collectflush" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'热门合集提取间隔'"`
	Nftcount       int    `json:"nftcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'热门nft选取数'"`
	Exchangerprv   string `json:"exchangerprv" gorm:"type:longtext;COMMENT:'交易所私钥'"`
	Exchangerauth  string `json:"exchangerauth" gorm:"type:longtext;COMMENT:'交易所签名'"`
	Transfersnft   string `json:"transfersnft" gorm:"type:varchar(10) ;comment:'snft是否自动导入'"`
	Allownft       string `json:"allownft" gorm:"type:varchar(10) ;comment:'nft是否允许创建'"`
	Autocommitsnft string `json:"autocommitsnft" gorm:"type:varchar(10) ;comment:'snft自动注入eth'"`
	Allowusermint  string `json:"allowusermint" gorm:"type:varchar(10);comment:'是否允许用户铸币'"`
	Uploadsize     uint64 `json:"uploadsize" gorm:"type:bigint unsigned  DEFAULT 0;COMMENT:'上传nft限制大小'"`

	//Nftlush                int    `json:"nftlush" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'热门nft提取时间间隔'"`
}

type HomePageNft struct {
	Contract string `json:"contract"`
	Tokenid  string `json:"tokenid"`
}

type HomePageCollections struct {
	Creator string `json:"creator"`
	Name    string `json:"name"`
}

type HomePage struct {
	Announcement []string              `json:"announcement"`
	NftLoop      []HomePageNft         `json:"nft_loop"`
	Collections  []HomePageCollections `json:"collections"`
	Nfts         []HomePageNft         `json:"nfts"`
}

type SysParams struct {
	gorm.Model
	SysParamsRec
}

func (v SysParams) TableName() string {
	return "sysparams"
}

type SysParamsInfo struct {
	NFT1155addr    string `json:"nft1155addr"`
	Adminaddr      string `json:"adminaddr"`
	Lowprice       string `json:"lowprice"`
	Blocknumber    string `json:"blocknumber"`
	Scannumber     string `json:"scannumber"`
	Royaltylimit   string `json:"royaltylimit"`
	Homepage       string `json:"homepage"`
	Exchangerinfo  string `json:"exchangerinfo"`
	Icon           string `json:"icon"`
	Data           string `json:"data"`
	Categories     string `json:"categories"`
	Nftaudit       string `json:"nftaudit"`
	Userkyc        string `json:"userkyc"`
	Deflanguage    string `json:"def_language"`
	Restrictcode   string `json:"restrictcode"`
	Sig            string `json:"sig"`
	AutoFlag       string `json:"autoflag"`
	Catchtime      string `json:"catchtime"`
	Nftloopcount   string `json:"nftloopcount"`
	Collectcount   string `json:"collectcount"`
	Nftcount       string `json:"nftcount"`
	Announcement   string `json:"announcement"`
	TransferNFT    string `json:"transfernft"`
	ExchangerAddr  string `json:"exchangeraddr"`
	ExchangerAuth  string `json:"exchangerauth"`
	AutocommitSnft string `json:"autocommitsnft"`
	AllowNft       string `json:"allownft"`
	AllowUserMint  string `json:"allowusermint"`
	Uploadsize     string `json:"uploadsize"`
}

func (nft NftDb) QuerySysParams() (*SysParamsInfo, error) {
	var params SysParams
	err := nft.db.Last(&params)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			params = SysParams{}
			params.NFT1155addr = strings.ToLower(NFT1155Addr)
			params.Adminaddr = strings.ToLower(AdminAddr)
			params.Lowprice = initLowprice
			params.Royaltylimit = initRoyaltylimit
			params.Categories = initCategories
			//params.Blocknumber = contracts.GetCurrentBlockNumber()
			params.Blocknumber = ExchangeBlocknumber
			params.Scannumber = params.Blocknumber
			params.Scansnftnumber = DefSNftStartBlock
			params.Homepage = HomePages
			params.Deflanguage = Deflanguage
			params.Exchangerinfo = ExchangeName
			params.Autoflag = DefAutoFlag
			params.Catchtime = DefCatchTime
			params.Nftloopcount = DefNftloopcount
			//params.Nftloopflush = DefNftloopflush
			params.Collectcount = DefCollectcount
			//params.Collectflush = DefCollectflush
			params.Nftcount = DefNftcount
			params.Transfersnft = DefAutoSnft
			params.Autocommitsnft = DefAutoSnft
			params.Allownft = DefAutoSnft
			params.Allowusermint = DefAutoSnft
			params.Uploadsize = DefUploadSize
			//params.Exchangerprv = key
			//params.Nftlush = DefNftlush
			err = nft.db.Model(&SysParams{}).Create(&params)
			if err.Error != nil {
				fmt.Println("SetSysParams() create SysParams err= ", err.Error)
				return nil, err.Error
			}
		} else {
			fmt.Println("QuerySysParams() not find err=", err.Error)
			return nil, err.Error
		}
	}
	var paraminfo SysParamsInfo
	paraminfo.NFT1155addr = params.NFT1155addr
	paraminfo.Adminaddr = params.Adminaddr
	paraminfo.Lowprice = strconv.FormatUint(params.Lowprice, 10)
	paraminfo.Blocknumber = strconv.FormatUint(params.Blocknumber, 10)
	paraminfo.Scannumber = strconv.FormatUint(params.Scannumber, 10)
	paraminfo.Royaltylimit = strconv.Itoa(params.Royaltylimit)
	paraminfo.Homepage = params.Homepage
	paraminfo.Exchangerinfo = params.Exchangerinfo
	paraminfo.Icon = params.Icon
	paraminfo.Data = params.Data
	paraminfo.Categories = params.Categories
	paraminfo.Nftaudit = params.Nftaudit
	paraminfo.Deflanguage = params.Deflanguage
	paraminfo.Restrictcode = params.Restrictcode
	paraminfo.Userkyc = params.Userkyc
	paraminfo.AutoFlag = params.Autoflag
	paraminfo.Catchtime = strconv.Itoa(params.Catchtime)
	paraminfo.Nftloopcount = strconv.Itoa(params.Nftloopcount)
	//paraminfo.Nftloopflush 	= strconv.Itoa(params.Nftloopflush)
	paraminfo.Collectcount = strconv.Itoa(params.Collectcount)
	//paraminfo.Collectflush 	= strconv.Itoa(params.Collectflush)
	paraminfo.Nftcount = strconv.Itoa(params.Nftcount)
	paraminfo.Announcement = strconv.FormatBool(AnnouncementRequired)
	paraminfo.ExchangerAddr = ExchangerAddr
	ExchangerAuth = params.Exchangerauth
	paraminfo.AutocommitSnft = params.Autocommitsnft
	AutocommitSnft, _ = strconv.ParseBool(params.Autocommitsnft)
	paraminfo.TransferNFT = params.Transfersnft
	TransferNFT, _ = strconv.ParseBool(params.Transfersnft)
	paraminfo.AllowNft = params.Allownft
	AllowNft, _ = strconv.ParseBool(params.Allownft)
	paraminfo.AllowUserMint = params.Allowusermint
	AllowUserMinit, _ = strconv.ParseBool(params.Allowusermint)
	if params.Uploadsize != 0 {
		UploadSize = params.Uploadsize
		beego.BConfig.MaxMemory = int64(params.Uploadsize)
	} else {
		UploadSize = DefUploadSize
		beego.BConfig.MaxMemory = DefUploadSize
	}
	paraminfo.Uploadsize = strconv.FormatUint(params.Uploadsize, 10)
	//paraminfo.Nftlush 		= strconv.Itoa(params.Nftlush)
	return &paraminfo, err.Error
}

func (nft NftDb) GetSysParam(parameter string) (string, error) {
	parameter = strings.ToLower(parameter)
	var params SysParams
	var param string
	err := nft.db.Model(&params).Select(parameter).Last(&param)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			params = SysParams{}
			params.NFT1155addr = strings.ToLower(NFT1155Addr)
			params.Adminaddr = strings.ToLower(AdminAddr)
			params.Lowprice = initLowprice
			params.Royaltylimit = initRoyaltylimit
			params.Categories = initCategories
			params.Blocknumber = contracts.GetCurrentBlockNumber()
			params.Scannumber = params.Blocknumber
			params.Homepage = HomePages
			params.Deflanguage = Deflanguage
			params.Autoflag = DefAutoFlag
			params.Catchtime = DefCatchTime
			params.Nftloopcount = DefNftloopcount
			//params.Nftloopflush = DefNftloopflush
			params.Collectcount = DefCollectcount
			//params.Collectflush = DefCollectflush
			params.Nftcount = DefNftcount
			//params.Nftlush = DefNftlush
			err = nft.db.Model(&SysParams{}).Create(&params)
			if err.Error != nil {
				fmt.Println("SetSysParams() create SysParams err= ", err.Error)
				return "", err.Error
			}
			err = nft.db.Model(&params).Select(parameter).Last(&param)
			if err.Error != nil {
				fmt.Println("GetSysParams() create SysParams not find err= ", err.Error)
				return "", err.Error
			}
		} else {
			fmt.Println("GetSysParams() not find err=", err.Error)
			return "", err.Error
		}
	}
	return param, err.Error
}

func (nft NftDb) SetSysParams(param SysParamsInfo) error {
	var paramRec, updateP SysParams
	err := nft.db.Last(&paramRec)
	if err.Error != nil {
		if nft.db.Error == gorm.ErrRecordNotFound {
			updateP.NFT1155addr = NFT1155Addr
			updateP.Adminaddr = AdminAddr
			updateP.Lowprice = initLowprice
			updateP.Royaltylimit = initRoyaltylimit
			updateP.Categories = initCategories
			updateP.Homepage = HomePages
			//key, _ := ExchangeGenerate()
			//updateP.Exchangerprv = key
		} else {
			fmt.Println("QuerySysParams() not find err=", err.Error)
			return err.Error
		}
	} else {
		updateP.SysParamsRec = paramRec.SysParamsRec
		if param.NFT1155addr != "" {
			updateP.NFT1155addr = param.NFT1155addr
		}
		if param.Adminaddr != "" {
			updateP.Adminaddr = param.Adminaddr
		}
		if IsUint64DataValid(param.Lowprice) {
			low, _ := strconv.ParseUint(param.Lowprice, 10, 64)
			updateP.Lowprice = low
		}
		if param.Nftaudit != "" {
			updateP.Nftaudit = param.Nftaudit
			audit, err := strconv.ParseBool(updateP.Nftaudit)
			if err != nil {
				return errors.New("NftAudit input  error.")
			}
			NFTUploadAuditRequired = audit
		}
		if param.TransferNFT != "" {
			updateP.Transfersnft = param.TransferNFT
			audit, err := strconv.ParseBool(param.TransferNFT)
			if err != nil {
				return errors.New("TransferNFT input  error.")
			}
			TransferNFT = audit
		}
		if param.Userkyc != "" {
			updateP.Userkyc = param.Userkyc
			audit, err := strconv.ParseBool(updateP.Userkyc)
			if err != nil {
				return errors.New("Userkyc  input  error.")
			}
			KYCUploadAuditRequired = audit
		}
		if param.Announcement != "" {
			audit, err := strconv.ParseBool(param.Announcement)
			if err != nil {
				return errors.New("Announcement  input  error.")
			}
			AnnouncementRequired = audit
		}
		if param.AutocommitSnft != "" {
			updateP.Autocommitsnft = param.AutocommitSnft
			audit, err := strconv.ParseBool(param.AutocommitSnft)
			if err != nil {
				return errors.New("AutocommitSnft  input  error.")
			}
			AutocommitSnft = audit
		}
		if param.AllowNft != "" {
			updateP.Allownft = param.AllowNft
			audit, err := strconv.ParseBool(param.AllowNft)
			if err != nil {
				return errors.New("AllowNft  input  error.")
			}
			AllowNft = audit
		}
		if param.AllowUserMint != "" {
			updateP.Allowusermint = param.AllowUserMint
			audit, err := strconv.ParseBool(param.AllowUserMint)
			if err != nil {
				return errors.New("AllowUserMint  input  error.")
			}
			AllowUserMinit = audit
		}
		if param.Uploadsize != "" {
			low, _ := strconv.ParseUint(param.Uploadsize, 10, 64)
			updateP.Uploadsize = low
			UploadSize = low
			beego.BConfig.MaxMemory = int64(low)
		}
		if param.Deflanguage != "" {
			updateP.Deflanguage = param.Deflanguage
		}
		if param.Restrictcode != "" {
			updateP.Restrictcode = param.Restrictcode
		}
		if param.Homepage != "" {
			updateP.Homepage = param.Homepage
		}
		//if IsUint64DataValid(param.Scannumber) {
		//	updateP.Scannumber, _ = strconv.ParseUint(param.Scannumber, 10, 64)
		//}
		if IsIntDataValid(param.Royaltylimit) {
			updateP.Royaltylimit, _ = strconv.Atoi(param.Royaltylimit)
		}
		if param.Exchangerinfo != "" {
			updateP.Exchangerinfo = param.Exchangerinfo
		}
		if param.Deflanguage != "" {
			updateP.Deflanguage = param.Deflanguage
		}
		if param.Restrictcode != "" {
			updateP.Restrictcode = param.Restrictcode
		}
		if param.Icon != "" {
			updateP.Icon = param.Icon
		}
		if param.Data != "" {
			updateP.Data = param.Data
		}
		if param.Categories != "" {
			updateP.Categories = param.Categories
		}
		if param.Sig != "" {
			updateP.Signdata = param.Sig
		}
		if param.AutoFlag != "" {
			af := strings.ToLower(param.AutoFlag)
			if af == "true" || af == "false" {
				updateP.Autoflag = af
			} else {
				fmt.Println("SetSysParams() AutoFlag err= ")
				return errors.New("AutoFlag error.")
			}
		}
		if IsIntDataValid(param.Catchtime) {
			time, terr := strconv.Atoi(param.Catchtime)
			if terr != nil || time < 0 {
				fmt.Println("SetSysParams() Catchtime err= ")
				return errors.New("Catchtime error.")
			}
			updateP.Catchtime = time
		}
		if IsIntDataValid(param.Nftloopcount) {
			Nftloopcount, terr := strconv.Atoi(param.Nftloopcount)
			if terr != nil || Nftloopcount < 0 {
				fmt.Println("SetSysParams() Nftloopcount err= ")
				return errors.New("Nftloopcount error.")
			}
			updateP.Nftloopcount = Nftloopcount
		}
		//if IsIntDataValid(param.Nftloopflush) {
		//	updateP.Nftloopflush, _ = strconv.Atoi(param.Nftloopflush)
		//}
		if IsIntDataValid(param.Collectcount) {
			Collectcount, terr := strconv.Atoi(param.Collectcount)
			if terr != nil || Collectcount < 0 {
				fmt.Println("SetSysParams() Collectcount err= ")
				return errors.New("Collectcount error.")
			}
			updateP.Collectcount = Collectcount
		}
		//if IsIntDataValid(param.Collectflush) {
		//	updateP.Collectflush, _ = strconv.Atoi(param.Collectflush)
		//}
		if IsIntDataValid(param.Nftcount) {
			Nftcount, terr := strconv.Atoi(param.Nftcount)
			if terr != nil || Nftcount < 0 {
				fmt.Println("SetSysParams() Nftcount err= ")
				return errors.New("Nftcount error.")
			}
			updateP.Nftcount = Nftcount
		}
		if param.ExchangerAuth != "" {
			updateP.Exchangerauth = param.ExchangerAuth
			fmt.Println("ExchangerAuth=", ExchangerAuth)
			ExchangerAuth = param.ExchangerAuth
		}
		//if IsIntDataValid(param.Collectflush) {
		//	updateP.Nftlush, _ = strconv.Atoi(param.Nftlush)
		//}
	}
	err = nft.db.Model(&SysParams{}).Create(&updateP)
	if err.Error != nil {
		fmt.Println("SetSysParams() create SysParams err= ", err.Error)
		return err.Error
	}
	return nil
}

func (nft NftDb) SetAnnouncementParam(param string) error {

	param = strings.ToLower(param)
	announce, err := strconv.ParseBool(param)
	if err != nil {
		fmt.Println("SetAnnouncementParam()  err=", err)
		return errors.New("input announcement switch params  error ")
	}
	AnnouncementRequired = announce
	return nil
}

func (nft NftDb) SetExchageSig(exchange string) error {
	ExchangerAuth = exchange
	err := nft.db.Last(&SysParams{}).Update("exchangerauth", exchange)
	if err.Error != nil {
		fmt.Println("SetExchageSig() update exchangerauth err= ", err.Error)
		return err.Error
	}
	return nil
}

func (nft NftDb) TranSnft() error {
	if TransferNFT == true {
		return nil
	}
	err := nft.db.Last(&SysParams{}).Updates(map[string]interface{}{"blocknumber": 1, "scannumber": 1, "transfersnft": "true"})
	if err.Error != nil {
		fmt.Println("TranSnft() update Blocknumber err= ", err.Error)
		return err.Error
	}

	return nil
}

func (nft NftDb) GetExchageSig() (bool, error) {
	var params SysParams
	var auth bool
	err := nft.db.Last(&params)
	if err.Error != nil {
		fmt.Println("SetExchageSig() update exchangerauth err= ", err.Error)
		return false, err.Error
	}
	if params.Exchangerauth == "" {
		auth = false
	} else {
		auth = true
	}
	return auth, nil
}
func InitSysParams(Sqldsndb string) error {
	fmt.Println("InitSysParams() TradeAddr=", TradeAddr)
	fmt.Println("InitSysParams() NFT1155Addr=", NFT1155Addr)
	fmt.Println("InitSysParams() AdminAddr=", AdminAddr)
	fmt.Println("InitSysParams() EthersNode=", EthersNode)
	fmt.Println("InitSysParams() EthersWsNode=", EthersWsNode)
	fmt.Println("InitSysParams() ImageDir=", ImageDir)
	fmt.Println("InitSysParams() Weth9Addr=", Weth9Addr)
	//fmt.Println("InitSysParams() AdminListPrv=", AdminListPrv)
	//fmt.Println("InitSysParams() TradeAuthAddrPrv=", TradeAuthAddrPrv)
	//fmt.Println("InitSysParams() AdminMintPrv=", AdminMintPrv)
	var auth ExchangerAuthrize
	err := json.Unmarshal([]byte(Authorize), &auth)
	if err != nil {
		fmt.Printf("InitSysParams() ExchangerAuthrize= %s    Unmarshal err = %s\n", Authorize, err)
		return err
	}
	_, err = IsValidAddr(auth.ExchangeOwner+auth.To+strconv.Itoa(auth.BlockNumber), auth.Sig, auth.ExchangeOwner)
	if err != nil {
		fmt.Printf("InitSysParams() isValidAddr err = %s\n", err)
		return err
	}
	ExchangeOwer = strings.ToLower(auth.ExchangeOwner)
	ExchangeName = auth.ExchangeName
	fmt.Printf("InitSysParams() ExchangeOwer = %s\n", ExchangeOwer)
	ExchangeBlocknumber = uint64(auth.BlockNumber)
	SuperAdminPrv = DefSuperAdminPrv
	privateKey, err := crypto.HexToECDSA(SuperAdminPrv)
	if err != nil {
		fmt.Printf("InitSysParams() AdminListPrv err = %s\n", err)
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Printf("InitSysParams() publicKey err = %s\n", err)
		return errors.New("InitSysParams() publicKey err")
	}
	SuperAdminAddr = strings.ToLower(crypto.PubkeyToAddress(*publicKeyECDSA).String())
	//hexExchangePrv := crypto.FromECDSA(ExchangerPrv)
	//hexExchangePrvStr := hexutil.Encode(hexExchangePrv)[2:]
	ExchangerAddr = SuperAdminAddr
	contracts.SetSysParams(EthersNode, BrowseNode, EthersWsNode, Weth9Addr, TradeAddr,
		NFT1155Addr, AdminAddr, AdminListPrv, TradeAuthAddrPrv, AdminMintPrv, SuperAdminPrv, ExchangeOwer)
	nd, err := NewNftDb(Sqldsndb)
	defer nd.Close()

	if err != nil {
		fmt.Printf("InitSysParams() connect database err = %s\n", err)
		return err
	}
	params, err := nd.QuerySysParams()
	if err != nil {
		fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
		return err
	} else {
		Lowprice, _ = strconv.ParseUint(params.Lowprice, 10, 64)
		RoyaltyLimit, _ = strconv.Atoi(params.Royaltylimit)
		dbAudit, err := strconv.ParseBool(params.Nftaudit)
		if err != nil {
			newParam := SysParamsInfo{
				Nftaudit: fmt.Sprintf("%t", NFTUploadAuditRequired),
			}
			nd.SetSysParams(newParam)
		} else {
			NFTUploadAuditRequired = dbAudit
		}
		userAudit, err := strconv.ParseBool(params.Userkyc)
		if err != nil {
			newParam := SysParamsInfo{
				Userkyc: fmt.Sprintf("%t", KYCUploadAuditRequired),
			}
			nd.SetSysParams(newParam)
		} else {
			KYCUploadAuditRequired = userAudit
		}
		go HomePageFlash(Sqldsndb)
	}
	if DebugAllowNft != "" {
		err = nd.TranSnft()
		if err != nil {
			fmt.Printf("InitSysParams() TranSnft() err = %v\n", err)
			return err
		}
		_, err = nd.QuerySysParams()
		if err != nil {
			fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
			return err
		}
	}
	go AutoPeriodEth(Sqldsndb)
	//hexExchangePrv := crypto.FromECDSA(ExchangerPrv)
	//hexExchangePrvStr := hexutil.Encode(hexExchangePrv)[2:]
	//contracts.SetSysParams(BrowseNode, EthersWsNode, Weth9Addr, TradeAddr,
	//	NFT1155Addr, AdminAddr, AdminListPrv, TradeAuthAddrPrv, AdminMintPrv, hexExchangePrvStr, ExchangeOwer)
	err = nd.SetExchangerAdmin(ExchangeOwer)
	if err != nil {
		fmt.Printf("InitSysParams() SetExchangerAdmin() ExchangeOwer err = %s\n", err)
		return err
	}
	err = nd.SetExchangerAdmin(SuperAdminAddr)
	if err != nil {
		fmt.Printf("InitSysParams() SetExchangerAdmin()  SuperAdminAddr err = %s\n", err)
		return err
	}
	err = nd.DefaultCountry()
	if err != nil {
		fmt.Printf("InitSysParams() DefaultCountry()  err = %s\n", err)
		return err
	}
	sysInfo := SysInfos{}
	dberr := nd.db.Model(&SysInfos{}).Last(&sysInfo)
	if dberr.Error != nil {
		if dberr.Error != gorm.ErrRecordNotFound {
			log.Println("InitSysParams() SysInfos err=", dberr)
			return ErrCollectionNotExist
		}
		dberr = nd.db.Model(&SysInfos{}).Create(&sysInfo)
		if dberr.Error != nil {
			log.Println("InitSysParams() SysInfos create err=", dberr)
			return ErrCollectionNotExist
		}
	}
	_, err = nd.QueryHomePage(true)
	return nil
}

func ScanNftLoop(sqldsn string, interval time.Duration, stop chan struct{}) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("HomePageFlash() connect database err = %s\n", err)
				continue
			}
			params, err := nd.QuerySysParams()
			if err != nil {
				fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
				continue
			}
			if params.AutoFlag == "true" {

			} else {

			}
			nd.Close()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func ScanCollection(sqldsn string, interval time.Duration, stop chan struct{}) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("HomePageFlash() connect database err = %s\n", err)
				continue
			}
			params, err := nd.QuerySysParams()
			if err != nil {
				fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
				continue
			}
			if params.AutoFlag == "true" {

			} else {

			}
			nd.Close()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func ScanNfts(sqldsn string, interval time.Duration, stop chan struct{}) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("HomePageFlash() connect database err = %s\n", err)
				continue
			}
			params, err := nd.QuerySysParams()
			if err != nil {
				fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
				continue
			}
			if params.AutoFlag == "true" {

			} else {

			}
			nd.Close()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

type HotTrans struct {
	Contract string `json:"contract" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Tokenid  string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
}

func ScanLoop(sqldsn string, interval int, stop chan struct{}, stoped chan struct{}) {
	ticker := time.NewTicker(time.Minute * 1)
	scan := time.NewTicker(time.Duration(interval) * time.Minute)

	for {
		select {
		case <-scan.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				log.Printf("ScanLoop() connect database err = %s\n", err)
				continue
			}
			params, err := nd.QuerySysParams()
			if err != nil {
				nd.Close()
				log.Printf("ScanLoop() QuerySysParams() err = %s\n", err)
				continue
			}
			var hp, homepage HomePage
			err = json.Unmarshal([]byte(params.Homepage), &hp)
			if err != nil {
				log.Println("ScanLoop() Unmarshal() err =", err)
				continue
			}
			homepage.Announcement = hp.Announcement
			trans := make([]HotTrans, 0, 20)
			limit, _ := strconv.Atoi(params.Nftcount)
			//dberr := nd.db.Group("tokenid as tc").Order("tc desc").Limit(limit).Find(&trans)
			sql := "SELECT  contract, tokenid, count(tokenid) as tc from trans group by contract, tokenid order by tc desc limit "
			sql = sql + params.Nftcount
			dberr := nd.db.Raw(sql).Scan(&trans)
			if dberr.Error != nil {
				if dberr.Error != gorm.ErrRecordNotFound {
					nd.Close()
					log.Printf("ScanLoop() Find(&trans) err = %s\n", dberr.Error)
					continue
				}
			} else {
				if len(trans) > 0 {
					for _, tran := range trans {
						var homenft HomePageNft
						homenft.Contract = tran.Contract
						homenft.Tokenid = tran.Tokenid
						homepage.Nfts = append(homepage.Nfts, homenft)
					}
				} else {
					homepage.Nfts = []HomePageNft{{"", ""}}
				}
			}
			collects := make([]Collects, 0, 20)
			limit, _ = strconv.Atoi(params.Nftcount)
			dberr = nd.db.Order("transcnt desc").Limit(limit).Find(&collects)
			if dberr.Error != nil {
				if dberr.Error != gorm.ErrRecordNotFound {
					nd.Close()
					log.Printf("ScanLoop() Find(&collects) err = %s\n", dberr.Error)
					continue
				}
			} else {
				if len(collects) > 0 {
					for _, collect := range collects {
						var hcollect HomePageCollections
						hcollect.Creator = collect.Createaddr
						hcollect.Name = collect.Name
						homepage.Collections = append(homepage.Collections, hcollect)
					}
				} else {
					homepage.Collections = []HomePageCollections{{"", ""}}
				}
			}
			var recCount int64
			countSql := `select max(id) from nfts`
			dberr = nd.db.Raw(countSql).Scan(&recCount)
			//dberr = nd.db.Model(Nfts{}).Count(&recCount)
			if dberr.Error != nil {
				//if dberr.Error != gorm.ErrRecordNotFound {
				//	nd.Close()
				//	log.Printf("ScanLoop() Count(&recCount) err = %s\n", dberr.Error)
				//	continue
				//}
				log.Printf("ScanLoop() Count(&recCount) err = %s\n", dberr.Error)
				recCount = 0
			}
			if recCount != 0 {
				rand.Seed(time.Now().UnixNano())
				limit, _ = strconv.Atoi(params.Nftloopcount)
				scaned := make(map[int64]bool)
				for i := 0; i < limit && int64(i) < recCount; i++ {
					index := rand.Int63() % recCount
					if index == 0 {
						index = 1
					}
					flag := scaned[index]
					if flag {
						i = i - 1
						continue
					}
					scaned[index] = true
					var nftRec Nfts
					dberr := nd.db.Where("id = ?", index).First(&nftRec)
					if dberr.Error != nil {
						nd.Close()
						log.Println("ScanLoop() First(&nftRec) err = ", dberr.Error)
						continue
					}
					var hpnft HomePageNft
					hpnft.Contract = nftRec.Contract
					hpnft.Tokenid = nftRec.Tokenid
					homepage.NftLoop = append(homepage.NftLoop, hpnft)
				}
			} else {
				homepage.NftLoop = []HomePageNft{{"", ""}}
			}
			homestr, err := json.Marshal(&homepage)
			if err != nil {
				nd.Close()
				log.Println("ScanLoop() Marshal(&homepage) err = ", err)
				continue
			}
			newParam := SysParamsInfo{
				Homepage: string(homestr),
			}
			nd.SetSysParams(newParam)
			HomePageCatchs.HomePageFlashLock()
			HomePageCatchs.HomePageFlashFlag = true
			HomePageCatchs.HomePageFlashUnLock()
			nd.Close()
		case <-ticker.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				log.Printf("ScanLoop() connect database err = %s\n", err)
				continue
			}
			params, err := nd.QuerySysParams()
			nd.Close()
			if err != nil {
				log.Printf("ScanLoop() QuerySysParams() err = %s\n", err)
				continue
			}
			t, _ := strconv.Atoi(params.Catchtime)
			if interval != t {
				interval = t
				scan.Reset(time.Duration(interval) * time.Minute)
			}
		case <-stop:
			ticker.Stop()
			scan.Stop()
			stoped <- struct{}{}
			return
		}
	}
}

func HomePageFlash(sqldsn string) {
	autoFlag := "false"
	ScanStop := make(chan struct{})
	ScanStoped := make(chan struct{})
	ScanAutoFlag := time.NewTicker(FlushTime)
	for {
		select {
		case <-ScanAutoFlag.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("HomePageFlash() connect database err = %s\n", err)
				continue
			}
			params, err := nd.QuerySysParams()
			nd.Close()
			if err != nil {
				fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
				continue
			}
			if params.AutoFlag == "true" {
				if autoFlag == "false" {
					autoFlag = params.AutoFlag
					catchtime, _ := strconv.Atoi(params.Catchtime)
					go ScanLoop(sqldsn, catchtime, ScanStop, ScanStoped)
				}
			} else {
				if autoFlag != "false" {
					autoFlag = "false"
					ScanStop <- struct{}{}
					<-ScanStoped
				}
			}
		}
	}
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
