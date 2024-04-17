package Routes

import (
	"Clarity_go/Controller"
	"Clarity_go/Interfaces"
	"Clarity_go/Repository"
	"Clarity_go/Services"
	"github.com/gin-gonic/gin"
)

type AuthorizationRoutes struct {
	mongodb                 Interfaces.IMongoService
	token                   Interfaces.IToken
	routes                  *gin.Engine
	AuthorizationController *Controller.AuthorizationController
}

func NewAuthorizationRoutes(pMongoDb Interfaces.IMongoService, pToken Interfaces.IToken, pAuthenticationMiddleware Interfaces.IAuthenticationMiddleware, pRoutes *gin.Engine) *AuthorizationRoutes {
	xAuthorizationRepository := Repository.NewAuthorizationRepository(pMongoDb, pToken)
	xAuthorizationService := Services.NewAuthorizationService(xAuthorizationRepository)
	xAuthorizationController := Controller.NewAuthorizationController(xAuthorizationService)

	xRouter := pRoutes
	xAuthorizationGroup := xRouter.Group("/auth")
	{
		xAuthorizationGroup.POST("/register", xAuthorizationController.Register)
		xAuthorizationGroup.POST("/login", xAuthorizationController.Login)
		xAuthorizationGroup.POST("/refresh", xAuthorizationController.RefreshTokens)
		xAuthorizationGroup.POST("/logout", xAuthorizationController.Logout)
		xAuthorizationGroup.GET("/validate", pAuthenticationMiddleware.RequireAuth(pToken), xAuthorizationController.Validate)
	}

	return &AuthorizationRoutes{
		mongodb:                 pMongoDb,
		token:                   pToken,
		routes:                  xRouter,
		AuthorizationController: xAuthorizationController,
	}
}
