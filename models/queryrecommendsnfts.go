package models

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

var (
	L0AveragePrice uint64 = 30000000
	L1AveragePrice uint64 = 143000000
	L2AveragePrice uint64 = 271000000
	L3AveragePrice uint64 = 650000000
)

const (
	RecommendTypeSell   = "sell"
	RecommendTypeBuying = "buying"
)

type RecommendResp struct {
	Ownaddr   string `json:"ownaddr"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Url       string `json:"source_url"`
	Contract  string `json:"nft_contract_addr"`
	Tokenid   string `json:"nft_token_id"`
	Nftaddr   string `json:"nft_address"`
	TransType string `json:"trans_type"`
	Snftnum   int64  `json:"snftnum"`
}

type RecommendBuyingSell struct {
	Sell   []RecommendResp
	Buying []RecommendResp
}

func (nft *NftDb) QueryRecommendSnfts(userAddr string) (RecommendBuyingSell, error) {

	var recommend RecommendBuyingSell
	var recomwg sync.WaitGroup
	cerr := GetRedisCatch().GetCatchData("QueryRecommendSnfts", userAddr, &recommend)
	if cerr == nil {
		log.Printf("QueryRecommendSnfts()  default spend  time.now=%s\n", time.Now())
		return recommend, nil
	}
	var err1 error
	recomwg.Add(1)
	go func(err1 error) {
		defer recomwg.Done()
		buyrecomm, err1 := nft.BuyingRecommend(userAddr)
		if err1 != nil {
			log.Println("QueryRecommendSnfts BuyingRecommend err=", err1)
			//return RecommendBuyingSell{}, err
		}
		recommend.Buying = buyrecomm
	}(err1)
	//recomwg.Add(1)
	//go func(err2 error) {
	//	defer recomwg.Done()
	//	sellrecomm, err2 := nft.SellRecommend(userAddr)
	//	if err2 != nil {
	//		log.Println("QueryRecommendSnfts SellRecommend err=", err2)
	//	}
	//	recommend.Sell = sellrecomm
	//}(err2)
	recomwg.Wait()

	if err1 != nil {
		log.Println("QueryRecommendSnfts BuyingRecommend err=", err1)
		return RecommendBuyingSell{}, err1
	}
	//if err2 != nil {
	//	log.Println("QueryRecommendSnfts BuyingRecommend err=", err2)
	//	return RecommendBuyingSell{}, err2
	//}
	GetRedisCatch().CatchQueryData("QueryRecommendSnfts", userAddr, &recommend)

	return recommend, nil
}

type RecommendSnftAddr struct {
	Addr  string `json:"addr"`
	Csnft int64  `json:"csnft"`
}

func (nft *NftDb) BuyingRecommend(useraddr string) ([]RecommendResp, error) {
	if useraddr == "" {
		var selectaddr []string
		//sql := `select addr from (SELECT left(nftaddr,41) as addr FROM auctions where deleted_at is null) as mm group by addr order by count(*) desc limit 0, 5`
		//sql := `SELECT nftaddr as addr FROM auctions where deleted_at is null group by addr order by count(*) desc limit 0, 5`
		sql := `SELECT left(nftaddr,41) as addr  FROM auctions where deleted_at is null 
		and locate('m',nftaddr)=0 and selltype = ? GROUP BY addr order by count(*) desc limit 0, 5`
		db := nft.db.Raw(sql, SellTypeFixPrice.String()).Scan(&selectaddr)
		if db.Error != nil {
			log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
			return nil, db.Error
		}
		var snftaddr []string
		for _, value := range selectaddr {
			snftaddr = append(snftaddr, value+"0")
			//if value[len(value)-1:] != "m" {
			//	snftaddr = append(snftaddr, value)
			//}
			if len(snftaddr) == 3 {
				break
			}
		}
		if len(snftaddr) != 3 {
			sql = `select addr from (SELECT *,left(nftaddr,41) as addr FROM trans where nftaddr <> '' ) 
		as mm GROUP BY  addr order by count(*) desc`
			db = nft.db.Raw(sql).Scan(&selectaddr)
			if db.Error != nil {
				log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
				return nil, db.Error
			}
			for _, value := range selectaddr {
				if value[len(value)-1:] != "m" {
					snftaddr = append(snftaddr, value+"0")
				}
				if len(snftaddr) == 3 {
					break
				}
			}
			if len(snftaddr) != 3 {
				sql = `select addr from (SELECT left(nftaddr,41) as addr FROM auctions where deleted_at is null) as mm group by addr order by count(*) desc`
				db = nft.db.Raw(sql).Scan(&selectaddr)
				if db.Error != nil {
					log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
					return nil, db.Error
				}
				for _, value := range selectaddr {
					if value[len(value)-1:] != "m" {
						snftaddr = append(snftaddr, value+"0")
					}
					if len(snftaddr) == 3 {
						break
					}
				}
				if len(snftaddr) < 3 {
					var maxCount int64
					countSql := `select max(id) from nfts`
					db = nft.db.Raw(countSql).Scan(&maxCount)
					rand.Seed(time.Now().UnixNano())
					for {
						index := rand.Intn(int(maxCount))

						var nftRec Nfts
						db = nft.db.Where("id = ? ", index).First(&nftRec)
						if db.Error != nil {
							//nd.Close()
							log.Println("SellRecommend() index=", index, "First(&nftRec) err = ", db.Error)
							continue
						}

						if nftRec.Nftaddr != "" {
							if nftRec.Nftaddr[41:] != "m" {
								snftaddr = append(snftaddr, nftRec.Nftaddr)
							}
							if len(snftaddr) == 3 {
								break
							}
						}
					}
				}
			}

		}

		findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeBuying)
		if err != nil {
			log.Println("SellRecommend FindSnftRecommend err=", err)
			return nil, err
		}
		return findrecom, nil
	} else {
		var selectaddr []RecommendSnftAddr
		weight := make(map[string]int64)
		var snftaddr []string
		//sql := `select addr from (select left(nftaddr,41) as addr  from nfts where ownaddr = ?
		//and mergetype=0 and mergelevel =0 and exchange =0) as mm group by addr order by count(*) limit 0, 5`
		sql := `select left(nftaddr,41) as addr, count(*) as csnft from nfts where ownaddr = ? and mergetype=0 and mergelevel =mergetype and exchange =0 group by addr order by csnft desc limit 0, 3 `
		db := nft.db.Raw(sql, useraddr).Scan(&selectaddr)
		if db.Error != nil {
			log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
			return nil, db.Error
		}
		for _, value := range selectaddr {
			//snftaddr = append(snftaddr, value.addr+"0")
			//if len(snftaddr) == 5 {
			//	break
			//}
			weight[value.Addr] += value.Csnft * 2
		}
		var auctionaddr []RecommendSnftAddr
		sql = `SELECT left(nftaddr,41) as addr ,count(*) as csnft FROM auctions where deleted_at is null 
		and locate('m',nftaddr)=0 and selltype = ? and ownaddr != ? GROUP BY addr order by count(*) desc limit 0, 3`
		db = nft.db.Raw(sql, SellTypeFixPrice.String(), useraddr).Scan(&auctionaddr)
		if db.Error != nil {
			log.Println("SellRecommend shards recommended for sale err=", db.Error)
			return nil, db.Error
		}
		for _, value := range auctionaddr {
			//snftaddr = append(snftaddr, value+"0")
			weight[value.Addr] += value.Csnft

		}
		type peroson struct {
			Addr  string
			Count int64
		}
		var lstPerson []peroson
		for k, v := range weight {
			lstPerson = append(lstPerson, peroson{k, v})
		}
		sort.Slice(lstPerson, func(i, j int) bool {
			return lstPerson[i].Count > lstPerson[j].Count
		})
		for _, vaule := range lstPerson {
			snftaddr = append(snftaddr, vaule.Addr+"0")
		}
		if len(snftaddr) > 3 {
			snftaddr = snftaddr[:3]
			findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeBuying)
			if err != nil {
				log.Println("SellRecommend FindSnftRecommend err=", err)
				return nil, err
			}
			return findrecom, nil
		}
		switch len(snftaddr) {
		case 0:
			recom, err := nft.BuyingRecommend("")
			if err != nil {
				log.Println("BuyingRecommend BuyingRecommend query 0 chip err=", err)
				return nil, err
			}
			return recom, nil
		case 1, 2:
			findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeBuying)
			if err != nil {
				log.Println("BuyingRecommend FindSnftRecommend err=", err)
				return nil, err
			}
			return findrecom, nil
		default:
			snftaddr = snftaddr[:3]
			findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeBuying)
			if err != nil {
				log.Println("SellRecommend FindSnftRecommend err=", err)
				return nil, err
			}
			return findrecom, nil
		}
	}
}

