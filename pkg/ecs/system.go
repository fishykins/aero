package ecs

type System interface {
	RequiredComponents() []string
	Init(world *World)
	Run(world *World, entities []*Entity)
}
