package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

func (nft NftDb) NewSnftCollections(useraddr, name, img, contract_type, contract_addr,
	desc, categories, sig, exchanger string) error {
	useraddr = strings.ToLower(useraddr)
	contract_addr = strings.ToLower(contract_addr)
	fmt.Println("NewCollections() user_addr=", useraddr, "      time=", time.Now().String())
	UserSync.Lock(useraddr)
	defer UserSync.UnLock(useraddr)
	//fmt.Println("NewCollections() useraddr=", useraddr )
	fmt.Println("NewCollections() contract_addr=", contract_addr)
	var snftcollectrec SnftCollect
	err := nft.db.Where("name = ? ", name).First(&snftcollectrec)
	if err.Error == nil {
		fmt.Println("NewSnftCollections() err=Collection already exist.")
		return ErrCollectionExist
	} else if err.Error == gorm.ErrRecordNotFound {
		snftcollectrec = SnftCollect{}
		snftcollectrec.Createaddr = useraddr
		snftcollectrec.Name = name
		snftcollectrec.Desc = desc
		snftcollectrec.Exchanger = exchanger
		if contract_addr != "" {
			snftcollectrec.Contract = contract_addr
		} else {
			//collectRec.Contract = strings.ToLower(NFT1155Addr)
			snftcollectrec.Contract = strings.ToLower(ExchangeOwer)
		}
		snftcollectrec.Contracttype = contract_type
		snftcollectrec.Categories = categories
		snftcollectrec.SigData = sig
		newtoken, terr := nft.NewCollectTokenGen()
		if terr != nil {
			fmt.Printf("newtokengen err=%s", terr)
			return terr
		}
		file, serr := saveIpfsjpgImage(newtoken, img)
		if serr != nil {
			fmt.Println("SaveToIpfs() save collection image err=", serr)
			return serr
		}
		collectImageUrl, serr := SaveToIpfs(file)
		if serr != nil {
			fmt.Println("SaveToIpfs() save collection image err=", serr)
			return serr
		}
		snftcollectrec.Img = "/ipfs/" + collectImageUrl
		snftcollectrec.Tokenid = newtoken
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&SnftCollect{}).Create(&snftcollectrec)
			if err.Error != nil {
				fmt.Println("NewCollections() err=", err.Error)
				return err.Error
			}
			imagerr := SaveSnftCollectionsImage(ImageDir, useraddr, name, img)
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

func (nft NftDb) SetSnftCollection(collecton SnftCollection) error {
	var snftcollectrec SnftCollect
	err := nft.db.Model(&snftcollectrec).Where("tokenid = ?", collecton.TokenID).First(&snftcollectrec)
	if err.Error != nil {
		fmt.Println("SetSnftPeriod() err= not find period.")
		return err.Error
	} else {
		if collecton.Name != "" {
			err = nft.db.Where(" name = ? ", collecton.Name).First(&snftcollectrec)
			if err.Error == nil {
				fmt.Println("SetSnftPeriod() err=name already exist.")
				return ErrCollectionExist
			} else {
				snftcollectrec.Name = collecton.Name
			}
		}
		if collecton.Desc != "" {
			snftcollectrec.Desc = collecton.Desc
		}
		if collecton.Categories != "" {
			snftcollectrec.Categories = collecton.Categories
		}

		err = nft.db.Model(&SnftPhase{}).Where("tokenid = ?", collecton.TokenID).Updates(&snftcollectrec)
		if err.Error != nil {
			fmt.Println("SetSnftCollection() update err= ", err.Error)
			return err.Error
		}

		fmt.Println("SetSnftCollection() Ok.")
		return nil
	}
}

func (nft NftDb) SetCollectSnft(collecton, collectid string) error {
	var collectsnft []ModifyPeriodCollect
	err := json.Unmarshal([]byte(collecton), &collectsnft)
	if err != nil {
		fmt.Println("setcollectSnft  Unmrashal input err=", err)
		return err
	}
	if len(collectsnft) > 16 {
		fmt.Println("setCollectSnft data err")
		return errors.New("Snft data err")
	}
	snftcollects := SnftCollect{}
	ferr := nft.db.Model(&snftcollects).Where("tokenid = ?", collectid).Find(&snftcollects)
	if ferr.Error != nil {
		fmt.Println("SetSnftPeriod() err= not find period.")
		return ferr.Error
	}

	return nft.db.Transaction(func(tx *gorm.DB) error {
		collectstr := ""
		err := tx.Model(&Snfts{}).Unscoped().Where("collection =? ", collectid).Delete(&Snfts{})
		if err.Error != nil {
			fmt.Println(" SnftCollect  delete err= ", err.Error)
			return err.Error
		}
		for i, snft := range collectsnft {
			//existsnft := Snfts{}
			existnft := Nfts{}
			err = tx.Model(&Nfts{}).Where("tokenid = ? ", snft.Collect).First(&existnft)
			if err.Error != nil {
				fmt.Printf("input nft err=%s", err.Error)
				return errors.New("input nft err")
			}
			insnft := Snfts{}
			insnft.Name = existnft.Name
			insnft.Desc = existnft.Desc
			insnft.Ownaddr = existnft.Ownaddr
			insnft.Image = existnft.Image
			insnft.Md5 = existnft.Md5
			insnft.Meta = existnft.Meta
			insnft.Nftmeta = existnft.Nftmeta
			insnft.Url = existnft.Url
			insnft.Contract = existnft.Contract
			insnft.Tokenid = existnft.Tokenid
			insnft.Nftaddr = existnft.Nftaddr
			insnft.Count = 1
			insnft.Approve = existnft.Approve
			insnft.Categories = existnft.Categories
			insnft.Hide = existnft.Hide
			insnft.Signdata = existnft.Signdata
			insnft.Createaddr = existnft.Createaddr
			insnft.Verifyaddr = existnft.Verifyaddr
			insnft.Currency = existnft.Currency
			insnft.Price = existnft.Price
			insnft.Royalty = existnft.Royalty
			insnft.Collection = collectid
			insnft.Local = snft.Local
			err = tx.Model(&Snfts{}).Create(&insnft)
			if err.Error != nil {
				fmt.Printf("SetCollectSnt() create  snft err=%v", err.Error)
				return err.Error
			}
			collectstr += snft.Collect
			if i < len(collectsnft)-1 {
				collectstr += ","
			}
		}
		totalcount := len(collectsnft)
		fmt.Println(totalcount)
		snftcollects.Snft = collectstr
		snftcollects.Totalcount = totalcount
		err = tx.Model(&SnftCollect{}).Where("tokenid = ?", collectid).Updates(&snftcollects)
		if err.Error != nil {
			fmt.Println("SetCollectSnft() update  collect  err= ", err.Error)
			return err.Error
		}
		if totalcount != 16 {
			snftcollectperiod := []SnftCollectPeriod{}
			err := tx.Model(&SnftCollectPeriod{}).Where("collect=? ", collectid).Find(&snftcollectperiod)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
				return err.Error
			}
			for _, period := range snftcollectperiod {
				snftphase := SnftPhase{}
				err := tx.Model(&SnftPhase{}).Where("tokenid=? ", period.Period).Find(&snftphase)
				if err.Error != nil {
					fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
					return err.Error
				}
				snftphase.Accedvote = ""
				err = tx.Model(&SnftPhase{}).Where("tokenid = ?", period.Period).Updates(&snftphase)
				if err.Error != nil {
					fmt.Println("SetCollectSnft() update  collect  err= ", err.Error)
					return err.Error
				}
			}
		} else {
			//snftcollectperiod := []SnftCollectPeriod{}
			//err := tx.Model(&SnftCollectPeriod{}).Where("collect=? ", collectid).Find(&snftcollectperiod)
			//if err.Error != nil {
			//	fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			//	return err.Error
			//}
			//for _, period := range snftcollectperiod {
			//	snftphase := SnftPhase{}
			//	err := tx.Model(&SnftPhase{}).Where("tokenid=? ", period.Period).Find(&snftphase)
			//	if err.Error != nil {
			//		fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			//		return err.Error
			//	}
			//	if snftphase.Accedvote != "" {
			//		break
			//	}
			//	snftperup := []SnftCollectPeriod{}
			//	err = tx.Model(&SnftCollectPeriod{}).Where("period=? ", period.Period).Find(&snftperup)
			//	if err.Error != nil {
			//		fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			//		return err.Error
			//	}
			//	total := 0
			//	for _, collect := range snftperup {
			//		snft := SnftCollect{}
			//		err = tx.Model(&SnftCollect{}).Where("tokenid = ?", collect.Collect).Find(&snft)
			//		if err.Error != nil {
			//			fmt.Println("SetCollectSnft() update  collect  err= ", err.Error)
			//			return err.Error
			//		}
			//		fmt.Println(collect)
			//		//for _, totalsnft := range snft {
			//		//	param := strings.Split(totalsnft.Snft, ",")
			//		//	fmt.Println("snft:", totalsnft.Snft, ",collect: ", totalsnft.Tokenid, "collectif: ", collectid, "param: ", (len(param) == 16) || (totalsnft.Tokenid == collectid))
			//		//
			//		//
			//		//}
			//		param := strings.Split(snft.Snft, ",")
			//		if (len(param) == 16) || (snft.Tokenid == collectid) {
			//			total++
			//		} else {
			//			break
			//		}
			//
			//	}
			//	accedvote := ""
			//	if total == 16 {
			//		accedvote = "false"
			//	}
			//	fmt.Println("accedvote =", accedvote, ",total=", total)
			//	err = tx.Model(&SnftPhase{}).Where("tokenid = ?", period.Period).Update("accedvote", accedvote)
			//	if err.Error != nil {
			//		fmt.Println("SetCollectSnft() update  collect  err= ", err.Error)
			//		return err.Error
			//	}
			//
			//}

			snftcollectperiod := []SnftPhase{}
			err := tx.Model(&SnftCollectPeriod{}).Select("snftphase.*").Joins("left join snftphase on snftphase.tokenid =  snftcollectperiod.period").
				Where("snftcollectperiod.collect = ?", collectid).Find(&snftcollectperiod)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
				return err.Error

			}
			for _, period := range snftcollectperiod {
				if period.Accedvote != "" {
					continue
				}
				snftperup := []SnftCollect{}
				err = tx.Model(&SnftCollectPeriod{}).Select("snftcollect.*").Joins("left join snftcollect on snftcollect.tokenid =  snftcollectperiod.collect").
					Where("snftcollectperiod.period = ?", period.Tokenid).Find(&snftperup)
				if err.Error != nil {
					fmt.Println("SetSnftPeriod() find SnftCollect err=.", err.Error)
					return err.Error
				}
				total := 0
				for _, collect := range snftperup {
					fmt.Println(collect)
					param := strings.Split(collect.Snft, ",")
					if (len(param) == 16) || (collect.Tokenid == collectid) {
						total++
					} else {
						break
					}

				}
				accedvote := ""
				if total == 16 {
					accedvote = "false"
				}
				fmt.Println("accedvote =", accedvote, ",total=", total)
				err = tx.Model(&SnftPhase{}).Where("tokenid = ?", period.Tokenid).Update("accedvote", accedvote)
				if err.Error != nil {
					fmt.Println("SetCollectSnft() update  collect  err= ", err.Error)
					return err.Error
				}

			}
		}

		fmt.Println("SetSnftCollect()  Ok")
		return nil
	})
}

