package main

import (
	"github.com/sazor/bittrex-wallet/config"
	"github.com/sazor/bittrex-wallet/gui"
)

func main() {
	config.Load("")
	gui.Launch()
}
