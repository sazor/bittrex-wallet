package wallet

import (
	"errors"
	"log"

	"github.com/sazor/bittrex-wallet/services/client"
	"github.com/sazor/go-bittrex"
)

const buyOrder = "LIMIT_BUY"

type WalletCoin struct {
	Ticker       string
	Balance      float64
	AvgPrice     float64
	Last         float64
	Bid          float64
	Ask          float64
	updatesCh    chan bittrex.SummaryState
	stopUpdateCh chan bool
	listeners    []chan *WalletCoin
}

func NewWalletCoin(ticker string, balance float64) *WalletCoin {
	return &WalletCoin{Ticker: ticker,
		Balance:      balance,
		updatesCh:    make(chan bittrex.SummaryState),
		stopUpdateCh: make(chan bool),
	}
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
			price, _ := order.Price.Float64()
			commission, _ := order.Commission.Float64()
			quantity, _ := order.Quantity.Float64()
			quantityRem, _ := order.QuantityRemaining.Float64()
			totalCost += price + commission
			totalUnits += quantity - quantityRem
			if totalUnits >= coin.Balance {
				break
			}
		}
	}
	coin.AvgPrice = totalCost / totalUnits
}

func (coin *WalletCoin) RefreshPrices() {
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
	coin.Last, _ = ticker.Last.Float64()
	coin.Bid, _ = ticker.Bid.Float64()
	coin.Ask, _ = ticker.Ask.Float64()
}

func (coin *WalletCoin) FetchInfo() {
	avgDone := make(chan bool)
	pricesDone := make(chan bool)
	go func() {
		coin.RefreshAvgPrice()
		avgDone <- true
	}()
	go func() {
		coin.RefreshPrices()
		pricesDone <- true
	}()
	<-avgDone
	<-pricesDone
}

func (coin *WalletCoin) SubscribeToUpdates() chan bittrex.SummaryState {
	go func() {
		for {
			select {
			case <-coin.stopUpdateCh:
				return
			case state := <-coin.updatesCh:
				coin.updateInfo(state)
			}
		}
	}()
	return coin.updatesCh
}

func (coin *WalletCoin) StopSubscription() {
	coin.stopUpdateCh <- true
}

func (coin *WalletCoin) GetUpdatesChannel() (chan bittrex.SummaryState, error) {
	if coin.updatesCh != nil {
		return coin.updatesCh, nil
	} else {
		return nil, errors.New("Updates channel isnt initialized")
	}
}

func (coin *WalletCoin) AddListener(listener chan *WalletCoin) {
	coin.listeners = append(coin.listeners, listener)
}

func (coin *WalletCoin) NewListener() chan *WalletCoin {
	listener := make(chan *WalletCoin)
	coin.AddListener(listener)
	return listener
}

func (coin *WalletCoin) updateInfo(updatedState bittrex.SummaryState) {
	coin.Last = updatedState.Last
	coin.Bid = updatedState.Bid
	coin.Ask = updatedState.Ask
	for _, listener := range coin.listeners {
		listener <- coin
	}
}
