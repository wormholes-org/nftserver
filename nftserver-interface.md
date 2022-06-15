# frontend

### POST  /v2/login  用户登录

```json
请求值
user_addr: 用户地址
approve_addr: 临时验证地址
sig: 验证用签名

Request body （application/json）
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "approve_addr": "0xabcdef",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
返回值
Hash:md5
Secret:题目
TimeStamp:时间戳
Token:令牌

Responses
200    Successful Response
Example Value
{ 
    "code": "200",
    "msg": "string",
    "data": { 
        Hash:"0xABCD", 
        Secret:"1234", 
        TimeStamp:123456, 
        Token:"abcdefg" 
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserInfo  获取用户基本信息

```json
请求值
user_addr: 用户地址

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900"
}
```

```json
返回值
{
 "user_name":"name",
 "portrait:"......"
 "user_mail":"abcd@mail.com",
 "user_info":"none"
 "nft_count":5,
 "create_count":5,
 "owner_count":5,
 "trade_amount":150000,
 "trade_avg_price":100000,
 "trade_floor_price":10000,
 "verified":"Passed"
}
user_name:用户名
portrait:用户头像
user_mail:邮箱
user_info:自我描述
nft_count: 用户持有的NFT总数
create_count: 用户创作的NFT总数
owner_count: 用户创作的NFT的拥有者数量
trade_amount: 用户NFT的成交额,
trade_avg_price: 用户参与交易的NFT均价,
trade_floor_price: 用户参与交易的NFT最低价
verified: 用户身份认证 [Passed:已经认证 NoPass:未认证]


Responses
200 Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}


422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserNFTList  获取用户NFT列表

```json
请求值
user_addr: 用户地址
start_index: 查询序号
count: 查询数量

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 "creator_addr":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "owner_addr":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "md5":"2128ba8f7822c9cde4a3737bfd6be582",
 "name":"noname",
 "desc":"none",
 "meta":"http://....",
 "source_url":'QmNLofb3sty95CXxkfJPVVJfpf6f97LkQaZ1kSwR8gnZJH',
 "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id":"1631753648255",
 "categories":"Art",
 "collection_creator_addr":"0x1234",
 "collections":"album01",
 "asset_sample":"",
 "hide":"false"
 }
 ],
 "total_count":999
}

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserCollectionList  获取用户NFT合集列表

```json
请求值
user_addr: 用户地址
start_index: 查询序号
count: 查询数量

Request body
application/json
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 "collection_creator_addr":"0x1234",
 "name":"New Collection",
 "img":"......",
 "contract_addr":"0x12345",
 "desc":'none',
 "categories":"Art",
 "total_count":10
 }
 ],
 "total_count":999
}
collection_creator_addr: 合集创建者地址
name: 合集名
img: 合集封面
contract_addr: 合约地址
desc: 合集描述
categories: 分类
total_count: 合集NFT总数

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
application/json
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserTradingHistory  获取用户交易历史

```json
请求值
- user_addr: 用户地址
- start_index: 查询序号
- count: 查询数量

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
        "data":[
                {
                        "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
                        "nft_token_id":"1631753648255",
                        "name":"My NFT",
                        "price":1000000,
                        "count":"1",
                        "from":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
                        "to":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
                        "selltype":"FixPrice",
                        "date":123456789,
                        "trade_hash":"0xdfdsfdf"
                }
        ],
        "total_count":999
}
返回交易结构
nft_contract_addr: NFT合约地址
nft_token_id： NFT token ID
nft_name:NFT 名称
price:成交价
count:成交数量
from:卖家地址
to:买家地址
selltype: 交易类型
date:成交时间日期 UTC时间戳
trade_hash: 交易Hash


Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserOfferList  获取用户NFT出价列表

```json
请求值
user_addr: 用户地址
start_index: 查询序号
count: 查询数量


Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id":"1631753648255",
 "name":"name",
 "price":100000,
 "count":1,
 "date":123456789
 }
 ],
 "total_count":999
}

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserBidList  获取用户出价列表

