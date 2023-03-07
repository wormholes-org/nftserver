package main

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	_ "github.com/nftexchange/nftserver/routers"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

func TimeProc(sqldsn string) {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			nd, err := models.NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("connect database err = %s\n", err)
				continue
			}
			CallContractsNew(nd)
			nd.Close()
		}
	}
}

func clearMergeAuction(nft *models.NftDb, auctionRec *models.Auction) {
	nft.GetDB().Transaction(func(tx *gorm.DB) error {
		nftrecord := models.Nfts{}
		err := tx.Model(&models.Nfts{}).Where("contract = ? AND tokenid =?",
			auctionRec.Contract, auctionRec.Tokenid).First(&nftrecord)
		if err.Error != nil {
			log.Println("clearMergeAuction() update record err=", err.Error)
			return err.Error
		}
		if nftrecord.Mergetype != nftrecord.Mergelevel || nftrecord.Exchange == 1 || nftrecord.Pledgestate == models.Pledge.String() {
			err = nft.GetDB().Model(&models.Bidding{}).Where("auctionid = ?", auctionRec.ID).Delete(&models.Bidding{})
			if err.Error != nil {
				log.Println("clearMergeAuction() delete bidding record err=", err.Error)
				return err.Error
			}
			err = tx.Model(&models.Auction{}).Where("id = ? ", auctionRec.ID).Delete(&models.Auction{})
			if err.Error != nil {
				log.Println("clearMergeAuction() delete auction record err=", err.Error)
				return err.Error
			}
			//nftrecord := models.Nfts{}
			//nftrecord.Selltype = models.SellTypeNotSale.String()
			//err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nftrecord)
			//if err.Error != nil {
			//	fmt.Println("clearMergeAuction() update record err=", err.Error)
			//	return err.Error
			//}
			nfttab := map[string]interface{}{
				"Selltype":    models.SellTypeNotSale.String(),
				"Sellprice":   0,
				"Offernum":    0,
				"Maxbidprice": 0,
			}
			err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nfttab)
			if err.Error != nil {
				fmt.Println("ClearAuction() update record err=", err.Error)
				return err.Error
			}
		}
		return nil
	})
}

