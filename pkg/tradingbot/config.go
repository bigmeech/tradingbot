package tradingbot

type BotConfig struct {
	Connectors []ConnectorConfig // Configurations for each connector
	Indicators []IndicatorConfig // Configurations for each indicator
	Strategies []StrategyConfig  // Configurations for each strategy
}

type ConnectorConfig struct {
	Name      string // e.g., "Binance"
	APIKey    string
	APISecret string
	WSURL     string
	RestURL   string
}

type IndicatorConfig struct {
	MarketName  string // e.g., "Binance"
	TradingPair string // e.g., "BTC/USDT"
	Type        string // e.g., "SMA", "RSI"
	Period      int
}

type StrategyConfig struct {
	MarketName  string // e.g., "Binance"
	TradingPair string // e.g., "BTC/USDT"
	Type        string // e.g., "MA_Crossover", "RSI"
}
