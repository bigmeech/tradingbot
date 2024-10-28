package framework

import (
	"context"
	"github.com/bigmeech/tradingbot/pkg/types"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisFastStore struct {
	client *redis.Client
	limit  int
	ctx    context.Context
	key    string // Redis key for storing trading data
}

// NewRedisFastStore initializes a RedisFastStore with a Redis client.
func NewRedisFastStore(addr, password string, db, limit int, key string) *RedisFastStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisFastStore{
		client: client,
		limit:  limit,
		ctx:    context.Background(),
		key:    key,
	}
}

// RecordTick adds new market data to Redis with a capped list.
func (r *RedisFastStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	price := strconv.FormatFloat(marketData.Price, 'f', -1, 64)
	err := r.client.LPush(r.ctx, r.key, price).Err()
	if err != nil {
		return err
	}
	r.client.LTrim(r.ctx, r.key, 0, int64(r.limit-1)) // Trim list to maintain limit
	return nil
}

// QueryPriceHistory retrieves recent price data from Redis.
func (r *RedisFastStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	if period > r.limit {
		period = r.limit
	}

	data, err := r.client.LRange(r.ctx, r.key, 0, int64(period-1)).Result()
	if err != nil {
		return nil
	}

	prices := make([]float64, len(data))
	for i, str := range data {
		prices[i], _ = strconv.ParseFloat(str, 64)
	}
	return prices
}
