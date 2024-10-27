package tradingbot

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	"trading-bot/pkg/types"
	"trading-bot/testutils"

	"github.com/rs/zerolog"
)

type MockConnector struct {
	connected    bool
	streamDataFn func(handler func(ctx *types.TickContext))
	stopCh       chan struct{} // Channel to stop streaming data
}

func NewMockConnector() *MockConnector {
	return &MockConnector{
		stopCh: make(chan struct{}),
	}
}

func (m *MockConnector) Connect() error {
	m.connected = true
	fmt.Println("MockConnector: Connected")
	return nil
}

func (m *MockConnector) StreamMarketData(handler func(ctx *types.TickContext)) error {
	if !m.connected || m.streamDataFn == nil {
		return nil
	}

	// Stream data continuously with a delay to simulate real-time data
	go func() {
		for {
			select {
			case <-m.stopCh:
				fmt.Println("MockConnector: Stopping data stream")
				return
			default:
				// Simulate a tick with a random price and volume for each loop iteration
				m.streamDataFn(handler)
				time.Sleep(200 * time.Millisecond) // Short delay between ticks
			}
		}
	}()

	return nil
}

// Stop streaming data
func (m *MockConnector) Stop() {
	close(m.stopCh)
}

func (m *MockConnector) ExecuteAction(action types.ActionType, tradingPair string, amount float64) error {
	return nil
}

func TestBot_RegisterConnectorAndStart(t *testing.T) {
	// Initialize mock stores for testing
	fastStore := testutils.NewMockStore()
	largeStore := testutils.NewMockStore()

	// Set up log buffer and configure logger
	var logBuffer bytes.Buffer
	logger := zerolog.New(&logBuffer).With().Timestamp().Logger()

	// Initialize Bot with custom logger
	bot := NewBot(fastStore, largeStore, 10, logger)
	bot.EnableDebug()

	// Create and register the mock connector
	mockConnector := NewMockConnector()
	mockConnector.streamDataFn = func(handler func(ctx *types.TickContext)) {
		handler(&types.TickContext{
			TradingPair: "BTC/USDT",
			MarketData:  &types.MarketData{Price: 50000.0, Volume: 1.5},
			Actions: types.ActionAPI{
				MarketName:    "MockConnector",
				ExecuteAction: func(action types.ActionType, tradingPair string, amount float64) error { return nil },
			},
		})
	}

	bot.RegisterConnector("MockConnector", mockConnector)
	if len(bot.fw.Connectors()) != 1 {
		t.Fatalf("Expected 1 connector, got %v", len(bot.fw.Connectors()))
	}

	// Start the bot asynchronously
	go func() {
		err := bot.Start()
		if err != nil {
			t.Fatalf("Expected bot to start without error, got: %v", err)
		}
	}()

	// Allow time for several ticks to be processed
	time.Sleep(1 * time.Second)

	// Stop the mock connector's data stream
	mockConnector.Stop()

	// Check the log output for the tick details
	logOutput := logBuffer.String()
	if logOutput == "" {
		t.Fatal("Expected log output but found none")
	}
	if !bytes.Contains(logBuffer.Bytes(), []byte("Received tick")) {
		t.Errorf("Expected log to contain 'Received tick', got %v", logOutput)
	}
	if !bytes.Contains(logBuffer.Bytes(), []byte("BTC/USDT")) {
		t.Errorf("Expected log to contain 'BTC/USDT', got %v", logOutput)
	}
	if !bytes.Contains(logBuffer.Bytes(), []byte("Price\":50000")) {
		t.Errorf("Expected log to contain 'Price\":50000', got %v", logOutput)
	}
}
