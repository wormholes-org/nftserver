package main

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	_ "github.com/nftexchange/nftserver/routers"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func TimeProc(sqldsn string) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			nd, err := models.NewNftDb(sqldsn)
			if err != nil {
				fmt.Printf("connect database err = %s\n", err)
				continue
			}
			CallContracts(nd)
			nd.Close()
		}
	}
}

func CallContracts(nft *models.NftDb) {
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
		}
		if auctionRec.Selltype == models.SellTypeBidPrice.String() {
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
			}
		}
	}
	fmt.Println()
	fmt.Println(time.Now().String()[:20], "TimeProc() end +++++++++++++++++++")
}
