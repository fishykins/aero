package ecs

import (
	"reflect"

	log "github.com/fishykins/aero/pkg/logging"
)

func (w *World) RunSystem(id string, response chan<- string, manager *WorldManager, resources RMap, queries ...QueryResult) {
	if _, ok := w.Systems[id]; !ok {
		log.FatalWith("System not found", map[string]interface{}{"id": id, "queries": queries})
		return
	}
	system := w.Systems[id]

	// Run the system.
	log.TraceWith("Running system...", map[string]interface{}{"id": id})
	system.run(manager, resources, queries...)
	response <- id
}

func (w *World) BuildQuery(id uint32, query Query, output chan<- QueryResult) {
	// We are going to assume that all queries are valid at this stage of the process.
	entities := make(map[Entity][]interface{})
	initialComponentType := query[0]
	initialEntities := w.Components[initialComponentType]
	// Go through the first component type and find all entities that have the required component.
	// We will then go through the remaining component types and thin down this list to only include
	// entities that have all the required components.
	for entity, initialComponent := range initialEntities {
		entities[Entity(entity)] = []interface{}{initialComponent}
	}
	// Now we have a list of entities that have the first component- look for each of these in the remaining components.
	// If a component is not found, the entity is removed from the pool, thinning the list.
	for _, componentType := range query[1:] {
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
	output <- QueryResult{
		ID:     id,
		Result: entities,
	}
}

// Sends all pending data to the world, usually at the end of a frame.
func (m *WorldManager) UpdateWorld(w *World, concurency bool) {
	if !concurency {
		if len(m.pendingEntities) > 0 {
			m.updateEntities(w, nil)
		}
		if len(m.pendingComponents) > 0 {
			m.updateComponents(w, nil)
		}
		if len(m.pendingSystems) > 0 {
			m.updateSystems(w, nil)
		}
		if len(m.pendingResources) > 0 {
			m.updateResources(w, nil)
		}
	} else {
		m.updateWorldGo(w)
	}
}

func (m *WorldManager) updateWorldGo(w *World) {
	c := make(chan bool)
	j := 0
	if len(m.pendingEntities) > 0 {
		j++
		go m.updateEntities(w, c)
	}
	if len(m.pendingComponents) > 0 {
		j++
		go m.updateComponents(w, c)
	}
	if len(m.pendingSystems) > 0 {
		j++
		go m.updateSystems(w, c)
	}
	if len(m.pendingResources) > 0 {
		j++
		go m.updateResources(w, c)
	}
	// Wait for all updates to finish
	for i := 0; i < j; i++ {
		<-c
	}
}

func (m *WorldManager) updateEntities(w *World, c chan<- bool) {
	for _, eb := range m.pendingEntities {
		entity := Entity(w.Entities.Add())

		for _, component := range eb.components {
			typeId := GetTypeId(component)
			if _, ok := w.Components[typeId]; !ok {
				w.Components[typeId] = make(map[Entity]interface{})
			}
			w.Components[typeId][entity] = component
		}
	}
	m.pendingEntities = make([]*EntityBuilder, 0)
	if c != nil {
		c <- true
	}
}

func (m *WorldManager) updateComponents(w *World, c chan<- bool) {
	for _, cb := range m.pendingComponents {
		typeId := GetTypeId(cb.component)
		if _, ok := w.Components[typeId]; !ok {
			w.Components[typeId] = make(map[Entity]interface{})
		}
		w.Components[typeId][cb.entity] = cb.component
	}
	m.pendingComponents = make([]*ComponentBuilder, 0)
	if c != nil {
		c <- true
	}
}

func (m *WorldManager) updateSystems(w *World, c chan<- bool) {
	runtimeUpdates := make(map[string][]string)
	// Build all pending systems
	for _, sb := range m.pendingSystems {
		// Store queries as hashes. This helps us eliminate any duplicates that may occur between systems.
		queryKeys := make([]uint32, 0)
		for _, query := range sb.queries {
			h := HashTags(query)
			if _, ok := w.Queries[h]; !ok {
				w.Queries[h] = query
			}
			queryKeys = append(queryKeys, h)
		}
		// Store system data
		id := sb.getId()
		w.Systems[id] = System{
			run:       sb.systemFunc,
			queries:   queryKeys,
			resources: sb.resources,
			runsAfter: sb.runAfter,
		}
		log.Info("System added: ", w.Systems[id])
		// Reflect any "before" requirements into "after" requirements on the relevant systems.
		if len(sb.rubBefore) > 0 {
			for _, beforeId := range sb.rubBefore {
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
		runsAfter := append(system.runsAfter, extraRunsAfter...)
		system.runsAfter = RemoveDuplicateStr(runsAfter)
		w.Systems[id] = system
	}
	m.pendingSystems = make([]*SystemBuilder, 0)
	if c != nil {
		c <- true
	}
}

func (m *WorldManager) updateResources(w *World, c chan<- bool) {
	for _, rb := range m.pendingResources {
		typeId := reflect.TypeOf(rb).Name()
		log.InfoWith("resource added", map[string]interface{}{"type": typeId, "data": rb})
		w.Resources[typeId] = rb
	}
	m.pendingResources = make([]interface{}, 0)
	if c != nil {
		c <- true
	}
}
