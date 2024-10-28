package connectors

import (
	"encoding/json"
	"fmt"
	"time"
	"trading-bot/adapters"
	"trading-bot/clients"
	"trading-bot/pkg/types"
)

// BinanceConnector encapsulates Binance-specific streaming and order execution functionality.
type BinanceConnector struct {
	streamer *adapters.WebSocketStreamer
	executor *adapters.RestExecutor
}

// NewBinanceConnector initializes a BinanceConnector with Binance-specific WebSocket and REST clients.
func NewBinanceConnector(wsURL, restURL, apiKey string) *BinanceConnector {
	// Set up a WebSocket client with Binance constraints
	wsClient := clients.NewWebSocketClient(
		wsURL,
		24*time.Hour,   // Connection lifetime
		3*time.Minute,  // Ping interval
		10*time.Minute, // Pong timeout
		10,             // Rate limit: 10 messages per second
		200,            // Stream limit: 200 streams per connection
	)

	// Initialize WebSocketStreamer with Binance-specific parser and stream limit
	streamer := adapters.NewWebSocketStreamer(wsClient, binanceMessageParser, 200)

	// Initialize RestExecutor with Binance-specific request formatter and REST client
	restClient := clients.NewRestClient(restURL, apiKey)
	executor := adapters.NewRestExecutor(restClient, binanceRequestFormatter)

	return &BinanceConnector{
		streamer: streamer,
		executor: executor,
	}
}

// StartStreaming begins streaming Binance market data and processes each tick.
func (bc *BinanceConnector) StartStreaming(handler types.MarketDataHandler) error {
	return bc.streamer.StartStreaming(func(ctx *types.TickContext) {
		// Wrap bc.ExecuteOrder with a compatible function signature for ctx.ExecuteOrder
		ctx.ExecuteOrder = func(orderType types.OrderType, side types.OrderSide, amount, price float64) error {
			return bc.ExecuteOrder(orderType, side, ctx.TradingPair, amount, price)
		}
		handler(ctx)
	})
}

// StopStreaming stops the Binance data streaming.
func (bc *BinanceConnector) StopStreaming() error {
	return bc.streamer.StopStreaming()
}

// binanceMessageParser parses Binance WebSocket messages into MarketData.
func binanceMessageParser(message []byte) (*types.MarketData, string, error) {
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

// ExecuteOrder places an order on Binance with the specified type and side.
func (bc *BinanceConnector) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
	return bc.executor.ExecuteOrder(orderType, side, tradingPair, amount, price)
}

// binanceRequestFormatter formats requests for the Binance REST API.
func binanceRequestFormatter(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) (string, string, interface{}, error) {
	endpoint := "/api/v3/order"
	method := "POST"
	orderData := map[string]interface{}{
		"symbol":   tradingPair,
		"side":     string(side),      // Use the side directly as "BUY" or "SELL"
		"type":     string(orderType), // Order type (e.g., "MARKET", "LIMIT")
		"quantity": amount,
	}

	// Include price for order types that require it (e.g., LIMIT)
	if orderType == types.OrderTypeLimit || orderType == types.OrderTypeStopLossLimit || orderType == types.OrderTypeTakeProfitLimit {
		orderData["price"] = price
	}

	return endpoint, method, orderData, nil
}
