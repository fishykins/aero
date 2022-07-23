package ecs

import (
	"os"
	"reflect"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger  *zap.SugaredLogger
	World   World
	Systems map[reflect.Type]bool
	tick    int
}

type AppBuilder struct {
	App
}

type System interface{}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		App: App{
			World:   NewWorld(),
			Systems: make(map[reflect.Type]bool),
			tick:    0,
		},
	}
}

func (ab *AppBuilder) WithSystem(name string, system System) *AppBuilder {
	x := reflect.TypeOf(system)
	numIn := x.NumIn()
	numOut := x.NumOut()

	if numOut != 0 {
		ab.App.Logger.Error("System", name, "should not have any return values")
		os.Exit(21)
	}

	validInputs := true

	for i := 0; i < numIn; i++ {
		input := x.In(i)
		if !input.Implements(reflect.TypeOf((*Component)(nil)).Elem()) && !(reflect.TypeOf(input) == reflect.TypeOf((*Entity)(nil))) {
			ab.App.Logger.Warn("System", name, "has an invalid input type at index ", i, ": ", input)
			validInputs = false
			i = numIn
		}
	}

	if validInputs {
		ab.App.Systems[x] = true
		ab.App.Logger.Info("System '", name, "' added")
	}
	return ab
}

func (ab *AppBuilder) WithLogger(logger *zap.SugaredLogger) *AppBuilder {
	ab.App.Logger = logger
	return ab
}

func (ab *AppBuilder) Build() *App {
	return &ab.App
}

func (a *App) Run(time *time.Time) {
	if time == nil {
		a.Logger.Warn("time is nil")
	}

	for systemType, running := range a.Systems {
		if !running {
			continue
		}
		numIn := systemType.NumIn()
		for i := 0; i < numIn; i++ {
			compType := systemType.In(i)
			// We know it is valid as we checked when we added the system.
			if _, ok := a.World.Components[compType]; !ok {
				// No components of this type- return an empty slice.
			}
		}
		//newSystemCall := reflect.ValueOf(systemType)
		// Cant cast back- HOW TO DO THIS???? Grrr
	}
	a.tick += 1
}

func (a *App) Shutdown() {
	a.Logger.Info("App stopped")
}

func (a *App) ConsoleInput(input string) {
	a.Logger.Info(input)
}
