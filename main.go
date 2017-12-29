package main

import (
	"log"

	"github.com/sazor/bittrex-wallet/config"
	"github.com/sazor/bittrex-wallet/gui/tray"
	"github.com/sazor/bittrex-wallet/services/wallet"
)

func main() {
	config.Load("")
	wlt, err := wallet.LoadDetailedWallet()
	if err != nil {
		log.Fatalln(err)
		return
	}
	wallet.Sort(wlt, wallet.ChangeSort, wallet.Asc)
	tray.Launch(wlt)
}
