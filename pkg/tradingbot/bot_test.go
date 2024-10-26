package tradingbot

import (
	"fmt"
	"sync"
	"testing"
	"trading-bot/pkg/types"
)

// MockFastStore and MockLargeStore are simplified as placeholders.
type MockFastStore struct {
	mu           sync.Mutex
	recordedData []types.MarketData
}

func (m *MockFastStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordedData = append(m.recordedData, *marketData)
	return nil
}

func (m *MockFastStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	history := make([]float64, 0, period)
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

type MockLargeStore struct {
	mu           sync.Mutex
	recordedData []types.MarketData
}

func (m *MockLargeStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordedData = append(m.recordedData, *marketData)
	return nil
}

func (m *MockLargeStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	history := make([]float64, 0, period)
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

// MockConnector simulates a trading connector for testing.
type MockConnector struct {
	connected      bool
	StreamDataFunc func(handler func(ctx *types.TickContext)) // Customizable data stream function
}

func (m *MockConnector) Connect() error {
	m.connected = true
	return nil
}

func (m *MockConnector) StreamMarketData(handler func(ctx *types.TickContext)) error {
	if m.connected && m.StreamDataFunc != nil {
		m.StreamDataFunc(handler)
	}
	return nil
}

func (m *MockConnector) ExecuteAction(action types.ActionType, tradingPair string, amount float64) error {
	return nil
}

// MockMiddleware for asserting expected conditions on each tick.
func MockMiddleware(t *testing.T, expectedPrice float64) types.Middleware {
	return func(ctx *types.TickContext) error {
		if ctx.MarketData.Price != expectedPrice {
			t.Errorf("Expected price %v, got %v", expectedPrice, ctx.MarketData.Price)
		}
		fmt.Printf("Middleware processed price: %v\n", ctx.MarketData.Price)
		return nil
	}
}

func TestBot_EndToEndWithMiddleware(t *testing.T) {
	// Setup with mock stores and a threshold
	fastStore := &MockFastStore{}
	largeStore := &MockLargeStore{}
	bot := NewBot(fastStore, largeStore, 2)

	// Create and register a mock connector
	mockConnector := &MockConnector{
		StreamDataFunc: func(handler func(ctx *types.TickContext)) {
			handler(&types.TickContext{
				TradingPair: "BTC/USDT",
				MarketData:  &types.MarketData{Price: 50000.0, Volume: 1.5},
				Actions: types.ActionAPI{
					MarketName: "MockConnector",
					ExecuteAction: func(action types.ActionType, tradingPair string, amount float64) error {
						fmt.Printf("Executing %s for %s with amount %f\n", action, tradingPair, amount)
						return nil
					},
				},
			})
		},
	}
	bot.RegisterConnector("MockConnector", mockConnector)

	// Register mock middleware with an expected price check
	expectedPrice := 50000.0
	bot.fw.UsePair("MockConnector", "BTC/USDT", MockMiddleware(t, expectedPrice))

	// Start the bot
	err := bot.Start()
	if err != nil {
		t.Fatalf("Expected bot to start without error, got: %v", err)
	}

	// Verify connector connection status
	if !mockConnector.connected {
		t.Errorf("Expected connector to be connected after bot start")
	}
}
