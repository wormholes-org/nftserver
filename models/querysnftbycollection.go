package models

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nftexchange/nftserver/common/contracts"
	"log"
	"time"
)

func (nft NftDb) QuerySnftByCollectionOld(ownaddr, createaddr, name, startIndex, count string) ([]SnftChipInfo, error) {
	spendT := time.Now()
	queryCatchSql := ownaddr + createaddr + name + count
	nftCatch := SnftChipCatch{}
	cerr := GetRedisCatch().GetCatchData("QuerySnftByCollection", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QuerySnftByCollection() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, nil
	}
	snftInfo := []SnftChipInfo{}
	snftAddrs := []string{}
	//stIndex, _ := strconv.Atoi(startIndex)
	//nftCount, _ := strconv.Atoi(count)
	snftChips := []struct {
		SnftAddr  string
		Chipcount int
	}{}
	alltChips := []struct {
		SnftAddr  string
		Chipcount int
	}{}

	allsql := `select left(nftaddr,41) as SnftAddr, 0 as chipcount from nfts where ` +
		`collectcreator = ? and collections = ? and deleted_at is null and nftaddr not like "%m"  GROUP BY SnftAddr`
	err := nft.db.Raw(allsql, createaddr, name).Find(&alltChips)
	if err.Error != nil {
		log.Println("QuerySnftByCollection() Select(all) err=", err.Error)
		return nil, ErrDataBase
	}
	ownsql := `select left(nftaddr,41) as SnftAddr, count(nftaddr) as chipcount from nfts where ` +
		`ownaddr = ? and  collectcreator = ? and collections = ? and deleted_at is null and nftaddr not like "%m" and exchange =0  GROUP BY SnftAddr`
	err = nft.db.Raw(ownsql, ownaddr, createaddr, name).Find(&snftChips)
	if err.Error != nil {
		log.Println("QuerySnftByCollection() Select(snft) err=", err.Error)
		return nil, ErrDataBase
	}
	k := 0
	for j, snft := range alltChips {
		for i := k; i < len(snftChips); i++ {
			if snft.SnftAddr == snftChips[i].SnftAddr {
				alltChips[j].Chipcount = snftChips[i].Chipcount
				k++
				break
			}
		}
		snftAddrs = append(snftAddrs, snft.SnftAddr+"0")
	}
	//sql := `select left(nftaddr,41) as SnftAddr, count(nftaddr) as chipcount from nfts where `
	//if ownaddr != "" {
	//	sql = sql + `ownaddr = ? and  collectcreator = ? and collections = ?  GROUP BY snft`
	//	err := nft.db.Raw(sql, ownaddr, createaddr, name).Scan(&snftChips)
	//	if err.Error != nil {
	//		log.Println("QuerySnftByCollection() Select(snft) err=", err.Error)
	//		return nil, ErrDataBase
	//	}
	//} else {
	//	sql = sql + `collectcreator = ? and collections = ?  GROUP BY snft`
	//	err := nft.db.Raw(sql, createaddr, name).Scan(&snftChips)
	//	if err.Error != nil {
	//		log.Println("QuerySnftByCollection() Select(snft) err=", err.Error)
	//		return nil, ErrDataBase
	//	}
	//}
	//for _, chip := range snftChips {
	//	snftAddrs = append(snftAddrs, chip.SnftAddr)
	//}
	fmt.Printf("QuerySnftByCollection() min(nftaddr)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	nftsql := "select * from nfts where nftaddr in ?"
	err = nft.db.Model(&Nfts{}).Raw(nftsql, snftAddrs).Scan(&snftInfo)
	if err.Error != nil {
		log.Println("QuerySnftByCollection() Scan(&snftInfo) err=", err)
		return nil, ErrDataBase
	}

	for i, chip := range alltChips {
		snftInfo[i].Chipcount = chip.Chipcount
		addr := common.HexToAddress(chip.SnftAddr + "0")
		accountInfo, err := contracts.GetAccountInfo(addr, nil)
		if err != nil {
			log.Println("QuerySnftByCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
			return nil, ErrBlockchain
		}
		snftInfo[i].MergeLevel = accountInfo.MergeLevel
		snftInfo[i].BlockNumber = accountInfo.NFTPledgedBlockNumber
	}
	GetRedisCatch().CatchQueryData("QuerySnftByCollection", queryCatchSql, &SnftChipCatch{snftInfo, uint64(0)})
	fmt.Printf("QuerySnftByCollection()  nftaddr in()  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return snftInfo, nil
}

func (nft NftDb) QuerySnftByCollection(ownaddr, createaddr, name, startIndex, count string) ([]SnftChipInfo, error) {
	spendT := time.Now()
	queryCatchSql := ownaddr + createaddr + name + count
	nftCatch := SnftChipCatch{}
	cerr := GetRedisCatch().GetCatchData("QuerySnftByCollection", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QuerySnftByCollection() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.NftInfo, nil
	}
	ownSnfts := []string{}
	sqlOwnCollects := `select min(snft) as snft from nfts where ` +
		`ownaddr = ? and collectcreator = ? and collections = ? and snft != "" and mergetype = 0 and exchange = 0 and deleted_at is null GROUP BY collectcreator, collections, snft`
	err := nft.db.Raw(sqlOwnCollects, ownaddr, createaddr, name).Find(&ownSnfts)
	if err.Error != nil {
		log.Println("QuerySnftByCollection() Select(all) err=", err.Error)
		return nil, ErrDataBase
	}
	for i, snft := range ownSnfts {
		ownSnfts[i] = snft + "m"
	}
	snftInfo := []SnftChipInfo{}
	//allsql := `SELECT * FROM nfts where ownaddr = ? and Collectcreator = ?  and collections = ? and mergetype = 1 and snft != "" and deleted_at is null`
	allsql := `SELECT * FROM nfts where snft in ? and mergetype = 1 and deleted_at is null`
	err = nft.db.Raw(allsql, ownSnfts).Find(&snftInfo)
	if err.Error != nil {
		log.Println("QuerySnftByCollection() Select(all) err=", err.Error)
		return nil, ErrDataBase
	}
	fmt.Printf("QuerySnftByCollection() min(nftaddr)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	for i, snft := range snftInfo {
		addr := common.HexToAddress(snft.Nftaddr[:len(snft.Nftaddr)-1] + "0")
		accountInfo, err := contracts.GetAccountInfo(addr, nil)
		if err != nil {
			log.Println("QuerySnftByCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
			return nil, ErrBlockchain
		}
		if accountInfo.Owner.String() == ZeroAddr {
			addr = common.HexToAddress(snft.Nftaddr[:len(snft.Nftaddr)-2] + "00")
			accountInfo, err = contracts.GetAccountInfo(addr, nil)
			if err != nil {
				log.Println("QuerySnftByCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
				return nil, ErrBlockchain
			}
			if accountInfo.Owner.String() == ZeroAddr {
				addr = common.HexToAddress(snft.Nftaddr[:len(snft.Nftaddr)-3] + "000")
				accountInfo, err = contracts.GetAccountInfo(addr, nil)
				if err != nil {
					log.Println("QuerySnftByCollection() GetAccountInfo err =", err, "NftAddress= ", addr)
					return nil, ErrBlockchain
				}
			}
		}
		snftInfo[i].Nftaddr = snft.Nftaddr[:len(snft.Nftaddr)-1] + "0"
		snftInfo[i].MergeLevel = accountInfo.MergeLevel
		snftInfo[i].MergeNumber = accountInfo.MergeNumber
		snftInfo[i].BlockNumber = accountInfo.NFTPledgedBlockNumber
	}
	GetRedisCatch().CatchQueryData("QuerySnftByCollection", queryCatchSql, &SnftChipCatch{snftInfo, uint64(0)})
	fmt.Printf("QuerySnftByCollection()  nftaddr in()  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return snftInfo, nil
}
