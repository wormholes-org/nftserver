package models

import (
	"gorm.io/gorm"
)

type SysInfoRec struct {
	Nfttotal     uint64 `json:"nfttotal" gorm:"type:bigint unsigned DEFAULT 0;comment:'nft总数量'"`
	Snfttotal    uint64 `json:"Snfttotal" gorm:"type:bigint unsigned DEFAULT 0;comment:'snft总数量'"`
	Nfthour      uint64 `json:"nfthour" gorm:"type:int unsigned DEFAULT 0 ;comment:'每小时nft添加数'"`
	Snfthour     uint64 `json:"snfthour" gorm:"type:bigint unsigned DEFAULT 0;comment:'每小时snft增加数'"`
	Transcnthour uint64 `json:"transcnthour" gorm:"type:int unsigned DEFAULT 0 ;comment:'每小时交易次数'"`
	Transamthour uint64 `json:"transamthour" gorm:"type:bigint unsigned DEFAULT 0;comment:'每小时交易总金额'"`
	Transcntday  uint64 `json:"transcntday" gorm:"type:int unsigned DEFAULT 0 ;comment:'日交易次数'"`
	Transamtday  uint64 `json:"transamtday" gorm:"type:bigint unsigned DEFAULT 0;comment:'日交易总金额'"`
}

type SysInfos struct {
	gorm.Model
	SysInfoRec
}

func (v SysInfos) TableName() string {
	return "sysinfos"
}

