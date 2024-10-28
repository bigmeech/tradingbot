package indicators

import (
	"testing"
)

func TestEMA_Calculate(t *testing.T) {
	ema := NewEMA(3)

	// Case 1: Not enough data to calculate EMA
	data := []float64{10, 20} // Only 2 points, while EMA requires 3
	result := ema.Calculate(data)
	if result != 0 {
		t.Errorf("Expected 0 for insufficient data, got %v", result)
	}

	// Case 2: Exactly enough data to calculate initial EMA
	data = []float64{10, 20, 30} // Start with SMA of (10+20+30) / 3 = 20, then apply EMA
	expectedInitialEma := 20.0
	result = ema.Calculate(data)
	if result != expectedInitialEma {
		t.Errorf("Expected initial EMA of %v, got %v", expectedInitialEma, result)
	}

	// Case 3: Adding more data for EMA continuation
	data = []float64{10, 20, 30, 40} // Expected EMA of 30, approximately
	expectedContinuationEma := 30.0
	result = ema.Calculate(data)
	if result < expectedContinuationEma-0.1 || result > expectedContinuationEma+0.1 {
		t.Errorf("Expected EMA of around %v, got %v", expectedContinuationEma, result)
	}
}
