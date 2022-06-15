package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/common/contracts/trade"
	"math/big"
	"strconv"
	"time"
)


type ResponseSell struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Totalcount int		`json:"total_count"`
}

func Sell(userAddr string, sellType ,contract string, tokenId string, price string, token string, workKey, userKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "sell"
	datam := make(map[string]string)
	datam["currency_type"] = "eth"
	datam["day"] = "1"
	datam["nft_contract_addr"] = contract
	datam["nft_token_id"] = tokenId
	datam["pay_channel"] = "eth"
	datam["price1"] = price
	datam["price2"] = ""
	datam["selltype"] = sellType
	datam["trade_sig"] = ""
	datam["user_addr"] = userAddr
	nonce, err := contracts.GetNonce(contract, userAddr, tokenId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	sign, err := contracts.TradeSign(1, contract, tokenId, "1",
		 price + "000000000", nonce.String(), userKey)
	if err != nil {
		fmt.Println("tradesign err=", err)
		return err
	}
	datam["trade_sig"] = sign
	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("Sell() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("Sell() HttpSendRev() err=", err)
		return err
	}
	var revData ResponseSell
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("Sell() Unmarshal err=", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

type ResponseBuy struct {
	Code 	string 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Totalcount int		`json:"total_count"`
}

func MakeOffer(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "buy"
	datam := make(map[string]string)
	datam["currency_type"] = "weth"
	datam["dead_time"] = strconv.FormatInt(time.Now().Unix() + 100000, 10)
	datam["nft_contract_addr"] = contract
	datam["nft_token_id"] = tokenId
	datam["pay_channel"] = "weth"
	datam["price"] = price
	datam["user_addr"] = userAddr

	nonce, err := contracts.GetNonce(contract, userAddr, tokenId)
	if err != nil {
		fmt.Println(err)
	}
	sign, err := contracts.TradeSign(3, contract, tokenId, "1",
		price + "000000000", nonce.String(), userKey)
	if err != nil {
		fmt.Println("tradesign err=", err)
	}
	datam["trade_sig"] = sign
	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("Sell() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("Sell() HttpSendRev() err=", err)
		return err
	}
	var revData ResponseBuy
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("Sell() Unmarshal err=", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

func CancelMakeOffer(userAddr string, contract string, tokenId string, token string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "cancelBuy"
	datam := make(map[string]string)
	datam["nft_contract_addr"] = contract
	datam["nft_token_id"] = tokenId
	datam["trade_sig"] = ""
	datam["user_addr"] = userAddr
	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("Sell() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("Sell() HttpSendRev() err=", err)
		return err
	}
	var revData ResponseBuy
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("Sell() Unmarshal err=", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

func CancelSell(userAddr string, contract string, tokenId string, token string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "cancelSell"
	datam := make(map[string]string)
	datam["user_addr"] = userAddr
	datam["nft_contract_addr"] = contract
	datam["nft_token_id"] = tokenId

	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("CancelSell() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("CancelSell() HttpSendRev() err=", err)
		return err
	}
	var revData ResponseSell
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("CancelSell() Unmarshal err=", err)
		return err
	}
	if revData.Code != "200" {
		fmt.Println("CancelSell() response err=", err)
		return errors.New(revData.Msg)
	}
	return nil
}

func BuyMint(contract, from, to, tokenId, royalty, price string, mintSig, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	//client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		fmt.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress("0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"), client)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	//if oldNonce == 0 {
	//	oldNonce = nonce
	//} else {
	//	for {
	//		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//		if err != nil {
	//			fmt.Println(err)
	//			fmt.Println("Buy() PendingNonceAt() err=", err)
	//			return nil, err
	//		}
	//		if oldNonce != nonce {
	//			oldNonce = nonce
	//			break
	//		}
	//	}
	//}
	fmt.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() SuggestGasPrice() err=", err)
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
	//Pricing1155Mint(opts *bind.TransactOpts, _addr common.Address, _from common.Address,
	//_to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16,
	//_tokenURI string, _minerSig []byte, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	n, err := instance.Pricing1155Mint(auth, common.HexToAddress(contract),
		common.HexToAddress(from), common.HexToAddress(to), big.NewInt(int64(tokenid)), big.NewInt(int64(1)),
		uint16(r), "", mintSig, tradeSig, []byte{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func Buy(contract, from, to, tokenId, price string, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	//client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		fmt.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress("0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"), client)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	fmt.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromKey)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)

	n, err := instance.Pricing1155(auth, common.HexToAddress(contract),common.HexToAddress(from),
		common.HexToAddress(to), big.NewInt(int64(tokenid)), big.NewInt(int64(1)), tradeSig, []byte{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func BuyBidding(contract, from, to, tokenId, price string, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(EthNode)
	//client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println("Buy() err=", err)
		return nil, err
	}
	//privateKey, err := crypto.HexToECDSA(fromKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		fmt.Println("Buy() publickey err=", err)
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := trade.NewTrade(common.HexToAddress("0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"), client)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() NewTrade() err=", err)
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() PendingNonceAt() err=", err)
		return nil, err
	}
	fmt.Println("Buy() nonce=", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromKey)
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := strconv.ParseUint(price, 10, 64)
	//auth.Value = big.NewInt(int64(p))
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice
	tokenid, _ := strconv.ParseUint(tokenId, 10, 64)
	//Biding1155(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address,
	//_id *big.Int, _amount *big.Int, _price *big.Int, _toSig []byte, _data []byte)
	n, err := instance.Biding1155(auth, common.HexToAddress(contract),common.HexToAddress(from),
		common.HexToAddress(to), big.NewInt(int64(tokenid)), big.NewInt(int64(1)), big.NewInt(int64(p)), tradeSig, []byte{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Buy() SuggestGasPrice() err=", err)
		return nil, err
	}
	return n, nil
}

func Buying(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "buying"
	datam := make(map[string]string)
	datam["currency_type"] = "weth"
	datam["dead_time"] = strconv.FormatInt(time.Now().Unix() + 100000, 10)
	datam["nft_contract_addr"] = contract
	datam["nft_token_id"] = tokenId
	datam["pay_channel"] = "weth"
	datam["price"] = price
	datam["user_addr"] = userAddr

	nonce, err := contracts.GetNonce(contract, userAddr, tokenId)
	if err != nil {
		fmt.Println(err)
	}
	sign, err := contracts.TradeSign(3, contract, tokenId, "1",
		price + "000000000", nonce.String(), userKey)
	if err != nil {
		fmt.Println("tradesign err=", err)
	}
	datam["trade_sig"] = sign
	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("Sell() sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("Sell() HttpSendRev() err=", err)
		return err
	}
	var revData ResponseBuy
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("Sell() Unmarshal err=", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}
