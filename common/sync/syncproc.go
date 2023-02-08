package sync

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	waitTime            = time.Second * 60
	ScanBlockTime       = time.Second * 1
	ScanSnftBlockTime   = time.Second * 5
	WaitIpfsFailTime    = time.Second * 1
	SaveIpfsToLocalTime = time.Minute * 30
	ScanIpfsFlagTime    = time.Second * 10
	ZeroAddr            = "0x0000000000000000000000000000000000000000"
	DefaultSnft         = "0x80000000000000000000000000000000000000"
)

func SyncBlockTxs(sqldsn string, block uint64, blockTxs []*contracts.NftTx) error {
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		fmt.Printf("SyncBlockTxs() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	for _, nftTx := range blockTxs {
		if nftTx.From == "" {
			err = nd.BuyResultRoyalty(nftTx.From, nftTx.To, nftTx.Contract, nftTx.TokenId, "", nftTx.Ratio, nftTx.TxHash, nftTx.Ts)
			if err != nil {
				break
			}
		}
		if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" && nftTx.Price != "" &&
			nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
			fmt.Println("SyncBlockTxs() nftTx.Value=", nftTx.Value)
			var price string
			if len(nftTx.Price) >= 9 {
				price = nftTx.Price[:len(nftTx.Price)-9]
			} else {
				continue
				//price = "0"
			}
			err = nd.BuyResultWithAmount(nftTx.From, nftTx.To, nftTx.Contract, nftTx.TokenId,
				nftTx.Value, price, nftTx.Ratio, nftTx.TxHash, nftTx.Ts)
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		var params models.SysParams
		dbErr := nd.GetDB().Last(&params)
		if dbErr.Error != nil {
			fmt.Println("SyncBlockTxs() params err=", dbErr.Error)
			return dbErr.Error
		}
		dbErr = nd.GetDB().Model(&models.SysParams{}).Where("id = ?", params.ID).Update("scannumber", block+1)
		if dbErr.Error != nil {
			fmt.Println("SyncBlockTxs() update params err=", dbErr.Error)
			return dbErr.Error
		}
		fmt.Println("SyncBlockTxs() update block=", block)
	}
	return err
}

func SyncBlockTxsNew(sqldsn string, block uint64, blockTrans contracts.NftTrans) error {
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		fmt.Printf("SyncBlockTxs() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	for _, mintTx := range blockTrans.Minttxs {
		if mintTx.From == "" {
			err = nd.BuyResultRoyalty(mintTx.From, mintTx.To, mintTx.Contract, mintTx.TokenId, "", mintTx.Ratio, mintTx.TxHash, mintTx.Ts)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultRoyalty() err=", err)
				return err
			}
		}
	}
	for _, nftTx := range blockTrans.Nfttxs {
		if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" && nftTx.Price != "" &&
			nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
			fmt.Println("SyncBlockTxs() nftTx.Value=", nftTx.Value)
			var price string
			if len(nftTx.Price) >= 9 {
				price = nftTx.Price[:len(nftTx.Price)-9]
			} else {
				continue
				//price = "0"
			}
			err = nd.BuyResultWithAmount(nftTx.From, nftTx.To, nftTx.Contract, nftTx.TokenId,
				nftTx.Value, price, nftTx.Ratio, nftTx.TxHash, nftTx.Ts)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWithAmount() err=", err)
				return err
			}

		}
	}
	for _, mintTx := range blockTrans.Wmintxs {
		if mintTx.From == "" {
			err = nd.BuyResultWRoyalty(mintTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWRoyalty() err=", err)
				return err
			}
		}
	}
	for _, nftTx := range blockTrans.Wnfttxs {
		log.Println("nfttx :", nftTx)
		if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" /*&& nftTx.Price != ""*/ &&
			nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
			err = nd.BuyResultWithWAmount(nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWithWAmount() err=", err)
				return err
			}

			models.GetRedisCatch().SetDirtyFlag(models.NftCacheDirtyName)
			models.GetRedisCatch().SetDirtyFlag(models.TradingDirtyName)
		}
		if nftTx.From == ZeroAddr && nftTx.To != ZeroAddr {
			err = nd.BuyResultWTransfer(nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWTransfer() err=", err)
				return err
			}

			models.GetRedisCatch().SetDirtyFlag(models.NftCacheDirtyName)
			models.GetRedisCatch().SetDirtyFlag(models.TradingDirtyName)
		}
		if nftTx.From != ZeroAddr && nftTx.To == ZeroAddr {
			err = nd.BuyResultExchange(nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultExchange() err=", err)
				return err
			}
			models.GetRedisCatch().SetDirtyFlag(models.NftCacheDirtyName)
			models.GetRedisCatch().SetDirtyFlag(models.TradingDirtyName)
		}
	}

	if err == nil {
		var params models.SysParams
		dbErr := nd.GetDB().Last(&params)
		if dbErr.Error != nil {
			fmt.Println("SyncBlockTxs() params err=", dbErr.Error)
			return dbErr.Error
		}
		dbErr = nd.GetDB().Model(&models.SysParams{}).Where("id = ?", params.ID).Update("scannumber", block+1)
		if dbErr.Error != nil {
			fmt.Println("SyncBlockTxs() update params err=", dbErr.Error)
			return dbErr.Error
		}
		fmt.Println("SyncBlockTxs() update block=", block)
	}
	return err
}

