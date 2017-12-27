package wallet

import (
	"log"
	"sync"

	"github.com/sazor/bittrex-wallet/services/client"
)

const buyOrder = "LIMIT_BUY"

type WalletCoin struct {
	Ticker   string
	Balance  float64
	AvgPrice float64
	Last     float64
	Bid      float64
	Ask      float64
}

func NewWalletCoin(ticker string, balance float64) *WalletCoin {
	return &WalletCoin{Ticker: ticker, Balance: balance}
}

func (coin *WalletCoin) DiffPrice() float64 {
	return coin.Last - coin.AvgPrice
}

func (coin *WalletCoin) PercentDiffPrice() float64 {
	return (coin.DiffPrice() / coin.AvgPrice) * 100
}

func (coin *WalletCoin) BtcBalance() float64 {
	return coin.Last * coin.Balance
}

var once sync.Once

func (coin *WalletCoin) CalcAvgPrice() float64 {
	once.Do(func() {
		coin.RefreshAvgPrice()
	})
	return coin.AvgPrice
}

func (coin *WalletCoin) RefreshAvgPrice() {
	clnt, err := client.GetClient()
	if err != nil {
		log.Printf("Cant get history for %s", coin.Ticker)
		return
	}
	orders, err := clnt.GetOrderHistory("BTC-" + coin.Ticker)
	if err != nil {
		log.Printf("Cant get history for %s", coin.Ticker)
		return
	}
	var totalCost, totalUnits float64
	for _, order := range orders {
		if order.OrderType == buyOrder {
			totalCost += order.Price + order.Commission
			totalUnits += order.Quantity - order.QuantityRemaining
			if totalUnits >= coin.Balance {
				break
			}
		}
	}
	coin.AvgPrice = totalCost / totalUnits
}

func (coin *WalletCoin) GetPrices() {
	clnt, err := client.GetClient()
	if err != nil {
		log.Printf("Cant get prices for %s", coin.Ticker)
		return
	}
	ticker, err := clnt.GetTicker("BTC-" + coin.Ticker)
	if err != nil {
		log.Printf("Cant get prices for %s", coin.Ticker)
		return
	}
	coin.Last = ticker.Last
	coin.Bid = ticker.Bid
	coin.Ask = ticker.Ask
}

func (coin *WalletCoin) FetchInfo() {
	avgDone := make(chan bool)
	pricesDone := make(chan bool)
	go func() {
		coin.RefreshAvgPrice()
		avgDone <- true
	}()
	go func() {
		coin.GetPrices()
		pricesDone <- true
	}()
	<-avgDone
	<-pricesDone
}
