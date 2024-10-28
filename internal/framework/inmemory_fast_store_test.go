package framework

import (
	"github.com/bigmeech/tradingbot/pkg/types"
	"testing"
)

func TestInMemoryFastStore_RecordTickAndQuery(t *testing.T) {
	store := NewInMemoryFastStore(3) // Limit of 3 for testing circular buffer

	// Add three ticks
	store.RecordTick("BTC/USDT", &types.MarketData{Price: 100.0})
	store.RecordTick("BTC/USDT", &types.MarketData{Price: 200.0})
	store.RecordTick("BTC/USDT", &types.MarketData{Price: 300.0})

	// Retrieve all three prices
	prices := store.QueryPriceHistory("BTC/USDT", 3)
	expected := []float64{100.0, 200.0, 300.0}

	if len(prices) != len(expected) {
		t.Fatalf("Expected %v prices, got %v", len(expected), len(prices))
	}

	for i, price := range prices {
		if price != expected[i] {
			t.Errorf("Expected price %v at index %v, got %v", expected[i], i, price)
		}
	}

	// Add another tick, which should overwrite the oldest entry
	store.RecordTick("BTC/USDT", &types.MarketData{Price: 400.0})
	prices = store.QueryPriceHistory("BTC/USDT", 3)
	expected = []float64{200.0, 300.0, 400.0} // Circular buffer should now contain the last three entries

	for i, price := range prices {
		if price != expected[i] {
			t.Errorf("Expected price %v at index %v, got %v", expected[i], i, price)
		}
	}
}