func SelfSyncBlockTxs(sqldsn string, block uint64, blockTrans contracts.NftTrans) error {
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		fmt.Printf("SyncBlockTxs() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	for _, mintTx := range blockTrans.Minttxs {
		if mintTx.From == "" {
			err = nd.BuyResultRoyalty(mintTx.From, mintTx.To, mintTx.Contract, mintTx.TokenId, "", mintTx.Ratio, mintTx.TxHash, mintTx.Ts)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultRoyalty() err=", err)
				return err
			}
		}
	}
	for _, nftTx := range blockTrans.Nfttxs {
		if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" && nftTx.Price != "" &&
			nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
			fmt.Println("SyncBlockTxs() nftTx.Value=", nftTx.Value)
			var price string
			if len(nftTx.Price) >= 9 {
				price = nftTx.Price[:len(nftTx.Price)-9]
			} else {
				continue
				//price = "0"
			}
			err = nd.BuyResultWithAmount(nftTx.From, nftTx.To, nftTx.Contract, nftTx.TokenId,
				nftTx.Value, price, nftTx.Ratio, nftTx.TxHash, nftTx.Ts)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWithAmount() err=", err)
				return err
			}

		}
	}
	for _, mintTx := range blockTrans.Wmintxs {
		if mintTx.From == "" {
			err = nd.BuyResultWRoyalty(mintTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWRoyalty() err=", err)
				return err
			}
		}
	}
	for _, nftTx := range blockTrans.Wnfttxs {
		log.Println("nfttx :", nftTx)
		if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" /*&& nftTx.Price != ""*/ &&
			nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
			err = nd.BuyResultWithWAmount(nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWithWAmount() err=", err)
				return err
			}

			models.GetRedisCatch().SetDirtyFlag(models.NftCacheDirtyName)
			models.GetRedisCatch().SetDirtyFlag(models.TradingDirtyName)
		}
		if nftTx.From == ZeroAddr && nftTx.To != ZeroAddr {
			err = nd.BuyResultWTransfer(nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWTransfer() err=", err)
				return err
			}

			models.GetRedisCatch().SetDirtyFlag(models.NftCacheDirtyName)
			models.GetRedisCatch().SetDirtyFlag(models.TradingDirtyName)
		}
		if nftTx.From != ZeroAddr && nftTx.To == ZeroAddr {
			err = nd.BuyResultExchange(nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultExchange() err=", err)
				return err
			}
			models.GetRedisCatch().SetDirtyFlag(models.NftCacheDirtyName)
			models.GetRedisCatch().SetDirtyFlag(models.TradingDirtyName)
		}
	}
	return err
}

