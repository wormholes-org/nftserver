package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type AnnounceRec struct {
	Title   string `json:"title" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'title'"`
	Content string `json:"content" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'content'"`
}

type Announcements struct {
	gorm.Model
	AnnounceRec
}

func (v Announcements) TableName() string {
	return "announcements"
}

func (nft NftDb) QueryAnnouncement() ([]Announcements, error) {
	var announce []Announcements
	err := nft.db.Model(&Announcements{}).Select([]string{"title", "content"}).Order("id desc").Limit(5).Find(&announce)
	if err.Error != nil {
		fmt.Println("QueryAnnouncement() not find err=", err.Error)
		return nil, err.Error
	}
	return announce, err.Error
}

type Announces struct {
	ID uint `json:"id"`
	AnnounceRec
}

func (nft NftDb) QueryAnnouncements(start_index, count string) (int, []Announces, error) {
	if IsIntDataValid(start_index) != true {
		return 0, nil, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return 0, nil, ErrDataFormat
	}
	var recCount int64
	err := nft.db.Model(Announcements{}).Count(&recCount)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			return 0, nil, nil
		}
		fmt.Println("QuerySingleAnnouncement() recCount err=", err)
		return 0, nil, err.Error
	}

	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)

	var announces []Announcements
	err = nft.db.Offset(startIndex).Limit(nftCount).Find(&announces)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("QuerySingleAnnouncement() dbase err=", err.Error)
			return 0, nil, err.Error
		} else {
			fmt.Println("QuerySingleAnnouncement() not find err=", err.Error)
			return 0, nil, nil
		}
	}

	var retAns []Announces
	for _, announce := range announces {
		var ans Announces
		ans.AnnounceRec = announce.AnnounceRec
		ans.ID = announce.ID
		retAns = append(retAns, ans)
	}
	return int(recCount), retAns, err.Error
}

func (nft NftDb) SetAnnouncement(title, content string) error {
	var announce Announcements
	announce.Title = title
	announce.Content = content
	db := nft.db.Model(&Announcements{}).Create(&announce)
	if db.Error != nil {
		fmt.Println("SetAnnouncement()->create() err=", db.Error)
		return db.Error
	}
	HomePageCatchs.AnnouncesLock()
	HomePageCatchs.AnnouncesFlag = true
	HomePageCatchs.AnnouncesUnLock()
	return nil
}

func (nft NftDb) DelAnnounce(delAnnouncelist string) error {
	var dellst []int
	err := json.Unmarshal([]byte(delAnnouncelist), &dellst)
	if err != nil {
		fmt.Println("DelAnnounce() Unmarshal err=", err)
		return err
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range dellst {
			db := tx.Model(&Announcements{}).Where("id = ?", id).Delete(&Announcements{})
			if db.Error != nil {
				if db.Error != gorm.ErrRecordNotFound {
					fmt.Println("DelAnnounce() delete auction record err=", db.Error)
					return db.Error
				}
			}
		}
		return nil
	})
}
