package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type OverviewData struct {
	Ownaddr        string          `json:"ownaddr"`
	Name           string          `json:"name"`
	Users          OverviewUserNum `json:"users"`
	Trans          OverviewTran    `json:"trans"`
	NftCollection  OverviewNft     `json:"nft_collection"`
	Nft            OverviewNft     `json:"nft"`
	Snft           OverviewNft     `json:"snft"`
	SnftCollection OverviewNft     `json:"snft_collection"`
	Auction        OverviewAuction `json:"auction"`
}

type OverviewUserNum struct {
	KYCUser      int `json:"kyc_user"`
	UnKYCUser    int `json:"un_kyc_user"`
	CreatorUser  int `json:"creator_user"`
	OtherUser    int `json:"other_user"`
	ActiveUser   int `json:"active_user"`
	UnActiveUser int `json:"un_active_user"`
	AuctionUser  int `json:"auction_user"`
}

type OverviewTran struct {
	NftTrans     uint64 `json:"nft_trans"`
	NftTransNum  uint64 `json:"nft_trans_num"`
	SnftTrans    uint64 `json:"snft_trans"`
	SNftTransNum uint64 `json:"s_nft_trans_num"`
}

type OverviewNft struct {
	Total     int `json:"total"`
	MintTotal int `json:"mint_total"`
	Day       int `json:"day"`
	Week      int `json:"week"`
	TwoWeek   int `json:"two_week"`
	Month     int `json:"month"`
}

type DayMaketinfo struct {
	NftTrans     int `json:"nft_trans"`
	NftTransNum  int `json:"nft_trans_num"`
	NftAdd       int `json:"nft_add"`
	MintNftAdd   int `json:"mint_nft_add"`
	SnftTrans    int `json:"snft_trans"`
	SnftTransNum int `json:"snft_trans_num"`
}

type OverviewAuction struct {
	NftSell     int `json:"nft_sell"`
	NftDaySell  int `json:"nft_day_sell"`
	SnftSell    int `json:"snft_sell"`
	SnftDaySell int `json:"snft_day_sell"`
}

