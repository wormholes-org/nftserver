package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
AdminPrv = "8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed"
BatchBuyAddr = "0x5E83e4c3Bc80769B4d67Fc4CB577b352C7B658Bf"
BatchBuyPrv = "5f7407151b0539359d216d38534893b203546aedbd7e8dc05e2fddc8423f2ab1"
BatchTokenId = "2597785300040"
SellPrice = "1000"
BuyPrice = "1000"
Nft1155Contract = "0xa1e67a33e090afe696d7317e05c506d7687bb2e5"

)

func TestSellSingle(t *testing.T) {
	testCount := 1
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestSellSingle() err=", err)
		return
	}
	fmt.Println("TestSellSingle() login end.")
	fmt.Println("start Test TestSellSingle.")
	userAddr := crypto.PubkeyToAddress(tKey[0].LogKey.PublicKey).String()
	userNft, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[0], tKey[0].WorkKey)
	if err != nil {
		fmt.Println("TestSellSingle() QueryUserNFTList() err=", err)
		return
	}
	err = Sell(userAddr, "FixPrice", userNft[0].NftContractAddr, userNft[0].NftTokenId, SellPrice, tokens[0], tKey[0].WorkKey, tKey[0].LogKey)
	if err != nil {
		fmt.Println("TestSellSingle()  err=", err, "userAddr=", userAddr)
		//return
	}

	mintSign, err := contracts.MintSign(userNft[0].NftContractAddr, userAddr,
		userNft[0].NftTokenId, "1", "200", "", AdminPrv);
	if err != nil {
		fmt.Println("TestSellSingle() MintSign()  err=", err)
	}
	fmt.Println("TestSellSingle() mintSign0=", mintSign)
	nonce, err := contracts.GetNonce(userNft[0].NftContractAddr, rechargeAddr, userNft[0].NftTokenId)
	if err != nil {
		fmt.Println(err)
	}
	key, err := crypto.HexToECDSA(rechargePrv)
	if err != nil {
		fmt.Println("TestSellSingle() key err=", err)
		return
	}
	tradeSign, err := contracts.TradeSign(1, userNft[0].NftContractAddr,
		userNft[0].NftTokenId, "1", SellPrice + "000000000", nonce.String(), key)
	if err != nil {
		fmt.Println("TestSellSingle() tradesign err=", err)
		return
	}
	fmt.Println("TestSellSingle() tradeSign=", tradeSign)
	key, err = crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("Sell() key err=", err)
		return
	}
	//BuyMint(contract, from, to, tokenId, royalty, price string, mintSig, tradeSig []byte, fromKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	tx, err := BuyMint(userNft[0].NftContractAddr, userAddr,
		BatchBuyAddr, userNft[0].NftTokenId,
		"200", SellPrice + "000000000", common.FromHex(mintSign), common.FromHex(tradeSign), key)
	if err != nil {
		fmt.Println("Sell() Buy() err=", err, "userAddr=", userAddr)
		return
	}
	fmt.Println("Sell() tx.hash=", tx.Hash())
	fmt.Println("end test TestSellSingle().")
}

func TestMakeOfferSingle(t *testing.T) {
	testCount := 1
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestMakeOfferSingle() err=", err)
		return
	}
	fmt.Println("TestMakeOfferSingle() login end.")
	fmt.Println("start Test TestMakeOfferSingle.")
	BiddingAddr := crypto.PubkeyToAddress(tKey[0].LogKey.PublicKey).String()
	err = MakeOffer(BiddingAddr, "0xa1e67a33e090afe696d7317e05c506d7687bb2e5",
		"9981050309826", BuyPrice, tokens[0], tKey[0].LogKey, tKey[0].WorkKey)
	if err != nil {
		fmt.Println("Sell() key err=", err)
		return
	}
	key, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("Sell() key err=", err)
		return
	}
	ts, _ := hexutil.Decode("0xe9ee13e8619a459e11fd78aac3eb2f9e20094281e8c8de4c25c0ea5758da572c39bd3038e6bacd74c22e49943b28114a85aa614d382dd4db3579a0ca631640661c")
	tx, err := BuyBidding("0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		BiddingAddr,"9981050309826", BuyPrice + "000000000", ts, key)
	if err != nil {
		fmt.Println("TestMakeOfferSingle() CancelMakeOffer() err=", err)
		return
	}
	fmt.Println("TestMakeOfferSingle() tx.hash=", tx.Hash())
	fmt.Println("end test TestMakeOfferSingle().")
}

func TestHighestToOtherAddrsSingle(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestHighestToOtherAddrsSingle() err=", err)
		return
	}
	fmt.Println("TestHighestToOtherAddrsSingle() login end.")
	fmt.Println("start Test TestHighestToOtherAddrsSingle.")
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestHighestToOtherAddrsSingle() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestHighestToOtherAddrsSingle() login err=", err)
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	//Sell(userAddr string, sellType ,contract string, tokenId string, price string, token string, workKey, userKey *ecdsa.PrivateKey) error {
	err = Sell(useraddr, "HighestBid", Nft1155Contract, BatchTokenId,
		SellPrice, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestHighestToOtherAddrsSingle()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := crypto.PubkeyToAddress(tKey[0].LogKey.PublicKey).String()
	//MakeOffer(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tokens[0], tKey[0].LogKey, tKey[0].WorkKey)
	if err != nil {
		fmt.Println("TestHighestToOtherAddrsSingle() key err=", err)
		return
	}
	BiddingAddr = crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1002", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestHighestToOtherAddrsSingle() key err=", err)
		return
	}
	fmt.Println("end test TestHighestToOtherAddrsSingle().")
}

