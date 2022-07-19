package internal

import (
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger    *zap.SugaredLogger
	tickCount int
}

func NewApp(logger *zap.SugaredLogger) *App {
	return &App{
		Logger: logger,
	}
}

func (a *App) Init() {
	a.Logger.Info("App started")
}

func (a *App) Run(time *time.Time) {
	if time == nil {
		a.Logger.Warn("time is nil")
	}

	a.tickCount += 1
}

func (a *App) Shutdown() {
	a.Logger.Info("App stopped")
}

func (a *App) ConsoleInput(input string) {
	a.Logger.Info(input)
}