func CallContracts(nft *models.NftDb) {
	spendT := time.Now()
	fmt.Println(time.Now().String()[:20], "TimeProc begin+++++++++++++++++++++++++.")
	rows, err := nft.GetDB().Model(&models.Auction{}).Rows()
	if err != nil {
		fmt.Println("TimeProc() Rows err=", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var auctionRec models.Auction
		nft.GetDB().ScanRows(rows, &auctionRec)
		if auctionRec.Selltype == models.SellTypeHighestBid.String() &&
			auctionRec.SellState == models.SellStateStart.String() &&
			time.Now().Unix() >= auctionRec.Enddate {
			var bidRecs []models.Bidding
			err := nft.GetDB().Order("price desc").Where("Auctionid = ?", auctionRec.ID).Find(&bidRecs)
			if err.Error != nil || err.RowsAffected == 0 {
				if err.Error == gorm.ErrRecordNotFound || err.RowsAffected == 0 {
					nft.GetDB().Transaction(func(tx *gorm.DB) error {
						nftrecord := models.Nfts{}
						nftrecord.Selltype = models.SellTypeNotSale.String()
						err = tx.Model(&models.Nfts{}).Where("contract = ? AND tokenid =? AND ownaddr = ?",
							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(&nftrecord)
						if err.Error != nil {
							fmt.Println("TimeProc() update record err=", err.Error)
							return err.Error
						}
						err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Delete(&models.Auction{})
						if err.Error != nil {
							fmt.Println("TimeProc() delete auction record err=", err.Error)
							return err.Error
						}
						return nil
					})
				}
				continue
			}
			var bidRec *models.Bidding
			for _, rec := range bidRecs {
				valid, _, cerr := models.WormsAmountValid(rec.Price, rec.Bidaddr)
				if cerr != nil {
					continue
				}
				if !valid {
					err = nft.GetDB().Model(&models.Bidding{}).Where("id = ?", rec.ID).Delete(&models.Bidding{})
					if err.Error != nil {
						fmt.Println("TimeProc() delete bidding record err=", err.Error)
						continue
					}
					continue
				}
				bidRec = &rec
				break
			}
			if bidRec == nil {
				continue
			}
			//fmt.Println("TimeProc() bidRecs.Price=", bidRecs.Price, "controllers.Lowprice=",
			//	models.Lowprice,"auctionRec.Startprice=", auctionRec.Startprice, "valid=", valid)
			if bidRec.Price >= models.Lowprice && bidRec.Price >= auctionRec.Startprice {
				var nftrecord models.Nfts
				err := nft.GetDB().Where("contract = ? AND tokenid = ? AND ownaddr = ?",
					auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).First(&nftrecord)
				if err.Error != nil {
					fmt.Println("TimeProc() get nftrecord err= ", err.Error)
					continue
				}
				if nftrecord.Mintstate == models.Minted.String() {
					fmt.Println("TimeProc() OwnAndAprove() == TRUE ")
					price := strconv.FormatUint(bidRec.Price, 10) + "000000000"
					fmt.Println("TimeProc() call  ethhelper.Auction()")
					fmt.Println("TimeProc() auction.id=", auctionRec.ID)
					fmt.Println("TimeProc() auctionRec.Ownaddr=", auctionRec.Ownaddr)
					fmt.Println("TimeProc() bidRecs.Bidaddr=", bidRec.Bidaddr)
					fmt.Println("TimeProc() auctionRec.Contract=", auctionRec.Contract)
					fmt.Println("TimeProc() auctionRec.Tokenid=", auctionRec.Tokenid)
					fmt.Println("TimeProc() auctionRec.Count=", auctionRec.Count)
					fmt.Println("TimeProc() price=", price)
					fmt.Println("TimeProc() bidRecs.Tradesig=", bidRec.Tradesig)

					buyer := contracts.Buyer{}
					err := json.Unmarshal([]byte(bidRec.Tradesig), &buyer)
					if err != nil {
						fmt.Println("TimeProc() ethhelper.Auction() err=", err)
						continue
					}
					seller := contracts.Seller1{}
					err = json.Unmarshal([]byte(auctionRec.Tradesig), &seller)
					if err != nil {
						fmt.Println("TimeProc() ethhelper.Auction() err=", err)
						continue
					}
					txhash, err := contracts.AuthExchangeTrans(seller, buyer, models.ExchangerAuth, contracts.SuperAdminAddr)
					if err != nil {
						fmt.Println("TimeProc() contracts.AuthExchangeTrans() err=", err)
						continue
					}
					fmt.Println("TimeProc() contracts.AuthExchangeTrans() Ok tx.hash=", txhash)
					nft.GetDB().Transaction(func(tx *gorm.DB) error {
						/*err = nft.GetDB().Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
							auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
						if err.Error != nil {
							fmt.Println("TimeProc() delete bid record err=", err.Error)
							return err.Error
						}*/
						auctRec := models.Auction{}
						//auctRec.Selltype = models.SellTypeWaitSale.String()
						auctRec.SellState = models.SellStateWait.String()
						auctRec.Price = bidRec.Price
						err := nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
						if err.Error != nil {
							fmt.Println("TimeProc() update auction record err=", err.Error)
							return err.Error
						}
						return nil
					})
				} else {
					price := strconv.FormatUint(bidRec.Price, 10) + "000000000"
					fmt.Println("TimeProc() call ethhelper.AuctionAndMint()")
					fmt.Println("TimeProc() auction.id=", auctionRec.ID)
					fmt.Println("TimeProc() auctionRec.Ownaddr=", auctionRec.Ownaddr)
					fmt.Println("TimeProc() bidRecs.Bidaddr=", bidRec.Bidaddr)
					fmt.Println("TimeProc() auctionRec.Contract=", auctionRec.Contract)
					fmt.Println("TimeProc() auctionRec.Tokenid=", auctionRec.Tokenid)
					fmt.Println("TimeProc() price=", price)
					fmt.Println("TimeProc() nftrecord.Meta=", nftrecord.Meta)
					fmt.Println("TimeProc() bidRecs.Tradesig=", bidRec.Tradesig)
					Royalty := strconv.Itoa(nftrecord.Royalty)
					count := strconv.Itoa(nftrecord.Count)
					fmt.Println("TimeProc() Royalty=", Royalty)
					fmt.Println("TimeProc() count=", count)
					fmt.Println("TimeProc() params=", auctionRec.Ownaddr, bidRec.Bidaddr, auctionRec.Contract,
						auctionRec.Tokenid, price, count, Royalty, nftrecord.Meta, bidRec.Tradesig)

					seller := contracts.Seller2{}
					err := json.Unmarshal([]byte(auctionRec.Tradesig), &seller)
					if err != nil {
						fmt.Println("TimeProc() Unmarshal() err=", err)
						continue
					}
					buyer := contracts.Buyer1{}
					err = json.Unmarshal([]byte(bidRec.Tradesig), &buyer)
					if err != nil {
						fmt.Println("TimeProc() Unmarshal() err=", err)
						continue
					}

					txhash, err := contracts.AuthExchangerMint(seller, buyer, models.ExchangerAuth, contracts.SuperAdminAddr)
					if err != nil {
						fmt.Println("TimeProc() ExchangerMint() err=", err)
						continue
					}
					fmt.Println("TimeProc() contracts.AuthExchangerMint() Ok contracts.AuthExchangerMint txhash", txhash)
					nft.GetDB().Transaction(func(tx *gorm.DB) error {
						/*err = nft.GetDB().Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
							auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
						if err.Error != nil {
							fmt.Println("TimeProc() AuctionAndMint delete bid record err=", err.Error)
							return err.Error
						}*/
						auctRec := models.Auction{}
						//auctRec.Selltype = models.SellTypeWaitSale.String()
						auctRec.SellState = models.SellStateWait.String()
						auctRec.Price = bidRec.Price
						err := nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
						if err.Error != nil {
							fmt.Println("TimeProc() update AuctionAndMint record err=", err.Error)
							return err.Error
						}
						return nil
					})
				}
			}
			/*else {
				fmt.Println("TimeProc() auth balance error.")
				var nftrecord models.Nfts
				err := nft.GetDB().Where("contract = ? AND tokenid = ? AND Ownaddr = ?",
					auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).First(&nftrecord)
				if err.Error != nil {
					fmt.Println("TimeProc() get nftrecord err= ", err.Error )
					continue
				}
				nft.GetDB().Transaction(func(tx *gorm.DB) error {
					trans := models.Trans{}
					trans.Auctionid = auctionRec.ID
					trans.Contract = auctionRec.Contract
					trans.Fromaddr = auctionRec.Ownaddr
					trans.Toaddr = bidRecs.Bidaddr
					trans.Signdata = bidRecs.Signdata
					trans.Tokenid = auctionRec.Tokenid
					trans.Paychan = auctionRec.Paychan
					trans.Currency = auctionRec.Currency
					trans.Price = bidRecs.Price
					trans.Transtime = time.Now().Unix()
					trans.Selltype = models.SellTypeError.String()
					trans.Error = auctionRec.Selltype + "," + errmsg
					trans.Name = nftrecord.Name
					trans.Meta = nftrecord.Meta
					trans.Desc = nftrecord.Desc
					err := tx.Model(&trans).Create(&trans)
					if err.Error != nil {
						fmt.Println("TimeProc() error create trans record err=", err.Error)
						return err.Error
					}
					//auctRec := models.Auction{}
					////auctRec.Selltype = models.SellTypeWaitSale.String()
					//auctRec.SellState = models.SellStateWait.String()
					//auctRec.Price = bidRecs.Price
					//err = nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ?",
					//	auctionRec.Contract, auctionRec.Tokenid).Updates(auctRec)
					//if err.Error != nil {
					//	fmt.Println("TimeProc() error  update auction record err=", err.Error)
					//	return err.Error
					//}
					nftrecord := models.Nfts{}
					nftrecord.Selltype = models.SellTypeNotSale.String()
					err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nftrecord)
					if err.Error != nil {
						fmt.Println("TimeProc() update record err=", err.Error)
						return err.Error
					}
					err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
						auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Delete(&models.Auction{})
					if err.Error != nil {
						fmt.Println("BuyResult() delete auction record err=", err.Error)
						return err.Error
					}
					err = tx.Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
						auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
					if err.Error != nil {
						fmt.Println("TimeProc() error  update auction record err=", err.Error)
						return err.Error
					}
					return nil
				})
			}*/
			clearMergeAuction(nft, &auctionRec)
		}
		if auctionRec.Selltype == models.SellTypeBidPrice.String() {
			sp := time.Now()
			var bidRecs []models.Bidding
			err := nft.GetDB().Order("price desc").Where("auctionid = ?", auctionRec.ID).Find(&bidRecs)
			if err.Error == nil {
				if err.RowsAffected != 0 {
					nft.GetDB().Transaction(func(tx *gorm.DB) error {
						var i int
						for i = 0; i < len(bidRecs); i++ {
							if bidRecs[i].Deadtime <= time.Now().Unix() {
								fmt.Println("TimeProc() BidPrice end. useraddr=", bidRecs[i].Bidaddr)
								err = tx.Model(&models.Bidding{}).Where("id = ?", bidRecs[i].ID).Delete(&models.Bidding{})
								if err.Error != nil {
									fmt.Println("TimeProc() delete bidding record err=", err.Error)
									return err.Error
								}
							}
						}
						fmt.Println("TimeProc() len(bidRecs)=", len(bidRecs), "i=", i)
						/*if i >= len(bidRecs) {
							err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ?",
								auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Auction{})
							if err.Error != nil {
								fmt.Println("TimeProc() delete auction record err=", err.Error)
								return err.Error
							}
						}*/
						return nil
					})
				} else {
					nft.GetDB().Transaction(func(tx *gorm.DB) error {
						nftrecord := models.Nfts{}
						nftrecord.Selltype = models.SellTypeNotSale.String()
						err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nftrecord)
						if err.Error != nil {
							fmt.Println("TimeProc() update record err=", err.Error)
							return err.Error
						}
						err = tx.Model(&models.Auction{}).Where("id = ?", auctionRec.ID).Delete(&models.Auction{})
						if err.Error != nil {
							fmt.Println("TimeProc() delete auction record err=", err.Error)
							return err.Error
						}
						return nil
					})
				}
				clearMergeAuction(nft, &auctionRec)
			}
			fmt.Println(time.Now().String()[:20], "TimeProc() SellTypeBidPrice spend time=", time.Now().Sub(sp))
		}
		//sp := time.Now()
		//clearMergeAuction(nft, &auctionRec)
		//fmt.Println(time.Now().String()[:20], "TimeProc() clearMergeAuction spend time=", time.Now().Sub(sp))
	}
	fmt.Println()
	fmt.Println(time.Now().String()[:20], "TimeProc() end +++++++++++++++++++ spend time=", time.Now().Sub(spendT))
}

//func CallContractsNew(nft *models.NftDb) {
//	spendT := time.Now()
//	fmt.Println(time.Now().String()[:20], "TimeProc begin+++++++++++++++++++++++++.")
//	aucRec := models.Auction{}
//	err := nft.GetDB().Model(&models.Auction{}).Where("1=1").First(&aucRec)
//	if err.Error != nil {
//		log.Println("TimeProc() First err=", err)
//		return
//	}
//	id := aucRec.ID
//	for {
//		var auctionRec models.Auction
//		err := nft.GetDB().Model(&models.Auction{}).Where("id >= ? ", id).First(&auctionRec)
//		if err.Error != nil {
//			log.Println("TimeProc() First err=", err)
//			break
//		}
//		id = auctionRec.ID + 1
//		if auctionRec.Selltype == models.SellTypeHighestBid.String() &&
//			auctionRec.SellState == models.SellStateStart.String() &&
//			time.Now().Unix() >= auctionRec.Enddate {
//			var bidRecs []models.Bidding
//			err := nft.GetDB().Order("price desc").Where("Auctionid = ?", auctionRec.ID).Find(&bidRecs)
//			if err.Error != nil || err.RowsAffected == 0 {
//				if err.Error == gorm.ErrRecordNotFound || err.RowsAffected == 0 {
//					nft.GetDB().Transaction(func(tx *gorm.DB) error {
//						nftrecord := models.Nfts{}
//						nftrecord.Selltype = models.SellTypeNotSale.String()
//						err = tx.Model(&models.Nfts{}).Where("contract = ? AND tokenid =? AND ownaddr = ?",
//							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(&nftrecord)
//						if err.Error != nil {
//							fmt.Println("TimeProc() update record err=", err.Error)
//							return err.Error
//						}
//						err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
//							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Delete(&models.Auction{})
//						if err.Error != nil {
//							fmt.Println("TimeProc() delete auction record err=", err.Error)
//							return err.Error
//						}
//						return nil
//					})
//				}
//				continue
//			}
//			var bidRec *models.Bidding
//			for _, rec := range bidRecs {
//				valid, _, cerr := models.WormsAmountValid(rec.Price, rec.Bidaddr)
//				if cerr != nil {
//					continue
//				}
//				if !valid {
//					err = nft.GetDB().Model(&models.Bidding{}).Where("id = ?", rec.ID).Delete(&models.Bidding{})
//					if err.Error != nil {
//						fmt.Println("TimeProc() delete bidding record err=", err.Error)
//						continue
//					}
//					continue
//				}
//				bidRec = &rec
//				break
//			}
//			if bidRec == nil {
//				continue
//			}
//			//fmt.Println("TimeProc() bidRecs.Price=", bidRecs.Price, "controllers.Lowprice=",
//			//	models.Lowprice,"auctionRec.Startprice=", auctionRec.Startprice, "valid=", valid)
//			if bidRec.Price >= models.Lowprice && bidRec.Price >= auctionRec.Startprice {
//				var nftrecord models.Nfts
//				err := nft.GetDB().Where("contract = ? AND tokenid = ? AND ownaddr = ?",
//					auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).First(&nftrecord)
//				if err.Error != nil {
//					fmt.Println("TimeProc() get nftrecord err= ", err.Error)
//					continue
//				}
//				if nftrecord.Mintstate == models.Minted.String() {
//					fmt.Println("TimeProc() OwnAndAprove() == TRUE ")
//					price := strconv.FormatUint(bidRec.Price, 10) + "000000000"
//					fmt.Println("TimeProc() call  ethhelper.Auction()")
//					fmt.Println("TimeProc() auction.id=", auctionRec.ID)
//					fmt.Println("TimeProc() auctionRec.Ownaddr=", auctionRec.Ownaddr)
//					fmt.Println("TimeProc() bidRecs.Bidaddr=", bidRec.Bidaddr)
//					fmt.Println("TimeProc() auctionRec.Contract=", auctionRec.Contract)
//					fmt.Println("TimeProc() auctionRec.Tokenid=", auctionRec.Tokenid)
//					fmt.Println("TimeProc() auctionRec.Count=", auctionRec.Count)
//					fmt.Println("TimeProc() price=", price)
//					fmt.Println("TimeProc() bidRecs.Tradesig=", bidRec.Tradesig)
//
//					buyer := contracts.Buyer{}
//					err := json.Unmarshal([]byte(bidRec.Tradesig), &buyer)
//					if err != nil {
//						fmt.Println("TimeProc() ethhelper.Auction() err=", err)
//						continue
//					}
//					seller := contracts.Seller1{}
//					err = json.Unmarshal([]byte(auctionRec.Tradesig), &seller)
//					if err != nil {
//						fmt.Println("TimeProc() ethhelper.Auction() err=", err)
//						continue
//					}
//					txhash, err := contracts.AuthExchangeTrans(seller, buyer, models.ExchangerAuth, contracts.SuperAdminAddr)
//					if err != nil {
//						fmt.Println("TimeProc() contracts.AuthExchangeTrans() err=", err)
//						continue
//					}
//					fmt.Println("TimeProc() contracts.AuthExchangeTrans() Ok tx.hash=", txhash)
//					nft.GetDB().Transaction(func(tx *gorm.DB) error {
//						/*err = nft.GetDB().Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
//							auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
//						if err.Error != nil {
//							fmt.Println("TimeProc() delete bid record err=", err.Error)
//							return err.Error
//						}*/
//						auctRec := models.Auction{}
//						//auctRec.Selltype = models.SellTypeWaitSale.String()
//						auctRec.SellState = models.SellStateWait.String()
//						auctRec.Price = bidRec.Price
//						err := nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
//							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
//						if err.Error != nil {
//							fmt.Println("TimeProc() update auction record err=", err.Error)
//							return err.Error
//						}
//						return nil
//					})
//				} else {
//					price := strconv.FormatUint(bidRec.Price, 10) + "000000000"
//					fmt.Println("TimeProc() call ethhelper.AuctionAndMint()")
//					fmt.Println("TimeProc() auction.id=", auctionRec.ID)
//					fmt.Println("TimeProc() auctionRec.Ownaddr=", auctionRec.Ownaddr)
//					fmt.Println("TimeProc() bidRecs.Bidaddr=", bidRec.Bidaddr)
//					fmt.Println("TimeProc() auctionRec.Contract=", auctionRec.Contract)
//					fmt.Println("TimeProc() auctionRec.Tokenid=", auctionRec.Tokenid)
//					fmt.Println("TimeProc() price=", price)
//					fmt.Println("TimeProc() nftrecord.Meta=", nftrecord.Meta)
//					fmt.Println("TimeProc() bidRecs.Tradesig=", bidRec.Tradesig)
//					Royalty := strconv.Itoa(nftrecord.Royalty)
//					count := strconv.Itoa(nftrecord.Count)
//					fmt.Println("TimeProc() Royalty=", Royalty)
//					fmt.Println("TimeProc() count=", count)
//					fmt.Println("TimeProc() params=", auctionRec.Ownaddr, bidRec.Bidaddr, auctionRec.Contract,
//						auctionRec.Tokenid, price, count, Royalty, nftrecord.Meta, bidRec.Tradesig)
//
//					seller := contracts.Seller2{}
//					err := json.Unmarshal([]byte(auctionRec.Tradesig), &seller)
//					if err != nil {
//						fmt.Println("TimeProc() Unmarshal() err=", err)
//						continue
//					}
//					buyer := contracts.Buyer1{}
//					err = json.Unmarshal([]byte(bidRec.Tradesig), &buyer)
//					if err != nil {
//						fmt.Println("TimeProc() Unmarshal() err=", err)
//						continue
//					}
//
//					txhash, err := contracts.AuthExchangerMint(seller, buyer, models.ExchangerAuth, contracts.SuperAdminAddr)
//					if err != nil {
//						fmt.Println("TimeProc() ExchangerMint() err=", err)
//						continue
//					}
//					fmt.Println("TimeProc() contracts.AuthExchangerMint() Ok contracts.AuthExchangerMint txhash", txhash)
//					nft.GetDB().Transaction(func(tx *gorm.DB) error {
//						/*err = nft.GetDB().Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
//							auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
//						if err.Error != nil {
//							fmt.Println("TimeProc() AuctionAndMint delete bid record err=", err.Error)
//							return err.Error
//						}*/
//						auctRec := models.Auction{}
//						//auctRec.Selltype = models.SellTypeWaitSale.String()
//						auctRec.SellState = models.SellStateWait.String()
//						auctRec.Price = bidRec.Price
//						err := nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
//							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
//						if err.Error != nil {
//							fmt.Println("TimeProc() update AuctionAndMint record err=", err.Error)
//							return err.Error
//						}
//						return nil
//					})
//				}
//			}
//			/*else {
//				fmt.Println("TimeProc() auth balance error.")
//				var nftrecord models.Nfts
//				err := nft.GetDB().Where("contract = ? AND tokenid = ? AND Ownaddr = ?",
//					auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).First(&nftrecord)
//				if err.Error != nil {
//					fmt.Println("TimeProc() get nftrecord err= ", err.Error )
//					continue
//				}
//				nft.GetDB().Transaction(func(tx *gorm.DB) error {
//					trans := models.Trans{}
//					trans.Auctionid = auctionRec.ID
//					trans.Contract = auctionRec.Contract
//					trans.Fromaddr = auctionRec.Ownaddr
//					trans.Toaddr = bidRecs.Bidaddr
//					trans.Signdata = bidRecs.Signdata
//					trans.Tokenid = auctionRec.Tokenid
//					trans.Paychan = auctionRec.Paychan
//					trans.Currency = auctionRec.Currency
//					trans.Price = bidRecs.Price
//					trans.Transtime = time.Now().Unix()
//					trans.Selltype = models.SellTypeError.String()
//					trans.Error = auctionRec.Selltype + "," + errmsg
//					trans.Name = nftrecord.Name
//					trans.Meta = nftrecord.Meta
//					trans.Desc = nftrecord.Desc
//					err := tx.Model(&trans).Create(&trans)
//					if err.Error != nil {
//						fmt.Println("TimeProc() error create trans record err=", err.Error)
//						return err.Error
//					}
//					//auctRec := models.Auction{}
//					////auctRec.Selltype = models.SellTypeWaitSale.String()
//					//auctRec.SellState = models.SellStateWait.String()
//					//auctRec.Price = bidRecs.Price
//					//err = nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ?",
//					//	auctionRec.Contract, auctionRec.Tokenid).Updates(auctRec)
//					//if err.Error != nil {
//					//	fmt.Println("TimeProc() error  update auction record err=", err.Error)
//					//	return err.Error
//					//}
//					nftrecord := models.Nfts{}
//					nftrecord.Selltype = models.SellTypeNotSale.String()
//					err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nftrecord)
//					if err.Error != nil {
//						fmt.Println("TimeProc() update record err=", err.Error)
//						return err.Error
//					}
//					err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
//						auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Delete(&models.Auction{})
//					if err.Error != nil {
//						fmt.Println("BuyResult() delete auction record err=", err.Error)
//						return err.Error
//					}
//					err = tx.Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
//						auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
//					if err.Error != nil {
//						fmt.Println("TimeProc() error  update auction record err=", err.Error)
//						return err.Error
//					}
//					return nil
//				})
//			}*/
//			clearMergeAuction(nft, &auctionRec)
//		}
//		if auctionRec.Selltype == models.SellTypeBidPrice.String() {
//			sp := time.Now()
//			var nftrecord models.Nfts
//			err := nft.GetDB().Where("contract = ? AND tokenid =?", auctionRec.Contract, auctionRec.Tokenid).First(&nftrecord)
//			if err.Error != nil {
//				log.Println("TimeProc()() not find nft err= ", err.Error)
//				return
//			}
//			var bidRecs []models.Bidding
//			err = nft.GetDB().Order("price desc").Where("auctionid = ?", auctionRec.ID).Find(&bidRecs)
//			if err.Error == nil {
//				if err.RowsAffected != 0 {
//					nft.GetDB().Transaction(func(tx *gorm.DB) error {
//						var i int
//						var lastBids []models.Bidding
//						for _, bid := range bidRecs {
//							if bid.Deadtime <= time.Now().Unix() {
//								fmt.Println("TimeProc() BidPrice end. useraddr=", bid.Bidaddr)
//								err = tx.Model(&models.Bidding{}).Where("id = ?", bid.ID).Delete(&models.Bidding{})
//								if err.Error != nil {
//									fmt.Println("TimeProc() delete bidding record err=", err.Error)
//									return err.Error
//								}
//							} else {
//								lastBids = append(lastBids, bid)
//							}
//						}
//						fmt.Println("TimeProc() len(bidRecs)=", len(bidRecs), "i=", i)
//						nfttab := map[string]interface{}{
//							"Offernum":    0,
//							"Maxbidprice": 0,
//						}
//						if len(lastBids) > 0 {
//							nfttab["Offernum"] = len(lastBids)
//							nfttab["Maxbidprice"] = lastBids[0].Price
//						} else {
//							nfttab = map[string]interface{}{
//								"Offernum":    0,
//								"Maxbidprice": 0,
//								"Selltype":    models.SellTypeNotSale.String(),
//							}
//							err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ?",
//								auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Auction{})
//							if err.Error != nil {
//								fmt.Println("TimeProc() delete auction record err=", err.Error)
//								return err.Error
//							}
//						}
//						err = tx.Model(&models.Nfts{}).Where("contract = ? AND tokenid =?",
//							auctionRec.Contract, auctionRec.Tokenid).Updates(&nfttab)
//						if err.Error != nil {
//							log.Println("TimeProc() update record err=", err.Error)
//							return err.Error
//						}
//						return nil
//					})
//				} else {
//					nft.GetDB().Transaction(func(tx *gorm.DB) error {
//						nfttab := map[string]interface{}{
//							"Offernum":    0,
//							"Maxbidprice": 0,
//							"Selltype":    models.SellTypeNotSale.String(),
//						}
//						err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nfttab)
//						if err.Error != nil {
//							fmt.Println("TimeProc() update record err=", err.Error)
//							return err.Error
//						}
//						err = tx.Model(&models.Auction{}).Where("id = ?", auctionRec.ID).Delete(&models.Auction{})
//						if err.Error != nil {
//							fmt.Println("TimeProc() delete auction record err=", err.Error)
//							return err.Error
//						}
//						return nil
//					})
//				}
//				clearMergeAuction(nft, &auctionRec)
//			}
//			fmt.Println(time.Now().String()[:20], "TimeProc() SellTypeBidPrice spend time=", time.Now().Sub(sp))
//		}
//		//sp := time.Now()
//		if auctionRec.Selltype == models.SellTypeFixPrice.String() {
//			clearMergeAuction(nft, &auctionRec)
//		}
//		//fmt.Println(time.Now().String()[:20], "TimeProc() clearMergeAuction spend time=", time.Now().Sub(sp))
//	}
//	fmt.Println()
//	fmt.Println(time.Now().String()[:20], "TimeProc() end +++++++++++++++++++ spend time=", time.Now().Sub(spendT))
//}

func CallContractsNew(nft *models.NftDb) {
	spendT := time.Now()
	fmt.Println(time.Now().String()[:20], "TimeProc begin+++++++++++++++++++++++++.")
	//aucRec := models.Auction{}
	//err := nft.GetDB().Model(&models.Auction{}).Where("1=1").First(&aucRec)
	//if err.Error != nil {
	//	log.Println("TimeProc() First err=", err)
	//	return
	//}
	//id := aucRec.ID
	var aucRecs []models.Auction
	err := nft.GetDB().Model(&models.Auction{}).Where("Enddate <= ? and Selltype = ? and Sell_State = ?",
		time.Now().Unix(), models.SellTypeHighestBid.String(), models.SellStateStart.String()).Limit(100).Find(&aucRecs)
	if err.Error != nil {
		log.Println("TimeProc() First err=", err)
		return
	}
	fmt.Println("TimeProc() HighestBid record=", len(aucRecs))
	if len(aucRecs) > 0 {
		for _, auctionRec := range aucRecs {
			if auctionRec.Selltype == models.SellTypeHighestBid.String() &&
				auctionRec.SellState == models.SellStateStart.String() &&
				time.Now().Unix() >= auctionRec.Enddate {
				var bidRecs []models.Bidding
				err := nft.GetDB().Order("price desc").Where("Auctionid = ?", auctionRec.ID).Find(&bidRecs)
				if err.Error != nil || err.RowsAffected == 0 {
					if err.Error == gorm.ErrRecordNotFound || err.RowsAffected == 0 {
						nft.GetDB().Transaction(func(tx *gorm.DB) error {
							nftrecord := models.Nfts{}
							nftrecord.Selltype = models.SellTypeNotSale.String()
							err = tx.Model(&models.Nfts{}).Where("contract = ? AND tokenid =? AND ownaddr = ?",
								auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(&nftrecord)
							if err.Error != nil {
								fmt.Println("TimeProc() update record err=", err.Error)
								return err.Error
							}
							err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
								auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Delete(&models.Auction{})
							if err.Error != nil {
								fmt.Println("TimeProc() delete auction record err=", err.Error)
								return err.Error
							}
							return nil
						})
					}
					continue
				}
				var bidRec *models.Bidding
				for _, rec := range bidRecs {
					valid, _, cerr := models.WormsAmountValid(rec.Price, rec.Bidaddr)
					if cerr != nil {
						continue
					}
					if !valid {
						err = nft.GetDB().Model(&models.Bidding{}).Where("id = ?", rec.ID).Delete(&models.Bidding{})
						if err.Error != nil {
							fmt.Println("TimeProc() delete bidding record err=", err.Error)
							continue
						}
						continue
					}
					bidRec = &rec
					break
				}
				if bidRec == nil {
					continue
				}
				//fmt.Println("TimeProc() bidRecs.Price=", bidRecs.Price, "controllers.Lowprice=",
				//	models.Lowprice,"auctionRec.Startprice=", auctionRec.Startprice, "valid=", valid)
				if bidRec.Price >= models.Lowprice && bidRec.Price >= auctionRec.Startprice {
					var nftrecord models.Nfts
					err := nft.GetDB().Where("contract = ? AND tokenid = ? AND ownaddr = ?",
						auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).First(&nftrecord)
					if err.Error != nil {
						fmt.Println("TimeProc() get nftrecord err= ", err.Error)
						continue
					}
					if nftrecord.Mintstate == models.Minted.String() {
						fmt.Println("TimeProc() OwnAndAprove() == TRUE ")
						price := strconv.FormatUint(bidRec.Price, 10) + "000000000"
						fmt.Println("TimeProc() call  ethhelper.Auction()")
						fmt.Println("TimeProc() auction.id=", auctionRec.ID)
						fmt.Println("TimeProc() auctionRec.Ownaddr=", auctionRec.Ownaddr)
						fmt.Println("TimeProc() bidRecs.Bidaddr=", bidRec.Bidaddr)
						fmt.Println("TimeProc() auctionRec.Contract=", auctionRec.Contract)
						fmt.Println("TimeProc() auctionRec.Tokenid=", auctionRec.Tokenid)
						fmt.Println("TimeProc() auctionRec.Count=", auctionRec.Count)
						fmt.Println("TimeProc() price=", price)
						fmt.Println("TimeProc() bidRecs.Tradesig=", bidRec.Tradesig)

						buyer := contracts.Buyer{}
						err := json.Unmarshal([]byte(bidRec.Tradesig), &buyer)
						if err != nil {
							fmt.Println("TimeProc() ethhelper.Auction() err=", err)
							continue
						}
						seller := contracts.Seller1{}
						err = json.Unmarshal([]byte(auctionRec.Tradesig), &seller)
						if err != nil {
							fmt.Println("TimeProc() ethhelper.Auction() err=", err)
							continue
						}
						txhash, err := contracts.AuthExchangeTrans(seller, buyer, models.ExchangerAuth, contracts.SuperAdminAddr)
						if err != nil {
							fmt.Println("TimeProc() contracts.AuthExchangeTrans() err=", err)
							continue
						}
						fmt.Println("TimeProc() contracts.AuthExchangeTrans() Ok tx.hash=", txhash)
						nft.GetDB().Transaction(func(tx *gorm.DB) error {
							/*err = nft.GetDB().Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
								auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
							if err.Error != nil {
								fmt.Println("TimeProc() delete bid record err=", err.Error)
								return err.Error
							}*/
							auctRec := models.Auction{}
							//auctRec.Selltype = models.SellTypeWaitSale.String()
							auctRec.SellState = models.SellStateWait.String()
							auctRec.Price = bidRec.Price
							err := nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
								auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
							if err.Error != nil {
								fmt.Println("TimeProc() update auction record err=", err.Error)
								return err.Error
							}
							return nil
						})
					} else {
						price := strconv.FormatUint(bidRec.Price, 10) + "000000000"
						fmt.Println("TimeProc() call ethhelper.AuctionAndMint()")
						fmt.Println("TimeProc() auction.id=", auctionRec.ID)
						fmt.Println("TimeProc() auctionRec.Ownaddr=", auctionRec.Ownaddr)
						fmt.Println("TimeProc() bidRecs.Bidaddr=", bidRec.Bidaddr)
						fmt.Println("TimeProc() auctionRec.Contract=", auctionRec.Contract)
						fmt.Println("TimeProc() auctionRec.Tokenid=", auctionRec.Tokenid)
						fmt.Println("TimeProc() price=", price)
						fmt.Println("TimeProc() nftrecord.Meta=", nftrecord.Meta)
						fmt.Println("TimeProc() bidRecs.Tradesig=", bidRec.Tradesig)
						Royalty := strconv.Itoa(nftrecord.Royalty)
						count := strconv.Itoa(nftrecord.Count)
						fmt.Println("TimeProc() Royalty=", Royalty)
						fmt.Println("TimeProc() count=", count)
						fmt.Println("TimeProc() params=", auctionRec.Ownaddr, bidRec.Bidaddr, auctionRec.Contract,
							auctionRec.Tokenid, price, count, Royalty, nftrecord.Meta, bidRec.Tradesig)

						seller := contracts.Seller2{}
						err := json.Unmarshal([]byte(auctionRec.Tradesig), &seller)
						if err != nil {
							fmt.Println("TimeProc() Unmarshal() err=", err)
							continue
						}
						buyer := contracts.Buyer1{}
						err = json.Unmarshal([]byte(bidRec.Tradesig), &buyer)
						if err != nil {
							fmt.Println("TimeProc() Unmarshal() err=", err)
							continue
						}

						txhash, err := contracts.AuthExchangerMint(seller, buyer, models.ExchangerAuth, contracts.SuperAdminAddr)
						if err != nil {
							fmt.Println("TimeProc() ExchangerMint() err=", err)
							continue
						}
						fmt.Println("TimeProc() contracts.AuthExchangerMint() Ok contracts.AuthExchangerMint txhash", txhash)
						nft.GetDB().Transaction(func(tx *gorm.DB) error {
							/*err = nft.GetDB().Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
								auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
							if err.Error != nil {
								fmt.Println("TimeProc() AuctionAndMint delete bid record err=", err.Error)
								return err.Error
							}*/
							auctRec := models.Auction{}
							//auctRec.Selltype = models.SellTypeWaitSale.String()
							auctRec.SellState = models.SellStateWait.String()
							auctRec.Price = bidRec.Price
							err := nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
								auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Updates(auctRec)
							if err.Error != nil {
								fmt.Println("TimeProc() update AuctionAndMint record err=", err.Error)
								return err.Error
							}
							return nil
						})
					}
				}
				/*else {
					fmt.Println("TimeProc() auth balance error.")
					var nftrecord models.Nfts
					err := nft.GetDB().Where("contract = ? AND tokenid = ? AND Ownaddr = ?",
						auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).First(&nftrecord)
					if err.Error != nil {
						fmt.Println("TimeProc() get nftrecord err= ", err.Error )
						continue
					}
					nft.GetDB().Transaction(func(tx *gorm.DB) error {
						trans := models.Trans{}
						trans.Auctionid = auctionRec.ID
						trans.Contract = auctionRec.Contract
						trans.Fromaddr = auctionRec.Ownaddr
						trans.Toaddr = bidRecs.Bidaddr
						trans.Signdata = bidRecs.Signdata
						trans.Tokenid = auctionRec.Tokenid
						trans.Paychan = auctionRec.Paychan
						trans.Currency = auctionRec.Currency
						trans.Price = bidRecs.Price
						trans.Transtime = time.Now().Unix()
						trans.Selltype = models.SellTypeError.String()
						trans.Error = auctionRec.Selltype + "," + errmsg
						trans.Name = nftrecord.Name
						trans.Meta = nftrecord.Meta
						trans.Desc = nftrecord.Desc
						err := tx.Model(&trans).Create(&trans)
						if err.Error != nil {
							fmt.Println("TimeProc() error create trans record err=", err.Error)
							return err.Error
						}
						//auctRec := models.Auction{}
						////auctRec.Selltype = models.SellTypeWaitSale.String()
						//auctRec.SellState = models.SellStateWait.String()
						//auctRec.Price = bidRecs.Price
						//err = nft.GetDB().Model(&models.Auction{}).Where("contract = ? AND tokenid = ?",
						//	auctionRec.Contract, auctionRec.Tokenid).Updates(auctRec)
						//if err.Error != nil {
						//	fmt.Println("TimeProc() error  update auction record err=", err.Error)
						//	return err.Error
						//}
						nftrecord := models.Nfts{}
						nftrecord.Selltype = models.SellTypeNotSale.String()
						err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nftrecord)
						if err.Error != nil {
							fmt.Println("TimeProc() update record err=", err.Error)
							return err.Error
						}
						err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ? AND ownaddr = ?",
							auctionRec.Contract, auctionRec.Tokenid, auctionRec.Ownaddr).Delete(&models.Auction{})
						if err.Error != nil {
							fmt.Println("BuyResult() delete auction record err=", err.Error)
							return err.Error
						}
						err = tx.Model(&models.Bidding{}).Where("contract = ? AND tokenid = ?",
							auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Bidding{})
						if err.Error != nil {
							fmt.Println("TimeProc() error  update auction record err=", err.Error)
							return err.Error
						}
						return nil
					})
				}*/
				clearMergeAuction(nft, &auctionRec)
			}
		}
	}
	var bidRecs []models.Bidding
	err = nft.GetDB().Model(&models.Bidding{}).Where("Deadtime <= ?", time.Now().Unix()).Limit(100).Find(&bidRecs)
	if err.Error != nil {
		log.Println("TimeProc() Bidding err=", err)
		return
	}
	fmt.Println("TimeProc() BidPrice record=", len(bidRecs))
	if len(bidRecs) > 0 {
		for _, bidRec := range bidRecs {
			auctionRec := models.Auction{}
			err := nft.GetDB().Model(&models.Auction{}).Where("id = ? ", bidRec.Auctionid).First(&auctionRec)
			if err.Error != nil {
				if err.Error == gorm.ErrRecordNotFound {
					continue
				}
				log.Println("TimeProc() Auction First err=", err)
				return
			}
			if auctionRec.Selltype == models.SellTypeBidPrice.String() {
				sp := time.Now()
				var nftrecord models.Nfts
				err := nft.GetDB().Where("contract = ? AND tokenid =?", auctionRec.Contract, auctionRec.Tokenid).First(&nftrecord)
				if err.Error != nil {
					log.Println("TimeProc()() not find nft err= ", err.Error)
					return
				}
				var bidRecs []models.Bidding
				err = nft.GetDB().Order("price desc").Where("auctionid = ?", auctionRec.ID).Find(&bidRecs)
				if err.Error == nil {
					if err.RowsAffected != 0 {
						nft.GetDB().Transaction(func(tx *gorm.DB) error {
							var i int
							var lastBids []models.Bidding
							for _, bid := range bidRecs {
								if bid.Deadtime <= time.Now().Unix() {
									fmt.Println("TimeProc() BidPrice end. useraddr=", bid.Bidaddr)
									err = tx.Model(&models.Bidding{}).Where("id = ?", bid.ID).Delete(&models.Bidding{})
									if err.Error != nil {
										fmt.Println("TimeProc() delete bidding record err=", err.Error)
										return err.Error
									}
								} else {
									lastBids = append(lastBids, bid)
								}
							}
							fmt.Println("TimeProc() len(bidRecs)=", len(bidRecs), "i=", i)
							nfttab := map[string]interface{}{
								"Offernum":    0,
								"Maxbidprice": 0,
							}
							if len(lastBids) > 0 {
								nfttab["Offernum"] = len(lastBids)
								nfttab["Maxbidprice"] = lastBids[0].Price
							} else {
								nfttab = map[string]interface{}{
									"Offernum":    0,
									"Maxbidprice": 0,
									"Selltype":    models.SellTypeNotSale.String(),
								}
								err = tx.Model(&models.Auction{}).Where("contract = ? AND tokenid = ?",
									auctionRec.Contract, auctionRec.Tokenid).Delete(&models.Auction{})
								if err.Error != nil {
									fmt.Println("TimeProc() delete auction record err=", err.Error)
									return err.Error
								}
							}
							err = tx.Model(&models.Nfts{}).Where("contract = ? AND tokenid =?",
								auctionRec.Contract, auctionRec.Tokenid).Updates(&nfttab)
							if err.Error != nil {
								log.Println("TimeProc() update record err=", err.Error)
								return err.Error
							}
							return nil
						})
					} else {
						nft.GetDB().Transaction(func(tx *gorm.DB) error {
							nfttab := map[string]interface{}{
								"Offernum":    0,
								"Maxbidprice": 0,
								"Selltype":    models.SellTypeNotSale.String(),
							}
							err = tx.Model(&models.Nfts{}).Where("id = ?", auctionRec.Nftid).Updates(&nfttab)
							if err.Error != nil {
								fmt.Println("TimeProc() update record err=", err.Error)
								return err.Error
							}
							err = tx.Model(&models.Auction{}).Where("id = ?", auctionRec.ID).Delete(&models.Auction{})
							if err.Error != nil {
								fmt.Println("TimeProc() delete auction record err=", err.Error)
								return err.Error
							}
							return nil
						})
					}
					clearMergeAuction(nft, &auctionRec)
				}
				fmt.Println(time.Now().String()[:20], "TimeProc() SellTypeBidPrice spend time=", time.Now().Sub(sp))
			}
		}
	}
	fmt.Println()
	fmt.Println(time.Now().String()[:20], "TimeProc() end +++++++++++++++++++ spend time=", time.Now().Sub(spendT))
}
