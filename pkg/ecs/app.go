package ecs

import (
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger  *zap.SugaredLogger
	World   World
	Systems map[string]System
	tick    int
}

func (a *App) Init() {
	a.Logger.Info("App started")
}

func (a *App) Run(time *time.Time) {
	if time == nil {
		a.Logger.Warn("time is nil")
	}

	// for _, system := range a.Systems {
	// 	var requiredComponents = system.RequiredComponents()
	// 	var entities = a.World.FilterEntities(requiredComponents)
	// 	system.Run(&a.World, entities)
	// }

	a.tick += 1
}

func (a *App) Shutdown() {
	a.Logger.Info("App stopped")
}

func (a *App) ConsoleInput(input string) {
	a.Logger.Info(input)
}
