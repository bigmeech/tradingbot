package adapters

import (
	"fmt"
	"log"
	"trading-bot/clients"
	"trading-bot/pkg/types"
)

// MessageParser defines a function type for parsing WebSocket messages into MarketData and trading pairs.
type MessageParser func(message []byte) (*types.MarketData, string, error)

type WebSocketStreamer struct {
	wsClient     *clients.WebSocketClient
	stopCh       chan struct{}
	parseMessage MessageParser
}

// NewWebSocketStreamer initializes a WebSocketStreamer with a WebSocket client and a parser function.
func NewWebSocketStreamer(wsClient *clients.WebSocketClient, parseMessage MessageParser) *WebSocketStreamer {
	return &WebSocketStreamer{
		wsClient:     wsClient,
		stopCh:       make(chan struct{}),
		parseMessage: parseMessage,
	}
}

// StartStreaming begins streaming market data to the provided handler function.
func (ws *WebSocketStreamer) StartStreaming(handler types.MarketDataHandler) error {
	go func() {
		for {
			select {
			case <-ws.stopCh:
				fmt.Println("WebSocketStreamer: Stopping data stream")
				return
			default:
				// Read message from WebSocket
				message, err := ws.wsClient.ReadMessage()
				if err != nil {
					log.Printf("Error reading from WebSocket: %v\n", err)
					continue
				}

				// Use the provided parser to extract MarketData and trading pair
				marketData, tradingPair, parseErr := ws.parseMessage(message)
				if parseErr != nil {
					log.Printf("Error parsing WebSocket message: %v\n", parseErr)
					continue
				}

				// Pass parsed data to the handler
				handler(&types.TickContext{
					TradingPair: tradingPair,
					MarketData:  marketData,
				})
			}
		}
	}()
	return nil
}

// StopStreaming stops the WebSocket data streaming by closing the stop channel.
func (ws *WebSocketStreamer) StopStreaming() error {
	close(ws.stopCh)
	return ws.wsClient.Close()
}
