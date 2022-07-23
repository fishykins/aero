package ecs

import "reflect"

// World is the main object of the ECS system. It contains all the entities and components,
// and should be used to instantiate new ones.
type World struct {
	Entities   []Entity
	Components map[reflect.Type]map[int]Component
	nextId     int
}

type Entity struct {
	uuid int
}

type Component interface{}

func NewWorld() World {
	var world World = World{
		Entities:   make([]Entity, 0),
		Components: make(map[reflect.Type]map[int]Component),
		nextId:     0,
	}
	return world
}

func (w *World) SpawnEntity() Entity {
	var uuid int = w.nextId
	w.nextId += 1
	var entity = Entity{uuid: uuid}
	w.Entities = append(w.Entities, entity)
	return entity
}

func (w *World) GetEntity(uuid int) Entity {
	return w.Entities[uuid]
}

func (w *World) AddComponent(entity Entity, component Component) {
	var compType = reflect.TypeOf(component)

	if _, ok := w.Components[compType]; !ok {
		w.Components[compType] = make(map[int]Component)
	}
	w.Components[compType][entity.UUID()] = component
}

func (w *World) GetComponent(entity Entity, component *Component) (Component, bool) {
	var compType = reflect.TypeOf(component)
	if _, ok := w.Components[compType]; !ok {
		return nil, false
	}

	if _, ok := w.Components[compType][entity.UUID()]; !ok {
		return nil, false
	}
	return w.Components[compType][entity.UUID()], true
}

func (e *Entity) UUID() int {
	return e.uuid
}
