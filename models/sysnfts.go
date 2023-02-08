package models

type SysnftRecord struct {
	Ownaddr        string `json:"ownaddr" gorm:"type:char(42) ;comment:'nft owner address'"`
	Md5            string `json:"md5" gorm:"type:longtext ;comment:'Picture md5 value'"`
	Name           string `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft name'"`
	Desc           string `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft description'"`
	Meta           string `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta information'"`
	Nftmeta        string `json:"nftmeta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta information,tokenid'"`
	Url            string `json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc raw data hold address'"`
	Contract       string `json:"nft_contract_addr" gorm:"type:char(42) ;comment:'contract address'"`
	Tokenid        string `json:"nft_token_id" gorm:"type:char(42) ;comment:'Uniquely identifies the nft flag'"`
	Nftaddr        string `json:"nft_address" gorm:"type:char(42) DEFAULT NULL;comment:'Chain of wormholes uniquely identifies the nft'"`
	Snftstage      string `json:"snftstage" gorm:"type:char(42) DEFAULT NULL;comment:'wormholes chain snft period'"`
	Snftcollection string `json:"snftcollection" gorm:"type:char(42) DEFAULT NULL;comment:'Wormholes chain snft collection'"`
	Snft           string `json:"snft" gorm:"type:char(42) ;comment:'wormholes chain snft'"`
	Count          int    `json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft sellable quantity'"`
	Approve        string `json:"approve" gorm:"type:longtext ;comment:'Authorize'"`
	Categories     string `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:''"`
	Collectcreator string `json:"collection_creator_addr" gorm:"type:char(42) ;comment:''"`
	Collections    string `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft classification'"`
	Image          string `json:"asset_sample" gorm:"type:longtext ;comment:'Collection creator address'"`
	Hide           string `json:"hide" gorm:"type:char(20) ;comment:'Whether to let others see'"`
	Signdata       string `json:"sig" gorm:"type:longtext ;comment:'Signature data, generated when created'"`
	Createaddr     string `json:"user_addr" gorm:"type:char(42) ;comment:'Create nft address'"`
	Verifyaddr     string `json:"vrf_addr" gorm:"type:char(42) ;comment:'Validator address'"`
	Currency       string `json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'Transaction currency'"`
	Price          uint64 `json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'Price at creation time'"`
	Royalty        int    `json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'royalty'"`
	Paychan        string `json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'trading channel'"`
	TransCur       string `json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'Transaction currency'"`
	Transprice     uint64 `json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'transaction price'"`
	Transtime      int64  `json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'Last trading time'"`
	Createdate     int64  `json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft creation time'"`
	Favorited      int    `json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'Follow count'"`
	Transcnt       int    `json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'The number of transactions, plus one for each transaction'"`
	Transamt       uint64 `json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'total transaction amount'"`
	Verified       string `json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'Whether the nft work has passed the review'"`
	Verifieddesc   string `json:"Verifieddesc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'Review description: Failed review description'"`
	Verifiedtime   int64  `json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'Review time'"`
	Selltype       string `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft transaction type'"`
	Mintstate      string `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'minting status'"`
	Pledgestate    string `json:"pledgestate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'Pledgestate status'"`
	Chipcount      int    `json:"chipcount" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'snft slice count.'"`
}

//type Sysnfts struct {
//	gorm.Model
//	SysnftRecord
//}
//
//func (v Sysnfts) TableName() string {
//	return "sysnfts"
//}
