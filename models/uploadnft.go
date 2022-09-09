package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"io/ioutil"
	"log"
	"path"

	//"github.com/nftexchange/nftserver/ethhelper"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) UploadNft(
	user_addr string,
	creator_addr string,
	owner_addr string,
	md5 string,
	name string,
	desc string,
	meta string,
	source_url string,
	nft_contract_addr string,
	nft_token_id string,
	categories string,
	collections string,
	asset_sample string,
	hide string,
	royalty string,
	count string,
	sig string) error {

	user_addr = strings.ToLower(user_addr)
	creator_addr = strings.ToLower(creator_addr)
	owner_addr = strings.ToLower(owner_addr)
	nft_contract_addr = strings.ToLower(nft_contract_addr)
	spendT := time.Now()
	fmt.Println("UploadNft() user_addr=", user_addr, "      time=", time.Now().String())
	UserSync.Lock(user_addr)
	defer UserSync.UnLock(user_addr)
	fmt.Printf("UploadNft() UserSync.Lock Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	fmt.Println("UploadNft() begin ->> time = ", time.Now().String()[:22])
	fmt.Println("UploadNft() user_addr = ", user_addr)
	fmt.Println("UploadNft() creator_addr = ", creator_addr)
	fmt.Println("UploadNft() owner_addr = ", owner_addr)
	fmt.Println("UploadNft() md5 = ", md5)
	fmt.Println("UploadNft() name = ", name)
	fmt.Println("UploadNft() desc = ", desc)
	fmt.Println("UploadNft() meta = ", meta)
	fmt.Println("UploadNft() source_url = ", source_url)
	fmt.Println("UploadNft() nft_contract_addr = ", nft_contract_addr)
	fmt.Println("UploadNft() nft_token_id = ", nft_token_id)
	fmt.Println("UploadNft() categories = ", categories)
	fmt.Println("UploadNft() collections = ", collections)
	//fmt.Println("UploadNft() asset_sample = ", asset_sample)
	fmt.Println("UploadNft() hide = ", hide)
	fmt.Println("UploadNft() royalty = ", royalty)
	//fmt.Println("UploadNft() sig = ", sig)

	if IsIntDataValid(count) != true {
		return ErrDataFormat
	}
	if IsIntDataValid(royalty) != true {
		return ErrDataFormat
	}
	if !nft.UserKYCAduit(user_addr) {
		return ErrUserNotVerify
	}
	r, _ := strconv.Atoi(royalty)
	fmt.Println("UploadNft() royalty=", r, "SysRoyaltylimit=", SysRoyaltylimit, "RoyaltyLimit", RoyaltyLimit)
	if r > SysRoyaltylimit || r > RoyaltyLimit {
		return ErrRoyalty
	}
	if count == "" {
		count = "1"
	}
	if c, _ := strconv.Atoi(count); c < 1 {
		fmt.Println("UploadNft() contract count < 1.")
		return ErrContractCountLtZero
	}
	if nft.IsValidCategory(categories) {
		return ErrNoCategory
	}

	var collectRec Collects
	if collections != "" {
		err := nft.db.Model(&Collects{}).Select([]string{"id", "contract", "createaddr", "totalcount"}).Where("createaddr = ? AND name =?",
			creator_addr, collections).First(&collectRec)
		if err.Error != nil {
			fmt.Println("UploadNft() err=Collection not exist.")
			return ErrCollectionNotExist
		}
	} else {
		return ErrCollectionNotExist
	}
	fmt.Printf("UploadNft() preprocess Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	if nft_contract_addr == "" && nft_token_id == "" {
		var NewTokenid string
		rand.Seed(time.Now().UnixNano())
		var i int
		for i = 0; i < genTokenIdRetry; i++ {
			s := fmt.Sprintf("%d", rand.Int63())
			if len(s) < 15 {
				continue
			}
			s = s[len(s)-13:]
			NewTokenid = s
			if s[0] == '0' {
				continue
			}
			fmt.Println("UploadNft() NewTokenid=", NewTokenid)
			spendT = time.Now()
			nfttab := Nfts{}
			err := nft.db.Model(&Nfts{}).Select("id").Where("tokenid = ?", NewTokenid).First(&nfttab)
			if err.Error == gorm.ErrRecordNotFound {
				fmt.Printf("UploadNft() Nfts{} Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
				break
			}
			fmt.Println("UploadNft() Tokenid repetition.", NewTokenid)
		}
		if i >= 20 {
			fmt.Println("UploadNft() generate tokenId error.")
			return ErrGenerateTokenId
		}
		spendT = time.Now()
		imagerr := SaveNftImage(ImageDir, collectRec.Contract, NewTokenid, asset_sample)
		if imagerr != nil {
			fmt.Println("UploadNft() save image err=", imagerr)
			return ErrNftImage
		}
		fmt.Printf("UploadNft() SaveNftImage Spend time=%s filesize=%d time.now=%s\n", time.Now().Sub(spendT), len(asset_sample), time.Now())
		nftmeta := contracts.NftMeta{}
		nftmeta.Meta = meta
		nftmeta.TokenId = NewTokenid
		nftmetaJson, _ := json.Marshal(nftmeta)
		nfttab := Nfts{}
		nfttab.Tokenid = NewTokenid
		//nfttab.Contract = strings.ToLower(ExchangAddr) //nft_contract_addr
		nfttab.Contract = collectRec.Contract
		nfttab.Createaddr = creator_addr
		nfttab.Ownaddr = owner_addr
		nfttab.Name = name
		nfttab.Desc = desc
		nfttab.Meta = meta
		nfttab.Nftmeta = string(nftmetaJson)
		nfttab.Categories = categories
		nfttab.Collectcreator = collectRec.Createaddr
		nfttab.Collections = collections
		nfttab.Signdata = sig
		nfttab.Url = source_url
		//nfttab.Image = asset_sample
		nfttab.Md5 = md5
		nfttab.Selltype = SellTypeNotSale.String()

		if NFTUploadAuditRequired {
			nfttab.Verified = NoVerify.String()
		} else {
			nfttab.Verified = Passed.String()
			nfttab.Verifiedtime = time.Now().Unix()
		}
		nfttab.Mintstate = NoMinted.String()
		/*if collectRec.Contract == strings.ToLower(NFT1155Addr) {
			nfttab.Mintstate = NoMinted.String()
		} else {

			nfttab.Mintstate = Minted.String()
		}*/
		nfttab.Createdate = time.Now().Unix()
		nfttab.Royalty, _ = strconv.Atoi(royalty)
		//nfttab.Royalty /= 100
		nfttab.Count, _ = strconv.Atoi(count)
		nfttab.Hide = hide
		//err0, approve := ethhelper.GenCreateNftSign(NFT1155Addr, nfttab.Ownaddr, nfttab.Meta,
		//	nfttab.Tokenid, count, royalty)
		//if err0 != nil {
		//	fmt.Println("UploadNft() GenCreateNftSign() err=", err0)
		//	return err0
		//}
		//MintSign(contract string, toAddr string, tokenId string, count string, royalty string, tokenUri string, prv string)
		/*approve, err0 := contracts.MintSign(nfttab.Contract, nfttab.Ownaddr,
			nfttab.Tokenid, count, royalty, "", contracts.AdminMintPrv)
		if err0 != nil {
			fmt.Println("UploadNft() MintSign() err=", err0)
			return err0
		}
		fmt.Println("UploadNft() MintSign() approve=", approve)
		nfttab.Approve = approve*/
		spendT = time.Now()
		sysInfo := SysInfos{}
		dberr := nft.db.Model(&SysInfos{}).Last(&sysInfo)
		if dberr.Error != nil {
			if dberr.Error != gorm.ErrRecordNotFound {
				log.Println("UploadNft() SysInfos err=", dberr)
				return ErrCollectionNotExist
			}
			dberr = nft.db.Model(&SysInfos{}).Create(&sysInfo)
			if dberr.Error != nil {
				log.Println("UploadNft() SysInfos create err=", dberr)
				return ErrCollectionNotExist
			}
		}
		fmt.Println("UploadNft() SysInfos nfttotal count=", sysInfo.Nfttotal)
		err := nft.db.Transaction(func(tx *gorm.DB) error {
			spendT := time.Now()
			err := tx.Model(&Nfts{}).Create(&nfttab)
			if err.Error != nil {
				fmt.Println("UploadNft() err=", err.Error)
				return ErrDataBase
			}
			fmt.Printf("UploadNft() Nfts Create record Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			spendT = time.Now()
			if collections != "" {
				var collectListRec CollectLists
				collectListRec.Collectsid = collectRec.ID
				collectListRec.Nftid = nfttab.ID
				err = tx.Model(&CollectLists{}).Create(&collectListRec)
				if err.Error != nil {
					fmt.Println("UploadNft() create CollectLists err=", err.Error)
					return ErrDataBase
				}
				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
					collections, creator_addr).Update("totalcount", collectRec.Totalcount+1)
				if err.Error != nil {
					fmt.Println("UploadNft() add collectins totalcount err= ", err.Error)
					return ErrDataBase
				}
			}
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("nfttotal", sysInfo.Nfttotal+1)
			if err.Error != nil {
				fmt.Println("UploadNft() add  SysInfos nfttotal err=", err.Error)
				return ErrDataBase
			}
			fmt.Printf("UploadNft() collections Create record Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			//HomePageCatchs.NftCountLock()
			//HomePageCatchs.NftCountFlag = true
			//HomePageCatchs.NftCountUnLock()
			return nil
		})
		//NftCatch.SetFlushFlag()
		GetRedisCatch().SetDirtyFlag(UploadNftDirtyName)

		fmt.Printf("UploadNft() Create record Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return err
	} else {
		/*var nfttab Nfts
		dberr := nft.db.Where("contract = ? AND tokenid = ? ", nft_contract_addr, nft_token_id).First(&nfttab)
		if dberr.Error == nil {
			fmt.Println("UploadNft() err=nft already exist.")
			return ErrNftAlreadyExist
		}*/
		/*ownAddr, royalty, err := func(contract, tokenId string) (string, string, error) {
			return "ownAddr", "200", nil
		}(nft_contract_addr, nft_token_id)
		if ownAddr == user_addr {
			var nfttab Nfts
			nfttab.Tokenid = nft_token_id
			nfttab.Contract = nft_contract_addr //nft_contract_addr
			nfttab.Createaddr = creator_addr
			nfttab.Ownaddr = ownAddr
			nfttab.Name = name
			nfttab.Desc = desc
			nfttab.Meta = meta
			nfttab.Categories = categories
			nfttab.Collections = collections
			nfttab.Signdata = sig
			nfttab.Url = source_url
			nfttab.Image = asset_sample
			nfttab.Md5 = md5
			nfttab.Selltype = SellTypeNotSale.String()
			nfttab.Verified = NoVerify.String()
			nfttab.Mintstate = Minted.String()
			nfttab.Royalty, _ = strconv.Atoi(royalty)
			nfttab.Royalty = nfttab.Royalty / 100
			nfttab.Createdate = time.Now().Unix()
			nfttab.Hide = hide
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&Nfts{}).Create(&nfttab)
				if err.Error != nil {
					fmt.Println("UploadNft() create exist nft err=", err.Error)
					return err.Error
				}
				if collections != "" {
					var collectListRec CollectLists
					collectListRec.Collectsid = collectRec.ID
					collectListRec.Nftid = nfttab.ID
					err = tx.Model(&CollectLists{}).Create(&collectListRec)
					if err.Error != nil {
						fmt.Println("UploadNft() create CollectLists err=", err.Error)
						return err.Error
					}
				}
				return nil
			})
		}*/
		//IsAdminAddr, err := IsAdminAddr(user_addr)
		//if err != nil {
		//	fmt.Println("UploadNft() upload address is not admin.")
		//	return ErrNftUpAddrNotAdmin
		//}
		//if IsAdminAddr {
		//	var nfttab Nfts
		//	nfttab.Tokenid = nft_token_id
		//	nfttab.Contract = nft_contract_addr //nft_contract_addr
		//	nfttab.Createaddr = creator_addr
		//	nfttab.Ownaddr = owner_addr
		//	nfttab.Name = name
		//	nfttab.Desc = desc
		//	nfttab.Meta = meta
		//	nfttab.Categories = categories
		//	nfttab.Collectcreator = creator_addr
		//	nfttab.Collections = collections
		//	nfttab.Signdata = sig
		//	nfttab.Url = source_url
		//	//nfttab.Image = asset_sample
		//	nfttab.Md5 = md5
		//	nfttab.Selltype = SellTypeNotSale.String()
		//	nfttab.Verified = Passed.String()
		//	nfttab.Mintstate = Minted.String()
		//	/*nfttab.Royalty, _ = strconv.Atoi(royalty)
		//	nfttab.Royalty = nfttab.Royalty / 100*/
		//	nfttab.Createdate = time.Now().Unix()
		//	nfttab.Royalty, _ = strconv.Atoi(royalty)
		//	//nfttab.Royalty /= 100
		//	nfttab.Count, _ = strconv.Atoi(count)
		//	nfttab.Hide = hide
		//	return nft.db.Transaction(func(tx *gorm.DB) error {
		//		err := tx.Model(&Nfts{}).Create(&nfttab)
		//		if err.Error != nil {
		//			fmt.Println("UploadNft() admin create nft err=", err.Error)
		//			return err.Error
		//		}
		//		if collections != "" {
		//			/*var collectListRec CollectLists
		//			collectListRec.Collectsid = collectRec.ID
		//			collectListRec.Nftid = nfttab.ID
		//			err = tx.Model(&CollectLists{}).Create(&collectListRec)
		//			if err.Error != nil {
		//				fmt.Println("UploadNft() create CollectLists err=", err.Error)
		//				return err.Error
		//			}*/
		//			err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
		//				collections, creator_addr).Update("totalCount", collectRec.Totalcount+1)
		//			if err.Error != nil {
		//				fmt.Println("UploadNft() add collectins totalcount err= ", err.Error)
		//				return err.Error
		//			}
		//		}
		//		return nil
		//	})
		//} else {
		//	fmt.Println("UploadNft() upload address is not admin.")
		//	return ErrNftUpAddrNotAdmin
		//}
		fmt.Println("UploadNft() upload address is not admin.")
		return ErrData
	}
	return nil
}

func (nft NftDb) DelNft(useraddr, contract, tokenid string) error {

	nfts := Nfts{}
	err := nft.db.Model(&Nfts{}).Where("tokenid =? and contract=? and createaddr=?", tokenid, contract, useraddr).First(&nfts)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("DelNft() delete subscribe record err=", err.Error)
			return ErrDataBase
		}
		fmt.Println("DelNft() RecordNotFound")
		return ErrNotFound
	}
	if nfts.Mintstate != NoMinted.String() {
		fmt.Println("DelNft() delete nft Minstate cannot deleted")
		return ErrDeleteNft
	}
	rerr := nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Nfts{}).Where("tokenid=?", tokenid).Delete(&Nfts{})
		if err.Error != nil {
			fmt.Println("DelNft() delete subscribe record err=", err.Error)
			return ErrDataBase
		}
		err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
			nfts.Collections, nfts.Collectcreator).Update("totalcount", gorm.Expr("totalcount - ?", 1))
		if err.Error != nil {
			fmt.Println("DelNft() add collectins totalcount err= ", err.Error)
			return ErrDataBase
		}
		sysInfo := SysInfos{}
		err = nft.db.Model(&SysInfos{}).Last(&sysInfo)
		if err.Error != nil {
			if err.Error != gorm.ErrRecordNotFound {
				log.Println("DelNft() SysInfos err=", err)
				return ErrDataBase
			}
			err = nft.db.Model(&SysInfos{}).Create(&sysInfo)
			if err.Error != nil {
				log.Println("DelNft() SysInfos create err=", err)
				return ErrDataBase
			}
		}
		err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("nfttotal", sysInfo.Nfttotal-1)
		if err.Error != nil {
			fmt.Println("DelNft() add  SysInfos nfttotal err=", err.Error)
			return ErrNotExist
		}
		err = tx.Model(&NftFavorited{}).Where("tokenid = ?", tokenid).Delete(&NftFavorited{})
		if err.Error != nil {
			fmt.Println("DelNft() del  NftFavorited  err=", err.Error)
			return ErrNotExist
		}
		GetRedisCatch().SetDirtyFlag(UploadNftDirtyName)
		return nil
	})
	if rerr != nil {
		log.Println("DelNft Transaction err =", rerr)
		return rerr
	}
	homeerr := HomePageRenew()
	if homeerr != nil {
		log.Println("DelNft() HomePageRenew err=", homeerr)
		return homeerr
	}
	return nil
}

