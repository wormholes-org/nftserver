module github.com/nftexchange/nftserver

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	github.com/ethereum/go-ethereum v1.10.9
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nftexchange/nftserver/common/contracts v0.0.0
	github.com/nftexchange/nftserver/common/sync v0.0.0
	//github.com/nftexchange/nftserver/ethhelper v0.0.0
	github.com/nftexchange/nftserver/models v0.0.0
	github.com/nftexchange/nftserver/routers v0.0.0
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/swaggo/swag v1.8.10
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
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
	//github.com/nftexchange/nftserver/ethhelper v0.0.0 => ./ethhelper
	//github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ./ethhelper/common
	//github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ./ethhelper/database
	github.com/nftexchange/nftserver/models v0.0.0 => ./models
	github.com/nftexchange/nftserver/routers v0.0.0 => ./routers
)