func (nft NftDb) GetOverview() (OverviewData, error) {
	var overview OverviewData
	var users []Users
	var usernft Nfts
	var nftlist []Nfts
	var trans []Trans
	var collects []Collects
	var auctions []Auction
	overview.Ownaddr = ExchangeOwer
	overview.Name = ExchangeName
	err := nft.db.Find(&users)
	if err.Error != nil {
		fmt.Println("GetOverview() find user err= ", err.Error)
		return OverviewData{}, ErrNftNotExist
	}
	for _, user := range users {
		if user.UpdatedAt.After(time.Now().AddDate(0, 0, -14)) {
			overview.Users.ActiveUser++
		} else {
			overview.Users.UnActiveUser++
		}
		if user.Verified == Passed.String() && user.Certifyimg != "" {
			overview.Users.KYCUser++
		} else {
			overview.Users.UnKYCUser++
		}
		err = nft.db.Model(&Nfts{}).Where("createaddr=?", user.Useraddr).First(&usernft)
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				overview.Users.OtherUser++
			}
		} else {
			overview.Users.CreatorUser++
		}
	}
	err = nft.db.Find(&trans)
	if err.Error != nil {
		fmt.Println("GetOverview() find trans err=  ", err.Error)
		return OverviewData{}, ErrNftNotExist
	}
	for _, trandata := range trans {
		switch trandata.Nfttype {
		case "snft":
			overview.Trans.SNftTransNum++
			overview.Trans.SnftTrans += trandata.Price
		case "nft":
			overview.Trans.NftTransNum++
			overview.Trans.NftTrans += trandata.Price
		}

	}
	err = nft.db.Find(&nftlist)
	if err.Error != nil {
		fmt.Println("GetOverview() find nft err=  ", err.Error)
		return OverviewData{}, ErrNftNotExist
	}
	for _, nftdata := range nftlist {
		nftjudge := true
		if nftdata.Nftaddr != "" && nftdata.Nftaddr[0:3] == "0x8" {
			nftjudge = false
		}
		if !nftjudge {
			if nftdata.Mintstate == Minted.String() {
				overview.Snft.MintTotal++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, 0, -1)) {
				overview.Snft.Day++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, 0, -7)) {
				overview.Snft.Week++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, 0, -14)) {
				overview.Snft.TwoWeek++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, -1, 0)) {
				overview.Snft.Month++
			}
			overview.Snft.Total++
		} else {
			if nftdata.Mintstate == Minted.String() {
				overview.Nft.MintTotal++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, 0, -1)) {
				overview.Nft.Day++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, 0, -7)) {
				overview.Nft.Week++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, 0, -14)) {
				overview.Nft.TwoWeek++
			}
			if nftdata.UpdatedAt.After(time.Now().AddDate(0, -1, 0)) {
				overview.Nft.Month++
			}
			overview.Nft.Total++

		}

	}

	err = nft.db.Find(&collects)
	if err.Error != nil {
		fmt.Println("GetOverview() find collects err=  ", err.Error)
		return OverviewData{}, ErrNftNotExist
	}
	for _, collect := range collects {
		nftjudge := true
		if collect.Contracttype == "snft" {
			nftjudge = false
		}
		if !nftjudge {
			if collect.UpdatedAt.After(time.Now().AddDate(0, 0, -1)) {
				overview.SnftCollection.Day++
			}
			if collect.UpdatedAt.After(time.Now().AddDate(0, 0, -7)) {
				overview.SnftCollection.Week++
			}
			if collect.UpdatedAt.After(time.Now().AddDate(0, 0, -14)) {
				overview.SnftCollection.TwoWeek++
			}
			if collect.UpdatedAt.After(time.Now().AddDate(0, -1, 0)) {
				overview.SnftCollection.Month++
			}
			overview.SnftCollection.Total++
		} else {
			if collect.UpdatedAt.After(time.Now().AddDate(0, 0, -1)) {
				overview.NftCollection.Day++
			}
			if collect.UpdatedAt.After(time.Now().AddDate(0, 0, -7)) {
				overview.NftCollection.Week++
			}
			if collect.UpdatedAt.After(time.Now().AddDate(0, 0, -14)) {
				overview.NftCollection.TwoWeek++
			}
			if collect.UpdatedAt.After(time.Now().AddDate(0, -1, 0)) {
				overview.NftCollection.Month++
			}
			overview.NftCollection.Total++
		}

	}

	err = nft.db.Find(&auctions)
	if err.Error != nil {
		fmt.Println("GetOverview() find auctions err=  ", err.Error)
		return OverviewData{}, ErrNftNotExist
	}
	for _, auction := range auctions {
		nftjudge := true
		if auction.Nftaddr != "" && auction.Nftaddr[0:3] == "0x8" {
			nftjudge = false
		}
		if !nftjudge {
			overview.Auction.SnftSell++
			if auction.UpdatedAt.After(time.Now().AddDate(0, 0, -1)) {
				overview.Auction.SnftDaySell++
			}
		} else {
			overview.Auction.NftSell++
			if auction.UpdatedAt.After(time.Now().AddDate(0, 0, -1)) {
				overview.Auction.NftDaySell++
			}
		}
	}

	sql := "select toaddr from (select * from trans)as aa group by toaddr"
	var addr []string
	err = nft.db.Raw(sql).Scan(&addr)
	if err.Error != nil {
		fmt.Println("GetOverview() find trans err=  ", err.Error)
		return OverviewData{}, ErrNftNotExist
	}
	overview.Users.AuctionUser = len(addr)
	return overview, nil

}

