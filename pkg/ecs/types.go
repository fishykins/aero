package ecs

// These are core types used by the ecs package.
type ECS struct {
	World World
	tick  int64
}

type World struct {
	// The list of entities in the world.
	Entities []Entity
	// The list of systems in the world.
	Systems []SystemPair
	// The list of components in the world.
	Components map[string]map[Entity]Component
}

type Entity struct {
	id string
}

type Component interface {
	Type() string
}

type System func(w *World, q ...QueryResult)

type SystemPair struct {
	queries []Query
	system  System
}

type Query struct {
	With    []QueryComponent
	Without []string
}

type QueryComponent struct {
	Type string
}

type QueryResult struct {
	Entities map[Entity][]Component
}

// Less important helper structs used internally by the ecs package.
type systemExecutable struct {
	system  func(w *World, q ...QueryResult)
	queries []QueryResult
}
