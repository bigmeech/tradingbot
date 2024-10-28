package indicators

import "fmt"

// In indicators/rsi.go

type RSI struct {
	period int
	name   string
}

func NewRSI(period int) *RSI {
	return &RSI{
		period: period,
		name:   fmt.Sprintf("RSI_%d", period),
	}
}

func (r *RSI) Calculate(data []float64) float64 {
	if len(data) < r.period {
		return 0.0 // Insufficient data to calculate RSI
	}
	gain, loss := 0.0, 0.0
	for i := 1; i < r.period; i++ {
		change := data[len(data)-r.period+i] - data[len(data)-r.period+i-1]
		if change > 0 {
			gain += change
		} else {
			loss -= change
		}
	}
	if loss == 0 {
		return 100
	}
	avgGain := gain / float64(r.period)
	avgLoss := loss / float64(r.period)
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

func (r *RSI) Name() string {
	return r.name
}

func (r *RSI) Period() int {
	return r.period
}
