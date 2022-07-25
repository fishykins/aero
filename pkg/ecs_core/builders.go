package ecscore

import (
	"hash/fnv"
	"sort"
	"strings"
)

type EntityBuilder struct {
	id         string
	components map[string]Component
}

func NewEntityBuilder(id string) *EntityBuilder {
	return &EntityBuilder{
		id:         id,
		components: make(map[string]Component),
	}
}

func (eb *EntityBuilder) Named(id string) *EntityBuilder {
	eb.id = id
	return eb
}

func (eb *EntityBuilder) With(components ...Component) *EntityBuilder {
	for _, component := range components {
		eb.components[component.Type()] = component
	}
	return eb
}

func (eb *EntityBuilder) Build() (string, map[string]Component) {
	return eb.id, eb.components
}

type SystemBuilder struct {
	name    string
	queries []Query
	system  System
	after   []string
	before  []string
}

func NewSystemBuilder(name string, system System, queries []Query) *SystemBuilder {
	return &SystemBuilder{
		name:    name,
		queries: queries,
		system:  system,
		after:   []string{},
		before:  []string{},
	}
}

func (sb *SystemBuilder) Named(name string) *SystemBuilder {
	sb.name = name
	return sb
}

func (sb *SystemBuilder) After(tags ...string) *SystemBuilder {
	sb.after = append(sb.after, tags...)
	return sb
}

func (sb *SystemBuilder) Before(tags ...string) *SystemBuilder {
	sb.before = append(sb.before, tags...)
	return sb
}

func (sb *SystemBuilder) Build() (string, []Query, System, []string, []string) {
	for _, query := range sb.queries {
		query.Format()
	}
	return sb.name, sb.queries, sb.system, sb.after, sb.before
}

func (q *Query) Format() {
	comps := RemoveDuplicateStr(q.Components)
	sort.Strings(comps)
	q.Components = comps
}

func (q *Query) Hash() uint32 {
	h := fnv.New32a()
	s := strings.Join(q.Components[:], ",")
	h.Write([]byte(s))
	return h.Sum32()
}
