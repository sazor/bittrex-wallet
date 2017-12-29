package wallet

import "sort"

const (
	CurPriceSort = iota
	AvgPriceSort
	ChangeSort
	BalanceSort
)

const (
	Asc = iota
	Desc
)

func Sort(w Wallet, field int, direction int) {
	switch field {
	case CurPriceSort:
		sortByCurrPrice(w, direction == Asc)
	case AvgPriceSort:
		sortByAvgPrice(w, direction == Asc)
	case ChangeSort:
		sortByChange(w, direction == Asc)
	case BalanceSort:
		sortByBalance(w, direction == Asc)
	default:
		sortByBalance(w, direction == Asc)
	}
}

func sortByCurrPrice(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].Last < wallet[j].Last
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
		return wallet[i].BtcBalance() < wallet[j].BtcBalance()
	})
}

func sortByChange(wallet Wallet, asc bool) {
	sort.Slice(wallet, func(i, j int) bool {
		if !asc {
			i, j = j, i
		}
		return wallet[i].PercentDiffPrice() < wallet[j].PercentDiffPrice()
	})
}
