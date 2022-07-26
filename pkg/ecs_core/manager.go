package ecscore

import (
	"fmt"
)

// ============================================================================
// ============================= WORLD MANAGER ================================
// ============================================================================

// The world manager is a user-facing endpoint for instantiating ECS data, and any logic kept to the bare minimum.
// Integration is handled by the native 'World' struct inside the 'inecs' package.
type WorldManager struct {
	PendingEntities   []*EntityBuilder
	PendingSystems    []*SystemBuilder
	PendingComponents map[string]map[Entity]Component
	nextId            int
	// Unlike the ECS data, resources are stored in the world manager.
	// This is because they are quick and easy to access, and require little to no marshalling.
	// Putting them in the world would require a massive work around to get access to them in systems!
	resources map[string]Component
}

func NewWorldManager() *WorldManager {
	return &WorldManager{
		PendingEntities:   []*EntityBuilder{},
		PendingSystems:    []*SystemBuilder{},
		PendingComponents: map[string]map[Entity]Component{},
		nextId:            0,
	}
}

func (m *WorldManager) AddEntity(args ...string) *EntityBuilder {
	var name string = fmt.Sprintf("entity%d", m.nextId)
	if len(args) > 0 {
		name = args[0]
	} else {
		m.nextId++
	}
	eb := NewEntityBuilder(name)
	m.PendingEntities = append(m.PendingEntities, eb)
	return eb
}

func (m *WorldManager) AddComponent(entity *Entity, component Component) {
	if m.PendingComponents[component.Type()] == nil {
		m.PendingComponents[component.Type()] = map[Entity]Component{}
	}
	m.PendingComponents[component.Type()][*entity] = component
}

func (m *WorldManager) AddSystem(system System, queries ...Query) *SystemBuilder {
	sb := NewSystemBuilder(GetFunctionName(system), system, queries)
	m.PendingSystems = append(m.PendingSystems, sb)
	return sb
}

func (m *WorldManager) AddResource(name string, component Component) {
	if m.resources == nil {
		m.resources = map[string]Component{}
	}
	m.resources[name] = component
}

func (m *WorldManager) GetResource(name string) (Component, error) {
	if m.resources == nil {
		return nil, fmt.Errorf("no resource named '%s'", name)
	}
	return m.resources[name], nil
}

func (m *WorldManager) Resources() map[string]Component {
	return m.resources
}
