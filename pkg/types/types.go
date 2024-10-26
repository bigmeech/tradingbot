package types

type Connector interface {
	Connect() error
	StreamMarketData(func(*TickContext)) error
	ExecuteAction(action ActionType, tradingPair string, amount float64) error
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

type TickContext struct {
	TradingPair string
	MarketData  *MarketData
	Actions     ActionAPI
	Store       Store
	Indicators  map[string]float64
}

type ActionAPI struct {
	MarketName    string
	ExecuteAction func(action ActionType, tradingPair string, amount float64) error
}

type MarketData struct {
	Price  float64
	Volume float64
}

type ActionType string

const (
	ActionBuy  ActionType = "BUY"
	ActionSell ActionType = "SELL"
)
