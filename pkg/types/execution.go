package types

type OrderExecutor interface {
	ExecuteOrder(action ActionType, tradingPair string, amount float64) error
}
