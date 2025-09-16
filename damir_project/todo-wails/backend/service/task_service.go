package service

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"todo-wails/backend/entity"
	"todo-wails/backend/repo"
)

type Filter string
const (
	FilterAll   Filter = "all"
	FilterActive Filter = "active"
	FilterDone   Filter = "done"
)

type Sort string
const (
	SortCreatedAt Sort = "createdAt"
	SortPriority  Sort = "priority"
	SortDueAt     Sort = "dueAt"
)

type TaskService struct {
	repo  repo.TaskRepo
	cache []entity.Task
}

func NewTaskService(r repo.TaskRepo) (*TaskService, error) {
	s := &TaskService{repo: r}
	tasks, err := r.Load()
	if err != nil {
		return nil, err
	}
	s.cache = tasks
	return s, nil
}

func (s *TaskService) Add(title string, due *time.Time, p entity.Priority) (entity.Task, error) {
	t := entity.Task{
		ID:        uuid.NewString(),
		Title:     title,
		CreatedAt: time.Now(),
		DueAt:     due,
		Priority:  p,
		Done:      false,
	}
	s.cache = append(s.cache, t)
	return t, s.repo.Save(s.cache)
}

func (s *TaskService) ToggleDone(id string) error {
	for i := range s.cache {
		if s.cache[i].ID == id {
			s.cache[i].Done = !s.cache[i].Done
			return s.repo.Save(s.cache)
		}
	}
	return nil
}

func (s *TaskService) Delete(id string) error {
	out := s.cache[:0]
	for _, t := range s.cache {
		if t.ID != id {
			out = append(out, t)
		}
	}
	s.cache = out
	return s.repo.Save(s.cache)
}

func (s *TaskService) List(filter Filter, sortBy Sort) []entity.Task {
	var res []entity.Task
	for _, t := range s.cache {
		switch filter {
		case FilterActive:
			if !t.Done { res = append(res, t) }
		case FilterDone:
			if t.Done { res = append(res, t) }
		default:
			res = append(res, t)
		}
	}

	switch sortBy {
	case SortPriority:
		order := map[entity.Priority]int{entity.PriorityHigh:0, entity.PriorityMedium:1, entity.PriorityLow:2}
		sort.Slice(res, func(i, j int) bool {
			return order[res[i].Priority] < order[res[j].Priority]
		})
	case SortDueAt:
		sort.Slice(res, func(i, j int) bool {
			var ti, tj time.Time
			if res[i].DueAt != nil { ti = *res[i].DueAt }
			if res[j].DueAt != nil { tj = *res[j].DueAt }
			return ti.Before(tj)
		})
	default: // createdAt
		sort.Slice(res, func(i, j int) bool {
			return res[i].CreatedAt.Before(res[j].CreatedAt)
		})
	}
	return res
}
