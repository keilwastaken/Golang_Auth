package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RefreshToken struct {
	_id       primitive.ObjectID `bson:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId" binding:"required"`
	Token     string             `bson:"token" binding:"required"`
	CreatedAt time.Time          `bson:"createdAt" binding:"required"`
	ExpiresAt time.Time          `bson:"expiresAt" binding:"required"`
}
