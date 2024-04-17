package Routes

import (
	"Clarity_go/Controller"
	"Clarity_go/Interfaces"
	"Clarity_go/Repository"
	"Clarity_go/Services"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	mongodb        Interfaces.IMongoService
	token          Interfaces.IToken
	routes         *gin.Engine
	UserController *Controller.UsersController
}

func NewUserRoutes(pMongoDb Interfaces.IMongoService, pToken Interfaces.IToken, pAuthenticationMiddleware Interfaces.IAuthenticationMiddleware, pRoutes *gin.Engine) *UserRoutes {
	xUsersRepository := Repository.NewUsersRepository(pMongoDb, pToken)
	xUsersService := Services.NewUsersService(xUsersRepository)
	xUsersController := Controller.NewUsersController(xUsersService)

	xRouter := pRoutes
	xUserGroup := xRouter.Group("/user")
	{
		xUserGroup.POST("/register", xUsersController.Register)
		xUserGroup.POST("/login", xUsersController.Login)
		xUserGroup.POST("/refresh", xUsersController.RefreshTokens)
		xUserGroup.POST("/logout", xUsersController.Logout)
		xUserGroup.GET("/validate", pAuthenticationMiddleware.RequireAuth(pToken), xUsersController.Validate)
	}

	return &UserRoutes{
		mongodb:        pMongoDb,
		token:          pToken,
		routes:         xRouter,
		UserController: xUsersController,
	}
}
