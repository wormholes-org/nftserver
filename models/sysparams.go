package models

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/crypto"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	initNFT1155Addr = "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5"
	initAdminAddr   = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"

	initNFT1155 = "0x53d76f1988B50674089e489B5ad1217AaC08CC85"
	initTrade   = "0x3dE836C28a578da26D846f27353640582761909f"

	initLowprice     = 1000000000
	initRoyaltylimit = 50 * 100
	SysRoyaltylimit  = 50 * 100
	ZeroAddr         = "0x0000000000000000000000000000000000000000"
	genTokenIdRetry  = 20
	initCategories   = "art,music,domain_names,virtual_worlds,trading_cards,collectibles,sports,utility"
	LenName          = 60
	LenEmail         = 60
	LenLink          = 2000
	LenPriceStr      = 9
	LowPrice         = 0
	ToolongAuciton   = 365
	//HomePages         = "{\"announcement\":[\"m1\",\"m2\",\"m3\",\"m4\",\"m5\"],\"nft_loop\":[{\"contract\":\"\",\"tokenid\":\"\"}],\"collections\":[{\"creator\":\"\",\"name\":\"\"}],\"nfts\":[{\"contract\":\"\",\"tokenid\":\"\"}]}"
	HomePages         = "{\"announcement\":[\"m1\",\"m2\",\"m3\",\"m4\",\"m5\"],\"nft_loop\":[{\"contract\":\"\",\"tokenid\":\"\"},{\"contract\":\"\",\"tokenid\":\"\"},{\"contract\":\"\",\"tokenid\":\"\"}],\"collections\":[{\"creator\":\"\",\"name\":\"\"},{\"creator\":\"\",\"name\":\"\"},{\"creator\":\"\",\"name\":\"\"},{\"creator\":\"\",\"name\":\"\"}],\"nfts\":[{\"contract\":\"\",\"tokenid\":\"\"},{\"contract\":\"\",\"tokenid\":\"\"},{\"contract\":\"\",\"tokenid\":\"\"},{\"contract\":\"\",\"tokenid\":\"\"}]}"
	DefExchangeLink   = "{\"github\":\"https://github.com/wormholes-org\",\"discord\":\"https://discord.com/invite/AbmTrrAmuN\",\"twitter\":\"https://twitter.com/WormholesChain\"}"
	DefAutoFlag       = "true"
	DefAutoSnft       = "false"
	DefAudit          = "false"
	DefUploadSize     = 100 * 1024 * 1024
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
	DefPartlogo       = "wormholes"
	DefCaptchaNum     = 16
	DefaultWormBlack  = "/ipfs/QmSQf4rm7C2riGffS6YFxrPgiVBAaSwLMfdhubPo1jP2A5/worm_black.png"
	DefaultWormBlue   = "/ipfs/QmSQf4rm7C2riGffS6YFxrPgiVBAaSwLMfdhubPo1jP2A5/worm_blue.jpg"
	DefaultCollection = "mycollection"
)

var (
	TradeAddr               string
	NFT1155Addr             string
	Weth9Addr               string
	AdminAddr               string
	BrowseNode              string
	EthersNode              string
	NftIpfsServerIP         string
	NftstIpfsServerPort     string
	EthersWsNode            string
	ImageDir                string
	AdminListPrv            string
	SuperAdminPrv           string
	SuperAdminAddr          string
	TradeAuthAddrPrv        string
	AdminMintPrv            string
	Lowprice                uint64
	RoyaltyLimit            int
	NFTUploadAuditRequired  bool
	KYCUploadAuditRequired  bool
	Authorize               string
	ExchangeOwer            string
	ExchangeName            string
	ExchangeBlocknumber     uint64
	AnnouncementRequired    bool
	TransferSNFT            bool
	AutocommitSnft          bool
	ExchangerPrv            *ecdsa.PrivateKey
	ExchangerAddr           string
	ExchangerAuth           string
	DebugPort               string
	DebugAllowNft           string
	AllowNft                bool
	AllowUserMinit          bool
	UploadSize              uint64
	Backupipfs              bool
	BackupIpfsUrl           string
	QueryRedisCatchSvr      string
	QueryRedisSvrPasswd     string
	MainRedisCatchSvr       string
	MainRedisCatchSvrPasswd string
	DefaultCaptcha          string
	DefaultMask             string
	DefaultMaskFrame        string
	DefaultCaptchaNum       int
	LimitWritesDatabase     bool
	NftScanServer           string
	AgentExchangePrv        string
	LimitTotalSize          bool
	LimitFileSize           string
	TransSnft               bool
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
	NFT1155addr    string `json:"nft1155addr" gorm:"type:char(42) ;comment:'nft1155 contract address'"`
	Adminaddr      string `json:"adminaddr" gorm:"type:char(42) ;comment:'Administrator contract address'"`
	Lowprice       uint64 `json:"lowprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Reserve price'"`
	Blocknumber    uint64 `json:"blocknumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'block height'"`
	Scannumber     uint64 `json:"scannumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Scanned block height'"`
	Scansnftnumber uint64 `json:"scansnftnumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Scanned snft block height'"`
	Savedsnft      string `json:"snft" gorm:"type:char(42) ;comment:'snft backed up to ipfs'"`
	Royaltylimit   int    `json:"royaltylimit" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	Signdata       string `json:"sig" gorm:"type:longtext ;comment:'sign data'"`
	Homepage       string `json:"homepage" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'homepage data'"`
	//Exchangerinfo  string `json:"exchangerInfo" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'Exchange information data'"`
	//Icon           string `json:"icon" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'picture data'"`
	//Data           string `json:"data" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'slideshow, Kanban, etc.'"`
	//Categories     string `json:"categories" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'nft category'"`
	Extend string `json:"extend" gorm:"type:longtext ;comment:'extend'"`
	//Nftaudit       string `json:"nftaudit" gorm:"type:varchar(10) ;comment:'Does nft upload need to be reviewed?'"`
	//Userkyc        string `json:"userkyc" gorm:"type:varchar(10) ;comment:'Does kyc need to be audited?'"`
	//Deflanguage    string `json:"def_language" gorm:"type:varchar(50) CHARACTER SET utf8mb4 ;comment:'Exchange default language'"`
	//Restrictcode  string `json:"restrictcode" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'Exchange country code restrictions'"`
	Autoflag      string `json:"autoflag" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'Automatically extract flag from data such as carousel and kanban'"`
	Catchtime     int    `json:"catchtime" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'HomePage fetch interval'"`
	Nftloopcount  int    `json:"nftloopcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'The number of nfts selected by the carousel'"`
	Collectcount  int    `json:"collectcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'Popular Collection Picks'"`
	Nftcount      int    `json:"nftcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'Popular nft Picks'"`
	Exchangerprv  string `json:"exchangerprv" gorm:"type:longtext;COMMENT:'Exchange private key'"`
	Exchangerauth string `json:"exchangerauth" gorm:"type:longtext;COMMENT:'Exchange signature'"`
	//Transfersnft   string `json:"transfersnft" gorm:"type:varchar(10) ;comment:'Whether snft is automatically imported'"`
	//Allownft       string `json:"allownft" gorm:"type:varchar(10) ;comment:'Does nft allow to create'"`
	//Autocommitsnft string `json:"autocommitsnft" gorm:"type:varchar(10) ;comment:'snft automatically injects chain'"`
	//Allowusermint  string `json:"allowusermint" gorm:"type:varchar(10);comment:'Whether to allow users to mint coins'"`
	//Uploadsize     uint64 `json:"uploadsize" gorm:"type:bigint unsigned  DEFAULT 0;COMMENT:'upload nft limit size'"`
	//Backupipfs string `json:"backupipfs" gorm:"type:varchar(10) ;comment:'whether to backup to ipfs'"`

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
	Backupipfs     string `json:"backupipfs"`
}

