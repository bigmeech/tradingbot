package indicators

import "fmt"

// EMA represents an Exponential Moving Average indicator.
type EMA struct {
	period     int
	name       string
	multiplier float64
}

// NewEMA creates a new EMA instance with the specified period.
func NewEMA(period int) *EMA {
	multiplier := 2.0 / (float64(period) + 1.0)
	return &EMA{
		period:     period,
		name:       fmt.Sprintf("EMA_%d", period),
		multiplier: multiplier,
	}
}

// Calculate computes the EMA value based on the price data.
func (e *EMA) Calculate(data []float64) float64 {
	if len(data) < e.period {
		return 0.0 // Insufficient data to calculate EMA
	}

	// Start EMA calculation with the SMA of the first 'period' values
	ema := 0.0
	for _, price := range data[:e.period] {
		ema += price
	}
	ema /= float64(e.period)

	// Apply EMA formula to remaining prices
	for _, price := range data[e.period:] {
		ema = ((price - ema) * e.multiplier) + ema
	}
	return ema
}

// Name returns the name of the indicator.
func (e *EMA) Name() string {
	return e.name
}

// Period returns the period of the EMA.
func (e *EMA) Period() int {
	return e.period
}
