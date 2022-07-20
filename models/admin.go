package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

type Adminrec struct {
	Adminaddr string `json:"adminaddr" gorm:"type:char(42) NOT NULL;comment:'administrator address'"`
	Admintype string `json:"admintype" gorm:"type:char(42) NOT NULL;comment:'administrator type,nft,kyc,admin'"`
	AdminAuth string `json:"adminauth" gorm:"type:char(42) NOT NULL;comment:'admin rights'"`
}

type Admins struct {
	gorm.Model
	Adminrec
}

func (v Admins) TableName() string {
	return "admins"
}

type AdminType int

const (
	AdminTypeNft AdminType = iota
	AdminTypeKyc
	AdminTypeAdmin
)

func (this AdminType) String() string {
	switch this {
	case AdminTypeNft:
		return "nft"
	case AdminTypeKyc:
		return "kyc"
	case AdminTypeAdmin:
		return "admin"
	default:
		return "Unknow"
	}
}

type AdminAuthType int

const (
	AdminBrowse AdminAuthType = iota
	AdminEdit
	AdminAudit
	_
	_
	_
	AdminBrowseEditAudit
)

func (this AdminAuthType) String() string {
	switch this {
	case AdminBrowse:
		return "AuthBrowse"
	case AdminEdit:
		return "AuthEdit"
	case AdminAudit:
		return "AuthAudit"
	case AdminBrowseEditAudit:
		return "AuthBrowseEditAudit"
	default:
		return "Unknow"
	}
}

type DelAdmiList struct {
	DelAdmins []string `json:"del_admins"`
}

func (nft NftDb) AdminTypeValid(adminType string) error {
	switch adminType {
	case "nft", "kyc", "admin":
		return nil
	default:
		return errors.New("Admin type error.")
	}
}

func (nft NftDb) AdminAuthValid(adminAuth string) error {
	auth, err := strconv.Atoi(adminAuth)
	if err != nil {
		return err
	}
	if AdminAuthType(auth) > AdminAudit {
		return errors.New("admin auth error.")
	}
	return nil
}

type QueryAdminsCache struct {
	Admin []Adminrec
	Total int
}

func (nft NftDb) QueryAdmins(adminType, start_index, count string) (int, []Adminrec, error) {
	err := nft.AdminTypeValid(adminType)
	if err != nil {
		fmt.Println("ModifyAdmin() input data error.")
		return 0, nil, ErrAminType
	}
	if IsIntDataValid(start_index) != true {
		return 0, nil, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return 0, nil, ErrDataFormat
	}

	admin := QueryAdminsCache{}
	//cerr := GetRedisCatch().GetCatchData("QueryAdmins", adminType+start_index+count, &admin)
	//if cerr == nil {
	//	log.Printf("QueryAdmins() default  time.now=%s\n", time.Now())
	//	return admin.Total, admin.Admin, nil
	//}
	var recCount int64
	dberr := nft.db.Model(Admins{}).Where("admintype = ? AND adminaddr != ? AND adminaddr != ? AND  admin_auth !=?",
		adminType, SuperAdminAddr, ExchangeOwer, AdminBrowseEditAudit).Count(&recCount)
	if dberr.Error != nil {
		if dberr.Error == gorm.ErrRecordNotFound {
			return 0, nil, nil
		}
		fmt.Println("QuerySingleAnnouncement() recCount err=", err)
		return 0, nil, ErrDataBase
	}

	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)
	admins := make([]Admins, 0, 20)
	db := nft.db.Model(&Admins{}).Where("admintype = ? AND adminaddr != ? AND adminaddr != ? AND  admin_auth !=?",
		adminType, SuperAdminAddr, ExchangeOwer, AdminBrowseEditAudit).Offset(startIndex).Limit(nftCount).Find(&admins)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return 0, nil, nil
		} else {
			fmt.Println("ModifyAdmin() dbase err=", db.Error)
			return 0, nil, ErrDataBase
		}
	} else {
		adminrec := make([]Adminrec, 0, 20)
		for _, admin := range admins {
			adminrec = append(adminrec, admin.Adminrec)
		}
		admin.Admin = adminrec
		admin.Total = int(recCount)
		//GetRedisCatch().CatchQueryData("QueryAdmins", adminType+start_index+count, &admin)

		return int(recCount), adminrec, nil
	}
}

