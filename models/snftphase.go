package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type QueryPeriod struct {
	Createaddr     string        `json:"collection_creator_addr"`
	Name           string        `json:"name"`
	Desc           string        `json:"desc"`
	Vote           int           `json:"vote"`
	Accedvote      string        `json:"accedvote"`
	Categories     string        `json:"categories"`
	Collect        string        `json:"collect"`
	Extend         string        `json:"extend"`
	Tokenid        string        `json:"tokenid"`
	CollectionList []SnftCollect `json:"collectionlist"`
}

type ModifyPeriodCollect struct {
	Collect string `json:"collect"`
	Local   string `json:"local"`
}

type ToSnftPeriodCollect struct {
	Createaddr   string `json:"collection_creator_addr"`
	Contract     string `json:"nft_contract_addr"`
	Contracttype string `json:"contracttype"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Categories   string `json:"categories"`
	Totalcount   int    `json:"total_count"`
	Transcnt     int    `json:"transcnt"`
	Transamt     uint64 `json:"transamt"`
	SigData      string `json:"sig"`
	Img          string `json:"img"`
	Tokenid      string `json:"tokenid"`
	Snft         string `json:"snft"`
	//Period       string `json:"period" gorm:"type:char(42) ;comment:'期号'"`
	Local     string `json:"local"`
	Exchanger string `json:"exchanger"`
	Extend    string `json:"extend"`
}

func (nft NftDb) NewSnftPhase(useraddr, name, desc string) error {
	useraddr = strings.ToLower(useraddr)
	fmt.Println("NewCollections() user_addr=", useraddr, "      time=", time.Now().String())
	UserSync.Lock(useraddr)
	defer UserSync.UnLock(useraddr)
	var snftcollectrec SnftPhase
	err := nft.db.Where("name = ? ", name).First(&snftcollectrec)
	if err.Error == nil {
		fmt.Println("NewSnftPhase() err=name already exist.")
		return ErrCollectionExist
	} else if err.Error == gorm.ErrRecordNotFound {
		snftcollectrec = SnftPhase{}
		snftcollectrec.Createaddr = useraddr
		snftcollectrec.Name = name
		snftcollectrec.Desc = desc
		newtoken, terr := nft.NewPeriodTokenGen()
		if terr != nil {
			fmt.Printf("newtokengen err=%s", terr)
			return terr
		}
		snftcollectrec.Tokenid = newtoken
		err = nft.db.Model(&SnftPhase{}).Create(&snftcollectrec)
		if err.Error != nil {
			fmt.Println("NewSnftPhase() create SnftPeriod err= ", err.Error)
			return err.Error
		}
		return nil
	}
	fmt.Println("NewSnftPhase() dbase err=.", err)
	return err.Error
}

func (nft NftDb) SetSnftPeriod(period Period) error {
	var snftcollectrec SnftPhase
	//sql := "select * from snftcollect where id = ?"
	//err := nft.db.Raw(sql, period.ID).Scan(&snftcollectrec)
	err := nft.db.Model(&snftcollectrec).Where("tokenid = ?", period.TokenID).First(&snftcollectrec)
	if err.Error != nil {
		fmt.Println("SetSnftPeriod() err= not find period.")
		return err.Error
	} else {
		if period.Name != "" {
			err = nft.db.Where(" name = ? ", period.Name).First(&snftcollectrec)
			if err.Error == nil {
				fmt.Println("SetSnftPeriod() err=name already exist.")
				return ErrCollectionExist
			} else {
				snftcollectrec.Name = period.Name
			}
		}
		if period.Desc != "" {
			snftcollectrec.Desc = period.Desc
		}
		if period.Categories != "" {
			snftcollectrec.Categories = period.Categories
		}

		if period.Accedvote != "" {
			snftcollectrec.Accedvote = period.Accedvote
		}
		if period.Collect != "" {
			var percoll []ModifyPeriodCollect
			uerr := json.Unmarshal([]byte(period.Collect), &percoll)
			fmt.Println(percoll)
			if uerr != nil {
				log.Println("ModifyPeriodCollect()  Unmarshal() err=", err)
			}
			if len(percoll) > 16 {
				fmt.Println("select collect to Period too long ")
				return errors.New("select collect to Period too long")
			}
			return nft.db.Transaction(func(tx *gorm.DB) error {
				snftphase := SnftPhase{}
				err := tx.Model(&SnftPhase{}).Where("tokenid = ?", period.TokenID).First(&snftphase)
				if err.Error != nil {
					fmt.Println(" SnftCollect  get err= ", err.Error)
					return err.Error
				}
				err = tx.Model(&SnftCollectPeriod{}).Where("period =? ", period.TokenID).Delete(&SnftCollectPeriod{})
				if err.Error != nil {
					fmt.Println(" SnftCollect  delete err= ", err.Error)
					return err.Error
				}

				total := 0
				collect := make(map[string]string)
				for _, modify := range percoll {
					_, ok := collect[modify.Collect]
					if !ok {
						collect[modify.Collect] = modify.Collect
					} else {
						fmt.Println(" SnftCollect  set input data repeat ")
						return errors.New("SnftCollect  set input data repeat")
					}
					//collectperiod := SnftCollectPeriod{}
					snftperiod := &SnftCollectPeriod{}
					snftperiod.Period = period.TokenID
					snftperiod.Collect = modify.Collect
					snftperiod.Local = modify.Local
					err = tx.Model(&SnftCollectPeriod{}).Create(&snftperiod)
					if err.Error != nil {
						fmt.Printf("SetSnftPeriod() create SnftCollectPeriod err=%s", err.Error)
						return err.Error
					}

					collect := SnftCollect{}
					err = tx.Model(&SnftCollect{}).Where("tokenid=? ", modify.Collect).Find(&collect)
					if err.Error != nil {
						fmt.Println("SetVoteSnftPeriod() find SnftCollectPeriod err=", err.Error)
						return err.Error
					}
					snfts := strings.Split(collect.Snft, ",")
					if len(snfts) == 16 {
						total++
					}
				}
				fmt.Println(total, " and ", snftcollectrec.Accedvote)
				if total != 16 {
					snftcollectrec.Accedvote = ""
				} else {
					snftcollectrec.Accedvote = "false"
				}
				fmt.Println(snftcollectrec.Accedvote)
				err = tx.Model(&SnftPhase{}).Where("tokenid = ?", period.TokenID).Updates(&snftcollectrec)
				if err.Error != nil {
					fmt.Println("period collect  update err= ", err.Error)
					return err.Error
				}
				fmt.Println("SetSnftPeriod()  Ok")
				return nil
			})
		} else {
			fmt.Println("collect = null ", snftcollectrec.Accedvote)
			err = nft.db.Model(&SnftPhase{}).Where("tokenid = ?", period.TokenID).Updates(&snftcollectrec)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() update err= ", err.Error)
				return err.Error
			}
			fmt.Println("SetSnftPeriod() Ok.")
			return nil
		}
	}
}

func (nft NftDb) NewPeriodTokenGen() (string, error) {
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
		nfttab := SnftPhase{}
		err := nft.db.Model(&SnftPhase{}).Where("tokenid = ?", NewTokenid).First(&nfttab)
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

func (nft NftDb) GetSnftPeriod() ([]QueryPeriod, error) {
	//collectionInfo := SnftPhase{}
	snftphase := []SnftPhase{}
	snftphasecollect := []QueryPeriod{}
	//var snftcollectrec SnftPhase
	db := nft.db.Model(&SnftPhase{}).Where("accedeth is null or accedeth =? ", "").Find(&snftphase)
	if db.Error != nil {
		fmt.Println("GetSnftPeriod() dbase err=", db.Error)
		return nil, db.Error
	}
	for _, perod := range snftphase {
		collect := QueryPeriod{}
		collect.Name = perod.Name
		collect.Desc = perod.Desc
		collect.Vote = perod.Vote
		collect.Accedvote = perod.Accedvote
		collect.Categories = perod.Categories
		collect.Extend = perod.Extend
		collect.Tokenid = perod.Tokenid
		percoll := []SnftCollect{}
		collectlocal := []SnftCollectPeriod{}
		err := nft.db.Model(&SnftCollectPeriod{}).Where("period=? ", perod.Tokenid).Find(&collectlocal)
		if err.Error != nil {
			fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			return nil, err.Error
		}
		for _, local := range collectlocal {
			cocolloct := SnftCollect{}
			err := nft.db.Where(" tokenid = ? ", local.Collect).First(&cocolloct)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() eSnftCollect err=.", err.Error)
				return nil, err.Error
			}
			cocolloct.Img = ""
			cocolloct.Local = local.Local
			percoll = append(percoll, cocolloct)
		}
		collect.CollectionList = percoll
		snftphasecollect = append(snftphasecollect, collect)
	}
	return snftphasecollect, nil
}

func (nft NftDb) GetAllVotePeriod() ([]QueryPeriod, error) {
	//collectionInfo := SnftPhase{}
	snftphase := []SnftPhase{}
	snftphasecollect := []QueryPeriod{}
	//var snftcollectrec SnftPhase
	db := nft.db.Model(&SnftPhase{}).Where("accedvote = ? and (accedeth is null or accedeth =?)", "true", "").Find(&snftphase)
	if db.Error != nil {
		fmt.Println("GetSnftPeriod() dbase err=", db.Error)
		return nil, db.Error
	}
	for _, perod := range snftphase {
		collect := QueryPeriod{}
		collect.Name = perod.Name
		collect.Desc = perod.Desc
		collect.Vote = perod.Vote
		collect.Accedvote = perod.Accedvote
		collect.Categories = perod.Categories
		collect.Extend = perod.Extend
		collect.Tokenid = perod.Tokenid
		percoll := []SnftCollect{}
		collectlocal := []SnftCollectPeriod{}
		err := nft.db.Model(&SnftCollectPeriod{}).Where("period=? ", perod.Tokenid).Find(&collectlocal)
		if err.Error != nil {
			fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			return nil, err.Error
		}
		for _, local := range collectlocal {
			cocolloct := SnftCollect{}
			err := nft.db.Where(" tokenid = ? ", local.Collect).First(&cocolloct)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() eSnftCollect err=.", err.Error)
				return nil, err.Error
			}
			cocolloct.Img = ""
			cocolloct.Local = local.Local
			percoll = append(percoll, cocolloct)
		}
		collect.CollectionList = percoll
		snftphasecollect = append(snftphasecollect, collect)
	}
	return snftphasecollect, nil
}

func (nft NftDb) GetAccedVotePeriod() ([]QueryPeriod, error) {
	//collectionInfo := SnftPhase{}
	snftphase := []SnftPhase{}
	snftphasecollect := []QueryPeriod{}
	//var snftcollectrec SnftPhase
	db := nft.db.Model(&SnftPhase{}).Where("accedvote = ? and (accedeth is null or accedeth =?)", "true", "").Find(&snftphase)
	if db.Error != nil {
		fmt.Println("GetSnftPeriod() dbase err=", db.Error)
		return nil, db.Error
	}
	for _, perod := range snftphase {
		collect := QueryPeriod{}
		collect.Name = perod.Name
		collect.Desc = perod.Desc
		collect.Vote = perod.Vote
		collect.Accedvote = perod.Accedvote
		collect.Categories = perod.Categories
		collect.Extend = perod.Extend
		collect.Tokenid = perod.Tokenid
		percoll := []SnftCollect{}
		collectlocal := []SnftCollectPeriod{}
		err := nft.db.Model(&SnftCollectPeriod{}).Where("period=? ", perod.Tokenid).Find(&collectlocal)
		if err.Error != nil {
			fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			return nil, err.Error
		}
		for _, local := range collectlocal {
			cocolloct := SnftCollect{}
			err := nft.db.Where(" tokenid = ? ", local.Collect).First(&cocolloct)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() eSnftCollect err=.", err.Error)
				return nil, err.Error
			}
			cocolloct.Img = ""
			cocolloct.Local = local.Local
			percoll = append(percoll, cocolloct)
		}
		collect.CollectionList = percoll
		snftphasecollect = append(snftphasecollect, collect)
	}
	rand.Seed(time.Now().UnixNano())
	snftvoteperiod := []QueryPeriod{}
	snftlen := 0
	if len(snftphasecollect) < 4 {
		snftlen = len(snftphasecollect)
		for i := 0; i < snftlen; i++ {
			snftvoteperiod = append(snftvoteperiod, snftphasecollect[i])
		}
	} else {
		collect := make(map[int]int)
		for i := 0; i < 20; i++ {
			if len(snftvoteperiod) > 3 {
				break
			}
			number := rand.Intn(len(snftphasecollect))
			_, ok := collect[number]
			if ok {
				continue
			}
			collect[number] = number
			snftvoteperiod = append(snftvoteperiod, snftphasecollect[number])
		}
	}

	return snftvoteperiod, nil
}

func (nft NftDb) GetVoteSnftPeriod(startIndex string, count string) ([]QueryPeriod, int, error) {
	//collectionInfo := SnftPhase{}
	snftphase := []SnftPhase{}
	snftphasecollect := []QueryPeriod{}
	//var snftcollectrec SnftPhase
	start, _ := strconv.Atoi(startIndex)
	nftCount, _ := strconv.Atoi(count)
	db := nft.db.Model(&SnftPhase{}).Where("accedvote <> '' and (accedeth is null or accedeth =?)", "").Limit(nftCount).Offset(start).Find(&snftphase)
	if db.Error != nil {
		fmt.Println("GetSnftPeriod() dbase err=", db.Error)
		return nil, 0, db.Error
	}
	for _, perod := range snftphase {
		collect := QueryPeriod{}
		collect.Name = perod.Name
		collect.Desc = perod.Desc
		collect.Vote = perod.Vote
		collect.Accedvote = perod.Accedvote
		collect.Categories = perod.Categories
		collect.Extend = perod.Extend
		collect.Tokenid = perod.Tokenid
		percoll := []SnftCollect{}
		collectlocal := []SnftCollectPeriod{}
		err := nft.db.Model(&SnftCollectPeriod{}).Where("period=? ", perod.Tokenid).Find(&collectlocal)
		if err.Error != nil {
			fmt.Println("SetSnftPeriod() find SnftCollectPeriod err=.", err.Error)
			return nil, 0, err.Error
		}
		if len(collectlocal) != 16 {
			continue
		}
		for _, local := range collectlocal {
			cocolloct := SnftCollect{}
			err := nft.db.Where(" tokenid = ? ", local.Collect).First(&cocolloct)
			if err.Error != nil {
				fmt.Println("SetSnftPeriod() eSnftCollect err=.", err.Error)
				return nil, 0, err.Error
			}
			cocolloct.Img = ""
			cocolloct.Local = local.Local
			cocolloct.Snft = ""
			cocolloct.SigData = ""

			//param := strings.Split(cocolloct.Snft, ",")
			//if len(param) != 16 {
			//	continue
			//}
			percoll = append(percoll, cocolloct)
		}
		collect.CollectionList = percoll
		snftphasecollect = append(snftphasecollect, collect)
	}
	return snftphasecollect, len(snftphasecollect), nil
}

func (nft NftDb) SetVoteSnftPeriod(params string) error {
	if params == "" {
		fmt.Println("input param is null")
		return errors.New("input param is null")
	}

	param := strings.Split(params, ",")
	uerr := nft.db.Model(&SnftPhase{}).Where("tokenid not in ?  and accedvote <> '' ", param).Update("accedvote", "false")
	if uerr.Error != nil {
		fmt.Println("SetVoteSnftPeriod() update err=", uerr.Error)
		return uerr.Error
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		for _, period := range param {
			percoll := SnftPhase{}
			err := tx.Model(&SnftPhase{}).Where("tokenid=? ", period).Find(&percoll)
			if err.Error != nil {
				fmt.Println("SetVoteSnftPeriod() find SnftPeriod err=.", err.Error)
				return err.Error
			}
			percoll.Accedvote = "true"
			err = tx.Model(&SnftPhase{}).Where("tokenid=?", period).Updates(&percoll)
			if err.Error != nil {
				fmt.Println("SetVoteSnftPeriod() update vote err=", err.Error)
				return err.Error
			}
			collectperiod := []SnftCollectPeriod{}
			err = tx.Model(&SnftCollectPeriod{}).Where("period=? ", period).Find(&collectperiod)
			if err.Error != nil {
				fmt.Println("SetVoteSnftPeriod() find SnftCollectPeriod err=", err.Error)
				return err.Error
			}
			if len(collectperiod) != 16 {
				fmt.Println("collect data less than 16")
				return errors.New("collect data less than 16")
			}
			//total := 0
			for _, coll := range collectperiod {
				collect := SnftCollect{}
				err = tx.Model(&SnftCollect{}).Where("tokenid=? ", coll.Collect).Find(&collect)
				if err.Error != nil {
					fmt.Println("SetVoteSnftPeriod() find SnftCollectPeriod err=", err.Error)
					return err.Error
				}
				//collectImageUrl, serr := SaveToIpfs(collect.Img)
				//if serr != nil {
				//	log.Println("SaveToIpfs() save collection image err=", serr)
				//	return errors.New("SaveToIpfs() save collection image error.")
				//}
				snfts := []Snfts{}
				err = tx.Model(&Snfts{}).Where("collection=? ", coll.Collect).Find(&snfts)
				if err.Error != nil {
					fmt.Println("SetVoteSnftPeriod() find SnftCollectPeriod err=", err.Error)
					return err.Error
				}
				if len(snfts) != 16 {
					fmt.Println("snft data less than 16")
					return errors.New("snft data less than 16")
				}
				//for _, snft := range snfts {
				//	wg.Add(1)
				//	go nft.savemeta(snft, collect, total, collectImageUrl)
				//	total++
				//}

			}
		}
		//wg.Wait()
		return nil
	})

}

func (nft NftDb) SetPeriodEth(params string) error {
	if params == "" {
		fmt.Println("input param is null")
		return errors.New("input param is null")
	}
	percoll := SnftPhase{}
	uerr := nft.db.Model(&SnftPhase{}).Where("tokenid =  ?", params).Find(&percoll)
	if uerr.Error != nil {
		fmt.Println("SetPeriodEth() update err=", uerr.Error)
		return uerr.Error
	}
	//go contracts.SendNFTTrans()

	return nft.db.Transaction(func(tx *gorm.DB) error {

		collect := []*SnftCollect{}
		err := tx.Model(&SnftCollectPeriod{}).Select("snftcollect.*").Joins("left join snftcollect on snftcollect.tokenid =  snftcollectperiod.collect").
			Where("snftcollectperiod.period = ?", percoll.Tokenid).Find(&collect)
		if err.Error != nil {
			fmt.Println("SetPeriodEth() find SnftCollectPeriod err=", err.Error)
			return err.Error
		}
		var total *int
		num := 0
		total = &num
		if len(collect) != 16 {
			fmt.Println("collect data less than 16")
			return errors.New("collect data less than 16")
		}

		os.Mkdir("./snft", 0777)
		stop := make(chan string)
		collectch := make(chan SnftCollect, 16)
		go savecollect(stop, collect, collectch)
		//syncBlocks, merr := GetDbBlockNumber(sqldsn)
		//if merr != nil {
		//	fmt.Println("InitSyncBlockTs() get scan block num err=", merr)
		//	return merr
		//}
		//snftInfo, merr := contracts.GetNominatedNFTInfo(big.NewInt(int64(syncBlocks)))
		////snftInfo, merr := contracts.GetNominatedNFTInfo(big.NewInt(0).SetUint64(syncBlocks))
		//if merr != nil {
		//	fmt.Println(merr)
		//}
		//periodnum := snftInfo.StartIndex / snftInfo.Number
		//fmt.Println(snftInfo)
	FOR:
		for {
			select {
			case v, ok := <-collectch:
				if ok {
					snfts := []*Snfts{}
					err = tx.Model(&Snfts{}).Where("collection=? ", v.Tokenid).Find(&snfts)
					if err.Error != nil {
						fmt.Println("SetPeriodEth() find SnftCollectPeriod err=", err.Error)
						return err.Error
					}
					if len(snfts) != 16 {
						fmt.Println("snft data less than 16")
						return errors.New("snft data less than 16")
					}
					fmt.Println("savemeta")
					for _, snft := range snfts {
						wg.Add(1)
						go nft.savemeta(snft, &v, *total)
						num++
					}
				} else {
					break FOR
				}
			case v, ok := <-stop:
				if ok {
					return errors.New(v)
				}
			}
		}

		wg.Wait()

		meta, derr := SaveDirToIpfs("./snft")
		if derr != nil {
			fmt.Println("SetPeriodEth() save nftmeta info err=", derr)
			return derr
		}
		dmeta := "/ipfs/" + meta
		fmt.Println("meta=", dmeta)
		go func() {
			wg.Add(1)
			defer wg.Done()
			percoll.Accedeth = "false"
			percoll.Meta = dmeta
			err = tx.Model(&SnftPhase{}).Where("tokenid =? ", percoll.Tokenid).Updates(percoll)
		}()
		if err.Error != nil {
			fmt.Println("SetPeriodEth() update eth err=", err.Error)
			return err.Error
		}
		serr := contracts.SendSnftTrans(dmeta, ExchangerAuth)
		if serr != nil {
			fmt.Println("SetPeriodEth() SendSnftTrans() err=", serr)
			return serr
		}
		wg.Wait()
		return nil
	})
}

func savecollect(stop chan string, collect []*SnftCollect, collectch chan SnftCollect) {
	defer close(stop)
	defer close(collectch)
	//var tg sync.WaitGroup
	for _, coll := range collect {
		collectch <- *coll

	}
	//tg.Wait()

}

func saveIpfsjpgImage(name, image_base64 string) (string, error) {
	os.Mkdir("./snftcollect", 0777)
	newPath := "./snftcollect/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveCollectionsImage() create dir err=", err)
			return "", err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveCollectionsImage() ParseBase64Type() err=", err)
			return "", err
		}
		//hexname := hex.EncodeToString([]byte(name))
		var hexname string
		for _, c := range name {
			hexname += fmt.Sprintf("%02x", c)
		}
		file = newPath + hexname + "." + imagetype
	} else {
		fmt.Println("SaveCollectionsImage() image_base64==0 error.")
		return "", err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveCollectionsImage() imagetype error.")
		return "", err
	}
	err = base64tofile(file, img)
	if err != nil {
		fmt.Println("SaveCollectionsImage() base64toJpeg() err=", err)
		return "", err
	}
	f, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("base64toJpeg() readall err=", err)
		return "", err
	}
	return string(append([]byte{}, f...)), err
	//return string(f), err
}

func base64tofile(file, data string) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		fmt.Println("base64toJpeg() Decode() err=", err)
		return err
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("base64toJpeg() OpenFile() err=", err)
		return err
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		fmt.Println("base64toJpeg() jpeg.Encode() err=", err)
		return err
	}
	i := strings.LastIndex(file, "jpeg")
	if i != -1 {
		file = file[:i] + "jpg"
	} else {
		file = file[:i] + "jpeg"
	}
	f, err = os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("base64toJpeg() OpenFile() err=", err)
		return err
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		fmt.Println("base64toJpeg() jpeg.Encode() err=", err)
		return err
	}

	return nil
}

func SaveDirToIpfs(str string) (string, error) {
	url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	spendT := time.Now()
	s := shell.NewShell(url)
	s.SetTimeout(5 * time.Second)
	mhash, err := s.AddDir(str)
	if err != nil {
		fmt.Println("SaveToIpfs() err=", err)
		return "", err
	}
	fmt.Printf("SaveToIpfs  Spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return mhash, nil
}

func (nft NftDb) savemeta(snft *Snfts, collect *SnftCollect, total int) error {
	defer wg.Done()
	var nftMeta SnftInfo
	nftMeta.CreatorAddr = snft.Createaddr
	nftMeta.Contract = collect.Contract
	nftMeta.Name = snft.Name
	nftMeta.Desc = snft.Desc
	nftMeta.Category = snft.Categories
	nftMeta.Royalty = float64(snft.Royalty / 100)
	//nftMeta.SourceUrl = "image/" + snft.Tokenid + ".jpg"
	//nftMeta.FileType = path.Ext(fmt.Sprintf("%02x", total))[1:]
	//nftMeta.FileType = ""
	nftMeta.SourceUrl = snft.Url
	nftMeta.Md5 = snft.Md5
	nftMeta.CollectionsName = collect.Name
	nftMeta.CollectionsCreator = collect.Createaddr
	nftMeta.CollectionsExchanger = collect.Contract
	nftMeta.CollectionsCategory = collect.Categories
	nftMeta.CollectionsImgUrl = collect.Img
	file := "./snft/" + fmt.Sprintf("%02x", total)
	file6, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("creat file error")
		return errors.New("creat file error")
	}
	//defer file6.Close()

	metaStr, merr := json.Marshal(&nftMeta)
	if merr != nil {
		fmt.Println("SetVoteSnftPeriod() save nftmeta info err=", merr)
		return errors.New("SetVoteSnftPeriod() save marshal info error.")
	}
	file6.Write(metaStr)
	//snft.Meta = meta
	//err := nft.db.Model(&Snfts{}).Where("tokenid =? and collection=? and local=? ", snft.Tokenid, snft.Collection, snft.Local).Update("meta", snft.Meta)
	//if err.Error != nil {
	//	fmt.Println("SetVoteSnftPeriod() update vote err=", err.Error)
	//	return err.Error
	//}
	return nil
}

func (nft NftDb) DelSnftPeriod(delperiod string) error {
	if delperiod == "" {
		fmt.Println("params error")
		return errors.New("params error")
	}
	return nft.GetDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&SnftPhase{}).Where("tokenid= ?", delperiod).Delete(&SnftPhase{})
		if err.Error != nil {
			fmt.Println("delete snftPeriod err=", err.Error)
			return err.Error
		}
		err = nft.db.Model(&SnftCollectPeriod{}).Where(" period = ?", delperiod).Delete(&SnftCollectPeriod{})
		if err.Error != nil {
			fmt.Println("DelSnftPeriod() update  snftcollect err= ", err.Error)
			return err.Error
		}
		return nil
	})

}

func (nft NftDb) SnftPeriodVote(period string) error {
	if period == "" {
		fmt.Println("vote period id null")
		return errors.New("vote period is null")
	}
	var snftcollectrec SnftPhase
	err := nft.db.Model(&snftcollectrec).Where("tokenid = ?", period).First(&snftcollectrec)
	if err.Error != nil {
		fmt.Println("SnftPeriodVote() err= not find period.")
		return err.Error
	}
	snftcollectrec.Vote += 1
	err = nft.db.Model(&SnftPhase{}).Where("tokenid=", period).Updates(&snftcollectrec)
	if err.Error != nil {
		fmt.Println("SnftPeriodVote() update vote err=", err.Error)
		return err.Error
	}
	return nil
}

func (nft NftDb) SetPeriodAccedEth(meta string) error {
	err := nft.db.Model(&SnftPhase{}).Where("meta =? ", meta).Update("accedeth", "true")
	if err.Error != nil {
		fmt.Println("SetPeriodAccedEth() update true err=", err.Error)
		return err.Error
	}
	err = nft.db.Model(&SnftPhase{}).Where("accedeth =? ", "false").Update("accedeth", "")
	if err.Error != nil {
		fmt.Println("SetPeriodAccedEth() update false err=", err.Error)
		return err.Error
	}
	return nil
}

func AutoPeriodEth(sqldsn string) {
	ScanAutoFlag := time.NewTicker(time.Hour * 9)
	AccedEth := time.NewTicker(time.Hour * 5)
	for {
		select {
		case <-ScanAutoFlag.C:
			if !AutocommitSnft {
				continue
			}
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("AutoPeriodEth() connect database err = %s\n", err)
				continue
			}
			percoll := SnftPhase{}
			perr := nd.db.Model(&SnftPhase{}).Where("accedvote <> '' and (accedeth is null or accedeth =?)", "").Limit(1).Order("vote desc").Find(&percoll)
			if perr.Error != nil {
				fmt.Println("AutoPeriodEth() found snftperiod  err=", perr.Error)
				continue
			}
			serr := nd.SetPeriodEth(percoll.Tokenid)
			if serr != nil {
				fmt.Println("AutoPeriodEth() period injection  err=", serr)
				continue
			}
			nd.Close()
		case <-AccedEth.C:
			nd, err := NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("AutoPeriodEth() connect database err = %s\n", err)
				continue
			}
			syncBlocks, merr := GetDbBlockNumber(sqldsn)
			if merr != nil {
				fmt.Println("AutoPeriodEth() get scan block num err=", merr)
				continue
			}
			snftInfo, err := contracts.GetNominatedNFTInfo(big.NewInt(int64(syncBlocks)))
			if err != nil {
				fmt.Println("AutoPeriodEth() get NominatedNFTInfo err=", merr)
				continue
			}
			fmt.Println(snftInfo.Dir)
			serr := nd.SetPeriodAccedEth(snftInfo.Dir)
			if serr != nil {
				fmt.Println("AutoPeriodEth() SetPeriodAccedEth  err=", serr)
				continue
			}
			nd.Close()

		}
	}
}
