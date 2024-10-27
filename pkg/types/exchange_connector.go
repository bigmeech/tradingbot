package types

type ExchangeConnector interface {
	MarketDataStreamer
	OrderExecutor
	Connect() error
	Disconnect() error
}
