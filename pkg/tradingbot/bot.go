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

// NewBot initializes a new Bot instance with configurable Store instances and threshold.
func NewBot(fastStore, largeStore types.Store, threshold int, logger zerolog.Logger) *Bot {
	storeManager := framework.NewStoreManager(fastStore, largeStore, threshold)
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

// Start begins processing data from connectors and applying registered indicators and strategies.
func (b *Bot) Start() error {
	if len(b.fw.Connectors()) == 0 {
		return fmt.Errorf("no connectors registered; add at least one connector to start the bot")
	}
	b.fw.Start(b.ProcessTick) // Pass ProcessTick as the callback
	return nil
}

// ProcessTick handles incoming ticks with logging.
func (b *Bot) ProcessTick(ctx *types.TickContext) {
	if b.debugMode {
		b.logger.Debug().
			Str("TradingPair", ctx.TradingPair).
			Float64("Price", ctx.MarketData.Price).
			Float64("Volume", ctx.MarketData.Volume).
			Msg("Received tick")
	}

	// Execute any middleware or strategy logic here
	for _, middleware := range b.fw.GetMiddleware(ctx.MarketName, ctx.TradingPair) {
		if err := middleware(ctx); err != nil {
			b.logger.Error().Err(err).Msg("Error processing middleware")
		}
	}
}
