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

func NewUserRoutes(pMongoDb Interfaces.IMongoService, pToken Interfaces.IToken, pRoutes *gin.Engine) *UserRoutes {
	xUsersRepository := Repository.NewUsersRepository(pMongoDb)
	xUsersService := Services.NewUsersService(xUsersRepository)
	xUsersController := Controller.NewUsersController(xUsersService, pToken)

	xRouter := pRoutes
	xUserGroup := xRouter.Group("/user")
	{
		xUserGroup.POST("/register", xUsersController.Register)
		xUserGroup.POST("/login", xUsersController.Login)
		xUserGroup.POST("/logout", xUsersController.Logout)
	}

	return &UserRoutes{
		mongodb:        pMongoDb,
		token:          pToken,
		routes:         xRouter,
		UserController: xUsersController,
	}
}
