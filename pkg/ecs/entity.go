package ecs

type Entity struct {
	uuid int
}

func (e *Entity) UUID() int {
	return e.uuid
}
