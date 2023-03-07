package contracts

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/nftexchange/nftserver/common/contracts/nft1155"
	"github.com/nftexchange/nftserver/common/contracts/trade"
	"github.com/nftexchange/nftserver/common/contracts/weth9"
	"sync"

	//"github.com/nftexchange/nftserver/ethhelper"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	TransferSingleHash                 = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	RoyaltyHash                        = "0x611d12c0f8b2d9f4cfb23a30f560228db53e712dfcd34bb5b239e702efc2d22f"
	PricingHash                        = "0xb14cef5ea4cbf14922e663d47af6b5458327dda3f30ba4ffd69e1d63a305ed2c"
	BidingHash                         = "0xf64cc560cd0c25fe33108f5a79cd3104e22f65e64b260be4b2072104012fa60a"
	Erc1155Interface                   = "0xd9b67a26"
	WormHolesVerseion                  = "0.01"
	WormHolesMint                      = 0
	WormHolesTransfer                  = 1
	WormHolesExchange                  = 6
	WormHolesPledge                    = 7
	WormHolesUnPledge                  = 8
	WormHolesOpenExchanger             = 11
	WormHolesExToBuyTransfer           = 14
	WormHolesBuyFromSellTransfer       = 15
	WormHolesBuyFromSellMintTransfer   = 16
	WormHolesExToBuyMintToSellTransfer = 17
	WormHolesExAuthToExBuyTransfer     = 18
	WormHolesExAuthToExMintBuyTransfer = 19
	WormHolesExSellNoAuthTransfer      = 20
	WormHolesExSellBatchAuthTransfer   = 27
	WormHolesExForceBuyingAuthTransfer = 28

	WormHolesContract   = "0xffffffffffffffffffffffffffffffffffffffff"
	WormHolesNftCount   = "1"
	UserMintDeepDef     = "0x0000000000000000000000000000000000000001"
	TransTypeBuyerTrans = 14
	TransTypeMintTrans  = 16
	WormholesVersion    = "0"
	GasLimitTx1819      = 200000
)

type NftTx struct {
	Operator         string
	From             string
	To               string
	Contract         string
	TokenId          string
	Value            string
	Price            string
	Ratio            string
	TxHash           string
	Ts               string
	BlockNumber      string
	TransactionIndex string
	MetaUrl          string
	NftAddr          string
	Nonce            string
	Status           bool
	TransType        int
}

type NftTrans struct {
	Nfttxs     []*NftTx
	Minttxs    []*NftTx
	Wnfttxs    []*NftTx
	Wmintxs    []*NftTx
	Wethc      map[string]bool
	Wexchanger string
}

type WethTransfer struct {
	From  string
	To    string
	Value string
}

type WethChange map[string]bool

type Royalty struct {
	contract string
	Id       string
	Ratio    string
	Receiver string
}

type Mint struct {
	Operator string
	From     string
	To       string
	contract string
	Id       string
	Value    string
}

type Wormholes struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"`
}

type NftMeta struct {
	Meta    string `json:"meta"`
	TokenId string `json:"token_id"`
}

type WormholesMint struct {
	Version   string `json:"version"`
	Type      uint8  `json:"type"`
	Royalty   uint32 `json:"royalty"`
	MetaUrl   string `json:"meta_url"`
	Exchanger string `json:"exchanger"`
}

type WormholesOpenExchanger struct {
	Version     string `json:"version"`
	Type        uint8  `json:"type"`
	Feerate     uint32 `json:"fee_rate"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Blocknumber string `json:"block_number"`
}

type Buyer struct {
	Price       string `json:"price"`
	Nftaddress  string `json:"nft_address"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Seller      string `json:"seller"`
	Sig         string `json:"sig"`
}

type Buyer1 struct {
	Price       string `json:"price"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

//type Buyer2 struct {
//	Nftaddress  string `json:"nft_address"`
//	Exchanger   string `json:"exchanger"`
//	Blocknumber string `json:"block_number"`
//	Sig         string `json:"sig"`
//}

type WormholesFixTrans struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"` //14
	Buyer   `json:"buyer"`
}

type WormholesFixTransAuth struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyer         `json:"buyer"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type BuyerTrans struct {
	Worm WormholesFixTrans `json:"wormholes"`
}

type ExchangerTrans struct {
	Worm WormholesFixTrans `json:"wormholes"`
}

type WormholesAuthFixTrans struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyer         `json:"buyer"`
	Seller1       `json:"seller1"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerAuthTrans struct {
	Worm WormholesAuthFixTrans `json:"wormholes"`
}

