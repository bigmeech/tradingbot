package types

type Connector interface {
	Connect() error
	StreamMarketData(func(*TickContext)) error

	// ExecuteOrder places an order with the specified type, side, amount, and price.
	ExecuteOrder(orderType OrderType, side OrderSide, tradingPair string, amount, price float64) error
}

type Indicator interface {
	Calculate(data []float64) float64
	Name() string
}

type Middleware func(*TickContext) error

type Store interface {
	RecordTick(tradingPair string, marketData *MarketData) error
	QueryPriceHistory(tradingPair string, period int) []float64
}

// OrderSide represents the side of the order, either "buy" or "sell"
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents the type of order being placed
type OrderType string

const (
	// OrderTypeMarket represents a market order
	// which is executed immediately at the best available price.
	OrderTypeMarket OrderType = "market"

	// OrderTypeLimit represents a limit order
	// where the order is placed at a specific price or better.
	// The order will only execute if the market reaches the limit price.
	OrderTypeLimit OrderType = "limit"

	// OrderTypeStopLoss represents a stop-loss order
	// which triggers a market order when the price reaches a specified stop price
	// to limit losses from an unfavorable price movement.
	OrderTypeStopLoss OrderType = "stop-loss"

	// OrderTypeStopLossLimit represents a stop-loss limit order
	// which triggers a limit order at a specified price when the price reaches
	// a designated stop price. This limits losses, but only executes at or above
	// the specified limit price.
	OrderTypeStopLossLimit OrderType = "stop-loss-limit"

	// OrderTypeTakeProfit represents a take-profit order
	// which triggers a market order when the price reaches a specified level.
	// It is used to lock in profits from a favorable price movement.
	OrderTypeTakeProfit OrderType = "take-profit"

	// OrderTypeTakeProfitLimit represents a take-profit limit order
	// which triggers a limit order at a specified price when the price reaches
	// a designated level. It is used to secure profits, but only executes
	// at or above the specified limit price.
	OrderTypeTakeProfitLimit OrderType = "take-profit-limit"

	// OrderTypeTrailingStop represents a trailing stop order
	// where a stop price is set at a fixed amount below the market price
	// for a sell or above the market price for a buy. The stop price
	// adjusts with favorable market movements.
	OrderTypeTrailingStop OrderType = "trailing-stop"

	// OrderTypeTrailingStopLimit represents a trailing stop-limit order
	// which sets a stop price at a fixed distance from the market price,
	// moving with favorable price changes, and triggers a limit order
	// when the stop price is reached.
	OrderTypeTrailingStopLimit OrderType = "trailing-stop-limit"

	// OrderTypeIceberg represents an iceberg order
	// where only a portion of the total order size is displayed in the order book.
	// The full order is divided into smaller visible orders.
	OrderTypeIceberg OrderType = "iceberg"

	// OrderTypeSettlePosition represents a settle position order
	// typically used to settle or close an open position.
	OrderTypeSettlePosition OrderType = "settle-position"
)

type TickContext struct {
	MarketName  string
	TradingPair string
	MarketData  *MarketData
	Store       Store
	Indicators  map[string]float64

	// ExecuteOrder function to place orders with order_type and side
	ExecuteOrder func(orderType OrderType, side OrderSide, amount, price float64) error
}

type MarketData struct {
	Price  float64
	Volume float64
}
