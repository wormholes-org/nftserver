module github.com/nftexchange/nftserver/models

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	github.com/ethereum/go-ethereum v1.10.9
	github.com/ipfs/go-ipfs-api v0.3.0
	github.com/nftexchange/nftserver/common/contracts v0.0.0
	github.com/nftexchange/nftserver/ethhelper v0.0.0
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
	//github.com/fatih/structs v1.1.0
)

replace (
	github.com/nftexchange/nftserver/common/contracts v0.0.0 => ../common/contracts
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../ethhelper/database
)
