package models

import (
	"fmt"
	"log"
	//"github.com/nftexchange/nftserver/ethhelper"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type WnftRecord struct {
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) ;comment:'nft拥有者地址'"`
	Md5            string `json:"md5" gorm:"type:longtext ;comment:'图片md5值'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft名称'"`
	Desc           string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft描述'"`
	Meta           string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'元信息'"`
	Nftmeta        string `json:"nftmeta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'元信息、tokenid'"`
	Url            string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc原始数据保持地址'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) ;comment:'合约地址'"`
	Tokenid        string `json:"nft_token_id" gorm:"type:char(42) ;comment:'唯一标识nft标志'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) DEFAULT NULL;comment:'wormholes链唯一标识nft标志'"`
	Snftstage      string `json:"snftstage" gorm:"type:char(42) DEFAULT NULL;comment:'wormholes链snft期号'"`
	Snftcollection string `json:"snftcollection" gorm:"type:char(42) DEFAULT NULL;comment:'wormholes链snft合集'"`
	Snft           string `json:"snft" gorm:"type:char(42) ;comment:'wormholes链snft'"`
	Count          int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Approve        string `json:"approve" gorm:"type:longtext ;comment:'授权'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft分类'"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) ;comment:'合集创建者地址'"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'NFT合集名'"`
	Image          string `json:"asset_sample" gorm:"type:longtext ;comment:'缩略图二进制数据'"`
	//Imageid			string 		`json:"imageid" gorm:"type:char(42) ;comment:'图片存储索引'"`
	Hide         string `json:"hide" gorm:"type:char(20) ;comment:'是否让其他人看到'"`
	Signdata     string `json:"sig" gorm:"type:longtext ;comment:'签名数据，创建时产生'"`
	Createaddr   string `json:"user_addr" gorm:"type:char(42) ;comment:'创建nft地址'"`
	Verifyaddr   string `json:"vrf_addr" gorm:"type:char(42) ;comment:'验证人地址'"`
	Currency     string `json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Price        uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'创建时定的价格'"`
	Royalty      int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Paychan      string `json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'交易通道'"`
	TransCur     string `json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Transprice   uint64 `json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'交易成交价格'"`
	Transtime    int64  `json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'最后交易时间'"`
	Createdate   int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft创建时间'"`
	Favorited    int    `json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'被关注计数'"`
	Transcnt     int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt     uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	Verified     string `json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'nft作品是否通过审核'"`
	Verifieddesc string `json:"Verifieddesc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'审核描述：未通过审核描述'"`
	Verifiedtime int64  `json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'审核时间'"`
	Selltype     string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft交易类型'"`
	//Sellprice		uint64		`json:"sellingprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'正在销售价格'"`
	Mintstate string `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'铸币状态'"`
	//Extend    string `json:"extend" gorm:"type:longtext ;comment:'扩展字段'"`
}

type Wnfts struct {
	gorm.Model
	NftRecord
}

func (v Wnfts) TableName() string {
	return "wnfts"
}

