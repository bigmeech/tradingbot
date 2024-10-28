package tradingbot

import (
	"fmt"
	"github.com/bigmeech/tradingbot/internal/framework"
	"github.com/bigmeech/tradingbot/pkg/types"
	"github.com/rs/zerolog"
)

type Bot struct {
	fw        *framework.Framework
	logger    zerolog.Logger
	debugMode bool
}

// NewBot initializes a new Bot instance with a StoreManager.
func NewBot(largeStore types.Store, bufferSize, threshold int, logger zerolog.Logger) *Bot {
	// Create StoreManager
	storeManager := framework.NewStoreManager(largeStore, bufferSize, threshold)

	return &Bot{
		fw:        framework.NewFramework(storeManager),
		logger:    logger,
		debugMode: false,
	}
}

// EnableDebug enables debug mode for the bot, allowing detailed logging.
func (b *Bot) EnableDebug() {
	b.debugMode = true
	b.logger = b.logger.Level(zerolog.DebugLevel)
}

// RegisterConnector registers a trading connector.
func (b *Bot) RegisterConnector(name string, connector types.Connector) {
	b.fw.RegisterConnector(name, connector)
}

// RegisterMiddleware adds middleware for a specific market and trading pair.
func (b *Bot) RegisterMiddleware(marketName, tradingPair string, mw types.Middleware) {
	b.fw.RegisterMiddleware(marketName, tradingPair, mw)
}

// RegisterIndicator registers an indicator for a specific market and trading pair.
func (b *Bot) RegisterIndicator(marketName, tradingPair string, indicator types.Indicator) {
	b.fw.RegisterIndicator(marketName, tradingPair, indicator)
}

// Start begins processing data from connectors and applying registered indicators and strategies.
func (b *Bot) Start() error {
	if len(b.fw.Connectors()) == 0 {
		return fmt.Errorf("no connectors registered; add at least one connector to start the bot")
	}
	b.fw.Start(b.ProcessTick) // Pass ProcessTick as the callback
	return nil
}

// ProcessTick handles incoming ticks with logging and runs indicators on recent data.
func (b *Bot) ProcessTick(ctx *types.TickContext) {
	if b.debugMode {
		b.logger.Debug().
			Str("TradingPair", ctx.TradingPair).
			Float64("Price", ctx.MarketData.Price).
			Float64("Volume", ctx.MarketData.Volume).
			Msg("Received tick")
	}

	// Compute and update indicators in the context
	indicators := b.fw.GetIndicators(ctx.MarketName, ctx.TradingPair)
	for _, indicator := range indicators {
		period := indicator.Period()                                                    // Get the period from each indicator
		priceHistory := b.fw.QueryPriceHistory(ctx.MarketName, ctx.TradingPair, period) // Fetch recent or historical data
		ctx.Indicators[indicator.Name()] = indicator.Calculate(priceHistory)
	}

	// Execute any middleware or strategy logic here
	for _, middleware := range b.fw.GetMiddleware(ctx.MarketName, ctx.TradingPair) {
		if err := middleware(ctx); err != nil {
			b.logger.Error().Err(err).Msg("Error processing middleware")
		}
	}
}