func SelfSync(sqldsn string) error {
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		log.Printf("SelfSync() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	blockS := uint64(0)
	if models.TransferSNFT {
		blockS, err = models.GetDbSnftBlockNumber(sqldsn)
		if err != nil {
			log.Println("SelfSync() get scan block num err=", err)
			return err
		}
	} else {
		blockS, err = models.GetDbBlockNumber(sqldsn)
		if err != nil {
			log.Println("SelfSync() get scan block num err=", err)
			return err
		}
	}
	for curBlock := contracts.GetCurrentBlockNumber(); blockS <= curBlock; {
		if models.TransferSNFT {
			log.Println("SelfSync() call ScanWorkerNft() blockNum=", blockS)
			err := SelfScanWorkerNft(sqldsn, blockS)
			if err != nil {
				log.Println("SelfSync() call SyncWorkerNft() err=", err)
				return err
			}
			fmt.Println("SelfSync() sync ScanWorkerNft ok.  blockNum=", blockS)
		}
		//txs, err := contracts.GetBlockTxsNew(blockS)
		txs, err := contracts.SelfGetBlockTxs(blockS)
		if err != nil {
			fmt.Println("SelfSync() call GetBlockTxs() err=", err)
			return err
		}
		//err = SelfSyncBlockTxs(sqldsn, blockS, *txs)
		err = models.SyncTxs(nd, txs)
		if err != nil {
			log.Println("SelfSync() SyncBlockTxs err=", err)
			return err
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
			return err
		}
		fmt.Println("SelfSync() update OK block=", blockS)
		if blockS >= curBlock {
			curBlock = contracts.GetCurrentBlockNumber()
		}
	}
	return err
}

func InitSyncBlockTs(sqldsn string) error {
	if !models.LimitWritesDatabase {
		if models.NftScanServer != "" {
			for {
				err := models.SyncBlock(sqldsn)
				if err == nil {
					break
				}
			}
			go models.SyncChain(sqldsn)
		} else {
			return errors.New("nftscan server is not set.")
			/*	for {
					err := SelfSync(sqldsn)
					if err == nil {
						break
					}
					time.Sleep(time.Second)
				}
				go func() {
					ticker := time.NewTicker(ScanBlockTime)
					for {
						select {
						case <-ticker.C:
							SelfSync(sqldsn)
						}
					}
				}()*/
		}
	}

	go BackupIpfsSnft(sqldsn)
	return nil
}

func HttpGetSendRev(url string, data string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	if strings.Index(url, "https") != -1 {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("token", token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type ResponseNftInfo struct {
	Code        string                  `json:"code"`
	Msg         string                  `json:"msg"`
	Data        []models.NftSysMintInfo `json:"data"`
	Total_count int                     `json:"total_count"`
}

func GetNftSysMintInfo(blockNumber uint64) ([]models.NftSysMintInfo, error) {
	srcurl := contracts.BrowseNode + "/v2/querymetaurl/"
	bn := strconv.FormatUint(blockNumber, 10)
	//srcurl = srcurl + "blocknumber=" + bn
	srcurl = srcurl + bn
	b, err := HttpGetSendRev(srcurl, "", "")
	if err != nil {
		fmt.Println("GetNftSysMintInfo() err=", err)
		return nil, err
	}
	var revData ResponseNftInfo
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("AuditKYC() get resp failed, err", err)
		return nil, err
	}
	if revData.Code != "200" {
		return nil, errors.New(revData.Msg)
	}
	return revData.Data, nil
}

func GetSnftInfoFromIPFS(hash string) (*models.SnftInfo, error) {
	//url := models.NftIpfsServerIP + ":" + models.NftstIpfsServerPort
	url := "http://api.wormholestest.com" + ":" + "8666"
	url = url + hash
	b, err := HttpGetSendRev(url, "", "")
	if err != nil {
		log.Println("GetSnftInfoFromIPFS() err=", err)
		return nil, err
	}
	var snft models.SnftInfo
	err = json.Unmarshal([]byte(b), &snft)
	if err != nil {
		log.Println("GetSnftInfoFromIPFS() get resp failed, err", err)
		return nil, err
	}
	return &snft, nil
}

func GetSnftInfoFromIPFSWithShell(hash string) (*models.SnftInfo, error) {
	url := models.NftIpfsServerIP + ":" + models.NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(100 * time.Second)
	rc, err := s.Cat(hash)
	if err != nil {
		log.Println("GetSnftInfoFromIPFSWithShell() err=", err)
		return nil, err
	}
	var snft models.SnftInfo
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Println("GetSnftInfoFromIPFSWithShell() ReadAll() err=", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(b), &snft)
	if err != nil {
		log.Println("GetSnftInfoFromIPFSWithShell() Unmarshal, err=", err)
		return nil, err
	}
	return &snft, nil
}

/*func SyncWorkerNft(sqldsn string, blockS uint64) error {
	nftInfo, err := GetNftSysMintInfo(blockS)
	if err != nil {
		fmt.Println("SyncWorkerNft() err=", err)
		return err
	}
	if len(nftInfo) != 0 {
		nd, err := models.NewNftDb(sqldsn)
		if err != nil {
			fmt.Printf("SyncWorkerNft() connect database err = %s\n", err)
			return err
		}
		defer nd.Close()
		for _, info := range nftInfo {
			err = nd.UploadWNft(&info)
			if err != nil {
				fmt.Println("SyncWorkerNft() upload err=", err)
				return err
			}
		}
	}
	return nil
}
*/

func AddDirIpfs(dir string) (string, error) {
	url := models.NftIpfsServerIP + ":" + models.NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(20 * time.Second)
	cid, err := s.AddDir(dir)
	if err != nil {
		log.Println("AddDirIpfs() err=", err)
		return "", err
	}
	return cid, err
}

func SaveIpfsToLocal(ipfsHash string) error {
	url := models.BackupIpfsUrl
	s := shell.NewShell(url)
	s.SetTimeout(SaveIpfsToLocalTime)
	err := s.Pin(ipfsHash)
	if err != nil {
		log.Println("SaveIpfsToLocal() err=", err)
		return err
	}
	return err
}

type SnftInfoData struct {
	SnftInfo *models.SnftInfo
	TimeTag  time.Time
}

type IpfsCatch struct {
	Mux      sync.Mutex
	SnftInfo map[string]*SnftInfoData
}

func (n *IpfsCatch) GetByHash(hash string) *models.SnftInfo {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	fmt.Println("IpfsCatch-GetByHash() GetByHash n.NftInfo catch len=", len(n.SnftInfo))
	if len(n.SnftInfo) == 0 {
		n.SnftInfo = make(map[string]*SnftInfoData)
	}
	if nftinfo := n.SnftInfo[hash]; nftinfo != nil {
		fmt.Println("IpfsCatch-GetByHash() NftFilterCatch hash=", hash)
		s := *nftinfo.SnftInfo
		return &s
	}
	return nil
}

func (n *IpfsCatch) SetByHash(hash string, snftinfo *models.SnftInfo) *models.SnftInfo {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if len(n.SnftInfo) == 0 {
		fmt.Println("IpfsCatch-SetByHash() NftFilterCatch len ==0 ")
		n.SnftInfo = make(map[string]*SnftInfoData)
	}
	s := *snftinfo
	n.SnftInfo[hash] = &SnftInfoData{&s, time.Now().Add(time.Minute * 30)}
	fmt.Println("IpfsCatch-SetByHash() NftFilterCatch", "len=", len(n.SnftInfo), " hash=", hash)
	for s, info := range n.SnftInfo {
		if info.TimeTag.Before(time.Now()) {
			delete(n.SnftInfo, s)
		}
	}
	return nil
}

var ScanIpfsCatch IpfsCatch

func ScanWorkerNft(sqldsn string, blockS uint64) error {
	snftAddr, err := contracts.GetSnftAddressList(big.NewInt(0).SetUint64(blockS), true)
	if err != nil {
		log.Println("ScanWorkerNft() GetSnftAddressList err =", err, "blocks", blockS)
		return err
	}
	snftInfos := make([]models.SnftInfo, len(snftAddr))
	if len(snftAddr) > 0 {
		for i, address := range snftAddr {
			if address.NftAddress.String() == ZeroAddr {
				continue
			}
			accountInfo, err := contracts.GetAccountInfo(address.NftAddress, big.NewInt(0).SetUint64(blockS))
			if err != nil {
				log.Println("ScanWorkerNft() GetAccountInfo err =", err, "NftAddress= ", address.NftAddress, "blocks", blockS)
				return err
			}
			fmt.Println("ScanWorkerNft() MetaUrl=", accountInfo.MetaURL, "blockS=", blockS)
			index := strings.Index(accountInfo.MetaURL, "/ipfs/")
			if index == -1 {
				log.Printf("ScanWorkerNft() Index ipfs error.\n")
				continue
				return errors.New("ScanWorkerNft(): MetaUrl error.")
			}
			index = strings.LastIndex(accountInfo.MetaURL, "/")
			if index == -1 {
				log.Printf("ScanWorkerNft() LastIndex error.\n")
				continue
				return errors.New("ScanWorkerNft(): MetaUrl error.")
			}
			/*if accountInfo.MetaURL[:index] == "/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By" ||
				accountInfo.MetaURL[:index] == "/ipfs/QmaiReZpUeWcSRvhWhHwQ4PN2NbggYdZt7hKFAoM8kTVF7" {
				continue
			}*/
			var metaUrl, metaHash string
			if accountInfo.MetaURL[:index] == "/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7" ||
				accountInfo.MetaURL[:index] == "/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By" {
				//metaUrl = "/ipfs/QmVyVJTMQVbHRz8dr8RHrW4c1pgnspcM3Ee1pj9vae2oo8" //1.237
				metaUrl = "/ipfs/QmNbNvhW1StGPQaXhXMQcfT6W7HqEXDY6MfZijuRLf7Roa" //云服务器
				//metaUrl = "/ipfs/QmWpDcyU287P3bgw74nmUmWGDcaRYGud51y8xxQkiK5zDR" //云服务器
				//metaHash = metaUrl + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-4:len(accountInfo.MetaURL)-2])
				metaHash = metaUrl + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-3:len(accountInfo.MetaURL)-1])
			} else {
				//metaHash = accountInfo.MetaURL[:index] + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-4:len(accountInfo.MetaURL)-2])
				metaHash = accountInfo.MetaURL[:index] + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-3:len(accountInfo.MetaURL)-1])
			}
			fmt.Println("ScanWorkerNft() metaHash=", metaHash)
			var snftinfo *models.SnftInfo
			if snftinfo = ScanIpfsCatch.GetByHash(metaHash); snftinfo == nil {
				retry := 0
				for {
					snftinfo, err = GetSnftInfoFromIPFSWithShell(metaHash)
					if err != nil {
						log.Println("ScanWorkerNft() GetSnftInfoFromIPFS count=", retry, " err =", err, "ipfs hash=", metaHash)
						errflag := strings.Index(err.Error(), "context deadline exceeded")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "connection refused")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "502 Bad Gateway")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "403 Forbidden")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
					}
					break
				}
				if err != nil {
					continue
				}
				ScanIpfsCatch.SetByHash(metaHash, snftinfo)
			}
			snftinfo.Ownaddr = strings.ToLower(accountInfo.Owner.String())
			snftinfo.Contract = models.ExchangeOwer
			snftinfo.Nftaddr = strings.ToLower(address.NftAddress.String())
			snftinfo.Meta = accountInfo.MetaURL
			b, _ := big.NewInt(0).SetString(snftinfo.Nftaddr[models.SnftCollectionsStageIndex:models.SnftStageOffset], 16)
			snftinfo.CollectionsName = b.String() + "-" + snftinfo.CollectionsName
			snftInfos[i] = *snftinfo
		}
	}
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		log.Printf("ScanWorkerNft() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	if len(snftInfos) != 0 {
		for _, info := range snftInfos {
			if info.Nftaddr == "" {
				continue
			}
			err = nd.UploadWNft(&info)
			if err != nil {
				log.Println("ScanWorkerNft() upload err=", err)
				return err
			}
		}
	}
	if err == nil {
		var params models.SysParams
		dbErr := nd.GetDB().Select([]string{"scansnftnumber", "id"}).Last(&params)
		if dbErr.Error != nil {
			log.Println("ScanWorkerNft() params err=", dbErr.Error)
			return dbErr.Error
		}
		dbErr = nd.GetDB().Model(&models.SysParams{}).Where("id = ?", params.ID).Update("scansnftnumber", blockS+1)
		if dbErr.Error != nil {
			log.Println("ScanWorkerNft() update params err=", dbErr.Error)
			return dbErr.Error
		}
		fmt.Println("ScanWorkerNft() update block=", blockS)
	}
	return err
}

