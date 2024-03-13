package Interfaces

import (
	"Clarity_go/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IToken interface {
	GenerateToken(pId primitive.ObjectID) (*string, *Models.ResponseError)
	ValidateToken(tokenString string) (primitive.ObjectID, *Models.ResponseError)
}
