package ecs

import "github.com/fishykins/aero/pkg/nesting"

type EntityBuilder struct {
	components map[string]interface{}
}

func NewEntityBuilder() *EntityBuilder {
	return &EntityBuilder{
		components: map[string]interface{}{},
	}
}

func (eb *EntityBuilder) WithComponent(c interface{}) *EntityBuilder {
	cid := GetTypeId(c)
	eb.components[cid] = c
	return eb
}

func (eb *EntityBuilder) Named(name string) *EntityBuilder {
	component := Name(name)
	cid := GetTypeId(component)
	eb.components[cid] = component
	return eb
}

type ComponentBuilder struct {
	entity    Entity
	component interface{}
}

func NewComponentBuilder(entity Entity, component interface{}) *ComponentBuilder {
	return &ComponentBuilder{
		entity:    entity,
		component: component,
	}
}

type SystemBuilder struct {
	id         *string
	systemFunc SystemFunc
	queries    []Query
	runAfter   []string
	rubBefore  []string
}

func NewSystemBuilder(systemFunc SystemFunc) *SystemBuilder {
	return &SystemBuilder{
		id:         nil,
		systemFunc: systemFunc,
		queries:    []Query{},
	}
}

func (sb *SystemBuilder) WithQuery(query Query) *SystemBuilder {
	sb.queries = append(sb.queries, query)
	return sb
}

func (sb *SystemBuilder) Named(name string) *SystemBuilder {
	sb.id = &name
	return sb
}

func (sb *SystemBuilder) After(after ...interface{}) *SystemBuilder {
	for _, a := range after {
		sysId := GetTypeId(a)
		sb.runAfter = append(sb.runAfter, sysId)
	}
	return sb
}

func (sb *SystemBuilder) Before(before ...interface{}) *SystemBuilder {
	for _, b := range before {
		sysId := GetTypeId(b)
		sb.rubBefore = append(sb.rubBefore, sysId)
	}
	return sb
}

func (sb *SystemBuilder) getId() string {
	if sb.id == nil {
		GetTypeId(sb.systemFunc)
	}
	return *sb.id
}

func NewWorld() *World {
	return &World{
		Entities:   *nesting.NewNest(0),
		Components: map[string]map[Entity]interface{}{},
		Systems:    map[string]System{},
		Queries:    map[uint32]Query{},
		Resources:  map[string]interface{}{},
	}
}

func (w *World) BuildQuery(id uint32, query Query, output chan<- QueryResult) {
	// We are going to assume that all queries are valid at this stage of the process.
	entities := make(map[Entity][]interface{})
	initialComponentType := query.Tags()[0]
	initialEntities := w.Components[initialComponentType]
	// Go through the first component type and find all entities that have the required component.
	// We will then go through the remaining component types and thin down this list to only include
	// entities that have all the required components.
	for entity, initialComponent := range initialEntities {
		entities[Entity(entity)] = []interface{}{initialComponent}
	}
	// Now we have a list of entities that have the first component- look for each of these in the remaining components.
	// If a component is not found, the entity is removed from the pool, thinning the list.
	for _, componentType := range query.Tags()[1:] {
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
