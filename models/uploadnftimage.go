package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"log"
	"path"

	//"github.com/nftexchange/nftserver/ethhelper"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type nftInfo struct {
	CreatorAddr          string `json:"creator_addr"`
	Contract             string `json:"nft_contract_addr"`
	Name                 string `json:"name"`
	Desc                 string `json:"desc"`
	Category             string `json:"category"`
	Royalty              string `json:"royalty"`
	SourceImageName      string `json:"source_image_name"`
	FileType             string `json:"fileType"`
	SourceUrl            string `json:"source_url"`
	Md5                  string `json:"md5"`
	CollectionsName      string `json:"collections_name"`
	CollectionsCreator   string `json:"collections_creator"`
	CollectionsExchanger string `json:"collections_exchanger"`
	CollectionsCategory  string `json:"collections_category"`
	CollectionsImgUrl    string `json:"collections_img_url"`
	CollectionsDesc      string `json:"collections_desc"`
	Attributes           string `json:"attributes"`
}

type NftImage struct {
	Tokenid string `json:"tokenid"`
	Meta    string `json:"meta"`
	NftMeta string `json:"nftmeta"`
}

func SaveToIpfs(str string) (string, error) {
	url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	spendT := time.Now()
	s := shell.NewShell(url)
	s.SetTimeout(500 * time.Second)
	mhash, err := s.Add(bytes.NewBufferString(str))
	if err != nil {
		log.Println("SaveToIpfs() err=", err)
		return "", err
	}
	fmt.Printf("SaveToIpfs  Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return mhash, nil
}

func (nft NftDb) UploadNftImage(
	user_addr string,
	creator_addr string,
	owner_addr string,
	md5 string,
	name string,
	desc string,
	SourceUrl string,
	imageName string,
	nft_contract_addr string,
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
	creator_addr = strings.ToLower(creator_addr)
	owner_addr = strings.ToLower(owner_addr)
	nft_contract_addr = strings.ToLower(nft_contract_addr)
	spendT := time.Now()
	fmt.Println("UploadNftImage() user_addr=", user_addr, "      time=", time.Now().String())
	UserSync.Lock(user_addr)
	defer UserSync.UnLock(user_addr)
	fmt.Printf("UploadNftImage() UserSync.Lock Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	fmt.Println("UploadNftImage() begin ->> time = ", time.Now().String()[:22])
	fmt.Println("UploadNftImage() user_addr = ", user_addr)
	fmt.Println("UploadNftImage() creator_addr = ", creator_addr)
	fmt.Println("UploadNftImage() owner_addr = ", owner_addr)
	fmt.Println("UploadNftImage() md5 = ", md5)
	fmt.Println("UploadNftImage() name = ", name)
	fmt.Println("UploadNftImage() desc = ", desc)
	fmt.Println("UploadNftImage() SourceUrl = ", SourceUrl)
	fmt.Println("UploadNftImage() imageName = ", imageName)
	fmt.Println("UploadNftImage() nft_contract_addr = ", nft_contract_addr)
	fmt.Println("UploadNftImage() nft_token_id = ", nft_token_id)
	fmt.Println("UploadNftImage() categories = ", categories)
	fmt.Println("UploadNftImage() collections = ", collections)
	//fmt.Println("UploadNftImage() asset_sample = ", asset_sample)
	fmt.Println("UploadNftImage() hide = ", hide)
	fmt.Println("UploadNftImage() royalty = ", royalty)
	//fmt.Println("UploadNftImage() sig = ", sig)

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
	fmt.Println("UploadNftImage() royalty=", r, "SysRoyaltylimit=", SysRoyaltylimit, "RoyaltyLimit", RoyaltyLimit)
	if r > SysRoyaltylimit || r > RoyaltyLimit {
		return NftImage{}, ErrRoyalty
	}
	if count == "" {
		count = "1"
	}
	if c, _ := strconv.Atoi(count); c < 1 {
		fmt.Println("UploadNftImage() contract count < 1.")
		return NftImage{}, ErrContractCountLtZero
	}
	if nft.IsValidCategory(categories) {
		return NftImage{}, ErrNoCategory
	}

	var collectRec Collects
	if collections != "" {
		err := nft.db.Model(&Collects{}).Where("createaddr = ? AND name =?",
			creator_addr, collections).First(&collectRec)
		if err.Error != nil {
			log.Println("UploadNftImage() err=Collection not exist.")
			return NftImage{}, ErrCollectionNotExist
		}

	} else {
		return NftImage{}, ErrCollectionNotExist
	}
	collectImageUrl, serr := SaveToIpfs(collectRec.Img)
	if serr != nil {
		log.Println("UploadNftImage() save collection image err=", serr)
		return NftImage{}, errors.New("UploadNftImage() save collection image error.")
	}
	fmt.Printf("UploadNftImage() preprocess Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
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
			fmt.Println("UploadNftImage() NewTokenid=", NewTokenid)
			spendT = time.Now()
			nfttab := Nfts{}
			err := nft.db.Model(&Nfts{}).Select("id").Where("tokenid = ? ", NewTokenid).First(&nfttab)
			if err.Error == gorm.ErrRecordNotFound {
				fmt.Printf("UploadNftImage() Nfts{} Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
				break
			}
			fmt.Println("UploadNftImage() Tokenid repetition.", NewTokenid)
		}
		if i >= 20 {
			fmt.Println("UploadNftImage() generate tokenId error.")
			return NftImage{}, ErrGenerateTokenId
		}
		spendT = time.Now()
		imagerr := SaveNftImage(ImageDir, collectRec.Contract, NewTokenid, asset_sample)
		if imagerr != nil {
			fmt.Println("UploadNftImage() save image err=", imagerr)
			return NftImage{}, ErrNftImage
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
			log.Println("UploadNftImage() save nftmeta info err=", serr)
			return NftImage{}, errors.New("UploadNftImage() save nftmeta info error.")
		}
		meta = "/ipfs/" + meta
		fmt.Printf("UploadNftImage() SaveNftImage Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
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
		nfttab.Url = "/ipfs/" + SourceUrl
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
		nfttab.Createdate = time.Now().Unix()
		nfttab.Royalty, _ = strconv.Atoi(royalty)
		nfttab.Count, _ = strconv.Atoi(count)
		nfttab.Hide = hide
		nftimage := NftImage{}

		nftimage.NftMeta = hex.EncodeToString(nftmetaJson)
		nftimage.Meta = meta
		nftimage.Tokenid = NewTokenid
		//err0, approve := ethhelper.GenCreateNftSign(NFT1155Addr, nfttab.Ownaddr, nfttab.Meta,
		//	nfttab.Tokenid, count, royalty)
		//if err0 != nil {
		//	fmt.Println("UploadNftImage() GenCreateNftSign() err=", err0)
		//	return err0
		//}
		//MintSign(contract string, toAddr string, tokenId string, count string, royalty string, tokenUri string, prv string)
		/*approve, err0 := contracts.MintSign(nfttab.Contract, nfttab.Ownaddr,
			nfttab.Tokenid, count, royalty, "", contracts.AdminMintPrv)
		if err0 != nil {
			fmt.Println("UploadNftImage() MintSign() err=", err0)
			return err0
		}
		fmt.Println("UploadNftImage() MintSign() approve=", approve)
		nfttab.Approve = approve*/
		spendT = time.Now()
		sysInfo := SysInfos{}
		dberr := nft.db.Model(&SysInfos{}).Last(&sysInfo)
		if dberr.Error != nil {
			if dberr.Error != gorm.ErrRecordNotFound {
				log.Println("UploadNft() SysInfos err=", dberr)
				return NftImage{}, ErrCollectionNotExist
			}
			dberr = nft.db.Model(&SysInfos{}).Create(&sysInfo)
			if dberr.Error != nil {
				log.Println("UploadNft() SysInfos create err=", dberr)
				return NftImage{}, ErrCollectionNotExist
			}
		}
		fmt.Println("UploadNft() SysInfos nfttotal count=", sysInfo.Nfttotal)
		err = nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Create(&nfttab)
			if err.Error != nil {
				fmt.Println("UploadNftImage() err=", err.Error)
				return err.Error
			}
			if collections != "" {
				/*var collectListRec CollectLists
				collectListRec.Collectsid = collectRec.ID
				collectListRec.Nftid = nfttab.ID
				err = tx.Model(&CollectLists{}).Create(&collectListRec)
				if err.Error != nil {
					fmt.Println("UploadNftImage() create CollectLists err=", err.Error)
					return err.Error
				}*/
				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
					collections, creator_addr).Update("totalcount", collectRec.Totalcount+1)
				if err.Error != nil {
					fmt.Println("UploadNftImage() add collectins totalcount err= ", err.Error)
					return err.Error
				}
			}
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("nfttotal", sysInfo.Nfttotal+1)
			if err.Error != nil {
				fmt.Println("UploadNft() add  SysInfos nfttotal err=", err.Error)
				return err.Error
			}
			HomePageCatchs.NftCountLock()
			HomePageCatchs.NftCountFlag = true
			HomePageCatchs.NftCountUnLock()
			return nil
		})
		NftCatch.SetFlushFlag()
		fmt.Printf("UploadNftImage() Create record Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftimage, err
	} else {
		log.Println("UploadNftImage() contract!=NULL tokenid != 0 ")
		return NftImage{}, errors.New("UploadNftImage() contract!=NULL tokenid != 0")
	}
}
