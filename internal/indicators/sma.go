package indicators

type SMA struct {
	Period int
}

func (s *SMA) Calculate(data []float64) float64 {
	if len(data) < s.Period {
		return 0
	}
	sum := 0.0
	for i := len(data) - s.Period; i < len(data); i++ {
		sum += data[i]
	}
	return sum / float64(s.Period)
}

func (s *SMA) Name() string {
	return "SMA"
}
