package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

/*var NftFields = "Ownaddr, Md5, Name, Desc, Meta, Url, Contract, Tokenid, Count, " +
"Ownaddr, Md5, Name, Desc, Meta, Url, Contract, Tokenid, Count, " +
"Approve, Categories, Collectcreator, Collections, Image, Hide, " +
"Signdata, Createaddr, Verifyaddr, Currency, Price, Royalty, " +
"Paychan, TransCur, Transprice, Transtime, Createdate, Favorited, " +
"Transcnt, Transamt, Verified, Verifiedtime, Selltype, Mintstate, " +
"Extend"
*/
var NftFields = "ownaddr, md5, name, desc, meta, url, contract, tokenid, count, " +
	"ownaddr, md5, name, desc, meta, url, contract, tokenid, count, " +
	"approve, categories, collectcreator, collections, image, hide, " +
	"signdata, createaddr, verifyaddr, currency, price, royalty, " +
	"paychan, transcur, transprice, transtime, createdate, favorited, " +
	"transcnt, transamt, verified, verifiedtime, selltype, mintstate, " +
	"extend"

/*var TransFields = "auctioni, contract, createaddr, " +
"fromaddr, toaddr, tradesig, signdata, txhash, tokenid, " +
"count, transtime, paychan, currency, price, name, desc, " +
"meta, selltype, error"
*/
var TransFields = "auctioni, contract, createaddr, " +
	"fromaddr, toaddr, tradesig, signdata, txhash, tokenid, " +
	"count, transtime, paychan, currency, price, name, desc, " +
	"meta, selltype, error, sellprice"

type MarketTradingCatch struct {
	MarketTrading []TradingHistory
	Total         int
}

func (nft NftDb) QueryMarketTradingHistory(nfttype string, filter []StQueryField, sort []StSortField,
	start_index string, count string) ([]TradingHistory, int, error) {
	var tranRecs []Trans
	var recCount int64
	var queryWhere string
	var orderBy string

	spendT := time.Now()
	nfttype = strings.ToLower(nfttype)
	sql := "SELECT trans.*, trans.price as sellprice FROM trans LEFT JOIN nfts ON trans.contract = nfts.contract AND trans.tokenid = nfts.tokenid"
	countSql := "SELECT count(*) FROM trans "
	var mergetypeflag string
	if len(filter) > 0 {
		for k, v := range filter {
			if strings.Contains(TransFields, strings.ToLower(v.Field)) {
				filter[k].Field = "trans." + filter[k].Field
				if v.Field == "sellprice" {
					filter[k].Field = "trans.price"
				}
			} else if strings.Contains(NftFields, strings.ToLower(v.Field)) {
				//filter[k].Field = "nfts." + filter[k].Field
				//filter[k].Field = "trans.transtime"
				if strings.ToLower(v.Field) == "collections" || strings.ToLower(v.Field) == "collectcreator" {
					filter[k].Field = "nfts." + filter[k].Field
					countSql = "SELECT count(*) FROM trans LEFT JOIN nfts ON trans.contract = nfts.contract AND trans.tokenid = nfts.tokenid"
				} else {
					filter[k].Field = "trans.transtime"
				}

			}
			if filter[k].Field == "mergetype" {
				mergetypeflag = filter[k].Value
			}
		}
		queryWhere = nft.joinFilters(filter)
		if mergetypeflag != "" {
			switch mergetypeflag {
			case "0":
				queryWhere = strings.Replace(queryWhere, "mergetype='0'", "mergetype = mergelevel and mergelevel = 0", -1)
			case "1":
				queryWhere = strings.Replace(queryWhere, "mergetype='0'", "mergetype = mergelevel and mergelevel = 1", -1)
			case "2":
				queryWhere = strings.Replace(queryWhere, "mergetype='0'", "mergetype = mergelevel and mergelevel = 2", -1)
			case "3":
				queryWhere = strings.Replace(queryWhere, "mergetype='0'", "mergetype = mergelevel and mergelevel = 3", -1)
			}
		}
		if len(queryWhere) > 0 {
			sql = sql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' and" + queryWhere
			countSql = countSql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' and" + queryWhere
		} else {
			sql = sql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
			countSql = countSql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
		}
	} else {
		sql = sql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
		countSql = countSql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
	}
	if nfttype == "nft" || nfttype == "snft" {
		sql = sql + " and trans.nfttype = " + "'" + nfttype + "' "
		countSql = countSql + " and trans.nfttype = " + "'" + nfttype + "' "
	}
	if len(sort) > 0 {
		for k, v := range sort {
			if k > 0 {
				orderBy = orderBy + ", "
			}
			if strings.Contains(TransFields, strings.ToLower(v.By)) {
				if strings.ToLower(v.By) == "sellprice" {
					orderBy += "trans.price" + " " + v.Order
				} else {
					orderBy += "trans." + v.By + " " + v.Order
				}
			} else if strings.Contains(NftFields, v.By) {
				orderBy += "nfts." + v.By + " " + v.Order
			}
		}
	}
	if len(orderBy) > 0 {
		orderBy = orderBy + ", trans.id desc"
	} else {
		orderBy = "trans.id desc"
	}
	sql = sql + " order by " + orderBy
	//countSql = countSql + " order by " + orderBy
	fmt.Println("QueryMarketTradingHistory() sql=", sql)
	fmt.Println("QueryMarketTradingHistory() countSql=", countSql)
	if len(start_index) > 0 {
		sql = sql + " limit " + start_index + ", " + count
	}
	queryCatchSql := sql
	nftCatch := MarketTradingCatch{}
	cerr := GetRedisCatch().GetCatchData("QueryMarketTradingHistory", queryCatchSql, &nftCatch)
	if cerr == nil {
		log.Printf("QueryMarketTradingHistory() catch spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
		return nftCatch.MarketTrading, nftCatch.Total, nil
	}
	fmt.Println("QueryMarketTradingHistory() sql=", sql)
	err := nft.db.Raw(sql).Scan(&tranRecs)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return nil, 0, ErrDataBase
	}
	fmt.Println("QueryMarketTradingHistory() countSql=", countSql)
	err = nft.db.Raw(countSql).Scan(&recCount)
	if err.Error != nil {
		return nil, 0, ErrDataBase
	}

	trans := make([]TradingHistory, 0, 20)
	for i := 0; i < len(tranRecs); i++ {
		var tran TradingHistory
		tran.NftContractAddr = tranRecs[i].Contract
		tran.NftTokenId = tranRecs[i].Tokenid
		tran.Nftaddr = tranRecs[i].Nftaddr
		tran.Url = tranRecs[i].Url
		tran.NftName = tranRecs[i].Name
		tran.Price = tranRecs[i].Price
		tran.Count = 1
		tran.From = tranRecs[i].Fromaddr
		tran.To = tranRecs[i].Toaddr
		tran.Date = tranRecs[i].Transtime
		tran.Selltype = tranRecs[i].Selltype
		tran.Txhash = tranRecs[i].Txhash
		trans = append(trans, tran)
	}
	GetRedisCatch().CatchQueryData("QueryMarketTradingHistory", queryCatchSql, &MarketTradingCatch{trans, int(recCount)})
	fmt.Printf("QueryMarketTradingHistory() spend time=%s time.now=%s\n", time.Now().Sub(spendT), time.Now())
	return trans, int(recCount), nil
}
