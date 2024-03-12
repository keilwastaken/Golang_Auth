package Helpers

import (
	"Clarity_go/Interfaces"
	"Clarity_go/Models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
)

type Token struct {
}

var _ Interfaces.IToken = (*Token)(nil)

func (t Token) GenerateToken(pId primitive.ObjectID) (*string, *Models.ResponseError) {

	//TODO TOKEN LIFE SPAN IN ENV FILE
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": pId,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	SignedToken, SigningError := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if SigningError != nil {

		return nil, &Models.ResponseError{
			Message: "Failed to create token",
			Status:  http.StatusBadRequest,
		}
	}

	return &SignedToken, nil
}