type Buyauth struct {
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type Sellerauth struct {
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type WormholesBatchAuthFixTrans struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyauth       `json:"buyer_auth"`
	Buyer         `json:"buyer"`
	Sellerauth    `json:"seller_auth"`
	Seller1       `json:"seller1"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerBatchAuthTrans struct {
	Worm WormholesBatchAuthFixTrans `json:"wormholes"`
}

type WormholesForceBuyingTrans struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyauth       `json:"buyer_auth"`
	Buyer         Buyer         `json:"buyer"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerForceBuyingAuthTrans struct {
	Worm WormholesForceBuyingTrans `json:"wormholes"`
}

type Seller1 struct {
	Price       string `json:"price"`
	Nftaddress  string `json:"nft_address"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type WormholesBuyFromSellTrans struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"` //15
	Seller1 `json:"seller"`
}

type Seller2 struct {
	Price         string `json:"price"`
	Royalty       string `json:"royalty"`
	Metaurl       string `json:"meta_url"`
	Exclusiveflag string `json:"exclusive_flag"`
	Exchanger     string `json:"exchanger"`
	Blocknumber   string `json:"block_number"`
	Sig           string `json:"sig"`
}

type WormholesBuyFromSellMintTrans struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"`
	Seller2 `json:"seller2"`
}

type Seller struct {
	Price       string `json:"price"`
	Royalty     string `json:"royalty"`
	Metaurl     string `json:"meta_url"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type WormholesSellerMintTrans struct {
	Version string  `json:"version"`
	Type    uint8   `json:"type"`
	Seller  Seller2 `json:"seller2"`
	Buyer   Buyer1  `json:"buyer"`
}

type MintTrans struct {
	Worm WormholesSellerMintTrans `json:"wormholes"`
}

type ExchangerAuth struct {
	Exchangerowner string `json:"exchanger_owner"`
	To             string `json:"to"`
	Blocknumber    string `json:"block_number"`
	Sig            string `json:"sig"`
}

type WormholesAuthMintTrans struct {
	Version       string        `json:"version"`
	Type          uint8         `json:"type"`
	Seller        Seller2       `json:"seller2"`
	Buyer         Buyer1        `json:"buyer"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type AuthMintTrans struct {
	Worm WormholesAuthMintTrans `json:"wormholes"`
}

type ExchangerMintTrans struct {
	Version string  `json:"version"`
	Type    uint8   `json:"type"`
	Seller  Seller2 `json:"seller2"`
	Buyer   Buyer1  `json:"buyer"`
}

type ExchangerAuthMintTrans struct {
	Version       string        `json:"version"`
	Type          uint8         `json:"type"`
	Seller        Seller2       `json:"seller2"`
	Buyer         Buyer1        `json:"buyer"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type WormholesTransfer struct {
	Version    string `json:"version"`
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
}

type WormholesExchange struct {
	Version    string `json:"version"`
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
}

type WormholesPledge struct {
	Version    string `json:"version"`
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
}

/*type AccountNFT struct {
	//Account
	Name                  string
	Symbol                string
	Price                 *big.Int
	Direction             uint8 // 0:未交易,1:买入,2:卖出
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel uint8

	Creator   common.Address
	Royalty   uint32
	Exchanger common.Address
	MetaURL   string
}
*/

type AccountNFT struct {
	//Account
	Name   string
	Symbol string
	//Price                 *big.Int
	//Direction             uint8 // 0:un_tx,1:buy,2:sell
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel            uint8
	MergeNumber           uint32
	PledgedFlag           bool
	NFTPledgedBlockNumber *big.Int

	Creator   common.Address
	Royalty   uint32
	Exchanger common.Address
	MetaURL   string
}

type Account struct {
	Nonce   uint64
	Balance *big.Int
	// *** modify to support nft transaction 20211220 begin ***
	//NFTCount uint64		// number of nft who account have
	// *** modify to support nft transaction 20211220 end ***
	Root           common.Hash // merkle root of the storage trie
	CodeHash       []byte
	PledgedBalance *big.Int
	// *** modify to support nft transaction 20211215 ***
	//Owner common.Address
	// whether the account has a NFT exchanger
	ExchangerFlag bool
	BlockNumber   *big.Int
	// The ratio that exchanger get.
	FeeRate       uint32
	ExchangerName string
	ExchangerURL  string
	// ApproveAddress have the right to handle all nfts of the account
	ApproveAddressList []common.Address
	// NFTBalance is the nft number that the account have
	NFTBalance            uint64
	NFTPledgedBlockNumber uint64
	AccountNFT
}

type TransLock struct {
	Mux         sync.Mutex
	blocknumber uint64
	nonce       uint64
	initTime    time.Time
}

func (t *TransLock) Init(c *ethclient.Client, blocknumber uint64, address common.Address) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	nonce, err := c.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Println("Init() err=", err)
		return
	}
	t.nonce = nonce
	t.blocknumber = blocknumber + 5
	t.initTime = time.Now().Add(5 * time.Second)
}

func (t *TransLock) GetNonce(c *ethclient.Client, blocknumber uint64, address common.Address) (uint64, error) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	if t.blocknumber <= blocknumber || t.initTime.Before(time.Now()) {
		nonce, err := c.PendingNonceAt(context.Background(), address)
		if err != nil {
			log.Println("Init() err=", err)
			return t.nonce, err
		}
		fmt.Println("GetNonce() get from chain  t.nonce=", t.nonce, " chain nonce=", nonce)
		if nonce > t.nonce {
			t.nonce = nonce
		}
		t.blocknumber = blocknumber + 1
		t.initTime = time.Now().Add(5 * time.Second)
	} else {
		t.nonce = t.nonce + 1
	}
	fmt.Println("GetNonce() nonce=", t.nonce)
	return t.nonce, nil
}

var transLock TransLock

func init() {

}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func GetAccountInfo(nftaddr common.Address, blockNumber *big.Int) (*Account, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetAccountInfo() err=", err)
		return nil, err
	}
	var result Account
	err = client.CallContext(context.Background(), &result, "eth_getAccountInfo", nftaddr, toBlockNumArg(blockNumber))
	if err != nil {
		log.Println("GetAccountInfo() err=", err)
		return nil, err
	}
	return &result, err
}

func GetForcedSaleAmount(nftaddr common.Address) (string, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetForcedSaleAmount() err=", err)
		return "", err
	}
	var result hexutil.Big
	err = client.CallContext(context.Background(), &result, "eth_getForcedSaleAmount", nftaddr)
	if err != nil {
		log.Println("GetForcedSaleAmount() err=", err)
		return "", err
	}
	return result.String(), err
}

func GetLatestAccountInfo(nftaddr common.Address) (*Account, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetLatestAccountInfo() err=", err)
		return nil, err
	}
	defer client.Close()
	var result Account
	err = client.CallContext(context.Background(), &result, "eth_getAccountInfo", nftaddr, "latest")
	if err != nil {
		log.Println("GetLatestAccountInfo() err=", err)
		return nil, err
	}
	return &result, err
}

type BeneficiaryAddress struct {
	Address    common.Address
	NftAddress common.Address
}
type BeneficiaryAddressList []*BeneficiaryAddress

func GetSnftAddressList(blockNumber *big.Int, fulltx bool) ([]*BeneficiaryAddress, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		fmt.Println("GetSnftAddressList() err=", err)
		return nil, err
	}
	var result BeneficiaryAddressList
	err = client.CallContext(context.Background(), &result, "eth_getBlockBeneficiaryAddressByNumber", toBlockNumArg(blockNumber), fulltx)
	if err != nil {
		fmt.Println("GetSnftAddressList() err=", err)
		return nil, err
	}
	return result, err
}

func SendTrans(to, price string, prv *ecdsa.PrivateKey) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	//privateKey, err := crypto.HexToECDSA(prv)
	//if err != nil {
	//	log.Println("SendTrans() err=", err)
	//	return err
	//}
	publicKey := prv.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("SendTrans() err=", err)
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	log.Println("SendTrans() nonce=", nonce)
	p, _ := strconv.ParseUint(price, 10, 64)
	value := big.NewInt(int64(p)) // in wei (1 eth)
	gasLimit := uint64(21000)     // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), prv)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	return nil
}

func Deposit(to, price string, prv *ecdsa.PrivateKey) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	//privateKey, err := crypto.HexToECDSA(prv)
	//if err != nil {
	//	log.Println("SendTrans() err=", err)
	//	return err
	//}
	instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	if err != nil {
		log.Println(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(to))
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return err
	}
	log.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("Buy() SuggestGasPrice() err=", err)
		return err
	}
	auth := bind.NewKeyedTransactor(prv)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	n, err := instance.Deposit(auth)
	if err != nil {
		log.Println("Buy() SuggestGasPrice() err=", err)
		return err
	}
	log.Println("tx.hash=", n.Hash())
	return nil
}

func Approve(to, price string, prv *ecdsa.PrivateKey) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return err
	}
	instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	if err != nil {
		log.Println(err)
	}
	fromAddr := crypto.PubkeyToAddress(prv.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return err
	}
	log.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("Buy() SuggestGasPrice() err=", err)
		return err
	}
	auth := bind.NewKeyedTransactor(prv)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	//auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	n, err := instance.Approve(auth, common.HexToAddress(TradeAddr), big.NewInt(int64(p)))
	if err != nil {
		log.Println("Buy() SuggestGasPrice() err=", err)
		return err
	}
	log.Println("tx.hash=", n.Hash())
	return nil
}

func SetApprove(contract string, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("SetApprove() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(prv)
	//if err != nil {
	//	log.Println("SetApprove() err=", err)
	//	return nil, err
	//}
	publicKey := prv.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("SetApprove() addr=", fromAddress)
	instance, err := nft1155.NewNft1155(common.HexToAddress(Nft1155Addr), client)
	if err != nil {
		log.Println("SetApprove() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("SetApprove() err=", err)
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("SetApprove() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(prv)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	t, err := instance.SetApprovalForAll(auth, common.HexToAddress(contract), true)
	if err != nil {
		log.Println("SetApprove() err=", err)
		return nil, err
	}
	log.Println("tx.hash=", t.Hash())
	return t, nil
}

//func GetCurrentBlockNumber() uint64 {
//	var client *ethclient.Client
//	var err error
//	for {
//		for {
//			fmt.Println("GetCurrentBlockNumber() dial ", "EthNode=", EthNode)
//			client, err = ethclient.Dial(EthNode)
//			if err != nil {
//				log.Println("GetCurrentBlockNumber()", "EthNode=", EthNode, " connect err=", err)
//				time.Sleep(ReDialDelyTime * time.Second)
//			} else {
//				//log.Println("GetCurrentBlockNumber() connect OK!")
//				break
//			}
//		}
//		fmt.Println("GetCurrentBlockNumber() get HeaderByNumber")
//		header, err := client.HeaderByNumber(context.Background(), nil)
//		if err != nil {
//			log.Println("GetCurrentBlockNumber() get HeaderByNumber err=", err)
//			client.Close()
//			time.Sleep(ReDialDelyTime * time.Second)
//		} else {
//			log.Println("GetCurrentBlockNumber() header.Number=", header.Number.String())
//			client.Close()
//			return header.Number.Uint64()
//		}
//	}
//}

func GetCurrentBlockNumber() uint64 {
	var client *ethclient.Client
	var err error
	for {
		fmt.Println("GetCurrentBlockNumber() dial ", "EthNode=", EthNode)
		client, err = ethclient.Dial(EthNode)
		if err != nil {
			time.Sleep(ReDialDelyTime * time.Second)
			continue
		}
		blocknum, err := client.BlockNumber(context.Background())
		client.Close()
		if err != nil {
			log.Println("GetCurrentBlockNumber() err=", err)
			time.Sleep(ReDialDelyTime * time.Second)
			continue
		} else {
			fmt.Println("GetCurrentBlockNumber() blocknumber=", blocknum)
			return blocknum
		}
	}
}

func GetNonce(contract, userAddr, tokenId string) (*big.Int, error) {
	//return big.NewInt(1), nil
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("err=", err)
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(TradeAuthAddrPrv)
	if err != nil {
		log.Println("err=", err)
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	if err != nil {
		log.Println("err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	opts := &bind.CallOpts{Context: context.Background()}
	//opts := &bind.CallOpts{From: common.HexToAddress("0xBAaeeab54cDFF708a8dCc51F56f4e2A4CE7c2ABc"), Context: context.Background()}
	n, err := instance.Nonce(opts, common.HexToAddress(contract), big.NewInt(int64(tokenid)), common.HexToAddress(userAddr))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return n, nil
}

func BalanceOfWeth(userAddr string) (uint64, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("BalanceOfWeth() Dial err=", err)
		return 0, err
	}
	instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	if err != nil {
		log.Println("BalanceOfWeth() NewWeth9 err=", err)
		return 0, err
	}
	opts := &bind.CallOpts{Context: context.Background()}
	n, err := instance.BalanceOf(opts, common.HexToAddress(userAddr))
	if err != nil {
		log.Println("BalanceOfWeth() BalanceOf err=", err)
		return 0, err
	}
	return n.Uint64(), nil
}

func WormsBalance(userAddr string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("BalanceOfWeth() Dial err=", err)
		return "", err
	}
	fromAddress := common.HexToAddress(userAddr)
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Println("SendTrans() err=", err)
		return "", err
	}
	return balance.String(), nil
}

func AllowanceOfWeth(userAddr string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("AllowanceOfWeth() Dial err=", err)
		return "", err
	}
	instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	if err != nil {
		log.Println("AllowanceOfWeth() NewWeth9 err=", err)
		return "", err
	}
	opts := &bind.CallOpts{Context: context.Background()}
	n, err := instance.Allowance(opts, common.HexToAddress(userAddr), common.HexToAddress(TradeAddr))
	if err != nil {
		log.Println("AllowanceOfWeth() Allowance err=", err)
		return "", err
	}
	return n.String(), nil
}

func IsErc721(userAddr string) (int, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("IsErc721() Dial err=", err)
		return 0, err
	}
	instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	if err != nil {
		log.Println("IsErc721() NewWeth9 err=", err)
		return 0, err
	}
	opts := &bind.CallOpts{Context: context.Background()}
	n, err := instance.Allowance(opts, common.HexToAddress(userAddr), common.HexToAddress(TradeAddr))
	if err != nil {
		log.Println("IsErc721() Allowance err=", err)
		return 0, err
	}
	return int(n.Int64()), nil
}

func IsOwnerOfNFT1155(owner, contract, tokenId string) (bool, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("IsOwnerOfNFT1155() Dial err=", err)
		return false, err
	}
	instance, err := nft1155.NewNft1155(common.HexToAddress(Nft1155Addr), client)
	if err != nil {
		log.Println("IsOwnerOfNFT1155() NewNft1155 err=", err)
		return false, err
	}
	id, _ := strconv.ParseUint(tokenId, 10, 64)
	opts := &bind.CallOpts{From: common.HexToAddress(owner), Context: context.Background()}
	n, err := instance.BalanceOf(opts, common.HexToAddress(owner), big.NewInt(int64(id)))
	if err != nil {
		log.Println("IsOwnerOfNFT1155() Allowance err=", err)
		return false, err
	}
	return n.Uint64() >= 1, nil
}

func IsApprovedNFT1155(owner, contract string) (bool, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("IsApprovedNFT1155() Dial err=", err)
		return false, err
	}
	instance, err := nft1155.NewNft1155(common.HexToAddress(Nft1155Addr), client)
	if err != nil {
		log.Println("IsApprovedNFT1155() NewNft1155 err=", err)
		return false, err
	}
	opts := &bind.CallOpts{From: common.HexToAddress(owner), Context: context.Background()}
	b, err := instance.IsApprovedForAll(opts, common.HexToAddress(contract), common.HexToAddress(TradeAddr))
	if err != nil {
		log.Println("IsApprovedNFT1155() Allowance err=", err)
		return false, err
	}
	return b, nil
}

func IsErcNFT1155(contract string) (bool, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("IsApprovedNFT1155() Dial err=", err)
		return false, err
	}
	instance, err := nft1155.NewNft1155(common.HexToAddress(Nft1155Addr), client)
	if err != nil {
		log.Println("IsApprovedNFT1155() NewNft1155 err=", err)
		return false, err
	}
	opts := &bind.CallOpts{Context: context.Background()}
	var iId [4]byte
	m, _ := hexutil.Decode(Erc1155Interface)
	copy(iId[:], m)
	b, err := instance.SupportsInterface(opts, iId)
	if err != nil {
		log.Println("IsApprovedNFT1155() Allowance err=", err)
		return false, err
	}
	return b, nil
}

/*func OwnAndAprove(owner, contract, tokenId string) (bool, error) {
	b, err := IsErcNFT1155(contract)
	if err != nil {
		log.Println("OwnAndAprove() err=", err)
		return false, err
	}
	if !b {
		isOwner, err := ethhelper.IsOwnerOfNFT721(owner, contract, tokenId)
		if err != nil {
			log.Println("OwnAndAprove() NFT721 err=", err)
			return false, err
		}
		approve, err := ethhelper.IsApprovedNFT721(owner, contract, tokenId)
		if err != nil {
			log.Println("OwnAndAprove() NFT721 err=", err)
			return false, err
		}
		return isOwner && approve, nil
	} else {
		isOwner, err := IsOwnerOfNFT1155(owner, contract, tokenId)
		if err != nil {
			log.Println("OwnAndAprove() NFT1155 err=", err)
			return false, err
		}
		approve, err := IsApprovedNFT1155(owner, contract)
		if err != nil {
			log.Println("OwnAndAprove() NFT1155 err=", err)
			return false, err
		}
		return isOwner && approve, nil
	}
}*/

func Sign(data []byte, prv *ecdsa.PrivateKey) (string, error) {
	//hash := crypto.Keccak256(data)
	//log.Printf("hash=%s\n", hash)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", crypto.Keccak256(data))
	dmsg := []byte(msg)
	log.Println(dmsg)
	sig, err := crypto.Sign(crypto.Keccak256([]byte(msg)), prv)
	if err != nil {
		log.Println("signature error: ", err)
		return "", err
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	return sigstr, err
}

func PricePack(priceT byte, contract string, tokenId string, count string, price string, nonce string) []byte {
	var data []byte
	data = append(data, priceT)
	ctract, _ := hexutil.Decode(contract)
	data = append(data, ctract...)
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	data = append(data, common.LeftPadBytes(big.NewInt(int64(tokenid)).Bytes(), 32)...)
	c, _ := strconv.ParseUint(count, 10, 64)
	amount := common.LeftPadBytes(big.NewInt(int64(c)).Bytes(), 32)
	data = append(data, amount...)
	p, _ := strconv.ParseUint(price, 10, 64)
	data = append(data, common.LeftPadBytes(big.NewInt(int64(p)).Bytes(), 32)...)
	n, _ := strconv.ParseInt(nonce, 10, 64)
	data = append(data, common.LeftPadBytes(big.NewInt(n).Bytes(), 32)...)
	return data
}

func MintPack(contract string, toAddr string, tokenId string, count string, royalty string, tokenUri string) []byte {
	var data []byte
	ctract, _ := hexutil.Decode(contract)
	data = append(data, ctract...)
	toaddr, _ := hexutil.Decode(toAddr)
	data = append(data, toaddr...)
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	data = append(data, common.LeftPadBytes(big.NewInt(int64(tokenid)).Bytes(), 32)...)
	c, _ := strconv.ParseUint(count, 10, 64)
	amount := common.LeftPadBytes(big.NewInt(int64(c)).Bytes(), 32)
	data = append(data, amount...)
	r, _ := strconv.ParseUint(royalty, 10, 64)
	data = append(data, common.LeftPadBytes(big.NewInt(int64(r)).Bytes(), 2)...)
	data = append(data, []byte(tokenUri)...)
	return data
}

func TradeSign(tradeType byte, contract string, tokenId string, count string, price string, nonce string, key *ecdsa.PrivateKey) (string, error) {
	//key, err := crypto.HexToECDSA(prv)
	//if err != nil {
	//	log.Printf("TradeSign() key err=", err)
	//	return "", err
	//}
	data := PricePack(tradeType, contract, tokenId, count, price, nonce)
	sign, err := Sign(data, key)
	if err != nil {
		log.Println("TradeSign() sign err=", err)
		return "", err
	}
	return sign, nil
}

func MintSign(contract string, toAddr string, tokenId string, count string, royalty string, tokenUri string, prv string) (string, error) {
	key, err := crypto.HexToECDSA(prv)
	if err != nil {
		log.Println("MintSign() key err=", err)
		return "", err
	}
	data := MintPack(contract, toAddr, tokenId, count, royalty, tokenUri)
	sign, err := Sign(data, key)
	if err != nil {
		log.Println("MintSign() sign err=", err)
		return "", err
	}
	return sign, nil
}

func BuyMint(contract, from, to, tokenId, count, royalty, price string, mintSig, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	log.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		log.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	if err != nil {
		log.Println(err)
		log.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	//if oldNonce == 0 {
	//	oldNonce = nonce
	//} else {
	//	for {
	//		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//		if err != nil {
	//			log.Println(err)
	//			log.Println("Buy() PendingNonceAt() err=", err)
	//			return nil, err
	//		}
	//		if oldNonce != nonce {
	//			oldNonce = nonce
	//			break
	//		}
	//	}
	//}
	log.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromKey)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	r, _ := strconv.ParseUint(royalty, 10, 64)
	c, _ := strconv.ParseUint(count, 10, 64)

	//Pricing1155Mint(opts *bind.TransactOpts, _addr common.Address, _from common.Address,
	//_to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16,
	//_tokenURI string, _minerSig []byte, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	n, err := instance.Pricing1155Mint(auth, common.HexToAddress(contract),
		common.HexToAddress(from), common.HexToAddress(to), big.NewInt(int64(tokenid)), big.NewInt(int64(c)),
		uint16(r), "", mintSig, tradeSig, []byte{})
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func Buy(contract, from, to, tokenId, count, price string, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	log.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		log.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	if err != nil {
		log.Println(err)
		log.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	log.Println("Buy() nonce=", nonce)
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	log.Println("Buy() balance=", balance.String())
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromKey)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	c, _ := strconv.ParseUint(count, 10, 64)
	auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)

	n, err := instance.Pricing1155(auth, common.HexToAddress(contract), common.HexToAddress(from),
		common.HexToAddress(to), big.NewInt(int64(tokenid)), big.NewInt(int64(c)), tradeSig, []byte{})
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func BuyBidding(contract, from, to, tokenId, count, price string, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	log.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		log.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	if err != nil {
		log.Println(err)
		log.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	log.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromKey)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	c, _ := strconv.ParseUint(count, 10, 64)
	//auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	//Biding1155(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address,
	//_id *big.Int, _amount *big.Int, _price *big.Int, _toSig []byte, _data []byte)
	n, err := instance.Biding1155(auth, common.HexToAddress(contract), common.HexToAddress(from),
		common.HexToAddress(to), big.NewInt(int64(tokenid)), big.NewInt(int64(c)), big.NewInt(int64(p)), tradeSig, []byte{})
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func BuyBidingMint(contract, from, to, tokenId, count, royalty, price string, mintSig, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	log.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		log.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	if err != nil {
		log.Println(err)
		log.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		log.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	//if oldNonce == 0 {
	//	oldNonce = nonce
	//} else {
	//	for {
	//		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//		if err != nil {
	//			log.Println(err)
	//			log.Println("Buy() PendingNonceAt() err=", err)
	//			return nil, err
	//		}
	//		if oldNonce != nonce {
	//			oldNonce = nonce
	//			break
	//		}
	//	}
	//}
	log.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromKey)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	r, _ := strconv.ParseUint(royalty, 10, 64)
	c, _ := strconv.ParseUint(count, 10, 64)

	//Biding1155Mint(opts *bind.TransactOpts, _addr common.Address, _from common.Address,
	//_to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int,
	//_royaltyRatio uint16, _tokenURI string, _minerSig []byte, _toSig []byte, _data []byte) (*types.Transaction, error)
	n, err := instance.Biding1155Mint(auth, common.HexToAddress(contract),
		common.HexToAddress(from), common.HexToAddress(to), big.NewInt(int64(tokenid)),
		big.NewInt(int64(c)), big.NewInt(int64(p)), uint16(r), "", mintSig, tradeSig, []byte{})
	if err != nil {
		log.Println(err)
		log.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func GetBlockTxs(blockNum uint64) ([]*NftTx, map[string]bool, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("GetBlockTxs() err=", err)
		return nil, nil, err
	}
	tradeInstance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	if err != nil {
		log.Println("GetBlockTxs() NewTrade() err=", err)
		return nil, nil, err
	}
	nft1155Instance, err := nft1155.NewNft1155(common.HexToAddress(Nft1155Addr), client)
	if err != nil {
		log.Println("GetBlockTxs() NewNft1155() err=", err)
		return nil, nil, err
	}
	//weth9Instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	//if err != nil {
	//	log.Println("GetBlockTxs() NewWeth9() err=", err)
	//	return nil, nil, err
	//}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
	if err != nil {
		log.Println("GetBlockTxs() BlockByNumber() err=", err)
		return nil, nil, err
	}
	transT := block.Time()
	log.Println(time.Unix(int64(transT), 0))
	nfttxs := make([]*NftTx, 0, 20)
	wethchanges := make(WethChange)
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		if tx.To().Hex() == TradeAddr {
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				log.Println("GetBlockTxs() TransactionReceipt() err=", err)
				return nil, nil, err
			}
			for _, logger := range receipt.Logs {
				if logger.Address == common.HexToAddress(Nft1155Addr) {
					if logger.Topics[0] == common.HexToHash(TransferSingleHash) {
						nftx := NftTx{}
						m, err := nft1155Instance.ParseTransferSingle(*logger)
						if err != nil {
							log.Println("GetBlockTxs() ParseTransferSingle() err=", err)
							return nil, nil, err
						}
						nftx.Operator = m.Operator.String()
						nftx.From = m.From.String()
						nftx.To = m.To.String()
						nftx.Value = m.Value.String()
						nftx.Contract = Nft1155Addr
						nftx.TokenId = m.Id.String()
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
						nfttxs = append(nfttxs, &nftx)
					} else if logger.Topics[0] == common.HexToHash(RoyaltyHash) {
						r, err := nft1155Instance.ParseRoyalty(*logger)
						if err != nil {
							log.Println("GetBlockTxs() ParseRoyalty() err=", err)
							return nil, nil, err
						}
						nftx := NftTx{}
						nftx.To = r.Receiver.String()
						nftx.Contract = Nft1155Addr
						nftx.TokenId = r.Id.String()
						nftx.Ratio = strconv.Itoa(int(r.Ratio))
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
						nfttxs = append(nfttxs, &nftx)
					}
				} else if logger.Address == common.HexToAddress(TradeAddr) {
					nftx := NftTx{}
					if logger.Topics[0] == common.HexToHash(PricingHash) {
						p, err := tradeInstance.ParsePRICING(*logger)
						if err != nil {
							log.Println("GetBlockTxs() ParsePRICING() err=", err)
							return nil, nil, err
						}
						nftx.From = p.From.String()
						nftx.To = p.To.String()
						nftx.Contract = p.Addr.String()
						nftx.TokenId = p.Id.String()
						nftx.Value = p.Amount.String()
						nftx.Price = p.Price.String()
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)

					} else if logger.Topics[0] == common.HexToHash(BidingHash) {
						p, err := tradeInstance.ParseBIDING(*logger)
						if err != nil {
							log.Println("GetBlockTxs() TransactionReceipt() err=", err)
							return nil, nil, err
						}
						nftx.From = p.From.String()
						nftx.To = p.To.String()
						nftx.Contract = p.Addr.String()
						nftx.TokenId = p.Id.String()
						nftx.Value = p.Amount.String()
						nftx.Price = p.Price.String()
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
					} else {
						continue
					}
					nfttxs = append(nfttxs, &nftx)
					wethchanges[nftx.To] = true
				}
			}
		}
		/*if tx.To().Hex() == Nft1155Addr {
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				log.Println("GetBlockTxs() TransactionReceipt() err=", err)
				return nil, nil, err
			}
			for _, log := range receipt.Logs {
				t, err := nft1155Instance.ParseRoyalty(*log)
				if err != nil {
					log.Println("GetBlockTxs() TransactionReceipt() err=", err)
					continue
				}
				log.Println(t)
				ts, err := nft1155Instance.ParseTransferSingle(*log)
				if err != nil {
					log.Println("GetBlockTxs() ParseTransferSingle() err=", err)
					continue
				}
				log.Println(ts)
			}
		}
		if tx.To().Hex() == Weth9Addr {
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				log.Println("GetBlockTxs() TransactionReceipt() err=", err)
				return nil, nil, err
			}
			for _, log := range receipt.Logs {
				t, err := weth9Instance.ParseWithdrawal(*log)
				if err != nil {
					log.Println("GetBlockTxs() TransactionReceipt() err=", err)
					continue
				}
				log.Println(t)
			}
		}*/
	}
	return nfttxs, wethchanges, nil
}

func hashMsg(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func recoverAddress(msg string, sigStr string) (*common.Address, error) {
	sigData, err := hexutil.Decode(sigStr)
	if err != nil {
		log.Println("recoverAddress() err=", err)
		return nil, err
	}
	if len(sigData) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := hashMsg([]byte(msg))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(*rpk)
	return &addr, nil
}

func RecoverAddress(msg string, sigStr string) (*common.Address, error) {
	sigData, err := hexutil.Decode(sigStr)
	if err != nil {
		log.Println("recoverAddress() err=", err)
		return nil, err
	}
	if len(sigData) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := hashMsg([]byte(msg))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(*rpk)
	return &addr, nil
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func WormholesSign(msg string, prv *ecdsa.PrivateKey) (string, error) {
	sig, err := crypto.Sign(signHash([]byte(msg)), prv)
	if err != nil {
		fmt.Println("EthSign() err=", err)
		return "", err
	}
	sig[64] += 27
	return hexutil.Encode(sig), nil
}

func GenNftAddr(UserMintDeep *big.Int) error {
	//UserMintDeep = "0x" + UserMintDeep
	//nft, err := hexutil.DecodeBig(*UserMintDeep)
	//if err != nil {
	//	log.Println("GenNftAddr() err=", err)
	//	return err
	//}

	UserMintDeep = UserMintDeep.Add(UserMintDeep, big.NewInt(1))
	//*UserMintDeep = hexutil.EncodeBig(nft)
	//*UserMintDeep =  common.BytesToAddress(nft.Bytes()).String()
	//if len(UserMintDeep) >= 2 {
	//	UserMintDeep = UserMintDeep[2:]
	//}
	return nil
}

func GetBlockTxsNew(blockNum uint64) (*NftTrans, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("GetBlockTxs() err=", err)
		return nil, err
	}
	//tradeInstance, err := trade.NewTrade(common.HexToAddress(TradeAddr), client)
	//if err != nil {
	//	log.Println("GetBlockTxs() NewTrade() err=", err)
	//	return nil, err
	//}
	//nft1155Instance, err := nft1155.NewNft1155(common.HexToAddress(Nft1155Addr), client)
	//if err != nil {
	//	log.Println("GetBlockTxs() NewNft1155() err=", err)
	//	return nil, err
	//}
	//weth9Instance, err := weth9.NewWeth9(common.HexToAddress(Weth9Addr), client)
	//if err != nil {
	//	log.Println("GetBlockTxs() NewWeth9() err=", err)
	//	return nil, nil, err
	//}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
	if err != nil {
		log.Println("GetBlockTxs() BlockByNumber() err=", err)
		return nil, err
	}
	//UserMintDeep := UserMintDeepDef
	UserMintDeep := big.NewInt(0)
	if blockNum > 1 {
		mintdeep, err := GetUserMintDeep(blockNum - 1)
		if err != nil {
			log.Println("GetBlockTxs() GetUserMintDeep() err=", err)
			return nil, err
		}
		UserMintDeep, ok := UserMintDeep.SetString(mintdeep, 16)
		if !ok {
			log.Println("GetBlockTxs() UserMintDeep.SetString() errors.")
			return nil, err
		}
		log.Println("GetBlockTxs() UserMintDeep= ", UserMintDeep)
	}
	transT := block.Time()
	log.Println(time.Unix(int64(transT), 0))
	nfttxs := make([]*NftTx, 0, 20)
	wnfttxs := make([]*NftTx, 0, 20)
	minttxs := make([]*NftTx, 0, 20)
	wminttxs := make([]*NftTx, 0, 20)
	var exchangerInfo string
	wethchanges := make(WethChange)
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		nonce := tx.Nonce()
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Println("GetBlockTxs() TransactionReceipt() err=", err)
			return nil, err
		}
		transFlag := true
		if receipt.Status != 1 {
			log.Println("GetBlockTxs() receipt.Status != 1")
			transFlag = false
			//continue
		}
		/*if tx.To().Hex() == TradeAddr {
			for _, logger := range receipt.Logs {
				if logger.Address == common.HexToAddress(Nft1155Addr) {
					if logger.Topics[0] == common.HexToHash(TransferSingleHash) {
						nftx := NftTx{}
						m, err := nft1155Instance.ParseTransferSingle(*logger)
						if err != nil {
							log.Println("GetBlockTxs() ParseTransferSingle() err=", err)
							return nil, err
						}
						nftx.Operator = m.Operator.String()
						nftx.From = m.From.String()
						nftx.To = m.To.String()
						nftx.Value = m.Value.String()
						nftx.Contract = Nft1155Addr
						nftx.TokenId = m.Id.String()
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
						nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
						//nfttxs = append(nfttxs, &nftx)
					} else if logger.Topics[0] == common.HexToHash(RoyaltyHash) {
						r, err := nft1155Instance.ParseRoyalty(*logger)
						if err != nil {
							log.Println("GetBlockTxs() ParseRoyalty() err=", err)
							return nil, err
						}
						nftx := NftTx{}
						nftx.To = r.Receiver.String()
						nftx.Contract = Nft1155Addr
						nftx.TokenId = r.Id.String()
						nftx.Ratio = strconv.Itoa(int(r.Ratio))
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
						nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
						minttxs = append(minttxs, &nftx)
					}
				} else if logger.Address == common.HexToAddress(TradeAddr) {
					nftx := NftTx{}
					if logger.Topics[0] == common.HexToHash(PricingHash) {
						p, err := tradeInstance.ParsePRICING(*logger)
						if err != nil {
							log.Println("GetBlockTxs() ParsePRICING() err=", err)
							return nil, err
						}
						nftx.From = p.From.String()
						nftx.To = p.To.String()
						nftx.Contract = p.Addr.String()
						nftx.TokenId = p.Id.String()
						nftx.Value = p.Amount.String()
						nftx.Price = p.Price.String()
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
						nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
					} else if logger.Topics[0] == common.HexToHash(BidingHash) {
						p, err := tradeInstance.ParseBIDING(*logger)
						if err != nil {
							log.Println("GetBlockTxs() TransactionReceipt() err=", err)
							return nil, err
						}
						nftx.From = p.From.String()
						nftx.To = p.To.String()
						nftx.Contract = p.Addr.String()
						nftx.TokenId = p.Id.String()
						nftx.Value = p.Amount.String()
						nftx.Price = p.Price.String()
						nftx.TxHash = logger.TxHash.String()
						nftx.Ts = strconv.FormatUint(transT, 10)
						nftx.BlockNumber = strconv.FormatUint(logger.BlockNumber, 10)
						nftx.TransactionIndex = strconv.FormatUint(uint64(logger.TxIndex), 10)
						nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
					} else {
						continue
					}
					nfttxs = append(nfttxs, &nftx)
					wethchanges[nftx.To] = true
				}
			}
		}*/
		data := tx.Data()
		if len(data) > 10 && string(data[:10]) == "wormholes:" {
			var wormholes Wormholes
			jsonErr := json.Unmarshal(data[10:], &wormholes)
			if jsonErr != nil {
				log.Println("GetBlockTxs() wormholes type err=", err)
				continue
			}
			switch wormholes.Type {
			case WormHolesMint:
				wormMint := WormholesMint{}
				jsonErr := json.Unmarshal(data[10:], &wormMint)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormMint.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				var nftmeta NftMeta
				metabyte, _ := hex.DecodeString(wormMint.MetaUrl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesMint
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				nftx.Contract = strings.ToLower(wormMint.Exchanger)
				//nftx.TokenId = wormMint.NftAddress
				nftx.TokenId = nftmeta.TokenId
				nftx.Value = WormHolesNftCount
				nftx.Ratio = strconv.Itoa(int(wormMint.Royalty))
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftx.MetaUrl = nftmeta.Meta
				nftx.NftAddr = common.BytesToAddress(UserMintDeep.Bytes()).String()
				wminttxs = append(wminttxs, &nftx)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
			case WormHolesExchange:
				if !transFlag {
					continue
				}
				wormtrans := WormholesExchange{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.From = strings.ToLower(tx.To().String())
				nftx.To = ZeroAddr
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesPledge:
				/*if !transFlag {
					continue
				}*/
				wormtrans := WormholesPledge{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() Unmarshal err=", err)
					continue
				}
				/*from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					//return err
					continue
				}*/
				msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					//return err
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.TransType = WormHolesPledge
				nftx.From = strings.ToLower(msg.From().String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesUnPledge:
				/*if !transFlag {
					continue
				}*/
				wormtrans := WormholesPledge{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					//return err
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.TransType = WormHolesUnPledge
				nftx.From = strings.ToLower(msg.From().String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesOpenExchanger:
				wormOpen := WormholesOpenExchanger{}
				jsonErr := json.Unmarshal(data[10:], &wormOpen)
				if jsonErr != nil {
					log.Println("GetBlockTxs() Unmarshal err=", err)
					continue
				}
				wormOpen.Blocknumber = strconv.FormatUint(block.NumberU64(), 10)
				exInfo, err := json.Marshal(&wormOpen)
				if err != nil {
					log.Println("GetBlockTxs() Marshal err=", err)
					continue
				}
				if wormOpen.Name == "exchanger test." {
					exchangerInfo = string(exInfo)
					log.Println("GetBlockTxs() find open exchanger trans")
				}
			case WormHolesTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesTransfer{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesTransfer unmarshal type err=", err)
					continue
				}
				if wormtrans.NftAddress == "" {
					log.Println("GetBlockTxs() WormHolesTransfer nftaddress equal null.")
					continue
				}
				from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					log.Println("GetBlockTxs() WormHolesTransfer() err=", err)
					continue
					//return nil, err
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesTransfer
				nftx.From = strings.ToLower(from.String())
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				//nftx.Contract = wormtrans.Exchanger
				nftx.NftAddr = wormtrans.NftAddress
				nftx.Value = WormHolesNftCount
				nftx.Price = tx.Value().String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExToBuyTransfer:
				wormtrans := WormholesFixTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				fmt.Println("14=", wormtrans)
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber + wormtrans.Seller
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x" + wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					continue
					//return nil, err
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExToBuyTransfer
				nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				//nftx.Contract = wormtrans.Exchanger
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesBuyFromSellTransfer:
				wormtrans := WormholesBuyFromSellTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber
				fromaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					//todo
					//return nil, err
					continue
				}
				//if fromaddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	return nil, err
				//}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesBuyFromSellTransfer
				nftx.From = strings.ToLower(fromaddr.String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesBuyFromSellMintTransfer:
				wormtrans := WormholesBuyFromSellMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					continue
					//return nil, err
				}
				msg := wormtrans.Price + wormtrans.Royalty + wormtrans.Metaurl + wormtrans.Exclusiveflag +
					wormtrans.Exchanger + wormtrans.Blocknumber
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				/*if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					return nil, err
				}*/
				var nftmeta NftMeta
				metabyte, _ := hex.DecodeString(wormtrans.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesBuyFromSellMintTransfer
				nftxm.To = strings.ToLower(toaddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Exchanger)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := hexutil.DecodeUint64(wormtrans.Royalty)
				nftxm.Ratio = strconv.FormatUint(royalty, 10)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, &nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesBuyFromSellMintTransfer
				nftx.To = strings.ToLower(from.String())
				nftx.From = strings.ToLower(toaddr.String())
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.TokenId = nftmeta.TokenId
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExToBuyMintToSellTransfer:
				wormtrans := ExchangerMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				if ExchangeOwer != wormtrans.Seller.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x"+ wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				//2a95249bcbe73397f54562ff7a74d40b9d34a08b
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					//return nil, err
				}
				msg = wormtrans.Seller.Price + wormtrans.Seller.Royalty + wormtrans.Seller.Metaurl + wormtrans.Seller.Exclusiveflag +
					wormtrans.Seller.Exchanger + wormtrans.Seller.Blocknumber
				/*msghash = crypto.Keccak256([]byte(msg))
				hexsig, err = hexutil.Decode("0x"+ wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err =crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				fromAddr := crypto.PubkeyToAddress(*pub)
				*/
				fromAddr, err := recoverAddress(msg, wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				//if fromAddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	//return nil, err
				//}
				var nftmeta NftMeta
				metabyte, jsonErr := hex.DecodeString(wormtrans.Seller.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesExToBuyMintToSellTransfer
				nftxm.To = strings.ToLower(fromAddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Seller.Exchanger)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := strconv.Atoi(wormtrans.Seller.Royalty)
				nftxm.Ratio = strconv.Itoa(royalty)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, &nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExToBuyMintToSellTransfer
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(toaddr.String())
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				price, _ := hexutil.DecodeUint64(wormtrans.Buyer.Price)
				nftx.Price = strconv.FormatUint(price, 10)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExAuthToExBuyTransfer:
				wormtrans := WormholesFixTransAuth{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber + wormtrans.Seller
				//msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x" + wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					continue
					//return nil, errors.New("buyer address error.")
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExBuyTransfer
				nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.To = strings.ToLower(tx.To().String())
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExAuthToExMintBuyTransfer:
				wormtrans := ExchangerAuthMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				if ExchangeOwer != wormtrans.Seller.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x"+ wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				//2a95249bcbe73397f54562ff7a74d40b9d34a08b
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					//return nil, err
				}
				msg = wormtrans.Seller.Price + wormtrans.Seller.Royalty + wormtrans.Seller.Metaurl + wormtrans.Seller.Exclusiveflag +
					wormtrans.Seller.Exchanger + wormtrans.Seller.Blocknumber
				/*msghash = crypto.Keccak256([]byte(msg))
				hexsig, err = hexutil.Decode("0x"+ wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err =crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				fromAddr := crypto.PubkeyToAddress(*pub)
				*/
				fromAddr, err := recoverAddress(msg, wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				//if fromAddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	//return nil, err
				//}
				var nftmeta NftMeta
				metabyte, jsonErr := hex.DecodeString(wormtrans.Seller.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesExAuthToExMintBuyTransfer
				nftxm.To = strings.ToLower(fromAddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := strconv.Atoi(wormtrans.Seller.Royalty)
				nftxm.Ratio = strconv.Itoa(royalty)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, &nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExMintBuyTransfer
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(toaddr.String())
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Buyer.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Buyer.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			}
		}
	}
	return &NftTrans{nfttxs, minttxs, wnfttxs, wminttxs, wethchanges, exchangerInfo}, nil
}

func SelfGetBlockTxs(blockNum uint64) ([]NftTx, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("GetBlockTxs() err=", err)
		return nil, err
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
	if err != nil {
		log.Println("GetBlockTxs() BlockByNumber() err=", err)
		return nil, err
	}
	UserMintDeep := big.NewInt(0)
	if blockNum > 1 {
		mintdeep, err := GetUserMintDeep(blockNum - 1)
		if err != nil {
			log.Println("GetBlockTxs() GetUserMintDeep() err=", err)
			return nil, err
		}
		UserMintDeep, ok := UserMintDeep.SetString(mintdeep, 16)
		if !ok {
			log.Println("GetBlockTxs() UserMintDeep.SetString() errors.")
			return nil, err
		}
		log.Println("GetBlockTxs() UserMintDeep= ", UserMintDeep)
	}
	transT := block.Time()
	log.Println(time.Unix(int64(transT), 0))
	wnfttxs := make([]NftTx, 0, 20)
	wminttxs := make([]NftTx, 0, 20)
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		nonce := tx.Nonce()
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Println("GetBlockTxs() TransactionReceipt() err=", err)
			return nil, err
		}
		transFlag := true
		if receipt.Status != 1 {
			log.Println("GetBlockTxs() receipt.Status != 1")
			transFlag = false
			//continue
		}
		data := tx.Data()
		if len(data) > 10 && string(data[:10]) == "wormholes:" {
			var wormholes Wormholes
			jsonErr := json.Unmarshal(data[10:], &wormholes)
			if jsonErr != nil {
				log.Println("GetBlockTxs() wormholes type err=", err)
				continue
			}
			switch wormholes.Type {
			case WormHolesMint:
				wormMint := WormholesMint{}
				jsonErr := json.Unmarshal(data[10:], &wormMint)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormMint.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				var nftmeta NftMeta
				metabyte, _ := hex.DecodeString(wormMint.MetaUrl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesMint
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				nftx.Contract = strings.ToLower(wormMint.Exchanger)
				//nftx.TokenId = wormMint.NftAddress
				nftx.TokenId = nftmeta.TokenId
				nftx.Value = WormHolesNftCount
				nftx.Ratio = strconv.Itoa(int(wormMint.Royalty))
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftx.MetaUrl = nftmeta.Meta
				nftx.NftAddr = common.BytesToAddress(UserMintDeep.Bytes()).String()
				wminttxs = append(wminttxs, nftx)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
			case WormHolesExchange:
				if !transFlag {
					continue
				}
				wormtrans := WormholesExchange{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.From = strings.ToLower(tx.To().String())
				nftx.To = ZeroAddr
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesPledge:
				/*if !transFlag {
					continue
				}*/
				wormtrans := WormholesPledge{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() Unmarshal err=", err)
					continue
				}
				/*from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					//return err
					continue
				}*/
				msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					//return err
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.TransType = WormHolesPledge
				nftx.From = strings.ToLower(msg.From().String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesUnPledge:
				/*if !transFlag {
					continue
				}*/
				wormtrans := WormholesPledge{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					//return err
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.TransType = WormHolesUnPledge
				nftx.From = strings.ToLower(msg.From().String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesOpenExchanger:
				wormOpen := WormholesOpenExchanger{}
				jsonErr := json.Unmarshal(data[10:], &wormOpen)
				if jsonErr != nil {
					log.Println("GetBlockTxs() Unmarshal err=", err)
					continue
				}
				wormOpen.Blocknumber = strconv.FormatUint(block.NumberU64(), 10)
				exInfo, err := json.Marshal(&wormOpen)
				if err != nil {
					log.Println("GetBlockTxs() Marshal err=", err)
					continue
				}
				if wormOpen.Name == "exchanger test." {
					exchangerInfo := string(exInfo)
					log.Println("GetBlockTxs() find open exchanger trans", exchangerInfo)
				}
			case WormHolesTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesTransfer{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesTransfer unmarshal type err=", err)
					continue
				}
				if wormtrans.NftAddress == "" {
					log.Println("GetBlockTxs() WormHolesTransfer nftaddress equal null.")
					continue
				}
				from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					log.Println("GetBlockTxs() WormHolesTransfer() err=", err)
					continue
					//return nil, err
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesTransfer
				nftx.From = strings.ToLower(from.String())
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				//nftx.Contract = wormtrans.Exchanger
				nftx.NftAddr = wormtrans.NftAddress
				nftx.Value = WormHolesNftCount
				nftx.Price = tx.Value().String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesExToBuyTransfer:
				wormtrans := WormholesFixTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber + wormtrans.Seller
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x" + wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					continue
					//return nil, err
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExToBuyTransfer
				nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				//nftx.Contract = wormtrans.Exchanger
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesBuyFromSellTransfer:
				wormtrans := WormholesBuyFromSellTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber
				fromaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					//todo
					//return nil, err
					continue
				}
				//if fromaddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	return nil, err
				//}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesBuyFromSellTransfer
				nftx.From = strings.ToLower(fromaddr.String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesBuyFromSellMintTransfer:
				wormtrans := WormholesBuyFromSellMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					continue
					//return nil, err
				}
				msg := wormtrans.Price + wormtrans.Royalty + wormtrans.Metaurl + wormtrans.Exclusiveflag +
					wormtrans.Exchanger + wormtrans.Blocknumber
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				/*if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					return nil, err
				}*/
				var nftmeta NftMeta
				metabyte, _ := hex.DecodeString(wormtrans.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesBuyFromSellMintTransfer
				nftxm.To = strings.ToLower(toaddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Exchanger)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := hexutil.DecodeUint64(wormtrans.Royalty)
				nftxm.Ratio = strconv.FormatUint(royalty, 10)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesBuyFromSellMintTransfer
				nftx.To = strings.ToLower(from.String())
				nftx.From = strings.ToLower(toaddr.String())
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.TokenId = nftmeta.TokenId
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesExToBuyMintToSellTransfer:
				wormtrans := ExchangerMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				if ExchangeOwer != wormtrans.Seller.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x"+ wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				//2a95249bcbe73397f54562ff7a74d40b9d34a08b
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					//return nil, err
				}
				msg = wormtrans.Seller.Price + wormtrans.Seller.Royalty + wormtrans.Seller.Metaurl + wormtrans.Seller.Exclusiveflag +
					wormtrans.Seller.Exchanger + wormtrans.Seller.Blocknumber
				/*msghash = crypto.Keccak256([]byte(msg))
				hexsig, err = hexutil.Decode("0x"+ wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err =crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				fromAddr := crypto.PubkeyToAddress(*pub)
				*/
				fromAddr, err := recoverAddress(msg, wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				//if fromAddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	//return nil, err
				//}
				var nftmeta NftMeta
				metabyte, jsonErr := hex.DecodeString(wormtrans.Seller.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesExToBuyMintToSellTransfer
				nftxm.To = strings.ToLower(fromAddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Seller.Exchanger)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := strconv.Atoi(wormtrans.Seller.Royalty)
				nftxm.Ratio = strconv.Itoa(royalty)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExToBuyMintToSellTransfer
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(toaddr.String())
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				price, _ := hexutil.DecodeUint64(wormtrans.Buyer.Price)
				nftx.Price = strconv.FormatUint(price, 10)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesExAuthToExBuyTransfer:
				wormtrans := WormholesFixTransAuth{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
					continue
				}
				if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber + wormtrans.Seller
				//msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x" + wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					continue
					//return nil, errors.New("buyer address error.")
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExBuyTransfer
				nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.To = strings.ToLower(tx.To().String())
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			case WormHolesExAuthToExMintBuyTransfer:
				wormtrans := ExchangerAuthMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("GetBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				//from, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				//if err != nil {
				//	log.Println("GetBlockTxs() TransactionSender() err=", err)
				//	return nil, err
				//}
				if ExchangeOwer != wormtrans.Seller.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}
				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber
				/*msghash := crypto.Keccak256([]byte(msg))
				hexsig, err := hexutil.Decode("0x"+ wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err :=crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				//2a95249bcbe73397f54562ff7a74d40b9d34a08b
				toaddr := crypto.PubkeyToAddress(*pub)
				*/
				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
					//return nil, err
				}
				msg = wormtrans.Seller.Price + wormtrans.Seller.Royalty + wormtrans.Seller.Metaurl + wormtrans.Seller.Exclusiveflag +
					wormtrans.Seller.Exchanger + wormtrans.Seller.Blocknumber
				/*msghash = crypto.Keccak256([]byte(msg))
				hexsig, err = hexutil.Decode("0x"+ wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() Decode() err=", err)
					return nil, err
				}
				pub, err =crypto.SigToPub(msghash, hexsig)
				if err != nil {
					log.Println("GetBlockTxs() TransactionSender() err=", err)
					return nil, err
				}
				fromAddr := crypto.PubkeyToAddress(*pub)
				*/
				fromAddr, err := recoverAddress(msg, wormtrans.Seller.Sig)
				if err != nil {
					log.Println("GetBlockTxs() recoverAddress() err=", err)
					continue
					//return nil, err
				}
				//if fromAddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	//return nil, err
				//}
				var nftmeta NftMeta
				metabyte, jsonErr := hex.DecodeString(wormtrans.Seller.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesExAuthToExMintBuyTransfer
				nftxm.To = strings.ToLower(fromAddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := strconv.Atoi(wormtrans.Seller.Royalty)
				nftxm.Ratio = strconv.Itoa(royalty)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("GetBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExMintBuyTransfer
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(toaddr.String())
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Buyer.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Buyer.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, nftx)
			}
		}
	}
	wnfttxs = append(wnfttxs, wminttxs...)
	return wnfttxs, nil
}

func GetUserMintDeep(blockNumber uint64) (string, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetUserMintDeep() err=", err)
		return "", err
	}
	var result string
	blockN := hexutil.EncodeUint64(blockNumber)
	err = client.CallContext(context.Background(), &result, "eth_getUserMintDeep", blockN)
	if err != nil {
		log.Println("GetUserMintDeep() err=", err)
		return "", err
	}
	return result, err
}

func BuyerTransaction(to, price, nftaddr, exchanger, tosig, prv string) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	privateKey, err := crypto.HexToECDSA(prv)
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("SendTrans() err=", err)
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	blocknum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	gasLimit := uint64(51000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	p, _ := strconv.ParseUint(price, 10, 64)
	value := big.NewInt(int64(p))
	var trans BuyerTrans
	trans.Worm.Version = WormholesVersion
	trans.Worm.Type = TransTypeBuyerTrans
	trans.Worm.Price = hexutil.EncodeUint64(p)
	trans.Worm.Nftaddress = nftaddr
	trans.Worm.Sig = tosig
	trans.Worm.Exchanger = exchanger
	trans.Worm.Blocknumber = hexutil.EncodeUint64(blocknum)
	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println(data)
	toAddress := common.HexToAddress(to)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("BuyerTransaction() err=", err)
		return err
	}
	return nil
}

func MintTransaction(to string, seller Seller, buyer Buyer, prv string) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	privateKey, err := crypto.HexToECDSA(prv)
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("MintTransaction() err=", err)
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	blocknum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	gasLimit := uint64(51000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	sp, _ := strconv.ParseUint(seller.Price, 10, 64)
	value := big.NewInt(int64(sp))
	var trans MintTrans
	trans.Worm.Version = WormholesVersion
	trans.Worm.Type = TransTypeMintTrans
	p, _ := strconv.ParseUint(buyer.Price, 10, 64)
	trans.Worm.Buyer.Price = hexutil.EncodeUint64(p)
	//trans.Worm.Buyer.Nftaddress = buyer.Nftaddress
	trans.Worm.Buyer.Sig = buyer.Sig
	trans.Worm.Buyer.Exchanger = buyer.Exchanger
	trans.Worm.Buyer.Blocknumber = hexutil.EncodeUint64(blocknum)
	p, _ = strconv.ParseUint(seller.Price, 10, 64)
	trans.Worm.Seller.Price = hexutil.EncodeUint64(p)
	trans.Worm.Seller.Royalty = seller.Royalty
	trans.Worm.Seller.Metaurl = seller.Metaurl
	trans.Worm.Seller.Exchanger = seller.Exchanger
	trans.Worm.Seller.Blocknumber = hexutil.EncodeUint64(blocknum)
	trans.Worm.Seller.Sig = seller.Sig

	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println(data)
	toAddress := common.HexToAddress(to)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("MintTransaction() err=", err)
		return err
	}
	return nil
}

func ExchangerMint(seller Seller2, buyer Buyer1, fromprv string) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	fromKey, err := crypto.HexToECDSA(fromprv)
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	gasLimit := uint64(GasLimitTx1819)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	//sp, _ := strconv.ParseUint(seller.Price, 10, 64)
	value, err := hexutil.DecodeBig(buyer.Price)
	if err != nil {
		log.Println("ExchangerMint() DecodeBig err=", err)
		return err
	}
	var trans MintTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = WormHolesExToBuyMintToSellTransfer
	trans.Worm.Buyer.Price = buyer.Price
	trans.Worm.Buyer.Exchanger = buyer.Exchanger
	trans.Worm.Buyer.Blocknumber = buyer.Blocknumber
	trans.Worm.Buyer.Sig = buyer.Sig
	msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Exchanger + trans.Worm.Buyer.Blocknumber
	toAddress, err := recoverAddress(msg, trans.Worm.Buyer.Sig)
	if err != nil {
		log.Println("ExchangerMint() recoverAddress() err=", err)
		return err
	}
	trans.Worm.Seller.Price = seller.Price
	trans.Worm.Seller.Royalty = seller.Royalty
	trans.Worm.Seller.Metaurl = seller.Metaurl
	trans.Worm.Seller.Exchanger = seller.Exchanger
	trans.Worm.Seller.Blocknumber = seller.Blocknumber
	trans.Worm.Seller.Exclusiveflag = seller.Exclusiveflag
	trans.Worm.Seller.Sig = seller.Sig
	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println("ExchangerMint() value=", value.String())
	tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, data)

	log.Println("ExchangerMint() tx.Value().String()=", tx.Value().String())
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("ExchangerMint() err=", err)
		return err
	}
	return nil
}

func AuthExchangerMint(seller Seller2, buyer Buyer1, authSign string, fromprv string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	defer client.Close()
	fromKey, err := crypto.HexToECDSA(fromprv)
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}

	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//if err != nil {
	//	log.Println("AuthExchangerMint() err=", err)
	//	return "", err
	//}
	//blocknum, err := client.BlockNumber(context.Background())
	//if err != nil {
	//	log.Println("AuthExchangerMint() err=", err)
	//	return "", err
	//}
	//nonce, err := transLock.GetNonce(client, blocknum, fromAddress)
	//if err != nil {
	//	log.Println("AuthExchangeTrans() GetNonce err=", err)
	//	return "", err
	//}
	gasLimit := uint64(GasLimitTx1819)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	value, err := hexutil.DecodeBig(buyer.Price)
	if err != nil {
		log.Println("AuthExchangerMint() DecodeBig err=", err)
		return "", err
	}
	log.Println("AuthExchangerMint() price=", value.String())
	var trans AuthMintTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = WormHolesExAuthToExMintBuyTransfer
	err = json.Unmarshal([]byte(authSign), &trans.Worm.Exchangerauth)
	if err != nil {
		log.Println("AuthExchangerMint()  Unmarshal() err=", err)
		return "", err
	}
	trans.Worm.Buyer.Price = buyer.Price
	trans.Worm.Buyer.Exchanger = buyer.Exchanger
	trans.Worm.Buyer.Blocknumber = buyer.Blocknumber
	trans.Worm.Buyer.Sig = buyer.Sig
	msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Exchanger + trans.Worm.Buyer.Blocknumber
	toAddress, err := recoverAddress(msg, trans.Worm.Buyer.Sig)
	if err != nil {
		log.Println("AuthExchangerMint() recoverAddress() err=", err)
		return "", err
	}
	trans.Worm.Seller.Price = seller.Price
	trans.Worm.Seller.Royalty = seller.Royalty
	trans.Worm.Seller.Metaurl = seller.Metaurl
	trans.Worm.Seller.Exchanger = seller.Exchanger
	trans.Worm.Seller.Blocknumber = seller.Blocknumber
	trans.Worm.Seller.Exclusiveflag = seller.Exclusiveflag
	trans.Worm.Seller.Sig = seller.Sig
	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println("AuthExchangerMint() value=", value.String())
	blocknum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	nonce, err := transLock.GetNonce(client, blocknum, fromAddress)
	if err != nil {
		log.Println("AuthExchangeTrans() GetNonce err=", err)
		return "", err
	}
	tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, data)

	log.Println("AuthExchangerMint() data=", sstr)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	log.Println("AuthExchangerMint() chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
		return "", err
	}
	log.Println("AuthExchangerMint() OK")
	log.Println("AuthExchangerMint() blocknumber=", blocknum, "  txhash=", signedTx.Hash())
	return strings.ToLower(signedTx.Hash().String()), nil
}

func ExchangeTrans(buyer Buyer, fromprv string) error {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	privateKey, err := crypto.HexToECDSA(fromprv)
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	gasLimit := uint64(51000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	//p, _ := strconv.ParseUint(buyer.Price, 10, 64)
	value, err := hexutil.DecodeBig(buyer.Price)
	if err != nil {
		log.Println("ExchangeTrans() DecodeBig err=", err)
		return err
	}
	var trans ExchangerTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = WormHolesExToBuyTransfer
	trans.Worm.Buyer.Price = buyer.Price
	trans.Worm.Buyer.Exchanger = buyer.Exchanger
	trans.Worm.Buyer.Nftaddress = buyer.Nftaddress
	trans.Worm.Buyer.Blocknumber = buyer.Blocknumber
	trans.Worm.Buyer.Seller = buyer.Seller
	trans.Worm.Buyer.Sig = buyer.Sig
	msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger +
		trans.Worm.Buyer.Blocknumber + trans.Worm.Buyer.Seller
	toAddress, err := recoverAddress(msg, trans.Worm.Buyer.Sig)
	if err != nil {
		log.Println("ExchangerMint() recoverAddress() err=", err)
		return err
	}
	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println(data)
	tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("ExchangeTrans() err=", err)
		return err
	}
	return nil
}

func GetTransStatus(txHash string) (bool, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("GetTransStatus() err=", err)
		return false, err
	}
	defer client.Close()
	hash := common.HexToHash(txHash)
	log.Println("hash = ", hash, ",  txHash = ", txHash)
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Println("GetTransStatus() err=", err)
		return false, err
	}
	if receipt.Status == 1 {
		return true, nil
	} else {
		return false, err
	}
}

func AuthExchangeTrans(sell Seller1, buyer Buyer, authSign, fromprv string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	defer client.Close()
	privateKey, err := crypto.HexToECDSA(fromprv)
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("AuthExchangeTrans() err=", err)
		return "", nil
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//if err != nil {
	//	log.Println("AuthExchangeTrans() err=", err)
	//	return "", err
	//}
	//blocknum, err := client.BlockNumber(context.Background())
	//if err != nil {
	//	log.Println("AuthExchangeTrans() err=", err)
	//	return "", err
	//}
	//nonce, err := transLock.GetNonce(client, blocknum, fromAddress)
	//if err != nil {
	//	log.Println("AuthExchangeTrans() GetNonce err=", err)
	//	return "", err
	//}
	gasLimit := uint64(GasLimitTx1819)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	//p, _ := strconv.ParseUint(buyer.Price, 10, 64)
	log.Println("buyer price = ", buyer.Price)
	value, err := hexutil.DecodeBig(buyer.Price)
	if err != nil {
		log.Println("AuthExchangeTrans() DecodeBig err=", err)
		return "", err
	}
	var trans ExchangerAuthTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = WormHolesExAuthToExBuyTransfer
	err = json.Unmarshal([]byte(authSign), &trans.Worm.Exchangerauth)
	if err != nil {
		log.Println("AuthExchangeTrans() Minted Buyer Unmarshal() err=", err)
		return "", err
	}
	trans.Worm.Seller1 = sell
	trans.Worm.Buyer.Price = buyer.Price
	trans.Worm.Buyer.Exchanger = buyer.Exchanger
	trans.Worm.Buyer.Nftaddress = buyer.Nftaddress
	trans.Worm.Buyer.Blocknumber = buyer.Blocknumber
	trans.Worm.Buyer.Seller = buyer.Seller
	trans.Worm.Buyer.Sig = buyer.Sig
	msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger +
		trans.Worm.Buyer.Blocknumber + trans.Worm.Buyer.Seller
	/*msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger +
	trans.Worm.Buyer.Blocknumber*/
	toAddress, err := recoverAddress(msg, trans.Worm.Buyer.Sig)
	if err != nil {
		log.Println("AuthExchangeTrans() recoverAddress() err=", err)
		return "", err
	}
	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println("AuthExchangeTrans() data=", sstr)
	log.Println("AuthExchangeTrans() price=", value.String())
	blocknum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	nonce, err := transLock.GetNonce(client, blocknum, fromAddress)
	if err != nil {
		log.Println("AuthExchangeTrans() GetNonce err=", err)
		return "", err
	}
	tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	fmt.Println("AuthExchangeTrans() chainID=", chainID)
	fmt.Println("AuthExchangeTrans() nonce=", nonce)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("AuthExchangeTrans() err=", err)
		return "", err
	}
	log.Println("AuthExchangeTrans() OK")
	log.Println("AuthExchangeTrans() blocknumber=", blocknum, "  txhash=", signedTx.Hash().String())
	return strings.ToLower(signedTx.Hash().String()), nil
}

func BatchAuthExchangeTrans(sell Seller1, buyer Buyer, sellauthsign, buyauthsign, authSign, fromprv string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	defer client.Close()
	privateKey, err := crypto.HexToECDSA(fromprv)
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", nil
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasLimit := uint64(GasLimitTx1819)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	//p, _ := strconv.ParseUint(buyer.Price, 10, 64)
	log.Println("buyer price = ", buyer.Price)
	value, err := hexutil.DecodeBig(buyer.Price)
	if err != nil {
		log.Println("BatchAuthExchangeTrans() DecodeBig err=", err)
		return "", err
	}
	var trans ExchangerBatchAuthTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = WormHolesExSellBatchAuthTransfer
	err = json.Unmarshal([]byte(authSign), &trans.Worm.Exchangerauth)
	if err != nil {
		log.Println("BatchAuthExchangeTrans() Minted Buyer Unmarshal() err=", err)
		return "", err
	}
	if sellauthsign != "" {
		err = json.Unmarshal([]byte(sellauthsign), &trans.Worm.Sellerauth)
		if err != nil {
			log.Println("BatchAuthExchangeTrans() Minted Buyer Unmarshal() err=", err)
			return "", err
		}
	}
	if buyauthsign != "" {
		err = json.Unmarshal([]byte(buyauthsign), &trans.Worm.Buyauth)
		if err != nil {
			log.Println("BatchAuthExchangeTrans() Minted Buyer Unmarshal() err=", err)
			return "", err
		}
	}
	blocknum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	trans.Worm.Seller1 = sell
	if sellauthsign != "" {
		trans.Worm.Seller1.Blocknumber = hexutil.EncodeUint64(blocknum + 1000)
		msg := trans.Worm.Seller1.Price + trans.Worm.Seller1.Nftaddress + trans.Worm.Seller1.Exchanger +
			trans.Worm.Seller1.Blocknumber
		sig, err := WormholesSign(msg, privateKey)
		if err != nil {
			log.Println("BatchAuthExchangeTrans() WormholesSign() err=", err)
			return "", err
		}
		trans.Worm.Seller1.Sig = sig
	}
	var toAddress *common.Address
	if buyauthsign == "" {
		trans.Worm.Buyer.Price = buyer.Price
		trans.Worm.Buyer.Exchanger = buyer.Exchanger
		trans.Worm.Buyer.Nftaddress = buyer.Nftaddress
		trans.Worm.Buyer.Blocknumber = buyer.Blocknumber
		trans.Worm.Buyer.Seller = buyer.Seller
		trans.Worm.Buyer.Sig = buyer.Sig
		msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger +
			trans.Worm.Buyer.Blocknumber + trans.Worm.Buyer.Seller
		/*msg := trans.Worm.Buyer.Price + trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger +
		trans.Worm.Buyer.Blocknumber*/
		toAddress, err = recoverAddress(msg, trans.Worm.Buyer.Sig)
		if err != nil {
			log.Println("BatchAuthExchangeTrans() recoverAddress() err=", err)
			return "", err
		}
	} else {
		msg := trans.Worm.Buyauth.Exchanger + trans.Worm.Buyauth.Blocknumber
		toAddress, err = recoverAddress(msg, trans.Worm.Buyauth.Sig)
		if err != nil {
			log.Println("BatchAuthExchangeTrans() recoverAddress() err=", err)
			return "", err
		}
		trans.Worm.Buyer = buyer
		trans.Worm.Buyer.Blocknumber = hexutil.EncodeUint64(blocknum + 1000)
		msg = trans.Worm.Buyer.Price + trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger +
			trans.Worm.Buyer.Blocknumber + trans.Worm.Buyer.Seller
		sig, err := WormholesSign(msg, privateKey)
		if err != nil {
			log.Println("BatchAuthExchangeTrans() WormholesSign() err=", err)
			return "", err
		}
		trans.Worm.Buyer.Sig = sig
	}

	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println("BatchAuthExchangeTrans() data=", sstr)
	log.Println("BatchAuthExchangeTrans() price=", value.String())
	//blocknum, err := client.BlockNumber(context.Background())
	//if err != nil {
	//	log.Println("BatchAuthExchangeTrans() err=", err)
	//	return "", err
	//}
	nonce, err := transLock.GetNonce(client, blocknum, fromAddress)
	if err != nil {
		log.Println("AuthExchangeTrans() GetNonce err=", err)
		return "", err
	}
	tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	fmt.Println("BatchAuthExchangeTrans() chainID=", chainID)
	fmt.Println("BatchAuthExchangeTrans() nonce=", nonce)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("BatchAuthExchangeTrans() err=", err)
		return "", err
	}
	log.Println("BatchAuthExchangeTrans() OK")
	log.Println("BatchAuthExchangeTrans() blocknumber=", blocknum, "  txhash=", signedTx.Hash().String())
	return strings.ToLower(signedTx.Hash().String()), nil
}

func ForceBuyingAuthExchangeTrans(buyer Buyer, buyauthsign, authSign, fromprv string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	defer client.Close()
	privateKey, err := crypto.HexToECDSA(fromprv)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", nil
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasLimit := uint64(GasLimitTx1819)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	var trans ExchangerForceBuyingAuthTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = WormHolesExForceBuyingAuthTransfer
	err = json.Unmarshal([]byte(authSign), &trans.Worm.Exchangerauth)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() Minted Buyer Unmarshal() err=", err)
		return "", err
	}
	if buyauthsign == "" {
		log.Println("ForceBuyingAuthExchangeTrans() buyauthsign err=", "buyauthsign is null")
		return "", errors.New("buyauthsign is null")
	}
	err = json.Unmarshal([]byte(buyauthsign), &trans.Worm.Buyauth)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() Minted Buyer Unmarshal() err=", err)
		return "", err
	}
	blocknum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	var toAddress *common.Address
	msg := trans.Worm.Buyauth.Exchanger + trans.Worm.Buyauth.Blocknumber
	toAddress, err = recoverAddress(msg, trans.Worm.Buyauth.Sig)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() recoverAddress() err=", err)
		return "", err
	}
	trans.Worm.Buyer = buyer
	trans.Worm.Buyer.Blocknumber = hexutil.EncodeUint64(blocknum + 100000)
	msg = trans.Worm.Buyer.Nftaddress + trans.Worm.Buyer.Exchanger + trans.Worm.Buyer.Blocknumber + trans.Worm.Buyer.Seller
	sig, err := WormholesSign(msg, privateKey)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() WormholesSign() err=", err)
		return "", err
	}
	trans.Worm.Buyer.Sig = sig

	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	log.Println("ForceBuyingAuthExchangeTrans() data=", sstr)
	nonce, err := transLock.GetNonce(client, blocknum, fromAddress)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() GetNonce err=", err)
		return "", err
	}
	tx := types.NewTransaction(nonce, *toAddress, nil, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	fmt.Println("ForceBuyingAuthExchangeTrans() chainID=", chainID)
	fmt.Println("ForceBuyingAuthExchangeTrans() nonce=", nonce)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("ForceBuyingAuthExchangeTrans() err=", err)
		return "", err
	}
	log.Println("ForceBuyingAuthExchangeTrans() OK")
	log.Println("ForceBuyingAuthExchangeTrans() blocknumber=", blocknum, "  txhash=", signedTx.Hash().String())
	return strings.ToLower(signedTx.Hash().String()), nil
}

type ExchangerSnftTrans struct {
	Worm WormholesSnftTrans `json:"wormholes"`
}

type WormholesSnftTrans struct {
	Version       string        `json:"version"`
	Type          uint8         `json:"type"` //23
	Dir           string        `json:"dir"`
	StartIndex    string        `json:"start_index"`
	Number        uint64        `json:"number"`
	Royalty       uint32        `json:"royalty"`
	Creator       string        `json:"creator"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerNFTTrans struct {
	Worm WormholesMint `json:"wormholes"`
}
type ExchangerAuthNFTTrans struct {
	Worm WormholesBuyFromSellTrans `json:"wormholes"`
}
type NominatedNFT struct {
	Dir        string         `json:"dir"`
	StartIndex uint64         `json:"start_index"`
	Number     uint64         `json:"number"`
	Royalty    uint32         `json:"royalty"`
	Address    common.Address `json:"address"`
	Creator    common.Address `json:"creator"`
	VoteWeight uint64         `json:"vote_weight"`
}

func GetNominatedNFTInfo(blockNumber *big.Int) (*NominatedNFT, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetNominatedNFTInfo() err=", err)
		return nil, err
	}
	var result NominatedNFT
	err = client.CallContext(context.Background(), &result, "eth_getNominatedNFTInfo", toBlockNumArg(blockNumber))
	if err != nil {
		log.Println("GetAccountInfo() err=", err)
		return nil, err
	}
	return &result, err
}

func SendSnftTrans(dir, authSign string) error {
	fmt.Println(SuperAdminAddr)
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		fmt.Println("SendTrans() Dial err=", err)
		return err
	}
	privateKey, err := crypto.HexToECDSA(SuperAdminAddr)
	if err != nil {
		fmt.Println("ExchangeTrans() HexToECDSA err=", err)
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("ExchangeTrans() publicKeyECDSA err=", err)
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("SendTrans() PendingNonceAt err=", err)
		return err
	}
	fmt.Println("SendTrans() nonce=", nonce)
	p, _ := strconv.ParseUint("0", 10, 64)
	value := big.NewInt(int64(p))      // in wei (1 eth)
	gasLimit := uint64(GasLimitTx1819) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("SendTrans() SuggestGasPrice err=", err)
		return err
	}
	toAddress := common.HexToAddress(SuperAdminAddr)
	var trans ExchangerSnftTrans
	trans.Worm.Version = WormHolesVerseion
	trans.Worm.Type = 24
	trans.Worm.Number = 256
	trans.Worm.StartIndex = "00"
	trans.Worm.Dir = dir
	trans.Worm.Royalty = 100
	trans.Worm.Creator = ExchangeOwer
	err = json.Unmarshal([]byte(authSign), &trans.Worm.Exchangerauth)
	if err != nil {
		log.Println("SendTrans()  Unmarshal() err=", err)
		return err
	}
	str, err := json.Marshal(trans)
	sstr := strings.Replace(string(str), "\"wormholes\"", "wormholes", -1)
	sstr = sstr[1 : len(sstr)-1]
	data := []byte(sstr)
	fmt.Println(string(data))
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("SendTrans()NetworkID  err=", err)
		return err
	}
	fmt.Println("tx=", tx.Hash())
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("SendTrans() SignTx err=", err)
		return err
	}
	fmt.Println("signedTx=", signedTx.Hash())
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("SendTrans() SendTransaction err=", err)
		return err
	}
	return nil
}
