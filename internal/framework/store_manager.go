package framework

import (
	"trading-bot/pkg/types"
)

// StoreManager manages interactions with fast and large stores based on a time threshold.
type StoreManager struct {
	fastStore  types.Store // Fast, in-memory store for quick access
	largeStore types.Store // Large, persistent store for long-term storage
	threshold  int         // Threshold (in periods) for switching between stores
}

// NewStoreManager initializes a new StoreManager with given fast and large stores, and a threshold.
func NewStoreManager(fastStore, largeStore types.Store, threshold int) *StoreManager {
	return &StoreManager{
		fastStore:  fastStore,
		largeStore: largeStore,
		threshold:  threshold,
	}
}

// RecordTick records a tick in the fast store and optionally in the large store.
func (s *StoreManager) RecordTick(tradingPair string, marketData *types.MarketData) error {
	// Record in fast store for quick access
	if err := s.fastStore.RecordTick(tradingPair, marketData); err != nil {
		return err
	}

	// Optionally, record in large store for long-term storage
	return s.largeStore.RecordTick(tradingPair, marketData)
}

// QueryPriceHistory retrieves price history, using the fast store if within threshold,
// and the large store if beyond threshold.
func (s *StoreManager) QueryPriceHistory(tradingPair string, period int) []float64 {
	// Use fast store for recent data (within threshold)
	if period <= s.threshold {
		return s.fastStore.QueryPriceHistory(tradingPair, period)
	}

	// Use large store for older data (beyond threshold)
	return s.largeStore.QueryPriceHistory(tradingPair, period)
}
