package repo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"todo-wails/backend/entity"
)

type TaskRepo interface {
	Load() ([]entity.Task, error)
	Save([]entity.Task) error
}

type FileTaskRepo struct {
	path string
	mu   sync.Mutex
}

func NewFileTaskRepo(appDataDir string) *FileTaskRepo {
	_ = os.MkdirAll(appDataDir, 0o755)
	return &FileTaskRepo{path: filepath.Join(appDataDir, "tasks.json")}
}

func (r *FileTaskRepo) Load() ([]entity.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return []entity.Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []entity.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *FileTaskRepo) Save(tasks []entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}
