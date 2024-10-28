package indicators

import "fmt"

// In indicators/sma.go

type SMA struct {
	period int
	name   string
}

func NewSMA(period int) *SMA {
	return &SMA{
		period: period,
		name:   fmt.Sprintf("SMA_%d", period),
	}
}

func (s *SMA) Calculate(data []float64) float64 {
	if len(data) < s.period {
		return 0.0 // Insufficient data to calculate SMA
	}
	sum := 0.0
	for _, price := range data[len(data)-s.period:] {
		sum += price
	}
	return sum / float64(s.period)
}

func (s *SMA) Name() string {
	return s.name
}

func (s *SMA) Period() int {
	return s.period
}
