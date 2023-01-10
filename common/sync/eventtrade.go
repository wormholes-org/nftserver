package sync

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	"log"
	"time"
)

const (
	ReDialDelyTime = 5
)

var EventCh chan interface{}

//Royalty(uint256 id, uint256 royalty, address receiver);
//func EventRoyalty(sqldsn string) {
//	var client *ethclient.Client
//	var err error
//	for {
//		for  {
//			client, err = ethclient.Dial(contracts.EthWsNode)
//			if err != nil {
//				log.Println("EventRoyalty() connect err=", err)
//				fmt.Println("EventRoyalty() connect err=", err)
//				time.Sleep(ReDialDelyTime * time.Second)
//			} else {
//				log.Println("EventRoyalty() connect OK!")
//				fmt.Println("EventRoyalty() connect OK!")
//				break
//			}
//		}
//		RoyaltyCh := make(chan *nft1155.Nft1155Royalty)
//
//		address := common.HexToAddress(models.NFT1155Addr)
//		instance, err := nft1155.NewNft1155(address, client)
//		if err != nil {
//			log.Println("EventRoyalty() new NewNft1155 err=", err)
//			fmt.Println("EventRoyalty() new NewNft1155 err=", err)
//			client.Close()
//			continue
//		}
//		sub, err := instance.WatchRoyalty(&bind.WatchOpts{Start: nil, Context: context.Background()}, RoyaltyCh)
//		if err != nil {
//			log.Println("EventRoyalty() WatchRoyalty error=", err)
//			fmt.Println("EventRoyalty() WatchRoyalty error=", err)
//			client.Close()
//			continue
//		}
//		//var emptyaddr []common.Address
//		//TransferSingleCh := make(chan *nft1155.Nft1155TransferSingle)
//		//_, err = instance.WatchTransferSingle(&bind.WatchOpts{Start: nil, Context: context.Background()}, TransferSingleCh, emptyaddr, emptyaddr, emptyaddr)
//		//if err != nil {
//		//	log.Println("EventRoyalty() WatchTransferSingle error=", err)
//		//	fmt.Println("EventRoyalty() WatchTransferSingle error=", err)
//		//	client.Close()
//		//	continue
//		//}
//		te := make(chan struct{})
//		go func() {
//			ticker := time.NewTicker(waitTime)
//			for {
//				select {
//				case <-ticker.C:
//					header, err := client.HeaderByNumber(context.Background(), nil)
//					if err != nil {
//						log.Println("EventRoyalty() get HeaderByNumber err=", err)
//						fmt.Println("EventRoyalty() get HeaderByNumber err=", err)
//						continue
//					}
//					block, err := client.BlockByNumber(context.Background(), header.Number)
//					if err != nil {
//						log.Println("EventRoyalty() get HeaderByNumber err=", err)
//						fmt.Println("EventRoyalty() get HeaderByNumber err=", err)
//						continue
//					}
//					fmt.Println("EventRoyalty() block.Number()=", block.Number())
//				case <-te:
//					te <- struct{}{}
//					fmt.Println("EventRoyalty() ticker end.")
//					return
//				}
//			}
//		}()
//		fmt.Println("EventRoyalty start!")
//	loop:
//		for {
//			select {
//			case look := <- RoyaltyCh:
//				EventCh <- look
//			case err := <-sub.Err():
//				fmt.Println("EventRoyalty() sub.err()=", err)
//				sub.Unsubscribe()
//				te <- struct{}{}
//				<- te
//				client.Close()
//				log.Println("EventRoyalty() restart.")
//				break loop
//			}
//		}
//	}
//}

