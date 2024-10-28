package indicators

import (
	"testing"
)

func TestMACD_Calculate(t *testing.T) {
	macd := NewMACD(3, 6, 3) // Fast EMA=3, Slow EMA=6, Signal EMA=3

	// Case 1: Not enough data to calculate MACD
	data := []float64{10, 20, 30} // Only enough for fast EMA, not slow EMA
	result := macd.Calculate(data)
	if result != 0 {
		t.Errorf("Expected 0 for insufficient data, got %v", result)
	}

	// Case 2: Sufficient data for MACD
	data = []float64{10, 20, 30, 40, 50, 60}

	// Calculate fast and slow EMAs directly to validate the MACD line
	fastEMA := NewEMA(3).Calculate(data)
	slowEMA := NewEMA(6).Calculate(data)
	expectedMacd := fastEMA - slowEMA

	result = macd.Calculate(data)
	if result < expectedMacd-0.1 || result > expectedMacd+0.1 {
		t.Errorf("Expected MACD of around %v, got %v", expectedMacd, result)
	}
}
