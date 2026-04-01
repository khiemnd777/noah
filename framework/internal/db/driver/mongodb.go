package driver

import (
	"context"
	"time"

	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	config frameworkdb.MongoConfig
}

func NewMongoClient(cfg frameworkdb.MongoConfig) *MongoClient {
	return &MongoClient{config: cfg}
}

func (m *MongoClient) Provider() string {
	return "mongodb"
}

func (m *MongoClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.config.URI))
	if err != nil {
		return err
	}

	m.client = client
	return nil
}

func (m *MongoClient) Close() error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(context.Background())
}