func (nft NftDb) NewCollectTokenGen() (string, error) {
	var NewTokenid string
	spendT := time.Now()
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
		nfttab := SnftCollect{}
		err := nft.db.Model(&SnftCollect{}).Where("tokenid = ?", NewTokenid).First(&nfttab)
		if err.Error == gorm.ErrRecordNotFound {
			fmt.Printf("UploadNft() Nfts{} Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
			break
		}
		fmt.Println("UploadNft() Tokenid repetition.", NewTokenid)
	}
	if i >= 20 {
		fmt.Println("UploadNft() generate tokenId error.")
		return "", ErrGenerateTokenId
	}
	return NewTokenid, nil
}

func (nft NftDb) GetSnftCollection(collectid string) ([]Snfts, error) {
	snftcollect := SnftCollect{}
	db := nft.db.Model(&SnftCollect{}).Where("tokenid = ?", collectid).Find(&snftcollect)
	if db.Error != nil {
		fmt.Println("GetSnftCollection() dbase err=", db.Error)
		return nil, db.Error
	}
	snftlist := []Snfts{}
	db = nft.db.Model(&Snfts{}).Where(" collection = ? ", collectid).Find(&snftlist)
	if db.Error != nil {
		fmt.Printf("GetSnftCollection() err=%s", db.Error)
		return nil, db.Error
	}
	nftlist := []Snfts{}
	for _, snft := range snftlist {
		snft.Image = ""
		snft.Signdata = ""

		nftlist = append(nftlist, snft)
	}
	return nftlist, nil
	//if snftcollect.Snft != "" {
	//	coll := strings.Split(snftcollect.Snft, ",")
	//	for _, snft := range coll {
	//		if snft != "" {
	//			//snftint, err := strconv.Atoi(snft)
	//			//if err != nil {
	//			//	fmt.Printf("GetSnftCollection=%v, snft err =%s", snftcollect.ID, err)
	//			//
	//			//}
	//			cocolloct := Snfts{}
	//			serr := nft.db.Model(&Snfts{}).Where(" tokenid = ? ", snft).Find(&cocolloct)
	//			if serr.Error != nil {
	//				fmt.Printf("GetSnftCollection() err=%s", serr.Error)
	//				return nil, serr.Error
	//			}
	//			cocolloct.Image = ""
	//			snftlist = append(snftlist, cocolloct)
	//		}
	//	}
	//	return snftlist, nil
	//} else {
	//	return nil, nil
	//}
}

