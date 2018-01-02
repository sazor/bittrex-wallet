package tray

import (
	"fmt"
	"os"

	"github.com/sazor/bittrex-wallet/services/wallet"
	"github.com/skratchdot/open-golang/open"
	"github.com/therecipe/qt/core"
	gui "github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/quickcontrols2"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
)

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon

	_ func() `slot:"triggerSlot"`
}

type QMenuWithSlot struct {
	widgets.QMenu

	_ func(newTitle string) `slot:"changeTitleSlot"`
}

type QmlBridge struct {
	core.QObject
}

func Launch(wlt *wallet.Wallet) {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetQuitOnLastWindowClosed(false)

	systray := NewQSystemTrayIconWithCustomSlot(nil)
	systray.SetIcon(gui.NewQIcon5(":/qml/images/icon32.png"))
	systrayMenu := widgets.NewQMenu(nil)
	bitcoinMenu := NewQMenuWithSlot(nil)
	bitcoinMenu.SetTitle(fmt.Sprintf("%0.8f BTC", wlt.EstimatedBtcBalance()))
	bitcoinMenu.SetIcon(gui.NewQIcon5(":/qml/images/logo/USDT-BTC.png"))
	bitcoinMenu.ConnectChangeTitleSlot(func(newTitle string) {
		bitcoinMenu.SetTitle(newTitle)
	})
	for ticker, coin := range wlt.Altcoins {
		market := fmt.Sprintf("BTC-%s", ticker)
		coinMenu := NewQMenuWithSlot(nil)
		coinMenu.SetTitle(fmt.Sprintf("%s | %0.8f", market, coin.Last))
		coinMenu.SetIcon(gui.NewQIcon5(fmt.Sprintf(":/qml/images/logo/%s.png", market)))
		coinMenu.ConnectChangeTitleSlot(func(newTitle string) {
			coinMenu.SetTitle(newTitle)
		})
		systrayMenu.AddMenu(coinMenu)
		updateMenuTitle(coinMenu, coin)
		updateBitcoinBalance(bitcoinMenu, wlt, coin)
		addCoinAction(coinMenu, "Chart", coin, showChart)
		addCoinAction(coinMenu, "Exchange", coin, openBittrex)
	}
	systrayMenu.AddSeparator()
	systrayMenu.AddMenu(bitcoinMenu)
	systrayMenu.AddSeparator()
	settingsAction := systrayMenu.AddAction("Settings")
	settingsAction.ConnectTriggered(func(bool) {
		showSettings()
	})
	quitAction := systrayMenu.AddAction("Quit")
	quitAction.ConnectTriggered(func(bool) {
		app.Quit()
	})
	systray.SetContextMenu(systrayMenu)
	systrayMenu.ConnectAboutToShow(func() {
		go wlt.SubscribeToUpdates()
	})
	systrayMenu.ConnectAboutToHide(func() {
		go wlt.UnsubscribeFromUpdates()
	})
	systray.Show()
	widgets.QApplication_Exec()
}

func updateMenuTitle(menu *QMenuWithSlot, coin *wallet.WalletCoin) {
	listener := coin.NewListener()
	go func() {
		for {
			select {
			case <-listener:
				menu.ChangeTitleSlot(fmt.Sprintf("BTC-%s | %0.8f", coin.Ticker, coin.Last))
			}
		}
	}()
}

func updateBitcoinBalance(menu *QMenuWithSlot, wlt *wallet.Wallet, coin *wallet.WalletCoin) {
	listener := coin.NewListener()
	go func() {
		for {
			select {
			case <-listener:
				menu.ChangeTitleSlot(fmt.Sprintf("%0.8f BTC", wlt.EstimatedBtcBalance()))
			}
		}
	}()
}

func showSettings() {
	//var qmlBridge = NewQmlBridge(nil)
	quickcontrols2.QQuickStyle_SetStyle("Material")
	settings := quick.NewQQuickWidget3(core.NewQUrl3("qrc:/qml/settings.qml", 0), nil)
	//view.RootContext().SetContextProperty("QmlBridge", qmlBridge)
	settings.SetResizeMode(quick.QQuickWidget__SizeRootObjectToView)
	settings.SetWindowTitle("Settings")
	settings.SetWindowFlags(core.Qt__Window | core.Qt__WindowTitleHint | core.Qt__WindowCloseButtonHint | core.Qt__CustomizeWindowHint)
	settings.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
	settings.SetFixedHeight(380)
	settings.SetFixedWidth(350)
	settings.Show()
}

func addCoinAction(menu *QMenuWithSlot, name string, coin *wallet.WalletCoin, callback func(*wallet.WalletCoin) func(bool)) {
	action := menu.AddAction(name)
	action.ConnectTriggered(callback(coin))
}

func showChart(coin *wallet.WalletCoin) func(bool) {
	return func(bool) {
		html := fmt.Sprintf(`<!-- TradingView Widget BEGIN -->
					<script type="text/javascript" src="https://s3.tradingview.com/tv.js"></script>
					<script type="text/javascript">
					new TradingView.widget({
					  "autosize": true,
					  "symbol": "BITTREX:%sBTC",
						  "interval": "D",
						  "timezone": "Etc/UTC",
						  "theme": "Dark",
						  "style": "1",
						  "locale": "en",
						  "toolbar_bg": "#f1f3f6",
						  "enable_publishing": false,
						  "withdateranges": true,
						  "hide_side_toolbar": false,
						  "allow_symbol_change": true,
						  "hideideas": true
						});
						</script>
						<!-- TradingView Widget END -->`, coin.Ticker)
		view := webengine.NewQWebEngineView(nil)
		view.SetHtml(html, core.NewQUrl())
		view.SetWindowTitle(coin.Ticker + "BTC")
		view.Show()
	}
}

func openBittrex(coin *wallet.WalletCoin) func(bool) {
	return func(bool) {
		open.Start("https://bittrex.com/Market/Index?MarketName=BTC-" + coin.Ticker)
	}
}
