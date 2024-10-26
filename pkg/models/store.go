package models

import "trading-bot/pkg/types"

type FastStore interface {
	RecordTick(tradingPair string, marketData *types.MarketData) error
	QueryPriceHistory(tradingPair string, period int) []float64
}

type LargeStore interface {
	RecordTick(tradingPair string, marketData *types.MarketData) error
	QueryPriceHistory(tradingPair string, period int) []float64
}