func TestHighestMintToOtherAddrsSingle(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestHighestMintToOtherAddrsSingle() err=", err)
		return
	}
	fmt.Println("TestHighestMintToOtherAddrsSingle() login end.")
	fmt.Println("start Test TestHighestMintToOtherAddrsSingle.")
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("NewCollect() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("NewCollect() login err=", err)
		return
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	err = Sell(useraddr, "HighestBid", Nft1155Contract, BatchTokenId,
		SellPrice, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestHighestMintToOtherAddrsSingle()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	//MakeOffer(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestHighestMintToOtherAddrsSingle() key err=", err)
		return
	}
	BiddingAddr = crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String()
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1002", tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestHighestMintToOtherAddrsSingle() key err=", err)
		return
	}
	fmt.Println("end test TestHighestMintToOtherAddrsSingle().")
}

func TestSellToBatchAddr(t *testing.T) {
	const testCount = 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestSell() err=", err)
		return
	}
	fmt.Println("TestSell() login end.")
	fmt.Println("start Test TestSell.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			userNft, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestSell() QueryUserNFTList() err=", err)
				return
			}
			err = Sell(userAddr, "FixPrice", userNft[0].NftContractAddr, userNft[0].NftTokenId, SellPrice, tokens[i], tKey[i].WorkKey, tKey[i].LogKey)
			if err != nil {
				fmt.Println("TestSell()  err=", err, "userAddr=", userAddr)
				return
			}
			fmt.Println("TestSell() Sell() OK.", err, "userAddr=", userAddr)
		}(i)
	}
	wd.Wait()

	for i := 0; i < testCount; i++ {
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			userNft, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestSell() QueryUserNFTList() err=", err)
				continue
			}
			nftInfo, err := QueryNFT(userNft[0].NftContractAddr, userNft[0].NftTokenId, tokens[i])
			if err != nil {
				fmt.Println("TestSell() QueryNFT  err=", err, "userAddr=", userAddr)
				continue
			}
			r := strconv.Itoa(nftInfo.Royalty)
			p := strconv.FormatUint(nftInfo.Auction.Startprice, 10) + "000000000"
			key, err := crypto.HexToECDSA(BatchBuyPrv)
			if err != nil {
				fmt.Println("TestSell() key err=", err)
				continue
			}
			t, _ := hexutil.Decode(nftInfo.Auction.Tradesig)
			mint, _ := hexutil.Decode(nftInfo.Approve)
			if nftInfo.Mintstate == models.Minted.String() {
				tx, err := Buy(nftInfo.NftContractAddr,  nftInfo.OwnerAddr, BatchBuyAddr, nftInfo.NftTokenId,
					p, t, key)
				if err != nil {
					fmt.Println("TestSell() Buy() err=", err, "userAddr=", userAddr)
					continue
				}
				fmt.Println("TestSell() OK. tx.hash=", tx.Hash())
			} else {
				tx, err := BuyMint(nftInfo.NftContractAddr, nftInfo.OwnerAddr, BatchBuyAddr, nftInfo.NftTokenId,
					r, p, mint, t, key)
				if err != nil {
					fmt.Println("TestSell() BuyMint() err=", err, "userAddr=", userAddr)
					continue
				}
				fmt.Println("TestSell() OK. tx.hash=", tx.Hash())
			}
			time.Sleep(1*time.Second)
	}
	fmt.Println("end test TestSell().")
}

