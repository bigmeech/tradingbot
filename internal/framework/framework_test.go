package framework

import (
	"fmt"
	"testing"
	"trading-bot/pkg/types"
)

// MockConnector is a mock implementation of the Connector interface for testing.
type MockConnector struct {
	connected bool
}

func (m *MockConnector) Connect() error {
	m.connected = true
	return nil
}

func (m *MockConnector) StreamMarketData(handler func(*types.TickContext)) error {
	if m.connected {
		marketData := &types.MarketData{Price: 100.0, Volume: 1.5}
		handler(&types.TickContext{
			TradingPair: "BTC/USDT",
			MarketData:  marketData,
			Actions: types.ActionAPI{
				MarketName: "MockConnector",
				ExecuteAction: func(action types.ActionType, tradingPair string, amount float64) error {
					fmt.Printf("Executing %s for %s with amount %f\n", action, tradingPair, amount)
					return nil
				},
			},
		})
	}
	return nil
}

func (m *MockConnector) ExecuteAction(action types.ActionType, tradingPair string, amount float64) error {
	fmt.Printf("Executing %s for %s with amount %f\n", action, tradingPair, amount)
	return nil
}

// MockIndicator is a mock implementation of the Indicator interface for testing.
type MockIndicator struct {
	name string
}

func (mi *MockIndicator) Name() string {
	return mi.name
}

func (mi *MockIndicator) Calculate(data []float64) float64 {
	return 50.0 // Return a fixed value for simplicity
}

// MockMiddleware is a mock implementation of the Middleware function type for testing.
func MockMiddleware(ctx *types.TickContext) error {
	fmt.Println("Executing middleware for:", ctx.TradingPair)
	return nil
}

func TestFramework_RegisterConnectorAndStart(t *testing.T) {
	store := NewInMemoryFastStore(100) // Using InMemoryFastStore for simplicity
	framework := NewFramework(store)

	mockConnector := &MockConnector{}
	framework.RegisterConnector("MockConnector", mockConnector)

	if len(framework.Connectors()) != 1 {
		t.Fatalf("Expected 1 connector, got %v", len(framework.Connectors()))
	}

	framework.Start()
	if !mockConnector.connected {
		t.Errorf("Expected connector to be connected")
	}
}

func TestFramework_RegisterIndicator(t *testing.T) {
	store := NewInMemoryFastStore(100)
	framework := NewFramework(store)

	indicator := &MockIndicator{name: "MockIndicator"}
	framework.RegisterIndicator("MockConnector", "BTC/USDT", indicator)

	if len(framework.indicators["MockConnector"]["BTC/USDT"]) != 1 {
		t.Fatalf("Expected 1 indicator for BTC/USDT, got %v", len(framework.indicators["MockConnector"]["BTC/USDT"]))
	}
}

func TestFramework_RegisterMiddleware(t *testing.T) {
	store := NewInMemoryFastStore(100)
	framework := NewFramework(store)

	framework.UsePair("MockConnector", "BTC/USDT", MockMiddleware)

	if len(framework.middleware["MockConnector"]["BTC/USDT"]) != 1 {
		t.Fatalf("Expected 1 middleware for BTC/USDT, got %v", len(framework.middleware["MockConnector"]["BTC/USDT"]))
	}
}

func TestFramework_ProcessTick(t *testing.T) {
	store := NewInMemoryFastStore(100)
	framework := NewFramework(store)

	// Register a mock connector, indicator, and middleware
	mockConnector := &MockConnector{}
	framework.RegisterConnector("MockConnector", mockConnector)

	indicator := &MockIndicator{name: "MockIndicator"}
	framework.RegisterIndicator("MockConnector", "BTC/USDT", indicator)

	framework.UsePair("MockConnector", "BTC/USDT", MockMiddleware)

	// Simulate a tick being processed
	tickContext := &types.TickContext{
		TradingPair: "BTC/USDT",
		MarketData:  &types.MarketData{Price: 100.0, Volume: 1.5},
		Actions: types.ActionAPI{
			MarketName:    "MockConnector",
			ExecuteAction: mockConnector.ExecuteAction,
		},
	}

	framework.processTick(tickContext)
	indicators := tickContext.Indicators

	// Check that the indicator value is set correctly in TickContext
	if len(indicators) != 1 || indicators["MockIndicator"] != 50.0 {
		t.Fatalf("Expected indicator value 50.0 for MockIndicator, got %v", indicators["MockIndicator"])
	}
}
