package strategies

import (
	"trading-bot/pkg/types"
)

func MovingAverageCrossoverStrategy() types.Middleware {
	return func(ctx *types.TickContext) error {
		shortSMA := ctx.Indicators["SMA_50"]
		longSMA := ctx.Indicators["SMA_200"]

		if shortSMA > longSMA {
			ctx.Actions.ExecuteAction(types.ActionBuy, ctx.TradingPair, 1.0)
		} else if shortSMA < longSMA {
			ctx.Actions.ExecuteAction(types.ActionSell, ctx.TradingPair, 1.0)
		}

		return nil
	}
}
