package Controller

import (
	"Clarity_go/Models"
	"Clarity_go/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type UsersController struct {
	usersService *Services.UsersService
}

func NewUsersController(usersService *Services.UsersService) *UsersController {
	return &UsersController{
		usersService: usersService,
	}
}

func (uc UsersController) Register(ctx *gin.Context) {

	var xRegisterDto Models.UserRegisterDto

	if err := ctx.ShouldBindJSON(&xRegisterDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerResult, err := uc.usersService.RegisterUser(xRegisterDto)
	if err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	objectID, ok := registerResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the error if the type assertion fails
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	xAccessToken, xRefreshToken, err := uc.usersService.GenerateBothTokens(objectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Message})
		return

	}

	if _, err := uc.usersService.AddRefreshTokenToDb(objectID, xRefreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully.", "token": xAccessToken, "refreshToken": xRefreshToken})
}

func (uc UsersController) Login(ctx *gin.Context) {
	var xRegisterDto Models.UserRegisterDto

	if err := ctx.ShouldBindJSON(&xRegisterDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Look up for requested user
	xUser, err := uc.usersService.Login(xRegisterDto)
	if err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	xAccessToken, xRefreshToken, err := uc.usersService.GenerateBothTokens(xUser.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Message})
		return
	}

	if _, err := uc.usersService.DeleteRefreshTokenFromDb(*xRefreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove token from db"})
		return
	}

	if _, err := uc.usersService.AddRefreshTokenToDb(xUser.Id, xRefreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully.", "accessToken": xAccessToken, "refreshToken": xRefreshToken})
}

type TokenRequest struct {
	RefreshToken string `json:"refreshToken"` // Field names must be capitalized for the JSON decoder to see them
}

func (uc UsersController) RefreshTokens(ctx *gin.Context) {
	var xTokenRequest TokenRequest

	if err := ctx.ShouldBindJSON(&xTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Use the refresh token from the parsed JSON
	refreshToken := xTokenRequest.RefreshToken
	if refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}
	xUserId, err := uc.usersService.GetUserIdByToken(refreshToken)
	if err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	xAccessToken, xRefreshToken, err := uc.usersService.GenerateBothTokens(xUserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Message})
		return
	}

	if _, err := uc.usersService.DeleteRefreshTokenFromDb(refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove token from db"})
		return
	}

	if _, err := uc.usersService.AddRefreshTokenToDb(xUserId, xRefreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Token refresh successfully.", "accessToken": xAccessToken, "refreshToken": xRefreshToken})
}

func (uc UsersController) Logout(ctx *gin.Context) {

	refreshToken := ctx.Request.Header.Get("Authorization")
	strippedToken := strings.TrimPrefix(refreshToken, "Bearer ")

	if _, err := uc.usersService.DeleteRefreshTokenFromDb(strippedToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully."})
}

func (uc UsersController) Validate(c *gin.Context) {
	user, _ := c.Get("user")

	// user.(models.User).Email    -->   to access specific data

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
