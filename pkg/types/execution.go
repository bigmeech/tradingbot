package types

// OrderExecutor interface defines the method signature for executing an order.
type OrderExecutor interface {
	// ExecuteOrder places an order with the specified type, side, trading pair, amount, and price.
	ExecuteOrder(orderType OrderType, side OrderSide, tradingPair string, amount float64, price float64) error
}
