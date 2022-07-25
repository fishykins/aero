package ecs2

import (
	"fmt"
	"log"
)

type EcsBuilder struct {
}

func New() *ECS {
	return &ECS{
		tick: 0,
		World: World{
			Entities:   make([]Entity, 0),
			Systems:    make(map[string]SystemData),
			Components: make(map[string]map[Entity]Component),
		},
	}
}

func (e *ECS) AddEntity() *EntityBuilder {
	entity := Entity{id: fmt.Sprintf("Entity%d", len(e.World.Entities))}
	e.World.Entities = append(e.World.Entities, entity)
	return &EntityBuilder{
		ecs:    e,
		entity: entity,
	}
}

func (e *ECS) AddSystem(system System, queries ...Query) *SystemBuilder {
	defaultName := GetFunctionName(system)
	e.World.Systems[defaultName] = SystemData{
		queries: queries,
		system:  system,
		name:    defaultName,
	}
	log.Println("Added system", defaultName)
	// TODO: Add support for queries being reused (between systems).
	return &SystemBuilder{
		ecs:    e,
		system: defaultName,
	}
}

func (e *ECS) AddComponent(entity Entity, component Component) {
	if _, ok := e.World.Components[component.Type()]; !ok {
		e.World.Components[component.Type()] = make(map[Entity]Component)
	}
	e.World.Components[component.Type()][entity] = component
}