//event SALE(address from, address to, address nftAddr, uint256 nftId, uint256 price, bytes sig);
//func EventSale(sqldsn string) {
//	var client *ethclient.Client
//	var err error
//	for {
//		for {
//			client, err = ethclient.Dial(contracts.EthWsNode)
//			if err != nil {
//				log.Println("EventSale() connect err=", err)
//				fmt.Println("EventSale() connect err=", err)
//				time.Sleep(ReDialDelyTime * time.Second)
//			} else {
//				log.Println("EventSale() connect OK!")
//				fmt.Println("EventSale() connect OK!")
//				break
//			}
//		}
//
//		ch := make(chan *trade.TradePRICING)
//		address := common.HexToAddress(models.TradeAddr)
//		instance, err := trade.NewTrade(address, client)
//		if err != nil {
//			log.Println("EventSale() NewTrade err=", err)
//			fmt.Println("EventSale() NewTrade err=", err)
//			client.Close()
//			continue
//		}
//		sub, err := instance.WatchPRICING(&bind.WatchOpts{Start: nil, Context: context.Background()}, ch)
//		if err != nil {
//			log.Println("EventSale() WatchPRICING error", err)
//			fmt.Println("EventSale() WatchPRICING error", err)
//			client.Close()
//			continue
//		}
//		te := make(chan struct{})
//		go func() {
//			ticker := time.NewTicker(waitTime)
//			for {
//				select {
//				case <-ticker.C:
//					header, err := client.HeaderByNumber(context.Background(), nil)
//					if err != nil {
//						log.Println("EventSale() get header err=", err)
//						continue
//					}
//					block, err := client.BlockByNumber(context.Background(), header.Number)
//					if err != nil {
//						log.Println("EventSale() get BlockByNumber err=", err)
//						continue
//					}
//					fmt.Println("EventSale() block.Number()=", block.Number())
//				case <-te:
//					te <- struct{}{}
//					log.Println("EventSale() ticker end.")
//					return
//				}
//			}
//		}()
//		fmt.Println("EventSale() start.")
//	loop:
//		for {
//			select {
//			case look := <-ch:
//				EventCh <- look
//			case err := <-sub.Err():
//				fmt.Println("EventSale() sub.err()=", err)
//				sub.Unsubscribe()
//				te <- struct{}{}
//				<-te
//				fmt.Println("EventSale() restart.")
//				client.Close()
//				break loop
//			}
//		}
//	}
//}

//func EventAuction(sqldsn string) {
//	var client *ethclient.Client
//	var err error
//	for {
//		for  {
//			client, err = ethclient.Dial(contracts.EthWsNode)
//			if err != nil {
//				log.Println("EventAuction() connect err=", err)
//				fmt.Println("EventAuction() connect err=", err)
//				time.Sleep(ReDialDelyTime * time.Second)
//			} else {
//				log.Println("EventAuction() connect OK!")
//				fmt.Println("EventAuction() connect OK!")
//				break
//			}
//		}
//		ch := make(chan *trade.TradeBIDING)
//		address := common.HexToAddress(models.TradeAddr)
//		instance, err := trade.NewTrade(address, client)
//		if err != nil {
//			fmt.Println("EventAuction() NewTrade err=", err)
//			log.Println("EventAuction() NewTrade err=", err)
//			client.Close()
//			continue
//		}
//		sub, err := instance.WatchBIDING(&bind.WatchOpts{Start: nil, Context: context.Background()}, ch)
//		if err != nil {
//			log.Println("EventAuction() WatchBIDING err=", err)
//			fmt.Println("EventAuction() WatchBIDING err=", err)
//			client.Close()
//			continue
//		}
//		te := make(chan struct{})
//		go func() {
//			ticker := time.NewTicker(waitTime)
//			for {
//				select {
//				case <-ticker.C:
//					header, err := client.HeaderByNumber(context.Background(), nil)
//					if err != nil {
//						log.Println("EventAuction() get header err=", err)
//						fmt.Println("EventAuction() get header err=", err)
//						continue
//					}
//					block, err := client.BlockByNumber(context.Background(), header.Number)
//					if err != nil {
//						log.Println("EventAuction() get BlockByNumber err=", err)
//						fmt.Println("EventAuction() get BlockByNumber err=", err)
//						continue
//					}
//					fmt.Println("EventAuction() block.Number()=", block.Number())
//				case <-te:
//					te <- struct{}{}
//					log.Println("EventAuction() ticker end.")
//					return
//				}
//			}
//		}()
//		fmt.Println("EventAuction() start!")
//	loop:
//		for {
//			select {
//			case look := <-ch:
//				EventCh <- look
//			case err := <-sub.Err():
//				fmt.Println("EventAuction sub.err()=", err)
//				sub.Unsubscribe()
//				te <- struct{}{}
//				<- te
//				log.Println("EventAuction() restart.")
//				client.Close()
//				break loop
//			}
//		}
//	}
//}

