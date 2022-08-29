package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) NewCollections(useraddr, name, img, contract_type, contract_addr,
	desc, categories, sig string) error {
	useraddr = strings.ToLower(useraddr)
	contract_addr = strings.ToLower(contract_addr)
	fmt.Println("NewCollections() user_addr=", useraddr, "      time=", time.Now().String())
	UserSync.Lock(useraddr)
	defer UserSync.UnLock(useraddr)
	//fmt.Println("NewCollections() useraddr=", useraddr )
	fmt.Println("NewCollections() contract_addr=", contract_addr)
	if !nft.UserKYCAduit(useraddr) {
		return ErrUserNotVerify
	}
	var collectRec Collects
	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
	if err.Error == nil {
		fmt.Println("NewCollections() err=Collection already exist.")
		return ErrCollectionExist
	} else if err.Error == gorm.ErrRecordNotFound {
		collectRec = Collects{}
		collectRec.Createaddr = useraddr
		collectRec.Name = name
		collectRec.Desc = desc
		//collectRec.Img = img
		if contract_addr != "" {
			collectRec.Contract = contract_addr
		} else {
			//collectRec.Contract = strings.ToLower(NFT1155Addr)
			collectRec.Contract = strings.ToLower(ExchangeOwer)
		}
		collectRec.Contracttype = contract_type
		collectRec.Categories = categories
		collectRec.SigData = sig
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Collects{}).Create(&collectRec)
			if err.Error != nil {
				fmt.Println("NewCollections() err=", err.Error)
				return errors.New(ErrDataBase.Error() + err.Error.Error())
			}
			imagerr := SaveCollectionsImage(ImageDir, useraddr, name, img)
			if imagerr != nil {
				fmt.Println("NewCollections() SaveCollectionsImage() err=", imagerr)
				return ErrNftImage
			}
			GetRedisCatch().SetDirtyFlag(CollectionList)
			return nil
		})
	}

	fmt.Println("NewCollections() dbase err=.", err)
	return errors.New(ErrDataBase.Error() + err.Error.Error())
}

func (nft NftDb) NewUserCollection(useraddr, name, img, contract_type, contract_addr,
	desc, categories, sig string) error {
	useraddr = strings.ToLower(useraddr)
	contract_addr = strings.ToLower(contract_addr)
	fmt.Println("NewUserCollection() user_addr=", useraddr, "      time=", time.Now().String())
	//fmt.Println("NewCollections() useraddr=", useraddr )
	fmt.Println("NewUserCollection() contract_addr=", contract_addr)
	if !nft.UserKYCAduit(useraddr) {
		return ErrUserNotVerify
	}
	var collectRec Collects
	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
	if err.Error == nil {
		fmt.Println("NewUserCollection() err=Collection already exist.")
		return ErrCollectionExist
	} else if err.Error == gorm.ErrRecordNotFound {
		collectRec = Collects{}
		collectRec.Createaddr = useraddr
		collectRec.Name = name
		collectRec.Desc = desc
		//collectRec.Img = img
		if contract_addr != "" {
			collectRec.Contract = contract_addr
		} else {
			//collectRec.Contract = strings.ToLower(NFT1155Addr)
			collectRec.Contract = strings.ToLower(ExchangeOwer)
		}
		collectRec.Contracttype = contract_type
		collectRec.Categories = categories
		collectRec.SigData = sig
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Collects{}).Create(&collectRec)
			if err.Error != nil {
				fmt.Println("NewUserCollection() err=", err.Error)
				return errors.New(ErrDataBase.Error() + err.Error.Error())
			}
			imagerr := SaveCollectionsImage(ImageDir, useraddr, name, img)
			if imagerr != nil {
				fmt.Println("NewUserCollection() SaveCollectionsImage() err=", imagerr)
				return ErrNftImage
			}
			GetRedisCatch().SetDirtyFlag(CollectionList)
			return nil
		})
	}

	fmt.Println("NewCollections() dbase err=.", err)
	return errors.New(ErrDataBase.Error() + err.Error.Error())
}

