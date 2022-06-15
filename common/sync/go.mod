module github.com/nftexchange/nftserver/common/sync

go 1.15

require (
	github.com/ethereum/go-ethereum v1.10.9
	github.com/ipfs/go-ipfs-api v0.3.0 // indirect
	github.com/nftexchange/nftserver/common/contracts v0.0.0
	github.com/nftexchange/nftserver/models v0.0.0
)

replace (
	github.com/nftexchange/nftserver/common/contracts v0.0.0 => ../contracts
	github.com/nftexchange/nftserver/controllers v0.0.0 => ../../controllers
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../../ethhelper/database
	github.com/nftexchange/nftserver/models v0.0.0 => ../../models
)
