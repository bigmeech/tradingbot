package testutils

import (
	"github.com/bigmeech/tradingbot/pkg/types"
	"sync"
)

// MockStore simulates an in-memory store for testing purposes.
type MockStore struct {
	mu           sync.Mutex
	recordedData map[string][]types.MarketData
}

// NewMockStore initializes a new MockStore.
func NewMockStore() *MockStore {
	return &MockStore{
		recordedData: make(map[string][]types.MarketData),
	}
}

// RecordTick simulates recording a tick for a trading pair.
func (m *MockStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordedData[tradingPair] = append(m.recordedData[tradingPair], *marketData)
	return nil
}

// QueryPriceHistory simulates querying the price history for a trading pair.
func (m *MockStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	history := make([]float64, 0, period)
	data := m.recordedData[tradingPair]
	count := len(data)
	start := count - period
	if start < 0 {
		start = 0
	}
	for _, marketData := range data[start:] {
		history = append(history, marketData.Price)
	}
	return history
}
