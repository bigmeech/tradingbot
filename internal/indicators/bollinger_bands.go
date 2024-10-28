package indicators

import (
	"fmt"
	"math"
)

// BollingerBands represents the Bollinger Bands indicator.
type BollingerBands struct {
	period     int
	multiplier float64
	name       string
}

// NewBollingerBands creates a new Bollinger Bands instance with a given period and multiplier.
func NewBollingerBands(period int, multiplier float64) *BollingerBands {
	return &BollingerBands{
		period:     period,
		multiplier: multiplier,
		name:       fmt.Sprintf("BollingerBands_%d", period),
	}
}

// Calculate computes the upper, middle, and lower Bollinger Bands.
func (b *BollingerBands) Calculate(data []float64) []float64 {
	if len(data) < b.period {
		return nil // Return nil for insufficient data
	}

	// Calculate the SMA for the period
	sma := NewSMA(b.period).Calculate(data)
	sumSquaredDiffs := 0.0

	// Calculate standard deviation of the last N prices
	for _, price := range data[len(data)-b.period:] {
		sumSquaredDiffs += math.Pow(price-sma, 2)
	}
	standardDeviation := math.Sqrt(sumSquaredDiffs / float64(b.period))

	// Return a slice representing the upper, middle, and lower bands
	return []float64{
		sma + (b.multiplier * standardDeviation), // Upper Band
		sma,                                      // Middle Band
		sma - (b.multiplier * standardDeviation), // Lower Band
	}
}

// Name returns the name of the indicator.
func (b *BollingerBands) Name() string {
	return b.name
}

// Period returns the period of the Bollinger Bands.
func (b *BollingerBands) Period() int {
	return b.period
}
