package framework

import (
	"github.com/bigmeech/tradingbot/pkg/types"
	"log"
)

// Framework manages connectors, indicators, and middleware for the bot.
type Framework struct {
	storeManager *StoreManager              // Manages both fast and historical data storage
	connectors   map[string]types.Connector // Registered connectors
	indicators   map[string]map[string][]types.Indicator
	middleware   map[string]map[string][]types.Middleware
}

// NewFramework initializes a new Framework with StoreManager and configuration.
func NewFramework(storeManager *StoreManager) *Framework {
	return &Framework{
		storeManager: storeManager,
		connectors:   make(map[string]types.Connector),
		indicators:   make(map[string]map[string][]types.Indicator),
		middleware:   make(map[string]map[string][]types.Middleware),
	}
}

// RegisterConnector registers a connector and handles connection errors.
func (f *Framework) RegisterConnector(name string, connector types.Connector) {
	if err := connector.Connect(); err != nil {
		log.Printf("Failed to connect to %s: %v\n", name, err)
		return
	}
	f.connectors[name] = connector
}

// RegisterMiddleware adds middleware for a specific market and trading pair.
func (f *Framework) RegisterMiddleware(marketName, tradingPair string, mw types.Middleware) {
	if f.middleware[marketName] == nil {
		f.middleware[marketName] = make(map[string][]types.Middleware)
	}
	f.middleware[marketName][tradingPair] = append(f.middleware[marketName][tradingPair], mw)
}

// QueryPriceHistory retrieves price history for a specific market and trading pair.
func (f *Framework) QueryPriceHistory(market, tradingPair string, period int) []float64 {
	return f.storeManager.QueryPriceHistory(market, tradingPair, period)
}

// GetMiddleware retrieves middleware for a given market and trading pair.
func (f *Framework) GetMiddleware(marketName, tradingPair string) []types.Middleware {
	return f.middleware[marketName][tradingPair]
}

// RegisterIndicator registers an indicator for a specific market and trading pair.
func (f *Framework) RegisterIndicator(marketName, tradingPair string, indicator types.Indicator) {
	if f.indicators[marketName] == nil {
		f.indicators[marketName] = make(map[string][]types.Indicator)
	}
	f.indicators[marketName][tradingPair] = append(f.indicators[marketName][tradingPair], indicator)
}

// GetIndicators retrieves indicators for a given market and trading pair.
func (f *Framework) GetIndicators(marketName, tradingPair string) []types.Indicator {
	return f.indicators[marketName][tradingPair]
}

// executeMiddleware calculates indicators and then runs all middleware for a specific market and trading pair.
func (f *Framework) executeMiddleware(ctx *types.TickContext) error {
	// Calculate indicators for the trading pair and store in context
	for _, indicator := range f.GetIndicators(ctx.MarketName, ctx.TradingPair) {
		period := indicator.Period() // Use the indicator's period to get historical data
		priceHistory := ctx.Store.QueryPriceHistory(ctx.TradingPair, period)
		ctx.Indicators[indicator.Name()] = indicator.Calculate(priceHistory)
	}

	// Run middleware
	mws := f.GetMiddleware(ctx.MarketName, ctx.TradingPair)
	for _, mw := range mws {
		if err := mw(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Connectors returns all registered connectors.
func (f *Framework) Connectors() map[string]types.Connector {
	return f.connectors
}

// Start initiates the connectors, applying middleware to each tick through processTickFunc.
func (f *Framework) Start(processTickFunc func(ctx *types.TickContext)) {
	for _, connector := range f.connectors {
		go func(connector types.Connector) {
			connector.StreamMarketData(func(ctx *types.TickContext) {
				// Calculate indicators and run middleware, then process the tick
				if err := f.executeMiddleware(ctx); err != nil {
					log.Printf("Middleware error for %s: %v\n", ctx.TradingPair, err)
					return
				}
				processTickFunc(ctx)
			})
		}(connector)
	}
}
