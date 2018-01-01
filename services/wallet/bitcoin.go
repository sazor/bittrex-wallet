package wallet

import (
	"errors"
	"log"

	"github.com/sazor/bittrex-wallet/services/client"
	bittrex "github.com/sazor/go-bittrex"
)

type Bitcoin struct {
	Ticker       string
	Balance      float64
	Last         float64
	Bid          float64
	Ask          float64
	updatesCh    chan bittrex.SummaryState
	stopUpdateCh chan bool
	listeners    []chan *Bitcoin
}

func NewBitcoin(balance float64) *Bitcoin {
	return &Bitcoin{Ticker: "BTC",
		Balance:      balance,
		updatesCh:    make(chan bittrex.SummaryState),
		stopUpdateCh: make(chan bool),
	}
}

func (coin *Bitcoin) RefreshPrices() {
	clnt, err := client.GetClient()
	if err != nil {
		log.Printf("Cant get prices for %s", coin.Ticker)
		return
	}
	ticker, err := clnt.GetTicker("USDT-BTC")
	if err != nil {
		log.Printf("Cant get prices for %s", coin.Ticker)
		return
	}
	coin.Last, _ = ticker.Last.Float64()
	coin.Bid, _ = ticker.Bid.Float64()
	coin.Ask, _ = ticker.Ask.Float64()
}

func (coin *Bitcoin) SubscribeToUpdates() chan bittrex.SummaryState {
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

func (coin *Bitcoin) StopSubscription() {
	coin.stopUpdateCh <- true
}

func (coin *Bitcoin) GetUpdatesChannel() (chan bittrex.SummaryState, error) {
	if coin.updatesCh != nil {
		return coin.updatesCh, nil
	} else {
		return nil, errors.New("Updates channel isnt initialized")
	}
}

func (coin *Bitcoin) AddListener(listener chan *Bitcoin) {
	coin.listeners = append(coin.listeners, listener)
}

func (coin *Bitcoin) NewListener() chan *Bitcoin {
	listener := make(chan *Bitcoin)
	coin.AddListener(listener)
	return listener
}

func (coin *Bitcoin) updateInfo(updatedState bittrex.SummaryState) {
	coin.Last = updatedState.Last
	coin.Bid = updatedState.Bid
	coin.Ask = updatedState.Ask
	for _, listener := range coin.listeners {
		listener <- coin
	}
}
