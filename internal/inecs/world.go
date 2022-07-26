package inecs

import (
	"log"

	core "github.com/fishykins/aero/pkg/ecs_core"
	"github.com/fishykins/aero/pkg/slotmap"
)

// a World contains all the data corresponding to the ECS system.
type World struct {
	Entities   slotmap.SlotMap[string]
	Queries    map[uint32]core.Query
	Components map[string]map[core.Entity]core.Component
	Resources  map[string]core.Component
	Systems    map[string]SystemData
}

var defaultWorld *World

func NewWorld() *World {
	return &World{
		Entities:   slotmap.New[string]("entity"),
		Queries:    make(map[uint32]core.Query),
		Components: make(map[string]map[core.Entity]core.Component),
		Resources:  make(map[string]core.Component),
		Systems:    make(map[string]SystemData),
	}
}

func DefaultWorld() *World {
	if defaultWorld == nil {
		defaultWorld = NewWorld()
	}
	return defaultWorld
}

func (w *World) SetDefault() {
	defaultWorld = w
}

// Looks over the world manager and instantiates any pending entities/systems.
// This should be called once per frame, and probably at the end of the main loop.
// We can also push/pull resource data from the world manager.
func (w *World) Manage(manager *core.WorldManager) {
	w.manageEntities(manager)
	w.manageComponents(manager)
	w.manageSystems(manager)
	manager.PendingEntities = nil
	manager.PendingSystems = nil
	manager.PendingComponents = nil
	w.Resources = manager.Resources()
}

func (w *World) manageEntities(manager *core.WorldManager) {
	for _, eb := range manager.PendingEntities {
		id, components := eb.Build()
		entity := core.Entity(w.Entities.Add(id))

		for _, component := range components {
			if _, ok := w.Components[component.Type()]; !ok {
				w.Components[component.Type()] = make(map[core.Entity]core.Component)
			}
			w.Components[component.Type()][entity] = component
		}
	}
}

func (w *World) manageComponents(manager *core.WorldManager) {
	for _, componentType := range manager.PendingComponents {
		for entity, component := range componentType {
			if _, ok := w.Components[component.Type()]; !ok {
				w.Components[component.Type()] = make(map[core.Entity]core.Component)
			}
			w.Components[component.Type()][entity] = component
		}
	}
}

func (w *World) manageSystems(manager *core.WorldManager) {
	runtimeUpdates := make(map[string][]string)
	// Build all pending systems
	for _, sb := range manager.PendingSystems {
		id, queries, systemFunc, after, before := sb.Build()
		// Store queries as hashes. This helps us eliminate any duplicates that may occur between systems.
		queryKeys := make([]uint32, 0)
		for _, query := range queries {
			h := query.Hash()
			if _, ok := w.Queries[h]; !ok {
				w.Queries[h] = query
			}
			queryKeys = append(queryKeys, h)
		}
		// Store system data
		w.Systems[id] = SystemData{
			System:    systemFunc,
			Queries:   queryKeys,
			RunsAfter: after,
		}
		// Reflect any "before" requirements into "after" requirements on the relevant systems.
		if len(before) > 0 {
			for _, beforeId := range before {
				if _, ok := runtimeUpdates[beforeId]; !ok {
					runtimeUpdates[beforeId] = []string{}
				}
				runtimeUpdates[beforeId] = append(runtimeUpdates[beforeId], id)
			}
		}
	}
	// update reflected systems
	for id, extraRunsAfter := range runtimeUpdates {
		if _, ok := w.Systems[id]; !ok {
			continue
		}
		system := w.Systems[id]
		runsAfter := append(system.RunsAfter, extraRunsAfter...)
		system.RunsAfter = core.RemoveDuplicateStr(runsAfter)
		w.Systems[id] = system
	}
}

func (w *World) BuildQueries(output chan<- core.QueryResult) {
	for id, query := range w.Queries {
		go w.buildQuery(id, query, output)
	}
}

func (w *World) buildQuery(id uint32, query core.Query, output chan<- core.QueryResult) {
	// We are going to assume that all queries are valid at this stage of the process.
	entities := make(map[core.Entity][]core.Component)
	initialComponentType := query.Components[0]
	initialEntities := w.Components[initialComponentType]
	// Go through the first component type and find all entities that have the required component.
	// We will then go through the remaining component types and thin down this list to only include
	// entities that have all the required components.
	for entity, initialComponent := range initialEntities {
		entities[core.Entity(entity)] = []core.Component{initialComponent}
	}
	// Now we have a list of entities that have the first component- look for each of these in the remaining components.
	// If a component is not found, the entity is removed from the pool, thinning the list.
	for _, componentType := range query.Components[1:] {
		currentComponent := w.Components[componentType]
		for entity := range entities {
			if _, ok := currentComponent[entity]; !ok {
				delete(entities, entity)
			} else {
				entities[entity] = append(entities[entity], currentComponent[entity])
			}
		}
	}
	// At this point, we should have a good list- wrap it up and send it to the output channel.
	output <- core.QueryResult{
		ID:     id,
		Result: entities,
	}
}

func (w *World) RunSystem(id string, response chan<- string, manager *core.WorldManager, queries ...core.QueryResult) {
	if _, ok := w.Systems[id]; !ok {
		log.Fatal("System not found: ", id)
		return
	}
	system := w.Systems[id]

	// Run the system.
	system.System(manager, queries...)
	response <- id
}
