package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type HomePageReq struct {
	Announcement []string        `json:"announcement"`
	NftLoop      []NftLoopKey    `json:"nft_loop"`
	Collections  []CollectionKey `json:"collections"`
	NftList      []NftKey        `json:"nfts"`
}
type NftLoopKey struct {
	Contract string `json:"contract"`
	Tokenid  string `json:"tokenid"`
}
type CollectionKey struct {
	Creator string `json:"creator"`
	Name    string `json:"name"`
}
type NftKey struct {
	Contract string `json:"contract"`
	Tokenid  string `json:"tokenid"`
}

type HomePageResp struct {
	Announcement []Announcements `json:"announcement"`
	NftLoop      []Nfts          `json:"nft_loop"`
	Collections  []Collects      `json:"collections"`
	NftList      []Nfts          `json:"nfts"`
	Total        int64           `json:"total"`
}

type HomePageCatch struct {
	AnnouncesMux      sync.Mutex
	AnnouncesFlag     bool
	Announces         []Announcements
	HomePageFlashFlag bool
	HomePageFlashMux  sync.Mutex
	NftLoopMux        sync.Mutex
	NftLoopFlag       bool
	NftLoop           []Nfts
	CollectsMux       sync.Mutex
	CollectsFlag      bool
	Collects          []Collects
	NftListMux        sync.Mutex
	NftListFlag       bool
	NftList           []Nfts
	NftCountMux       sync.Mutex
	NftCountFlag      bool
	NftCount          int64
}

//func (h *HomePageCatch) HomePageFlashLock() {
//	h.HomePageFlashMux.Lock()
//}
//
//func (h *HomePageCatch) HomePageFlashUnLock() {
//	h.HomePageFlashMux.Unlock()
//}
//
//func (h *HomePageCatch) AnnouncesLock() {
//	h.AnnouncesMux.Lock()
//}
//
//func (h *HomePageCatch) AnnouncesUnLock() {
//	h.AnnouncesMux.Unlock()
//}
//
//func (h *HomePageCatch) NftLoopLock() {
//	h.NftLoopMux.Lock()
//}
//
//func (h *HomePageCatch) NftLoopUnLock() {
//	h.NftLoopMux.Unlock()
//}
//
//func (h *HomePageCatch) CollectsLock() {
//	h.NftLoopMux.Lock()
//}
//
//func (h *HomePageCatch) CollectsUnLock() {
//	h.CollectsMux.Unlock()
//}
//
//func (h *HomePageCatch) NftListLock() {
//	h.NftListMux.Lock()
//}
//
//func (h *HomePageCatch) NftListUnLock() {
//	h.NftListMux.Unlock()
//}
//
//func (h *HomePageCatch) NftCountLock() {
//	h.NftCountMux.Lock()
//}
//
//func (h *HomePageCatch) NftCountUnLock() {
//	h.NftCountMux.Unlock()
//}

var HomePageCatchs HomePageCatch

