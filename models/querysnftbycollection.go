package models

import (
	"fmt"
	"log"
	"time"
)

func (nft NftDb) QuerySnftByCollection(ownaddr, createaddr, name, startIndex, count string) ([]SnftChipInfo, error) {
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
	sql := `select min(nftaddr) as SnftAddr, count(nftaddr) as chipcount from nfts where `
	if ownaddr != "" {
		sql = sql + `ownaddr = ? and  collectcreator = ? and collections = ?  GROUP BY snft`
		err := nft.db.Raw(sql, ownaddr, createaddr, name).Scan(&snftChips)
		if err.Error != nil {
			log.Println("QuerySnftByCollection() Select(snft) err=", err.Error)
			return nil, ErrDataBase
		}
	} else {
		sql = sql + `collectcreator = ? and collections = ?  GROUP BY snft`
		err := nft.db.Raw(sql, createaddr, name).Scan(&snftChips)
		if err.Error != nil {
			log.Println("QuerySnftByCollection() Select(snft) err=", err.Error)
			return nil, ErrDataBase
		}
	}
	for _, chip := range snftChips {
		snftAddrs = append(snftAddrs, chip.SnftAddr)
	}
	fmt.Printf("QuerySnftByCollection() min(nftaddr)  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	spendT = time.Now()
	err := nft.db.Model(&Nfts{}).Where("nftaddr in ?", snftAddrs).Scan(&snftInfo)
	if err.Error != nil {
		log.Println("QuerySnftByCollection() Scan(&snftInfo) err=", err)
		return nil, ErrDataBase
	}
	for i, chip := range snftChips {
		snftInfo[i].Chipcount = chip.Chipcount
	}
	GetRedisCatch().CatchQueryData("QuerySnftByCollection", queryCatchSql, &SnftChipCatch{snftInfo, uint64(0)})
	fmt.Printf("QuerySnftByCollection()  nftaddr in()  spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return snftInfo, nil
}