func TestSellBiddingToOthers(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestSellBiddingToOthers() err=", err)
		return
	}
	fmt.Println("TestSellBiddingToOthers() login end.")
	fmt.Println("start Test TestSellBiddingToOthers.")
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestSellBiddingToOthers() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestSellBiddingToOthers() login err=", err)
	}
	userNftCnt := strconv.Itoa(testCount)
	userNft, err := QueryUserNFTList(BatchBuyAddr, userNftCnt, "0", batchToken, batchkey)
	if err != nil {
		fmt.Println("TestSellBiddingToOthers() QueryUserNFTList() err=", err)
		return
	}
	wd := sync.WaitGroup{}
	for i, nft := range userNft {
		wd.Add(1)
		go func(i int, contract, tokenId string) {
			defer wd.Done()
			BiddingAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			err = MakeOffer(BiddingAddr, contract, tokenId,
				BuyPrice, tokens[i], tKey[i].LogKey, tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestSellBiddingToOthers() MakeOffer() err=", err, "BiddingAddr=", BiddingAddr)
				return
			}
			fmt.Println("TestSellBiddingToOthers() MakeOffer() Ok.", "BiddingAddr=", BiddingAddr)
		}(i, nft.NftContractAddr, nft.NftTokenId)
	}
	wd.Wait()
	for i, nft := range userNft {
		userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
		nftInfo, err := QueryNFT(nft.NftContractAddr, nft.NftTokenId, batchToken)
		if err != nil {
			fmt.Println("TestSellBiddingToOthers() QueryNFT  err=", err, "userAddr=", userAddr)
			continue
		}
		if len(nftInfo.Bids) ==0 {
			continue
		}
		p := strconv.FormatUint(nftInfo.Bids[0].Price, 10) + "000000000"
		key, err := crypto.HexToECDSA(BatchBuyPrv)
		if err != nil {
			fmt.Println("TestSellBiddingToOthers() key err=", err)
			continue
		}
		ts, _ := hexutil.Decode(nftInfo.Bids[0].Tradesig)
		tx, err := BuyBidding(nft.NftContractAddr, BatchBuyAddr, nftInfo.Bids[0].Bidaddr,
			               nft.NftTokenId, p, ts, key)
		if err != nil {
			fmt.Println("TestSellBiddingToOthers() Buy() err=", err, "userAddr=", userAddr)
			continue
		}
		fmt.Println("TestSellBiddingToOthers() tx.hash=", tx.Hash())
		time.Sleep(1*time.Second)
	}
	fmt.Println("end test TestSellBiddingToOthers().")
}

func TestSellBiddingToBatchAddr(t *testing.T) {
	const testCount = 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestSellBiddingToBatchAddr() err=", err)
		return
	}
	fmt.Println("TestSellBiddingToBatchAddr() login end.")
	fmt.Println("start Test TestSellBidding.")
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestSellBiddingToBatchAddr() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestSellBiddingToBatchAddr() login err=", err)
	}
	//userNftCnt := strconv.Itoa(testCount)
	//userNft, err := QueryUserNFTList(BatchBuyAddr, userNftCnt, "0", batchToken, batchkey)
	//if err != nil {
	//	fmt.Println("TestSellBiddingToOtherAddr() QueryUserNFTList() err=", err)
	//	return
	//}
	wd := sync.WaitGroup{}
	var SellNft [testCount]*UserNft
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			userNft, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestSellBiddingToBatchAddr() QueryUserNFTList() err=", err)
				return
			}
			err = MakeOffer(BatchBuyAddr, userNft[0].NftContractAddr, userNft[0].NftTokenId,
				BuyPrice, batchToken, batchkey, batchkey)
			if err != nil {
				fmt.Println("TestSellBiddingToBatchAddr() MakeOffer() err=", err, "SellerAddr=", userAddr)
				return
			}
			SellNft[i] = &userNft[0]
		}(i)
	}
	wd.Wait()
	for i := 0; i < testCount; i++ {
		if SellNft[i] == nil {
			continue
		}
		userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
		userNft, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
		if err != nil {
			fmt.Println("TestSellBiddingToBatchAddr() QueryUserNFTList() err=", err)
			return
		}
		nftInfo, err := QueryNFT(userNft[0].NftContractAddr, userNft[0].NftTokenId, tokens[i])
		if err != nil {
			fmt.Println("TestSellBiddingToBatchAddr() QueryNFT  err=", err, "userAddr=", userAddr)
			continue
		}
		if len(nftInfo.Bids) ==0 {
			continue
		}
		p := strconv.FormatUint(nftInfo.Bids[0].Price, 10) + "000000000"
		ts, _ := hexutil.Decode(nftInfo.Bids[0].Tradesig)
		tx, err := BuyBidding(userNft[0].NftContractAddr, userAddr, nftInfo.Bids[0].Bidaddr,
			userNft[0].NftTokenId, p, ts, tKey[i].LogKey)
		if err != nil {
			fmt.Println("TestSellBiddingToBatchAddr() Buy() err=", err, "userAddr=", userAddr)
			continue
		}
		fmt.Println("TestSellBiddingToBatchAddr() tx.hash=", tx.Hash())
		fmt.Println("TestSellBiddingToBatchAddr() OK.", "userAddr", userAddr)
		time.Sleep(1*time.Second)
	}
	fmt.Println("end test TestSellBiddingToBatchAddr().")
}

func TestCancelBiddings(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestCancelBiddings() err=", err)
		return
	}
	fmt.Println("TestCancelBiddings() login end.")
	fmt.Println("start Test TestCancelBiddings.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			userNft, err := QueryUserBidList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestCancelBiddings() QueryUserNFTList() err=", err)
				return
			}
			err = CancelMakeOffer(userAddr, userNft[0].NftContractAddr, userNft[0].NftTokenId, tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestCancelBiddings()  err=", err, "userAddr=", userAddr)
				return
			}
			fmt.Println("TestCancelBiddings()  Ok.", "userAddr=", userAddr)
		}(i)
	}
	wd.Wait()
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestCancelBiddings() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestCancelBiddings() login err=", err)
		return
	}
	userNftCnt := strconv.Itoa(testCount)
	userNft, err := QueryUserNFTList(BatchBuyAddr, userNftCnt, "0", batchToken, batchkey)
	if err != nil {
		fmt.Println("TestCancelBiddings() QueryUserNFTList() err=", err)
		return
	}
	wd = sync.WaitGroup{}
	for _, nft := range userNft {
		//wd.Add(1)
		//go func(nft UserNft) {
		//	defer wd.Done()
			userAddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
			err = CancelMakeOffer(userAddr, nft.NftContractAddr, nft.NftTokenId, batchToken, batchkey)
			if err != nil {
				fmt.Println("TestCancelBiddings()  err=", err, "userAddr=", nft.OwnerAddr)
				//return
			}
			fmt.Println("TestCancelBiddings()  OK.", "userAddr=", nft.OwnerAddr)
		//}(nft)
	}
	wd.Wait()
	fmt.Println("end test TestCancelBiddings().")
}

