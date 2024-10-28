package framework

import (
	"github.com/bigmeech/tradingbot/pkg/types"
	"testing"
)

// MockLargeStore simulates a persistent store for testing.
type MockLargeStore struct {
	recordedData []types.MarketData
}

func (m *MockLargeStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	m.recordedData = append(m.recordedData, *marketData)
	return nil
}

func (m *MockLargeStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	history := make([]float64, 0, period)
	start := len(m.recordedData) - period
	if start < 0 {
		start = 0
	}
	for _, data := range m.recordedData[start:] {
		history = append(history, data.Price)
	}
	return history
}

func TestStoreManager_RecordTickAndQuery(t *testing.T) {
	// Setup
	bufferSize := 2 // Limit of recent data in circular buffer
	threshold := 2  // Period threshold for fastStore vs largeStore
	largeStore := &MockLargeStore{}
	manager := NewStoreManager(largeStore, bufferSize, threshold)

	// Record three ticks
	manager.RecordTick("Market1", "BTC/USDT", &types.MarketData{Price: 100.0, Volume: 1.0})
	manager.RecordTick("Market1", "BTC/USDT", &types.MarketData{Price: 200.0, Volume: 1.0})
	manager.RecordTick("Market1", "BTC/USDT", &types.MarketData{Price: 300.0, Volume: 1.0})

	// Verify fastStore holds only last `bufferSize` entries (circular buffer behavior)
	if buffer := manager.fastStore["Market1_BTC/USDT"]; buffer != nil {
		if len(buffer.GetData(bufferSize)) != bufferSize {
			t.Fatalf("Expected %v entries in fastStore, got %v", bufferSize, len(buffer.GetData(bufferSize)))
		}
	}

	// Verify fastStore content (should contain the two most recent prices: 200.0, 300.0)
	expectedFastStore := []float64{200.0, 300.0}
	fastPrices := manager.QueryPriceHistory("Market1", "BTC/USDT", threshold)
	for i, price := range fastPrices {
		if price != expectedFastStore[i] {
			t.Errorf("Expected price %v at index %d in fastStore, got %v", expectedFastStore[i], i, price)
		}
	}

	// Verify largeStore holds all three entries
	if len(largeStore.recordedData) != 3 {
		t.Fatalf("Expected 3 entries in largeStore, got %v", len(largeStore.recordedData))
	}

	// Query beyond threshold - should use largeStore
	expectedLargeStore := []float64{100.0, 200.0, 300.0}
	largePrices := manager.QueryPriceHistory("Market1", "BTC/USDT", 3)
	for i, price := range largePrices {
		if price != expectedLargeStore[i] {
			t.Errorf("Expected price %v at index %d in largeStore, got %v", expectedLargeStore[i], i, price)
		}
	}
}