type ExchangeinfoRec struct {
	Exchangerinfo  string `json:"exchangerInfo" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'Exchange information data'"`
	Icon           string `json:"icon" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'picture data'"`
	Data           string `json:"data" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'slideshow, Kanban, etc.'"`
	Categories     string `json:"categories" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'nft category'"`
	Restrictcode   string `json:"restrictcode" gorm:"type:longtext CHARACTER SET utf8mb4 ;comment:'Exchange country code restrictions'"`
	Extend         string `json:"extend" gorm:"type:longtext ;comment:'extend'"`
	Nftaudit       string `json:"nftaudit" gorm:"type:varchar(10) ;comment:'Does nft upload need to be reviewed?'"`
	Userkyc        string `json:"userkyc" gorm:"type:varchar(10) ;comment:'Does kyc need to be audited?'"`
	Deflanguage    string `json:"def_language" gorm:"type:varchar(50) CHARACTER SET utf8mb4 ;comment:'Exchange default language'"`
	Transfersnft   string `json:"transfersnft" gorm:"type:varchar(10) ;comment:'Whether snft is automatically imported'"`
	Allownft       string `json:"allownft" gorm:"type:varchar(10) ;comment:'Does nft allow to create'"`
	Autocommitsnft string `json:"autocommitsnft" gorm:"type:varchar(10) ;comment:'snft automatically injects chain'"`
	Allowusermint  string `json:"allowusermint" gorm:"type:varchar(10);comment:'Whether to allow users to mint coins'"`
	Uploadsize     uint64 `json:"uploadsize" gorm:"type:bigint unsigned  DEFAULT 0;COMMENT:'upload nft limit size'"`
	Backupipfs     string `json:"backupipfs" gorm:"type:varchar(10) ;comment:'whether to backup to ipfs'"`
	Announcements  string `json:"announcements" gorm:"type:varchar(10) ;comment:'announcement whether to open'"`
	Link           string `json:"link" gorm:"type:longtext NOT NULL;comment:'exchange social account'"`
	Partnerslogo   string `json:"partnerslogo" gorm:"type:longtext NOT NULL;comment:'exchange partnerslogo'"`
	Desc           string `json:"desc" gorm:"type:longtext ;comment:'exchange desc'"`
	Totalsize      uint64 `json:"totalsize" gorm:"type:bigint unsigned  DEFAULT 0;COMMENT:'upload nft total size'"`
	Limitsize      uint64 `json:"limitsize" gorm:"type:bigint unsigned  DEFAULT 0;COMMENT:'upload nft limit amount size'"`
}

type Exchangeinfos struct {
	gorm.Model
	ExchangeinfoRec
}

func (v Exchangeinfos) TableName() string {
	return "exchangeinfo"
}

type ExchangeInfo struct {
	Adminaddr      string   `json:"adminaddr"`
	Exchangerinfo  string   `json:"exchangerinfo"`
	Icon           string   `json:"icon"`
	Data           string   `json:"data"`
	Categories     string   `json:"categories"`
	Nftaudit       string   `json:"nftaudit"`
	Userkyc        string   `json:"userkyc"`
	Restrictcode   string   `json:"restrictcode"`
	TransferNFT    string   `json:"transfernft"`
	AllowNft       string   `json:"allownft"`
	AutocommitSnft string   `json:"autocommitsnft"`
	AllowUserMint  string   `json:"allowusermint"`
	Uploadsize     string   `json:"uploadsize"`
	Backupipfs     string   `json:"backupipfs"`
	Link           string   `json:"link"`
	Partnerslogo   string   `json:"partnerslogo"`
	ExchangerAddr  string   `json:"exchangeraddr"`
	Announcement   string   `json:"announcement"`
	Deflanguage    string   `json:"def_language"`
	ExtendLogo     []string `json:"extend_logo"`
	Homepage       string   `json:"homepage"`
	Desc           string   `json:"desc"`
	AutoFlag       string   `json:"autoflag"`
}