/*func EventQueue(sqldsn string) {
	var syncCh chan uint64
	var procEnd chan struct{}
	for {
		select {
		case event := <-EventCh:
			fmt.Println("get event ", "type=", reflect.TypeOf(event))
			switch value := event.(type) {
				case *nft1155.Nft1155Royalty:
					//fmt.Println("EventRoyalty()<-RoyaltyCh Ok", "look.Id=",value.Id.String(), "look.Royalty=", value.Royalty.String())
					nft, err := models.NewNftDb(sqldsn)
					if err == nil {
						//BuyResult(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty string)
						//err = nft.BuyResult("", value.Receiver.String(), models.NFT1155Addr,
						//	value.Id.String(), "", "", "", fmt.Sprintf("%v", value.Ratio),value.Raw.TxHash.String())
						err = nft.BuyResultRoyalty("", value.Receiver.String(), models.NFT1155Addr,
							value.Id.String(), "", fmt.Sprintf("%v", value.Ratio),value.Raw.TxHash.String(), strconv.FormatInt(time.Now().Unix(), 10))
						fmt.Println("EventQueue()-->EventRoyalty() <-RoyaltyCh BuyResult() err=", err)
						fmt.Println(err)
					}
					nft.Close()
				case *trade.TradePRICING:
					fmt.Println("EventQueue()-->TradePRICING <-ch",
						"look.FROM=", value.From.String(),
						"look.to=", value.To.String(),
						"look.ContractAddr=", value.Addr.String(),
						"look.tokenId=", value.Id.String(),
						"look.Amount=", value.Amount.String(),
						"look.Price=", value.Price.String(),
						"look.Raw.TxHash=", value.Raw.TxHash.String())
					nft, err := models.NewNftDb(sqldsn)
					if err == nil {
						//BuyResultNew(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty string)
						price := value.Price.String()
						fmt.Println("EventQueue()-->EventSale() price=", price)
						price = price[:len(price)-9]
						fmt.Println("EventQueue()-->EventSale() price=", price)
						//err := nft.BuyResult(value.From.Hex(), value.To.Hex(), value.Addr.Hex(),
						//	value.Id.String(), "", price, "", "",value.Raw.TxHash.String())
						count := value.Amount.String()
						err := nft.BuyResultWithAmount(value.From.Hex(), value.To.Hex(), value.Addr.Hex(),
							value.Id.String(), count, price, "",  value.Raw.TxHash.String(), strconv.FormatInt(time.Now().Unix(), 10))
						fmt.Println("EventQueue()-->EventSale() BuyResult() err=", err)
					}
					//time.Sleep(2*time.Second)
					nft.Close()
				case *nft1155.Nft1155TransferSingle:
					fmt.Println("EventQueue()--> Rev <-TransferSingleCh <-ch",
						"look.FROM=", value.From.String(),
						"look.to=", value.To.String(),
						"look.tokenId=", value.Id.String(),
						"look.value=", value.Value.String(),
						"look.Raw.TxHash=", value.Raw.TxHash.String())
					nft, err := models.NewNftDb(sqldsn)
					if err == nil {
						//BuyResult(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty string)
						//err = nft.BuyResult(value.From.String(), value.To.String(), models.NFT1155Addr,
						//	value.Id.String(), "", "", "", "",value.Raw.TxHash.String())
						err = nft.BuyResultRoyalty(value.From.String(), value.To.String(), models.NFT1155Addr,
							value.Id.String(), "", "",value.Raw.TxHash.String(), strconv.FormatInt(time.Now().Unix(), 10))
						fmt.Println("EventQueue()-->EventRoyalty() <-TransferSingleCh nft.BuyResult() err=", err)
					}
					nft.Close()
				case *trade.TradeBIDING:
					fmt.Println("EventQueue()--> rev TradeBIDING <-ch",
						"look.FROM=", value.From.String(),
						"look.to=", value.To.String(),
						"look.ContractAddr=", value.Addr.String(),
						"look.tokenId=", value.Id.String(),
						"look.Amount=", value.Amount.String(),
						"look.Price=", value.Price.String(),
						"look.Raw.TxHash=", value.Raw.TxHash.String())
					price := value.Price.String()
					price = price[:len(price)-9]
					nft := new(models.NftDb)
					nft, err := models.NewNftDb(sqldsn)
					if err == nil {
						//BuyResultNew(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty string)
						fmt.Println("EventQueue()-->EventAuction() look.Raw.TxHash=", value.Raw.TxHash.String())
						//err := nft.BuyResult(value.From.Hex(), value.To.Hex(), value.Addr.Hex(),
						//	value.Id.String(), "", price, "", "",value.Raw.TxHash.String())
						count := value.Amount.String()
						err := nft.BuyResultWithAmount(value.From.Hex(), value.To.Hex(), value.Addr.Hex(),
							value.Id.String(), count, price, "", value.Raw.TxHash.String(), strconv.FormatInt(time.Now().Unix(), 10))
						fmt.Println("EventQueue()-->EventAuction() BuyResult() err=", err)
					}
					nft.Close()
				case *types.Header:
					fmt.Println(time.Now().String()[:22],"EventQueue()-->SyncProc() EventBlockTxs() block.number=",  value.Number.Uint64())
					if syncCh == nil {
						syncCh = make(chan uint64, 200)
						procEnd = models.SyncProc(sqldsn, syncCh)
						syncCh <- value.Number.Uint64()
					} else {
						select {
							case <-procEnd:
								fmt.Println("EventQueue()-->SyncProc() <-procEnd:")
								syncCh <- value.Number.Uint64()
							default:
								fmt.Println("EventQueue()-->SyncProc() default:")
						}
					}
					fmt.Println(time.Now().String()[:22],"EventQueue()-->SyncProc() OK")
			}
		}
	}
}*/

