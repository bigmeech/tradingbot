package framework

import (
	"trading-bot/pkg/models"
	"trading-bot/pkg/types"
)

// StoreManager holds instances of FastStore and LargeStore.
type StoreManager struct {
	fastStore  models.FastStore
	largeStore models.LargeStore
	threshold  int
}

// NewStoreManager initializes StoreManager with specified FastStore, LargeStore, and threshold.
func NewStoreManager(fastStore models.FastStore, largeStore models.LargeStore, threshold int) *StoreManager {
	return &StoreManager{
		fastStore:  fastStore,
		largeStore: largeStore,
		threshold:  threshold,
	}
}

// RecordTick records data in both stores.
func (s *StoreManager) RecordTick(tradingPair string, marketData *types.MarketData) error {
	s.fastStore.RecordTick(tradingPair, marketData)
	s.largeStore.RecordTick(tradingPair, marketData)
	return nil
}

// QueryPriceHistory retrieves data from FastStore or LargeStore based on the period.
func (s *StoreManager) QueryPriceHistory(tradingPair string, period int) []float64 {
	if period <= s.threshold {
		return s.fastStore.QueryPriceHistory(tradingPair, period)
	}
	return s.largeStore.QueryPriceHistory(tradingPair, period)
}
