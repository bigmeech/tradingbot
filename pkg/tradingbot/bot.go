package tradingbot

import (
	"fmt"
	"trading-bot/internal/framework"
	"trading-bot/pkg/models"
	"trading-bot/pkg/types"
)

type Bot struct {
	fw *framework.Framework
}

// NewBot initializes a new Bot instance with configurable FastStore, LargeStore, and threshold.
// The threshold determines the cut-off point for when to switch between FastStore and LargeStore.
func NewBot(fastStore models.FastStore, largeStore models.LargeStore, threshold int) *Bot {
	storeManager := framework.NewStoreManager(fastStore, largeStore, threshold)
	return &Bot{
		fw: framework.NewFramework(storeManager),
	}
}

// RegisterConnector registers a trading connector.
func (b *Bot) RegisterConnector(name string, connector types.Connector) {
	b.fw.RegisterConnector(name, connector)
}

// RegisterIndicator registers an indicator for a specific market and trading pair.
func (b *Bot) RegisterIndicator(marketName, tradingPair string, indicator types.Indicator) {
	b.fw.RegisterIndicator(marketName, tradingPair, indicator)
}

// AddStrategy adds a strategy as middleware for a specific market and trading pair.
func (b *Bot) AddStrategy(marketName, tradingPair string, strategy types.Middleware) {
	b.fw.UsePair(marketName, tradingPair, strategy)
}

// Start begins processing data from connectors and applying registered indicators and strategies.
func (b *Bot) Start() error {
	if len(b.fw.Connectors()) == 0 {
		return fmt.Errorf("no connectors registered; add at least one connector to start the bot")
	}
	b.fw.Start()
	return nil
}
