package framework

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"trading-bot/pkg/types"
)

type MongoDBLargeStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
}

// NewMongoDBLargeStore initializes a connection to MongoDB.
func NewMongoDBLargeStore(uri, dbName, collectionName string) (*MongoDBLargeStore, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection(collectionName)

	return &MongoDBLargeStore{
		client:     client,
		collection: collection,
		ctx:        ctx,
	}, nil
}

// RecordTick stores new market data in MongoDB.
func (m *MongoDBLargeStore) RecordTick(tradingPair string, marketData *types.MarketData) error {
	_, err := m.collection.InsertOne(m.ctx, marketData)
	return err
}

// QueryPriceHistory retrieves price history from MongoDB.
func (m *MongoDBLargeStore) QueryPriceHistory(tradingPair string, period int) []float64 {
	cursor, err := m.collection.Find(m.ctx, nil, options.Find().SetSort(map[string]int{"_id": -1}).SetLimit(int64(period)))
	if err != nil {
		return nil
	}
	defer cursor.Close(m.ctx)

	prices := []float64{}
	for cursor.Next(m.ctx) {
		var data types.MarketData
		cursor.Decode(&data)
		prices = append(prices, data.Price)
	}
	return prices
}
