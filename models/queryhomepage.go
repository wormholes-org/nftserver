package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"sync"
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

func (h *HomePageCatch) HomePageFlashLock() {
	h.HomePageFlashMux.Lock()
}

func (h *HomePageCatch) HomePageFlashUnLock() {
	h.HomePageFlashMux.Unlock()
}

func (h *HomePageCatch) AnnouncesLock() {
	h.AnnouncesMux.Lock()
}

func (h *HomePageCatch) AnnouncesUnLock() {
	h.AnnouncesMux.Unlock()
}

func (h *HomePageCatch) NftLoopLock() {
	h.NftLoopMux.Lock()
}

func (h *HomePageCatch) NftLoopUnLock() {
	h.NftLoopMux.Unlock()
}

func (h *HomePageCatch) CollectsLock() {
	h.NftLoopMux.Lock()
}

func (h *HomePageCatch) CollectsUnLock() {
	h.CollectsMux.Unlock()
}

func (h *HomePageCatch) NftListLock() {
	h.NftListMux.Lock()
}

func (h *HomePageCatch) NftListUnLock() {
	h.NftListMux.Unlock()
}

func (h *HomePageCatch) NftCountLock() {
	h.NftCountMux.Lock()
}

func (h *HomePageCatch) NftCountUnLock() {
	h.NftCountMux.Unlock()
}

var HomePageCatchs HomePageCatch

func (nft *NftDb) QueryHomePage(flashFlag bool) ([]HomePageResp, error) {
	sysParams := SysParams{}

	homePageResp := HomePageResp{}
	HomePageCatchs.HomePageFlashLock()
	if HomePageCatchs.HomePageFlashFlag || flashFlag {
		result := nft.db.Model(&SysParams{}).Select("homepage").Last(&sysParams)
		if result.Error != nil {
			return nil, result.Error
		}
		var homePageReq HomePageReq
		err := json.Unmarshal([]byte(sysParams.Homepage), &homePageReq)
		if err != nil {
			return nil, err
		}

		for _, v := range homePageReq.NftLoop {
			nftData := Nfts{}
			result := nft.db.Model(&Nfts{}).Select([]string{"contract", "tokenid", "collectcreator", "Collections", "desc", "name",
				"ownaddr", "createaddr", "url", "snft"}).Where("contract = ? and tokenid = ?", v.Contract, v.Tokenid).
				First(&nftData)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return nil, result.Error
			}
			if nftData.Snft == "" {
				nftData.Url = ""
			}
			homePageResp.NftLoop = append(homePageResp.NftLoop, nftData)
		}
		if len(homePageResp.NftLoop) > 0 {
			HomePageCatchs.NftLoop = homePageResp.NftLoop
		}

		for _, v := range homePageReq.Collections {
			collectData := Collects{}
			result := nft.db.Model(&Collects{}).Select([]string{"categories", "createaddr", "desc", "name", "contract", "contracttype", "img", "totalcount"}).Where("createaddr = ? and name = ?", v.Creator, v.Name).
				First(&collectData)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return nil, result.Error
			}
			if collectData.Contracttype != "snft" {
				collectData.Img = ""
			}
			homePageResp.Collections = append(homePageResp.Collections, collectData)
		}
		if len(homePageResp.Collections) > 0 {
			HomePageCatchs.Collects = homePageResp.Collections
		}

		for _, v := range homePageReq.NftList {
			nftData := Nfts{}
			result := nft.db.Model(&Nfts{}).Select([]string{"contract", "tokenid", "collectcreator", "Collections", "desc", "name",
				"ownaddr", "createaddr", "url", "snft"}).Where("contract = ? and tokenid = ?", v.Contract, v.Tokenid).
				First(&nftData)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				return nil, result.Error
			}
			if nftData.Snft == "" {
				nftData.Url = ""
			}
			homePageResp.NftList = append(homePageResp.NftList, nftData)
		}
		if len(homePageResp.NftList) > 0 {
			HomePageCatchs.NftList = homePageResp.NftList
		}
		HomePageCatchs.HomePageFlashFlag = false
	} else {
		homePageResp.NftLoop = HomePageCatchs.NftLoop
		homePageResp.Collections = HomePageCatchs.Collects
		homePageResp.NftList = HomePageCatchs.NftList
	}
	HomePageCatchs.HomePageFlashUnLock()

	HomePageCatchs.AnnouncesLock()
	if HomePageCatchs.AnnouncesFlag || flashFlag {
		announcementList, err := nft.QueryAnnouncement()
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		homePageResp.Announcement = append(homePageResp.Announcement, announcementList...)
		HomePageCatchs.Announces = announcementList
		HomePageCatchs.AnnouncesFlag = false
	} else {
		homePageResp.Announcement = HomePageCatchs.Announces
	}
	HomePageCatchs.AnnouncesUnLock()

	HomePageCatchs.NftCountLock()
	if HomePageCatchs.NftCountFlag || flashFlag {
		nftRec := Nfts{}
		result := nft.db.Model(&Nfts{}).Select("id").Last(&nftRec)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		homePageResp.Total = int64(nftRec.ID)
		HomePageCatchs.NftCount = homePageResp.Total
		HomePageCatchs.NftCountFlag = false
	} else {
		homePageResp.Total = HomePageCatchs.NftCount
	}
	HomePageCatchs.NftCountUnLock()

	return []HomePageResp{homePageResp}, nil
}
