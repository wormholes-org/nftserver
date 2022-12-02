package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func copySnftToSysnft(snft *Nfts) Sysnfts {
	sysNft := Sysnfts{}
	sysNft.Ownaddr = snft.Ownaddr
	sysNft.Md5 = snft.Md5
	sysNft.Name = snft.Name
	sysNft.Desc = snft.Desc
	sysNft.Meta = snft.Meta
	sysNft.Nftmeta = snft.Nftmeta
	sysNft.Url = snft.Url
	sysNft.Contract = snft.Contract
	sysNft.Tokenid = snft.Tokenid
	sysNft.Nftaddr = snft.Nftaddr
	sysNft.Snftstage = snft.Snftstage
	sysNft.Snftcollection = snft.Snftcollection
	sysNft.Snft = snft.Snft
	sysNft.Count = snft.Count
	sysNft.Approve = snft.Approve
	sysNft.Categories = snft.Categories
	sysNft.Collectcreator = snft.Collectcreator
	sysNft.Collections = snft.Collections
	sysNft.Image = snft.Image
	sysNft.Hide = snft.Hide
	sysNft.Signdata = snft.Signdata
	sysNft.Createaddr = snft.Createaddr
	sysNft.Verifyaddr = snft.Verifyaddr
	sysNft.Currency = snft.Currency
	sysNft.Price = snft.Price
	sysNft.Royalty = snft.Royalty
	sysNft.Paychan = snft.Paychan
	sysNft.TransCur = snft.TransCur
	sysNft.Transprice = snft.Transprice
	sysNft.Transtime = snft.Transtime
	sysNft.Createdate = snft.Createdate
	sysNft.Favorited = snft.Favorited
	sysNft.Transcnt = snft.Transcnt
	sysNft.Transamt = snft.Transamt
	sysNft.Verified = snft.Verified
	sysNft.Verifieddesc = snft.Verifieddesc
	sysNft.Verifiedtime = snft.Verifiedtime
	sysNft.Selltype = snft.Selltype
	sysNft.Mintstate = snft.Mintstate
	sysNft.Pledgestate = snft.Pledgestate
	return sysNft
}