```json
请求值
user_addr: 用户地址
start_index: 查询序号
count: 查询数量

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id":"1631753648255",
 "name":"name",
 "price":100000,
 "count":1,
 "date":123456789
 }
 ],
 "total_count":999
}

Responses 
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryUserFavoriteList  获取用户NFT关注列表

```json
请求值
user_addr: 用户地址
start_index: 查询序号
count: 查询数量

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id":"1631753648255",
 "name":"noname",
 "asset_sample":"......",
 "collection_creator_addr":"0x1234"
 "collections:"noname",
 "img":"......"
 }
 ],
 "total_count":999
}
nft_contract_addr:
nft_token_id":"1631753648255",
name":"noname",
asset_sample":"......",
collections:"noname",
img:""

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/like  用户关注/取消关注NFT Toggle模式

```json
请求值
user_addr: 用户地址
nft_contract_addr:"所关注NFT的合约地址",
nft_token_id:"所关注NFT的tokenid"，
sig:"验证签名"

Request body
application/json
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "nft_contract_addr": "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id": "1631753648255",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/modifyUserInfo  编辑用户信息

```json
请求值
user_addr: 用户地址
user_name: 用户名称
portrait: 用户头像
user_mail: 用户邮箱
user_info: 用户信息
sig: 验证用签名

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "user_name": "user_name",
 "portrait": "......",
 "user_mail": "user@mail.com",
 "user_info": "user_info",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": [
        ""： "",
    ],
    "total_count": int #失败返回0
}
```

### POST  /v2/upload  上传NFT

```
请求值
user_addr":用户地址,
creator_addr:作者地址,
owner_addr:拥有者地址,
md5":NFT数据MD5,
name":NFT 名称，
desc":NFT 描述，
meta":元信息
source_url":源数据链接
nft_contract_addr":NFT合约地址
nft_token_id":NFT tokenID
categories":NFT分类
collections":NFT合集名
asset_sample":缩略信息
hide":"true"
royalty: NFT版税
count: NFT数量
sig":验证用签名
```

```json
Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "creator_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "owner_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "md5": "2128ba8f7822c9cde4a3737bfd6be582",
 "name": "noname",
 "desc": "none",
 "meta": "http://....",
 "source_url": "QmNLofb3sty95CXxkfJPVVJfpf6f97LkQaZ1kSwR8gnZJH",
 "nft_contract_addr": "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id": "1631753648255",
 "categories": "Art",
 "collections": "album01",
 "asset_sample": "",
 "hide": "false",
 "royalty": "10",
 "count": "1",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": [
        ""： "",
    ],
    "total_count": int #失败返回0
}
```

### POST  /v2/newCollections  新建合集

```json
请求值
user_addr: 用户地址,
name: 合集名,
img: 合集封面,
contract_type: 合约类型,
contract_addr: 合约地址,
desc: 描述简介,
categories: 合集分类,
sig: 数据验证签名

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "name": "New Collection",
 "img": "......",
 "contract_type": "721",
 "contract_addr": "0x12345",
 "desc": "none",
 "categories": "Art",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### 

### POST  /v2/modifyCollections  修改合集信息

```json
请求值
user_addr: 用户地址,
name: 合集名,
img: 合集封面,
contract_type: 合约类型,
contract_addr: 合约地址,
desc: 描述简介,
categories: 合集分类,
sig: 数据验证签名

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "name": "New Collection",
 "img": "......",
 "contract_type": "721",
 "contract_addr": "0x12345",
 "desc": "none",
 "categories": "Art",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryNFT  获取NFT信息

```json
请求值
nft_contract_addr: NFT合约地址
nft_token_id: NFT TokenID

