package connectors

import (
	"encoding/json"
	"fmt"
	"github.com/bigmeech/tradingbot/adapters"
	"github.com/bigmeech/tradingbot/clients"
	"github.com/bigmeech/tradingbot/pkg/types"
	"log"
	"time"
)

// LocalConnector encapsulates local WebSocket streaming and REST order execution.
type LocalConnector struct {
	streamer *adapters.WebSocketStreamer
	executor *adapters.RestExecutor
}

// NewLocalConnector initializes a LocalConnector with WebSocket and REST clients.
func NewLocalConnector(wsURL, restURL, apiKey string) *LocalConnector {
	// Set up WebSocket client for local data streaming
	wsClient := clients.NewWebSocketClient(
		wsURL,
		24*time.Hour,   // Connection lifetime
		3*time.Minute,  // Ping interval
		10*time.Minute, // Pong timeout
		10,             // Rate limit: 10 messages per second
		200,            // Stream limit: 200 streams per connection
	)

	// Establish the WebSocket connection immediately
	if err := wsClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// Initialize WebSocketStreamer with a message parser
	streamer := adapters.NewWebSocketStreamer(wsClient, localMessageParser, 200)

	// Set up RestExecutor with a REST client for local order execution
	restClient := clients.NewRestClient(restURL, apiKey)
	executor := adapters.NewRestExecutor(restClient, localRequestFormatter)

	return &LocalConnector{
		streamer: streamer,
		executor: executor,
	}
}

// StreamMarketData begins streaming local market data and processes each tick.
func (lc *LocalConnector) StreamMarketData(handler func(ctx *types.TickContext)) error {
	return lc.streamer.StartStreaming(func(ctx *types.TickContext) {
		// Wrap ExecuteOrder function in TickContext
		ctx.ExecuteOrder = func(orderType types.OrderType, side types.OrderSide, amount, price float64) error {
			return lc.ExecuteOrder(orderType, side, ctx.TradingPair, amount, price)
		}
		handler(ctx)
	})
}

// StopStreaming stops the local data streaming.
func (lc *LocalConnector) StopStreaming() error {
	return lc.streamer.StopStreaming()
}

// ExecuteOrder places an order locally with the specified type and side.
func (lc *LocalConnector) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
	return lc.executor.ExecuteOrder(orderType, side, tradingPair, amount, price)
}

// GetIdentifier returns the WebSocket URL as the unique identifier for LocalConnector.
func (lc *LocalConnector) GetIdentifier() string {
	return lc.streamer.Client.GetConnectionUrl()
}

// localMessageParser parses WebSocket messages from the local server into MarketData.
func localMessageParser(message []byte) (*types.MarketData, string, error) {
	var parsedData map[string]interface{}
	if err := json.Unmarshal(message, &parsedData); err != nil {
		return nil, "", err
	}

	price, _ := parsedData["price"].(float64)
	volume, _ := parsedData["volume"].(float64)
	tradingPair, _ := parsedData["symbol"].(string)

	if price == 0 || volume == 0 {
		return nil, "", fmt.Errorf("invalid data in WebSocket message")
	}

	return &types.MarketData{
		Price:  price,
		Volume: volume,
	}, tradingPair, nil
}

// localRequestFormatter formats requests for the local REST API.
func localRequestFormatter(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) (string, string, interface{}, error) {
	endpoint := "/api/v1/order"
	method := "POST"
	orderData := map[string]interface{}{
		"symbol":   tradingPair,
		"side":     string(side),      // Use "BUY" or "SELL"
		"type":     string(orderType), // Order type (e.g., "MARKET", "LIMIT")
		"quantity": amount,
	}

	// Include price for order types that require it (e.g., LIMIT)
	if orderType == types.OrderTypeLimit || orderType == types.OrderTypeStopLossLimit || orderType == types.OrderTypeTakeProfitLimit {
		orderData["price"] = price
	}

	return endpoint, method, orderData, nil
}