type OverviewExcel struct {
	Ownaddr    string `json:"ownaddr"`
	Name       string `json:"name"`
	Trans      string `json:"trans"`
	TransNum   int64  `json:"trans_num"`
	User       int64  `json:"user"`
	ActiveUser int64  `json:"active_user"`
	KYCUser    int64  `json:"kyc_user"`
	NftCreator int64  `json:"nft_creator"`
	Collection int64  `json:"collection"`
	Nft        int64  `json:"nft"`
	NftMint    int64  `json:"nft_mint"`
	NftAuction int64  `json:"nft_auction"`
}

func (nft NftDb) GetOverviewExcel() (OverviewExcel, error) {
	var overview OverviewExcel
	overview.Ownaddr = ExchangeOwer
	overview.Name = ExchangeName
	var trans []Trans
	db := nft.db.Model(&Trans{}).Where("1=1").Find(&trans)
	if db.Error != nil {
		log.Println("GetOverviewExcel find trans err=", db.Error)
		return OverviewExcel{}, db.Error
	}
	var transamount uint64
	for _, singetran := range trans {
		transamount += singetran.Price
	}
	overview.Trans = fmt.Sprintf("%.9f", float64(transamount)/1000000000)
	overview.TransNum = int64(len(trans))
	var users []Users
	db = nft.db.Model(&Users{}).Where("1=1").Find(&users)
	if db.Error != nil {
		log.Println("GetOverviewExcel find user err=", db.Error)
		return OverviewExcel{}, db.Error
	}
	for _, singeuser := range users {
		overview.User++
		if singeuser.UpdatedAt.After(time.Now().AddDate(0, 0, -14)) {
			overview.ActiveUser++
		}
		if singeuser.Verified == Passed.String() {
			overview.KYCUser++
		}
	}

	sql := `select createaddr from (SELECT * FROM nfts ) as mm GROUP BY createaddr`
	var creatorlist []string
	db = nft.db.Raw(sql).Scan(&creatorlist)
	overview.NftCreator = int64(len(creatorlist))
	var nfts []Nfts
	db = nft.db.Model(&Nfts{}).Where("1=1").Find(&nfts)
	if db.Error != nil {
		log.Println("GetOverviewExcel find nfts err=", db.Error)
		return OverviewExcel{}, db.Error
	}
	overview.Nft = int64(len(nfts))
	for _, singenft := range nfts {
		if singenft.Mintstate == Minted.String() {
			overview.NftMint++
		}
		if singenft.Selltype != "NotSale" {
			overview.NftAuction++
		}
	}
	var collectnum int64
	db = nft.db.Model(&Collects{}).Where("1=1").Count(&collectnum)
	if db.Error != nil {
		log.Println("GetOverviewExcel find collection err=", db.Error)
		return OverviewExcel{}, db.Error
	}
	overview.Collection = collectnum
	return overview, nil

}

func (nft NftDb) GetSnftPeriodNum(addr string) (int, error) {
	var resList []string
	sql := `select snftstage from (SELECT * FROM collects where createaddr = ?  )as aa group by snftstage`
	err := nft.db.Raw(sql, addr).Scan(&resList)
	if err.Error != nil {
		log.Println("GetSnftPeriodNum find snft err=", err.Error)
		return 0, err.Error
	}
	resCount := len(resList)
	return resCount, nil
}

func (nft NftDb) GetSnftPeledge() (int, int, error) {

	var plecount int64
	var resList []string
	sql := `select ownaddr from (SELECT * FROM nfts where pledgestate ="Pledge" and deleted_at is null) as mm GROUP BY ownaddr`
	err := nft.db.Raw(sql).Scan(&resList)
	if err.Error != nil {
		log.Println("GetSnftPeledge find snft err=", err.Error)
		return 0, 0, err.Error
	}
	resCount := len(resList)
	err = nft.db.Model(&Nfts{}).Where("pledgestate = ?", "Pledge").Count(&plecount)
	if err.Error != nil {
		log.Println("GetSnftPeledge find snft number err=", err.Error)
		return 0, 0, err.Error
	}
	return resCount, int(plecount), nil
}
