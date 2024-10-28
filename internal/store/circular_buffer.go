package store

import "github.com/bigmeech/tradingbot/pkg/types"

// CircularBuffer holds a fixed-size buffer of MarketData for quick, recent access.
type CircularBuffer struct {
	data   []types.MarketData // Slice holding the circular buffer data
	size   int                // Total size of the buffer
	index  int                // Current index in the buffer
	isFull bool               // Tracks if buffer has wrapped around to start
}

// NewCircularBuffer initializes a new CircularBuffer with a fixed size.
func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data:   make([]types.MarketData, size),
		size:   size,
		index:  0,
		isFull: false,
	}
}

// Add inserts a new data point into the circular buffer.
func (cb *CircularBuffer) Add(data types.MarketData) {
	cb.data[cb.index] = data
	cb.index = (cb.index + 1) % cb.size
	if cb.index == 0 {
		cb.isFull = true
	}
}

// GetData retrieves the most recent `count` data points, or fewer if insufficient data.
func (cb *CircularBuffer) GetData(count int) []types.MarketData {
	if count > cb.size {
		count = cb.size
	}
	if !cb.isFull {
		if count > cb.index {
			count = cb.index
		}
		return cb.data[:count]
	}

	// For a full buffer, return the most recent `count` data points in correct order
	start := (cb.index - count + cb.size) % cb.size
	if start+count <= cb.size {
		return cb.data[start : start+count]
	}

	// Handle wrap-around case for circular buffer
	return append(cb.data[start:], cb.data[:start+count-cb.size]...)
}
