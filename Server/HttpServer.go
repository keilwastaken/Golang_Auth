package Server

import (
	"Clarity_go/Config"
	"Clarity_go/Controller"
	"Clarity_go/Interfaces"
	"Clarity_go/MiddleWare"
	"Clarity_go/Routes"
	"github.com/gin-gonic/gin"
	"log"
)

type HttpServer struct {
	config          *Config.Config
	router          *gin.Engine
	usersController *Controller.UsersController
	taskController  *Controller.TaskController
}

func InitHttpServer(pConfig *Config.Config, pMongoDb Interfaces.IMongoService, pToken Interfaces.IToken) HttpServer {
	xRouter := gin.Default()
	xAuthenticationMiddleware := MiddleWare.Authentication{}
	xUserRoutes := Routes.NewUserRoutes(pMongoDb, pToken, xAuthenticationMiddleware, xRouter)
	xTaskRoutes := Routes.NewTaskRoutes(pMongoDb, xRouter)

	return HttpServer{
		config:          pConfig,
		router:          xRouter,
		usersController: xUserRoutes.UserController,
		taskController:  xTaskRoutes.TaskController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.ServerAddress)
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
