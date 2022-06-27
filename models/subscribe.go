package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type SubscribeRec struct {
	Useraddr    string 	`json:"useraddr" gorm:"type:char(42) NOT NULL;comment:'User address'"`
	Email		string	`json:"Email" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'Subscription email address'"`
}

type Subscribes struct {
	gorm.Model
	SubscribeRec
}

func (v Subscribes) TableName() string {
	return "subscribes"
}

func (nft NftDb) QuerySubscribeEmails(start_index, count string) (int, interface{}, error) {
	if IsIntDataValid(start_index) != true {
		return 0, nil, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return 0, nil, ErrDataFormat
	}
	var recCount int64
	err := nft.db.Model(Subscribes{}).Count(&recCount)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			return 0, nil , nil
		}
		fmt.Println("QuerySingleAnnouncement() recCount err=", err)
		return 0, nil, err.Error
	}

	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)

	var emails []struct{
				Useraddr string
				Email string
			}
	err = nft.db.Model(&Subscribes{}).Select([]string{"useraddr", "email"}).Offset(startIndex).Limit(nftCount).Scan(&emails)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("QuerySingleAnnouncement() dbase err=", err.Error)
			return 0, nil, err.Error
		} else {
			fmt.Println("QuerySingleAnnouncement() not find err=", err.Error)
			return 0, nil, nil
		}
	}
	return int(recCount), emails, err.Error
}

func (nft NftDb) SetSubscribeEmail(useraddr, email string) error {
	var subRec Subscribes
	db := nft.db.Model(&Subscribes{}).Where("useraddr = ? and email = ?", useraddr, email).First(&subRec)
	if db.Error == nil {
		fmt.Println("SetSubscribeEmail() email already exist")
		return db.Error
	}
	subRec = Subscribes{}
	subRec.Email = email
	subRec.Useraddr = useraddr
	db = nft.db.Model(&Subscribes{}).Create(&subRec)
	if db.Error != nil {
		fmt.Println("SetSubscribeEmail()->create() err=", db.Error)
		return db.Error
	}
	return nil
}

func (nft NftDb) DelSubscribeEmail(useraddr, email string) (error)  {
	db := nft.db.Model(&Subscribes{}).Where("useraddr = ? and email = ?", useraddr, email).Delete(&Subscribes{})
	if db.Error != nil {
		if db.Error != gorm.ErrRecordNotFound {
			fmt.Println("DelSubscribeEmail() delete subscribe record err=", db.Error)
			return db.Error
		}
	}
	return nil
}
