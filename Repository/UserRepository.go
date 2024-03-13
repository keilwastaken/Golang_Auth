package Repository

import (
	"Clarity_go/Interfaces"
	"Clarity_go/Models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type UsersRepository struct {
	mongodb Interfaces.IMongoService
}

func NewUsersRepository(pMongoDb Interfaces.IMongoService) *UsersRepository {
	return &UsersRepository{
		mongodb: pMongoDb,
	}
}

func (ur UsersRepository) RegisterUser(pUserRegisterDto Models.UserRegisterDto) (*mongo.InsertOneResult, *Models.ResponseError) {

	// Create a new User model from the UserRegisterDto
	newUser := Models.User{
		Id:       primitive.NewObjectID(), // Generate a new ObjectID
		Email:    pUserRegisterDto.Email,  // Assuming you want to use Username as Email
		Password: pUserRegisterDto.Password,
	}
	user, err := ur.isUserExists(newUser.Email)
	if err != nil {
		return nil, &Models.ResponseError{
			Message: "Failed to insert user",
			Status:  http.StatusInternalServerError,
		}
	}

	if user {
		return nil, &Models.ResponseError{
			Message: "User already has an account",
			Status:  http.StatusInternalServerError,
		}
	}

	// Insert the user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Correctly access the user collection from the mongodb service
	userCollection := ur.mongodb.GetUserCollection()
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, &Models.ResponseError{
			Message: "Failed to insert user",
			Status:  http.StatusInternalServerError,
		}
	}

	return result, nil
}

func (ur UsersRepository) Login(pUserRegisterDto Models.UserRegisterDto) (*Models.User, *Models.ResponseError) {
	xUser := Models.User{
		Email:    pUserRegisterDto.Email,
		Password: pUserRegisterDto.Password,
	}

	userCollection := ur.mongodb.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": xUser.Email}).Decode(&xUser)
	if err != nil {
		return nil,
			&Models.ResponseError{
				Message: "Email or password is incorrect",
				Status:  http.StatusUnauthorized,
			}
	}
	err = bcrypt.CompareHashAndPassword([]byte(xUser.Password), []byte(pUserRegisterDto.Password))
	if err != nil {
		return nil,
			&Models.ResponseError{
				Message: "Email or password is incorrect",
				Status:  http.StatusUnauthorized,
			}
	}

	return &xUser, nil
}

func (ur UsersRepository) isUserExists(email string) (bool, error) {
	collection := ur.mongodb.GetUserCollection()
	count, err := collection.CountDocuments(context.TODO(), bson.M{"email": strings.ToLower(email)})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