func (nft NftDb) QuerySysParams() (*SysParamsInfo, error) {
	var params SysParams
	var paraminfo SysParamsInfo
	err := nft.db.Last(&params)
	log.Println(params)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			GetRedisCatch().SetDirtyFlag(AllDirty)
			params = SysParams{}
			params.NFT1155addr = strings.ToLower(NFT1155Addr)
			params.Adminaddr = strings.ToLower(AdminAddr)
			params.Lowprice = initLowprice
			params.Royaltylimit = initRoyaltylimit
			//params.Categories = initCategories
			//params.Blocknumber = contracts.GetCurrentBlockNumber()
			params.Blocknumber = ExchangeBlocknumber
			params.Scannumber = params.Blocknumber
			params.Scansnftnumber = DefSNftStartBlock
			params.Homepage = HomePages
			//params.Deflanguage = Deflanguage
			//params.Exchangerinfo = ExchangeName
			params.Autoflag = DefAutoFlag
			params.Catchtime = DefCatchTime
			params.Nftloopcount = DefNftloopcount
			//params.Nftloopflush = DefNftloopflush
			params.Collectcount = DefCollectcount
			//params.Collectflush = DefCollectflush
			params.Nftcount = DefNftcount
			//params.Transfersnft = DefAutoSnft
			//params.Autocommitsnft = DefAutoSnft
			//params.Allownft = DefAutoSnft
			//params.Allowusermint = DefAutoSnft
			//params.Uploadsize = DefUploadSize
			//params.Backupipfs = DefAutoSnft
			//params.Exchangerprv = key
			//params.Nftlush = DefNftlush
			err = nft.db.Model(&SysParams{}).Create(&params)
			if err.Error != nil {
				fmt.Println("SetSysParams() create SysParams err= ", err.Error)
				return nil, ErrDataBase
			}
		} else {
			fmt.Println("QuerySysParams() not find err=", err.Error)
			return nil, ErrDataBase
		}
	}
	//err = nft.db.Model(&Exchangeinfos{}).Select("icon", "exchangerInfo", "data", "categories", "restrictcode",
	//	"extend", "nftaudit", "userkyc", "transfersnft", "allownft", "allownft", "autocommitsnft", "allowusermint", "uploadsize",
	//	"announcements", "link", "partnerslogo").First(&params)
	var exchangeinfo Exchangeinfos
	err = nft.db.Last(&exchangeinfo)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			exchange := Exchangeinfos{}
			exchange.Link = DefExchangeLink
			exchange.Categories = initCategories
			exchange.Userkyc = DefAutoSnft
			exchange.Nftaudit = DefAutoSnft
			exchange.Exchangerinfo = ExchangeName
			exchange.Transfersnft = DefAutoSnft
			exchange.Autocommitsnft = DefAutoSnft
			exchange.Allownft = DefAutoSnft
			exchange.Allowusermint = DefAutoSnft
			exchange.Uploadsize = DefUploadSize
			exchange.Backupipfs = DefAutoSnft
			exchange.Partnerslogo = DefPartlogo
			err = nft.db.Model(&Exchangeinfos{}).Create(&exchange)
			if err.Error != nil {
				fmt.Println("QuerySysParams() create Exchangeinfos err= ", err.Error)
				return nil, ErrDataBase
			}
		} else {
			fmt.Println("QuerySysParams() not find  Exchangeinfos err=", err.Error)
			return nil, ErrDataBase
		}
	}
	//if exchangeinfo.Partnerslogo == "" {
	//	err = nft.db.Model(&Exchangeinfos{}).Last(&exchangeinfo).Update("partnerslogo", DefPartlogo)
	//	if err.Error != nil {
	//		fmt.Println("QuerySysParams() update partnerslogo err= ", err.Error)
	//		return nil, ErrDataBase
	//	}
	//}
	log.Println("params =", params)
	paraminfo.NFT1155addr = params.NFT1155addr
	paraminfo.Adminaddr = params.Adminaddr
	paraminfo.Lowprice = strconv.FormatUint(params.Lowprice, 10)
	paraminfo.Blocknumber = strconv.FormatUint(params.Blocknumber, 10)
	paraminfo.Scannumber = strconv.FormatUint(params.Scannumber, 10)
	paraminfo.Royaltylimit = strconv.Itoa(params.Royaltylimit)
	paraminfo.Homepage = params.Homepage
	paraminfo.Exchangerinfo = exchangeinfo.Exchangerinfo
	paraminfo.Icon = exchangeinfo.Icon
	paraminfo.Data = exchangeinfo.Data
	paraminfo.Categories = exchangeinfo.Categories
	paraminfo.Nftaudit = exchangeinfo.Nftaudit
	paraminfo.Deflanguage = exchangeinfo.Deflanguage
	paraminfo.Restrictcode = exchangeinfo.Restrictcode
	paraminfo.Userkyc = exchangeinfo.Userkyc
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
	paraminfo.AutocommitSnft = exchangeinfo.Autocommitsnft
	AutocommitSnft, _ = strconv.ParseBool(exchangeinfo.Autocommitsnft)
	paraminfo.TransferNFT = exchangeinfo.Transfersnft
	TransferSNFT, _ = strconv.ParseBool(exchangeinfo.Transfersnft)
	paraminfo.AllowNft = exchangeinfo.Allownft
	AllowNft, _ = strconv.ParseBool(exchangeinfo.Allownft)
	paraminfo.AllowUserMint = exchangeinfo.Allowusermint
	AllowUserMinit, _ = strconv.ParseBool(exchangeinfo.Allowusermint)
	paraminfo.Backupipfs = exchangeinfo.Backupipfs
	Backupipfs, _ = strconv.ParseBool(exchangeinfo.Backupipfs)
	TransSnft, _ = strconv.ParseBool(exchangeinfo.Transfersnft)
	UploadSize = exchangeinfo.Uploadsize
	beego.BConfig.MaxUploadSize = int64(exchangeinfo.Uploadsize)
	paraminfo.Uploadsize = strconv.FormatUint(exchangeinfo.Uploadsize, 10)
	return &paraminfo, nil
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
			//params.Categories = initCategories
			params.Blocknumber = contracts.GetCurrentBlockNumber()
			params.Scannumber = params.Blocknumber
			params.Homepage = HomePages
			//params.Deflanguage = Deflanguage
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
			//updateP.Categories = initCategories
			updateP.Homepage = HomePages
			//key, _ := ExchangeGenerate()
			//updateP.Exchangerprv = key
		} else {
			fmt.Println("QuerySysParams() not find err=", err.Error)
			return ErrDataBase
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
		//if param.Nftaudit != "" {
		//	updateP.Nftaudit = param.Nftaudit
		//	audit, err := strconv.ParseBool(updateP.Nftaudit)
		//	if err != nil {
		//		log.Println("NftAudit input  error.")
		//		return ErrData
		//	}
		//	NFTUploadAuditRequired = audit
		//}
		//if param.TransferNFT != "" {
		//	updateP.Transfersnft = param.TransferNFT
		//	audit, err := strconv.ParseBool(param.TransferNFT)
		//	if err != nil {
		//		log.Println("TransferNFT input  error.")
		//		return ErrData
		//	}
		//	log.Println("TransferNFT =", param.TransferNFT)
		//	TransferSNFT = audit
		//}
		//if param.Userkyc != "" {
		//	updateP.Userkyc = param.Userkyc
		//	audit, err := strconv.ParseBool(updateP.Userkyc)
		//	if err != nil {
		//		log.Println("Userkyc input  error.")
		//		return ErrData
		//
		//	}
		//	KYCUploadAuditRequired = audit
		//}
		if param.Announcement != "" {
			audit, err := strconv.ParseBool(param.Announcement)
			if err != nil {
				log.Println("Announcement input  error.")
				return ErrData
			}
			AnnouncementRequired = audit
		}
		//if param.AutocommitSnft != "" {
		//	updateP.Autocommitsnft = param.AutocommitSnft
		//	audit, err := strconv.ParseBool(param.AutocommitSnft)
		//	if err != nil {
		//		log.Println("AutocommitSnft input  error.")
		//		return ErrData
		//	}
		//	AutocommitSnft = audit
		//}
		//if param.AllowNft != "" {
		//	updateP.Allownft = param.AllowNft
		//	audit, err := strconv.ParseBool(param.AllowNft)
		//	if err != nil {
		//		log.Println("AllowNft input  error.")
		//		return ErrData
		//	}
		//	AllowNft = audit
		//}
		//if param.AllowUserMint != "" {
		//	updateP.Allowusermint = param.AllowUserMint
		//	audit, err := strconv.ParseBool(param.AllowUserMint)
		//	if err != nil {
		//		log.Println("AllowUserMint input  error.")
		//		return ErrData
		//	}
		//	AllowUserMinit = audit
		//}
		//if param.Backupipfs != "" {
		//	updateP.Backupipfs = param.Backupipfs
		//	audit, err := strconv.ParseBool(param.Backupipfs)
		//	if err != nil {
		//		log.Println("Backupipfs input  error.")
		//		return ErrData
		//	}
		//	Backupipfs = audit
		////}
		//if param.Uploadsize != "" {
		//	low, _ := strconv.ParseUint(param.Uploadsize, 10, 64)
		//	updateP.Uploadsize = low
		//	UploadSize = low
		//	beego.BConfig.MaxMemory = int64(low)
		//}
		//if param.Deflanguage != "" {
		//	updateP.Deflanguage = param.Deflanguage
		//}
		//if param.Restrictcode != "" {
		//	updateP.Restrictcode = param.Restrictcode
		//}
		if param.Homepage != "" {
			updateP.Homepage = param.Homepage
		}
		//if IsUint64DataValid(param.Scannumber) {
		//	updateP.Scannumber, _ = strconv.ParseUint(param.Scannumber, 10, 64)
		//}
		if IsIntDataValid(param.Royaltylimit) {
			updateP.Royaltylimit, _ = strconv.Atoi(param.Royaltylimit)
		}
		//if param.Exchangerinfo != "" {
		//	updateP.Exchangerinfo = param.Exchangerinfo
		//}
		//if param.Deflanguage != "" {
		//	updateP.Deflanguage = param.Deflanguage
		//}
		//if param.Restrictcode != "" {
		//	updateP.Restrictcode = param.Restrictcode
		//}
		//if param.Icon != "" {
		//	updateP.Icon = param.Icon
		//}
		//if param.Data != "" {
		//	updateP.Data = param.Data
		//}
		//if param.Categories != "" {
		//	updateP.Categories = param.Categories
		//}
		if param.Sig != "" {
			updateP.Signdata = param.Sig
		}
		if param.AutoFlag != "" {
			af := strings.ToLower(param.AutoFlag)
			if af == "true" || af == "false" {
				updateP.Autoflag = af
			} else {
				log.Println("SetSysParams AutoFlag  error.")
				return ErrData
			}
		}
		if IsIntDataValid(param.Catchtime) {
			time, terr := strconv.Atoi(param.Catchtime)
			if terr != nil || time < 0 {
				log.Println("SetSysParams Catchtime  error.")
				return ErrData
			}
			updateP.Catchtime = time
		}
		if IsIntDataValid(param.Nftloopcount) {
			Nftloopcount, terr := strconv.Atoi(param.Nftloopcount)
			if terr != nil || Nftloopcount < 0 {
				log.Println("SetSysParams Nftloopcount  error.=", terr)
				return ErrData

			}
			updateP.Nftloopcount = Nftloopcount
		}
		//if IsIntDataValid(param.Nftloopflush) {
		//	updateP.Nftloopflush, _ = strconv.Atoi(param.Nftloopflush)
		//}
		if IsIntDataValid(param.Collectcount) {
			Collectcount, terr := strconv.Atoi(param.Collectcount)
			if terr != nil || Collectcount < 0 {
				log.Println("SetSysParams() Collectcount err=", terr)
				return ErrData
			}
			updateP.Collectcount = Collectcount
		}
		//if IsIntDataValid(param.Collectflush) {
		//	updateP.Collectflush, _ = strconv.Atoi(param.Collectflush)
		//}
		if IsIntDataValid(param.Nftcount) {
			Nftcount, terr := strconv.Atoi(param.Nftcount)
			if terr != nil || Nftcount < 0 {
				log.Println("SetSysParams() Nftcount err=", terr)
				return ErrData
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
		return ErrDataBase
	}
	log.Println("SetSysParams =", updateP)
	return nil
}

func (nft NftDb) QueryExchangeInfo() (*ExchangeInfo, error) {
	var params Exchangeinfos
	var paraminfo ExchangeInfo

	err := nft.db.First(&params)
	log.Println(params)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			exchange := Exchangeinfos{}
			exchange.Categories = initCategories
			exchange.Userkyc = DefAutoSnft
			exchange.Nftaudit = DefAutoSnft
			exchange.Exchangerinfo = ExchangeName
			exchange.Transfersnft = DefAutoSnft
			exchange.Autocommitsnft = DefAutoSnft
			exchange.Allownft = DefAutoSnft
			exchange.Allowusermint = DefAutoSnft
			exchange.Uploadsize = DefUploadSize
			exchange.Backupipfs = DefAutoSnft
			exchange.Link = DefExchangeLink
			err = nft.db.Model(&Exchangeinfos{}).Create(&exchange)
			if err.Error != nil {
				fmt.Println("QueryExchangeInfo() create exchange err= ", err.Error)
				return nil, ErrDataBase
			}
		} else {
			fmt.Println("QueryExchangeInfo() not find err=", err.Error)
			return nil, ErrDataBase
		}
	}
	sysparam := SysParams{}
	err = nft.db.Model(&SysParams{}).Select("homepage").Last(&sysparam)
	if err.Error != nil {
		fmt.Println("QueryExchangeInfo() SysParams not find err=", err.Error)
		return nil, ErrDataBase
	}

	paraminfo.Adminaddr = ExchangerAddr
	paraminfo.Exchangerinfo = params.Exchangerinfo
	paraminfo.Icon = params.Icon
	paraminfo.Data = params.Data
	paraminfo.Categories = params.Categories
	paraminfo.Nftaudit = params.Nftaudit
	paraminfo.Restrictcode = params.Restrictcode
	paraminfo.Userkyc = params.Userkyc
	paraminfo.Announcement = strconv.FormatBool(AnnouncementRequired)
	paraminfo.ExchangerAddr = ExchangerAddr
	paraminfo.AutocommitSnft = params.Autocommitsnft
	AutocommitSnft, _ = strconv.ParseBool(params.Autocommitsnft)
	paraminfo.TransferNFT = params.Transfersnft
	TransferSNFT, _ = strconv.ParseBool(params.Transfersnft)
	paraminfo.AllowNft = params.Allownft
	AllowNft, _ = strconv.ParseBool(params.Allownft)
	paraminfo.AllowUserMint = params.Allowusermint
	AllowUserMinit, _ = strconv.ParseBool(params.Allowusermint)
	paraminfo.Backupipfs = params.Backupipfs
	Backupipfs, _ = strconv.ParseBool(params.Backupipfs)
	partnerlogo := strings.Split(params.Partnerslogo, ",")
	paraminfo.ExtendLogo = partnerlogo
	UploadSize = params.Uploadsize
	beego.BConfig.MaxMemory = int64(params.Uploadsize)
	paraminfo.Uploadsize = strconv.FormatUint(params.Uploadsize, 10)
	paraminfo.Homepage = sysparam.Homepage
	paraminfo.Link = params.Link
	paraminfo.Deflanguage = params.Deflanguage
	paraminfo.Desc = params.Desc

	//paraminfo.Nftlush 		= strconv.Itoa(params.Nftlush)
	return &paraminfo, nil
}