func TestSellToOtherAddr(t *testing.T) {
	testCount := 10
	tKey, _, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestSellToOtherAddr() log() err=", err)
		return
	}
	fmt.Println("TestSellToOtherAddr() login end.")
	fmt.Println("start Test TestSellToOtherAddr().")
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestSellToOtherAddr() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestSellToOtherAddr() login err=", err)
	}
	userNftCnt := strconv.Itoa(testCount)
	userNft, err := QueryUserNFTList(BatchBuyAddr, userNftCnt, "0", batchToken, batchkey)
	if err != nil {
		fmt.Println("TestSellToOtherAddr() QueryUserNFTList() err=", err)
		return
	}
	wd := sync.WaitGroup{}
	for _, nft := range userNft {
		wd.Add(1)
		go func(nft UserNft) {
			defer wd.Done()
			err = Sell(BatchBuyAddr, "FixPrice", nft.NftContractAddr, nft.NftTokenId, SellPrice, batchToken, batchkey, batchkey)
			if err != nil {
				fmt.Println("TestSellToOtherAddr() Sell() err=", err, "NftTokenId=", nft.NftTokenId)
				return
			}
			fmt.Println("TestSellToOtherAddr() OK.", "NftTokenId=", nft.NftTokenId)
		}(nft)
	}
	wd.Wait()
	for i, nft := range userNft {
		wd.Add(1)
		go func(i int, nft UserNft) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			nftInfo, err := QueryNFT(nft.NftContractAddr, nft.NftTokenId, batchToken)
			if err != nil {
				fmt.Println("TestSellToOtherAddr() QueryNFT  err=", err, "userAddr=", userAddr)
				return
			}
			//r := strconv.Itoa(nftInfo.Royalty)
			p := strconv.FormatUint(nftInfo.Auction.Startprice, 10) + "000000000"
			t, _ := hexutil.Decode(nftInfo.Auction.Tradesig)
			//mint, _ := hexutil.Decode(nftInfo.Approve)
			buyAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey)
			tx, err := Buy(nftInfo.NftContractAddr,  nftInfo.OwnerAddr, buyAddr.String(), nftInfo.NftTokenId,
				p, t, tKey[i].LogKey)
			if err != nil {
				fmt.Println("TestSellToOtherAddr() Buy() err=", err, "userAddr=", userAddr)
				return
			}
			fmt.Println("TestSellToOtherAddr() Buy() OK.", "userAddr=", userAddr)
			fmt.Println("Sell() tx.hash=", tx.Hash())
		}(i, nft)
	}
	wd.Wait()
	fmt.Println("end test TestSellToOtherAddr().")
}

func TestCancelSell(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestCancelSell() err=", err)
		return
	}
	fmt.Println("TestCancelSell() login end.")
	fmt.Println("start Test TestCancelSell.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			userNft, err := QueryUserNFTList(userAddr, strconv.Itoa(testCount), "0", tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestCancelSell() QueryUserNFTList() err=", err, "userAddr=", userAddr)
				return
			}
			err = CancelSell(userAddr, userNft[0].NftContractAddr, userNft[0].NftTokenId, tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("TestCancelSell()  err=", err, "userAddr=", userAddr)
				return
			}
			fmt.Println("TestCancelSell()  OK.", err, "userAddr=", userAddr)
		}(i)
	}
	wd.Wait()
	batchkey, err := crypto.HexToECDSA(BatchBuyPrv)
	if err != nil {
		fmt.Println("TestSellToOtherAddr() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestCancelSell() login err=", err)
		return
	}
	userNftCnt := strconv.Itoa(testCount)
	userNft, err := QueryUserNFTList(BatchBuyAddr, userNftCnt, "0", batchToken, batchkey)
	if err != nil {
		fmt.Println("TestCancelSell() QueryUserNFTList() err=", err)
		return
	}
	wd = sync.WaitGroup{}
	for _, nft := range userNft {
		wd.Add(1)
		go func(nft UserNft) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
			err = CancelSell(userAddr, nft.NftContractAddr, nft.NftTokenId, batchToken, batchkey)
			if err != nil {
				fmt.Println("TestCancelSell()  err=", err, "userAddr=", nft.OwnerAddr)
				return
			}
			fmt.Println("TestCancelSell()  OK.", err, "userAddr=", nft.OwnerAddr)
		}(nft)
	}
	wd.Wait()
	fmt.Println("end test TestCancelSell().")
}

