package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

func (nft NftDb) QueryOwnerSnftChipAmount(owner, Categories string) (int64, error) {
	spendT := time.Now()
	var recount int64
	err := nft.db.Model(&Nfts{}).Where("Ownaddr = ? and snft != \"\"", owner).Count(&recount)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		log.Println("QueryOwnerSnftChipAmount() Scan(&stageCollection) err=", err)
		return 0, ErrDataBase
	}
	fmt.Printf("QueryOwnerSnftChipAmount() spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return recount, nil
}
