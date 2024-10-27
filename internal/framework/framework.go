package framework

import (
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

// RegisterConnector registers a connector.
func (f *Framework) RegisterConnector(name string, connector types.Connector) {
	f.connectors[name] = connector
	connector.Connect()
}

// GetMiddleware retrieves middleware for a given market and trading pair.
func (f *Framework) GetMiddleware(marketName, tradingPair string) []types.Middleware {
	return f.middleware[marketName][tradingPair]
}

// Connectors returns all registered connectors.
func (f *Framework) Connectors() map[string]types.Connector {
	return f.connectors
}

// Start initiates the connectors, using a provided tick-processing function.
func (f *Framework) Start(processTickFunc func(ctx *types.TickContext)) {
	for _, connector := range f.connectors {
		go func(connector types.Connector) {
			connector.StreamMarketData(processTickFunc)
		}(connector)
	}
}