Request body
Example Value
{
 "nft_contract_addr": "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id": "1631753648255"
}
```

```json
返回值
{
 "name":"noname",
 "creator_addr":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "creator_portrait":"0x123456",
 "owner_addr":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "owner_portrait":"0x123456",
 "md5":"2128ba8f7822c9cde4a3737bfd6be582",
 "asset_sample":"",
 "desc":"none",
 "meta":"http://....",
 "source_url":'QmNLofb3sty95CXxkfJPVVJfpf6f97LkQaZ1kSwR8gnZJH',
 "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id":"1631753648255",
 "categories":"Art",
 "collection_creator_addr":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "collections":"album01",
 "collection_desc":"abcd",
 "img":"......"，
 "verified": "Passed",
 "selltype": "NotSale",
 "mintstate": "NoMinted",
 "approve":" ",
 "royalty":10,
 "count":1,
 "likes":99,
 "auction":{
 "selltype": 售卖类型
 "ownaddr": 拥有者地址
 "nft_token_id: NFT TokenID
 "nft_contract_addr: NFT合约地址
 "paychan: 交易通道
 "currency: 货币类型
 "startprice: 起拍价
 "endprice": 底价
 "startdate": 起拍时间
 "enddate": 结束信息
 "tradesig":交易签名
 },
 "bids":[
 {
 "bidaddr": 出价者地址
 "nft_token_id": NFT TokenID,
 "nft_contract_addr": NFT合约地址
 "paychan": 交易通道
 "currency": 货币类型
 "price": 出价
 "bidtime": 出价时间
 "tradesig":交易签名
 }
 ]，
 "trans":[
 {
 "nft_contract_addr": NFT合约地址
 "fromaddr": "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169“
 "toaddr": "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
 "nft_token_id": NFT TokenID,
 "transtime": 交易时间,
 "paychan": 交易通道,
 "currency": 货币类型
 "price": 成交价
 "selltype": 交易类型
 "trade_hash":交易Hash
 }
 ]
}
参数说明
name: NFT名称
creator_addr: 创作者地址
creator_portrait: 创作者头像
owner_addr: 拥有者地址
owner_portrait: 拥有者头像
md5: NFT MD5
asset_sample: NFT采样信息（缩略图）
desc: NFT描述
meta: NFT元数据
source_url: NFT源数据链接
nft_contract_addr: NFT合约地址
nft_token_id: NFT TokenID
categories: 分类
collection_creator_addr: 合集创作者地址
collections: 合集名
collection_desc: 合集描述
img: 合集封面
auction: 拍卖信息
bids: 出价信息
trans: 交易记录
verified: 验证状态 "Passed" "NoPass",
selltype: 售卖状态 "NotSale SetPrice FixPrice HighestBid",
mintstate: 铸造状态 "NoMinted Minted",
approve: 授权
royalty: 版税
count: 数量
likes: 关注数
tradesig": 交易签名
auction
selltype": 售卖类型
ownaddr": 拥有者地址
nft_token_id: NFT TokenID
nft_contract_addr: NFT合约地址
paychan: 交易通道
currency: 货币类型
startprice: 起拍价
endprice": 底价
startdate": 起拍时间
enddate": 结束信息
tradesig":交易签名
trans
nft_contract_addr: NFT合约地址
fromaddr: "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169“
toaddr: "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
nft_token_id: NFT TokenID,
transtime: 交易时间,
paychan: 交易通道,
currency": 货币类型
price: 成交价
selltype: 交易类型
trade_hash: 交易Hash
bids
bidaddr: 出价者地址
nft_token_id: NFT TokenID,
nft_contract_addr: NFT合约地址
paychan: 交易通道
currency: 货币类型
price: 出价
bidtime: 出价时间

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/sell  卖出(定价，拍卖)

```json
请求值
user_addr": "字符串，用户地址",
nft_contract_addr": "NFT合约地址",
nft_token_id": "NFT TokenID",
selltype": "售卖方式 NotSale SetPrice FixPrice HighestBid",
pay_channel": "售卖渠道 paypal",
currency_type": "货币类型",
price1": "价格1:如果是定价，那就是售卖价格，如果是拍卖那就是起拍价",
price2": "价格2:如果定价，则不处理，如果是拍卖价那就是底价",
day: 拍卖时长,从当前时间开始持续拍卖多少天,整数
trade_sig": "交易签名",
sig": 数据验证签名,

Request body
Example Value
{
 "user_addr": "0x1234",
 "nft_contract_addr": "0x1234",
 "nft_token_id": "NFT TokenID",
 "selltype": "FixPrice",
 "pay_channel": "eth",
 "currency_type": "eth",
 "price1": "100000",
 "price2": "200000",
 "day": "3",
 "trade_sig": "0x1234",
 "sig": "0x1234"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/cancelSell  取消卖出，下架

```json
user_addr":"字符串，用户地址",
nft_contract_addr":"NFT合约地址",
nft_token_id":"NFT TokenID",
sig": 数据验证签名

Request body
Example Value
{
 "user_addr": "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
 "nft_contract_addr": "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
 "nft_token_id": "1631753648255",
 "sig": "0x123456"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": [
        ""： "",
    ],
    "total_count": int #失败返回0
}
```

### POST  /v2/userRequireKYC  用户申请KYC

```json
请求值
user_addr: 用户地址
country: 国籍
pic: 证件照片
sig: 数据验证签名

