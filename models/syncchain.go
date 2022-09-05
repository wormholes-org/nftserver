package models

import (
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

const (
	scanBlockTime  = time.Second * 1
	ErrorsWaitTime = time.Second * 1
)

func SyncTxs(nd *NftDb, txn []contracts.NftTx) error {
	var err error
	txs := []contracts.NftTx{}
	for _, nftTx := range txn {
		if nftTx.TransType == contracts.WormHolesMint ||
			nftTx.TransType == contracts.WormHolesBuyFromSellMintTransfer ||
			nftTx.TransType == contracts.WormHolesExToBuyMintToSellTransfer ||
			nftTx.TransType == contracts.WormHolesExAuthToExMintBuyTransfer {
			if nftTx.From == "" && nftTx.Contract == ExchangeOwer {
				err = nd.BuyResultWRoyalty(&nftTx)
				if err != nil {
					log.Println("syncTxs() BuyResultWRoyalty() err=", err)
					return err
				}
				GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
				GetRedisCatch().SetDirtyFlag(TradingDirtyName)
				continue
			}
		}
		if nftTx.TransType == contracts.WormHolesExchange ||
			nftTx.TransType == contracts.WormHolesUnPledge ||
			nftTx.TransType == contracts.WormHolesPledge {
			if !TransferSNFT {
				continue
			}
		}
		txs = append(txs, nftTx)
	}
	for _, nftTx := range txs {
		if nftTx.TransType == contracts.WormHolesBuyFromSellMintTransfer ||
			nftTx.TransType == contracts.WormHolesExToBuyMintToSellTransfer ||
			nftTx.TransType == contracts.WormHolesExAuthToExMintBuyTransfer {
			if nftTx.Contract == ExchangeOwer {
				err = nd.BuyResultWithWAmount(&nftTx)
				if err != nil {
					fmt.Println("SyncBlockTxs() BuyResultWithWAmount() err=", err)
					return err
				}
				GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
				GetRedisCatch().SetDirtyFlag(TradingDirtyName)
				continue
			}
		}
		if nftTx.TransType == contracts.WormHolesExToBuyTransfer ||
			nftTx.TransType == contracts.WormHolesBuyFromSellTransfer ||
			nftTx.TransType == contracts.WormHolesExAuthToExBuyTransfer {
			err = nd.BuyResultWithWAmount(&nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWithWAmount() err=", err)
				return err
			}
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(TradingDirtyName)
			continue
		}
		if nftTx.TransType == contracts.WormHolesTransfer {
			if !nftTx.Status {
				continue
			}
			err = nd.BuyResultWTransfer(&nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWTransfer() err=", err)
				return err
			}
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(TradingDirtyName)
			continue
		}
		if nftTx.TransType == contracts.WormHolesExchange {
			if !nftTx.Status && TransferSNFT {
				continue
			}
			err = nd.BuyResultExchange(&nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultExchange() err=", err)
				return err
			}
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(TradingDirtyName)
			continue
		}
		if nftTx.TransType == contracts.WormHolesPledge || nftTx.TransType == contracts.WormHolesUnPledge {
			if !nftTx.Status && TransferSNFT {
				continue
			}
			err = nd.BuyResultWPledge(&nftTx)
			if err != nil {
				fmt.Println("SyncBlockTxs() BuyResultWPledge() err=", err)
				return err
			}
			GetRedisCatch().SetDirtyFlag(NftCacheDirtyName)
			GetRedisCatch().SetDirtyFlag(TradingDirtyName)
		}
	}
	return nil
}

func SyncBlock(sqldsn string) error {
	nd, err := NewNftDb(sqldsn)
	if err != nil {
		log.Println("SyncBlock() connect database err=", err)
		time.Sleep(ErrorsWaitTime)
		return err
	}
	defer nd.Close()
	blockS, err := GetDbBlockNumber(sqldsn)
	if err != nil {
		time.Sleep(ErrorsWaitTime)
		log.Println("SyncBlock() GetDbBlockNumber err=", err)
		return err
	}
	curBlock := contracts.GetCurrentBlockNumber()
	for blockS <= curBlock /*curBlock = contracts.GetCurrentBlockNumber()*/ {
		blockStr := strconv.FormatUint(blockS, 10)
		if TransferSNFT {
			spendT := time.Now()
			snfts, err := GetBlockSnfts(blockStr)
			if err != nil {
				log.Println("SyncBlock() GetBlockSnfts blocknumber=", blockS, "err=", err)
				time.Sleep(ErrorsWaitTime)
				return err
			}
			fmt.Println("SyncBlock() GetBlockTrans snfts len =", len(snfts), "sync block =", blockStr)
			fmt.Println("SyncBlock() GetBlockSnfts spend time =", time.Now().Sub(spendT), "time.now=", time.Now())
			if len(snfts) != 0 {
				for _, info := range snfts {
					if info.Nftaddr == "" {
						continue
					}
					info.Contract = ExchangeOwer
					err = nd.UploadWNft(&info)
					if err != nil {
						log.Println("SyncBlock() upload snft err=", err)
						time.Sleep(ErrorsWaitTime)
						return err
					}
				}
			}
		}
		spendT := time.Now()
		txs, err := GetBlockTrans(blockStr)
		if err != nil {
			log.Println("SyncBlock() GetBlockTrans error blocknumber=", blockS, "err=", err)
			time.Sleep(ErrorsWaitTime)
			return err
		}
		fmt.Println("SyncBlock() GetBlockTrans txs len =", len(txs), "sync block =", blockStr)
		fmt.Println("SyncBlock() GetBlockTrans spend time =", time.Now().Sub(spendT), "time.now=", time.Now())
		//var newtxs []contracts.NftTx
		/*for i, tx := range txs {
			if tx.Contract != "" {
				if TransferSNFT {
					if ExchangeOwer != tx.Contract && tx.NftAddr[:3] != "0x8" {
						continue
					} else {
						newtxs = append(newtxs, txs[i])
					}
				} else {
					if ExchangeOwer != tx.Contract {
						continue
					} else {
						newtxs = append(newtxs, txs[i])
					}
				}
			} else {
				newtxs = append(newtxs, txs[i])
			}
		}*/
		err = SyncTxs(nd, txs)
		if err != nil {
			log.Println("SyncBlock() syncTxs error blocknumber=", blockS, "err=", err)
			time.Sleep(ErrorsWaitTime)
			return err
		}
		fmt.Println("SyncProc() get chain ok  blocknumber=", blockS)
		params := SysParams{}
		dberr := nd.db.Model(&SysParams{}).Last(&params)
		if dberr.Error != nil && dberr.Error != gorm.ErrRecordNotFound {
			time.Sleep(ErrorsWaitTime)
			return err
		}
		blockS++
		dberr = nd.db.Model(&SysParams{}).Where("id = ?", params.ID).Update("Scannumber", blockS)
		if dberr.Error != nil {
			log.Println("SyncBlock() update params.Scannumber err=", err)
			time.Sleep(ErrorsWaitTime)
			return dberr.Error
		}
		fmt.Println("SyncBlock() update record upload block number=", blockS)
		if blockS >= curBlock {
			curBlock = contracts.GetCurrentBlockNumber()
		}
	}
	return nil
}

func SyncChain(sqldsn string) {
	for {
		SyncBlock(sqldsn)
		time.Sleep(scanBlockTime)
	}
}
