package gui

import (
	"os"

	"github.com/sazor/bittrex-notifier/data"
	"github.com/therecipe/qt/charts"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon

	_ func() `slot:"triggerSlot"`
}

func Launch() {
	pumps, _ := data.PumpsDumps()
	app := widgets.NewQApplication(len(os.Args), os.Args)
	systray := NewQSystemTrayIconWithCustomSlot(nil)
	systray.SetIcon(gui.NewQIcon5(":/qml/images/icon.png"))
	systrayMenu := widgets.NewQMenu(nil)
	for _, coin := range pumps {
		action := widgets.NewQWidgetAction(nil)
		label := widgets.NewQLabel2(coin.MarketName, nil, 0)
		label.SetMargin(20)
		action.SetDefaultWidget(label)
		func(marketName string, lbl *widgets.QLabel) {
			lbl.ConnectMousePressEvent(func(*gui.QMouseEvent) {
				series := charts.NewQCandlestickSeries(nil)
				candles, _ := data.MarketCandles(marketName)
				for _, candle := range candles {
					c := charts.NewQCandlestickSet2(candle.Open, candle.High, candle.Low, candle.Close, float64(candle.TimeStamp.UnixNano()), nil)
					series.Append(c)
				}
				series.SetName(marketName)
				series.SetIncreasingColor(gui.QColor_FromRgb2(0, 255, 0, 255))
				series.SetDecreasingColor(gui.QColor_FromRgb2(255, 0, 0, 255))
				chart := charts.NewQChart(nil, 0)
				chart.AddSeries(series)
				chart.SetAnimationOptions(charts.QChart__SeriesAnimations)
				chart.CreateDefaultAxes()
				lgd := chart.Legend()
				lgd.SetVisible(true)
				chv := charts.NewQChartView2(chart, nil)
				chv.Show()
				// chv := charts.NewQChartView(nil)
				// chv.SetChart(chart.Chart())
				// chv.Show()
			})
		}(coin.MarketName, label)
		systrayMenu.AddActions([]*widgets.QAction{action.QAction_PTR()})
	}
	systrayMenu.AddSeparator()
	quitAction := systrayMenu.AddAction("Quit")
	quitAction.ConnectTriggered(func(bool) {
		app.Quit()
	})
	systray.SetContextMenu(systrayMenu)

	systray.Show()

	//works
	// var buttonSlot = widgets.NewQPushButton2("call from other thread with slot", nil)
	// systray.ConnectTriggerSlot(func() {
	// 	systray.ShowMessage("title", "other thread message with slot", widgets.QSystemTrayIcon__Information, 5000)
	// })
	// buttonSlot.ConnectClicked(func(bool) {
	// 	go func() {
	// 		systray.TriggerSlot()
	// 	}()
	// })
	// widgetLayout.AddWidget(buttonSlot, 0, 0)
	// //
	// widget.Show()

	widgets.QApplication_Exec()
}