Request body
Example Value
{
 "user_addr": "0x1234",
 "country": "cn",
 "pic": "......",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/buy  对NFT进行出价

```json
user_addr":"买方地址",
nft_contract_addr":"NFT 合约地址",
nft_token_id":"NFT token ID",
pay_channel":"支付通道",
currency_type":"金额类型",
price": gwei为单位的整型,
trade_sig": "交易签名",
dead_time: 失效时间,时间戳
sig":数据验证签名,

Request body
Example Value
{
 "user_addr": "0x1234",
 "nft_contract_addr": "0x1234",
 "nft_token_id": "0x1234",
 "pay_channel": "eth",
 "currency_type": "eth",
 "price": "1000000000",
 "trade_sig": "0x1234",
 "dead_time": "123456789",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/cancelBuy  撤销出价

```json
user_addr:"买方地址",
nft_contract_addr:"NFT 合约地址",
nft_token_id:"NFT token ID",
trade_sig: "交易签名",
sig":数据验证签名,

Request body
Example Value
{
 "user_addr": "买方地址",
 "nft_contract_addr": "NFT 合约地址",
 "nft_token_id": "NFT token ID",
 "trade_sig": "交易签名",
 "sig": "签名"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/search  模糊查询 NFT、合集、账户

```json
返回值
{
 "nfts":[], // NFT结果列表 结构参考queryNFT
 "collections":[], // 合集详情列表 结构参考queryCollectionInfo
 "user_addrs":[] // 用户地址列表
}

Request body
Example Value
{
 "match": ""
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryNFTList  获取NFT列表（定价或拍卖）

```json
请求值
match:"模糊查询关键字"
filter: "查询过滤条件",
sort:"查询排序条件 asc,desc",
start_index:"查询起始序号",
count:"查询数量",
Filter 示例
{
 filter:[
 {field:categories, operation:=, value:Art}, 包含美术分类
 {field:categories, operation:=, value:Music}, 包含音乐分类
 {field:categories, operation:=, value:Domain Names}, 
{field:categories, operation:=, value:Virtual Worlds},
 {field:categories, operation:=, value:Trading Cards}, 
{field:categories, operation:=, value:Collectibles}, 
{field:categories, operation:=, value:Sports}, 
{field:categories, operation:=, value:Utility}, 
{field:createdate, operation:>=, value:"12345678"}, 新NFT
 {field:offernum, operation:>, value:"0"}, 有出价的NFT
 {field:price, operation:>=, value:"1"}, 最低价
 {field:price, operation:<=, value:"5"}, 最高价
 {field:collections, operation:=, value:ABC}, 属于哪个合集
 {field:creator_addr, operation:=, value:"0x1234"} 作者地址
 ]
 "sort":[
 {"by":"price","order":"asc"}
 ],
}
最近上架: vrf_time
最近创建： create_time

Request body
Example Value
{
 "match": "",
 "filter": [
 {
 "field": "currency_type",
 "operation": "=",
 "value": "eth"
 }
 ],
 "sort": [
 {
 "by": "price",
 "order": "asc"
 }
 ],
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 nft_token_id:"0x1234",
 nft_contract_id:"abcd",
 vrf_addr:"0xabcd",
 user_addr:"0xabcd",
 ownaddr:"0xabcd",
 sig:"0xabcd",
 source_url:"https://...",
 asset_sample:"",
 md5:"0xabcd",
 desc:"desc",
 currency:"eth",
 price:10,
 Paychan:"eth",
 TransCur:"eth",
 transprice:100,
 createdate:1234567,
 Favorited:10,
 Transcnt:2,
 Transamt:10,
 Verified:"Passed",
 selltype:"NotSale",
 hide:"False",
 sellprice:100,
 offernum:10,
 maxbidprice:1000
 }
 ],
 "total_count":999
}
nft_token_id string #'唯一标识nft标志'"` nft_contract_id string #'合约地址'"`
vrf_addr string #'验证人地址'"` user_addr string #'创建nft地址'"`
ownaddr string #'nft拥有者地址'"` sig string #'签名数据，创建时产生'"`
source_url string #'nfc原始数据保持地址'"` asset_sample string #'缩略图二进制数据'"`
md5 string #'图片md5值'"` desc string #'审核描述：未通过审核描述'"`
currency string #'交易币种'"` price uint64 #'创建时定的价格'"`
Paychan string #'交易通道'"` TransCur string #'交易币种'"`
transprice uint64 #'交易成交价格'"` createdate time.Time #'nft创建时间'"`
Favorited int #'被关注计数'"` Transcnt int #'交易次数，每交易一次加一'"`
Transamt string #'交易总金额'"` Verified string #'nft作品是否通过审核'"`
selltype string #'nft交易类型:售卖方式 NotSale SetPrice FixPrice HighestBid
hide string #'是否让其他人看到'"`
sellprice: 定价：销售价，拍卖：起拍价
offernum: 出价数
maxbidprice: 最高出价

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryNFTCollectionList  获取NFT合集列表

