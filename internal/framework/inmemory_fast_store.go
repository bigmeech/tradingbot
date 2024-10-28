package framework

import (
	"github.com/bigmeech/tradingbot/pkg/types"
)

// InMemoryFastStore is an in-memory circular buffer optimized for fast data access.
type InMemoryFastStore struct {
	data  []types.MarketData
	limit int
	index int
	full  bool
}

// NewInMemoryFastStore initializes an InMemoryFastStore with a fixed size.
func NewInMemoryFastStore(limit int) *InMemoryFastStore {
	return &InMemoryFastStore{
		data:  make([]types.MarketData, limit),
		limit: limit,
	}
}

// RecordTick records new market data in constant time.
func (f *InMemoryFastStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	f.data[f.index] = *marketData
	f.index = (f.index + 1) % f.limit
	if f.index == 0 {
		f.full = true
	}
	return nil
}

// QueryPriceHistory retrieves recent price history.
func (f *InMemoryFastStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	if !f.full && period > f.index {
		period = f.index
	}
	if period > f.limit {
		period = f.limit
	}

	prices := make([]float64, period)
	start := (f.index - period + f.limit) % f.limit

	for i := 0; i < period; i++ {
		prices[i] = f.data[(start+i)%f.limit].Price
	}
	return prices
}
