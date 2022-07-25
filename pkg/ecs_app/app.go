package ecs

import (
	"log"

	inecs "github.com/fishykins/aero/internal/inecs"
	core "github.com/fishykins/aero/pkg/ecs_core"
)

type App struct {
	world   *inecs.World
	manager *core.Manager
}

func New() *App {
	return &App{
		world:   inecs.NewWorld(),
		manager: core.NewManager(),
	}
}

func (a *App) AddEntity(args ...string) *core.EntityBuilder {
	return a.manager.AddEntity(args...)
}

func (a *App) AddSystem(system core.System, args ...core.Query) *core.SystemBuilder {
	return a.manager.AddSystem(system, args...)
}

func (a *App) Run() {
	// TODO
}

func (a *App) Update() {
	log.Println("====================== UPDATE ======================")
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
			log.Println("Query finished:", query.ID)
		case system := <-systemChan:
			finishedSystems[system] = true
			log.Println("System finished:", system)
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