func TestWormsSellSingle(t *testing.T) {
	batchkey, err := crypto.HexToECDSA("2ABE62D35B09680F007B225C318D5A672CA3E956B91BEEE5A5BA004A22DAAC2C")
	if err != nil {
		fmt.Println("TestWormsSellSingle() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsSellSingle() login err=", err)
		return
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	userNft, err := QueryUserNFTList(useraddr, "1", "0", batchToken, batchkey)
	if err != nil {
		fmt.Println("TestWormsSellSingle() QueryUserNFTList() err=", err)
		return
	}
	err = Sell(useraddr, "FixPrice", userNft[0].NftContractAddr, userNft[0].NftTokenId, SellPrice, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsSellSingle()  err=", err, "userAddr=", useraddr)
		//return
	}

	fmt.Println("end test TestWormsSellSingle().")
}

//from: 0x7fbc8ad616177c6519228fca4a7d9ec7d1804900
//to: 0xbd1cd0b483628c111b14d2be601f93e38655d94a
func TestWormsHighestMintMakeOffer(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer() err=", err)
		return
	}
	/*prv := crypto.FromECDSA(tKey[1].LogKey)
	fmt.Println(hexutil.Encode(prv))
	Addr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	fmt.Println(Addr)
	prv = crypto.FromECDSA(tKey[2].LogKey)
	fmt.Println(hexutil.Encode(prv))
	Addr = crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String()
	fmt.Println(Addr)*/

	fmt.Println("TestWormsHighestMintMakeOffer() login end.")
	fmt.Println("start Test TestWormsHighestMintMakeOffer.")
	batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer() login err=", err)
		return
	}
	useraddr := (crypto.PubkeyToAddress(batchkey.PublicKey).String())
	tradesign := `{"price":"0x174876e800","nft_address":"0x8000000000000000000000000000000000000001","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0xe","seller":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","sig":"2bed8c15cb15ddedd1a92774f6b6e8cf22b41f40bbb108146809bda3e33f958a71f348bbe40073c3c304ad924c3852ea73c7d9c76a0d912a32e850e7e034ab4100"}`
	err = WormSell(useraddr, "HighestBid", Nft1155Contract, BatchTokenId,
		SellPrice, tradesign, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String())
	tradesign = `{"price":"0x174876e800","nft_address":"0x8000000000000000000000000000000000000001","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x68","seller":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","sig":"b9827df881944f6816266b8329cf667388163f5eb5114b83fce0a3f2f5585b6a0fc11ff73e933f5dcbbdced9cb5dffae9e8c7dfc52ba08950e77734d39c3ab3f00"}`
	err = WormMakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tradesign, tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer() key err=", err)
		return
	}
	BiddingAddr = (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	err = WormMakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1002", tradesign,tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer() key err=", err)
		return
	}
	fmt.Println("end test TestWormsHighestMintMakeOffer().")
}

//from: 0x7fbc8ad616177c6519228fca4a7d9ec7d1804900
//to: 0xbd1cd0b483628c111b14d2be601f93e38655d94a
func TestWormsFixPricetMintBuy(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy() err=", err)
		return
	}

	fmt.Println("TestWormsFixPricetMintBuy() login end.")
	fmt.Println("start Test TestWormsFixPricetMintBuy.")
	batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy() login err=", err)
		return
	}
	useraddr := (crypto.PubkeyToAddress(batchkey.PublicKey).String())
	tradesign := `{"price":"0x174876e800","nft_address":"0x8000000000000000000000000000000000000001","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0xe","seller":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","sig":"2bed8c15cb15ddedd1a92774f6b6e8cf22b41f40bbb108146809bda3e33f958a71f348bbe40073c3c304ad924c3852ea73c7d9c76a0d912a32e850e7e034ab4100"}`
	err = WormSell(useraddr, "FixPrice", Nft1155Contract, BatchTokenId,
		SellPrice, tradesign, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String())
	tradesign = `{"price":"0x174876e800","nft_address":"0x8000000000000000000000000000000000000001","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x68","seller":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","sig":"b9827df881944f6816266b8329cf667388163f5eb5114b83fce0a3f2f5585b6a0fc11ff73e933f5dcbbdced9cb5dffae9e8c7dfc52ba08950e77734d39c3ab3f00"}`
	//WormBuying(userAddr, buyer_Addr, buyer_sig, contract, tokenId, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = WormBuying(BiddingAddr, BiddingAddr, tradesign, Nft1155Contract,
		BatchTokenId, "1001", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetMintBuy().")
}

