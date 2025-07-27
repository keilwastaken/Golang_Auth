package Interfaces

import (
	"goauth/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IToken interface {
	GenerateAccessToken(pId primitive.ObjectID) (*string, *Models.ResponseError)
	GenerateRefreshToken(pId primitive.ObjectID) (*string, *Models.ResponseError)
	ValidateToken(tokenString string) (primitive.ObjectID, *Models.ResponseError)
}
