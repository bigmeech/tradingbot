// connectors/binance_connector.go
package connectors

import (
	"encoding/json"
	"fmt"
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
func NewBinanceConnector(wsClient *clients.WebSocketClient, restClient *clients.RestClient) *BinanceConnector {
	// Initialize WebSocketStreamer with Binance-specific parser
	streamer := adapters.NewWebSocketStreamer(wsClient, binanceMessageParser)

	// Initialize RestExecutor with Binance-specific request formatter
	executor := adapters.NewRestExecutor(restClient, binanceRequestFormatter)

	return &BinanceConnector{
		streamer: streamer,
		executor: executor,
	}
}

// StartStreaming begins streaming Binance market data and processes each tick.
func (bc *BinanceConnector) StartStreaming(handler types.MarketDataHandler) error {
	// Wrap the handler to add Binance-specific actions
	return bc.streamer.StartStreaming(func(ctx *types.TickContext) {
		ctx.Actions = types.ActionAPI{
			MarketName:    "Binance",
			ExecuteAction: bc.ExecuteAction,
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

// ExecuteAction executes an order on Binance using the RestExecutor.
// This function wraps `RestExecutor.ExecuteOrder` and provides Binance-specific formatting.
func (bc *BinanceConnector) ExecuteAction(action types.ActionType, tradingPair string, amount float64) error {
	return bc.executor.ExecuteOrder(action, tradingPair, amount)
}

// binanceRequestFormatter formats requests for the Binance REST API.
// This function prepares the endpoint, HTTP method, and request body to place an order.
func binanceRequestFormatter(action types.ActionType, tradingPair string, amount float64) (string, string, interface{}, error) {
	endpoint := "/api/v3/order"
	method := "POST"
	orderData := map[string]interface{}{
		"symbol":   tradingPair,
		"side":     actionToSide(action),
		"type":     "MARKET", // Assuming market order; modify as needed
		"quantity": amount,
	}

	return endpoint, method, orderData, nil
}

// actionToSide converts action type to the side string used by Binance ("BUY" or "SELL").
func actionToSide(action types.ActionType) string {
	if action == types.ActionBuy {
		return "BUY"
	}
	return "SELL"
}
