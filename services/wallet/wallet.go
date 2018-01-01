package wallet

import (
	"errors"
	"log"
	"sync"

	"github.com/sazor/bittrex-wallet/services/client"
	bittrex "github.com/sazor/go-bittrex"
)

type Wallet struct {
	Altcoins      map[string]*WalletCoin
	stopUpdatesCh chan bool
}

func (wlt *Wallet) SubscribeToUpdates() {
	clnt, err := client.GetClient()
	if err != nil {
		log.Println(err)
		return
	}
	priceChannels := make(map[string]chan bittrex.SummaryState, len(wlt.Altcoins))
	for ticker, coin := range wlt.Altcoins {
		priceChannels["BTC-"+ticker] = coin.SubscribeToUpdates()
	}
	errCh := make(chan error)
	stopWsCh := make(chan bool)
	go func() {
		errCh <- clnt.SubscribeSummaryDeltas(priceChannels, stopWsCh)
	}()
	for {
		select {
		case <-wlt.stopUpdatesCh:
			stopWsCh <- true
			return
		case err := <-errCh:
			log.Println(err)
			return
		}
	}
}

func (wlt *Wallet) UnsubscribeFromUpdates() {
	wlt.stopUpdatesCh <- true
	for _, coin := range wlt.Altcoins {
		coin.StopSubscription()
	}
}

func NewWallet(balances []bittrex.Balance) *Wallet {
	wlt := &Wallet{stopUpdatesCh: make(chan bool),
		Altcoins: make(map[string]*WalletCoin, len(balances)),
	}
	for _, coin := range balances {
		balance, _ := coin.Balance.Float64()
		if balance > 0.0 && coin.Currency != "BTC" && coin.Currency != "USDT" {
			wlt.Altcoins[coin.Currency] = NewWalletCoin(coin.Currency, balance)
		}
	}
	return wlt
}

func LoadWallet() (*Wallet, error) {
	clnt, err := client.GetClient()
	if err != nil {
		log.Fatal("Connection issues: %+v", err)
		return nil, errors.New("Connection issue")
	}
	balances, err := clnt.GetBalances()
	if err != nil {
		log.Fatal("Connection issues: %+v", err)
		return nil, errors.New("Connection issue")
	}
	return NewWallet(balances), nil
}

func LoadDetailedWallet() (*Wallet, error) {
	wlt, err := LoadWallet()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(len(wlt.Altcoins))
	for _, coin := range wlt.Altcoins {
		go func(c *WalletCoin) {
			c.FetchInfo()
			wg.Done()
		}(coin)
	}
	wg.Wait()
	return wlt, nil
}

func LoadDetailedCoin(wlt *Wallet, ticker string) (*WalletCoin, error) {
	coin, ok := wlt.Altcoins[ticker]
	if !ok {
		log.Println("No such coin in the wallet")
		return nil, errors.New("No such coin in the wallet")
	}
	coin.FetchInfo()
	return coin, nil
}
