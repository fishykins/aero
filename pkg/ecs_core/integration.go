package ecs

import "log"

func (w *World) RunSystem(id string, response chan<- string, manager *WorldManager, queries ...QueryResult) {
	if _, ok := w.Systems[id]; !ok {
		log.Fatal("System not found: ", id)
		return
	}
	system := w.Systems[id]

	// Run the system.
	system.run(manager, queries...)
	response <- id
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
		return
	}
	m.updateWorldGo(w)
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
			h := HashTags(query.Tags())
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
			runsAfter: sb.runAfter,
		}
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
	if c != nil {
		c <- true
	}
}
