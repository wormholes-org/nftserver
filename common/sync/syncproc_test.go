package sync

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
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
	//contracts.EthNode = "http://api.wormholestest.com:8561"
	contracts.EthNode = "http://43.129.181.130:8561"
	models.NftIpfsServerIP = "http://api.wormholestest.com"
	models.NftstIpfsServerPort = "8666"
	models.RoyaltyLimit = 1000
	models.TransferSNFT = true
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
	h = h.Add(h, big.NewInt(256))
	snft = common.BigToAddress(h).Hex()
	fmt.Println("BackupIpfsSnft() call SyncNftFromChain() blockNum=")
}

func TestSyncBlockTxsNew(t *testing.T) {
	const sqlsvrLcT = "admin:user123456@tcp(192.168.1.235:3306)/"
	const dbNameT = "c0x57ed0c503c40308e802414405ce3d399fe3a42c6"
	const localtimeT = "?parseTime=true&loc=Local"
	const sqldsnT = sqlsvrLcT + dbNameT + localtimeT
	//contracts.EthNode = "https://api.wormholestest.com"
	contracts.EthNode = "http://43.129.181.130:8561"
	contracts.ExchangeOwer = "0x57ed0c503c40308e802414405ce3d399fe3a42c6"
	blockS := uint64(18351)
	blockS = uint64(52384)
	blockS = uint64(87736)
	for blockS <= contracts.GetCurrentBlockNumber() {
		txs, err := contracts.GetBlockTxsNew(blockS)
		if err != nil {
			fmt.Println("SyncProc() call GetBlockTxs() err=", err)
			return
		}
		err = SyncBlockTxsNew(sqldsnT, blockS, *txs)
		if err != nil {
			fmt.Println("SyncProc() SyncBlockTxs err=", err)
			return
		}
		if len(txs.Wethc) != 0 {
			err = models.ScanBiddings(sqldsnT, txs.Wethc)
			if err != nil {
				fmt.Println("SyncProc() ScanBiddings err=", err)
				//return err
			}
		}
		blockS++
	}
}

func TestSync(t *testing.T) {
	const sqlsvrLcT = "admin:user123456@tcp(192.168.56.122:3306)/"
	const dbNameT = "snftdb"
	const localtimeT = "?parseTime=true&loc=Local"
	const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

	contracts.EthNode = "http://43.129.181.130:8561"
	contracts.ExchangeOwer = "0x62e0c8032fb51bc401558b58b1e7733276c1ec8a"
	models.ExchangeOwer = "0x62e0c8032fb51bc401558b58b1e7733276c1ec8a"
	nd, err := models.NewNftDb(sqldsnT)
	if err != nil {
		log.Printf("SelfSync() connect database err = %s\n", err)
		return
	}
	models.NewQueryCatch("192.168.1.235:6379", "user123456")
	defer nd.Close()
	blockS := uint64(0)
	if models.TransferSNFT {
		blockS, err = models.GetDbSnftBlockNumber(sqldsnT)
		if err != nil {
			log.Println("SelfSync() get scan block num err=", err)
			return
		}
	} else {
		blockS, err = models.GetDbBlockNumber(sqldsnT)
		if err != nil {
			log.Println("SelfSync() get scan block num err=", err)
			return
		}
	}
	blockS = 34032
	for curBlock := contracts.GetCurrentBlockNumber(); blockS <= curBlock; {
		if models.TransferSNFT {
			log.Println("SelfSync() call ScanWorkerNft() blockNum=", blockS)
			err := SelfScanWorkerNft(sqldsnT, blockS)
			if err != nil {
				log.Println("SelfSync() call SyncWorkerNft() err=", err)
				return
			}
			fmt.Println("SelfSync() sync ScanWorkerNft ok.  blockNum=", blockS)
		}
		//txs, err := contracts.GetBlockTxsNew(blockS)
		txs, err := contracts.SelfGetBlockTxs(blockS)
		if err != nil {
			fmt.Println("SelfSync() call GetBlockTxs() err=", err)
			return
		}
		//err = SelfSyncBlockTxs(sqldsn, blockS, *txs)
		err = models.SyncTxs(nd, txs)
		if err != nil {
			log.Println("SelfSync() SyncBlockTxs err=", err)
			return
		}
		blockS++
		err = nd.GetDB().Transaction(func(tx *gorm.DB) error {
			var params models.SysParams
			dbErr := nd.GetDB().Select("id").Last(&params)
			if dbErr.Error != nil {
				log.Println("SelfSync() params err=", dbErr.Error)
				return dbErr.Error
			}
			nparams := models.SysParams{}
			nparams.Scannumber = blockS
			nparams.Scansnftnumber = blockS
			dbErr = nd.GetDB().Model(&models.SysParams{}).Where("id = ?", params.ID).Updates(&nparams)
			if dbErr.Error != nil {
				log.Println("SelfSync() update params err=", dbErr.Error)
				return dbErr.Error
			}
			return nil
		})
		if err != nil {
			log.Println("SelfSync() update params err=", err)
			return
		}
		fmt.Println("SelfSync() update OK block=", blockS)
		if blockS >= curBlock {
			curBlock = contracts.GetCurrentBlockNumber()
		}
	}
	return
}
