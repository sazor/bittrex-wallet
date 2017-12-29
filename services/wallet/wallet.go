package wallet

import (
	"errors"
	"log"
	"sync"

	"github.com/sazor/bittrex-wallet/services/client"
	bittrex "github.com/toorop/go-bittrex"
)

type Wallet []*WalletCoin

func (wlt Wallet) Find(ticker string) (*WalletCoin, error) {
	for _, coin := range wlt {
		if coin.Ticker == ticker {
			return coin, nil
		}
	}
	return nil, errors.New("No such coin in wallet")
}

func NewWallet(balances []bittrex.Balance) Wallet {
	var wlt Wallet
	for _, coin := range balances {
		if coin.Balance > 0.0 && coin.Currency != "BTC" && coin.Currency != "USDT" {
			wlt = append(wlt, NewWalletCoin(coin.Currency, coin.Balance))
		}
	}
	return wlt
}

func LoadWallet() (Wallet, error) {
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

func LoadDetailedWallet() (Wallet, error) {
	wlt, err := LoadWallet()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(len(wlt))
	for _, coin := range wlt {
		go func(c *WalletCoin) {
			c.FetchInfo()
			wg.Done()
		}(coin)
	}
	wg.Wait()
	return wlt, nil
}

func LoadDetailedCoin(wlt Wallet, ticker string) (*WalletCoin, error) {
	coin, err := wlt.Find(ticker)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	coin.FetchInfo()
	return coin, nil
}