func (nft NftDb) SetExchangeInfo(param ExchangeInfo) error {
	var updateP Exchangeinfos
	err := nft.db.First(&updateP)
	if err.Error != nil {
		if nft.db.Error == gorm.ErrRecordNotFound {
			//updateP.NFT1155addr = NFT1155Addr
			//updateP.Adminaddr = AdminAddr
			//updateP.Lowprice = initLowprice
			//updateP.Royaltylimit = initRoyaltylimit
			//updateP.Categories = initCategories
			//updateP.Homepage = HomePages
			//key, _ := ExchangeGenerate()
			//updateP.Exchangerprv = key
			return ErrNotFound
		} else {
			fmt.Println("Exchangeinfos() not find err=", err.Error)
			return ErrDataBase
		}
	} else {
		if param.Nftaudit != "" {
			updateP.Nftaudit = param.Nftaudit
			audit, err := strconv.ParseBool(updateP.Nftaudit)
			if err != nil {
				log.Println("NftAudit input  error.")
				return ErrData
			}
			NFTUploadAuditRequired = audit
		}
		if param.TransferNFT != "" {
			updateP.Transfersnft = param.TransferNFT
			audit, err := strconv.ParseBool(param.TransferNFT)
			if err != nil {
				log.Println("TransferNFT input  error.")
				return ErrData
			}
			log.Println("TransferNFT =", param.TransferNFT)
			TransferSNFT = audit
		}
		if param.Userkyc != "" {
			updateP.Userkyc = param.Userkyc
			audit, err := strconv.ParseBool(updateP.Userkyc)
			if err != nil {
				log.Println("Userkyc input  error.")
				return ErrData

			}
			KYCUploadAuditRequired = audit
		}
		if param.Announcement != "" {
			audit, err := strconv.ParseBool(param.Announcement)
			if err != nil {
				log.Println("Announcement input  error.")
				return ErrData
			}
			AnnouncementRequired = audit
		}
		if param.AutocommitSnft != "" {
			updateP.Autocommitsnft = param.AutocommitSnft
			audit, err := strconv.ParseBool(param.AutocommitSnft)
			if err != nil {
				log.Println("AutocommitSnft input  error.")
				return ErrData
			}
			AutocommitSnft = audit
		}
		if param.AllowNft != "" {
			updateP.Allownft = param.AllowNft
			audit, err := strconv.ParseBool(param.AllowNft)
			if err != nil {
				log.Println("AllowNft input  error.")
				return ErrData
			}
			AllowNft = audit
		}
		if param.AllowUserMint != "" {
			updateP.Allowusermint = param.AllowUserMint
			audit, err := strconv.ParseBool(param.AllowUserMint)
			if err != nil {
				log.Println("AllowUserMint input  error.")
				return ErrData
			}
			AllowUserMinit = audit
		}
		if param.Backupipfs != "" {
			updateP.Backupipfs = param.Backupipfs
			audit, err := strconv.ParseBool(param.Backupipfs)
			if err != nil {
				log.Println("Backupipfs input  error.")
				return ErrData
			}
			Backupipfs = audit
		}
		if param.Uploadsize != "" {
			low, _ := strconv.ParseUint(param.Uploadsize, 10, 64)
			updateP.Uploadsize = low
			UploadSize = low
			beego.BConfig.MaxUploadSize = int64(low)
		}

		if param.Restrictcode != "" {
			updateP.Restrictcode = param.Restrictcode
		}

		if param.Exchangerinfo != "" {
			updateP.Exchangerinfo = param.Exchangerinfo
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
		if param.Desc != "" {
			updateP.Desc = param.Desc
		}
		if param.Partnerslogo != "" {
			tokenid := time.Now().UnixNano()
			name := strconv.Itoa(int(tokenid))
			err := SavePartnerslogoImage(ImageDir, name, param.Partnerslogo)
			if err != nil {
				fmt.Println("SetExchangeInfo() save image err=", err)
				return ErrNftImage
			}
			if updateP.Partnerslogo == "" {
				updateP.Partnerslogo = name
			} else {
				updateP.Partnerslogo = updateP.Partnerslogo + "," + name
			}
		}
		if param.Link != "" {
			updateP.Link = param.Link
		}
		if param.Deflanguage != "" {
			updateP.Deflanguage = param.Deflanguage
		}

	}

	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&Exchangeinfos{}).Last(&Exchangeinfos{}).Updates(&updateP)
		if err.Error != nil {
			fmt.Println("Exchangeinfos() update  err= ", err.Error)
			return ErrDataBase
		}
		if param.Homepage != "" {
			err = tx.Last(&SysParams{}).Updates(map[string]interface{}{"homepage": param.Homepage, "autoflag": "false"})
			if err.Error != nil {
				fmt.Println("Exchangeinfos() update SysParams err= ", err.Error)
				return ErrDataBase
			}
		}
		if param.AutoFlag != "" {
			err = tx.Last(&SysParams{}).Updates(map[string]interface{}{"autoflag": param.AutoFlag})
			if err.Error != nil {
				fmt.Println("Exchangeinfos() update autoflag err= ", err.Error)
				return ErrDataBase
			}
		}
		log.Println("Exchangeinfos =", updateP)
		return nil
	})

}

