package connectors

import (
	"encoding/json"
	"fmt"
	"time"
	"trading-bot/adapters"
	"trading-bot/clients"
	"trading-bot/pkg/types"
)

// KrakenConnector encapsulates Kraken-specific streaming and order execution functionality.
type KrakenConnector struct {
	streamer *adapters.WebSocketStreamer
	executor *adapters.RestExecutor
}

// NewKrakenConnector initializes a KrakenConnector with Kraken-specific WebSocket and REST clients.
func NewKrakenConnector(wsURL, restURL, apiKey string) *KrakenConnector {
	// Set up a WebSocket client with Kraken constraints
	wsClient := clients.NewWebSocketClient(
		wsURL,
		24*time.Hour,   // Connection lifetime
		3*time.Minute,  // Ping interval
		10*time.Minute, // Pong timeout
		10,             // Rate limit: 10 messages per second
		200,            // Stream limit: 200 streams per connection
	)

	// Initialize WebSocketStreamer with Kraken-specific parser and stream limit
	streamer := adapters.NewWebSocketStreamer(wsClient, krakenMessageParser, 200)

	// Initialize RestExecutor with Kraken-specific request formatter and REST client
	restClient := clients.NewRestClient(restURL, apiKey)
	executor := adapters.NewRestExecutor(restClient, krakenRequestFormatter)

	return &KrakenConnector{
		streamer: streamer,
		executor: executor,
	}
}

// StartStreaming begins streaming Kraken market data and processes each tick.
func (kc *KrakenConnector) StartStreaming(handler types.MarketDataHandler) error {
	// Wrap the handler to provide Kraken-specific order execution
	return kc.streamer.StartStreaming(func(ctx *types.TickContext) {
		ctx.ExecuteOrder = func(orderType types.OrderType, side types.OrderSide, amount, price float64) error {
			return kc.ExecuteOrder(orderType, side, ctx.TradingPair, amount, price)
		}
		handler(ctx)
	})
}

// StopStreaming stops the Kraken data streaming.
func (kc *KrakenConnector) StopStreaming() error {
	return kc.streamer.StopStreaming()
}

// krakenMessageParser parses Kraken WebSocket messages into MarketData.
func krakenMessageParser(message []byte) (*types.MarketData, string, error) {
	var parsedData map[string]interface{}
	if err := json.Unmarshal(message, &parsedData); err != nil {
		return nil, "", err
	}

	// Assuming Kraken’s WebSocket data includes price, volume, and symbol
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

// ExecuteOrder places an order on Kraken with the specified type and side.
func (kc *KrakenConnector) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
	return kc.executor.ExecuteOrder(orderType, side, tradingPair, amount, price)
}

// krakenRequestFormatter formats requests for the Kraken REST API.
// This function prepares the endpoint, HTTP method, and request body to place an order.
func krakenRequestFormatter(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) (string, string, interface{}, error) {
	endpoint := "/0/private/AddOrder"
	method := "POST"
	orderData := map[string]interface{}{
		"pair":      tradingPair,
		"type":      string(side),      // "buy" or "sell"
		"ordertype": string(orderType), // Order type, e.g., "market", "limit"
		"volume":    amount,
	}

	// Include price for order types that require it (e.g., LIMIT)
	if orderType == types.OrderTypeLimit || orderType == types.OrderTypeStopLossLimit || orderType == types.OrderTypeTakeProfitLimit {
		orderData["price"] = price
	}

	return endpoint, method, orderData, nil
}
