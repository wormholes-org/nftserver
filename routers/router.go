package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/controllers/nftexchangev1"
	"github.com/nftexchange/nftserver/controllers/nftexchangev2"
	"github.com/nftexchange/nftserver/models"
)

func init() {
	registFilters()
	registRouterV1()
	registRouterV2()
	nftdb := new(models.NftDb)
	nftdb.InitDb(models.SqlSvr, models.DbName)
}

func registRouterV1() {
	//user login
	beego.Router("/v1/login", &nftexchangev1.NftExchangeControllerV1{}, "post:UserLogin")
	//upload nft
	beego.Router("/v1/upload", &nftexchangev1.NftExchangeControllerV1{}, "post:UploadNft")
	//buy nft
	beego.Router("/v1/buy", &nftexchangev1.NftExchangeControllerV1{}, "post:BuyNft")
	//query nft
	beego.Router("/v1/queryNFTList", &nftexchangev1.NftExchangeControllerV1{}, "get:QueryAllNftProducts")
	//query user
	beego.Router("/v1/queryUser", &nftexchangev1.NftExchangeControllerV1{}, "post:QueryUserInfo")
	//
	//beego.Router("/v1/getimage", &controllers.NftExchangeController{}, "post:GetImageFromIPFS")
	beego.Router("/v1/ipfsaddtest", &nftexchangev1.NftExchangeControllerV1{}, "get:IpfsTest")
}

