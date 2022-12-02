package models

const (
	DefSuperAdminPrv = "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
	DefSuperAddr     = "0x7fbc8ad616177c6519228fca4a7d9ec7d1804900"

	SnftOffset                = 41
	snftCollectionOffset      = 40
	SnftStageOffset           = 39
	SnftCollectionsStageIndex = 30

	SnftExchangeChip      = 42
	SnftExchangeSnft      = 41
	SnftExchangeColletion = 40
	SnftExchangeStage     = 39
	//DefSuperAdminPrv = "17109ec24e570fe115f363553251f4d476ac4113905dd817eb4660c68d324aae"
	//DefSuperAddr = "0x01842a2cf56400a245a56955dc407c2c4137321e"
)

var (
	NftCacheDirtyName  = []string{"QueryNftByFilterNftSnft", "QueryHomePage", "SnftSearch"}
	UploadNftDirtyName = []string{"QueryNftByFilterNftSnft", "QueryNFTCollectionList", "QueryOwnerSnftChip",
		"querySnftChip", "queryStageSnft", "queryOwnerSnftCollections", "querySnftByCollection", "queryStageList",
		"queryStageCollection", "queryNFTList"}
	SetNftDirtyName  = []string{"QueryNftByFilterNftSnft", "QueryHomePage", "SnftSearch"}
	TradingDirtyName = []string{"QueryMarketTradingHistory", "QueryNFTCollectionList"}
	AnnouncementName = []string{"Announcement"}
	CollectionList   = []string{"QueryNFTCollectionList"}
	MarketNftInfo    = []string{"GetNftMarketInfo"}
	SnftExchange     = []string{"QuerySnftChip", "QueryStageCollection", "QueryStageSnft", "QueryStageList", "QuerySnftByCollection", "QueryOwnerSnftCollection", "QueryOwnerSnftChip"}
	AllDirty         = []string{"QuerySnftChip", "QueryStageCollection", "QueryStageSnft", "QueryStageList", "QuerySnftByCollection",
		"QueryOwnerSnftCollection", "QueryOwnerSnftChip", "QueryNftByFilterNftSnft", "SnftSearch", "QueryNFTCollectionList",
		"QueryOwnerSnftChip", "QueryMarketTradingHistory", "QueryNFTCollectionList", "Announcement", "QueryHomePage", "Search", "GetNftMarketInfo"}
)
