package Controller

import (
	"Clarity_go/Interfaces"
	"Clarity_go/Models"
	"Clarity_go/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type UsersController struct {
	usersService *Services.UsersService
	IToken       Interfaces.IToken
}

func NewUsersController(usersService *Services.UsersService, pToken Interfaces.IToken) *UsersController {
	return &UsersController{
		usersService: usersService,
		IToken:       pToken,
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

	xAccessToken, xError := uc.IToken.GenerateAccessToken(objectID)
	if xError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	xRefreshToken, xError := uc.IToken.GenerateRefreshToken(objectID)
	if xError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
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

	xAccessToken, xError := uc.IToken.GenerateAccessToken(xUser.Id)
	if xError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	//todo store refresh token in db
	xRefreshToken, xError := uc.IToken.GenerateRefreshToken(xUser.Id)
	if xError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
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
