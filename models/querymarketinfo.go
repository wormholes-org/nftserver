package models

import (
	"fmt"
	"time"
)

type Info struct {
	Tindex       string `json:"tindex"`       //Time interval, hours, days, months
	Nfttrans     int    `json:"nfttrans"`     //Transaction Information List (Number of Transactions)
	Nftsumprice  uint64 `json:"nftsumprice"`  //Exchange NFT turnover
	Nftearings   int    `json:"nftearings"`   //NFT historical returns
	Nftavprice   uint64 `json:"nftavprice"`   //Average price of NFT transactions on exchanges
	Nfthighprice int    `json:"nfthighprice"` //NFT historical price (highest transaction price)
	Nftlowprice  int    `json:"nftlowprice"`  //NFT historical price (lowest transaction price)
	//Nftgases	 	uint64		`json:"nftgases"`           //Gas consumption statistics
}

type MarketInfo struct {
	Nftliked      map[string]int `json:"nftliked"`      //The number of NFT followers
	Collectowners map[string]int `json:"collectowners"` //Collection owner address distribution
	Nftamount     int            `json:"nftamount"`     //Total number of NFTs on exchanges
	//Nftamountex	 	int  			`json:"nftamountex"`        //Number of NFTs on this exchange
	Nftowners   map[string]int `json:"nftowners"`   //Account distribution of exchange NFTs
	Nfttransamt uint64         `json:"nfttransamt"` //Exchange total NFT turnover
	Dayinfo     [24]Info       `json:"dayinfo"`     //Exchange data of the day
	Monthinfo   [31]Info       `json:"monthinfo"`   //Exchange data for the current month
	Yearinfo    [12]Info       `json:"yearinfo"`    //Exchange data for the year
}

type NFTMarketInfo struct {
	Dayinfo   [24]MInfo `json:"dayinfo"`   //Exchange data of the day
	Monthinfo [31]MInfo `json:"monthinfo"` //Exchange data for the current month
	Yearinfo  [12]MInfo `json:"yearinfo"`  //Exchange data for the year
}

type MInfo struct {
	Tindex   string `json:"tindex"`   //Time interval, hours, days, months
	Nfttrans int    `json:"nfttrans"` //Transaction Information List (Number of Transactions)
}

