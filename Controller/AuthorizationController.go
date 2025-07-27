package Controller

import (
	"goauth/Models"
	"goauth/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type AuthorizationController struct {
	AuthorizationService *Services.AuthorizationService
}

func NewAuthorizationController(usersService *Services.AuthorizationService) *AuthorizationController {
	return &AuthorizationController{
		AuthorizationService: usersService,
	}
}

func (ac AuthorizationController) Register(ctx *gin.Context) {

	var xRegisterDto Models.UserRegisterDto

	if err := ctx.ShouldBindJSON(&xRegisterDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerResult, err := ac.AuthorizationService.RegisterUser(xRegisterDto)
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

	xAccessToken, xRefreshToken, err := ac.AuthorizationService.GenerateBothTokens(objectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Message})
		return

	}

	if _, err := ac.AuthorizationService.AddRefreshTokenToDb(objectID, xRefreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully.", "token": xAccessToken, "refreshToken": xRefreshToken})
}

func (ac AuthorizationController) Login(ctx *gin.Context) {
	var xRegisterDto Models.UserRegisterDto

	if err := ctx.ShouldBindJSON(&xRegisterDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Look up for requested user
	xUser, err := ac.AuthorizationService.Login(xRegisterDto)
	if err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	xAccessToken, xRefreshToken, err := ac.AuthorizationService.GenerateBothTokens(xUser.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Message})
		return
	}

	if _, err := ac.AuthorizationService.DeleteRefreshTokenByUserId(xUser.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove token from db"})
		return
	}

	if _, err := ac.AuthorizationService.AddRefreshTokenToDb(xUser.Id, xRefreshToken); err != nil {
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

func (ac AuthorizationController) RefreshTokens(ctx *gin.Context) {
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
	xUserId, err := ac.AuthorizationService.GetUserIdByToken(refreshToken)
	if err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	xAccessToken, xRefreshToken, err := ac.AuthorizationService.GenerateBothTokens(xUserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Message})
		return
	}

	if _, err := ac.AuthorizationService.DeleteRefreshTokenFromDb(refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove token from db"})
		return
	}

	if _, err := ac.AuthorizationService.AddRefreshTokenToDb(xUserId, xRefreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Token refresh successfully.", "accessToken": xAccessToken, "refreshToken": xRefreshToken})
}

func (ac AuthorizationController) Logout(ctx *gin.Context) {

	refreshToken := ctx.Request.Header.Get("Authorization")

	if _, err := ac.AuthorizationService.DeleteRefreshTokenFromDb(refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully."})
}

func (ac AuthorizationController) Validate(c *gin.Context) {
	user, _ := c.Get("user")

	// user.(models.User).Email    -->   to access specific data

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
