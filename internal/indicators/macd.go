package indicators

import "fmt"

// MACD represents a Moving Average Convergence Divergence indicator.
type MACD struct {
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
	name         string
}

// NewMACD creates a new MACD instance with specified periods for fast, slow, and signal lines.
func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		fastPeriod:   fastPeriod,
		slowPeriod:   slowPeriod,
		signalPeriod: signalPeriod,
		name:         fmt.Sprintf("MACD_%d_%d_%d", fastPeriod, slowPeriod, signalPeriod),
	}
}

// Calculate computes the MACD value based on price data.
func (m *MACD) Calculate(data []float64) float64 {
	if len(data) < m.slowPeriod {
		return 0.0 // Insufficient data to calculate MACD
	}

	// Calculate the fast and slow EMAs
	fastEma := NewEMA(m.fastPeriod).Calculate(data)
	slowEma := NewEMA(m.slowPeriod).Calculate(data)

	// MACD line is the difference between fast and slow EMAs
	return fastEma - slowEma
}

// Name returns the name of the indicator.
func (m *MACD) Name() string {
	return m.name
}

// Period returns the slow period of the MACD as the minimum data required.
func (m *MACD) Period() int {
	return m.slowPeriod
}