func TestWormsFixPricetNotMintBuy(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintBuy() err=", err)
		return
	}
	fmt.Println("TestWormsFixPricetNotMintBuy() login end.")
	fmt.Println("start Test TestWormsFixPricetMintBuy.")

	useraddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	tradesign := `{"price":"0x174876e800","royalty":"200","meta_url":"7b226d657461223a222f697066732f69706673516d534150326575794546446b694b5a446258347a754c556132576275577963586f6773513452734156445a446d222c22746f6b656e5f6964223a2231343632323830353233353730227d","exclusive_flag":"1","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x16","sig":"0x37514c7955f52ea21ea69b39303751cd2a70989a10cd436b43ba55093ab0a1545387223b1b84bc2e30767209394a5563ab68f384504eddd3dcc94a708ae1937c1c"}`
	err = WormSell(useraddr, "FixPrice", Nft1155Contract, "1462280523570",
		SellPrice, tradesign, tokens[1], tKey[1].WorkKey, tKey[1].LogKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy()  err=", err, "userAddr=", useraddr)
		//return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	tradesign = `{"price":"0x174876e800","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x16","sig":"0xefe446c0963c08a19882f7cd9f8be815fe527770fb4673646618166adb741ca32a1cba787d67d09350cba71c3d9fb0589084041c3fd35ada75f813ca4ad2381d1b"}`
	//WormBuying(userAddr, buyer_Addr, buyer_sig, contract, tokenId, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = WormBuying(BiddingAddr, BiddingAddr, tradesign, Nft1155Contract,
		"1462280523570", "1001", tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintBuy() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetNotMintBuy().")
}

