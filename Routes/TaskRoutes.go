package Routes

import (
	"goauth/Controller"
	"goauth/Interfaces"
	"goauth/Repository"
	"goauth/Services"
	"github.com/gin-gonic/gin"
)

type TaskRoutes struct {
	mongoDb        Interfaces.IMongoService
	routes         *gin.Engine
	TaskController *Controller.TaskController
}

func NewTaskRoutes(pMongoDb Interfaces.IMongoService, pRoutes *gin.Engine) *TaskRoutes {
	xTaskRepository := Repository.NewTaskRepository(pMongoDb)
	xTaskService := Services.NewTaskService(xTaskRepository)
	xTaskController := Controller.NewTaskController(xTaskService)

	xRouter := pRoutes

	xRouter.GET("/task/all", xTaskController.GetAll)

	return &TaskRoutes{
		mongoDb:        pMongoDb,
		routes:         xRouter,
		TaskController: xTaskController,
	}
}