func SelfScanWorkerNft(sqldsn string, blockS uint64) error {
	snftAddr, err := contracts.GetSnftAddressList(big.NewInt(0).SetUint64(blockS), true)
	if err != nil {
		log.Println("ScanWorkerNft() GetSnftAddressList err =", err, "blocks", blockS)
		return err
	}
	snftInfos := make([]models.SnftInfo, len(snftAddr))
	if len(snftAddr) > 0 {
		for i, address := range snftAddr {
			if address.NftAddress.String() == ZeroAddr {
				continue
			}
			accountInfo, err := contracts.GetAccountInfo(address.NftAddress, big.NewInt(0).SetUint64(blockS))
			if err != nil {
				log.Println("ScanWorkerNft() GetAccountInfo err =", err, "NftAddress= ", address.NftAddress, "blocks", blockS)
				return err
			}
			fmt.Println("ScanWorkerNft() MetaUrl=", accountInfo.MetaURL, "blockS=", blockS)
			index := strings.Index(accountInfo.MetaURL, "/ipfs/")
			if index == -1 {
				log.Printf("ScanWorkerNft() Index ipfs error.\n")
				continue
				return errors.New("ScanWorkerNft(): MetaUrl error.")
			}
			index = strings.LastIndex(accountInfo.MetaURL, "/")
			if index == -1 {
				log.Printf("ScanWorkerNft() LastIndex error.\n")
				continue
				return errors.New("ScanWorkerNft(): MetaUrl error.")
			}
			/*if accountInfo.MetaURL[:index] == "/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By" ||
				accountInfo.MetaURL[:index] == "/ipfs/QmaiReZpUeWcSRvhWhHwQ4PN2NbggYdZt7hKFAoM8kTVF7" {
				continue
			}*/
			var metaUrl, metaHash string
			if accountInfo.MetaURL[:index] == "/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7" ||
				accountInfo.MetaURL[:index] == "/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By" {
				//metaUrl = "/ipfs/QmVyVJTMQVbHRz8dr8RHrW4c1pgnspcM3Ee1pj9vae2oo8" //1.237
				metaUrl = "/ipfs/QmNbNvhW1StGPQaXhXMQcfT6W7HqEXDY6MfZijuRLf7Roa" //云服务器
				//metaUrl = "/ipfs/QmWpDcyU287P3bgw74nmUmWGDcaRYGud51y8xxQkiK5zDR" //云服务器
				//metaHash = metaUrl + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-4:len(accountInfo.MetaURL)-2])
				metaHash = metaUrl + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-3:len(accountInfo.MetaURL)-1])
			} else {
				//metaHash = accountInfo.MetaURL[:index] + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-4:len(accountInfo.MetaURL)-2])
				metaHash = accountInfo.MetaURL[:index] + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-3:len(accountInfo.MetaURL)-1])
			}
			fmt.Println("ScanWorkerNft() metaHash=", metaHash)
			var snftinfo *models.SnftInfo
			if snftinfo = ScanIpfsCatch.GetByHash(metaHash); snftinfo == nil {
				retry := 0
				for {
					snftinfo, err = GetSnftInfoFromIPFSWithShell(metaHash)
					if err != nil {
						log.Println("ScanWorkerNft() GetSnftInfoFromIPFS count=", retry, " err =", err, "ipfs hash=", metaHash)
						errflag := strings.Index(err.Error(), "context deadline exceeded")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "connection refused")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "502 Bad Gateway")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "403 Forbidden")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
					}
					break
				}
				if err != nil {
					continue
				}
				ScanIpfsCatch.SetByHash(metaHash, snftinfo)
			}
			snftinfo.Ownaddr = strings.ToLower(accountInfo.Owner.String())
			snftinfo.Contract = models.ExchangeOwer
			snftinfo.Nftaddr = strings.ToLower(address.NftAddress.String())
			snftinfo.Meta = accountInfo.MetaURL
			b, _ := big.NewInt(0).SetString(snftinfo.Nftaddr[models.SnftCollectionsStageIndex:models.SnftStageOffset], 16)
			snftinfo.CollectionsName = b.String() + "-" + snftinfo.CollectionsName
			snftInfos[i] = *snftinfo
		}
	}
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		log.Printf("ScanWorkerNft() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	if len(snftInfos) != 0 {
		for _, info := range snftInfos {
			if info.Nftaddr == "" {
				continue
			}
			err = nd.UploadWNft(&info)
			if err != nil {
				log.Println("ScanWorkerNft() upload err=", err)
				return err
			}
		}
	}
	return err
}
func SyncWorkerNft(sqldsn string) error {
	blockS, err := models.GetDbSnftBlockNumber(sqldsn)
	if err != nil {
		fmt.Println("SyncWorkerNft() get scan block num err=", err)
		return err
	}
	for {
		if !models.TransferSNFT {
			time.Sleep(ScanSnftBlockTime)
			continue
		}
		if blockS < contracts.GetCurrentBlockNumber() {
			fmt.Println("SyncWorkerNft() call SyncNftFromChain() blockNum=", blockS)
			err := ScanWorkerNft(sqldsn, blockS)
			if err != nil {
				log.Println("SyncWorkerNft() call SyncWorkerNft() err=", err)
				continue
				//return err
			}
			blockS = blockS + 1
			fmt.Println("SyncWorkerNft() call SyncNftFromChain() blockNum=", blockS)
		} else {
			time.Sleep(ScanSnftBlockTime)
		}
	}
	return nil
}

