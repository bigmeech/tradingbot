package types

type MarketDataHandler func(ctx *TickContext)

type MarketDataStreamer interface {
	StartStreaming(handler MarketDataHandler) error
	StopStreaming() error
}
