package backend

import (
	"context"
	"os"
	"path/filepath"

	"todo-wails/backend/repo"
	"todo-wails/backend/service"
	"todo-wails/backend/usecase"
)

type App struct {
	ctx context.Context
	UC  *usecase.TaskUsecase
}

func NewApp() *App { 
	return &App{} 
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	appDataDir, _ := os.UserConfigDir()
	appDataDir = filepath.Join(appDataDir, "todo-wails")
	r := repo.NewFileTaskRepo(appDataDir)
	svc, _ := service.NewTaskService(r)
	a.UC = usecase.NewTaskUsecase(svc)
}

func (a *App) domReady(ctx context.Context) {
	// Optional: Add any DOM ready logic here
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// Optional: Add any cleanup logic here
	return false
}

func (a *App) shutdown(ctx context.Context) {
	// Optional: Add any shutdown logic here
}
