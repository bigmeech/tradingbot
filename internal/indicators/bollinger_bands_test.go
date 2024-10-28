package indicators

import (
	"math"
	"testing"
)

func TestBollingerBands_Calculate(t *testing.T) {
	bb := NewBollingerBands(3, 2.0) // Period=3, Multiplier=2

	// Case 1: Not enough data to calculate Bollinger Bands
	data := []float64{100, 105} // Only 2 points, Bollinger Bands need 3
	result := bb.Calculate(data)
	if result != nil {
		t.Errorf("Expected nil for insufficient data, got %v", result)
	}

	// Case 2: Sufficient data for Bollinger Bands
	data = []float64{100, 105, 110} // Expected SMA: (100+105+110)/3 = 105
	expectedSMA := 105.0

	result = bb.Calculate(data)
	if len(result) != 3 {
		t.Errorf("Expected 3 values (upper, middle, lower bands), got %d", len(result))
	}

	// Validate middle band (SMA)
	if math.Abs(result[1]-expectedSMA) > 0.1 {
		t.Errorf("Expected middle band (SMA) of %v, got %v", expectedSMA, result[1])
	}

	// Calculate standard deviation for further validation
	variance := (math.Pow(100-expectedSMA, 2) + math.Pow(105-expectedSMA, 2) + math.Pow(110-expectedSMA, 2)) / 3
	stdDev := math.Sqrt(variance)

	// Validate upper and lower bands
	expectedUpper := expectedSMA + (2 * stdDev)
	expectedLower := expectedSMA - (2 * stdDev)
	if math.Abs(result[0]-expectedUpper) > 0.1 {
		t.Errorf("Expected upper band of %v, got %v", expectedUpper, result[0])
	}
	if math.Abs(result[2]-expectedLower) > 0.1 {
		t.Errorf("Expected lower band of %v, got %v", expectedLower, result[2])
	}
}