```json
请求值
start_index":"查询起始序号",
count":"查询数量",

Request body
Example Value
{
 "": "",
 "start_index": "0",
 "count": "100"
}
```

```json
返回值
{
 "data":[
 {
 "creator_addr":"0x1234",
 "name":"New Collection",
 "img":"......",
 "contract_addr":"0x12345",
 "desc":'none',
 "royalty":0,
 "categories":"Art",
 "total_count":10
 }
 ],
 "total_count":999
}
creator_addr: 合集创建者地址
name: 合集名
img: 合集封面
contract_addr: 合约地址
desc: 合集描述
categories: 分类
total_count: 合集NFT总数

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryCollectionInfo  获取合集详情

```json
请求值
 - creator_addr: 合集创作者地址,
 - collection_name: 合集名字,

Request body
Example Value
{
 "creator_addr": "",
 "collection_name": "0"
}
```

```json
返回值
{
        "data":[{
                "collection_creator_addr"："0xabcd",
                "nft_contract_addr": "0xabcd",
                "contracttype" : "mment:合约类型"
                "name": "abcd",
                "desc": "abcd",
                "categories":"Art",
                "trade_amount":150000,
                "trade_avg_price":100000,
                "trade_floor_price":10000,
                "extend":""
        }],
        "total_count":1
}

collection_creator_addr: 用户地址
nft_contract_addr: 合约地址
contracttype: 合约类型
name: 合集名称
desc: 合集描述
categories:合集分类
trade_amount:合集内NFT的总交易额
trade_avg_price:合集内NFT成交均价
trade_floor_price:合集内NFT最低价
extend:扩展字段

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": [
        ""： "",
    ],
    "total_count": int #失败返回0
}
```

### POST  /v2/askForApprove  获取交易所授权

```json
请求值
nft_contract_addr: NFT合约地址
nft_token_id: NFT token ID

Request body
Example Value
{
 "nft_contract_addr": "0x1234",
 "nft_token_id": "1234"
}
```

```json
返回值
{
 data{
 "approve":""
 },
 "total_count":1
}
approve: 交易所授权


Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/queryMarketTradingHistory  获取NFT交易所的交易历史

```json
请求值
match:"模糊查询关键字"
filter: "查询过滤条件",
sort:"查询排序条件 asc,desc",
start_index:"查询起始序号",
count:"查询数量",

Request body
Example Value
{
    "match": "",
    "filter": [
    {
        "field": "currency_type",
        "operation": "=".
        "value": "eth"
    }
    ],
    "sort": [
    {
        "by": "price",
        "order": "asc"
    }
    ],
    "start_index": "0",
    "count": 100
}
```

```json
返回值
交易历史列表
{
        "data":[
                {
                        "nft_contract_addr":"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
                        "nft_token_id":"1631753648255",
                        "nft_name":"My NFT",
                        "price":1000000,
                        "count":1,
                        "from":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
                        "to":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
                        "selltype":"FixPrice",
                        "date":123456789,
                        "trade_hash":"0xdfdsfdf"
                }
        ],
        "total_count":999
}

返回交易结构
nft_contract_addr: NFT合约地址
nft_token_id： NFT token ID
nft_name:NFT 名称
price:成交价
count:成交数量
from:卖家地址
to:买家地址
selltype: 交易类型
date:成交时间日期 UTC时间戳
trade_hash: 交易Hash

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### GET  /v2/queryHomePage  获取NFT交易所主页信息

```json
返回值
{
 "code":"200",
 "msg":"success",
 "data":[{
 "announcement":[
 <announcement1>,
 <announcement2>,
 <announcement3>,
 <announcement4>,
 <announcement5>
 ],
 "nft_loop":[<nft>,<nft>,<nft>,<nft>,<nft>],
 "collections":[<collection>,<collection>,<collection>,<collection>,<collection>],
 "nfts":[<nft>,<nft>,<nft>,<nft>,<nft>],
 "total":99
 }],
 "totalcount":1
}
字段说明
announcement: 轮播公告,字符串数组
数据结构如下