func TestWormsAuthExFixPricetNotMintBuy(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintBuy() err=", err)
		return
	}
	fmt.Println("TestWormsFixPricetNotMintBuy() login end.")
	fmt.Println("start Test TestWormsFixPricetMintBuy.")

	useraddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	tradesign := `{"price":"0x174876e800","royalty":"200","meta_url":"7b226d657461223a222f697066732f69706673516d534150326575794546446b694b5a446258347a754c556132576275577963586f6773513452734156445a446d222c22746f6b656e5f6964223a2231343632323830353233353730227d","exclusive_flag":"1","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x16","sig":"0x37514c7955f52ea21ea69b39303751cd2a70989a10cd436b43ba55093ab0a1545387223b1b84bc2e30767209394a5563ab68f384504eddd3dcc94a708ae1937c1c"}`
	err = WormSell(useraddr, "FixPrice", Nft1155Contract, "1462280523570",
		SellPrice, tradesign, tokens[1], tKey[1].WorkKey, tKey[1].LogKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy()  err=", err, "userAddr=", useraddr)
		//return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	tradesign = `{"price":"0x174876e800","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x16","sig":"0xefe446c0963c08a19882f7cd9f8be815fe527770fb4673646618166adb741ca32a1cba787d67d09350cba71c3d9fb0589084041c3fd35ada75f813ca4ad2381d1b"}`
	//WormBuying(userAddr, buyer_Addr, buyer_sig, contract, tokenId, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = WormBuying(BiddingAddr, BiddingAddr, tradesign, Nft1155Contract,
		"1462280523570", "1001", tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintBuy() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetNotMintBuy().")
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func EthSign(msg string, prv *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(signHash([]byte(msg)), prv)
	if err != nil {
		fmt.Println("EthSign() err=", err)
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

func TestWormsFixPricetNotMintBuySign(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintBuy() err=", err)
		return
	}
	fmt.Println("TestWormsFixPricetNotMintBuy() login end.")
	fmt.Println("start Test TestWormsFixPricetMintBuy.")
	nftmeta := contracts.NftMeta{Meta: "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm",
		TokenId: "1462280523570",}
	nftmetastr, _ := json.Marshal(&nftmeta)
	metastr := hex.EncodeToString(nftmetastr)
	seller := contracts.Seller2{
		Price:"100000000000",
		Royalty:"200",
		Metaurl:metastr,
		Exclusiveflag: "1",
		Exchanger:"0xa1e67a33e090afe696d7317e05c506d7687bb2e5",
		Sig:"",
	}
	p, _ := strconv.ParseUint(seller.Price, 10, 64)
	r, _ := strconv.ParseUint(seller.Royalty, 10, 64)
	seller.Price	= hexutil.EncodeUint64(p)
	seller.Royalty = hexutil.EncodeUint64(r)
	seller.Metaurl = seller.Metaurl
	seller.Exclusiveflag = "1"
	seller.Exchanger = seller.Exchanger
	seller.Blocknumber = hexutil.EncodeUint64(20)
	seller.Sig = seller.Sig
	Sellprv, err := crypto.HexToECDSA("273c6cfec83a9d35d808bae69d7e67d5e21a3121c3fa58f03155413475a6e36e")
	if err != nil {
		fmt.Println("SendFixTrans() HexToECDSA err=", err)
	}
	msg := seller.Price + seller.Royalty + seller.Metaurl +
		seller.Exclusiveflag + seller.Exchanger + seller.Blocknumber
	sig, err := EthSign(msg, Sellprv)
	if err != nil {
		fmt.Println("ExchangeTrans() Sign err=", err)
	}
	seller.Sig = hexutil.Encode(sig)
	tradesign, _  := json.Marshal(&seller)
	useraddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	//tradesign := `{"price":"0x174876e800","royalty":"200","meta_url":"7b226d657461223a222f697066732f69706673516d534150326575794546446b694b5a446258347a754c556132576275577963586f6773513452734156445a446d222c22746f6b656e5f6964223a2231343632323830353233353730227d","exclusive_flag":"1","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x16","sig":"0x37514c7955f52ea21ea69b39303751cd2a70989a10cd436b43ba55093ab0a1545387223b1b84bc2e30767209394a5563ab68f384504eddd3dcc94a708ae1937c1c"}`
	err = WormSell(useraddr, "FixPrice", Nft1155Contract, "1462280523570",
		SellPrice, string(tradesign), tokens[1], tKey[1].WorkKey, tKey[1].LogKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuy()  err=", err, "userAddr=", useraddr)
		//return
	}
	buyer := contracts.Buyer1{
		Price:"100000000000",
		//Nftaddress:"0x8000000000000000000000000000000000000001",
		Exchanger:"0xa1e67a33e090afe696d7317e05c506d7687bb2e5",
		Blocknumber: "1",
		Sig:"",
	}
	p, _ = strconv.ParseUint("100000000000", 10, 64)
	buyer.Price = hexutil.EncodeUint64(p)
	buyer.Exchanger = buyer.Exchanger
	buyer.Blocknumber = hexutil.EncodeUint64(10)
	Buyprv, err := crypto.HexToECDSA("accab7212dfbe235817bed4c1052c40fc0978e0a0b94feecc6df9cf5f08d3485")
	if err != nil {
		fmt.Println("SendFixTrans() HexToECDSA err=", err)
	}
	msg = buyer.Price + buyer.Exchanger + buyer.Blocknumber
	sig, err = EthSign(msg, Buyprv)
	if err != nil {
		fmt.Println("ExchangeTrans() Sign err=", err)
	}
	buyer.Sig = hexutil.Encode(sig)
	tradesign, _ = json.Marshal(&buyer)
	BiddingAddr := (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	//tradesign = `{"price":"0x174876e800","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x16","sig":"0xefe446c0963c08a19882f7cd9f8be815fe527770fb4673646618166adb741ca32a1cba787d67d09350cba71c3d9fb0589084041c3fd35ada75f813ca4ad2381d1b"}`
	err = WormBuying(BiddingAddr, BiddingAddr, string(tradesign), Nft1155Contract,
		"1462280523570", "1001", tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintBuy() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetNotMintBuy().")
}

func TestWormsFixPricetNotMintSell(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintSell() err=", err)
		return
	}
	fmt.Println("TestWormsFixPricetNotMintSell() login end.")
	fmt.Println("start Test TestWormsFixPricetNotMintSell.")

	selladdr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	tradesign := `{"price":"0x174876e800","royalty":"200","meta_url":"{\"meta\":\"/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm\",\"token_id\":\"4285206002595\"}","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x659","sig":"83da8a793cb976e59323ac1086b78ef14261a1e9e8b31b49bf9c23a226bf0fe046f85152168b58751da9f3e3de46c7b091e20b7d733c9b7b8eee9577de710c3400"}`
	err = WormSell(selladdr, "FixPrice", Nft1155Contract, "4285206002595",
		SellPrice, tradesign, tokens[1], tKey[1].WorkKey, tKey[1].LogKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintSell()  err=", err, "userAddr=", selladdr)
		return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	tradesign = `{"price":"0x174876e800","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x659","sig":"d51a89b11a7fa3a5e1f16f91a16a8331be41235da344e20a857e8768ec13714f597e6e90dc12806b296d0e90a8a3c124398727d7ab3fa6344907e36b0ffe85a800"}`
	BiddingAddr = (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	err = WormMakeOffer(BiddingAddr, Nft1155Contract,
		"4285206002595", "1002", tradesign,tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintSell() key err=", err)
		return
	}
	err = WormBuying(selladdr, BiddingAddr, tradesign, Nft1155Contract,
		"4285206002595", "1001", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetNotMintSell() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetNotMintSell().")
}

//from: 0x7fbc8ad616177c6519228fca4a7d9ec7d1804900
//to: 0xbd1cd0b483628c111b14d2be601f93e38655d94a
func TestWormsFixPricetMintSell(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuyer() err=", err)
		return
	}

	fmt.Println("TestWormsFixPricetMintBuyer() login end.")
	fmt.Println("start Test TestWormsFixPricetMintBuyer.")
	batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuyer() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuyer() login err=", err)
		return
	}
	useraddr := (crypto.PubkeyToAddress(batchkey.PublicKey).String())
	tradesign := `{"price":"0x174876e800","nft_address":"0x8000000000000000000000000000000000000001","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0xe","seller":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","sig":"2bed8c15cb15ddedd1a92774f6b6e8cf22b41f40bbb108146809bda3e33f958a71f348bbe40073c3c304ad924c3852ea73c7d9c76a0d912a32e850e7e034ab4100"}`
	err = WormSell(useraddr, "FixPrice", Nft1155Contract, BatchTokenId,
		SellPrice, tradesign, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuyer()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String())
	err = WormMakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tradesign, tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsHighestMintMakeOffer() key err=", err)
		return
	}
	tradesign = `{"price":"0x174876e800","nft_address":"0x8000000000000000000000000000000000000001","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x68","seller":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","sig":"b9827df881944f6816266b8329cf667388163f5eb5114b83fce0a3f2f5585b6a0fc11ff73e933f5dcbbdced9cb5dffae9e8c7dfc52ba08950e77734d39c3ab3f00"}`
	err = WormBuying(useraddr, BiddingAddr, tradesign, Nft1155Contract,
		BatchTokenId, "1001", batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMintBuyer() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetMintBuyer().")
}
//from:0x2a95249bcbe73397f54562ff7a74d40b9d34a08b
//to:0x79fb137312b1c6d60044b1f06b80a17609bc05a0
func TestWormsHighestNotMintMakeOffer(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsHighestNotMintMakeOffer() err=", err)
		return
	}
	fmt.Println("TestWormsHighestNotMintMakeOffer() login end.")
	fmt.Println("start Test TestWormsHighestNotMintMakeOffer.")
	//batchkey, err := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
	batchkey, err := crypto.HexToECDSA("7EFB556A977824B34D45B7FC8A265ED5AA6D6C21A90E9BFA5C021FC23795594C")
	if err != nil {
		fmt.Println("TestWormsHighestNotMintMakeOffer() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsHighestNotMintMakeOffer() login err=", err)
		return
	}
	useraddr := (crypto.PubkeyToAddress(batchkey.PublicKey).String())
	tradesign := `{"price":"0x174876e800","royalty":"200","meta_url":"{\"meta\":\"/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm\",\"token_id\":\"2597785300040\"}","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x10","sig":"c79737563e0465d87d50230781248a1e49bee52b03cf8cc5bffb5a35067c13bf5c97ba07f8e906f64948b00b59c55780d3de14aff1ecbfe138de30a8ed8997b401"}`
	err = WormSell(useraddr, "HighestBid", Nft1155Contract, BatchTokenId,
		SellPrice, tradesign, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsHighestNotMintMakeOffer()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := (crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String())
	tradesign = `{"price":"0x174876e800","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","blocknumber":"0x10","sig":"bdd9aed5e753ad379e5899a87720e217b9808b11ade89a362312cd90387faa35752ab22b29eb1fea818369315043a6cbdcb9cbbfb6f29f253cc926ba50e6878f00"}`
	err = WormMakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tradesign, tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsHighestNotMintMakeOffer() key err=", err)
		return
	}
	BiddingAddr = (crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String())
	err = WormMakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1002", tradesign,tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsHighestNotMintMakeOffer() key err=", err)
		return
	}
	fmt.Println("end test TestWormsHighestNotMintMakeOffer().")
}

func TestWormsFixPricetMakeOffer(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsFixPricetMakeOffer() err=", err)
		return
	}
	fmt.Println("TestWormsFixPricetMakeOffer() login end.")
	fmt.Println("start Test TestWormsFixPricetMakeOffer.")
	batchkey, err := crypto.HexToECDSA("2ABE62D35B09680F007B225C318D5A672CA3E956B91BEEE5A5BA004A22DAAC2C")
	if err != nil {
		fmt.Println("TestWormsFixPricetMakeOffer() key err=", err)
		return
	}
	batchToken, err := Login(batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMakeOffer() login err=", err)
		return
	}
	useraddr := crypto.PubkeyToAddress(batchkey.PublicKey).String()
	err = Sell(useraddr, "FixPrice", Nft1155Contract, BatchTokenId,
		SellPrice, batchToken, batchkey, batchkey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMakeOffer()  err=", err, "userAddr=", useraddr)
		return
	}
	BiddingAddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	//MakeOffer(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMakeOffer() key err=", err)
		return
	}
	BiddingAddr = crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String()
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1002", tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsFixPricetMakeOffer() key err=", err)
		return
	}
	fmt.Println("end test TestWormsFixPricetMakeOffer().")
}