func registRouterV2() {
	if !models.LimitWritesDatabase {
		//Upload nft
		beego.Router("/v2/upload", &nftexchangev2.NftExchangeControllerV2{}, "post:UploadNft")
		//Upload nft  with original images
		beego.Router("/v2/uploadNftImage", &nftexchangev2.NftExchangeControllerV2{}, "post:UploadNftImage")
		//buy nft
		beego.Router("/v2/buy", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyNft")
		//To buy nft , the exchange initiates the transaction
		beego.Router("/v2/buying", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyingNft")
		//cancel purchase of nft works
		beego.Router("/v2/cancelBuy", &nftexchangev2.NftExchangeControllerV2{}, "post:CancelBuyNft")
		//Modify user information
		beego.Router("/v2/modifyUserInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:ModifyUserInfo")
		//For sale (on the shelf)
		beego.Router("/v2/sell", &nftexchangev2.NftExchangeControllerV2{}, "post:Sell")
		//Cancel the sale (off the shelf)
		beego.Router("/v2/cancelSell", &nftexchangev2.NftExchangeControllerV2{}, "post:CancelSell")
		//Audit NFTs
		beego.Router("/v2/vrfNFT", &nftexchangev2.NftExchangeControllerV2{}, "post:VerifyNft")
		//User-focused NFTs
		beego.Router("/v2/like", &nftexchangev2.NftExchangeControllerV2{}, "post:SetFavoriteNft")
		//New collection
		beego.Router("/v2/newCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:CreateCollection")
		//Modify collection information
		beego.Router("/v2/modifyCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:ModifyCollection")
		//Modify Collection Cover
		beego.Router("/v2/modifyCollectionsImage", &nftexchangev2.NftExchangeControllerV2{}, "post:ModifyCollectionsImage")
		//User apply for KYC
		beego.Router("/v2/userRequireKYC", &nftexchangev2.NftExchangeControllerV2{}, "post:UserRequireKYC")
		//Audit KYC
		beego.Router("/v2/userKYC", &nftexchangev2.NftExchangeControllerV2{}, "post:UserKYC")
		//set sysparams data
		beego.Router("/v2/modifySysParams", &nftexchangev2.NftExchangeControllerV2{}, "post:SetSysParams")
		//Modify the announcement release switch
		beego.Router("/v2/setAnnouncementParams", &nftexchangev2.NftExchangeControllerV2{}, "post:SetAnnouncementParams")
		//Notification of completion of purchase of nft works
		beego.Router("/v2/buyResultInterface", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyResultInterface")
		//Set signature data
		beego.Router("/v2/setExchangeSig", &nftexchangev2.NftExchangeControllerV2{}, "post:SetExchangeSig")
		//Add edit administrator
		beego.Router("/v2/modifyAdmins", &nftexchangev2.NftExchangeControllerV2{}, "post:SetAdmins")
		//remove admin
		beego.Router("/v2/delAdmins", &nftexchangev2.NftExchangeControllerV2{}, "post:DelAdmins")
		//Add announcement
		beego.Router("/v2/modifyAnnounce", &nftexchangev2.NftExchangeControllerV2{}, "post:SetAnnounce")
		//delete announcement
		beego.Router("/v2/delAnnounces", &nftexchangev2.NftExchangeControllerV2{}, "post:DelAnnounces")
		//Add and modify country information
		beego.Router("/v2/modifyCountry", &nftexchangev2.NftExchangeControllerV2{}, "post:SetCountrys")
		//upload files
		beego.Router("/v2/upLoadFile", &nftexchangev2.NftExchangeControllerV2{}, "post:UpLoadFile")
		//Upload snft collection
		beego.Router("/v2/createSnftCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:CreateSnftCollection")
		//modify snft collection
		beego.Router("/v2/modifySnftCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:SetSnftCollection")
		//delete snft collection
		beego.Router("/v2/delSnftCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:CreateSnftCollection")
		//New snft period
		beego.Router("/v2/creatssnftphase", &nftexchangev2.NftExchangeControllerV2{}, "post:CreatsSnftphase")
		//Modify period
		beego.Router("/v2/modifyPeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:SetSnftPeriod")
		//delete period
		beego.Router("/v2/delSnftPeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:DelSnftPeriod")
		//Modify snft collection
		beego.Router("/v2/modifySnftCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:SetSnftCollection")
		//delete snft collection
		beego.Router("/v2/delSnftCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:DelSnftCollection")
		//Modify snft
		beego.Router("/v2/modifySnft", &nftexchangev2.NftExchangeControllerV2{}, "post:DelSnftPeriod")
		//delete snft
		beego.Router("/v2/delSnft", &nftexchangev2.NftExchangeControllerV2{}, "post:DelSnftPeriod")
		//snft collect setsnft
		beego.Router("/v2/modifycollectsnft", &nftexchangev2.NftExchangeControllerV2{}, "post:SetCollectSnft")
		//Set a voting period
		beego.Router("/v2/setVotePeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:SetVotePeriod")
		//Select the betting chain
		beego.Router("/v2/setperiodeth", &nftexchangev2.NftExchangeControllerV2{}, "post:SetPeriodEth")
		//Add subscription email
		beego.Router("/v2/setSubscribeEmail", &nftexchangev2.NftExchangeControllerV2{}, "post:SetSubscribeEmail")
		//Delete subscription mailbox
		beego.Router("/v2/delSubscribeEmail", &nftexchangev2.NftExchangeControllerV2{}, "post:DelSubscribeEmail")
		//delete nft
		beego.Router("/v2/delnft", &nftexchangev2.NftExchangeControllerV2{}, "post:DelNft")
		//delete collection
		beego.Router("/v2/delcollect", &nftexchangev2.NftExchangeControllerV2{}, "post:DelCollect")
		//Bulk purchase of nft works, initiated by the exchange
		beego.Router("/v2/groupbuying", &nftexchangev2.NftExchangeControllerV2{}, "post:GroupBuyingNft")
		//Group sale (on the shelf)
		beego.Router("/v2/groupsell", &nftexchangev2.NftExchangeControllerV2{}, "post:GroupSell")
		//Mass Cancellation of Sale (on the shelf)
		beego.Router("/v2/groupcancelsell", &nftexchangev2.NftExchangeControllerV2{}, "post:GroupCancelSell")
		//modify nft
		beego.Router("/v2/modifynft", &nftexchangev2.NftExchangeControllerV2{}, "post:SetNft")
		//user certify
		beego.Router("/v2/usercertify", &nftexchangev2.NftExchangeControllerV2{}, "post:UserSubmitCertify")
		//delete partnerslogo
		beego.Router("/v2/delpartnerslogo", &nftexchangev2.NftExchangeControllerV2{}, "post:DelPartnersLogo")
		//query expired nft
		beego.Router("/v2/queryexpirednft", &nftexchangev2.NftExchangeControllerV2{}, "post:GetExpiredNft")
		//delete expired nft
		beego.Router("/v2/delexpirednft", &nftexchangev2.NftExchangeControllerV2{}, "post:DelExpiredNft")
		//query collection and under nft
		beego.Router("/v2/querycollectionnft", &nftexchangev2.NftExchangeControllerV2{}, "post:SNftCollectionSearch")
		//modify snftcollection snft
		beego.Router("/v2/setperiod", &nftexchangev2.NftExchangeControllerV2{}, "post:SetPeriod")
		//query period collection snft
		beego.Router("/v2/queryperiod", &nftexchangev2.NftExchangeControllerV2{}, "get:GetPeriod")
		//set user agreement and privacy policy
		beego.Router("/v2/setAgreement", &nftexchangev2.NftExchangeControllerV2{}, "post:SetAgreement")

	}
	//User login
	beego.Router("/v2/login", &nftexchangev2.NftExchangeControllerV2{}, "post:UserLogin")
	//Notification of completion of purchase of nft works
	//beego.Router("/v2/buyResult", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyResult")
	//Query nft
	beego.Router("/v2/queryNFTList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryNftList")
	//query user
	beego.Router("/v2/queryUserInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserInfo")
	//Query single NFT information
	beego.Router("/v2/queryNFT", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryNft")
	//Querying information about a single SNFT fragment
	beego.Router("/v2/querySnftChip", &nftexchangev2.NftExchangeControllerV2{}, "post:QuerySnftChip")
	//Query SNFT information in a certain period
	beego.Router("/v2/queryStageSnft", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryStageSnft")
	//Query the owner's collection list
	beego.Router("/v2/queryOwnerSnftCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryOwnerSnftCollections")
	//Query the total number of fragments owned by an account
	beego.Router("/v2/queryOwnerSnftChipAmount", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryOwnerSnftChipAmount")
	//Query snft in the collection according to the creator, collection name
	beego.Router("/v2/querySnftByCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:QuerySnftByCollection")
	//Inquiry period information
	beego.Router("/v2/queryStageList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryStageList")
	//Query period collection information
	beego.Router("/v2/queryStageCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryStageCollection")
	//Return snft fragments based on user address
	beego.Router("/v2/queryOwnerSnftChip", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryOwnerSnftChip")
	//Return snft data according to array
	beego.Router("/v2/queryArraySnft", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryArraySnft")

	//beego.Router("/v2/modifyNFT", &controllers.NftExchangeControllerV2{}, "post:ModifyNFT")
	//Query the NFT pending review list
	beego.Router("/v2/queryPendingVrfList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryPendingVerificationList")
	//Get market data
	beego.Router("/v2/queryMarketInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryMarketInfo")
	//Get user NFT list
	beego.Router("/v2/queryUserNFTList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserNftList")
	//Query user collection list
	beego.Router("/v2/queryUserCollectionList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserCollectionList")
	//Query user transaction history
	beego.Router("/v2/queryUserTradingHistory", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserTradingHistory")
	//Get a list of user NFT bids
	beego.Router("/v2/queryUserOfferList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserOfferList")
	//Get the list of users' NFTs being bid on
	beego.Router("/v2/queryUserBidList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserBidList")
	//Get user NFT watchlist
	beego.Router("/v2/queryUserFavoriteList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserFavoriteList")
	//Get a list of NFT collections
	beego.Router("/v2/queryNFTCollectionList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryNftCollectionList")
	//Get collection details
	beego.Router("/v2/queryCollectionInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryCollectionInfo")
	//Get transaction history of NFT exchanges
	beego.Router("/v2/queryMarketTradingHistory", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryMarketTradingHistory")
	//Get user KYC list
	beego.Router("/v2/queryPendingKYCList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryPendingKYCList")
	//Fuzzy query nft, collection, account
	beego.Router("/v2/search", &nftexchangev2.NftExchangeControllerV2{}, "post:Search")
	//Obtain the authorization of the exchange
	beego.Router("/v2/askForApprove", &nftexchangev2.NftExchangeControllerV2{}, "post:AskForApprove")
	//Get home page information
	beego.Router("/v2/queryHomePage", &nftexchangev2.NftExchangeControllerV2{}, "get:QueryHomePage")
	//get sysparams data
	beego.Router("/v2/querySysParams", &nftexchangev2.NftExchangeControllerV2{}, "get:GetSysParams")
	beego.Router("/v2/testRecover", &nftexchangev2.NftExchangeControllerV2{}, "post:Recover")
	beego.Router("/v2/getSysParam", &nftexchangev2.NftExchangeControllerV2{}, "post:GetSysParamByParams")
	//Query whether to sign
	beego.Router("/v2/getExchangeSig", &nftexchangev2.NftExchangeControllerV2{}, "get:GetExchangeSig")
	//Query
	//beego.Router("/v1/getimage", &controllers.NftExchangeController{}, "post:GetImageFromIPFS")
	beego.Router("/v2/ipfsaddtest", &nftexchangev2.NftExchangeControllerV2{}, "get:IpfsTest")
	//Query version number
	beego.Router("/v2/version", &nftexchangev2.NftExchangeControllerV2{}, "get:GetVersion")
	//Batch Import
	//beego.Router("/v2/batchNewCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:BatchCreateCollection")
	//beego.Router("/v2/batchUpload", &nftexchangev2.NftExchangeControllerV2{}, "post:BatchUploadNft")
	//beego.Router("/v2/batchBuyResultInterface", &nftexchangev2.NftExchangeControllerV2{}, "post:BatchBuyResultInterface")
	//query administrator
	beego.Router("/v2/queryAdmins", &nftexchangev2.NftExchangeControllerV2{}, "post:GetAdmins")
	//Address query administrator
	beego.Router("/v2/queryAdminsByAddr", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryAdminsByAddr")
	//Query announcement
	beego.Router("/v2/queryAnnounce", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryAnnounce")
	//Query country information
	beego.Router("/v2/queryCountry", &nftexchangev2.NftExchangeControllerV2{}, "get:GetCountrys")
	//query period
	beego.Router("/v2/getSnftPeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:GetSnftPeriod")
	//query snft  collection
	beego.Router("/v2/querySnftCollection", &nftexchangev2.NftExchangeControllerV2{}, "post:GetSnftCollection")
	//query snft
	//beego.Router("/v2/querySnft", &nftexchangev2.NftExchangeControllerV2{}, "post:DelSnftPeriod")
	//query snft
	beego.Router("/v2/snftSearch", &nftexchangev2.NftExchangeControllerV2{}, "post:SnftSearch")
	//search collect
	beego.Router("/v2/snftCollectSearch", &nftexchangev2.NftExchangeControllerV2{}, "post:SnftCollectSearch")
	//Get a voting period
	beego.Router("/v2/queryVotePeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:GetVotePeriod")
	//Get all voting periods
	beego.Router("/v2/queryAllvotePeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:GetAllVotePeriod")
	//Get three voting periods
	beego.Router("/v2/queryAccedvotePeriod", &nftexchangev2.NftExchangeControllerV2{}, "post:GetAccedVotePeriod")
	//Get total nft data
	beego.Router("/v2/queryNftMarketInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:GetNftMarketInfo")
	//Query subscription mailbox
	beego.Router("/v2/querySubscribeEmails", &nftexchangev2.NftExchangeControllerV2{}, "post:QuerySubscribeEmails")
	//query aptcha
	beego.Router("/v2/querycaptcha", &nftexchangev2.NftExchangeControllerV2{}, "post:GetCaptcha")
	//auth captcha
	beego.Router("/v2/authcaptcha", &nftexchangev2.NftExchangeControllerV2{}, "post:AuthCaptcha")
	//admin login
	beego.Router("/v2/adminlogin", &nftexchangev2.NftExchangeControllerV2{}, "post:AdminLogin")

}

func registFilters() {
	//Authorization check before routing
	beego.InsertFilter("/v2/upload", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/buy", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/cancelBuy", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/buyResult", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryNFTList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/modifyUserInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryNFT", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/sell", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/cancelSell", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryPendingVrfList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/vrfNFT", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryMarketInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/like", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/newCollections", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/modifyCollections", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/modifyCollectionsImage", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserNFTList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserCollectionList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserTradingHistory", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserOfferList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserBidList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserFavoriteList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryNFTCollectionList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryCollectionInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryMarketTradingHistory", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/userRequireKYC", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/userKYC", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryPendingKYCList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/search", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/askForApprove", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryHomePage", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/get_sys_para", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/set_sys_para", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/buyResultInterface", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/version", beego.BeforeRouter, nftexchangev2.CheckToken)

}
