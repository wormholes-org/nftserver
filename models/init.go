package models

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"log"
)

var sqldsn string
var sqldsndb string
var Sqldsndb string
var sqllocaldsndb string

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	DbName, _ = beego.AppConfig.String("dbname")
	dbUserName, _ := beego.AppConfig.String("dbusername")
	dbUserPassword, _ := beego.AppConfig.String("dbuserpassword")
	dbServerIP, _ := beego.AppConfig.String("dbserverip")
	dbServerPort, _ := beego.AppConfig.String("dbserverport")
	//const SqlSvr = "admin:user123456@tcp(192.168.1.238:3306)/"
	SqlSvr = dbUserName + ":" + dbUserPassword + "@tcp(" + dbServerIP + ":" + dbServerPort + ")/"
	fmt.Println("SqlSvr=", SqlSvr)
	sqldsn = SqlSvr + localtime
	sqldsndb = SqlSvr + DbName + localtime
	sqllocaldsndb = SqlSvr + DbName + localtime
	Sqldsndb = sqldsndb
	TradeAddr, _ = beego.AppConfig.String("TradeAddr")
	NFT1155Addr, _ = beego.AppConfig.String("NFT1155Addr")
	AdminAddr, _ = beego.AppConfig.String("AdminAddr")
	EthersNode, _ = beego.AppConfig.String("EthersNode")
	NftIpfsServerIP, _ = beego.AppConfig.String("nftIpfsServerIP")
	NftstIpfsServerPort, _ = beego.AppConfig.String("nftIpfsServerPort")
	BackupIpfsUrl, _ = beego.AppConfig.String("backupipfsurl")
	if BackupIpfsUrl == "" {
		BackupIpfsUrl = "127.0.0.1:5001"
	}
	BrowseNode, _ = beego.AppConfig.String("BrowseNode")
	EthersWsNode, _ = beego.AppConfig.String("EthersWsNode")
	ImageDir, _ = beego.AppConfig.String("ImageDir")
	Weth9Addr, _ = beego.AppConfig.String("Weth9Addr")
	AdminListPrv, _ = beego.AppConfig.String("AdminListPrv")
	//SuperAdminPrv, _ = beego.AppConfig.String("SuperAdminPrv")
	TradeAuthAddrPrv, _ = beego.AppConfig.String("TradeAuthAddrPrv")
	AdminMintPrv, _ = beego.AppConfig.String("AdminMintPrv")
	NFTUploadAuditRequired, _ = beego.AppConfig.Bool("NFTUploadAuditRequired")
	Authorize, _ = beego.AppConfig.String("Authorize")
	DebugPort, _ = beego.AppConfig.String("DebugPort")
	DebugAllowNft, _ = beego.AppConfig.String("AllowNft")
	DefaultCaptcha, _ = beego.AppConfig.String("captachaurl")
	AnnouncementRequired = true
}