func (nft *NftDb) SellRecommend(useraddr string) ([]RecommendResp, error) {
	if useraddr == "" {
		//var selectaddr []string
		//sql := `select addr from (SELECT left(nftaddr,41) as addr FROM auctions where deleted_at is null) as mm group by addr order by count(*) desc`
		//db := nft.db.Raw(sql).Scan(&selectaddr)
		//if db.Error != nil {
		//	log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
		//	return nil, db.Error
		//}
		//var snftaddr []string
		//for _, value := range selectaddr {
		//	snftaddr = append(snftaddr, value+"0")
		//}
		//if len(snftaddr) > 5 {
		//	snftaddr = snftaddr[:5]
		//}
		//findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeSell)
		//if err != nil {
		//	log.Println("SellRecommend FindSnftRecommend err=", err)
		//	return nil, err
		//}
		//return findrecom, nil

		var selectaddr []string
		sql := `select addr from (SELECT *,left(nftaddr,41) as addr FROM trans where nftaddr <> '' ) 
		as mm GROUP BY  addr order by count(*) desc`
		db := nft.db.Raw(sql).Scan(&selectaddr)
		if db.Error != nil {
			log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
			return nil, db.Error
		}
		var snftaddr []string
		for _, value := range selectaddr {
			//snftaddr = append(snftaddr, value+"0")
			if value[len(value)-1:] != "m" {
				snftaddr = append(snftaddr, value+"0")
			}
		}
		if len(snftaddr) >= 5 {
			snftaddr = snftaddr[:5]
		} else {
			sql = `select addr from (SELECT left(nftaddr,41) as addr FROM auctions where deleted_at is null) as mm group by addr order by count(*) desc`
			db = nft.db.Raw(sql).Scan(&selectaddr)
			if db.Error != nil {
				log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
				return nil, db.Error
			}
			for _, value := range selectaddr {
				if len(snftaddr) == 5 {
					break
				}
				if value[len(value)-1:] != "m" {
					snftaddr = append(snftaddr, value+"0")
				}
			}
			if len(snftaddr) < 5 {
				var maxCount int64
				countSql := `select max(id) from nfts`
				db = nft.db.Raw(countSql).Scan(&maxCount)
				rand.Seed(time.Now().UnixNano())
				for {
					index := rand.Intn(int(maxCount))

					var nftRec Nfts
					db = nft.db.Where("id = ? ", index).First(&nftRec)
					if db.Error != nil {
						//nd.Close()
						log.Println("SellRecommend() index=", index, "First(&nftRec) err = ", db.Error)
						continue
					}

					if nftRec.Nftaddr != "" {
						if nftRec.Nftaddr[41:] != "m" {
							snftaddr = append(snftaddr, nftRec.Nftaddr)
						}
						if len(snftaddr) == 5 {
							break
						}
					}
				}
			}
		}
		findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeSell)
		if err != nil {
			log.Println("SellRecommend FindSnftRecommend err=", err)
			return nil, err
		}
		return findrecom, nil
	} else {
		var selectaddr []string
		sql := `select addr from (select left(nftaddr,41) as addr  from nfts where ownaddr = ?
		and mergetype=0 and mergelevel =0 and exchange =0) as mm group by addr order by count(*)`
		db := nft.db.Raw(sql, useraddr).Scan(&selectaddr)
		if db.Error != nil {
			log.Println("SellRecommend query snft when the user address is empty err=", db.Error)
			return nil, db.Error
		}
		var snftaddr []string
		for _, value := range selectaddr {
			snftaddr = append(snftaddr, value+"0")
		}
		switch len(snftaddr) {
		case 0:
			recom, err := nft.SellRecommend("")
			if err != nil {
				log.Println("SellRecommend SellRecommend query 0 chip err=", db.Error)
				return nil, err
			}
			return recom, nil
		case 1, 2, 3, 4, 5:
			findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeSell)
			if err != nil {
				log.Println("SellRecommend FindSnftRecommend err=", err)
				return nil, err
			}
			return findrecom, nil
		default:
			snftaddr = snftaddr[:5]
			findrecom, err := nft.FindSnftRecommend(snftaddr, RecommendTypeSell)
			if err != nil {
				log.Println("SellRecommend FindSnftRecommend err=", err)
				return nil, err
			}
			return findrecom, nil
		}
	}
}

