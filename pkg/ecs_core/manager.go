package ecscore

import "fmt"

type Manager struct {
	PendingEntities []*EntityBuilder
	PendingSystems  []*SystemBuilder
	nextId          int
}

func NewManager() *Manager {
	return &Manager{
		PendingEntities: []*EntityBuilder{},
		PendingSystems:  []*SystemBuilder{},
		nextId:          0,
	}
}

func (m *Manager) AddEntity(args ...string) *EntityBuilder {
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

func (m *Manager) AddSystem(system System, queries ...Query) *SystemBuilder {
	sb := NewSystemBuilder(GetFunctionName(system), system, queries)
	m.PendingSystems = append(m.PendingSystems, sb)
	return sb
}
