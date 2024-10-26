package framework

import (
	"fmt"
	"trading-bot/pkg/types"
)

type Framework struct {
	store      types.Store
	connectors map[string]types.Connector
	indicators map[string]map[string][]types.Indicator
	middleware map[string]map[string][]types.Middleware
}

// NewFramework initializes the Framework with a provided StoreManager.
func NewFramework(store types.Store) *Framework {
	return &Framework{
		store:      store,
		connectors: make(map[string]types.Connector),
		indicators: make(map[string]map[string][]types.Indicator),
		middleware: make(map[string]map[string][]types.Middleware),
	}
}

// RegisterConnector registers a trading connector.
func (f *Framework) RegisterConnector(name string, connector types.Connector) {
	f.connectors[name] = connector
	connector.Connect()
}

// Connectors returns all registered connectors.
func (f *Framework) Connectors() map[string]types.Connector {
	return f.connectors
}

// RegisterIndicator associates indicators with a specific market and trading pair.
func (f *Framework) RegisterIndicator(marketName, tradingPair string, indicator types.Indicator) {
	if f.indicators[marketName] == nil {
		f.indicators[marketName] = make(map[string][]types.Indicator)
	}
	f.indicators[marketName][tradingPair] = append(f.indicators[marketName][tradingPair], indicator)
}

// UsePair registers middleware (e.g., strategies) for a specific market and trading pair.
func (f *Framework) UsePair(marketName, tradingPair string, middleware types.Middleware) {
	if f.middleware[marketName] == nil {
		f.middleware[marketName] = make(map[string][]types.Middleware)
	}
	f.middleware[marketName][tradingPair] = append(f.middleware[marketName][tradingPair], middleware)
}

// Start begins the framework and starts processing ticks.
func (f *Framework) Start() {
	for _, connector := range f.connectors {
		connector.StreamMarketData(func(ctx *types.TickContext) {
			f.processTick(ctx)
		})
	}
}

// processTick processes each tick by querying data from the store and applying indicators and strategies.
func (f *Framework) processTick(ctx *types.TickContext) {
	priceData := f.store.QueryPriceHistory(ctx.TradingPair, 200) // Query recent prices
	ctx.Indicators = f.calculateIndicators(ctx.Actions.MarketName, ctx.TradingPair, priceData)

	for _, m := range f.middleware[ctx.Actions.MarketName][ctx.TradingPair] {
		if err := m(ctx); err != nil {
			fmt.Println("Error in middleware:", err)
		}
	}
}

// calculateIndicators calculates all registered indicators for a trading pair.
func (f *Framework) calculateIndicators(marketName, tradingPair string, data []float64) map[string]float64 {
	results := make(map[string]float64)
	if indicators, exists := f.indicators[marketName][tradingPair]; exists {
		for _, indicator := range indicators {
			results[indicator.Name()] = indicator.Calculate(data)
		}
	}
	return results
}
