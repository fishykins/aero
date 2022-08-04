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
	resources  []string
	runAfter   []string
	rubBefore  []string
}

func NewSystemBuilder(systemFunc SystemFunc) *SystemBuilder {
	return &SystemBuilder{
		id:         nil,
		systemFunc: systemFunc,
		queries:    []Query{},
		resources:  []string{},
	}
}

// Query results will be provided in the order in which they are added to the system.
func (sb *SystemBuilder) WithQuery(tags ...string) *SystemBuilder {
	sb.queries = append(sb.queries, tags)
	return sb
}

// Adds a resource required by the system. All resources are provided in a map, indexed by their given type as requested here.
func (sb *SystemBuilder) WithResource(r string) *SystemBuilder {
	sb.resources = append(sb.resources, r)
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
		return GetTypeId(sb.systemFunc)
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
