package Server

import (
	"goauth/Config"
	"goauth/Controller"
	"goauth/Interfaces"
	"goauth/MiddleWare"
	"goauth/Routes"
	"github.com/gin-gonic/gin"
	"log"
)

type HttpServer struct {
	config                  *Config.Config
	router                  *gin.Engine
	AuthorizationController *Controller.AuthorizationController
	taskController          *Controller.TaskController
}

func InitHttpServer(pConfig *Config.Config, pMongoDb Interfaces.IMongoService, pToken Interfaces.IToken) HttpServer {
	xRouter := gin.Default()
	xAuthenticationMiddleware := MiddleWare.Authentication{}
	xUserRoutes := Routes.NewAuthorizationRoutes(pMongoDb, pToken, xAuthenticationMiddleware, xRouter)
	xTaskRoutes := Routes.NewTaskRoutes(pMongoDb, xRouter)

	return HttpServer{
		config:                  pConfig,
		router:                  xRouter,
		AuthorizationController: xUserRoutes.AuthorizationController,
		taskController:          xTaskRoutes.TaskController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.ServerAddress)
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
