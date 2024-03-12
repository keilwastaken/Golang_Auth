package Controller

import (
	"Clarity_go/Interfaces"
	"Clarity_go/Models"
	"Clarity_go/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
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
		// Log the error without terminating the program
		log.Printf("Failed to insert user: %v", err)

		// Respond with an appropriate error message
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
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

	token, responseErr := uc.IToken.GenerateToken(objectID)
	if responseErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": token})
}

func (uc UsersController) Login(ctx *gin.Context) {
	var xRegisterDto Models.UserRegisterDto

	if err := ctx.ShouldBindJSON(&xRegisterDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Look up for requested user
	xUser, err := uc.usersService.Login(ctx, xRegisterDto)
	if err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	token, responseErr := uc.IToken.GenerateToken(xUser.Id)
	if responseErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})
}

func (uc UsersController) Logout(ctx *gin.Context) {

	println("User logout route hit")
}
