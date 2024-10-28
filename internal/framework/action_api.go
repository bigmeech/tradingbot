package framework

import (
	"github.com/bigmeech/tradingbot/pkg/types"
)

type ActionAPI struct {
	MarketName   string
	ExecuteOrder func(orderType types.OrderType, side types.OrderSide, tradingPair string, amount float64, price float64) error
}

// Buy is a helper function to execute a market buy order.
func (a *ActionAPI) Buy(amount float64, price float64) error {
	return a.ExecuteOrder(types.OrderTypeMarket, types.OrderSideBuy, a.MarketName, amount, price)
}

// Sell is a helper function to execute a market sell order.
func (a *ActionAPI) Sell(amount float64, price float64) error {
	return a.ExecuteOrder(types.OrderTypeMarket, types.OrderSideSell, a.MarketName, amount, price)
}
