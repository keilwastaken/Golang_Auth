package Services

import (
	"Clarity_go/Models"
	"Clarity_go/Repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"html"
	"net/http"
	"strings"
)

type UsersService struct {
	usersRepository *Repository.UsersRepository
}

func NewUsersService(usersRepository *Repository.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us UsersService) RegisterUser(pUserRegisterDto Models.UserRegisterDto) (*mongo.InsertOneResult, *Models.ResponseError) {
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pUserRegisterDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil,
			&Models.ResponseError{
				Message: "Failed to hash password",
				Status:  http.StatusInternalServerError,
			}
	}
	pUserRegisterDto.Password = string(hashedPassword)

	//TODO CHECK IF THE USER IS ALREADY IN THE DB
	//TODO CHECK IF THE EMAIL IS VALID
	//TODO CHECK IF THE PASSWORD IS VALID

	//remove spaces in username
	pUserRegisterDto.Username = html.EscapeString(strings.TrimSpace(pUserRegisterDto.Username))

	result, Error := us.usersRepository.RegisterUser(pUserRegisterDto)
	if Error != nil {
		return nil,
			&Models.ResponseError{
				Message: "Failed to insert user",
				Status:  http.StatusInternalServerError,
			}
	}
	return result, nil
}

func (us UsersService) Login(ctx *gin.Context, UserRegisterDto Models.UserRegisterDto) (*Models.User, *Models.ResponseError) {
	xUser, responseErr := us.usersRepository.Login(ctx, UserRegisterDto)
	if responseErr != nil {
		return nil, responseErr
	}

	if xUser.Id == primitive.NilObjectID {
		return nil,
			&Models.ResponseError{
				Message: "Email or password is incorrect",
				Status:  http.StatusUnauthorized,
			}
	}

	return xUser, nil
}