func (nft NftDb) SetAnnouncementParam(param string) error {

	param = strings.ToLower(param)
	announce, err := strconv.ParseBool(param)
	if err != nil {
		fmt.Println("SetAnnouncementParam()  err=", err)
		return ErrDataFormat
	}
	AnnouncementRequired = announce
	return nil
}

func (nft NftDb) SetExchageSig(exchange string) error {
	ExchangerAuth = exchange
	err := nft.db.Last(&SysParams{}).Update("exchangerauth", exchange)
	if err.Error != nil {
		fmt.Println("SetExchageSig() update exchangerauth err= ", err.Error)
		return ErrDataBase
	}
	return nil
}

func (nft NftDb) TranSnft() error {
	if TransferSNFT == true {
		return nil
	}
	log.Println("update blocknumber,transferNFT=", TransferSNFT)

	return nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Last(&SysParams{}).Updates(map[string]interface{}{"scannumber": 1})
		if err.Error != nil {
			fmt.Println("TranSnft() update SysParams Blocknumber err= ", err.Error)
			return err.Error
		}
		err = nft.db.Last(&Exchangeinfos{}).Updates(map[string]interface{}{"transfersnft": "true"})
		if err.Error != nil {
			fmt.Println("TranSnft() update Exchangeinfos transfersnft err= ", err.Error)
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) TranNft() error {
	if TransferSNFT == false {
		return nil
	}
	log.Println("update transfersnft,transferNFT=", TransferSNFT)

	err := nft.db.Last(&Exchangeinfos{}).Updates(map[string]interface{}{"transfersnft": "false"})
	if err.Error != nil {
		fmt.Println("TranNft() update Blocknumber err= ", err.Error)
		return err.Error
	}

	return nil
}

func (nft NftDb) DelPartnerLogo(name string) error {
	if name == "" {
		fmt.Println("input data error")
		return ErrData
	}
	var exchange Exchangeinfos
	err := nft.db.First(&exchange)
	if err.Error != nil {
		log.Println("DelPartnerLogo Exchangeinfos not found")
		return ErrNotFound
	}
	logolist := strings.Split(exchange.Partnerslogo, ",")
	partner := []string{}
	for i, j := range logolist {
		if j == name {
			partner = append(logolist[:i], logolist[i+1:]...)
			break
		}
		continue
	}
	partnerstr := ""
	for i, value := range partner {
		partnerstr = partnerstr + value
		if i != len(partner)-1 {
			partnerstr = partnerstr + ","
		}
	}
	err = nft.db.Model(&Exchangeinfos{}).Where("id = ? ", "1").Update("partnerslogo", partnerstr)
	if err.Error != nil {
		fmt.Println("DelPartnerLogo() update  err= ", err.Error)
		return ErrDataBase
	}
	delerr := DelPartnerslogoImage(ImageDir, name)
	if delerr != nil {
		fmt.Println("DelPartnerLogo() delete image err=", err)
		return ErrDeleteImg
	}

	return nil
}

func (nft NftDb) GetExchageSig() (bool, error) {
	var params SysParams
	var auth bool
	err := nft.db.Last(&params)
	if err.Error != nil {
		fmt.Println("GetExchageSig() update exchangerauth err= ", err.Error)
		return false, ErrDataBase
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
		return ErrDataFormat
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
	if AgentExchangePrv == "" {
		SuperAdminPrv = DefSuperAdminPrv
	} else {
		SuperAdminPrv = AgentExchangePrv
	}
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

	err = NewQueryCatch(QueryRedisCatchSvr, QueryRedisSvrPasswd)
	if err != nil {
		log.Println("InitSysParams() NewQueryCatch err=", err)
		return errors.New("InitSysParams() NewQueryCatch error.")
	}
	if LimitWritesDatabase {
		err = NewQueryMainCatch(MainRedisCatchSvr, MainRedisCatchSvrPasswd)
		if err != nil {
			log.Println("InitSysParams() NewQueryMainCatch err=", err)
			return errors.New("InitSysParams() NewQueryMainCatch error.")
		}
		go GetRedisCatch().ScanDirtyQuerys(GetRedisMainCatch())
	}
	nd, err := NewNftDb(Sqldsndb)
	if err != nil {
		fmt.Printf("InitSysParams() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
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
	log.Println("DebugAllowNft =", DebugAllowNft)
	if DebugAllowNft == "true" {
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
	} else {
		err = nd.TranNft()
		if err != nil {
			fmt.Printf("InitSysParams() TranNft() err = %v\n", err)
			return err
		}
		_, err = nd.QuerySysParams()
		if err != nil {
			fmt.Printf("InitSysParams() QuerySysParams() err = %s\n", err)
			return err
		}
	}
	go AutoPeriodEth(Sqldsndb)
	//go CaptchaDefault()
	go WormLogoDefault()
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
	Contract string `json:"contract" gorm:"type:char(42) NOT NULL;comment:'contract address'"`
	Tokenid  string `json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'contract token id'"`
}

func ScanLoop(sqldsn string, interval int, stop chan struct{}, stoped chan struct{}) {
	ticker := time.NewTicker(time.Minute * 1)
	scan := time.NewTicker(time.Duration(interval) * time.Minute)
	log.Println("ScanLoop() start: ", "interval=", interval)
	defer log.Println("ScanLoop() end")
	for {
		select {
		case <-scan.C:
			log.Println("ScanLoop() <-scan.C: start")
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
					//homepage.Nfts = []HomePageNft{{"", ""}}
					homepage.Nfts = []HomePageNft{{"", ""}, {"", ""}, {"", ""}, {"", ""}}
				}
			}
			collects := make([]Collects, 0, 20)
			limit, _ = strconv.Atoi(params.Nftcount)
			dberr = nd.db.Where("name <> ?", DefaultCollection).Order("transcnt desc").Limit(limit).Find(&collects)
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
					homepage.Collections = []HomePageCollections{{"", ""}, {"", ""}, {"", ""}, {"", ""}}
				}
			}
			var recCount int64
			var maxCount int64
			countSql := `select max(id) from nfts`
			dberr = nd.db.Raw(countSql).Scan(&maxCount)
			if TransSnft {
				dberr = nd.db.Model(&SysInfos{}).Select("snfttotal").First(&recCount)
			} else {
				dberr = nd.db.Model(&SysInfos{}).Select("nfttotal").First(&recCount)
			}
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
				log.Println("ScanLoop() recCount= ", recCount)
				homepage.NftLoop = []HomePageNft{}
				for i := 0; i < limit && int64(i) < recCount; {
					index := rand.Int63()%maxCount + 1
					log.Println("ScanLoop() rand.Int63() index= ", index)
					/*if index == 0 {
						index = 1
					}*/
					flag := scaned[index]
					if flag {
						//i = i - 1
						log.Println("ScanLoop() scaned[index] index= ", index)
						//time.Sleep(time.Second)
						continue
					}
					scaned[index] = true
					var nftRec Nfts
					dberr := nd.db.Where("id = ?", index).First(&nftRec)
					if dberr.Error != nil {
						//nd.Close()
						log.Println("ScanLoop() index=", index, "First(&nftRec) err = ", dberr.Error)
						continue
					}
					i++
					var hpnft HomePageNft
					hpnft.Contract = nftRec.Contract
					hpnft.Tokenid = nftRec.Tokenid
					homepage.NftLoop = append(homepage.NftLoop, hpnft)
				}
			} else {
				homepage.NftLoop = []HomePageNft{{"", ""}, {"", ""}, {"", ""}}
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
			//HomePageCatchs.HomePageFlashLock()
			//HomePageCatchs.HomePageFlashFlag = true
			//HomePageCatchs.HomePageFlashUnLock()
			nd.Close()
			log.Println("ScanLoop() <-scan.C: end")
		case <-ticker.C:
			log.Println("ScanLoop() <-ticker.C: start")
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
			log.Println("ScanLoop() <-ticker.C: end")
		case <-stop:
			log.Println("ScanLoop() <-stop: start")
			ticker.Stop()
			scan.Stop()
			stoped <- struct{}{}
			log.Println("ScanLoop() <-stop: end")
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

func HomePageRenew() error {
	log.Println("HomePageRenew: start")
	nd, err := NewNftDb(sqldsn)
	if err != nil {
		log.Printf("HomePageRenew() connect database err = %s\n", err)
		return err
	}
	params, err := nd.QuerySysParams()
	if err != nil {
		nd.Close()
		log.Printf("HomePageRenew() QuerySysParams() err = %s\n", err)
		return err
	}
	if params.AutoFlag == "false" {
		return nil
	}
	var hp, homepage HomePage
	err = json.Unmarshal([]byte(params.Homepage), &hp)
	if err != nil {
		log.Println("HomePageRenew() Unmarshal() err =", err)
		return err
	}
	homepage.Announcement = hp.Announcement
	trans := make([]HotTrans, 0, 20)
	limit, _ := strconv.Atoi(params.Nftcount)
	sql := "SELECT  contract, tokenid, count(tokenid) as tc from trans group by contract, tokenid order by tc desc limit "
	sql = sql + params.Nftcount
	dberr := nd.db.Raw(sql).Scan(&trans)
	if dberr.Error != nil {
		if dberr.Error != gorm.ErrRecordNotFound {
			nd.Close()
			log.Printf("HomePageRenew() Find(&trans) err = %s\n", dberr.Error)
			return err
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
			homepage.Nfts = []HomePageNft{{"", ""}, {"", ""}, {"", ""}, {"", ""}}
		}
	}
	collects := make([]Collects, 0, 20)
	limit, _ = strconv.Atoi(params.Nftcount)
	dberr = nd.db.Where("name <> ?", DefaultCollection).Order("transcnt desc").Limit(limit).Find(&collects)
	if dberr.Error != nil {
		if dberr.Error != gorm.ErrRecordNotFound {
			nd.Close()
			log.Printf("HomePageRenew() Find(&collects) err = %s\n", dberr.Error)
			return err
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
			homepage.Collections = []HomePageCollections{{"", ""}, {"", ""}, {"", ""}, {"", ""}}
		}
	}
	var recCount int64
	var maxCount int64
	countSql := `select max(id) from nfts`
	dberr = nd.db.Raw(countSql).Scan(&maxCount)
	if TransSnft {
		dberr = nd.db.Model(&SysInfos{}).Select("snfttotal").First(&recCount)
	} else {
		dberr = nd.db.Model(&SysInfos{}).Select("nfttotal").First(&recCount)
	}
	if dberr.Error != nil {
		log.Printf("HomePageRenew() Count(&recCount) err = %s\n", dberr.Error)
		recCount = 0
	}
	if recCount != 0 {
		rand.Seed(time.Now().UnixNano())
		limit, _ = strconv.Atoi(params.Nftloopcount)
		scaned := make(map[int64]bool)
		log.Println("HomePageRenew() recCount= ", recCount)
		homepage.NftLoop = []HomePageNft{}
		for i := 0; i < limit && int64(i) < recCount; {
			index := rand.Int63()%maxCount + 1
			log.Println("HomePageRenew() rand.Int63() index= ", index)
			/*if index == 0 {
				index = 1
			}*/
			flag := scaned[index]
			if flag {
				//i = i - 1
				log.Println("HomePageRenew() scaned[index] index= ", index)
				//time.Sleep(time.Second)
				continue
			}
			scaned[index] = true
			var nftRec Nfts
			dberr := nd.db.Where("id = ?", index).First(&nftRec)
			if dberr.Error != nil {
				//nd.Close()
				log.Println("HomePageRenew() index=", index, "First(&nftRec) err = ", dberr.Error)
				continue
			}
			i++
			var hpnft HomePageNft
			hpnft.Contract = nftRec.Contract
			hpnft.Tokenid = nftRec.Tokenid
			homepage.NftLoop = append(homepage.NftLoop, hpnft)
			if i == int(recCount) {
				break
			}
		}
	} else {
		homepage.NftLoop = []HomePageNft{{"", ""}, {"", ""}, {"", ""}}
	}
	homestr, err := json.Marshal(&homepage)
	if err != nil {
		nd.Close()
		log.Println("HomePageRenew() Marshal(&homepage) err = ", err)
		return err
	}
	newParam := SysParamsInfo{
		Homepage: string(homestr),
	}
	nd.SetSysParams(newParam)
	nd.Close()
	log.Println("HomePageRenew() <-scan.C: end")
	return nil
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

func WormLogoDefault() {
	if DefaultCaptcha == "" {
		DefaultCaptcha = "/ipfs/QmQjTwQDAxJ6cNhW7fQRC8EnAbdpSTPiP859m9EbNSs6Cx"
	}
	if DefaultMask == "" {
		DefaultMask = "/ipfs/QmR9W6BUBPvaBiabCk6DybAk6mhxzGN9fsZjRzXkzvMb2y/mask.png"
		DefaultMaskFrame = "/ipfs/QmR9W6BUBPvaBiabCk6DybAk6mhxzGN9fsZjRzXkzvMb2y/maskframe.png"
	}
	DefaultWormholes := "/ipfs/QmTr3WF9CG2wWNTo3eEZeyLGLcN6gAs68qFb87uihutVu1/worm.png"
	DefaultCaptchaNum = DefCaptchaNum

	fmt.Println("default worm:", DefaultWormholes)
	url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(100 * time.Second)

	var err error
	var wormdata io.Reader
	for {
		wormdata, err = s.Cat(DefaultWormholes)
		if err != nil {
			log.Printf("wormdata cat  [%v] failed! %v", DefaultWormholes, err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	wormbody, err := ioutil.ReadAll(wormdata)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return
	}
	wormimg, _, err := image.Decode(bytes.NewReader(wormbody))
	if err != nil {
		fmt.Printf("image Decode  failed! %v", err)
		return
	}

	emptyBuff := bytes.NewBuffer(nil)
	jpeg.Encode(emptyBuff, wormimg, nil)
	wormstr := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	partlogo := "data:image/jpg;base64," + wormstr

	err = SavePartnerslogoImage(ImageDir, "wormholes", partlogo)
	if err != nil {
		fmt.Println("Wormholes  logo save image err=", err)
		return
	}
	fmt.Println("default captcha init ok")
}

func CaptchaDefault() {
	if DefaultCaptcha == "" {
		DefaultCaptcha = "/ipfs/QmQjTwQDAxJ6cNhW7fQRC8EnAbdpSTPiP859m9EbNSs6Cx"
	}
	if DefaultMask == "" {
		DefaultMask = "/ipfs/QmR9W6BUBPvaBiabCk6DybAk6mhxzGN9fsZjRzXkzvMb2y/mask.png"
		DefaultMaskFrame = "/ipfs/QmR9W6BUBPvaBiabCk6DybAk6mhxzGN9fsZjRzXkzvMb2y/maskframe.png"
	}
	DefaultWormholes := "/ipfs/QmR9W6BUBPvaBiabCk6DybAk6mhxzGN9fsZjRzXkzvMb2y/worm.png"

	fmt.Println("default captcha:", DefaultCaptcha)
	url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(100 * time.Second)
	var maskdata io.Reader
	var maskframe io.Reader
	var err error
	for {
		maskdata, err = s.Cat(DefaultMask)
		if err != nil {
			log.Printf("mask Http  [%v] failed! %v", DefaultMask, err)
			time.Sleep(5 * time.Second)
			continue
		}
		maskframe, err = s.Cat(DefaultMaskFrame)
		if err != nil {
			log.Printf("mask Http  [%v] failed! %v", DefaultMask, err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}

	}
	//maskdata, err := s.Cat(DefaultMask)
	//if err != nil {
	//	log.Printf("mask Http  [%v] failed! %v", DefaultMask, err)
	//	return
	//}
	//var snft nftInfo
	//b, err := ioutil.ReadAll(rc)
	//if err != nil {
	//	log.Println("GetnftInfoFromIPFSWithShell() ReadAll() err=", err)
	//	return nil, err
	//}
	//maskdata, err := http.Get(DefaultCaptcha)
	//if err != nil {
	//	fmt.Printf("mask Http get [%v] failed! %v", maskdata, err)
	//	return
	//}

	jbody, err := ioutil.ReadAll(maskdata)
	if err != nil {
		log.Printf("mask Read cat response failed! %v", err)
		return
	}
	framebody, err := ioutil.ReadAll(maskframe)
	if err != nil {
		log.Printf("maskframe Read cat response failed! %v", err)
		return
	}

	newPath := ImageDir + "/captcha/"
	_, err = os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("CaptchaDefault() create dir err=", err)
			return
		}
	}
	f, err := os.Create(newPath + "mask")
	img, _, err := image.Decode(bytes.NewReader(jbody))
	if err != nil {
		fmt.Printf("mask Decode  failed! %v", err)
		return
	}
	ff, err := os.Create(newPath + "maskframe")
	ffimg, _, err := image.Decode(bytes.NewReader(framebody))
	if err != nil {
		fmt.Printf("maskframe Decode  failed! %v", err)
		return
	}

	defer f.Close()
	defer ff.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Println("mask png encode err=", err)
		return
	}
	err = png.Encode(ff, ffimg)
	if err != nil {
		log.Println("maskframe png encode err=", err)
		return
	}
	var v io.Reader
	for {
		v, err = s.Cat(DefaultCaptcha)
		if err != nil {
			log.Printf("captcha Http [%v] failed! %v", DefaultCaptcha, err)
			time.Sleep(2 * time.Second)
			continue
		} else {
			break
		}
	}
	//v, err := http.Get(DefaultCaptcha)
	//if err != nil {
	//	fmt.Printf("Http get [%v] failed! %v", DefaultCaptcha, err)
	//	return
	//}
	//defer v.Body.Close()
	content, err := ioutil.ReadAll(v)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return
	}
	var data []map[string]string
	err = json.Unmarshal(content, &data)
	if err != nil {
		fmt.Println(err)
	}
	DefaultCaptchaNum = len(data)
	for i, j := range data {
		si := fmt.Sprintf("%05x", i)
		fmt.Println(j)
		var jdata io.Reader
		for {
			jdata, err = s.Cat(j["url"])
			if err != nil {
				log.Printf("captcha cat [%v] failed! %v", DefaultCaptcha, err)
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
		jbody, err := ioutil.ReadAll(jdata)
		if err != nil {
			fmt.Printf("Read http response failed! %v", err)
			return
		}
		f, err := os.Create(newPath + si)
		img, _, err := image.Decode(bytes.NewReader(jbody))
		if err != nil {
			fmt.Printf("image Decode  failed! %v", err)
			return
		}
		defer f.Close()
		err = jpeg.Encode(f, img, nil)
		if err != nil {
			fmt.Println("jpeg encode err=", err)
			return
		}

	}

	var wormdata io.Reader
	for {
		wormdata, err = s.Cat(DefaultWormholes)
		if err != nil {
			log.Printf("wormdata cat  [%v] failed! %v", DefaultWormholes, err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	wormbody, err := ioutil.ReadAll(wormdata)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return
	}
	wormimg, _, err := image.Decode(bytes.NewReader(wormbody))
	if err != nil {
		fmt.Printf("image Decode  failed! %v", err)
		return
	}

	emptyBuff := bytes.NewBuffer(nil)
	jpeg.Encode(emptyBuff, wormimg, nil)
	wormstr := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	partlogo := "data:image/jpg;base64," + wormstr

	err = SavePartnerslogoImage(ImageDir, "wormholes", partlogo)
	if err != nil {
		fmt.Println("Wormholes  logo save image err=", err)
		return
	}
	fmt.Println("default captcha init ok")
}
