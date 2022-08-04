package app

import (
	"os"
	"os/signal"
	"time"

	ecs "github.com/fishykins/aero/pkg/ecs_core"
	log "github.com/fishykins/aero/pkg/logging"
)

type App struct {
	// World stores critial ECS data. It should never be exposed to the end-user!
	world *ecs.World
	// The world manager is an endpoint for the user, and can be used to add/remove data from the world in a controlled manner.
	manager *ecs.WorldManager
}

func New() *App {
	return &App{
		world:   ecs.NewWorld(),
		manager: ecs.NewWorldManager(),
	}
}

func (a *App) AddEntity() *ecs.EntityBuilder {
	return a.manager.NewEntity()
}

func (a *App) AddComponent(entity ecs.Entity, component interface{}) {
	a.manager.AddComponent(entity, component)
}

func (a *App) AddSystem(system ecs.SystemFunc) *ecs.SystemBuilder {
	return a.manager.AddSystem(system)
}

func (a *App) AddResource(resource interface{}) {
	a.manager.AddResource(resource)
}

func (a *App) Run() {
	var ticker *time.Ticker

	if fps, err := a.world.GetResource("UpdateFrequency"); err != nil {
		ticker = time.NewTicker(time.Second)
	} else {
		ticker = time.NewTicker(fps.(*ecs.UpdateFrequency).FPS())
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
	a.manager.UpdateWorld(a.world, true)

	// Empty maps of finished queries and systems.
	queries := make(map[uint32]ecs.QueryResult)
	pendingSystems := make(map[string]bool, 0)
	finishedSystems := make(map[string]bool, 0)
	for id := range a.world.Systems {
		pendingSystems[id] = true
	}

	// Channels for communicating between routines.
	queryChan := make(chan ecs.QueryResult)
	systemChan := make(chan string)

	// Kick off all the routines.
	for id, query := range a.world.Queries {
		go a.world.BuildQuery(id, query, queryChan)
	}

	// Await responses.
	for i := 0; i < len(a.world.Queries)+len(a.world.Systems); i++ {
		// Blocking call to wait for either a query or a system to finish.
		select {
		case query := <-queryChan:
			queries[query.ID] = query
			log.TraceWith("Query: finished", map[string]interface{}{"id": query.ID, "result": query.Result})
		case system := <-systemChan:
			finishedSystems[system] = true
		}

		// Check if any pending systems can start
		for s, pending := range pendingSystems {
			system := a.world.Systems[s]
			if pending {
				requiredQueries := system.GetQueries()
				queryScore := 0
				for _, query := range requiredQueries {
					if _, ok := queries[query]; ok {
						queryScore++
					}
				}
				if queryScore == len(requiredQueries) {
					requiredSystems := system.SystemsBefore()
					systemScore := 0
					for _, system := range requiredSystems {
						if _, ok := finishedSystems[system]; ok {
							systemScore++
						}
					}
					if systemScore == len(requiredSystems) {
						// The system is ready to run- grab resources and lets go!
						resources := make(ecs.RMap)
						for _, resourceId := range system.GetResources() {
							r, err := a.world.GetResource(resourceId)
							if err != nil {
								log.ErrorWith("Failed to get resource", map[string]interface{}{"id": resourceId})
							}
							resources[resourceId] = r
						}
						// Run the system!
						pendingSystems[s] = false
						systemQueries := make([]ecs.QueryResult, 0)
						for _, query := range requiredQueries {
							systemQueries = append(systemQueries, queries[query])
						}
						go a.world.RunSystem(s, systemChan, a.manager, resources, systemQueries...)
					}
				}
			}
		}
	}
}

func (a *App) Shutdown() {
	log.Warn("Shutting down ECS App")
}