func (nft NftDb) ModifyCollections(useraddr, name, img, contract_type, contract_addr,
	desc, categories, sig string) error {
	useraddr = strings.ToLower(useraddr)
	contract_addr = strings.ToLower(contract_addr)
	if !nft.UserKYCAduit(useraddr) {
		return ErrUserNotVerify
	}
	var collectRec Collects
	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
	if err.Error != nil {
		fmt.Println("NewCollections() err=Collection not exist.")
		return ErrCollectionNotExist
	}
	collectRec = Collects{}
	if img != "" {
		collectRec.Img = img
		imagerr := SaveCollectionsImage(ImageDir, useraddr, name, img)
		if imagerr != nil {
			fmt.Println("ModifyCollections() SaveCollectionsImage() err=", imagerr)
			return ErrNftImage
		}

	}
	if contract_type != "" {
		collectRec.Contracttype = contract_type
	}
	if contract_addr != "" {
		collectRec.Contract = contract_addr
	}
	if desc != "" {
		collectRec.Desc = desc
	}
	if categories != "" {
		collectRec.Categories = categories
	}
	if sig != "" {
		collectRec.SigData = sig
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Collects{}).Where("Createaddr = ? AND name = ? ", useraddr, name).Updates(&collectRec)
		if err.Error != nil {
			fmt.Println("NewCollections() err=", err.Error)
			return errors.New(ErrDataBase.Error() + err.Error.Error())
		}
		GetRedisCatch().SetDirtyFlag(CollectionList)

		return nil
	})
}

type NFTCollectionListCatch struct {
	UserCollections []UserCollection
	Total           int
}

func (nft NftDb) QueryNFTCollectionList(start_index, count string) ([]UserCollection, int, error) {
	var collectRecs []Collects
	var recCount int64
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	spendT := time.Now()
	queryCatchSql := start_index + count
	nftCatch := NFTCollectionListCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryNFTCollectionList", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryNFTCollectionList() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.UserCollections, nftCatch.Total, nil
	}
	err := nft.db.Model(Collects{}).Where("totalcount > 0").Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryNFTCollectionList() recCount err=", err)
		return nil, 0, ErrNftNotExist
	}
	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)
	if int64(startIndex) >= recCount || recCount == 0 {
		return nil, 0, ErrNotMore
	} else {
		temp := recCount - int64(startIndex)
		if int64(nftCount) > temp {
			nftCount = int(temp)
		}
		err = nft.db.Model(Collects{}).Select([]string{"createaddr", "contract", "Contracttype", "name", "img", "totalcount", "transcnt"}).Where("totalcount > 0").Order("transamt desc, id desc").Limit(nftCount).Offset(startIndex).Find(&collectRecs)
		if err.Error != nil {
			fmt.Println("QueryNFTCollectionList() find record err=", err)
			return nil, 0, ErrNftNotExist
		}

		userCollects := make([]UserCollection, 0, 20)
		//var collectlist, collectaddr []string
		for i := 0; i < len(collectRecs); i++ {
			var userCollect UserCollection
			userCollect.CreatorAddr = collectRecs[i].Createaddr
			userCollect.Name = collectRecs[i].Name
			if collectRecs[i].Contracttype == "snft" {
				userCollect.Img = collectRecs[i].Img
			}
			//userCollect.Img = collectRecs[i].Img
			userCollect.ContractAddr = collectRecs[i].Contract
			userCollect.Desc = collectRecs[i].Desc
			//userCollect.Royalty = collectRecs[i].Royalty
			userCollect.Categories = collectRecs[i].Categories
			userCollect.Contracttype = collectRecs[i].Contracttype
			userCollect.Totalcount = collectRecs[i].Totalcount
			userCollect.Transcnt = collectRecs[i].Transcnt
			userCollects = append(userCollects, userCollect)
			//collectlist = append(collectlist, collectRecs[i].Name)
			//collectaddr = append(collectaddr, collectRecs[i].Createaddr)
		}

		//var tran []Tranhistory
		//err = nft.db.Table("nfts").Select("nfts.collections ,nfts.collectcreator,trans.txhash").
		//	Joins("left join trans on trans.selltype != ? and trans.selltype != ?  and trans.tokenid =nfts.tokenid and trans.deleted_at is null",
		//		SellTypeError.String(), SellTypeMintNft.String()).
		//	Where("nfts.collections in ? and  nfts.collectcreator  in  ?", collectlist, collectaddr).Find(&tran)
		//if err.Error != nil {
		//	fmt.Println("QueryNFTCollectionList() find trans err=", err)
		//	return nil, 0, err.Error
		//}
		//for _, v := range tran {
		//	if v.Txhash == "" {
		//		continue
		//	}
		//	for i, j := range userCollects {
		//		if j.Name == v.Collections && j.CreatorAddr == v.Collectcreator {
		//			userCollects[i].Transcount++
		//			break
		//		}
		//		continue
		//	}
		//}
		GetRedisCatch().CatchQueryData("QueryNFTCollectionList", queryCatchSql, &NFTCollectionListCatch{userCollects, int(recCount)})
		log.Printf("QueryNFTCollectionList() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return userCollects, int(recCount), nil
	}
}

