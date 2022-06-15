package main

import (
	"github.com/nftexchange/nftserver/common/contracts"
)

const (
	EthNode = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	//EthNode = "http://192.168.1.235:8546"
	BrowseNode = "http://192.168.1.235:8546"
	EthersWsNode = "wss://rinkeby.infura.io/ws/v3/97cb2119c79842b7818a7a37df749b2b"
	Weth9Addr = "0xf4bb2e28688e89fcce3c0580d37d36a7672e8a9f"
	TradeAddr = "0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"
	Nft1155Addr = "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5"
	adminAddr = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"
	AdminListPrv = "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
	TradeAuthAddrPrv="564ea566096d3de340fc5ddac98aef672f916624c8b0e4664a908cd2a2d156fe"
	AdminMintPrv = "8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed"
	SuperAdminPrv = "2ABE62D35B09680F007B225C318D5A672CA3E956B91BEEE5A5BA004A22DAAC2C"
	ExchangeOwer = "2ABE62D35B09680F007B225C318D5A672CA3E956B91BEEE5A5BA004A22DAAC2C"
)

func init()  {
	contracts.SetSysParams(EthNode, BrowseNode, EthersWsNode, Weth9Addr, TradeAddr,
		Nft1155Addr, adminAddr, AdminListPrv, TradeAuthAddrPrv, AdminMintPrv, SuperAdminPrv, ExchangeOwer)
}