func (nft NftDb) ModifyAdmin(Adminaddr, Admintype, AdminAuth string) error {
	Adminaddr = strings.ToLower(Adminaddr)
	UserSync.Lock(Adminaddr)
	defer UserSync.UnLock(Adminaddr)

	if Adminaddr == SuperAdminAddr || Adminaddr == ExchangeOwer {
		fmt.Println("ModifyAdmin() No permission to modify the administrator error.")
		return ErrPermission
	}
	err := nft.AdminAuthValid(AdminAuth)
	if err != nil {
		fmt.Println("ModifyAdmin() input admin auth data error.")
		return ErrAminType
	}
	err = nft.AdminTypeValid(Admintype)
	if err != nil {
		fmt.Println("ModifyAdmin() input admin type data error.")
		return ErrAminType
	}
	admin := Admins{}
	db := nft.db.Model(&admin).Where("adminaddr = ? AND admintype = ?", Adminaddr, Admintype).First(&admin)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			admin.Adminaddr = Adminaddr
			admin.Admintype = Admintype
			admin.AdminAuth = AdminAuth
			db = nft.db.Model(&Admins{}).Create(&admin)
			if db.Error != nil {
				fmt.Println("ModifyAdmin()->create() err=", db.Error)
				return errors.New(ErrDataBase.Error() + db.Error.Error())
			}
		} else {
			fmt.Println("ModifyAdmin() dbase err=", db.Error)
			return errors.New(ErrDataBase.Error() + db.Error.Error())
		}
	} else {
		admin := Admins{}
		//admin.Admintype = Admintype
		admin.AdminAuth = AdminAuth
		db = nft.db.Model(&Admins{}).Where("adminaddr = ? AND admintype = ?", Adminaddr, Admintype).Updates(&admin)
		if db.Error != nil {
			fmt.Printf("ModifyAdmin()->UPdate() users err=%s\n", db.Error)
			return errors.New(ErrDataBase.Error() + db.Error.Error())
		}
	}
	//GetRedisCatch().SetDirtyFlag(AdminDirtyName)

	return nil
}

type deladmins struct {
	Admins [][]string `json:"del_admins"`
}

func (nft NftDb) DelAdmins(delAdminlist string) error {
	var dellst [][]string
	err := json.Unmarshal([]byte(delAdminlist), &dellst)
	if err != nil {
		log.Println("DelAdmins() Unmarshal err=", err)
		return ErrDataFormat
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		for _, admin := range dellst {
			if admin[0] != SuperAdminAddr && admin[0] != ExchangeOwer {
				db := tx.Model(&Admins{}).Where("adminaddr = ? AND Admintype = ?", admin[0], admin[1]).Delete(&Admins{})
				if db.Error != nil {
					if db.Error != gorm.ErrRecordNotFound {
						fmt.Println("DelAdmins() delete auction record err=", db.Error)
						return errors.New(ErrDataBase.Error() + db.Error.Error())
					}
				}
			}
		}
		//GetRedisCatch().SetDirtyFlag(AdminDirtyName)
		return nil
	})
}

