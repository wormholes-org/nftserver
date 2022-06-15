module github.com/nftexchange/nftserver/common/contracts

go 1.15

require (
	//github.com/nftexchange/nftserver/models v0.0.0
	github.com/ethereum/go-ethereum v1.10.9
	github.com/nftexchange/nftserver/ethhelper v0.0.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
)

replace (
	//github.com/nftexchange/nftserver/models v0.0.0 => ../../models
	//github.com/nftexchange/nftserver/controllers v0.0.0 => ../../controllers
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../../ethhelper/database
)
