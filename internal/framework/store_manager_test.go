package framework

import (
	"testing"
	"trading-bot/pkg/types"
)

// MockFastStore is a mock implementation of FastStore for testing.
type MockFastStore struct {
	recordedData []types.MarketData
}

func (m *MockFastStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	m.recordedData = append(m.recordedData, *marketData)
	return nil
}

func (m *MockFastStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	history := []float64{}
	count := len(m.recordedData)
	start := count - period
	if start < 0 {
		start = 0
	}
	for _, data := range m.recordedData[start:] {
		history = append(history, data.Price)
	}
	return history
}

// MockLargeStore is a mock implementation of LargeStore for testing.
type MockLargeStore struct {
	recordedData []types.MarketData
}

func (m *MockLargeStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	m.recordedData = append(m.recordedData, *marketData)
	return nil
}

func (m *MockLargeStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	history := []float64{}
	count := len(m.recordedData)
	start := count - period
	if start < 0 {
		start = 0
	}
	for _, data := range m.recordedData[start:] {
		history = append(history, data.Price)
	}
	return history
}

func TestStoreManager_RecordTickAndQuery(t *testing.T) {
	fastStore := &MockFastStore{}
	largeStore := &MockLargeStore{}
	threshold := 2
	manager := NewStoreManager(fastStore, largeStore, threshold)

	// Add ticks to the StoreManager
	manager.RecordTick("BTC/USDT", &types.MarketData{Price: 100.0, Volume: 1.0})
	manager.RecordTick("BTC/USDT", &types.MarketData{Price: 200.0, Volume: 1.0})
	manager.RecordTick("BTC/USDT", &types.MarketData{Price: 300.0, Volume: 1.0})

	// Verify that data was recorded in both stores
	if len(fastStore.recordedData) != 3 {
		t.Fatalf("Expected 3 entries in FastStore, got %v", len(fastStore.recordedData))
	}
	if len(largeStore.recordedData) != 3 {
		t.Fatalf("Expected 3 entries in LargeStore, got %v", len(largeStore.recordedData))
	}

	// Query with period below threshold - should use FastStore
	prices := manager.QueryPriceHistory("BTC/USDT", 2)
	expectedFast := []float64{200.0, 300.0}

	if len(prices) != len(expectedFast) {
		t.Fatalf("Expected %v prices from FastStore, got %v", len(expectedFast), len(prices))
	}

	for i, price := range prices {
		if price != expectedFast[i] {
			t.Errorf("Expected price %v at index %v, got %v", expectedFast[i], i, price)
		}
	}

	// Query with period above threshold - should use LargeStore
	prices = manager.QueryPriceHistory("BTC/USDT", 3)
	expectedLarge := []float64{100.0, 200.0, 300.0}

	if len(prices) != len(expectedLarge) {
		t.Fatalf("Expected %v prices from LargeStore, got %v", len(expectedLarge), len(prices))
	}

	for i, price := range prices {
		if price != expectedLarge[i] {
			t.Errorf("Expected price %v at index %v, got %v", expectedLarge[i], i, price)
		}
	}
}