func SaveIpfsMeta(sqldsn string) error {
	nd, err := models.NewNftDb(sqldsn)
	if err != nil {
		log.Printf("SaveIpfsMeta() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	snft := ""
	params := models.SysParams{}
	dbErr := nd.GetDB().Select("Savedsnft").Last(&params)
	if dbErr.Error != nil {
		fmt.Println("SaveIpfsMeta() opendb err=", dbErr.Error)
		return dbErr.Error
	}
	if params.Savedsnft == "" || len(params.Savedsnft) != 40 {
		snft = DefaultSnft
	} else {
		snft = params.Savedsnft
	}
	for {
		nftData := models.Nfts{}
		result := nd.GetDB().Model(&models.Nfts{}).Select([]string{"meta"}).Where("snft = ?", snft).First(&nftData)
		if result.Error != nil {
			log.Println("SaveIpfsMeta() no snft save to ipfs err=", dbErr.Error)
			return result.Error
		}
		log.Println("SaveIpfsMeta() backup snft meta snft=", snft)
		index := strings.LastIndex(nftData.Meta, "/")
		if index == -1 {
			log.Printf("ScanWorkerNft() LastIndex error.\n")
			h, _ := big.NewInt(0).SetString(snft[2:]+"00", 16)
			h = h.Add(h, big.NewInt(256))
			snft = common.BigToAddress(h).Hex()
			snft = snft[:len(snft)-2]
			continue
			return errors.New("ScanWorkerNft(): MetaUrl error.")
		}
		for {
			fmt.Println("SaveIpfsMeta() backup to ipfs main hash=", nftData.Meta[:index])
			err := SaveIpfsToLocal(nftData.Meta[:index])
			if err != nil {
				time.Sleep(WaitIpfsFailTime)
				continue
			}
			break
		}
		for i := 0; i < 256; i++ {
			metaHash := nftData.Meta[:index] + "/" + hex.EncodeToString([]byte{byte(i)})
			fmt.Println("SaveIpfsMeta() backup to ipfs chip hash=", metaHash)
			var snftinfo *models.SnftInfo
			retry := 0
			for {
				snftinfo, err = GetSnftInfoFromIPFSWithShell(metaHash)
				if err != nil {
					log.Println("SaveIpfsMeta() GetSnftInfoFromIPFS count=", retry, " err =", err, "ipfs hash=", metaHash)
					time.Sleep(WaitIpfsFailTime)
					continue
				}
				break
			}
			metaHash = snftinfo.SourceUrl
			for {
				err := SaveIpfsToLocal(metaHash)
				if err != nil {
					time.Sleep(WaitIpfsFailTime)
					continue
				}
				break
			}
			metaHash = snftinfo.CollectionsImgUrl
			for {
				err := SaveIpfsToLocal(metaHash)
				if err != nil {
					time.Sleep(WaitIpfsFailTime)
					continue
				}
				break
			}
		}
		h, _ := big.NewInt(0).SetString(snft[2:]+"00", 16)
		h = h.Add(h, big.NewInt(256))
		snft = common.BigToAddress(h).Hex()
		snft = snft[:len(snft)-2]
		err := nd.GetDB().Transaction(func(tx *gorm.DB) error {
			var params models.SysParams
			dbErr := nd.GetDB().Select([]string{"savedsnft", "id"}).Last(&params)
			if dbErr.Error != nil {
				log.Println("SaveIpfsMeta() get savedsnft  err=", dbErr.Error)
				return dbErr.Error
			}
			params.Savedsnft = snft
			dbErr = nd.GetDB().Model(&models.SysParams{}).Where("id = ?", params.ID).Updates(params)
			if dbErr.Error != nil {
				log.Println("SaveIpfsMeta() update savedsnft err=", dbErr.Error)
				return dbErr.Error
			}
			return nil
		})
		if err != nil {
			log.Println("SaveIpfsMeta() update savedsnft err=", dbErr.Error)
			return err
		}
	}
	return err
}

func BackupIpfsSnft(sqldsn string) error {
	for {
		if !models.Backupipfs {
			time.Sleep(ScanIpfsFlagTime)
			continue
		}
		err := SaveIpfsMeta(sqldsn)
		if err != nil {
			fmt.Println("BackupIpfsSnft() no snft to SaveIpfsMeta")
			time.Sleep(ScanIpfsFlagTime)
		}
	}
	return nil
}

func SyncBlockNew(sqldsn string) error {
	blockS, err := models.GetDbBlockNumber(sqldsn)
	if err != nil {
		fmt.Println("SyncProc() get scan block num err=", err)
		return err
	}
	//blockS = 37215
	//blockS = 18351
	for blockS <= contracts.GetCurrentBlockNumber() {
		txs, err := contracts.GetBlockTxsNew(blockS)
		if err != nil {
			fmt.Println("SyncProc() call GetBlockTxs() err=", err)
			return err
		}
		err = SyncBlockTxsNew(sqldsn, blockS, *txs)
		if err != nil {
			fmt.Println("SyncProc() SyncBlockTxs err=", err)
			return err
		}
		if len(txs.Wethc) != 0 {
			err = models.ScanBiddings(sqldsn, txs.Wethc)
			if err != nil {
				fmt.Println("SyncProc() ScanBiddings err=", err)
				//return err
			}
		}
		blockS++
	}
	return nil
}
