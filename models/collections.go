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
				return err.Error
			}
			imagerr := SaveCollectionsImage(ImageDir, useraddr, name, img)
			if imagerr != nil {
				fmt.Println("NewCollections() SaveCollectionsImage() err=", imagerr)
				return ErrCollectionImage
			}
			return nil
		})
	}
	fmt.Println("NewCollections() dbase err=.", err)
	return err.Error
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
			return ErrCollectionImage
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
			return err.Error
		}
		return nil
	})
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
	err := nft.db.Model(Collects{}).Where("totalcount > 0").Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryNFTCollectionList() recCount err=", err)
		return nil, 0, ErrNftNotExist
	}
	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)
	if int64(startIndex) >= recCount || recCount == 0 {
		return nil, 0, ErrNftNotExist
	} else {
		temp := recCount - int64(startIndex)
		if int64(nftCount) > temp {
			nftCount = int(temp)
		}
		err = nft.db.Model(Collects{}).Select([]string{"createaddr", "contract", "Contracttype", "name", "img"}).Where("totalcount > 0").Order("transamt desc, id desc").Limit(nftCount).Offset(startIndex).Find(&collectRecs)
		if err.Error != nil {
			fmt.Println("QueryNFTCollectionList() find record err=", err)
			return nil, 0, ErrNftNotExist
		}
		userCollects := make([]UserCollection, 0, 20)
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
			userCollects = append(userCollects, userCollect)
		}
		return userCollects, int(recCount), nil
	}
}

func (nft NftDb) DelCollection(useraddr, contract, name string) error {

	collect := Collects{}
	err := nft.db.Model(&Collects{}).Where("name =? and contract=? and createaddr=?", name, contract, useraddr).First(&collect)
	if err.Error != nil {
		fmt.Println("DelCollection() RecordNotFound ,err=", err)
		return errors.New("collection RecordNotFound")
	}
	nfts := Nfts{}
	err = nft.db.Model(&Nfts{}).Where("collectcreator =? and collections=? and mintstate <> ? ", useraddr, name, "NoMinted").First(&nfts)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("DelCollection() delete subscribe record err=", err.Error)
			return err.Error
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Collects{}).Where("name =? and contract=? and createaddr=?", name, contract, useraddr).Delete(&Collects{})
			if err.Error != nil {
				fmt.Println("DelCollection() delete subscribe record err=", err.Error)
				return err.Error
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
				return err.Error
			}
			sysInfo := SysInfos{}
			err = nft.db.Model(&SysInfos{}).Last(&sysInfo)
			if err.Error != nil {
				if err.Error != gorm.ErrRecordNotFound {
					log.Println("DelCollection() SysInfos err=", err)
					return ErrCollectionNotExist
				}
				err = nft.db.Model(&SysInfos{}).Create(&sysInfo)
				if err.Error != nil {
					log.Println("DelCollection() SysInfos create err=", err)
					return ErrCollectionNotExist
				}
			}
			fmt.Println("total=", total)
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("nfttotal", sysInfo.Nfttotal-uint64(total))
			if err.Error != nil {
				fmt.Println("DelCollection() add  SysInfos nfttotal err=", err.Error)
				return err.Error
			}

			return nil
		})
	} else {
		fmt.Println("DelCollection() nft mintstate under the collection cannot be deleted")
		return errors.New("nfts mintstate under the collection  cannot be deleted")
	}

}

func (nft NftDb) SetCollection(useraddr, contract, name string) error {

	collect := Collects{}
	err := nft.db.Model(&Collects{}).Where("name =? and contract=? and createaddr=?", name, contract, useraddr).First(&collect)
	if err.Error != nil {
		fmt.Println("SetCollection() collection RecordNotFound")
		return errors.New("collection RecordNotFound")
	}
	nfts := Nfts{}
	err = nft.db.Model(&Nfts{}).Where("collectcreator =? and collections=? and mintstate <> ? ", useraddr, name, "NoMinted").Find(&nfts)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("DelCollection() delete subscribe record err=", err.Error)
			return err.Error
		}
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Collects{}).Where("name =? and contract=? and createaddr=?", name, contract, useraddr).Delete(&Collects{})
			if err.Error != nil {
				fmt.Println("DelCollection() delete subscribe record err=", err.Error)
				return err.Error
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
				return err.Error
			}
			sysInfo := SysInfos{}
			err = nft.db.Model(&SysInfos{}).Last(&sysInfo)
			if err.Error != nil {
				if err.Error != gorm.ErrRecordNotFound {
					log.Println("DelCollection() SysInfos err=", err)
					return ErrCollectionNotExist
				}
				err = nft.db.Model(&SysInfos{}).Create(&sysInfo)
				if err.Error != nil {
					log.Println("DelCollection() SysInfos create err=", err)
					return ErrCollectionNotExist
				}
			}
			fmt.Println("total=", total)
			err = tx.Model(&SysInfos{}).Where("id = ?", sysInfo.ID).Update("nfttotal", sysInfo.Nfttotal-uint64(total))
			if err.Error != nil {
				fmt.Println("DelCollection() add  SysInfos nfttotal err=", err.Error)
				return err.Error
			}
			NftCatch.SetFlushFlag()
			return nil
		})
	} else {
		fmt.Println("SetCollection() nft mintstate under the collection cannot be modify")
		return errors.New("nfts mintstate under the collection  cannot be modify")
	}

}