func (nft NftDb) SetNft(
	user_addr string,
	md5 string,
	name string,
	desc string,
	SourceUrl string,
	imageName string,
	nft_token_id string,
	categories string,
	collections string,
	asset_sample string,
	hide string,
	royalty string,
	count string,
	attributes,
	sig string) (NftImage, error) {
	user_addr = strings.ToLower(user_addr)
	spendT := time.Now()
	fmt.Println("SetNft() user_addr=", user_addr, "      time=", time.Now().String())
	UserSync.Lock(user_addr)
	defer UserSync.UnLock(user_addr)
	fmt.Printf("SetNft() UserSync.Lock Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	fmt.Println("SetNft() begin ->> time = ", time.Now().String()[:22])
	fmt.Println("SetNft() user_addr = ", user_addr)
	fmt.Println("SetNft() md5 = ", md5)
	fmt.Println("SetNft() name = ", name)
	fmt.Println("SetNft() desc = ", desc)
	fmt.Println("SetNft() SourceUrl = ", SourceUrl)
	fmt.Println("SetNft() imageName = ", imageName)
	fmt.Println("SetNft() nft_token_id = ", nft_token_id)
	fmt.Println("SetNft() categories = ", categories)
	fmt.Println("SetNft() collections = ", collections)
	//fmt.Println("UploadNftImage() asset_sample = ", asset_sample)
	fmt.Println("SetNft() hide = ", hide)
	fmt.Println("SetNft() royalty = ", royalty)
	//fmt.Println("UploadNftImage() sig = ", sig)
	setnfts := Nfts{}
	err := nft.db.Model(&Nfts{}).Where("tokenid = ? ", nft_token_id).First(&setnfts)
	if err.Error != nil {
		log.Println("SetNft() err = Nft not exist.")
		return NftImage{}, ErrNftNotExist
	}
	oldnfts := setnfts
	if setnfts.Mintstate != NoMinted.String() {
		log.Println("SetNft() err = Nft sell type not exist.")
		return NftImage{}, ErrNftNotExist
	}

	if IsIntDataValid(count) != true {
		return NftImage{}, ErrDataFormat
	}
	if IsIntDataValid(royalty) != true {
		return NftImage{}, ErrDataFormat
	}
	if !nft.UserKYCAduit(user_addr) {
		return NftImage{}, ErrUserNotVerify
	}
	r, _ := strconv.Atoi(royalty)
	fmt.Println("SetNft() royalty=", r, "SysRoyaltylimit=", SysRoyaltylimit, "RoyaltyLimit", RoyaltyLimit)
	if r > SysRoyaltylimit || r > RoyaltyLimit {
		return NftImage{}, ErrRoyalty
	}
	if count == "" {
		count = "1"
	}
	if c, _ := strconv.Atoi(count); c < 1 {
		fmt.Println("SetNft() contract count < 1.")
		return NftImage{}, ErrContractCountLtZero
	}
	if nft.IsValidCategory(categories) {
		return NftImage{}, ErrNoCategory
	}

	var collectRec Collects
	if collections != "" {
		err = nft.db.Model(&Collects{}).Where("createaddr = ? AND name =?",
			user_addr, collections).First(&collectRec)
		if err.Error != nil {
			log.Println("SetNft() err=Collection not exist.")
			return NftImage{}, ErrCollectionNotExist
		}

	} else {
		return NftImage{}, ErrCollectionNotExist
	}
	collectImageUrl, serr := SaveToIpfs(collectRec.Img)
	if serr != nil {
		log.Println("SetNft() save collection image err=", serr)
		return NftImage{}, ErrIpfsImage
	}
	fmt.Printf("SetNft() preprocess Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	if SourceUrl != "" {
		spendT = time.Now()
		//imagerr := SaveNftImage(ImageDir, collectRec.Contract, nft_token_id, asset_sample)
		//if imagerr != nil {
		//	fmt.Println("UploadNftImage() save image err=", imagerr)
		//	return NftImage{}, ErrNftImage
		//}
		if asset_sample != "" {
			imagerr := SaveNftImage(ImageDir, collectRec.Contract, nft_token_id, asset_sample)
			if imagerr != nil {
				fmt.Println("UploadNftImage() save image err=", imagerr)
				return NftImage{}, ErrNftImage
			}
		}
		var nftMeta nftInfo
		nftMeta.CreatorAddr = user_addr
		nftMeta.Contract = collectRec.Contract
		nftMeta.Name = name
		nftMeta.Desc = desc
		nftMeta.Category = categories
		nftMeta.Royalty = royalty
		nftMeta.SourceImageName = imageName
		nftMeta.FileType = path.Ext(imageName)[1:]
		nftMeta.SourceUrl = "/ipfs/" + SourceUrl
		nftMeta.Md5 = md5
		nftMeta.CollectionsName = collectRec.Name
		nftMeta.CollectionsCreator = collectRec.Createaddr
		nftMeta.CollectionsExchanger = collectRec.Contract
		nftMeta.CollectionsCategory = collectRec.Categories
		nftMeta.CollectionsImgUrl = collectImageUrl
		nftMeta.Attributes = attributes
		metaStr, err := json.Marshal(&nftMeta)
		meta, serr := SaveToIpfs(string(metaStr))
		if serr != nil {
			log.Println("SetNft() save nftmeta info err=", serr)
			return NftImage{}, errors.New("UploadNftImage() save nftmeta info error.")
		}
		meta = "/ipfs/" + meta
		fmt.Printf("SetNft() SaveNftImage Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		nftmeta := contracts.NftMeta{}
		nftmeta.Meta = meta
		nftmeta.TokenId = nft_token_id
		nftmetaJson, _ := json.Marshal(nftmeta)

		setnfts.Name = name
		setnfts.Desc = desc
		setnfts.Meta = meta
		setnfts.Nftmeta = string(nftmetaJson)
		setnfts.Categories = categories
		setnfts.Collectcreator = collectRec.Createaddr
		setnfts.Collections = collections
		setnfts.Signdata = sig
		setnfts.Url = "/ipfs/" + SourceUrl
		//nfttab.Image = asset_sample
		setnfts.Md5 = md5

		if NFTUploadAuditRequired {
			setnfts.Verified = NoVerify.String()
		} else {
			setnfts.Verified = Passed.String()
			setnfts.Verifiedtime = time.Now().Unix()
		}
		setnfts.Royalty, _ = strconv.Atoi(royalty)
		setnfts.Count, _ = strconv.Atoi(count)
		setnfts.Hide = hide
		nftimage := NftImage{}
		nftimage.NftMeta = hex.EncodeToString(nftmetaJson)
		nftimage.Meta = meta

		spendT = time.Now()

		err = nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Where("tokenid=?", nft_token_id).Updates(&setnfts)
			if err.Error != nil {
				fmt.Println("SetNft() err=", err.Error)
				return ErrDataBase
			}
			err = tx.Model(&Collects{}).Where("name =? and createaddr=?", oldnfts.Collections, oldnfts.Collectcreator).Update("totalcount", gorm.Expr("totalcount - ?", 1))
			if err.Error != nil {
				fmt.Println("SetNft() update collection err=", err.Error)
				return ErrDataBase
			}
			err = tx.Model(&Collects{}).Where("name =? and createaddr=?", setnfts.Collections, setnfts.Collectcreator).Update("totalcount", gorm.Expr("totalcount + ?", 1))
			if err.Error != nil {
				fmt.Println("SetNft() update collection err=", err.Error)
				return ErrDataBase
			}
			//HomePageCatchs.NftCountLock()
			//HomePageCatchs.NftCountFlag = true
			//HomePageCatchs.NftCountUnLock()
			return nil
		})
		//NftCatch.SetFlushFlag()
		GetRedisCatch().SetDirtyFlag(SetNftDirtyName)

		fmt.Printf("SetNft() Create record Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftimage, err
	} else {
		spendT = time.Now()
		getMeta, gerr := GetNftInfoFromIPFSWithShell(setnfts.Meta)
		if gerr != nil {
			fmt.Println("SetNft() get meta err =", gerr)
			return NftImage{}, ErrIpfsImage
		}

		var nftMeta nftInfo
		nftMeta.CreatorAddr = user_addr
		nftMeta.Contract = collectRec.Contract
		nftMeta.Name = name
		nftMeta.Desc = desc
		nftMeta.Category = categories
		nftMeta.Royalty = royalty
		nftMeta.SourceImageName = getMeta.SourceImageName
		nftMeta.FileType = getMeta.FileType
		nftMeta.SourceUrl = getMeta.SourceUrl
		nftMeta.Md5 = getMeta.Md5
		nftMeta.CollectionsName = collectRec.Name
		nftMeta.CollectionsCreator = collectRec.Createaddr
		nftMeta.CollectionsExchanger = collectRec.Contract
		nftMeta.CollectionsCategory = collectRec.Categories
		nftMeta.CollectionsImgUrl = collectImageUrl
		nftMeta.Attributes = attributes
		metaStr, err := json.Marshal(&nftMeta)
		if err != nil {
			fmt.Println("SetNft() marshal nftmeta =", err)
			return NftImage{}, err
		}
		meta, serr := SaveToIpfs(string(metaStr))
		if serr != nil {
			log.Println("SetNft() save nftmeta info err=", serr)
			return NftImage{}, ErrIpfsImage
		}
		meta = "/ipfs/" + meta
		fmt.Printf("SetNft() SaveNftImage Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		nftmeta := contracts.NftMeta{}
		nftmeta.Meta = meta
		nftmeta.TokenId = nft_token_id
		nftmetaJson, _ := json.Marshal(nftmeta)

		setnfts.Name = name
		setnfts.Desc = desc
		setnfts.Meta = meta
		setnfts.Nftmeta = string(nftmetaJson)
		setnfts.Categories = categories
		setnfts.Collectcreator = collectRec.Createaddr
		setnfts.Collections = collections
		setnfts.Signdata = sig
		setnfts.Md5 = md5

		if NFTUploadAuditRequired {
			setnfts.Verified = NoVerify.String()
		} else {
			setnfts.Verified = Passed.String()
			setnfts.Verifiedtime = time.Now().Unix()
		}
		setnfts.Royalty, _ = strconv.Atoi(royalty)
		setnfts.Count, _ = strconv.Atoi(count)
		setnfts.Hide = hide
		nftimage := NftImage{}
		nftimage.NftMeta = hex.EncodeToString(nftmetaJson)
		nftimage.Meta = meta

		spendT = time.Now()
		err = nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Where("tokenid=?", nft_token_id).Updates(&setnfts)
			if err.Error != nil {
				fmt.Println("SetNft() err=", err.Error)
				return ErrDataBase
			}
			err = tx.Model(&Collects{}).Where("name =? and createaddr=?", oldnfts.Collections, oldnfts.Collectcreator).Update("totalcount", gorm.Expr("totalcount - ?", 1))
			if err.Error != nil {
				fmt.Println("SetNft() update collection err=", err.Error)
				return ErrDataBase
			}
			err = tx.Model(&Collects{}).Where("name =? and createaddr=?", setnfts.Collections, setnfts.Collectcreator).Update("totalcount", gorm.Expr("totalcount + ?", 1))
			if err.Error != nil {
				fmt.Println("SetNft() update collection err=", err.Error)
				return ErrDataBase
			}
			//HomePageCatchs.NftCountLock()
			//HomePageCatchs.NftCountFlag = true
			//HomePageCatchs.NftCountUnLock()
			return nil
		})
		//NftCatch.SetFlushFlag()
		GetRedisCatch().SetDirtyFlag(SetNftDirtyName)
		fmt.Printf("SetNft() Create record Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftimage, nil
	}
}

func GetNftInfoFromIPFSWithShell(hash string) (*nftInfo, error) {
	url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	//url = "http://192.168.1.235:5001"
	s := shell.NewShell(url)
	s.SetTimeout(100 * time.Second)
	rc, err := s.Cat(hash)
	if err != nil {
		log.Println("GetnftInfoFromIPFSWithShell() err=", err)
		return nil, err
	}
	var snft nftInfo
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Println("GetnftInfoFromIPFSWithShell() ReadAll() err=", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(b), &snft)
	if err != nil {
		log.Println("GetnftInfoFromIPFSWithShell() Unmarshal, err=", err)
		return nil, err
	}
	return &snft, nil
}
