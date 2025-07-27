package Controller

import (
	"goauth/Services"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskService *Services.TaskService
}

func NewTaskController(pTaskService *Services.TaskService) *TaskController {
	return &TaskController{
		TaskService: pTaskService,
	}
}

func (uc TaskController) GetAll(ctx *gin.Context) {

	println("GetAll Tasks hit")
}
