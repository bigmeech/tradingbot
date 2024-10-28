package strategies

import (
	"trading-bot/pkg/types"
)

func MovingAverageCrossoverStrategy() types.Middleware {
	return func(ctx *types.TickContext) error {
		shortSMA := ctx.Indicators["SMA_50"]
		longSMA := ctx.Indicators["SMA_200"]

		if shortSMA > longSMA {
			// Execute a market buy order if the short SMA crosses above the long SMA
			err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideBuy, 1.0, ctx.MarketData.Price)
			if err != nil {
				return err
			}
		} else if shortSMA < longSMA {
			// Execute a market sell order if the short SMA crosses below the long SMA
			err := ctx.ExecuteOrder(types.OrderTypeMarket, types.OrderSideSell, 1.0, ctx.MarketData.Price)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
