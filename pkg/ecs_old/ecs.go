package ecs2

import "log"

func (e *ECS) Update() {
	systems := e.buildSystemExecutables()
	delayedSystems := make(map[string]systemExecutable)
	systemsComplete := make(chan string)
	for name, system := range systems {
		if !system.delayed {
			go executeSystem(&e.World, name, system, systemsComplete)
		} else {
			delayedSystems[name] = system
		}
	}
	for i := 0; i < len(systems); i++ {
		finishedSystem := <-systemsComplete
		for delayedSystemName, delayedSystem := range delayedSystems {
			for j, trigger := range delayedSystem.triggers {
				if trigger == finishedSystem {
					delayedSystem.triggers = append(delayedSystem.triggers[:j], delayedSystem.triggers[j+1:]...)
					delayedSystems[delayedSystemName] = delayedSystem
				}
				if len(delayedSystem.triggers) == 0 {
					delete(delayedSystems, delayedSystemName)
					go executeSystem(&e.World, delayedSystemName, delayedSystem, systemsComplete)
				}
			}
		}
	}
	e.tick++
}

func executeSystem(world *World, id string, system systemExecutable, done chan<- string) {
	system.system(world, system.queries...)
	done <- id
}

func (e *ECS) buildSystemExecutables() map[string]systemExecutable {
	systems := make(map[string]systemExecutable, len(e.World.Systems))
	// A list of systems that have been implicitly delayed.
	delayedSystems := make(map[string][]string)

	for systemIndex, s := range e.World.Systems {
		results := make([]QueryResult, len(s.queries))
		for _, query := range s.queries {
			// get all the entities that match this query.
			entities := make(map[Entity][]Component)
			for componentIndex, with := range query.With {
				lookingFor := with.Type
				if _, ok := e.World.Components[lookingFor]; !ok {
					continue
				}
				for entity, component := range e.World.Components[lookingFor] {
					if componentIndex == 0 {
						entities[entity] = append(entities[entity], component)
					} else {
						// by this point we have already filtered so if not found, discard.
						if _, ok := entities[entity]; !ok {
							continue
						}
						if len(entities[entity]) == componentIndex {
							// This entity has the right number of components already so we can add this one.
							entities[entity] = append(entities[entity], component)
						} else {
							// This entity does not have the right number of components- discard
							delete(entities, entity)
						}
					}
				}
			}
			// Now we have a list of entities that have their components pulled out.
			// Filter out the entities that don't have the right number of components.
			filteredEntities := make(map[Entity][]Component)
			for entity, components := range entities {
				if len(components) == len(query.With) {
					filteredEntities[entity] = components
				}
			}

			results = append(results, QueryResult{
				Entities: filteredEntities,
			})
		}
		delayed := len(s.after) > 0
		systems[systemIndex] = systemExecutable{
			system:   s.system,
			queries:  results,
			triggers: s.after,
			delayed:  delayed,
		}
		// Add all system.before requirements to the delayed systems.
		for _, trigger := range s.before {
			if _, ok := delayedSystems[trigger]; !ok {
				delayedSystems[trigger] = []string{}
			}
			delayedSystems[trigger] = append(delayedSystems[trigger], systemIndex)
		}
	}
	// Final pass for delayed systems.
	for name, delayedSystem := range delayedSystems {
		// system := systems[name]
		// system.delayed = true
		// system.triggers = delayedSystem
		// systems[name] = system
		log.Println("Delayed system:", name, "-> triggers:", delayedSystem)
	}
	return systems
}
