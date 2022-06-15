package contracts

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nftexchange/nftserver/common/contracts/trade"
	//"github.com/nftexchange/nftserver/models"

	//"github.com/nftexchange/nftserver/models"
	"math/big"
	"strconv"
)

//SendDealAuctionTx(auctionRec.Ownaddr, bidRecs.Bidaddr, auctionRec.Contract,
//				auctionRec.Tokenid, price, bidRecs.Tradesig)
func Auction(from, to, nftAddr, tokenId, amount, price, sig string) (string, error) {
	client, err := ethclient.Dial(EthNode)
	//client, err := ethclient.Dial("https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//privateKey, err := crypto.HexToECDSA("564ea566096d3de340fc5ddac98aef672f916624c8b0e4664a908cd2a2d156fe")
	privateKey, err := crypto.HexToECDSA(TradeAuthAddrPrv)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", errors.New("publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	address := common.HexToAddress(TradeAddr)
	instance, err := trade.NewTrade(address, client)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	netid, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//auth := bind.NewKeyedTransactor(privateKey)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, netid)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	_from := common.HexToAddress(from)
	_to := common.HexToAddress(to)
	_nftaddr := common.HexToAddress(nftAddr)
	iprice, _ := strconv.ParseInt(price, 10, 64)
	_price := big.NewInt(int64(iprice))
	inftid, _ := strconv.ParseInt(tokenId, 10, 64)
	_nftid := big.NewInt(int64(inftid))
	_sig, _ := hexutil.Decode(sig)
	iamount, _ := new(big.Int).SetString(amount, 10)
	trans, err := instance.Biding1155(auth, _nftaddr, _from, _to, _nftid, iamount, _price, _sig, []byte{})
	if err != nil {
		fmt.Println(err)
		return "end", err
	}
	fmt.Println("Auction() txhash=", trans.Hash().String())
	return trans.Hash().String(), nil
}

//contract,owner,metaUrl,tokenId,amount,royalty
func AuctionAndMint(from, to, nftAddr, tokenId, price, amount, royaltyRatio, tokenURI, sig string) (string, error) {
	//err, createSig := ethhelper.GenCreateNftSign(nftAddr, from, tokenURI, tokenId, amount, royaltyRatio)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//MintSign(contract string, toAddr string, tokenId string, count string, royalty string, tokenUri string, prv string) (string, error) {
	createSig, err := MintSign(nftAddr, from, tokenId, amount, royaltyRatio, tokenURI, AdminMintPrv)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//client, err := ethclient.Dial(EthNode)
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(TradeAuthAddrPrv)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", errors.New("publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	address := common.HexToAddress(TradeAddr)
	instance, err := trade.NewTrade(address, client)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	netid, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//auth := bind.NewKeyedTransactor(privateKey)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, netid)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice
	_from := common.HexToAddress(from)
	_to := common.HexToAddress(to)
	_nftaddr := common.HexToAddress(nftAddr)
	iprice, _ := strconv.ParseInt(price, 10, 64)
	_price := big.NewInt(int64(iprice))
	inftid, _ := strconv.ParseInt(tokenId, 10, 64)
	_nftid := big.NewInt(int64(inftid))
	_sig, _ := hexutil.Decode(sig)
	rayalty, _ := strconv.ParseInt(royaltyRatio, 10, 16)
	_minerSig, _ := hexutil.Decode(createSig)
	amount1, _ := new(big.Int).SetString(amount, 0)
	trans, err := instance.Biding1155Mint(auth, _nftaddr, _from, _to, _nftid, amount1, _price, uint16(rayalty), tokenURI, _minerSig, _sig, []byte{})
	if err != nil {
		fmt.Println(err)
		return "end", err
	}
	fmt.Println("AuctionAndMint() txhash=", trans.Hash().String())
	return trans.Hash().String(), nil
}