func (nft NftDb) DelCollection(useraddr, contract, name string) error {

	collect := Collects{}
	err := nft.db.Model(&Collects{}).Where("name =? and contract=? and createaddr=?", name, contract, useraddr).First(&collect)
	if err.Error != nil {
		fmt.Println("DelCollection() RecordNotFound ,err=", err)
		return errors.New(ErrNotFound.Error() + err.Error.Error())
	}
	nfts := Nfts{}
	err = nft.db.Model(&Nfts{}).Where("collectcreator =? and collections=? and mintstate <> ? ", useraddr, name, "NoMinted").First(&nfts)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("DelCollection() delete subscribe record err=", err.Error)
			return errors.New(ErrDataBase.Error() + err.Error.Error())
		}
		rerr := nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Collects{}).Where("name =? and contract=? and createaddr=?", name, contract, useraddr).Delete(&Collects{})
			if err.Error != nil {
				fmt.Println("DelCollection() delete subscribe record err=", err.Error)
				return errors.New(ErrDataBase.Error() + err.Error.Error())
			}
			var total int64
			//err = tx.Model(&Nfts{}).Count(&total).Where("collectcreator =? and collections=? and mintstate =? ", useraddr, name, "NoMinted").Find(&Nfts{})
			//if err.Error != nil {
			//	fmt.Println("DelCollection() count nft err= ", err.Error)
			//	return err.Error
			//}
			err = tx.Model(&Nfts{}).Where("collectcreator =? and collections=? and mintstate =? ", useraddr, name, "NoMinted").Count(&total).Delete(&Nfts{})
			if err.Error != nil {
				fmt.Println("DelCollection() delete collection under nfts err= ", err.Error)
				return errors.New(ErrDataBase.Error() + err.Error.Error())
			}
			sysInfo := SysInfos{}
			err = nft.db.Model(&SysInfos{}).Last(&sysInfo)
			if err.Error != nil {
				if err.Error != gorm.ErrRecordNotFound {
					log.Println("DelCollection() SysInfos err=", err)
					return errors.New(ErrDataBase.Error() + err.Error.Error())
				}
				err = nft.db.Model(&SysInfos{}).Create(&sysInfo)
				if err.Error != nil {
					log.Println("DelCollection() SysInfos create err=", err)
					return errors.New(ErrDataBase.Error() + err.Error.Error())
				}
			}
			fmt.Println("total=", total)
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("nfttotal", sysInfo.Nfttotal-uint64(total))
			if err.Error != nil {
				fmt.Println("DelCollection() add  SysInfos nfttotal err=", err.Error)
				return errors.New(ErrDataBase.Error() + err.Error.Error())
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
	} else {
		fmt.Println("DelCollection() nft mintstate under the collection cannot be deleted")
		return ErrDeleteCollection
	}

}
