package framework

import (
	"github.com/bigmeech/tradingbot/pkg/types"
	"sync"
	"testing"
	"time"
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

// MockConnector simulates a trading connector for testing.
type MockConnector struct {
	connected    bool
	streamDataFn func(handler func(ctx *types.TickContext))
}

// Connect simulates a successful connection for the mock connector.
func (m *MockConnector) Connect() error {
	m.connected = true
	return nil
}

// StreamMarketData simulates streaming market data for the mock connector.
func (m *MockConnector) StreamMarketData(handler func(ctx *types.TickContext)) error {
	if m.connected && m.streamDataFn != nil {
		handler(&types.TickContext{
			MarketName:  "MockConnector",
			TradingPair: "BTC/USDT",
			MarketData:  &types.MarketData{Price: 50000.0, Volume: 1.5},
			ExecuteOrder: func(orderType types.OrderType, side types.OrderSide, amount, price float64) error {
				// Simulate order execution for testing
				return nil
			},
		})
	}
	return nil
}

// ExecuteOrder simulates executing an order for the mock connector.
func (m *MockConnector) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
	// Simulate order execution logic for testing
	return nil
}

func TestFramework_RegisterConnectorAndStreamTicks(t *testing.T) {
	store := NewMockStore() // Use the mock store directly
	framework := NewFramework(store)

	// Set up a channel to capture processed ticks
	processedTicks := make(chan *types.TickContext, 1)

	// Mock tick processing function
	processTickFunc := func(ctx *types.TickContext) {
		t.Log("processTickFunc called") // Log for debugging
		processedTicks <- ctx           // Send tick context to channel
	}

	// Create and register a mock connector with simulated data streaming
	mockConnector := &MockConnector{
		streamDataFn: func(handler func(ctx *types.TickContext)) {
			t.Log("streamDataFn called") // Log for debugging
			handler(&types.TickContext{
				MarketName:  "MockConnector",
				TradingPair: "BTC/USDT",
				MarketData:  &types.MarketData{Price: 50000.0, Volume: 1.5},
				ExecuteOrder: func(orderType types.OrderType, side types.OrderSide, amount, price float64) error {
					// Simulate order execution for testing
					return nil
				},
			})
		},
	}

	framework.RegisterConnector("MockConnector", mockConnector)
	if len(framework.Connectors()) != 1 {
		t.Fatalf("Expected 1 connector, got %v", len(framework.Connectors()))
	}

	// Start the framework with the mock tick processing function
	go framework.Start(processTickFunc)

	// Wait for the tick to be processed and received in the channel with a timeout
	select {
	case tick := <-processedTicks:
		// Validate the received tick
		if tick.MarketName != "MockConnector" {
			t.Errorf("Expected MarketName MockConnector, got %v", tick.MarketName)
		}
		if tick.TradingPair != "BTC/USDT" {
			t.Errorf("Expected trading pair BTC/USDT, got %v", tick.TradingPair)
		}
		if tick.MarketData.Price != 50000.0 {
			t.Errorf("Expected price 50000.0, got %v", tick.MarketData.Price)
		}
		if tick.MarketData.Volume != 1.5 {
			t.Errorf("Expected volume 1.5, got %v", tick.MarketData.Volume)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Expected tick to be processed but received none within the timeout period")
	}
}
