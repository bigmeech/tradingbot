package indicators

import (
	"math"
	"testing"
)

func TestRSI_Calculate(t *testing.T) {
	rsi := NewRSI(3)

	// Case 1: Not enough data to calculate RSI
	data := []float64{50, 51} // Only 2 points, while RSI requires 3
	result := rsi.Calculate(data)
	if result != 0 {
		t.Errorf("Expected 0 for insufficient data, got %v", result)
	}

	// Case 2: Exactly enough data to calculate RSI
	data = []float64{45, 50, 55} // Known upward trend, should yield a high RSI
	result = rsi.Calculate(data)
	expected := 100.0 // Since there's only gain
	if math.Abs(result-expected) > 1e-2 {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Case 3: Alternating prices leading to a mixed RSI
	data = []float64{50, 55, 45} // Mixed gain and loss
	result = rsi.Calculate(data)
	if result <= 0 || result >= 100 {
		t.Errorf("Expected RSI between 0 and 100, got %v", result)
	}
}
