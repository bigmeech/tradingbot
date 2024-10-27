// adapters/websocket_streamer.go
package adapters

import (
	"fmt"
	"trading-bot/clients"
	"trading-bot/pkg/types"
)

// MessageParser is a function type for parsing WebSocket messages into MarketData and trading pairs.
type MessageParser func(message []byte) (*types.MarketData, string, error)

// WebSocketStreamer manages the streaming of market data through a WebSocket connection.
type WebSocketStreamer struct {
	client        *clients.WebSocketClient
	messageParser MessageParser
	activeStreams int
	maxStreams    int
}

// NewWebSocketStreamer initializes a WebSocketStreamer with a WebSocket client, a message parser, and a max stream limit.
func NewWebSocketStreamer(client *clients.WebSocketClient, messageParser MessageParser, maxStreams int) *WebSocketStreamer {
	return &WebSocketStreamer{
		client:        client,
		messageParser: messageParser,
		maxStreams:    maxStreams,
		activeStreams: 0,
	}
}

// StartStreaming begins streaming data to the handler function, using the provided message parser.
func (ws *WebSocketStreamer) StartStreaming(handler types.MarketDataHandler) error {
	if ws.activeStreams >= ws.maxStreams {
		return fmt.Errorf("maximum stream limit reached")
	}

	ws.activeStreams++
	return ws.client.StartStreaming(func(data []byte) {
		// Parse the message using the provided message parser
		marketData, tradingPair, err := ws.messageParser(data)
		if err != nil {
			// Log or handle parsing errors if necessary
			return
		}

		// Call the handler with the parsed market data and trading pair
		handler(&types.TickContext{
			TradingPair: tradingPair,
			MarketData:  marketData,
		})
	})
}

// StopStreaming decreases the active stream count and stops the client if no streams are active.
func (ws *WebSocketStreamer) StopStreaming() error {
	ws.activeStreams--
	if ws.activeStreams == 0 {
		return ws.client.StopStreaming()
	}
	return nil
}