func (nft NftDb) SetExchangerAdmin(Adminaddr string) error {
	Adminaddr = strings.ToLower(Adminaddr)
	UserSync.Lock(Adminaddr)
	defer UserSync.UnLock(Adminaddr)

	admin := Admins{}
	db := nft.db.Model(&admin).Where("adminaddr = ? AND admintype = ?", Adminaddr, AdminTypeAdmin.String()).First(&admin)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			admin.Adminaddr = Adminaddr
			admin.Admintype = AdminTypeAdmin.String()
			auth := strconv.Itoa(int(AdminBrowseEditAudit))
			admin.AdminAuth = auth
			db = nft.db.Model(&Admins{}).Create(&admin)
			if db.Error != nil {
				fmt.Println("SetExchangerAdmin()->create() admin err=", db.Error)
				return errors.New(ErrDataBase.Error() + db.Error.Error())
			}
		} else {
			fmt.Println("SetExchangerAdmin() dbase err=", db.Error)
			return errors.New(ErrDataBase.Error() + db.Error.Error())
		}
	}
	admin = Admins{}
	db = nft.db.Model(&admin).Where("adminaddr = ? AND admintype = ?", Adminaddr, AdminTypeNft.String()).First(&admin)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			admin.Adminaddr = Adminaddr
			admin.Admintype = AdminTypeNft.String()
			auth := strconv.Itoa(int(AdminBrowseEditAudit))
			admin.AdminAuth = auth
			db = nft.db.Model(&Admins{}).Create(&admin)
			if db.Error != nil {
				fmt.Println("SetExchangerAdmin()->create() admin err=", db.Error)
				return errors.New(ErrDataBase.Error() + db.Error.Error())
			}
		} else {
			fmt.Println("SetExchangerAdmin() dbase err=", db.Error)
			return errors.New(ErrDataBase.Error() + db.Error.Error())
		}
	}
	admin = Admins{}
	db = nft.db.Model(&admin).Where("adminaddr = ? AND admintype = ?", Adminaddr, AdminTypeKyc.String()).First(&admin)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			admin.Adminaddr = Adminaddr
			admin.Admintype = AdminTypeKyc.String()
			auth := strconv.Itoa(int(AdminBrowseEditAudit))
			admin.AdminAuth = auth
			db = nft.db.Model(&Admins{}).Create(&admin)
			if db.Error != nil {
				fmt.Println("SetExchangerAdmin()->create() admin err=", db.Error)
				return errors.New(ErrDataBase.Error() + db.Error.Error())
			}
		} else {
			fmt.Println("SetExchangerAdmin() dbase err=", db.Error)
			return errors.New(ErrDataBase.Error() + db.Error.Error())
		}
	}
	return nil
}

func (nft NftDb) QueryAdminByAddr(adminaddr string) ([]Adminrec, error) {
	adminaddr = strings.ToLower(adminaddr)
	if adminaddr == "" {
		return nil, errors.New(adminaddr + " input params error")
	}
	adminrec := make([]Adminrec, 0, 20)

	//cerr := GetRedisCatch().GetCatchData("QueryAdmins", "QueryAdminByAddr"+adminaddr, &adminrec)
	//if cerr == nil {
	//	log.Printf("QueryAdminByAddr() default  time.now=%s\n", time.Now())
	//	return adminrec, nil
	//}
	admins := make([]Admins, 0, 20)
	db := nft.db.Model(&Admins{}).Where("adminaddr = ? ",
		adminaddr).Find(&admins)
	if db.Error != nil {
		fmt.Println("QueryAdminByAddr() dbase err=", db.Error)
		return nil, errors.New(ErrDataBase.Error() + db.Error.Error())
	} else {
		for _, admin := range admins {
			adminrec = append(adminrec, admin.Adminrec)
		}
		//GetRedisCatch().CatchQueryData("QueryAdmins", "QueryAdminByAddr"+adminaddr, &adminrec)

		return adminrec, nil
	}
}

func (nft NftDb) AdminLogin(adminaddr string) error {
	adminaddr = strings.ToLower(adminaddr)
	if adminaddr == "" {
		return errors.New(adminaddr + " input params error")
	}
	admins := Admins{}
	db := nft.db.Model(&Admins{}).Where("adminaddr = ? ",
		adminaddr).First(&admins)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			fmt.Println("AdminLogin() not found err=")
			return ErrPermission
		}
		fmt.Println("AdminLogin() dbase err=", db.Error)
		return errors.New(ErrDataBase.Error() + db.Error.Error())
	}
	return nil

}