func EventContract(sqldsn string) {
	EventCh = make(chan interface{}, 100)
	//go EventQueue(sqldsn)
	//go EventSale(sqldsn)
	//go EventRoyalty(sqldsn)
	//go EventAuction(sqldsn)
	//go EventBlockTxs(sqldsn)
	select {}
}

func EventBlockTxs(sqldsn string) {
	var client *ethclient.Client
	var err error
	for {
		for {
			client, err = ethclient.Dial(contracts.EthWsNode)
			if err != nil {
				log.Println("EventBlockTxs() connect err=", err)
				fmt.Println("EventBlockTxs() connect err=", err)
				time.Sleep(ReDialDelyTime * time.Second)
			} else {
				log.Println("EventBlockTxs() connect OK!")
				fmt.Println("EventBlockTxs() connect OK!")
				break
			}
		}
		headers := make(chan *types.Header)
		sub, err := client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			log.Println("EventBlockTxs() SubscribeNewHead err=", err)
			fmt.Println("EventBlockTxs() SubscribeNewHead err=", err)
			client.Close()
			continue
		}
		nd, err := models.NewNftDb(sqldsn)
		if err != nil {
			fmt.Printf("EventBlockTxs() connect database err = %s\n", err)
			client.Close()
			continue
		}
	loop:
		for {
			select {
			case err := <-sub.Err():
				fmt.Println("EventBlockTxs() sub.err()=", err)
				sub.Unsubscribe()
				client.Close()
				log.Println("EventBlockTxs() restart.")
				break loop
			case header := <-headers:
				headNumber := header.Number.String()
				fmt.Println("EventBlockTxs() headNumber=", headNumber)
				EventCh <- header
			}
		}
		nd.Close()
		fmt.Println("UpdateBlockNumber() db close.")
	}
}
