package sync

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestGetNftSysMintInfo(t *testing.T) {
	contracts.BrowseNode = "http://192.168.1.237:8090"
	models.NftIpfsServerIP = "http://api.wormholestest.com"
	models.NftstIpfsServerPort = "8666"
	//7250
	blocknu := uint64(7218)
	minfo, err := GetNftSysMintInfo(blocknu)
	if err != nil {
		t.Fatal("err= ", err)
	}
	fmt.Println(minfo)
}

func TestGetSnftInfo(t *testing.T) {
	contracts.EthNode = "http://api.wormholestest.com:8561"
	models.NftIpfsServerIP = "http://api.wormholestest.com"
	models.NftstIpfsServerPort = "8666"
	models.RoyaltyLimit = 1000
	SyncWorkerNft(sqldsnT)
}

func TestGetIpfsInfo(t *testing.T) {
	models.NftIpfsServerIP = "http://api.wormholestest.com"
	models.NftstIpfsServerPort = "8668"
	snftinfo, err := GetSnftInfoFromIPFSWithShell("/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By/01")
	if err != nil {
		errflag := strings.Index(err.Error(), "context deadline exceeded")
		if errflag == -1 {
			fmt.Println(errflag)
		}
		t.Fatal(err)
	}
	fmt.Println(snftinfo)
}

func TestAddDirIpfs(t *testing.T) {
	models.NftIpfsServerIP = "http://api.wormholestest.com"
	models.NftstIpfsServerPort = "8666"
	hash, err := AddDirIpfs("D:\\temp\\demo")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}

func TestPinIpfs(t *testing.T) {
	//models.NftIpfsServerIP = "https://www.wormholestest.com"
	models.NftIpfsServerIP = "192.168.1.237"
	models.NftstIpfsServerPort = "5001"
	url := models.NftIpfsServerIP + ":" + models.NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(20 * time.Second)
	pins, err := s.Pins()
	if err != nil {
		log.Println("AddDirIpfs() err=", err)
		return
	}
	err = s.Pin("QmNbNvhW1StGPQaXhXMQcfT6W7HqEXDY6MfZijuRLf7Roa")
	if err != nil {
		log.Println("AddDirIpfs() err=", err)
		return
	}
	fmt.Println(pins)
}

func TestName(t *testing.T) {
	/*if snft == "" {
		snft = DefaultSnft
	}*/
	tt := hex.EncodeToString([]byte{0})
	fmt.Println(tt)
	tt = hex.EncodeToString([]byte{0xff})
	fmt.Println(tt)
	snft := "0x8000000000000000000000000000000000000000"
	//addr := common.HexToAddress(snft)
	h := big.NewInt(0)
	h, err := big.NewInt(0).SetString(snft[2:], 16)
	fmt.Println(err)
	h = h.Add(h,big.NewInt(256))
	snft = common.BigToAddress(h).Hex()
	fmt.Println("BackupIpfsSnft() call SyncNftFromChain() blockNum=")
}