package ecs

type World struct {
	Entities   []Entity
	Components map[*Component]map[int]Component
	nextId     int
}

func NewWorld() World {
	var world World = World{
		Components: make(map[*Component]map[int]Component),
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
	if _, ok := w.Components[&component]; !ok {
		w.Components[&component] = make(map[int]Component)
	}
	w.Components[&component][entity.UUID()] = component
}

func (w *World) GetComponent(entity Entity, component *Component) (Component, bool) {
	if _, ok := w.Components[component]; !ok {
		return nil, false
	}

	if _, ok := w.Components[component][entity.UUID()]; !ok {
		return nil, false
	}
	return w.Components[component][entity.UUID()], true
}

func (w *World) FilterEntities(with []*Component) []Entity {
	var entities []Entity = make([]Entity, 0)
	var scores map[*Entity]int = make(map[*Entity]int)
	var requiredScore int = len(with)

	for _, Component := range with {
		if _, ok := w.Components[Component]; !ok {
			return nil
		}
		for uuid, _ := range w.Components[Component] {
			var entity = w.GetEntity(uuid)

			if _, ok := scores[&entity]; !ok {
				scores[&entity] = 0
			}
			scores[&entity] += 1
			if scores[&entity] == requiredScore {
				entities = append(entities, entity)
			}
		}

	}
	return entities
}
