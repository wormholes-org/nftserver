package models

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type SearchData struct {
	NftsRecords     []Nfts     `json:"nfts"`
	CollectsRecords []Collects `json:"collections"`
	UserAddrs       []string   `json:"user_addrs"`
}

type SearchCatch struct {
	Searchs []SearchData
}

func (nft *NftDb) Search(cond string) ([]SearchData, error) {
	var searchData SearchData
	spendT := time.Now()
	queryCatchSql := cond
	searchCatch := SearchCatch{}
	cerr := GetRedisCatch().GetCatchData("Search", queryCatchSql, &searchCatch)
	if cerr == nil {
		log.Printf("Search() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return searchCatch.Searchs, nil
	}
	nfts := []Nfts{}
	findNftsResult := nft.db.Model(&Nfts{}).Where("name like ?", "%"+cond+"%").
		Order("name asc").Offset(0).Limit(5).Find(&nfts)
	if findNftsResult.Error != nil && findNftsResult.Error != gorm.ErrRecordNotFound {
		return nil, findNftsResult.Error
	}
	for k, _ := range nfts {
		nfts[k].Image = ""
	}
	searchData.NftsRecords = append(searchData.NftsRecords, nfts...)

	collects := []Collects{}
	findCollectsResult := nft.db.Model(&Collects{}).Where("createaddr like ? or name like ? and  name <> ?", "%"+cond+"%", "%"+cond+"%", DefaultCollection).
		Order("name asc").Offset(0).Limit(5).Find(&collects)
	if findCollectsResult.Error != nil && findCollectsResult.Error != gorm.ErrRecordNotFound {
		return nil, findCollectsResult.Error
	}
	for k, _ := range collects {
		if collects[k].Contracttype != "snft" {
			collects[k].Img = ""
		}
	}
	searchData.CollectsRecords = append(searchData.CollectsRecords, collects...)

	users := []Users{}
	findUsersResult := nft.db.Model(&Users{}).Where("useraddr like ? or username like ?", "%"+cond+"%", "%"+cond+"%").
		Order("username asc").Offset(0).Limit(5).Find(&users)
	if findUsersResult.Error != nil && findUsersResult.Error != gorm.ErrRecordNotFound {
		return nil, findUsersResult.Error
	}

	searchData.UserAddrs = make([]string, 0)
	for _, user := range users {
		user.Portrait = ""
		user.Kycpic = ""
		searchData.UserAddrs = append(searchData.UserAddrs, user.Useraddr)
	}
	GetRedisCatch().CatchQueryData("Search", queryCatchSql, &SearchCatch{[]SearchData{searchData}})
	log.Printf("Search() no catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return []SearchData{searchData}, nil
}
