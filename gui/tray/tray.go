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

type QmlBridge struct {
	core.QObject
}

func Launch(wlt wallet.Wallet) {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetQuitOnLastWindowClosed(false)

	systray := NewQSystemTrayIconWithCustomSlot(nil)
	systray.SetIcon(gui.NewQIcon5(":/qml/images/icon.png"))
	systrayMenu := widgets.NewQMenu(nil)
	for _, coin := range wlt {
		coinMenu := systrayMenu.AddMenu2(coin.Ticker)
		addCoinAction(coinMenu, "Chart", showChart)
		addCoinAction(coinMenu, "Exchange", openBittrex)
	}
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
	systray.Show()
	widgets.QApplication_Exec()
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
	settings.SetFocus2()
}

func addCoinAction(menu *widgets.QMenu, coin string, callback func(string) func(bool)) {
	action := menu.AddAction(coin)
	action.ConnectTriggered(callback(coin))
}

func showChart(coin string) func(bool) {
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
						<!-- TradingView Widget END -->`, coin)
		view := webengine.NewQWebEngineView(nil)
		view.SetHtml(html, core.NewQUrl())
		view.SetWindowTitle(coin + "BTC")
		view.Show()
	}
}

func openBittrex(coin string) func(bool) {
	return func(bool) {
		open.Start("https://bittrex.com/Market/Index?MarketName=BTC-" + coin)
	}
}
