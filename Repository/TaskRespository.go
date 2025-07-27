package Repository

import (
	"goauth/Interfaces"
)

type TaskRepository struct {
	mongodb Interfaces.IMongoService
}

func NewTaskRepository(pMongoDb Interfaces.IMongoService) *TaskRepository {
	return &TaskRepository{
		mongodb: pMongoDb,
	}
}