{
 "title":"title_string",
 "content":"content_string"
}
nft_loop: 轮播图,nft结构数组，结构参考queryNFT
collections: 热门合集, collections结构数组，结构参考queryNFTCollectionList
nfts: 热门NFT,nft结构数组，结构参考queryNFT
total: Discover总数
```

```json
后端数据结构(前端不管这个)
{
 "announcement":[
 {
 "title":"title_string",
 "content":"content_string"
 },{
 "title":"title_string",
 "content":"content_string"
 },{
 "title":"title_string",
 "content":"content_string"
 },{
 "title":"title_string",
 "content":"content_string"
 },{
 "title":"title_string",
 "content":"content_string"
 }
 ],
 "nft_loop":[{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"}],
 "collections":[{"collection_creator_addr":"0x572bcacb7ae32db658c8dee49e156d455ad59ec8","name":"sai"},{"collection_creator_addr":"0x572bcacb7ae32db658c8dee49e156d455ad59ec8","name":"sai"},{"collection_creator_addr":"0x572bcacb7ae32db658c8dee49e156d455ad59ec8","name":"sai"},{"collection_creator_addr":"0x572bcacb7ae32db658c8dee49e156d455ad59ec8","name":"sai"},{"collection_creator_addr":"0x572bcacb7ae32db658c8dee49e156d455ad59ec8","name":"sai"}],
 "nfts":[{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"},{"contract":"0x53d76f1988B50674089e489B5ad1217AaC08CC85","tokenid":"1335808443097"}]
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

# Query

### GET  /v2/queryPendingKYCList  获取用户KYC列表

```json
返回值
[
 {
 "user_addr":"0x1234",
 "country":"cn",
 "pic":"......",
 }
]
user_addr: 待KYC用户地址
pic:证件照片
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### GET  /v2/queryPendingVrfList  获取NFT待审核列表

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### GET  /v2/version  版本号

```json
返回值
    版本号字符串

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

# backend

### POST  /v2/userKYC  用户KYC

```json
请求值
vrf_addr: 验证者地址
user_addr: KYC用户地址
desc": KYC结果原因描述
kyc_res: KYC结果["Passed","NoPass"]
sig": 数据验证签名

Request body
Example Value
{
 "vrf_addr": "0x1234",
 "user_addr": "0x1234",
 "desc": "none",
 "kyc_res": "Passed",
 "sig": "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
}
```

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": [
        ""： "",
    ],
    "total_count": int #失败返回0
}
```

### POST  /v2/vrfNFT  审核NFT

```json
Request body
Example Value
{
    ”vrf_addr“: "string",
    "Owner": "string",
    "nft_contract_addr": "string",
    "nft_token_id": "string",
    "desc": "string",
    "vrf_res": "string",
    "sig": "string"
}

Responses
200    Successful Response
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": [
        ""： "",
    ],
    "total_count": int #失败返回0
}
```

### GET  /v2/get_sys_para  获取系统信息

```json
Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/set_sys_para  设置系统信息

```json
Request body
Example Value
    "string"

Responses
200    Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}

422    Validation Error
Example 
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```

### POST  /v2/buyResultInterface  管理员数据导入

```json
输入参数
admin_addr: 管理员地址
from: 交易from
to: 交易to
nft_contract_addr: 合约地址
nft_token_id: TokenID
trade_sig: 交易签名
price:交易价格
royalty: 交易版税
txhash: 交易Hash
sig: 数据验证签名
admin_sig: 管理员数据验证签名

Request body
Example Value
{
 "admin_addr": "",
 "from": "0xabcd",
 "to": "0xabcd",
 "nft_contract_addr": "0xabcd",
 "nft_token_id": "0xabcd",
 "trade_sig": "0xabcd",
 "price": "9999",
 "royalty": "3",
 "txhash": "0xabcd",
 "sig": "0xabcd",
 "admin_sig": "0xabcd"
}
```

```json
Responses
200  Successful Response
Example Value
{
    "code": 200,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}


422  Validation Error
Example Value
{
    "code": 422,
    "msg": "string",
    "data": {
        ""： "",
    },
    "total_count": int #失败返回0
}
```