func (nft NftDb) NewTokenId() (string, error) {
	rand.Seed(time.Now().UnixNano())
	var Tokenid string
	var i int
	for i = 0; i < genTokenIdRetry; i++ {
		//NewTokenid := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		s := fmt.Sprintf("%d", rand.Int63())
		if len(s) < 15 {
			continue
		}
		s = s[len(s)-13:]
		Tokenid = s
		if s[0] == '0' {
			continue
		}
		fmt.Println("UploadWNft() NewTokenid=", Tokenid)
		nfttab := Nfts{}
		err := nft.db.Where("tokenid = ? ", Tokenid).First(&nfttab)
		if err.Error == gorm.ErrRecordNotFound {
			break
		}
	}
	if i >= 20 {
		fmt.Println("UploadWNft() generate tokenId error.")
		return "", ErrGenerateTokenId
	}
	return Tokenid, nil
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
	royalty := strconv.FormatInt(int64(nftInfo.Royalty*100), 10)
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
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
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
			collectRec.Contract = /*nftInfo.CollectionsExchanger*/ "0x0000000000000000000000000000000000000000"
			collectRec.Contracttype = "snft"
			collectRec.Categories = nftInfo.CollectionsCategory

			err := nft.db.Model(&Collects{}).Create(&collectRec)
			if err.Error != nil {
				log.Println("UploadWNft create Collections() err=", err.Error)
				return err.Error
			}
		}
	}
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
			"Createaddr":     creator_addr,
			"Ownaddr":        owner_addr,
			"Name":           name,
			"Desc":           desc,
			"Meta":           meta,
			"Categories":     categories,
			"Collectcreator": collectRec.Createaddr,
			"Collections":    collections,
			"Url":            source_url,
			"Selltype":       SellTypeNotSale.String(),
			"Verifiedtime":   time.Now().Unix(),
			"Createdate":     time.Now().Unix(),
			"Transprice":     0,
			"Price":          0,
			"Transtime":      0,
			"Transcnt":       0,
			"Transamt":       0,
			"Favorited":      0,
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
			err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
				nftInfo.CollectionsCreator, collections).Update("totalcount", collectRec.Totalcount+1)
			if err.Error != nil {
				fmt.Println("UploadWNft() add collectins totalcount err= ", err.Error)
				return err.Error
			}
			return nil
		})
	} else {
		var nerr error
		NewTokenid, nerr := nft.NewTokenId()
		if nerr != nil {
			fmt.Println("UploadWNft() generate tokenid err= ", nerr)
			return nerr
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
		nfttab.Pledgestate = NoPledge.String()
		nfttab.Createdate = time.Now().Unix()
		nfttab.Royalty, _ = strconv.Atoi(royalty)
		msnft := Nfts{}
		err := nft.db.Select([]string{"id", "Chipcount"}).Where("nftaddr = ?", nftaddress[:len(nftaddress)-1]+"m").First(&msnft)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("UploadWNft() database err=", err.Error)
			return ErrCollectionNotExist
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
		sysNft := Sysnfts{}
		err = nft.db.Where("snft = ?", nfttab.Snft).First(&sysNft)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			log.Println("UploadWNft() database err=", err.Error)
			return ErrCollectionNotExist
		}
		if err.Error == gorm.ErrRecordNotFound {
			sysNft = copySnftToSysnft(&nfttab)
			err := nft.db.Model(&Sysnfts{}).Create(&sysNft)
			if err.Error != nil {
				log.Println("UploadWNft() err=", err.Error)
				return err.Error
			}
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Create(&nfttab)
			if err.Error != nil {
				log.Println("UploadWNft() err=", err.Error)
				return err.Error
			}
			if nftaddress[len(nftaddress)-1:] == "0" {
				nfttab.ID = 0
				Tokenid, nerr := nft.NewTokenId()
				if nerr != nil {
					fmt.Println("UploadWNft() generate tokenid err= ", nerr)
					return nerr
				}
				nfttab.Tokenid = Tokenid
				nfttab.Chipcount = 1
				nfttab.Mergetype = 1
				nfttab.Ownaddr = ZeroAddr
				nfttab.Nftaddr = nftaddress[:len(nftaddress)-1] + "m"
				nfttab.Snftstage = nftaddress[:SnftStageOffset] + "m"
				nfttab.Snftcollection = nftaddress[:snftCollectionOffset] + "m"
				nfttab.Snft = nftaddress[:SnftOffset] + "m"
				err = tx.Model(&Nfts{}).Create(&nfttab)
				if err.Error != nil {
					log.Println("UploadWNft() err=", err.Error)
					return err.Error
				}
			}
			if nftaddress[len(nftaddress)-2:] == "00" {
				nfttab.ID = 0
				Tokenid, nerr := nft.NewTokenId()
				if nerr != nil {
					fmt.Println("UploadWNft() generate tokenid err= ", nerr)
					return nerr
				}
				nfttab.Tokenid = Tokenid
				nfttab.Chipcount = 256
				nfttab.Mergetype = 2
				nfttab.Ownaddr = ZeroAddr
				nfttab.Nftaddr = nftaddress[:len(nftaddress)-2] + "mm"
				nfttab.Snftstage = nftaddress[:SnftStageOffset] + "m"
				nfttab.Snftcollection = nftaddress[:snftCollectionOffset] + "m"
				nfttab.Snft = nftaddress[:SnftOffset] + "m"
				err = tx.Model(&Nfts{}).Create(&nfttab)
				if err.Error != nil {
					log.Println("UploadWNft() err=", err.Error)
					return err.Error
				}
			}
			if nftaddress[len(nftaddress)-3:] == "000" {
				nfttab.ID = 0
				Tokenid, nerr := nft.NewTokenId()
				if nerr != nil {
					fmt.Println("UploadWNft() generate tokenid err= ", nerr)
					return nerr
				}
				nfttab.Tokenid = Tokenid
				nfttab.Chipcount = 4096
				nfttab.Mergetype = 3
				nfttab.Ownaddr = ZeroAddr
				nfttab.Nftaddr = nftaddress[:len(nftaddress)-3] + "mmm"
				nfttab.Snftstage = nftaddress[:SnftStageOffset] + "m"
				nfttab.Snftcollection = nftaddress[:snftCollectionOffset] + "m"
				nfttab.Snft = nftaddress[:SnftOffset] + "m"
				err = tx.Model(&Nfts{}).Create(&nfttab)
				if err.Error != nil {
					log.Println("UploadWNft() err=", err.Error)
					return err.Error
				}
			}
			if nftaddress[len(nftaddress)-1:] != "0" {
				err = tx.Model(&Nfts{}).Where("id = ?", msnft.ID).Update("chipcount", msnft.Chipcount+1)
				if err.Error != nil {
					log.Println("UploadWNft() nft totalcount err= ", err.Error)
					return err.Error
				}
			}
			err = tx.Model(&Collects{}).Where("createaddr = ? AND  name=?",
				nftInfo.CollectionsCreator, collections).Update("totalcount", collectRec.Totalcount+1)
			if err.Error != nil {
				log.Println("UploadWNft() add collectins totalcount err= ", err.Error)
				return err.Error
			}
			err = tx.Model(&Sysnfts{}).Where("id = ?", sysNft.ID).Update("Chipcount", sysNft.Chipcount+1)
			if err.Error != nil {
				log.Println("UploadWNft() add Sysnfts chip count err= ", err.Error)
				return err.Error
			}
			if (sysNft.Chipcount + 1) == 16 {
				err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("Snfttotal", sysInfo.Snfttotal+1)
				if err.Error != nil {
					fmt.Println("UploadWNft() add  SysInfos snfttotal err=", err.Error)
					return err.Error
				}
				GetRedisCatch().SetDirtyFlag(UploadNftDirtyName)
			}
			return nil
		})
	}
	return nil
}
