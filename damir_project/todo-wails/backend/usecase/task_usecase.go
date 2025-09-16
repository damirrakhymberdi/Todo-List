package usecase

import (
	"time"

	"todo-wails/backend/entity"
	"todo-wails/backend/service"
)

type TaskDTO = entity.Task

type TaskUsecase struct {
	svc *service.TaskService
}

func NewTaskUsecase(s *service.TaskService) *TaskUsecase { 
	return &TaskUsecase{svc: s} 
}

func (u *TaskUsecase) Add(title string, dueISO *string, priority string) (TaskDTO, error) {
	if title == "" {
		return TaskDTO{}, nil
	}
	var due *time.Time
	if dueISO != nil && *dueISO != "" {
		t, err := time.Parse(time.RFC3339, *dueISO)
		if err == nil { 
			due = &t 
		}
	}
	p := entity.Priority(priority)
	if p == "" { 
		p = entity.PriorityMedium 
	}
	return u.svc.Add(title, due, p)
}

func (u *TaskUsecase) ToggleDone(id string) error { 
	return u.svc.ToggleDone(id) 
}

func (u *TaskUsecase) Delete(id string) error { 
	return u.svc.Delete(id) 
}

func (u *TaskUsecase) List(filter string, sortBy string) []TaskDTO {
	return u.svc.List(service.Filter(filter), service.Sort(sortBy))
}
