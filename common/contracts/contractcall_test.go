package contracts

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"testing"
)

func init() {
	//EthNode = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	EthNode = "https://api.wormholestest.com"
	Weth9Addr = "0xf4bb2e28688e89fcce3c0580d37d36a7672e8a9f"
	TradeAddr = "0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"
	Nft1155Addr = "0xa1e67a33e090afe696d7317e05c506d7687bb2e5"
	//SuperAdminAddr = "2ABE62D35B09680F007B225C318D5A672CA3E956B91BEEE5A5BA004A22DAAC2C"
	SuperAdminAddr = "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
}

func TestGetBlockTxs(t *testing.T) {
	_, _, err := GetBlockTxs(2006)
	if err != nil {
		fmt.Println("GetBlockTxs error.")
	}
	//_, _, err = GetBlockTxs(9766060)
	//if err != nil {
	//	fmt.Println("GetBlockTxs error.")
	//}
	//_, _, err = GetBlockTxs(9508909)
	//if err != nil {
	//	fmt.Println("GetBlockTxs error.")
	//}
	//_, _, err = GetBlockTxs(9766060)
	//if err != nil {
	//	fmt.Println("GetBlockTxs error.")
	//}
}

func TestGetBlockTxsNew(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, err := GetBlockTxsNew(184926 + uint64(i))
		if err != nil {
			fmt.Println("GetBlockTxs error.")
		}
	}

}

func TestAuctionAndMint(t *testing.T) {
	//AuctionAndMint(from, to, nftAddr, tokenId, price, amount, royaltyRatio, tokenURI, sig string) (string, error) {
	_, _ = AuctionAndMint("0xc9a9caa0147adc101138920ac7905ca6b62e9a2a", "0x5e83e4c3bc80769b4d67fc4cb577b352c7b658bf",
		"0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "9529124519747", "100000000000",
		"1000", "200", "", "")

}

func TestIsApproveOwn(t *testing.T) {
	EthNode = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	Weth9Addr = "0xf4bb2e28688e89fcce3c0580d37d36a7672e8a9f"
	TradeAddr = "0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82"
	Nft1155Addr = "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5"
	//_, err := IsApprovedNFT1155("0xc9a9caa0147adc101138920ac7905ca6b62e9a2a", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5")
	//_, err := IsOwnerOfNFT1155("0x5e83e4c3bc80769b4d67fc4cb577b352c7b658bf", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "5201739360460")
	//_, err = IsErcNFT1155("0xa1e67a33e090afe696d7317e05c506d7687bb2e5")
	//b, err := OwnAndAprove("0xc9a9caa0147adc101138920ac7905ca6b62e9a2a", "0xa1e67a33e090afe696d7317e05c506d7687bb2e5", "5201739360460")
	//if err != nil {
	//	fmt.Println("OwnAndAprove() err=", err)
	//}
	//fmt.Println("OwnAndAprove() =", b)
}