func TestWormsBidPricetMakeOffer(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsBidPricetMakeOffer() err=", err)
		return
	}
	fmt.Println("TestWormsBidPricetMakeOffer() login end.")
	fmt.Println("start Test TestWormsBidPricetMakeOffer.")

	BiddingAddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	//MakeOffer(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsBidPricetMakeOffer() key err=", err)
		return
	}
	BiddingAddr = crypto.PubkeyToAddress(tKey[2].LogKey.PublicKey).String()
	err = MakeOffer(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1002", tokens[2], tKey[2].LogKey, tKey[2].WorkKey)
	if err != nil {
		fmt.Println("TestWormsBidPricetMakeOffer() key err=", err)
		return
	}
	fmt.Println("end test TestWormsBidPricetMakeOffer().")
}

func TestWormsBuying(t *testing.T) {
	testCount := 10
	tKey, tokens, err := Mlogin(testCount)
	if err != nil {
		fmt.Println("TestWormsBuying() err=", err)
		return
	}
	fmt.Println("TestWormsBuying() login end.")
	fmt.Println("start Test TestWormsBuying.")

	BiddingAddr := crypto.PubkeyToAddress(tKey[1].LogKey.PublicKey).String()
	//MakeOffer(userAddr string, contract string, tokenId string, price string, token string, userKey, workKey *ecdsa.PrivateKey) error {
	err = Buying(BiddingAddr, Nft1155Contract,
		BatchTokenId, "1001", tokens[1], tKey[1].LogKey, tKey[1].WorkKey)
	if err != nil {
		fmt.Println("TestWormsBuying() key err=", err)
		return
	}
	fmt.Println("end test TestWormsBuying().")
}
