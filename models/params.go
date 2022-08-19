package models

const (
	DefSuperAdminPrv = "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
	DefSuperAddr     = "0x7fbc8ad616177c6519228fca4a7d9ec7d1804900"

	SnftOffset                = 40
	SnftStageOffset           = 38
	SnftCollectionsStageIndex = 30
	snftCollectionOffset      = 39
	//DefSuperAdminPrv = "17109ec24e570fe115f363553251f4d476ac4113905dd817eb4660c68d324aae"
	//DefSuperAddr = "0x01842a2cf56400a245a56955dc407c2c4137321e"
)

var (
	NftCacheDirtyName    = []string{"QueryNftByFilterNftSnft", "QueryHomePage", "SnftSearch"}
	UploadNftDirtyName   = []string{"QueryNftByFilterNftSnft", "SnftSearch", "QueryNFTCollectionList", "QueryOwnerSnftChip"}
	SetNftDirtyName      = []string{"QueryNftByFilterNftSnft", "QueryHomePage", "SnftSearch"}
	TradingDirtyName     = []string{"QueryMarketTradingHistory", "QueryNFTCollectionList"}
	NewSnftCollect       = []string{"CollectSearch"}
	NewSnftPeriod        = []string{"GetSnftPeriod"}
	ModifySnftCollect    = []string{"CollectSearch", "GetSnftCollection", "GetSnftPeriod", "GetAllVotePeriod", "GetVoteSnftPeriod"}
	ModifySnftPeriod     = []string{"GetSnftPeriod", "GetAllVotePeriod", "GetVoteSnftPeriod"}
	VoteSnftPeriod       = []string{"GetAllVotePeriod", "GetVoteSnftPeriod"}
	AdminDirtyName       = []string{"QueryAdmins"}
	SysParamsDirtyName   = []string{"QuerySysParams"}
	KYCListDirtyName     = []string{"QueryPendingKYCList"}
	NftVerifiedDirtyName = []string{"QueryUnverifiedNfts"}
	AnnouncementName     = []string{"Announcement"}
	CollectionList       = []string{"QueryNFTCollectionList"}
	SnftExchange         = []string{"QuerySnftChip", "QueryStageCollection", "QueryStageSnft", "QueryStageList", "QuerySnftByCollection", "QueryOwnerSnftCollection", "QueryOwnerSnftChip"}
)
