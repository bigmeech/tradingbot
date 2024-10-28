# TickContext Documentation

`TickContext` is a central data structure within the trading bot framework, representing the context for each incoming market data tick. It provides essential information for trading strategies (middleware), including details about the market, trading pair, real-time price and volume, computed technical indicators, an interface for querying historical data, and a function to execute trades.

## Purpose

The `TickContext` is designed to:
- Pass market data (ticks) from connectors to registered strategies.
- Provide strategies with real-time market information, technical indicators, and access to historical data.
- Enable strategies to place buy and sell orders based on defined logic.

This structure encapsulates all necessary data and functions, making it easy for developers to implement, test, and modify trading strategies.

## TickContext Structure

```go
// TickContext represents the context for a single market data tick.
type TickContext struct {
    MarketName   string
    TradingPair  string
    MarketData   *MarketData
    Store        Store
    Indicators   map[string]float64
    ExecuteOrder func(orderType OrderType, side OrderSide, amount, price float64) error
}
```

### Field Descriptions

1. **`MarketName`** (`string`):
    - **Description**: The name of the exchange or market (e.g., "Binance", "Kraken") from which the tick data originated.
    - **Usage**: Useful for identifying the source of the data, allowing strategies to implement market-specific logic if needed.

2. **`TradingPair`** (`string`):
    - **Description**: The asset pair associated with this tick, such as "BTC/USDT".
    - **Usage**: Allows strategies to apply specific logic based on the asset pair being traded.

3. **`MarketData`** (`*MarketData`):
    - **Description**: Contains real-time tick data for the trading pair, including fields like `Price` and `Volume`.
    - **Structure**:
      ```go
      type MarketData struct {
          Price  float64 // Latest price for the trading pair
          Volume float64 // Trading volume for the pair
      }
      ```
    - **Usage**: Provides the latest market conditions, which are essential for making informed trading decisions.

4. **`Store`** (`Store`):
    - **Description**: An interface to access historical market data for the trading pair. The store provides methods to retrieve historical price data, which can be valuable for indicators and strategy calculations.
    - **Example Usage**:
      ```go
      priceHistory := ctx.Store.QueryPriceHistory(ctx.TradingPair, 50) // Get the last 50 prices
      ```
    - **Purpose**: Allows strategies to access past market data, useful for calculating indicators or applying strategies that depend on historical trends.

5. **`Indicators`** (`map[string]float64`):
    - **Description**: A map of calculated technical indicator values relevant to the trading pair.
    - **Example**:
      ```go
      ctx.Indicators["SMA_50"] // Access the 50-period Simple Moving Average
      ```
    - **Usage**: Stores precomputed indicator values (e.g., moving averages) to help strategies analyze trends and make trade decisions based on those indicators.

6. **`ExecuteOrder`** (`func(orderType OrderType, side OrderSide, amount, price float64) error`):
    - **Description**: A function that enables strategies to execute buy or sell orders based on specific conditions.
    - **Parameters**:
        - `orderType` (`OrderType`): The type of order to place (e.g., `MARKET`, `LIMIT`).
        - `side` (`OrderSide`): The trade direction (`BUY` or `SELL`).
        - `amount` (`float64`): The amount to trade.
        - `price` (`float64`): The price at which to execute the order (used for limit orders).
    - **Usage**: Abstracts order execution, making it easy for strategies to place orders without needing direct access to the connector.
    - **Example**:
      ```go
      ctx.ExecuteOrder(OrderTypeMarket, OrderSideBuy, 1.0, 0) // Executes a market buy order
      ```

---

## Example Usage in a Strategy

Here’s an example of how `TickContext` might be used within a Moving Average Crossover strategy. This strategy buys when the short-term moving average crosses above the long-term moving average and sells when it crosses below.

```go
package strategies

import (
    "fmt"
    "trading-bot/pkg/types"
)

// MovingAverageCrossoverStrategy is a strategy that places buy or sell orders
// based on the crossover of short and long moving averages.
func MovingAverageCrossoverStrategy(ctx *types.TickContext) error {
    // Log the current market data
    fmt.Printf("Market: %s | Pair: %s | Price: %f | Volume: %f\n",
        ctx.MarketName, ctx.TradingPair, ctx.MarketData.Price, ctx.MarketData.Volume)
    
    // Access indicators
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
```

In this strategy:
- **`MarketName`** and **`TradingPair`** identify the data source and asset pair, respectively.
- **`MarketData`** provides current price and volume, which can be logged or used in strategy logic.
- **`Indicators`** is accessed to get short-term and long-term moving averages for the trading pair.
- **`ExecuteOrder`** is used to place market orders based on crossover conditions.

---

## Summary

The `TickContext` structure provides:
- **Real-Time Data**: Access to the latest price and volume data.
- **Historical Data**: Via the `Store`, enabling strategies to retrieve price history for indicator calculations.
- **Technical Indicators**: A map of computed indicators (like SMA) available for decision-making.
- **Trade Execution**: An easy-to-use API for placing buy/sell orders, abstracting the underlying order execution logic.

This design allows developers to implement trading strategies that are both powerful and easy to manage, with all necessary data and functions available in one structure.