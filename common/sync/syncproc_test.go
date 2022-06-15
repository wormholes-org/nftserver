package sync

import (
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	"strings"
	"testing"
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