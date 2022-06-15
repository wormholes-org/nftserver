module github.com/nftexchange/nftserver

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	github.com/ethereum/go-ethereum v1.10.9
	github.com/fatih/structs v1.1.0 // indirect
	github.com/ipfs/go-ipfs-api v0.3.0 // indirect
	github.com/nftexchange/nftserver/common/contracts v0.0.0
	github.com/nftexchange/nftserver/common/signature v0.0.0
	github.com/nftexchange/nftserver/common/sync v0.0.0
	github.com/nftexchange/nftserver/controllers/nftexchangev2 v0.0.0 // indirect
	github.com/nftexchange/nftserver/ethhelper v0.0.0
	github.com/nftexchange/nftserver/models v0.0.0
	github.com/nftexchange/nftserver/routers v0.0.0 // indirect
	//github.com/status-im/keycard-go v0.0.0-20190316090335-8537d3370df4 // indirect
	gorm.io/gorm v1.21.15
)

replace (
	github.com/nftexchange/nftserver/common/contracts v0.0.0 => ./common/contracts
	github.com/nftexchange/nftserver/common/signature v0.0.0 => ./common/signature
	github.com/nftexchange/nftserver/common/sync v0.0.0 => ./common/sync
	github.com/nftexchange/nftserver/controllers v0.0.0 => ./controllers
	github.com/nftexchange/nftserver/controllers/nftexchangev1 v0.0.0 => ./controllers/nftexchangev1
	github.com/nftexchange/nftserver/controllers/nftexchangev2 v0.0.0 => ./controllers/nftexchangev2
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ./ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ./ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ./ethhelper/database
	github.com/nftexchange/nftserver/models v0.0.0 => ./models
	github.com/nftexchange/nftserver/routers v0.0.0 => ./routers
)
