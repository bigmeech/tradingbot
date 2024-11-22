## Quick Start Guide: Trading Bot Framework

This guide will walk you through setting up and running the trading bot, including configuring connectors, defining middleware and strategies, and executing trades.

### Prerequisites

- **Golang**: Ensure you have Go installed (version 1.16 or higher recommended).
- **Dependencies**: Install dependencies via `go get` as required by your setup.

### Step 1: Define Connectors

Connectors provide access to different exchanges (e.g., Binance, Kraken) and are defined in the `connectors` folder. Each connector implements `types.Connector`, with methods for connecting to the exchange and executing orders.

Example:
```go
// Create a new Binance connector instance
binanceConnector := connectors.NewBinanceConnector("wss://binance-stream-url", "https://binance-api-url", "your-api-key")

// Register the connector with the bot
bot.RegisterConnector("Binance", binanceConnector)
```

### Step 2: Configure Middleware and Strategies

Middleware and strategies can be used to implement custom logic for processing market data ticks. Middleware is registered per market and trading pair. Here’s an example strategy based on a moving average crossover.

Example strategy:
```go
import "trading-bot/strategies"

// Register a moving average crossover strategy as middleware
bot.RegisterMiddleware("Binance", "BTC/USDT", strategies.MovingAverageCrossoverStrategy())
```

### Step 3: Initialize the Bot

Initialize the bot with fast and large stores, a threshold, and a logger. Stores handle tick data storage and retrieval.

```go
import (
    "bytes"
    "trading-bot/internal/framework"
    "trading-bot/pkg/types"
    "github.com/rs/zerolog"
)

// Initialize stores for tick data
fastStore := framework.NewInMemoryStore()
largeStore := framework.NewInMemoryStore()

// Set up a logger
var logBuffer bytes.Buffer
logger := zerolog.New(&logBuffer).With().Timestamp().Logger()

// Initialize the bot
bot := tradingbot.NewBot(fastStore, largeStore, 10, logger)
bot.EnableDebug()
```

### Step 4: Start Streaming Market Data

With connectors and strategies in place, start streaming data and processing ticks.

```go
// Start the bot
err := bot.Start()
if err != nil {
    log.Fatalf("Failed to start the bot: %v", err)
}
```

### Step 5: Monitor Log Output

The bot logs each tick received, including trading pair, price, and volume. Middleware and strategies log their actions, allowing you to track buy/sell executions.

```shell
cat logBuffer
```

### Example `main.go`

Here’s a simple `main.go` to bring it all together:

```go
package main

import (
    "bytes"
    "log"
    "trading-bot/connectors"
    "trading-bot/internal/framework"
    "trading-bot/pkg/types"
    "trading-bot/strategies"
    "trading-bot/tradingbot"

    "github.com/rs/zerolog"
)

func main() {
    // Initialize stores
    fastStore := framework.NewInMemoryStore()
    largeStore := framework.NewInMemoryStore()

    // Set up logger
    var logBuffer bytes.Buffer
    logger := zerolog.New(&logBuffer).With().Timestamp().Logger()

    // Initialize bot
    bot := tradingbot.NewBot(fastStore, largeStore, 10, logger)
    bot.EnableDebug()

    // Set up and register a Binance connector
    binanceConnector := connectors.NewBinanceConnector("wss://binance-stream-url", "https://binance-api-url", "your-api-key")
    bot.RegisterConnector("Binance", binanceConnector)

    // Register a moving average crossover strategy as middleware
    bot.RegisterMiddleware("Binance", "BTC/USDT", strategies.MovingAverageCrossoverStrategy())

    // Start bot
    if err := bot.Start(); err != nil {
        log.Fatalf("Failed to start bot: %v", err)
    }
}
```

---

Happy trading!