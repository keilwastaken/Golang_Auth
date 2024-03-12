package Services

import "Clarity_go/Repository"

type TaskService struct {
	taskRepository *Repository.TaskRepository
}

func NewTaskService(pTaskRepository *Repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: pTaskRepository,
	}
}
