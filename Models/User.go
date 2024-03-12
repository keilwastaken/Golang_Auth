package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"` //MongoDb document ID
	Email    string             `bson:"email" binding:"required"`
	Password string             `bson:"password" binding:"required"`
}
