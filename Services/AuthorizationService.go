package Services

import (
	"Clarity_go/Models"
	"Clarity_go/Repository"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

type AuthorizationService struct {
	AuthorizationRepository *Repository.AuthorizationRepository
}

func NewAuthorizationService(AuthorizationRepository *Repository.AuthorizationRepository) *AuthorizationService {
	return &AuthorizationService{
		AuthorizationRepository: AuthorizationRepository,
	}
}

func (as AuthorizationService) RegisterUser(pUserRegisterDto Models.UserRegisterDto) (*mongo.InsertOneResult, *Models.ResponseError) {

	if !isValidEmail(pUserRegisterDto.Email) {
		return nil,
			&Models.ResponseError{
				Message: "Invalid email",
				Status:  http.StatusBadRequest,
			}
	}

	valid, errMsg := isValidPassword(pUserRegisterDto.Password)
	if !valid {
		return nil,
			&Models.ResponseError{
				Message: errMsg,
				Status:  http.StatusBadRequest,
			}
	}

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

	result, Error := as.AuthorizationRepository.RegisterUser(pUserRegisterDto)
	if Error != nil {
		return nil,
			&Models.ResponseError{
				Message: Error.Message,
				Status:  http.StatusInternalServerError,
			}
	}
	return result, nil
}

func (as AuthorizationService) Login(UserRegisterDto Models.UserRegisterDto) (*Models.User, *Models.ResponseError) {
	xUser, responseErr := as.AuthorizationRepository.Login(UserRegisterDto)
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

func (as AuthorizationService) GenerateBothTokens(userId primitive.ObjectID) (*string, *string, *Models.ResponseError) {
	accessToken, err := as.AuthorizationRepository.GenerateAccessToken(userId)
	if err != nil {
		return nil, nil,
			&Models.ResponseError{
				Message: err.Message,
				Status:  http.StatusInternalServerError,
			}
	}

	refreshToken, err := as.AuthorizationRepository.GenerateRefreshToken(userId)
	if err != nil {
		return nil, nil,
			&Models.ResponseError{
				Message: err.Message,
				Status:  http.StatusInternalServerError,
			}
	}

	return accessToken, refreshToken, nil
}

func (as AuthorizationService) GetUserIdByToken(ptoken string) (primitive.ObjectID, *Models.ResponseError) {

	userId, err := as.AuthorizationRepository.GetUserIdByToken(ptoken)
	if err != nil {
		return primitive.NilObjectID,
			&Models.ResponseError{
				Message: err.Message,
				Status:  http.StatusInternalServerError,
			}
	}

	return userId, nil
}

func (as AuthorizationService) AddRefreshTokenToDb(userId primitive.ObjectID, refreshToken *string) (*mongo.InsertOneResult, *Models.ResponseError) {
	result, Error := as.AuthorizationRepository.AddRefreshTokenToDb(userId, *refreshToken)
	if Error != nil {
		return nil,
			&Models.ResponseError{
				Message: Error.Message,
				Status:  http.StatusInternalServerError,
			}
	}
	return result, nil
}

func (as AuthorizationService) DeleteRefreshTokenFromDb(pRefreshToken string) (*mongo.DeleteResult, *Models.ResponseError) {
	strippedToken := strings.TrimPrefix(pRefreshToken, "Bearer ")

	result, Error := as.AuthorizationRepository.DeleteRefreshToken(strippedToken)
	if Error != nil {
		return nil,
			&Models.ResponseError{
				Message: Error.Message,
				Status:  http.StatusInternalServerError,
			}
	}
	return result, nil
}

func (as AuthorizationService) DeleteRefreshTokenByUserId(pUserid primitive.ObjectID) (*mongo.DeleteResult, *Models.ResponseError) {
	result, Error := as.AuthorizationRepository.DeleteRefreshTokenByUserId(pUserid)
	if Error != nil {
		return nil,
			&Models.ResponseError{
				Message: Error.Message,
				Status:  http.StatusInternalServerError,
			}
	}
	return result, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(password string) (bool, string) {
	var (
		minLen      = 8
		upperCase   = regexp.MustCompile(`[A-Z]`)
		lowerCase   = regexp.MustCompile(`[a-z]`)
		number      = regexp.MustCompile(`[0-9]`)
		specialChar = regexp.MustCompile(`[^a-zA-Z0-9]`)
	)

	if len(password) < minLen {
		return false, fmt.Sprintf("Password must be at least %d characters long.", minLen)
	}
	if !upperCase.MatchString(password) {
		return false, "Password must contain at least one uppercase letter."
	}
	if !lowerCase.MatchString(password) {
		return false, "Password must contain at least one lowercase letter."
	}
	if !number.MatchString(password) {
		return false, "Password must contain at least one digit."
	}
	if !specialChar.MatchString(password) {
		return false, "Password must contain at least one special character."
	}

	return true, ""
}