//Get market data
func (nft *NftDb) QueryMarketInfo() (*MarketInfo, error) {
	mInfo := MarketInfo{}
	mInfo.Nftowners = make(map[string]int)
	mInfo.Nftliked = make(map[string]int)
	mInfo.Collectowners = make(map[string]int)
	t := time.Now()
	//marketcache := t.Format("2006-01-02")
	//cerr := GetRedisCatch().GetCatchData("QueryMarketInfo", marketcache, &mInfo)
	//if cerr == nil {
	//	log.Printf("QueryUnverifiedNfts() default  time.now=%s\n", time.Now())
	//	return &mInfo, nil
	//}
	//GetRedisCatch().SetDirtyFlag([]string{"QueryMarketInfo"})
	t = t.AddDate(0, 0, -1)
	for i := 0; i < 24; i++ {
		//mInfo.Dayinfo[i].Tindex = t.Hour()
		t = t.Add(time.Hour)
		mInfo.Dayinfo[i].Tindex = t.Format("02/15:00")
	}
	var eInft []Info
	rsql := "select tindex, count(*) as nfttrans,sum(price) as nftsumprice, round(sum(price) * 0.025) as nftearings, " +
		"round(avg(price)) as nftavprice, max(price) as nfthighprice, min(price) as nftlowprice " +
		"from (select id, updated_at, DATE_FORMAT(updated_at,\"%d/%H:00\") as tindex, (price) " +
		"from trans " +
		"where updated_at >= subdate(sysdate(), 1)  " + "and updated_at <= sysdate()" +
		//"where updated_at >= subdate(sysdate(), 23)  " + "and updated_at <= subdate(sysdate(), 22)" +
		"and (selltype != \"MintNft\" and selltype != \"Error\")  " +
		"ORDER BY id )as mm group by tindex"
	err := nft.db.Raw(rsql).Scan(&eInft)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Dayinfo err=", err)
		return nil, ErrDataBase
	}
	j := 0
	for _, info := range eInft {
		for i := j; i < 24; i++ {
			if mInfo.Dayinfo[i].Tindex == info.Tindex {
				mInfo.Dayinfo[i] = info
				j = i + 1
				break
			}
		}
	}
	t = time.Now()
	t = t.AddDate(0, 0, -31)
	for i := 0; i < 31; i++ {
		//mInfo.Monthinfo[i].Tindex = t.Day()
		t = t.AddDate(0, 0, 1)
		mInfo.Monthinfo[i].Tindex = t.Format("01-02")
	}
	eInft = []Info{}
	rsql = "select tindex, count(*) as nfttrans,sum(price) as nftsumprice, round(sum(price) * 0.025) as nftearings, " +
		"round(avg(price)) as nftavprice, max(price) as nfthighprice, min(price) as nftlowprice " +
		"from (select id, updated_at, DATE_FORMAT(updated_at,\"%m-%d\") as tindex, (price) " +
		"from trans " +
		"where updated_at > subdate(sysdate(), 31) and (selltype != \"MintNft\" and selltype != \"Error\")  " +
		"ORDER BY id )as mm group by tindex"
	err = nft.db.Raw(rsql).Scan(&eInft)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Monthinfo err=", err)
		return nil, ErrDataBase
	}
	j = 0
	for _, info := range eInft {
		for i := j; i < 31; i++ {
			if mInfo.Monthinfo[i].Tindex == info.Tindex {
				mInfo.Monthinfo[i] = info
				j = i + 1
				break
			}
		}
	}
	t = time.Now()
	t = t.AddDate(-1, 0, 0)
	for i := 0; i < 12; i++ {
		t = t.AddDate(0, 1, 0)
		mInfo.Yearinfo[i].Tindex = t.Format("2006-01")

	}
	eInft = []Info{}
	rsql = "select tindex, count(*) as nfttrans,sum(price) as nftsumprice, round(sum(price) * 0.025) as nftearings, " +
		"round(avg(price)) as nftavprice, max(price) as nfthighprice, min(price) as nftlowprice " +
		"from (select id, updated_at, DATE_FORMAT(updated_at,\"%Y-%m\") as tindex, (price) " +
		"from trans " +
		"where updated_at > subdate(sysdate(), 365) and (selltype != \"MintNft\" and selltype != \"Error\")  " +
		"ORDER BY id )as mm group by tindex"
	err = nft.db.Raw(rsql).Scan(&eInft)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Yearinfo err=", err)
		return nil, ErrDataBase
	}
	j = 0
	for _, info := range eInft {
		for i := j; i < 12; i++ {
			if mInfo.Yearinfo[i].Tindex == info.Tindex {
				mInfo.Yearinfo[i] = info
				j = i + 1
				break
			}
		}
	}
	var nftcount int64
	err = nft.db.Model(Nfts{}).Count(&nftcount)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() nftcount err=", err)
		return nil, ErrDataBase
	}
	mInfo.Nftamount = int(nftcount)

	var recCount int64
	dberr := nft.db.Model(Trans{}).Where(
		"selltype != ? and selltype != ?", SellTypeMintNft.String(), SellTypeError.String()).Count(&recCount)
	if dberr.Error != nil {
		fmt.Println("QueryMarketInfo() nfttransamt err=", err)
		return nil, ErrDataBase
	}
	if recCount > 0 {
		var nfttransamt int64
		rsql = "select sum(price) as nfttransamt  from trans where (selltype != \"MintNft\" and selltype != \"Error\")"
		err = nft.db.Raw(rsql).Scan(&nfttransamt)
		if err.Error != nil {
			fmt.Println("QueryMarketInfo() nfttransamt err=", err)
			return nil, ErrDataBase
		}
		mInfo.Nfttransamt = uint64(nfttransamt)
	}

	type nftaccount struct {
		Ownaddr string
		Count   int
	}
	nftac := []nftaccount{}
	rsql = "select ownaddr, count(*) as count from nfts where snft=\"\" group by ownaddr ORDER BY count(*) desc "
	//rsql = "select ownaddr, count(*) as count from nfts group by ownaddr ORDER BY count(*) desc limit 0, 20"
	err = nft.db.Raw(rsql).Scan(&nftac)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() nfttransamt err=", err)
		return nil, ErrDataBase
	}
	for _, n := range nftac {
		mInfo.Nftowners[n.Ownaddr] = n.Count
	}
	type collects struct {
		Createaddr string
		Count      int
	}
	collect := []collects{}
	rsql = "select createaddr, count(*) as count from collects group by createaddr ORDER BY count(*) desc limit 0, 20"
	err = nft.db.Raw(rsql).Scan(&collect)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Collectowners err=", err)
		return nil, ErrDataBase
	}
	for _, n := range collect {
		mInfo.Collectowners[n.Createaddr] = n.Count
	}
	type likes struct {
		Tokenid string
		Count   int
	}
	like := []likes{}
	rsql = "select tokenid, count(*) as count from nftfavoriteds where deleted_at is null group by tokenid ORDER BY count(*) desc limit 0, 20"
	err = nft.db.Raw(rsql).Scan(&like)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() likes err=", err)
		return nil, ErrDataBase
	}
	for _, n := range like {
		mInfo.Nftliked[n.Tokenid] = n.Count
	}
	//GetRedisCatch().CatchQueryData("QueryMarketInfo", marketcache, &mInfo)

	return &mInfo, nil
}

