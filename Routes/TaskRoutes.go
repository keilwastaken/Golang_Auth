package Routes

import (
	"Clarity_go/Controller"
	"Clarity_go/Interfaces"
	"Clarity_go/Repository"
	"Clarity_go/Services"
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
