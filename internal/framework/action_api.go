package framework

import (
	"trading-bot/pkg/types"
)

type ActionAPI struct {
	MarketName    string
	ExecuteAction func(action types.ActionType, tradingPair string, amount float64) error
}

// Buy is a helper function to execute a buy action.
func (a *ActionAPI) Buy(amount float64) error {
	return a.ExecuteAction(types.ActionBuy, a.MarketName, amount)
}

// Sell is a helper function to execute a sell action.
func (a *ActionAPI) Sell(amount float64) error {
	return a.ExecuteAction(types.ActionSell, a.MarketName, amount)
}