func (nft *NftDb) QueryHomePage(flashFlag bool) ([]HomePageResp, error) {
	sysParams := SysParams{}
	spendT := time.Now()

	homePageResp := HomePageResp{}

	//HomePageCatchs.HomePageFlashLock()
	result := nft.db.Model(&SysParams{}).Select("homepage").Last(&sysParams)
	if result.Error != nil {
		log.Println("QueryHomePage() select homepage err=", result.Error)
		return nil, ErrData
	}
	var homePageReq HomePageReq
	err := json.Unmarshal([]byte(sysParams.Homepage), &homePageReq)
	if err != nil {
		log.Println("homePageReq  unmarshal err =", err)
		return nil, ErrDataFormat
	}
	fmt.Println("Unmarshal homepage:", homePageReq)
	cerr := GetRedisCatch().GetCatchData("QueryHomePage", sysParams.Homepage, &homePageResp)
	if cerr == nil {
		log.Printf("QueryHomePage() NftLoop  default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	} else {
		for _, v := range homePageReq.NftLoop {
			nftData := Nfts{}
			result := nft.db.Model(&Nfts{}).Select([]string{"contract", "tokenid", "collectcreator", "Collections", "desc", "name",
				"ownaddr", "createaddr", "url", "snft", "transprice"}).Where("contract = ? and tokenid = ?", v.Contract, v.Tokenid).
				First(&nftData)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				log.Println("QueryHomePage() NftLoop select err =", result.Error)
				return nil, ErrDataBase
			}
			if nftData.Snft == "" {
				nftData.Url = ""
			}
			homePageResp.NftLoop = append(homePageResp.NftLoop, nftData)
		}
		fmt.Println("homePageResp NftLoop:", len(homePageResp.NftLoop))
		if len(homePageResp.NftLoop) > 0 {
			HomePageCatchs.NftLoop = homePageResp.NftLoop
		} else {
			homePageResp.NftLoop = []Nfts{
				Nfts{gorm.Model{}, NftRecord{Tokenid: "", Contract: ""}},
			}
			HomePageCatchs.NftLoop = homePageResp.NftLoop
		}

		for _, v := range homePageReq.Collections {
			collectData := Collects{}
			result := nft.db.Model(&Collects{}).Select([]string{"categories", "createaddr", "desc", "name", "contract", "contracttype", "img", "totalcount"}).Where("createaddr = ? and name = ?", v.Creator, v.Name).
				First(&collectData)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return nil, ErrDataBase
			}
			if collectData.Contracttype != "snft" {
				collectData.Img = ""
			}
			homePageResp.Collections = append(homePageResp.Collections, collectData)
		}
		if len(homePageResp.Collections) > 0 {
			HomePageCatchs.Collects = homePageResp.Collections
		} else {
			homePageResp.Collections = []Collects{
				Collects{gorm.Model{}, CollectRec{Createaddr: "", Name: ""}},
			}
			HomePageCatchs.Collects = homePageResp.Collections
		}

		for _, v := range homePageReq.NftList {
			nftData := Nfts{}
			result := nft.db.Model(&Nfts{}).Select([]string{"contract", "tokenid", "collectcreator", "Collections", "desc", "name",
				"ownaddr", "createaddr", "url", "snft", "transprice"}).Where("contract = ? and tokenid = ?", v.Contract, v.Tokenid).
				First(&nftData)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return nil, ErrDataBase
			}
			if nftData.Snft == "" {
				nftData.Url = ""
			}
			homePageResp.NftList = append(homePageResp.NftList, nftData)
		}
		if len(homePageResp.NftList) > 0 {
			HomePageCatchs.NftList = homePageResp.NftList
		} else {
			homePageResp.NftList = []Nfts{
				Nfts{gorm.Model{}, NftRecord{Tokenid: "", Contract: ""}},
			}
			HomePageCatchs.NftList = homePageResp.NftList
		}

		sysInfo := SysInfos{}
		results := nft.db.Model(&SysInfos{}).Last(&sysInfo)
		if results.Error != nil {
			if results.Error != gorm.ErrRecordNotFound {
				log.Println("homepage() SysInfos err=", results.Error)
				return nil, ErrDataBase
			}
		}
		homePageResp.Total = int64(sysInfo.Nfttotal)
		//GetRedisCatch().CatchQueryData("QueryHomePage", "Nfttotal", &homePageResp.Total)

		GetRedisCatch().CatchQueryData("QueryHomePage", sysParams.Homepage, &homePageResp)
	}

	//marshaldata, err = json.Marshal(homePageReq.NftList)
	//if err != nil {
	//	log.Println("homePageReq.NftList marshal err =", err)
	//	return nil, ErrDataFormat
	//}
	//cerr = GetRedisCatch().GetCatchData("QueryHomePage", string(marshaldata), &homePageResp.NftList)
	//if cerr == nil {
	//	log.Printf("QueryHomePage() Collections default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	//} else {
	//
	//}

	//HomePageCatchs.HomePageFlashFlag = false

	//if HomePageCatchs.HomePageFlashFlag || flashFlag {
	//
	//} else {
	//	homePageResp.NftLoop = HomePageCatchs.NftLoop
	//	homePageResp.Collections = HomePageCatchs.Collects
	//	homePageResp.NftList = HomePageCatchs.NftList
	//}
	//HomePageCatchs.HomePageFlashUnLock()

	cerr = GetRedisCatch().GetCatchData("Announcement", "Announcement", &homePageResp.Announcement)
	if cerr == nil {
		log.Printf("QueryHomePage() Announcement default spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	} else {
		announcementList, err := nft.QueryAnnouncement()
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, ErrDataBase
		}
		homePageResp.Announcement = append(homePageResp.Announcement, announcementList...)
		GetRedisCatch().CatchQueryData("Announcement", "Announcement", &homePageResp.Announcement)

	}

	//HomePageCatchs.Announces = announcementList
	//HomePageCatchs.AnnouncesFlag = false
	//HomePageCatchs.AnnouncesLock()
	//if HomePageCatchs.AnnouncesFlag || flashFlag {
	//
	//} else {
	//	homePageResp.Announcement = HomePageCatchs.Announces
	//}
	//HomePageCatchs.AnnouncesUnLock()
	//
	//HomePageCatchs.NftCountLock()

	//HomePageCatchs.NftCount = homePageResp.Total
	//HomePageCatchs.NftCountFlag = false
	//if HomePageCatchs.NftCountFlag || flashFlag {
	//	//nftRec := Nfts{}
	//	//result := nft.db.Model(&Nfts{}).Select("id").Last(&nftRec)
	//	//if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
	//	//	return nil, result.Error
	//	//}
	//
	//} else {
	//	homePageResp.Total = HomePageCatchs.NftCount
	//}
	//HomePageCatchs.NftCountUnLock()

	return []HomePageResp{homePageResp}, nil
}
