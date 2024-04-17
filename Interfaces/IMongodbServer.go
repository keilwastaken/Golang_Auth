package Interfaces

import "go.mongodb.org/mongo-driver/mongo"

type IMongoService interface {
	GetUserCollection() *mongo.Collection
	GetRefreshTokenCollection() *mongo.Collection
	// Add other necessary methods here
}
