package framework

import (
	"log"
	"trading-bot/pkg/types"
)

type Framework struct {
	store      types.Store
	connectors map[string]types.Connector
	indicators map[string]map[string][]types.Indicator
	middleware map[string]map[string][]types.Middleware
}

// NewFramework initializes a new Framework instance with a Store.
func NewFramework(store types.Store) *Framework {
	return &Framework{
		store:      store,
		connectors: make(map[string]types.Connector),
		indicators: make(map[string]map[string][]types.Indicator),
		middleware: make(map[string]map[string][]types.Middleware),
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

// GetMiddleware retrieves middleware for a given market and trading pair.
func (f *Framework) GetMiddleware(marketName, tradingPair string) []types.Middleware {
	return f.middleware[marketName][tradingPair]
}

// executeMiddleware runs all middleware for a specific market and trading pair.
func (f *Framework) executeMiddleware(ctx *types.TickContext) error {
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
				// Run middleware for each tick, then process the tick
				if err := f.executeMiddleware(ctx); err != nil {
					log.Printf("Middleware error for %s: %v\n", ctx.TradingPair, err)
					return
				}
				processTickFunc(ctx)
			})
		}(connector)
	}
}