func (nft NftDb) UploadWNft(nftInfo *SnftInfo) error {
	//user_addr := strings.ToLower(nftInfo.User_addr)
	creator_addr := strings.ToLower(nftInfo.CreatorAddr)
	//CollectionsCreator := strings.ToLower(nftInfo.CollectionsCreator)
	owner_addr := strings.ToLower(nftInfo.Ownaddr)
	nftaddress := strings.ToLower(nftInfo.Nftaddr)
	contract := strings.ToLower(nftInfo.Contract)
	md5 := nftInfo.Md5
	meta := nftInfo.Meta
	desc := nftInfo.Desc
	name := nftInfo.Name
	source_url := nftInfo.SourceUrl
	//nft_token_id := nftInfo.Nft_token_id
	collections := nftInfo.CollectionsName
	categories := nftInfo.Category
	//hide := nftInfo.Hide
	royalty := strconv.FormatInt(int64(nftInfo.Royalty * 100), 10)
	//royalty := nftInfo.Royalty
	//count := nftInfo.Count
	//asset_sample := "nftInfo.Image"
	//fmt.Println("UploadNft() user_addr=", user_addr,"      time=", time.Now().String())
	//UserSync.Lock(user_addr)
	//defer UserSync.UnLock(user_addr)
	fmt.Println("UploadWNft() begin ->> time = ", time.Now().String()[:22])
	//fmt.Println("UploadNft() user_addr = ", user_addr)
	fmt.Println("UploadWNft() creator_addr = ", creator_addr)
	fmt.Println("UploadWNft() owner_addr = ", owner_addr)
	fmt.Println("UploadWNft() md5 = ", md5)
	fmt.Println("UploadWNft() name = ", name)
	fmt.Println("UploadWNft() desc = ", desc)
	fmt.Println("UploadWNft() meta = ", meta)
	fmt.Println("UploadWNft() source_url = ", source_url)
	fmt.Println("UploadWNft() nft_contract_addr = ", contract)
	//fmt.Println("UploadNft() nft_token_id = ", nft_token_id)
	fmt.Println("UploadWNft() categories = ", categories)
	fmt.Println("UploadWNft() collections = ", collections)
	//fmt.Println("UploadNft() asset_sample = ", asset_sample)
	//fmt.Println("UploadNft() hide = ", hide)
	fmt.Println("UploadWNft() royalty = ", royalty)
	//fmt.Println("UploadNft() sig = ", sig)

	//if IsIntDataValid(count) != true {
	//	return ErrDataFormat
	//}
	nftExistFlag := false
	nftRec := Nfts{}
	err := nft.db.Select([]string{"id", "ownaddr"}).Where("nftaddr = ?", nftaddress).First(&nftRec)
	if err.Error == nil {
		if nftRec.Ownaddr == ZeroAddr {
			nftExistFlag = true
		} else {
			log.Println("UploadWNft() nft exist.")
			return nil
		}
	} else {
		if err.Error != gorm.ErrRecordNotFound {
			log.Println("UploadWNft() nfts dbase error.")
			return err.Error
		}
	}
	if IsIntDataValid(royalty) != true {
		return ErrDataFormat
	}
	r, _ := strconv.Atoi(royalty)
	fmt.Println("UploadWNft() royalty=", r, "SysRoyaltylimit=", SysRoyaltylimit, "RoyaltyLimit", RoyaltyLimit)
	if r > SysRoyaltylimit || r > RoyaltyLimit {
		return ErrRoyalty
	}
	count := 1
	if nft.IsValidCategory(categories) {
		return ErrNoCategory
	}

	var collectRec Collects
	snftCollection := ""
	if collections != "" {
		//collections = nftaddress[:snftCollectionOffset] + "." + collections
		snftStage := nftaddress[:SnftStageOffset]
		snftCollection = nftaddress[:snftCollectionOffset]
		err := nft.db.Where("createaddr = ? AND  name=?",
			nftInfo.CollectionsCreator, collections).First(&collectRec)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound{
			log.Println("UploadWNft() database err=", err.Error)
			return ErrCollectionNotExist
		}
		if err.Error == gorm.ErrRecordNotFound {
			collectRec = Collects{}
			collectRec.Createaddr = nftInfo.CollectionsCreator
			collectRec.Snftstage = snftStage
			collectRec.Snftcollection = snftCollection
			collectRec.Name = collections
			collectRec.Desc = nftInfo.CollectionsDesc
			collectRec.Img = nftInfo.CollectionsImgUrl
			collectRec.Contract = /*nftInfo.CollectionsExchanger*/"0x0000000000000000000000000000000000000000"
			collectRec.Contracttype = "snft"
			collectRec.Categories = nftInfo.CollectionsCategory

			err := nft.db.Model(&Collects{}).Create(&collectRec)
			if err.Error != nil {
				log.Println("UploadWNft create Collections() err=", err.Error)
				return err.Error
			}
		}
	}
	var NewTokenid string
	if nftExistFlag {
		log.Println("UploadWNft() snft exist!")
		/*nfttab := Nfts{}
		nfttab.Createaddr = creator_addr
		nfttab.Ownaddr = owner_addr
		nfttab.Name = name
		nfttab.Desc = desc
		nfttab.Meta = meta
		nfttab.Categories = categories
		nfttab.Collectcreator = collectRec.Createaddr
		nfttab.Collections = collections
		nfttab.Url = source_url
		nfttab.Selltype = SellTypeNotSale.String()
		nfttab.Verifiedtime = time.Now().Unix()
		nfttab.Createdate = time.Now().Unix()
		nfttab.Transcnt = 0
		nfttab.Transamt = 0*/
		nfttab := map[string]interface{}{
			"Createaddr": creator_addr,
			"Ownaddr": owner_addr,
			"Name": name,
			"Desc": desc,
			"Meta": meta,
			"Categories": categories,
			"Collectcreator": collectRec.Createaddr,
			"Collections": collections,
			"Url": source_url,
			"Selltype": SellTypeNotSale.String(),
			"Verifiedtime": time.Now().Unix(),
			"Createdate": time.Now().Unix(),
			"Transprice": 0,
			"Price": 0,
			"Transtime": 0,
			"Transcnt": 0,
			"Transamt": 0,
			"Favorited": 0,
		}
		sysInfo := SysInfos{}
		dberr := nft.db.Model(&SysInfos{}).Last(&sysInfo)
		if dberr.Error != nil {
			if dberr.Error != gorm.ErrRecordNotFound {
				log.Println("UploadWNft() SysInfos err=", dberr)
				return ErrCollectionNotExist
			}
		}
		log.Println("UploadWNft() SysInfos snfttotal count=", sysInfo.Snfttotal)
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Where("id = ?", nftRec.ID).Updates(&nfttab)
			if err.Error != nil {
				log.Println("UploadWNft() err=", err.Error)
				return err.Error
			}
			//if collectRec.Snftcollection != snftCollection {
			err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
				nftInfo.CollectionsCreator, collections).Update("totalcount", collectRec.Totalcount+1)
			if err.Error != nil {
				fmt.Println("UploadWNft() add collectins totalcount err= ", err.Error)
				return err.Error
			}
			//}
			if nftaddress[len(nftaddress)-2:] == "00" {
				NftCatch.SetFlushFlag()
				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal",sysInfo.Snfttotal + 1)
				if err.Error != nil {
					fmt.Println("UploadWNft() add  SysInfos snfttotal err=", err.Error)
					return err.Error
				}
			}
			return nil
		})
	} else {
		rand.Seed(time.Now().UnixNano())
		var i int
		for i = 0; i < genTokenIdRetry; i++ {
			//NewTokenid := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			s := fmt.Sprintf("%d", rand.Int63())
			if len(s) < 15 {
				continue
			}
			s = s[len(s)-13:]
			NewTokenid = s
			if s[0] == '0' {
				continue
			}
			fmt.Println("UploadWNft() NewTokenid=", NewTokenid)
			nfttab := Nfts{}
			err := nft.db.Where("tokenid = ? ", NewTokenid).First(&nfttab)
			if err.Error == gorm.ErrRecordNotFound {
				break
			}
		}
		if i >= 20 {
			fmt.Println("UploadWNft() generate tokenId error.")
			return ErrGenerateTokenId
		}

		nfttab := Nfts{}
		nfttab.Tokenid = NewTokenid
		nfttab.Nftaddr = nftaddress
		nfttab.Snftstage = nftaddress[:SnftStageOffset]
		nfttab.Snftcollection = nftaddress[:snftCollectionOffset]
		nfttab.Snft = nftaddress[:SnftOffset]
		//nfttab.Contract = strings.ToLower(ExchangAddr) //nft_contract_addr
		nfttab.Contract = contract
		nfttab.Createaddr = creator_addr
		nfttab.Ownaddr = owner_addr
		nfttab.Name = name
		nfttab.Desc = desc
		nfttab.Meta = meta
		//nfttab.Nftmeta = string(nftmetaJson)
		nfttab.Categories = categories
		nfttab.Collectcreator = collectRec.Createaddr
		nfttab.Collections = collections
		nfttab.Url = source_url
		//nfttab.Image = asset_sample
		nfttab.Md5 = md5
		nfttab.Selltype = SellTypeNotSale.String()
		nfttab.Count = count
		nfttab.Verified = Passed.String()
		nfttab.Verifiedtime = time.Now().Unix()

		nfttab.Mintstate = Minted.String()
		nfttab.Createdate = time.Now().Unix()
		nfttab.Royalty, _ = strconv.Atoi(royalty)
		sysInfo := SysInfos{}
		dberr := nft.db.Model(&SysInfos{}).Last(&sysInfo)
		if dberr.Error != nil {
			if dberr.Error != gorm.ErrRecordNotFound {
				log.Println("UploadWNft() SysInfos err=", dberr)
				return ErrCollectionNotExist
			}
		}
		log.Println("UploadWNft() SysInfos snfttotal count=", sysInfo.Snfttotal)
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Create(&nfttab)
			if err.Error != nil {
				fmt.Println("UploadWNft() err=", err.Error)
				return err.Error
			}
			//if collectRec.Snftcollection != snftCollection {
			err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
				nftInfo.CollectionsCreator, collections).Update("totalcount", collectRec.Totalcount+1)
			if err.Error != nil {
				fmt.Println("UploadWNft() add collectins totalcount err= ", err.Error)
				return err.Error
			}
			//}
			if nftaddress[len(nftaddress)-2:] == "00" {
				NftCatch.SetFlushFlag()
				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal",sysInfo.Snfttotal + 1)
				if err.Error != nil {
					fmt.Println("UploadWNft() add  SysInfos snfttotal err=", err.Error)
					return err.Error
				}
			}
			return nil
		})
	}

	return nil
}