func (nft *NftDb) FindSnftRecommend(addr []string, types string) ([]RecommendResp, error) {

	var recommend []RecommendResp
	sql := `select nftaddr,ownaddr,url,contract,tokenid,` + "`name`,`desc`" + `,count(*) as snftnum from (
	select n1.*,n2.nftaddr as addr from nfts n1 
	left join nfts n2 on left(n1.nftaddr,41) =n2.snft and n2.exchange =0 
	where n1.nftaddr = ? ) as mm GROUP BY nftaddr,ownaddr,url,contract,tokenid,` + "`name`,`desc`"
	for _, singeaddr := range addr {
		var singerecommend RecommendResp
		db := nft.db.Raw(sql, singeaddr).First(&singerecommend)
		if db.Error != nil {
			log.Println("FindSnftRecommend query snft data err=", db.Error)
			return nil, db.Error
		}
		recommend = append(recommend, singerecommend)
	}

	for key := range recommend {
		recommend[key].TransType = types
	}
	return recommend, nil
}

func (nft *NftDb) SetRecommendSnftCatch() error {
	var user []Users
	GetRedisCatch().SetDirtyFlag(RecommendSnft)

	db := nft.db.Model(&Users{}).Where("deleted_at is null").Find(&user)
	if db.Error != nil {
		log.Println("SetRecommendSnftCatch query user data err=", db.Error)
		return db.Error
	}
	_, err := nft.QueryRecommendSnfts("")
	if err != nil {
		log.Println("SetRecommendSnftCatch set nil user catch err=", err)
		return err
	}
	for _, singe := range user {
		_, err = nft.QueryRecommendSnfts(singe.Useraddr)
		if err != nil {
			log.Println("SetRecommendSnftCatch set  user =", singe.Useraddr, "  catch err=", err)
			return err
		}
	}
	return nil
}

func RecommendTimeProc(sqldsn string) {
	ticker := time.NewTicker(time.Minute * 15)
	nd, err := NewNftDb(sqldsn)
	if err != nil {
		log.Printf("RecommendTimeProc() connect database err = %s\n", err)
	}
	for {
		select {
		case <-ticker.C:
			err = nd.SetRecommendSnftCatch()
			if err != nil {
				fmt.Println("RecommendTimeProc : err = ", err)
				continue
			}
		}
	}
	nd.Close()

}
