package ecs

import (
	"fmt"
	"log"
)

func New() *ECS {
	return &ECS{
		tick: 0,
		World: World{
			Entities:   make([]Entity, 0),
			Systems:    make([]SystemPair, 0),
			Components: make(map[string]map[Entity]Component),
		},
	}
}

func (e *ECS) AddEntity() Entity {
	entity := Entity{id: fmt.Sprintf("Entity%d", len(e.World.Entities))}
	e.World.Entities = append(e.World.Entities, entity)
	return entity
}

func (e *ECS) AddSystem(system System, queries ...Query) {
	e.World.Systems = append(e.World.Systems, SystemPair{
		queries: queries,
		system:  system,
	})
	// TODO: Add support for queries being reused (between systems).
}

func (e *ECS) AddComponent(entity Entity, component Component) {
	if _, ok := e.World.Components[component.Type()]; !ok {
		e.World.Components[component.Type()] = make(map[Entity]Component)
	}
	e.World.Components[component.Type()][entity] = component
}

func (e *ECS) Update() {
	systems := e.BuildSystemExecutables()
	systemsComplete := make(chan int)
	for i, system := range systems {
		go executeSystem(&e.World, i, system, systemsComplete)
	}
	for i := 0; i < len(systems); i++ {
		ok := <-systemsComplete
		log.Println("System", ok, "completed")
	}
	e.tick++
}

func executeSystem(world *World, id int, system systemExecutable, done chan<- int) {
	system.system(world, system.queries...)
	done <- id
}

func (e *ECS) BuildSystemExecutables() []systemExecutable {
	systems := make([]systemExecutable, len(e.World.Systems))
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
		systems[systemIndex] = systemExecutable{
			system:  s.system,
			queries: results,
		}
	}
	return systems
}
