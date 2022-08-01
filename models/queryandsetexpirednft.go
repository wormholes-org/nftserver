package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type NftsExpired struct {
	Name       string `json:"name"`
	Tokenid    string `json:"tokenid"`
	UpdataTime string `json:"updata_time"`
}

func (nft NftDb) QueryExpireNft(start_index, count, param string) ([]Nfts, int, error) {
	var collectRecs []Nfts
	var recCount int64
	//t := time.Now()
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	if param == "" {
		err := nft.db.Model(Nfts{}).Where("mintstate =?", NoMinted.String()).Count(&recCount)
		if err.Error != nil {
			fmt.Println("QueryExpireNft() recCount err=", err)
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
			err = nft.db.Model(Nfts{}).Where("mintstate =?", NoMinted.String()).Limit(nftCount).Offset(startIndex).Find(&collectRecs)
			if err.Error != nil {
				fmt.Println("QueryExpireNft() find record err=", err)
				return nil, 0, ErrNftNotExist
			}
			expiredlist := []NftsExpired{}
			for _, value := range collectRecs {
				expired := NftsExpired{}
				expired.Name = value.Name
				expired.Tokenid = value.Tokenid
				expired.UpdataTime = value.UpdatedAt.Format("2006-01-02")
				expiredlist = append(expiredlist, expired)
			}
			return collectRecs, int(recCount), nil
		}
	} else {
		err := nft.db.Model(Nfts{}).Where("updated_at < date_sub(now(), interval ? day) and mintstate = ?", param, NoMinted.String()).Count(&recCount)
		if err.Error != nil {
			fmt.Println("QueryExpireNft() recCount err=", err)
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
			err = nft.db.Model(Nfts{}).Where("updated_at < date_sub(now(), interval ? day) and mintstate = ?", param, NoMinted.String()).Limit(nftCount).Offset(startIndex).Find(&collectRecs)
			if err.Error != nil {
				fmt.Println("QueryExpireNft() find record err=", err)
				return nil, 0, ErrNftNotExist
			}
			expiredlist := []NftsExpired{}
			for _, value := range collectRecs {
				expired := NftsExpired{}
				expired.Name = value.Name
				expired.Tokenid = value.Tokenid
				expired.UpdataTime = value.UpdatedAt.Format("2006-01-02")
				expiredlist = append(expiredlist, expired)
			}
			return collectRecs, int(recCount), nil
		}
	}

}

func (nft NftDb) DelExpiredNft(param string) error {

	if param == "" {
		log.Println("input param err")
		return ErrDataFormat
	}
	var total int64
	//err := nft.db.Model(Nfts{}).Where("updated_at < date_sub(now(), interval ? day)", param).Count(&total).Delete(&Nfts{})
	//if err.Error != nil {
	//	fmt.Println("DelExpiredNft() delete nft err=", err)
	//	return ErrNftNotExist
	//}
	return nft.db.Transaction(func(tx *gorm.DB) error {

		nftlist := []Nfts{}
		//err := nft.db.Model(Nfts{}).Where("updated_at < date_sub(now(), interval ? day)", param).Find(&nftlist)
		//if err.Error != nil {
		//	fmt.Println("DelExpiredNft() delete nft err=", err)
		//	return ErrNftNotExist
		//}
		err := nft.db.Debug().Model(Nfts{}).Where("updated_at < date_sub(now(), interval ? day) and mintstate = ?", param, NoMinted.String()).Count(&total).Find(&nftlist).Delete(&Nfts{})
		if err.Error != nil {
			fmt.Println("DelExpiredNft() delete nft err=", err)
			return ErrNftNotExist
		}
		fmt.Println("total =", total)
		fmt.Println("nftlist =", nftlist)
		for _, valule := range nftlist {
			err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
				valule.Collections, valule.Collectcreator).Update("totalcount", gorm.Expr("totalcount - ?", 1))
			if err.Error != nil {
				fmt.Println("DelExpiredNft() delete collectins totalcount err= ", err.Error)
				return ErrDataBase
			}
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
}
