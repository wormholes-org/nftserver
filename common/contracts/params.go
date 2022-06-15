package contracts

import "fmt"

const (
	ReDialDelyTime = 5
	ZeroAddr          = "0x0000000000000000000000000000000000000000"
	//AdminListPrv = "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
	//TradeAuthAddrPrv="564ea566096d3de340fc5ddac98aef672f916624c8b0e4664a908cd2a2d156fe"
	//AdminMintPrv = "8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed"
)

var (
	EthNode string
	BrowseNode string
	EthWsNode string
	Weth9Addr string
	TradeAddr string
	Nft1155Addr string
	AdminAddr string
	AdminListPrv string
	TradeAuthAddrPrv string
	AdminMintPrv string
	SuperAdminAddr string
	ExchangeOwer string
)

func SetSysParams(ethNode, browseNode, ethWsNode, weth9addr, tradeaddr, nft1155addr, adminAddr, adminListPrv, tradeAuthAddrPrv, adminMintPrv, superadminaddr, exchangeOwer string)  {
	EthNode = ethNode
	BrowseNode = browseNode
	EthWsNode = ethWsNode
	Weth9Addr = weth9addr
	TradeAddr = tradeaddr
	Nft1155Addr = nft1155addr
	AdminAddr = adminAddr
	AdminListPrv = adminListPrv
	TradeAuthAddrPrv = tradeAuthAddrPrv
	AdminMintPrv = adminMintPrv
	SuperAdminAddr = superadminaddr
	ExchangeOwer = exchangeOwer
	fmt.Println("SetSysParams() EthNode=", EthNode)
	fmt.Println("SetSysParams() BrowseNode=", BrowseNode)
	fmt.Println("SetSysParams() EthWsNode=", EthWsNode)
	fmt.Println("SetSysParams() SuperAdminAddr=", SuperAdminAddr)
	fmt.Println("SetSysParams() ExchangeOwer=", ExchangeOwer)
}

func SetEthNode(ethNode string)  {
	EthNode = ethNode
	fmt.Println("SetSysParams() EthNode=", EthNode)
}