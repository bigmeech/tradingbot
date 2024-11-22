# Indicators Documentation

The `Indicators` component in the trading bot framework provides tools to calculate and track technical indicators based on market data. These indicators serve as valuable inputs for strategies, helping to identify market trends, momentum, and potential buy/sell signals.

## Purpose

Indicators:
- Offer calculated metrics (e.g., Simple Moving Average, Exponential Moving Average, Relative Strength Index) that analyze historical market data.
- Serve as inputs for strategies, allowing strategies to base decisions on precomputed values rather than recalculating them.
- Help smooth out market noise and highlight potential trading opportunities.

## Indicator Interface

To allow flexibility and reuse, each indicator implements the following `Indicator` interface:

```go
type Indicator interface {
    Calculate(data []float64) float64
    Name() string
}
```

### Indicator Interface Methods

1. **`Calculate(data []float64) float64`**:
    - Takes a slice of historical price data as input and returns the calculated value for the indicator.
    - For instance, a moving average might take the last N prices and return their average.

2. **`Name() string`**:
    - Returns the name of the indicator (e.g., "SMA_50" for a 50-period Simple Moving Average).
    - This name is used as the key in the `Indicators` map within `TickContext`.

---

## Example Indicator Implementations

Below are implementations of two common indicators: **Simple Moving Average (SMA)** and **Relative Strength Index (RSI)**. Each indicator is implemented as a struct that fulfills the `Indicator` interface.

### 1. Simple Moving Average (SMA)

The **Simple Moving Average (SMA)** is one of the most basic and widely used indicators. It calculates the average price over a specified number of periods, helping to smooth out price fluctuations and identify trends.

#### `sma.go`

```go
package indicators

import "trading-bot/pkg/types"

// SMA represents a Simple Moving Average indicator.
type SMA struct {
    period int
    name   string
}

// NewSMA creates a new SMA instance with the specified period.
func NewSMA(period int) *SMA {
    return &SMA{
        period: period,
        name:   fmt.Sprintf("SMA_%d", period),
    }
}

// Calculate computes the SMA value based on the last N periods of data.
func (s *SMA) Calculate(data []float64) float64 {
    if len(data) < s.period {
        return 0.0 // Insufficient data to calculate the SMA
    }

    // Sum the last N prices
    sum := 0.0
    for _, price := range data[len(data)-s.period:] {
        sum += price
    }

    // Return the average
    return sum / float64(s.period)
}

// Name returns the name of the indicator.
func (s *SMA) Name() string {
    return s.name
}
```

- **Parameters**: `period` defines the number of periods over which to calculate the average.
- **Usage**: The `Calculate` method takes the last N price data points and returns the average, providing a smoothed view of price trends.

### 2. Relative Strength Index (RSI)

The **Relative Strength Index (RSI)** is a momentum indicator that measures the speed and change of price movements. It oscillates between 0 and 100, with values above 70 often indicating overbought conditions and values below 30 indicating oversold conditions.

#### `rsi.go`

```go
package indicators

import (
    "trading-bot/pkg/types"
)

// RSI represents a Relative Strength Index indicator.
type RSI struct {
    period int
    name   string
}

// NewRSI creates a new RSI instance with the specified period.
func NewRSI(period int) *RSI {
    return &RSI{
        period: period,
        name:   fmt.Sprintf("RSI_%d", period),
    }
}

// Calculate computes the RSI value based on the last N periods of data.
func (r *RSI) Calculate(data []float64) float64 {
    if len(data) < r.period {
        return 0.0 // Insufficient data to calculate RSI
    }

    gain, loss := 0.0, 0.0
    for i := 1; i < r.period; i++ {
        change := data[len(data)-r.period+i] - data[len(data)-r.period+i-1]
        if change > 0 {
            gain += change
        } else {
            loss -= change // Make it positive
        }
    }

    if loss == 0 {
        return 100 // Prevent division by zero
    }

    avgGain := gain / float64(r.period)
    avgLoss := loss / float64(r.period)
    rs := avgGain / avgLoss

    // RSI formula
    return 100 - (100 / (1 + rs))
}

// Name returns the name of the indicator.
func (r *RSI) Name() string {
    return r.name
}
```

- **Parameters**: `period` defines the number of periods over which to calculate RSI.
- **Usage**: The `Calculate` method analyzes recent price changes, identifying if the market is overbought or oversold.

---

## Using Indicators in Strategies

Indicators are typically precomputed for each tick and stored in the `TickContext`’s `Indicators` map, making them readily accessible to strategies. Below is an example of how strategies can access and use these indicators.

### Example Usage in a Strategy

```go
package strategies

import (
    "fmt"
    "trading-bot/pkg/types"
)

// MovingAverageCrossoverStrategy is a simple strategy that places buy or sell orders
// based on the crossover of short and long moving averages.
func MovingAverageCrossoverStrategy(ctx *types.TickContext) error {
    // Access indicators
    shortSMA := ctx.Indicators["SMA_50"]
    longSMA := ctx.Indicators["SMA_200"]

    // Determine crossover and execute trade
    if shortSMA > longSMA {
        // Buy signal
        if err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideBuy, 1.0, 0); err != nil {
            return fmt.Errorf("failed to execute buy order: %w", err)
        }
    } else if shortSMA < longSMA {
        // Sell signal
        if err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideSell, 1.0, 0); err != nil {
            return fmt.Errorf("failed to execute sell order: %w", err)
        }
    }
    return nil
}
```

- **Indicator Access**: The strategy accesses `SMA_50` and `SMA_200` directly from `ctx.Indicators`.
- **Trading Logic**: The strategy generates buy/sell signals based on the crossover of the two SMAs.

---

## Registering and Using Indicators in the Bot

Indicators are typically computed and updated periodically based on the historical data available in the store. Here’s an example of how to register and calculate indicators in the bot.

### Example Integration in `main.go`

```go
package main

import (
    "trading-bot/internal/connectors"
    "trading-bot/internal/indicators"
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

    // Create and register indicators
    sma50 := indicators.NewSMA(50)
    sma200 := indicators.NewSMA(200)
    rsi14 := indicators.NewRSI(14)

    bot.RegisterIndicator("Binance", "BTC/USDT", sma50)
    bot.RegisterIndicator("Binance", "BTC/USDT", sma200)
    bot.RegisterIndicator("Binance", "BTC/USDT", rsi14)

    // Register a strategy that uses these indicators
    bot.RegisterMiddleware("Binance", "BTC/USDT", strategies.MovingAverageCrossoverStrategy())

    // Start the bot
    if err := bot.Start(); err != nil {
        logger.Fatal().Err(err).Msg("Failed to start bot")
    }
}
```

---

## Summary

### What are Indicators?

Indicators are calculated values based on historical data, used to identify trading opportunities and trends. They smooth out data fluctuations and help strategies make data-driven decisions.

### Common Indicators

- **Simple Moving Average (SMA)**: Average price over a set period, used to identify trend direction.
- **Relative Strength Index (RSI)**: Measures momentum, identifying overbought and oversold conditions.

### Adding Indicators

- Implement indicators in the `indicators` package by defining structs that fulfill the `Indicator` interface.
- Register indicators with the bot to make them available in `TickContext` for strategies.

This design ensures modularity and flexibility, allowing developers to add, customize, and use indicators effectively within the trading bot.

---