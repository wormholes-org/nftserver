package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type UserInfo struct {
	Name string `json:"user_name"` //用户名
	//Portrait     	string 	`json:"portrait"`           //用户头像
	Email           string `json:"user_mail"`         //邮箱
	Bio             string `json:"user_info"`         //自我描述
	Verified        string `json:"verified"`          //是否通过验证
	NftCount        int    `json:"nft_count"`         //用户持有的NFT总数
	CreateCount     int    `json:"create_count"`      //用户创作的NFT总数
	OwnerCount      int    `json:"owner_count"`       //用户创作的NFT的拥有者数量
	TradeAmount     uint64 `json:"trade_amount"`      //用户创作的NFT的成交额,
	TradeAvgPrice   uint64 `json:"trade_avg_price"`   //用户创作的NFT均价,
	TradeFloorPrice uint64 `json:"trade_floor_price"` //用户创作的NFT最低价
	Identity        string `json:"identity"`          //用户身份
}

func (nft NftDb) QueryUserInfo(userAddr string) (UserInfo, error) {
	userAddr = strings.ToLower(userAddr)

	var uinfo UserInfo
	user := Users{}
	err := nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return UserInfo{}, nil
		} else {
			fmt.Println("QueryUserInfo() query users err=", err)
			return UserInfo{}, err.Error
		}
	}

	uinfo.Name = user.Username
	//uinfo.Portrait = user.Portrait
	uinfo.Email = user.Email
	uinfo.Bio = user.Bio
	uinfo.Verified = user.Verified
	var recCount int64
	err = nft.db.Model(Nfts{}).Where("ownaddr = ?", userAddr).Count(&recCount)
	if err.Error == nil {
		uinfo.NftCount = int(recCount)
	}
	err = nft.db.Model(Nfts{}).Where("createaddr = ?", userAddr).Count(&recCount)
	if err.Error == nil {
		uinfo.CreateCount = int(recCount)
	}
	err = nft.db.Model(Nfts{}).Where("createaddr = ? AND ownaddr != ?",
		userAddr, userAddr).Count(&recCount)
	if err.Error == nil {
		uinfo.OwnerCount = int(recCount)
	}

	/*type SumInfo struct {
		SumCount int
		SumPrice uint64
	}
	sum := SumInfo{}
	err = nft.db.Raw("SELECT SUM(Transcnt) as SumCount, SUM(Transamt) as SumPrice FROM nfts WHERE createaddr = ?", userAddr).Scan(&sum)
	if err.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
		return UserInfo{}, err.Error
	}
	uinfo.TradeAmount = sum.SumPrice
	if sum.SumCount != 0 {
		uinfo.TradeAvgPrice = sum.SumPrice / uint64(sum.SumCount)
	}

	var nftRec Nfts
	err = nft.db.Order("transprice desc").Where("createaddr = ?", userAddr).Last(&nftRec)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("QueryUserInfo() query statistics err=", err)
			return UserInfo{}, err.Error
		}
	}
	uinfo.TradeFloorPrice = nftRec.Transprice*/

	type TransInfo struct {
		TradeAmount     uint64
		TradeAvgPrice   float64
		TradeFloorPrice uint64
		TradeMaxPrice   uint64
		TradeCount      uint64
	}
	tInfo := TransInfo{}
	sql := "SELECT sum(trans.price) as TradeAmount, avg(trans.price) as TradeAvgPrice, " +
		"min(trans.price) as TradeFloorPrice, max(trans.price) as TradeMaxPrice, " +
		"COUNT(trans.price) AS TradeCount " +
		//"FROM trans" +" WHERE createaddr = ? AND selltype != ? AND selltype != ?"
		"FROM trans" + " WHERE ( trans.fromaddr = ? OR trans.toaddr = ?) AND selltype != ? AND selltype != ?"
	err = nft.db.Raw(sql, userAddr, userAddr, SellTypeMintNft.String(), SellTypeError.String()).Scan(&tInfo)
	if err.Error == nil {
		uinfo.TradeAmount = tInfo.TradeAmount
		uinfo.TradeAvgPrice = uint64(tInfo.TradeAvgPrice)
		uinfo.TradeFloorPrice = tInfo.TradeFloorPrice
	}

	admin := Admins{}
	err = nft.db.Model(&Admins{}).Where("adminaddr= ? and admintype= ?", userAddr, AdminTypeAdmin.String()).First(&admin)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			uinfo.Identity = "Normal"
		} else {
			fmt.Println("QueryUserInfo() query admin err=", err.Error)
			return UserInfo{}, err.Error
		}
	} else {
		switch admin.AdminAuth {
		case strconv.Itoa(int(AdminBrowse)):
			uinfo.Identity = "Normal"
		case strconv.Itoa(int(AdminEdit)):
			uinfo.Identity = "Admin"
		case strconv.Itoa(int(AdminAudit)):
			uinfo.Identity = "Admin"
		case strconv.Itoa(int(AdminBrowseEditAudit)):
			uinfo.Identity = "Owner"
		}
	}
	return uinfo, nil
}

func (nft NftDb) ModifyUserInfo(user_addr, user_name, portrait, background, user_mail, user_info, sig string) error {
	user_addr = strings.ToLower(user_addr)
	if len(user_name) > LenName {
		return ErrLenName
	}
	if len(user_mail) > LenEmail {
		return ErrLenEmail
	}
	fmt.Println("ModifyUserInfo() start.")
	user := Users{}
	err := nft.db.Model(&user).Where("useraddr = ?", user_addr).First(&user)
	if err.Error != nil {
		fmt.Println("ModifyUserInfo() err= not find user.")
		return err.Error
	}
	if !nft.UserKYCAduit(user_addr) {
		return ErrUserNotVerify
	}
	if user_name != "" {
		user.Username = user_name
	}
	if user_info != "" {
		user.Bio = user_info
	}
	if user_mail != "" {
		user.Email = user_mail
	}
	if portrait != "" {
		imagerr := SavePortrait(ImageDir, user_addr, portrait)
		if imagerr != nil {
			fmt.Println("ModifyUserInfo() SavePortrait() err=", imagerr)
			return ErrPortraitImage
		}
		user.Portrait = portrait
	}
	if background != "" {
		imagerr := SaveBackground(ImageDir, user_addr, background)
		if imagerr != nil {
			fmt.Println("ModifyUserInfo() SaveBackground() err=", imagerr)
			return ErrBackgroudImage
		}
		user.Background = background
	}
	if sig != "" {
		user.Signdata = sig
	}
	err = nft.db.Model(&Users{}).Where("useraddr = ?", user_addr).Updates(user)
	if err.Error != nil {
		fmt.Println("ModifyUserInfo() update err= ", err.Error)
		return err.Error
	}
	fmt.Println("ModifyUserInfo() Ok.")
	return err.Error
}
