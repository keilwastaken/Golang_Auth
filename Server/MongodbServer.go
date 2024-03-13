package Server

import (
	"Clarity_go/Config"
	"Clarity_go/Interfaces"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoService struct {
	Client         *mongo.Client
	UserCollection *mongo.Collection
}

// Ensure MongoService implements IMongoService
var _ Interfaces.IMongoService = (*MongoService)(nil)

func NewMongoService(pEnvironment *Config.Config) *MongoService {
	xClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(pEnvironment.MongodbURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = xClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Fetching Collections")
	xUserCollection := xClient.Database("Clarity").Collection("users")

	fmt.Println("Connected to MongoDB and this is working")
	return &MongoService{
		Client:         xClient,
		UserCollection: xUserCollection,
	}
}

func (ms MongoService) GetUserCollection() *mongo.Collection {
	return ms.UserCollection
}
