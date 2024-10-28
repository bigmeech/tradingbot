package indicators

import (
	"testing"
)

func TestSMA_Calculate(t *testing.T) {
	sma := NewSMA(3)

	// Case 1: Not enough data to calculate SMA
	data := []float64{10, 20} // Only 2 points, while SMA requires 3
	result := sma.Calculate(data)
	if result != 0 {
		t.Errorf("Expected 0 for insufficient data, got %v", result)
	}

	// Case 2: Exactly enough data to calculate SMA
	data = []float64{10, 20, 30} // Expected average: (10+20+30) / 3 = 20
	result = sma.Calculate(data)
	expected := 20.0
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Case 3: More than enough data, should use the last 3 points
	data = []float64{5, 15, 25, 35} // Last 3 points: (15+25+35) / 3 = 25
	result = sma.Calculate(data)
	expected = 25.0
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
