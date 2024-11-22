# Connector Documentation

The `Connector` component in the trading bot framework is responsible for interacting with exchanges (e.g., Binance, Kraken) to stream market data and execute trade orders. Each connector manages the connection with a specific exchange, fetches real-time market data, and provides order execution capabilities.

## Purpose

Connectors:
- Establish a link to exchanges for streaming market data in real time.
- Format and send orders (buy/sell) to the exchange based on strategy logic.
- Serve as the interface between the trading bot and each specific exchange.

## Connector Interface

The `Connector` interface defines the standard methods each connector must implement, allowing the bot framework to interact with exchanges in a consistent manner.

```go
type Connector interface {
    // Connect establishes a connection to the exchange, preparing it for streaming data and executing orders.
    Connect() error
    
    // StreamMarketData starts streaming market data, with each tick data passed to the provided handler.
    // The handler receives a TickContext that includes market details, trading pair, and data.
    StreamMarketData(handler func(ctx *TickContext)) error
    
    // ExecuteOrder places an order on the exchange using the specified parameters.
    // The function parameters include the order type, side (buy/sell), trading pair, amount, and price.
    ExecuteOrder(orderType OrderType, side OrderSide, tradingPair string, amount, price float64) error
}
```

### Connector Methods

1. **`Connect()`**:
    - Establishes a connection to the exchange, setting up WebSocket and/or REST clients as necessary.
    - Ensures that the connector is ready to stream data and execute trades.
    - Returns an error if the connection fails.

2. **`StreamMarketData(handler func(ctx *TickContext))`**:
    - Starts streaming market data (price, volume, etc.) for the configured trading pairs.
    - Each tick of market data is passed to the `handler`, which receives a `TickContext` containing the tick’s details.
    - Returns an error if data streaming fails.

3. **`ExecuteOrder(orderType OrderType, side OrderSide, tradingPair string, amount, price float64)`**:
    - Places a trade order (buy/sell) on the exchange.
    - Parameters:
        - `orderType`: Specifies the type of order (e.g., `LIMIT`, `MARKET`).
        - `side`: Defines the trade direction (`BUY` or `SELL`).
        - `tradingPair`: The trading pair for the order (e.g., "BTC/USDT").
        - `amount`: The amount of the asset to trade.
        - `price`: The price for the order (used for limit orders).
    - Returns an error if the order fails to execute.

---

## Example Connector Implementations

### 1. BinanceConnector

The `BinanceConnector` handles Binance-specific streaming and order execution. It uses Binance’s WebSocket API to stream market data and the REST API for executing orders.

#### `binance_connector.go`

```go
package connectors

import (
    "trading-bot/adapters"
    "trading-bot/clients"
    "trading-bot/pkg/types"
    "time"
)

// BinanceConnector encapsulates Binance-specific streaming and order execution functionality.
type BinanceConnector struct {
    streamer *adapters.WebSocketStreamer
    executor *adapters.RestExecutor
}

// NewBinanceConnector initializes a BinanceConnector with WebSocket and REST clients.
func NewBinanceConnector(wsURL, restURL, apiKey string) *BinanceConnector {
    wsClient := clients.NewWebSocketClient(wsURL, 24*time.Hour, 3*time.Minute, 10*time.Minute, 10, 200)
    streamer := adapters.NewWebSocketStreamer(wsClient, binanceMessageParser, 200)
    restClient := clients.NewRestClient(restURL, apiKey)
    executor := adapters.NewRestExecutor(restClient, binanceRequestFormatter)

    return &BinanceConnector{
        streamer: streamer,
        executor: executor,
    }
}

// Connect establishes the connection to Binance.
func (bc *BinanceConnector) Connect() error {
    return bc.streamer.Connect()
}

// StreamMarketData starts streaming Binance market data and passes it to the provided handler.
func (bc *BinanceConnector) StreamMarketData(handler func(ctx *types.TickContext)) error {
    return bc.streamer.StartStreaming(func(ctx *types.TickContext) {
        ctx.MarketName = "Binance"
        ctx.ExecuteOrder = bc.ExecuteOrder
        handler(ctx)
    })
}

// ExecuteOrder places an order on Binance.
func (bc *BinanceConnector) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
    return bc.executor.ExecuteOrder(orderType, side, tradingPair, amount, price)
}
```

### 2. KrakenConnector

The `KrakenConnector` is similar to `BinanceConnector`, but it implements Kraken-specific WebSocket and REST clients.

#### `kraken_connector.go`

```go
package connectors

import (
    "trading-bot/adapters"
    "trading-bot/clients"
    "trading-bot/pkg/types"
    "time"
)

// KrakenConnector encapsulates Kraken-specific streaming and order execution functionality.
type KrakenConnector struct {
    streamer *adapters.WebSocketStreamer
    executor *adapters.RestExecutor
}

// NewKrakenConnector initializes a KrakenConnector with WebSocket and REST clients.
func NewKrakenConnector(wsURL, restURL, apiKey string) *KrakenConnector {
    wsClient := clients.NewWebSocketClient(wsURL, 24*time.Hour, 3*time.Minute, 10*time.Minute, 10, 200)
    streamer := adapters.NewWebSocketStreamer(wsClient, krakenMessageParser, 200)
    restClient := clients.NewRestClient(restURL, apiKey)
    executor := adapters.NewRestExecutor(restClient, krakenRequestFormatter)

    return &KrakenConnector{
        streamer: streamer,
        executor: executor,
    }
}

// Connect establishes the connection to Kraken.
func (kc *KrakenConnector) Connect() error {
    return kc.streamer.Connect()
}

// StreamMarketData starts streaming Kraken market data and passes it to the provided handler.
func (kc *KrakenConnector) StreamMarketData(handler func(ctx *types.TickContext)) error {
    return kc.streamer.StartStreaming(func(ctx *types.TickContext) {
        ctx.MarketName = "Kraken"
        ctx.ExecuteOrder = kc.ExecuteOrder
        handler(ctx)
    })
}

// ExecuteOrder places an order on Kraken.
func (kc *KrakenConnector) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
    return kc.executor.ExecuteOrder(orderType, side, tradingPair, amount, price)
}
```

---

## Using Connectors with the Bot

Once initialized, connectors are registered with the bot and automatically handle data streaming and order execution based on market conditions and strategies.

### Example Usage in `main.go`

```go
package main

import (
    "trading-bot/internal/connectors"
    "trading-bot/pkg/tradingbot"
    "github.com/rs/zerolog"
    "os"
)

func main() {
    // Initialize logger
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Initialize bot
    bot := tradingbot.NewBot(logger)

    // Set up and register Binance connector
    binanceConnector := connectors.NewBinanceConnector("wss://binance-stream-url", "https://binance-api-url", "your-api-key")
    bot.RegisterConnector("Binance", binanceConnector)

    // Set up and register Kraken connector
    krakenConnector := connectors.NewKrakenConnector("wss://kraken-stream-url", "https://kraken-api-url", "your-api-key")
    bot.RegisterConnector("Kraken", krakenConnector)

    // Start the bot
    if err := bot.Start(); err != nil {
        logger.Fatal().Err(err).Msg("Failed to start bot")
    }
}
```

---

## Summary

- **Purpose**: Connectors establish connections to exchanges, stream market data, and execute orders.
- **Interface**: The `Connector` interface provides a standard structure for connecting, streaming data, and executing orders.
- **Implementation**: Each exchange requires a custom connector (e.g., `BinanceConnector`, `KrakenConnector`) to handle exchange-specific data formats and order processes.
- **Integration with Bot**: Connectors are registered with the bot, which uses them to receive real-time data and execute trades based on strategy outputs.