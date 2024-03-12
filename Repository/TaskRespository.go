package Repository

import (
	"Clarity_go/Interfaces"
)

type TaskRepository struct {
	mongodb Interfaces.IMongoService
}

func NewTaskRepository(pMongoDb Interfaces.IMongoService) *TaskRepository {
	return &TaskRepository{
		mongodb: pMongoDb,
	}
}
