# Strategies Documentation

The `Strategies` component in the trading bot framework defines the logic for making trading decisions based on incoming market data (ticks). Strategies analyze real-time data and technical indicators to identify trading signals, such as buy or sell opportunities, which are then executed via the bot's connectors.

## Purpose

Strategies:
- Act as middleware within the bot framework, processing each tick to determine if any trading actions are needed.
- Use technical indicators, price history, and other data provided in the `TickContext` to make trading decisions.
- Trigger orders by calling `ExecuteOrder` in `TickContext` based on custom logic (e.g., a Moving Average Crossover or RSI strategy).

## Strategy Structure

In the bot framework, a strategy is implemented as a middleware function that accepts a `TickContext` as its input and returns an error if any issues occur during execution. This design allows strategies to easily access all the necessary information (e.g., market data, indicators, order functions) to make trading decisions.

### Strategy Interface

Each strategy is essentially a function of type `types.Middleware`, which is defined as follows:

```go
type Middleware func(ctx *TickContext) error
```

- **Input**: `*TickContext` (provides real-time data, indicators, and order execution functionality).
- **Output**: `error` (returns `nil` if successful, or an error if something goes wrong).

### Key Components in `TickContext` for Strategies

- **`MarketData`**: Real-time price and volume.
- **`Indicators`**: Computed technical indicators (e.g., SMA, EMA).
- **`Store`**: Access to historical data.
- **`ExecuteOrder`**: Function to place buy or sell orders.

---

## Example Strategy Implementations

Below are examples of two strategies: a **Moving Average Crossover Strategy** and a **Relative Strength Index (RSI) Strategy**. Each strategy is designed to make trading decisions based on specific technical indicators.

### 1. Moving Average Crossover Strategy

The **Moving Average Crossover Strategy** generates buy or sell signals based on the crossover of two moving averages (e.g., short-term and long-term SMAs). When the short SMA crosses above the long SMA, it signals a potential uptrend (buy), and when the short SMA crosses below the long SMA, it signals a downtrend (sell).

#### `ma_crossover.go`

```go
package strategies

import (
    "fmt"
    "trading-bot/pkg/types"
)

// MovingAverageCrossoverStrategy is a simple strategy that places buy or sell orders
// based on the crossover of short and long moving averages.
func MovingAverageCrossoverStrategy() types.Middleware {
    return func(ctx *types.TickContext) error {
        // Access indicators for short and long moving averages
        shortSMA := ctx.Indicators["SMA_50"]
        longSMA := ctx.Indicators["SMA_200"]

        // Determine crossover and execute trade
        if shortSMA > longSMA {
            // Buy signal: Place a market buy order
            if err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideBuy, 1.0, 0); err != nil {
                return fmt.Errorf("failed to execute buy order: %w", err)
            }
        } else if shortSMA < longSMA {
            // Sell signal: Place a market sell order
            if err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideSell, 1.0, 0); err != nil {
                return fmt.Errorf("failed to execute sell order: %w", err)
            }
        }
        return nil
    }
}
```

- **Indicator Usage**: `SMA_50` and `SMA_200` (short and long Simple Moving Averages).
- **Trading Logic**: Buy when `SMA_50` > `SMA_200`, and sell when `SMA_50` < `SMA_200`.

### 2. Relative Strength Index (RSI) Strategy

The **RSI Strategy** generates buy signals when the RSI is below a certain threshold (indicating oversold conditions) and sell signals when the RSI is above a threshold (indicating overbought conditions).

#### `rsi_strategy.go`

```go
package strategies

import (
    "fmt"
    "trading-bot/pkg/types"
)

// RSIThresholdStrategy places buy/sell orders based on RSI thresholds.
func RSIThresholdStrategy(buyThreshold, sellThreshold float64) types.Middleware {
    return func(ctx *types.TickContext) error {
        // Access the RSI indicator value
        rsi := ctx.Indicators["RSI"]

        // Buy if RSI is below the buy threshold (oversold)
        if rsi < buyThreshold {
            if err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideBuy, 1.0, 0); err != nil {
                return fmt.Errorf("failed to execute buy order: %w", err)
            }
        }

        // Sell if RSI is above the sell threshold (overbought)
        if rsi > sellThreshold {
            if err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideSell, 1.0, 0); err != nil {
                return fmt.Errorf("failed to execute sell order: %w", err)
            }
        }
        return nil
    }
}
```

- **Indicator Usage**: `RSI` (Relative Strength Index).
- **Trading Logic**: Buy when RSI < `buyThreshold` (e.g., 30) and sell when RSI > `sellThreshold` (e.g., 70).

---

## Registering Strategies with the Bot

Once you have defined strategies, they must be registered with the bot for specific connectors and trading pairs. This allows the bot to apply each strategy as market data ticks are received.

### Example Usage in `main.go`

```go
package main

import (
    "trading-bot/internal/connectors"
    "trading-bot/internal/strategies"
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

    // Register Moving Average Crossover strategy for BTC/USDT on Binance
    bot.RegisterMiddleware("Binance", "BTC/USDT", strategies.MovingAverageCrossoverStrategy())

    // Register RSI Strategy with thresholds for BTC/USDT on Binance
    bot.RegisterMiddleware("Binance", "BTC/USDT", strategies.RSIThresholdStrategy(30, 70))

    // Start the bot
    if err := bot.Start(); err != nil {
        logger.Fatal().Err(err).Msg("Failed to start bot")
    }
}
```

---

## Summary

### What are Strategies?

Strategies are functions that act on incoming market data (ticks) to make trading decisions. They access real-time data, indicators, and order execution methods through the `TickContext`.

### Strategy Components

- **Middleware Type**: A function that takes `TickContext` as input and returns an error.
- **Indicators**: Precomputed values like SMA, EMA, or RSI, accessed through `TickContext`.
- **Order Execution**: The `ExecuteOrder` function in `TickContext` allows strategies to place buy/sell orders directly.

### Common Strategies

- **Moving Average Crossover**: Buy when a short SMA crosses above a long SMA; sell when it crosses below.
- **RSI Strategy**: Buy when RSI is below a certain threshold (oversold); sell when RSI is above a certain threshold (overbought).

### Adding Strategies

- Define strategies in the `strategies` directory as functions that match the `Middleware` type.
- Register strategies with the bot for specific connectors and trading pairs.

This modular design allows for flexibility, making it easy to create, test, and modify strategies to adapt to different trading scenarios.