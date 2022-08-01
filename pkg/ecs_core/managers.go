package ecs

// The WorldManager is an interface between the end user and the ECS world.
type WorldManager struct {
	pendingEntities   []*EntityBuilder
	pendingComponents []*ComponentBuilder
	pendingSystems    []*SystemBuilder
	pendingResources  []interface{}
}

func NewWorldManager() *WorldManager {
	return &WorldManager{
		pendingEntities:   []*EntityBuilder{},
		pendingComponents: []*ComponentBuilder{},
		pendingSystems:    []*SystemBuilder{},
		pendingResources:  []interface{}{},
	}
}

func (wm *WorldManager) NewEntity() *EntityBuilder {
	return NewEntityBuilder()
}

func (wm *WorldManager) AddComponent(entity Entity, component interface{}) *ComponentBuilder {
	cb := NewComponentBuilder(entity, component)
	wm.pendingComponents = append(wm.pendingComponents, cb)
	return cb
}

func (wm *WorldManager) AddSystem(systemFunc SystemFunc) *SystemBuilder {
	sb := NewSystemBuilder(systemFunc)
	wm.pendingSystems = append(wm.pendingSystems, sb)
	return sb
}

func (wm *WorldManager) AddResource(resource interface{}) {
	wm.pendingResources = append(wm.pendingResources, resource)
}