func (nft NftDb) CollectSearch(categories, param string) ([]SnftCollect, error) {
	snftsearch := []SnftCollect{}
	if param == "" && categories == "" {
		err := nft.db.Model(&SnftCollect{}).Find(&snftsearch)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			fmt.Printf("search nft err=%s", err.Error)
			return nil, err.Error
		}
		for i, _ := range snftsearch {
			snftsearch[i].Img = ""
		}
		return snftsearch, nil
	}

	if categories == "" {
		err := nft.db.Model(&SnftCollect{}).Where("name like ?", "%"+param+"%").Find(&snftsearch)
		if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
			fmt.Printf("search nft err=%s", err.Error)
			return nil, err.Error
		}
		for i, _ := range snftsearch {
			snftsearch[i].Img = ""
		}
		return snftsearch, nil
	} else {
		catestr := strings.Split(categories, ",")
		catesql := "select * from snftcollect where deleted_at IS NULL and ( "
		for i, str := range catestr {
			if i < len(catestr)-1 {
				catesql = catesql + " categories = " + "'" + str + "'" + " or "
			}
			if i == len(catestr)-1 {
				catesql = catesql + " categories =  " + "'" + str + "'"
			}
		}
		catesql += " ) and name like ?"
		err := nft.db.Raw(catesql, "%"+param+"%").Scan(&snftsearch)
		if err.Error != nil {
			fmt.Println("SnftSearch() Dayinfo err=", err)
			return nil, err.Error
		}
		for i, _ := range snftsearch {
			snftsearch[i].Img = ""
		}
		return snftsearch, nil
	}
}

func (nft NftDb) DelSnftCollect(delcollect string) error {
	if delcollect == "" {
		fmt.Println("params error")
		return errors.New("params error")
	}
	return nft.GetDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&SnftCollect{}).Where("tokenid= ?", delcollect).Delete(&SnftPhase{})
		if err.Error != nil {
			fmt.Println("delete snftPeriod err=", err.Error)
			return err.Error
		}
		err = nft.db.Model(&SnftCollectPeriod{}).Where(" collect = ?", delcollect).Delete(&SnftCollectPeriod{})
		if err.Error != nil {
			fmt.Println("DelSnftCollect() delete  snftcollect err= ", err.Error)
			return err.Error
		}
		return nil
	})

}
