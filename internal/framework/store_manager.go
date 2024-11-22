package framework

import (
	"github.com/bigmeech/tradingbot/internal/store" // For CircularBuffer
	"github.com/bigmeech/tradingbot/pkg/types"
	"sync"
)

// StoreManager manages both fast and persistent storage for market data.
type StoreManager struct {
	fastStore  map[string]*store.CircularBuffer // Fast in-memory buffer for recent data
	largeStore *store.MongoDBStore              // Persistent store for historical data
	bufferSize int                              // Configurable buffer size for each trading pair
	threshold  int                              // Threshold period for fastStore vs largeStore
	storeLock  sync.Mutex                       // For thread-safe access
}

// NewStoreManager initializes a StoreManager with a persistent store and buffer configuration.
func NewStoreManager(largeStore *store.MongoDBStore, bufferSize int, threshold int) *StoreManager {
	return &StoreManager{
		fastStore:  make(map[string]*store.CircularBuffer),
		largeStore: largeStore,
		bufferSize: bufferSize,
		threshold:  threshold,
	}
}

// RecordTick records a new market data point, adding it to both fastStore and largeStore.
func (s *StoreManager) RecordTick(market, tradingPair string, data *types.MarketData) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	// Ensure fastStore is initialized
	if s.fastStore == nil {
		s.fastStore = make(map[string]*store.CircularBuffer)
	}

	// Create a unique key for the market/trading pair combination
	key := market + ":" + tradingPair

	// Initialize a new CircularBuffer if one doesn't exist for the key
	if _, exists := s.fastStore[key]; !exists {
		s.fastStore[key] = store.NewCircularBuffer(s.bufferSize)
	}

	// Record the tick in the fast circular buffer
	s.fastStore[key].Add(*data)

	// Also store the tick in the largeStore for long-term storage
	return s.largeStore.RecordTick(tradingPair, data)
}

// QueryPriceHistory fetches data based on the period from either fastStore or largeStore.
func (s *StoreManager) QueryPriceHistory(market, tradingPair string, period int) []float64 {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	// Ensure fastStore is initialized
	if s.fastStore == nil {
		s.fastStore = make(map[string]*store.CircularBuffer)
	}

	// Define the key and check if recent data can be fetched from fastStore
	key := market + ":" + tradingPair
	if buffer, exists := s.fastStore[key]; exists {
		recentData := buffer.GetData(period)
		priceHistory := make([]float64, len(recentData))
		for i, entry := range recentData {
			priceHistory[i] = entry.Price
		}
		return priceHistory
	} else {
		// Handle the case where the key does not exist
		return []float64{}
	}

	// For periods beyond threshold, fall back to the largeStore
	return s.largeStore.QueryPriceHistory(tradingPair, period)
}
