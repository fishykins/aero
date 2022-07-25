package ecs2

type EntityBuilder struct {
	ecs    *ECS
	entity Entity
}

func (eb *EntityBuilder) Named(name string) *EntityBuilder {
	eb.entity.id = name
	return eb
}

func (eb *EntityBuilder) WithComponent(component ...Component) *EntityBuilder {
	for _, c := range component {
		eb.ecs.AddComponent(eb.entity, c)
	}
	return eb
}

func (eb *EntityBuilder) Unwrap() Entity {
	return eb.entity
}
