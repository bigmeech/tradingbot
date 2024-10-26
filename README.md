Here’s a `README.md` for the trading bot framework, designed as a library. This README provides an overview of the project, instructions for setup, usage examples, and guidance on extending the framework with additional components.

---

# Trading Bot Framework

A modular, extensible trading bot framework in Go for building automated trading strategies across multiple cryptocurrency exchanges. This framework is designed as a library that allows developers to configure exchange connectors, indicators, and strategies and run their custom trading bots.

## Features

- **Multi-Exchange Support**: Easily integrate multiple exchanges using connectors.
- **Indicators as Configurable Components**: Supports indicators like Simple Moving Average (SMA) and Relative Strength Index (RSI).
- **Middleware-Based Strategies**: Use strategies as middleware for flexibility in processing tick data.
- **Data Store for Historical Analysis**: Built-in support for historical data storage and retrieval.
- **Library Format**: Use the framework as a library in other applications, enabling full customization.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/bigmeech/trading-bot.git
   cd trading-bot
   ```

2. **Initialize Go Modules** (if not already done):
   ```bash
   go mod tidy
   ```

3. **Import the Library**: Import the `tradingbot` package in your project.

## Usage

### 1. Initialize the Trading Bot

The following code demonstrates how to set up and run the bot in a new project’s `main.go` file:

```go
package main

import (
    "trading-bot/internal/connectors"
    "trading-bot/internal/indicators"
    "trading-bot/internal/strategies"
    "trading-bot/pkg/models"
    "trading-bot/pkg/tradingbot"
)

func main() {
    // Initialize data store
    store := NewInMemoryStore()
    bot := tradingbot.NewBot(store)

    // Configure and register a Binance connector
    binance := connectors.NewBinanceConnector("api-key", "api-secret", "wss://binance.ws", "https://binance.rest")
    bot.RegisterConnector("Binance", binance)

    // Register indicators for trading pair BTC/USDT on Binance
    bot.RegisterIndicator("Binance", "BTC/USDT", &indicators.SMA{Period: 50})
    bot.RegisterIndicator("Binance", "BTC/USDT", &indicators.SMA{Period: 200})
    bot.RegisterIndicator("Binance", "BTC/USDT", &indicators.RSI{Period: 14})

    // Add strategies (middleware) for the trading pair
    bot.AddStrategy("Binance", "BTC/USDT", strategies.MovingAverageCrossoverStrategy())
    bot.AddStrategy("Binance", "BTC/USDT", strategies.RSIOverboughtOversoldStrategy())

    // Start the bot
    if err := bot.Start(); err != nil {
        panic(err)
    }
}
```

### 2. Configuration

The framework relies on three main configuration areas:
- **Connectors**: Define exchange API keys, URLs, and WebSocket URLs.
- **Indicators**: Attach indicators to specific trading pairs for strategy use.
- **Strategies**: Register strategies as middleware that responds to tick data and indicator values.

### 3. Example Strategies

1. **Moving Average Crossover Strategy**: Buys when a short SMA crosses above a long SMA and sells when it crosses below.

   ```go
   package strategies

   import (
       "trading-bot/pkg/models"
   )

   func MovingAverageCrossoverStrategy() models.Middleware {
       return func(ctx *models.TickContext) error {
           shortSMA := ctx.Indicators["SMA_50"]
           longSMA := ctx.Indicators["SMA_200"]

           if shortSMA > longSMA {
               ctx.Actions.Buy(1.0)
           } else if shortSMA < longSMA {
               ctx.Actions.Sell(1.0)
           }
           return nil
       }
   }
   ```

2. **RSI Overbought/Oversold Strategy**: Buys when RSI is below 30 (oversold) and sells when RSI is above 70 (overbought).

   ```go
   package strategies

   import (
       "trading-bot/pkg/models"
   )

   func RSIOverboughtOversoldStrategy() models.Middleware {
       return func(ctx *models.TickContext) error {
           rsi := ctx.Indicators["RSI_14"]

           if rsi < 30 {
               ctx.Actions.Buy(1.0) // Buy when oversold
           } else if rsi > 70 {
               ctx.Actions.Sell(1.0) // Sell when overbought
           }
           return nil
       }
   }
   ```

### 4. Adding Custom Indicators

To add custom indicators, implement the `models.Indicator` interface and register it in the bot.

```go
package indicators

type CustomIndicator struct {
    Period int
}

func (c *CustomIndicator) Calculate(data []float64) float64 {
    // Custom calculation logic here
    return result
}

func (c *CustomIndicator) Name() string {
    return "CustomIndicator"
}
```

Then, register it:

```go
bot.RegisterIndicator("Binance", "BTC/USDT", &indicators.CustomIndicator{Period: 14})
```

## Extending the Framework

### Adding a New Connector

1. Implement the `models.Connector` interface in `internal/connectors/`.
2. Register the connector using `bot.RegisterConnector()`.

### Adding New Strategies

1. Create a new strategy in `internal/strategies/` as a `models.Middleware`.
2. Register the strategy using `bot.AddStrategy()`.

## Contributing

If you'd like to contribute:
1. Fork the repository.
2. Create a new branch for your feature.
3. Open a pull request with a detailed description of your changes.

## License

This project is licensed under the MIT License.

---

This README provides essential information on setting up and configuring the library, allowing users to quickly integrate, extend, and run their own trading bot configurations.