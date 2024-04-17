package Helpers

import (
	"Clarity_go/Interfaces"
	"Clarity_go/Models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strings"
	"time"
)

type Token struct {
}

// TODO BEST TO ROTATE REFRESH TOKENS? WE SHOULD LOOK INTO IT.
var _ Interfaces.IToken = (*Token)(nil)

func (t Token) GenerateAccessToken(pId primitive.ObjectID) (*string, *Models.ResponseError) {

	//TODO TOKEN LIFE SPAN IN ENV FILE
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": pId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * 15).Unix(),
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

func (t Token) GenerateRefreshToken(pId primitive.ObjectID) (*string, *Models.ResponseError) {

	//TODO TOKEN LIFE SPAN IN ENV FILE
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": pId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
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

// ValidateToken parses and validates a given JWT token string.
// Returns the user ID from the token's claims if the token is valid,
// and an error if it's not valid or an error occurs during validation.
func (t Token) ValidateToken(tokenString string) (primitive.ObjectID, *Models.ResponseError) {
	// Check bearer token format
	// because its validation we want to make sure the token is in the correct format.
	if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
		tokenString = tokenString[7:]
	}

	// Define a claims structure to extract the user ID (sub) from the token.
	var claims jwt.MapClaims

	// check token signature
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's algorithm matches the expected algorithm.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,
				fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the key used to sign the token for verification.
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	//if there is an error parsing the token
	if err != nil {
		return primitive.NilObjectID, &Models.ResponseError{
			Message: "Failed to parse token, Error: " + err.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	//if the token is not valid
	if !token.Valid {
		return primitive.NilObjectID,
			&Models.ResponseError{
				Message: "Invalid token",
				Status:  http.StatusUnauthorized,
			}
	}

	// Convert the user ID from the token (sub claim) into a primitive.ObjectID.
	subj, err := claims.GetSubject()
	if err != nil {
		return primitive.NilObjectID, &Models.ResponseError{
			Message: "Invalid user ID in token",
			Status:  http.StatusUnauthorized,
		}
	}

	userID, err := primitive.ObjectIDFromHex(subj)
	if err != nil {
		return primitive.NilObjectID, &Models.ResponseError{
			Message: "Invalid user ID in token",
			Status:  http.StatusUnauthorized,
		}
	}

	return userID, nil
}