func TestExchangerMint(t *testing.T) {
	/*seller := Seller{
		Price:"0x174876e800",
		Royalty:"200",
		Metaurl:"{\"meta\":\"/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm\",\"token_id\":\"2597785300040\"}",
		Exchanger:"0xa1e67a33e090afe696d7317e05c506d7687bb2e5",
		Blocknumber: "0x59",
		Sig:"3dd8dafd8ae007c373d41fb084c605bd1b0953e5c87a2eca329ccb69e019064c4afaa8cfb4e330a0da43e12a628ee394e4e9238d2bb50e8e885e163b3836a16f01",
	}*/
	//seller := Seller{}
	//sellerbyte := "{\"price\":\"0x174876e800\",\"royalty\":\"200\",\"meta_url\":\"{\\\"meta\\\":\\\"/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm\\\",\\\"token_id\\\":\\\"2597785300040\\\"}\",\"exchanger\":\"0xa1e67a33e090afe696d7317e05c506d7687bb2e5\",\"blocknumber\":\"0x59\",\"sig\":\"3dd8dafd8ae007c373d41fb084c605bd1b0953e5c87a2eca329ccb69e019064c4afaa8cfb4e330a0da43e12a628ee394e4e9238d2bb50e8e885e163b3836a16f01\"}"
	//err := json.Unmarshal([]byte(sellerbyte), &seller)
	//if err != nil {
	//	t.Fatal(err)
	//}

	/*buyer := Buyer{
		Price:"0x174876e800",
		Exchanger:"0xa1e67a33e090afe696d7317e05c506d7687bb2e5",
		Blocknumber: "0x59",
		Sig:"06207a84f457e28f7637df9be62ccb63dd7c9ac59fd1b0a8934e2467ca6bdb645a6fa48eb4ba65e27259c73c9f4cb248c02cae681f07fefb920e86d63d40b96501",
	}*/
	/*{
		seller := Seller{}
		sellerbyte := "{\"price\":\"0x174876e800\",\"royalty\":\"200\",\"meta_url\":\"{\\\"meta\\\":\\\"/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm\\\",\\\"token_id\\\":\\\"2597785300040\\\"}\",\"exchanger\":\"0xa1e67a33e090afe696d7317e05c506d7687bb2e5\",\"blocknumber\":\"0x59\",\"sig\":\"3dd8dafd8ae007c373d41fb084c605bd1b0953e5c87a2eca329ccb69e019064c4afaa8cfb4e330a0da43e12a628ee394e4e9238d2bb50e8e885e163b3836a16f01\"}"
		err := json.Unmarshal([]byte(sellerbyte), &seller)
		if err != nil {
			t.Fatal(err)
		}
		buyer := Buyer{}
		buybyte := "{\"price\":\"0x174876e800\",\"exchanger\":\"0xa1e67a33e090afe696d7317e05c506d7687bb2e5\",\"blocknumber\":\"0x59\",\"sig\":\"06207a84f457e28f7637df9be62ccb63dd7c9ac59fd1b0a8934e2467ca6bdb645a6fa48eb4ba65e27259c73c9f4cb248c02cae681f07fefb920e86d63d40b96501\"}"
		err = json.Unmarshal([]byte(buybyte), &buyer)
		if err != nil {
			t.Fatal(err)
		}
		err = ExchangerMint(WormHolesExMintTransfer, WormHolesVerseion, seller, buyer, SuperAdminAddr)
		if err != nil {
			t.Fatal(err)
		}
	}*/

	{
		buyer := Buyer{}
		buybyte := "{\"price\":\"0x174876e800\",\"nft_address\":\"0x8000000000000000000000000000000000000001\",\"exchanger\":\"0xa1e67a33e090afe696d7317e05c506d7687bb2e5\",\"blocknumber\":\"0x85\",\"seller\":\"0xcf146d17a7086d0eb23faaa455b2418906ee5169\",\"sig\":\"6588def6ba887116dd426174e8f43a00898b4885a4bbd5539111cf93591a6f3f138f3f3c7a97690ea8319e23da53beb4ba39afa6de5c929d0dd3d5bdb6941a5201\"}"
		err := json.Unmarshal([]byte(buybyte), &buyer)
		if err != nil {
			t.Fatal(err)
		}
		err = ExchangeTrans(buyer, SuperAdminAddr)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestUnmarshal(t *testing.T) {
	data := []byte(`wormholes:{"version":"v0.0.1","type":16,"seller2":{"price":"0xf4240","royalty":"0x64","meta_url":{"meta":"/ipfs/QmPQVzUcgnGiPx6wQLkkjkHkxcZw2oUStRWmfwiW1Tw3qg","token_id":"9580750891422"},"exclusive_flag":"1","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0x3b9aca00","sig":"0x05636fa9a40ae621c2d987849787b805a91914ca17428cf31d87c15963da2fbe5d554d68a7d88cb0e97b623b065cbfbb009d5838a015a26968d6c3a187d84c591b"}}`)
	wormMint := WormholesMint{}
	jsonErr := json.Unmarshal(data[10:], &wormMint)
	if jsonErr != nil {
		t.Fatal("GetBlockTxs() wormholes mint type err=")
	}
	metastr := hex.EncodeToString([]byte(`{"meta":"/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm","token_id":"4285206002595"}`))
	fmt.Println(metastr)
	metabyte, _ := hex.DecodeString(metastr)
	fmt.Println(metabyte)
	var nftmeta NftMeta
	jsonErr = json.Unmarshal(metabyte, &nftmeta)
	if jsonErr != nil {
		fmt.Println("GetBlockTxs() NftMeta unmarshal type err=")
	}
}

func TestUmarshalAuth(t *testing.T) {
	var authTrans ExchangerAuth
	//authSign := `{"exchanger_owner":"0xa9ff5a316c84c40796a234a87f04a0d7989414f1","to":"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900","block_number":"0xb530","sig":"0x41ef367b0bde7a02d41cf435570aa5f067943b98dfe35d01c26a1470d6d25e32255a6c4222122833a25ba155a3d9bf267cec9c8cf658e13026f477edfdb978d01c"}`
	authSign := `{"exchanger_owner":"0x836Acd33188A8C2e2c94C112FD15657EEBFA223c","to":"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900","block_number":"0xbbcc","sig":"0x48f9f5c83d4ef8616d82b375c87df63e012d65b5969b73f4c373daeba0fdb4b7796019a2b3d644823c85c67777dac8d4d63a8340174cd039c8ce162c2fbef43b1c"}`
	//authSign = strings.Replace(authSign, "\\", "", -1)
	err := json.Unmarshal([]byte(authSign), &authTrans)
	if err != nil {
		log.Println("AuthExchangerMint()  Unmarshal() err=", err)
		return
	}
	msg := authTrans.Exchangerowner + authTrans.To + authTrans.Blocknumber
	toAddress, err := recoverAddress(msg, authTrans.Sig)
	if err != nil {
		log.Println("AuthExchangerMint() recoverAddress() err=", err)
		return
	}
	fmt.Println("toAddress=", toAddress)
}

func TestRecover(t *testing.T) {
	msg := `{"exchangerauth":"{\"exchanger_owner\":\"0x4571022790985add90f84e2ad96a8fc01569043e\",\"to\":\"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900\",\"block_number\":\"0xba5b\",\"sig\":\"0x7a02798034adc7182ca797d3eb9706b254d14ac0c43fd92a5ef178498138c910312bd8226215c9c9832fb2fb7c165f2ebcd2cc07b481fd8571d96467a6edce7e1c\"}"`
	Sig := `0xd5085bec920775c3b94f2fbbec2205b7753d44958c979be93f8f90d52a78296f11b34e28459f1e405c27830ec68fdfdd8d81d02e1a076c909c727dd1638f678a1c`
	toAddress, err := recoverAddress(msg, Sig)
	if err != nil {
		log.Println("AuthExchangerMint() recoverAddress() err=", err)
		return
	}
	fmt.Println("toAddress=", toAddress)
}

func TestUnMarshalSeller(t *testing.T) {
	seller := Seller2{}
	sellerSig := `{"price":"0xf4240","royalty":"0x64","meta_url":"7b226d657461223a222f697066732f516d52357756784a467574685577626d77585a5570394a6b5069314c5051644a596a716834754a77446d6d5a6758222c22746f6b656e5f6964223a2235323036333831343038313133227d","exclusive_flag":"1","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0xae24","sig":"0x6170e56c83bb2324d608af2036dd7cdb329982d4a1cb38e2b3221625af4896c8233619c8cbb6f8f3f5b2c1a13a7f21b942e40658d31c5713a43459382b9007151b"}`
	err := json.Unmarshal([]byte(sellerSig), &seller)
	if err != nil {
		t.Fatal("WormTrans() Unmarshal() err=", err)
	}
	var nftmeta NftMeta
	metabyte, jsonErr := hex.DecodeString(seller.Metaurl)
	if jsonErr != nil {
		fmt.Println("GetBlockTxs() hex.DecodeString err=", err)
	}
	jsonErr = json.Unmarshal(metabyte, &nftmeta)
	if jsonErr != nil {
		fmt.Println("GetBlockTxs() NftMeta unmarshal type err=", err)
	}
	msg := seller.Price + seller.Royalty + seller.Metaurl + seller.Exclusiveflag +
		seller.Exchanger + seller.Blocknumber

	fromAddr, err := recoverAddress(msg, seller.Sig)
	fmt.Println("toaddr=", fromAddr.String())
}

//func signHash(data []byte) []byte {
//	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
//	return crypto.Keccak256([]byte(msg))
//}

func EthSign(msg string, prv *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(signHash([]byte(msg)), prv)
	if err != nil {
		fmt.Println("EthSign() err=", err)
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

/*
v{"price":"0xf4240","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0xc6bf","sig":"0xb4f1d68a664ea0cc04fcec4a7ee2f2b43d9ef59fb3f7cf78cea9b35c78b838462e571727d64e60fa69f21fc077fec9ebd362fd3bb2c2ee584acc86b57f4478851b"}
{"price":"0xf4240","exchanger":"0xa1e67a33e090afe696d7317e05c506d7687bb2e5","block_number":"0xc68d","sig":"0xd58c4bbdc432be5bec753c211b11b79eb1fd3b7d510e2af7f7333f95b2a82a951160b59609753247cae1b57c823f9b22909ac1c23dcddc4b694bf25e3e9c54361c"}
*/

func TestUnMarshalBuyer1(t *testing.T) {
	buyerSig := "{\\\"price\\\":\\\"0xf4240\\\",\\\"exchanger\\\":\\\"0xa1e67a33e090afe696d7317e05c506d7687bb2e5\\\",\\\"block_number\\\":\\\"0xc713\\\",\\\"sig\\\":\\\"0x5b576aaf0bbd2135d9e98d44d5600c4f623dfa494da27e963016ee7d2891d02b35abb019fbd932104347e27a00a74267f8fc5317d64f6151f984a6efbad504b01b\\\"}"
	buyer := Buyer1{}
	err := json.Unmarshal([]byte(buyerSig), &buyer)
	if err != nil {
		t.Fatal("WormTrans() Unmarshal() err=", err)
	}
	msg := buyer.Price + buyer.Exchanger + buyer.Blocknumber
	toaddr, err := recoverAddress(msg, buyer.Sig)
	if err != nil {
		t.Fatal("GetBlockTxs() recoverAddress() err=", err)
	}
	fmt.Println("toaddr=", toaddr.String())
}

func TestUnMarshalBuyer(t *testing.T) {
	p, _ := hexutil.DecodeUint64("0x174876e800")
	fmt.Println(p)
	buyerSig := "{\"price\":\"0xf4240\",\"nft_address\":\"0x8000000000000000000000000000000000000001\",\"exchanger\":\"0xa1e67a33e090afe696d7317e05c506d7687bb2e5\",\"block_number\":\"0xc14b\",\"seller\":\"0x49b89e4fe5404ccecb7c6095032fe2ec94c4e1e7\",\"sig\":\"0xfe67a06aec03262da41b5bc3ee9cb6ed2afc0032b6ce148bb69afa306c955ecc247e1b423677de72f3e402e52d1ebdc9e2ed1f64097aedee169b37c2c238906e1c\"}"
	buyer := Buyer{}
	err := json.Unmarshal([]byte(buyerSig), &buyer)
	if err != nil {
		t.Fatal("WormTrans() Unmarshal() err=", err)
	}
	msg := buyer.Price + buyer.Exchanger + buyer.Blocknumber
	toaddr, err := recoverAddress(msg, buyer.Sig)
	if err != nil {
		t.Fatal("GetBlockTxs() recoverAddress() err=", err)
	}
	fmt.Println("toaddr=", toaddr.String())
}

func TestForceBuy(t *testing.T) {
	EthNode = "http://192.168.4.240:8561"
	buyer := Buyer{}
	buyer.Nftaddress = "0x80000000000000000000000000000000000002f"
	buyer.Exchanger = "0x0109CC44df1C9ae44Bac132eD96f146Da9A26B88"
	buyer.Seller = "0x68B14e0F18C3EE322d3e613fF63B87E56D86Df60"
	//buyerAuth := `{"exchanger":"0x0109cc44df1c9ae44bac132ed96f146da9a26b88","block_number":"0x7cc2","sig":"0x8cf32ee8ce40862ecaf4c658d7643dce73e3cba23d0998b611ebcdc587734b3463721d3e794cdada475fad4d14f2e343f22e94c33de563a7a4f5554138b279dd1c"}`
	buyerAuth := `{"exchanger":"0x0109cc44df1c9ae44bac132ed96f146da9a26b88","block_number":"0x603f8d","sig":"0x762923206bfa359f493b6332f8e8aab183d3f4378201f419038c5904d8fda396079736580d14c465b6c07f6b305a1c98c8da8b481a63a3ec67634d553eb5c98e1c"}`
	exchangeAuth := `{"exchanger_owner":"0x0109CC44df1C9ae44Bac132eD96f146Da9A26B88","to":"0x7fbc8ad616177c6519228fca4a7d9ec7d1804900","block_number":"0x2540be400","sig":"0x6f7508e28d3479326926c62ab3963f0efbfb6b24c9899af9e607a8a5465a4ac3590fa6d0c418ccc2fcc2d3e4d9bb687a1d9e5b201414fd60e1636cacc9aeef811c"}`
	fromprv := "501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d"
	txHash, blockn, err := ForceBuyingAuthExchangeTrans(buyer, buyerAuth, exchangeAuth, fromprv)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(txHash, blockn)
}

func TestDecodeHex(t *testing.T) {
	nftmeta := NftMeta{Meta: "/ipfs/ipfsQmSAP2euyEFDkiKZDbX4zuLUa2WbuWycXogsQ4RsAVDZDm",
		TokenId: "1462280523570"}
	nftmetastr, _ := json.Marshal(&nftmeta)
	metastr := hex.EncodeToString(nftmetastr)
	metastr = hex.EncodeToString([]byte(`wormholes:{"type":2,"nft_address":"0x800000000000000000000000000000000000000d","version":"v0.0.1"}`))
	metastr = "776f726d686f6c65733a7b2274797065223a322c226e66745f61646472657373223a22307838303030303030303030303030303030303030303030303030303030303030303030303030303064222c2276657273696f6e223a2276302e302e31227d"
	mstr, err := hex.DecodeString(metastr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mstr)
}

func TestPrice(t *testing.T) {
	sp, _ := hexutil.DecodeUint64("0x8ac7230489e80000")
	fmt.Println(sp)
	value := big.NewInt(int64(sp))
	fmt.Println(value)
	valu1, _ := hexutil.DecodeBig("0x8ac7230489e80000")
	fmt.Println(valu1.String())

}

func TestGetSnftInfo(t *testing.T) {
	//EthNode = "https://api.wormholestest.com"
	EthNode = "http://192.168.4.240:8561"
	//snftAddr, err := GetSnftAddressList(big.NewInt(53668), true)
	//if err != nil {
	//	t.Fatal("GetSnftInfo() err=", err)
	//}
	//fmt.Println(snftAddr)
	addr := common.HexToAddress("0x8000000000000000000000000000000000000390")
	snftInfo, err := GetAccountInfo(addr, big.NewInt(926))
	fmt.Println(snftInfo, err)
}

func TestGetNominatedNFTInfo(t *testing.T) {
	EthNode = "http://api.wormholestest.com:8561"
	snftInfo, err := GetNominatedNFTInfo(big.NewInt(53870))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(snftInfo)
	fmt.Println(snftInfo.StartIndex / snftInfo.Number)
}
