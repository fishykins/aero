package ecs

import (
	"os"
	"os/signal"
	"time"

	inecs "github.com/fishykins/aero/internal/inecs"
	core "github.com/fishykins/aero/pkg/ecs_core"
	log "github.com/fishykins/aero/pkg/logging"
)

var Log log.Logger

type App struct {
	// World stores critial ECS data. It should never be exposed to the end-user!
	world *inecs.World
	// The world manager is an endpoint for the user, and can be used to add/remove data from the world in a controlled manner.
	manager *core.WorldManager
}

func New() *App {
	return &App{
		world:   inecs.NewWorld(),
		manager: core.NewWorldManager(),
	}
}

func (a *App) AddEntity(args ...string) *core.EntityBuilder {
	return a.manager.AddEntity(args...)
}

func (a *App) AddComponent(entity *core.Entity, component core.Component) {
	a.manager.AddComponent(entity, component)
}

func (a *App) AddSystem(system core.System, args ...core.Query) *core.SystemBuilder {
	return a.manager.AddSystem(system, args...)
}

func (a *App) AddResource(resource core.Component) {
	a.manager.AddResource(resource.Type(), resource)
}

func (a *App) Run() {
	var ticker *time.Ticker

	if fps, err := a.manager.GetResource("UpdateFrequency"); err != nil {
		ticker = time.NewTicker(time.Second)
	} else {
		ticker = time.NewTicker(fps.(core.UpdateFrequency).FPS())
	}

	exit := make(chan bool)
	c := make(chan os.Signal)
	log.Info("Starting ECS App")

	// Catch Ctrl-C command
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			exit <- true
			return
		}
	}()

	// Run the main loop
	go func() {
		for {
			select {
			case <-ticker.C:
				a.Update()
			case <-exit:
				ticker.Stop()
				return
			}
		}
	}()

	// Wait for exit
	<-exit
	a.Shutdown()
}

func (a *App) Update() {
	// Checks the manager for new entities, components and systems.
	a.world.Manage(a.manager)

	// Empty maps of finished queries and systems.
	queries := make(map[uint32]core.QueryResult)
	pendingSystems := make(map[string]bool, 0)
	finishedSystems := make(map[string]bool, 0)
	for id := range a.world.Systems {
		pendingSystems[id] = true
	}

	// Channels for communicating between routines.
	queryChan := make(chan core.QueryResult)
	systemChan := make(chan string)

	// Kick off all the routines.
	a.world.BuildQueries(queryChan)

	// Await responses.
	for i := 0; i < len(a.world.Queries)+len(a.world.Systems); i++ {
		// Blocking call to wait for either a query or a system to finish.
		select {
		case query := <-queryChan:
			queries[query.ID] = query
		case system := <-systemChan:
			finishedSystems[system] = true
		}

		// Check if any pending systems can start
		for system, pending := range pendingSystems {
			if pending {
				requiredQueries := a.world.Systems[system].Queries
				queryScore := 0
				for _, query := range requiredQueries {
					if _, ok := queries[query]; ok {
						queryScore++
					}
				}
				if queryScore == len(requiredQueries) {
					requiredSystems := a.world.Systems[system].RunsAfter
					systemScore := 0
					for _, system := range requiredSystems {
						if _, ok := finishedSystems[system]; ok {
							systemScore++
						}
					}
					if systemScore == len(requiredSystems) {
						// Run the system!
						pendingSystems[system] = false
						systemQueries := make([]core.QueryResult, 0)
						for _, query := range requiredQueries {
							systemQueries = append(systemQueries, queries[query])
						}
						go a.world.RunSystem(system, systemChan, a.manager, systemQueries...)
					}
				}
			}
		}
	}
}

func (a *App) Shutdown() {
	log.Warn("Shutting down ECS App")
}