//获取nft数据
func (nft *NftDb) GetNftMarketInfo() (*NFTMarketInfo, error) {
	mInfo := NFTMarketInfo{}

	t := time.Now()
	//marketcache := t.Format("2006-01-02")
	//cerr := GetRedisCatch().GetCatchData("GetNftMarketInfo", marketcache, &mInfo)
	//if cerr == nil {
	//	log.Printf("GetNftMarketInfo() default  time.now=%s\n", time.Now())
	//	return &mInfo, nil
	//}
	//GetRedisCatch().SetDirtyFlag([]string{"QueryMarketInfo"})
	t = t.AddDate(0, 0, -31)
	for i := 0; i < 31; i++ {
		t = t.AddDate(0, 0, 1)
		mInfo.Monthinfo[i].Tindex = t.Format("01-02")
	}
	mInft := []MInfo{}
	rsql := "select tindex,count(*) as nfttrans " +
		"from (select id, created_at, DATE_FORMAT(created_at,\"%m-%d\") as tindex " +
		"from nfts where created_at >= subdate(sysdate(), 31)  and snft =\"\" and deleted_at is null " +
		"ORDER BY id ) as mm group by tindex"
	err := nft.db.Raw(rsql).Scan(&mInft)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Monthinfo err=", err)
		return nil, ErrDataBase
	}

	j := 0
	for _, info := range mInft {
		for i := j; i < 31; i++ {
			if mInfo.Monthinfo[i].Tindex == info.Tindex {
				mInfo.Monthinfo[i] = info
				j = i + 1
				break
			}
		}
	}
	mInft = []MInfo{}

	t = time.Now()
	t = t.AddDate(0, 0, -1)
	for i := 0; i < 24; i++ {
		t = t.Add(time.Hour)
		//mInfo.Dayinfo[i].Tindex = t.Hour()
		mInfo.Dayinfo[i].Tindex = t.Format("02/15:00")

	}
	rsql = "select tindex,count(*) as nfttrans " +
		"from (select id, created_at, DATE_FORMAT(created_at,\"%d/%H:00\") as tindex " +
		"from nfts where created_at >= subdate(sysdate(), 1)  and snft =\"\" and deleted_at is null " +
		"ORDER BY id ) as mm group by tindex"
	err = nft.db.Raw(rsql).Scan(&mInft)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Dayinfo err=", err)
		return nil, ErrDataBase
	}
	j = 0
	for _, info := range mInft {
		for i := j; i < 24; i++ {
			if mInfo.Dayinfo[i].Tindex == info.Tindex {
				mInfo.Dayinfo[i] = info
				j = i + 1
				break
			}
		}
	}
	mInft = []MInfo{}
	t = time.Now()
	t = t.AddDate(-1, 0, 0)
	for i := 0; i < 12; i++ {
		t = t.AddDate(0, 1, 0)
		//mInfo.Yearinfo[i].Tindex = int(t.Month())
		mInfo.Yearinfo[i].Tindex = t.Format("2006-01")

	}

	rsql = "select tindex,count(*) as nfttrans " +
		"from (select id, created_at, DATE_FORMAT(created_at,\"%Y-%m\") as tindex " +
		"from nfts where created_at >= subdate(sysdate(), 365)  and snft =\"\" and deleted_at is null " +
		"ORDER BY id ) as mm group by tindex"
	err = nft.db.Raw(rsql).Scan(&mInft)
	if err.Error != nil {
		fmt.Println("QueryMarketInfo() Dayinfo err=", err)
		return nil, ErrDataBase
	}
	j = 0
	for _, info := range mInft {
		for i := j; i < 12; i++ {
			if mInfo.Yearinfo[i].Tindex == info.Tindex {
				mInfo.Yearinfo[i] = info
				j = i + 1
				break
			}
		}
	}
	//GetRedisCatch().CatchQueryData("GetNftMarketInfo", marketcache, &mInfo)

	return &mInfo, nil
}
