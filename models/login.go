package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

type UserSyncMapList struct {
	Mux   sync.Mutex
	Users map[string]*sync.Mutex
	Trans map[string]*sync.Mutex
}

func (u *UserSyncMapList) Lock(userAddr string) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	if len(u.Users) == 0 {
		u.Users = make(map[string]*sync.Mutex)
	}
	user, ok := u.Users[userAddr]
	if !ok {
		u.Users[userAddr] = new(sync.Mutex)
		user, _ = u.Users[userAddr]
	}
	user.Lock()
}

func (u *UserSyncMapList) LockTran(userAddr string) bool {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	if len(u.Trans) == 0 {
		u.Trans = make(map[string]*sync.Mutex)
	}
	user, ok := u.Trans[userAddr]
	if !ok {
		u.Trans[userAddr] = new(sync.Mutex)
		user, _ = u.Trans[userAddr]
		user.Lock()
		return false
	} else {
		return true
	}

}
func (u *UserSyncMapList) UnLockTran(userAddr string) {
	user, ok := u.Trans[userAddr]
	if ok {
		user.Unlock()
		delete(u.Trans, userAddr)
	}
}
func (u *UserSyncMapList) UnLock(userAddr string) {
	user, ok := u.Users[userAddr]
	if ok {
		user.Unlock()
	}
}

var UserSync UserSyncMapList

func (nft NftDb) Login(userAddr, sigData string) error {
	userAddr = strings.ToLower(userAddr)
	UserSync.Lock(userAddr)
	defer UserSync.UnLock(userAddr)
	user := Users{}
	db := nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			if KYCUploadAuditRequired {
				user.Verified = NoVerify.String()
			} else {
				user.Verified = Passed.String()
			}
			user.Useraddr = userAddr
			user.Signdata = sigData
			user.Userlogin = time.Now().Unix()
			user.Userlogout = time.Now().Unix()
			user.Username = ""
			user.Userregd = time.Now().Unix()
			/*user.Portrait = PortraitImg
			imagerr := SavePortrait(ImageDir, userAddr, PortraitImg)
			if imagerr != nil {
				fmt.Println("loging() SavePortrait() err=", imagerr)
				return ErrPortraitImage
			}*/
			db = nft.db.Model(&user).Create(&user)
			if db.Error != nil {
				fmt.Println("loging()->create() err=", db.Error)
				return ErrLoginFailed
			}
			//GetRedisCatch().SetDirtyFlag(KYCListDirtyName)

		}
	} else {
		newUser := Users{}
		/*if user.Portrait == "" {
			imagerr := SavePortrait(ImageDir, userAddr, PortraitImg)
			if imagerr != nil {
				fmt.Println("loging() SavePortrait() err=", imagerr)
				return ErrPortraitImage
			}
			newUser.Portrait = PortraitImg
		}*/
		newUser.Userlogin = time.Now().Unix()
		db = nft.db.Model(&Users{}).Where("useraddr = ?", userAddr).Updates(&newUser)
		if db.Error != nil {
			fmt.Printf("login()->UPdate() users err=%s\n", db.Error)
			return ErrLoginFailed
		}
	}
	return db.Error
}

func (nft NftDb) LoginNew(userAddr, sigData string) error {
	userAddr = strings.ToLower(userAddr)

	userOld := Users{}
	db := nft.db.Model(&Users{}).Last(&userOld)
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		fmt.Println("login() look last err=", db.Error)
		return db.Error
	}
	fmt.Println("login()", "userOld.id= ", userOld.ID)
	user := Users{}
	db = nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			fmt.Println("log()", "userOld.id= ", userOld.ID)
			user.ID = userOld.ID + 1
			user.Useraddr = userAddr
			user.Signdata = sigData
			user.Userlogin = time.Now().Unix()
			user.Userlogout = time.Now().Unix()
			user.Username = ""
			user.Userregd = time.Now().Unix()
			db = nft.db.Model(&user).Create(&user)
			if db.Error != nil {
				fmt.Println("\"log num=\", num, loging()->create() err=", db.Error)
				return ErrLoginFailed
			}
			fmt.Println("log()", "user.id= ", user.ID)
			fmt.Println("log()", "userOld.id= ", userOld.ID)
		}
	} else {
		fmt.Println("log()", "find user.id= ", user.ID)
		db = nft.db.Model(&Users{}).Where("useraddr = ?", userAddr).Update("userlogin", time.Now().Unix())
		if db.Error != nil {
			fmt.Printf("login()->UPdate() users err=%s\n", db.Error)
			return ErrLoginFailed
		}
	}
	return db.Error
}
