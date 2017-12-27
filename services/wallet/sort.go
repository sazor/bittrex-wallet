package wallet

import "sort"

func Sort(w Wallet, field string, direction string) {
	ascSort := direction != "desc"
	switch field {
	case "curprice":
		sortByCurrPrice(w, ascSort)
	case "avgprice":
		sortByAvgPrice(w, ascSort)
	case "change":
		sortByChange(w, ascSort)
	case "balance":
		sortByBalance(w, ascSort)
	default:
		sortByBalance(w, ascSort)
	}
}

func sortByCurrPrice(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].CurrPrice < wallet[j].CurrPrice
	})
}

func sortByAvgPrice(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].AvgPrice < wallet[j].AvgPrice
	})
}

func sortByBalance(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].btcBalance() < wallet[j].btcBalance()
	})
}

func sortByChange(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].percentDiffPrice() < wallet[j].percentDiffPrice()
	})
}
