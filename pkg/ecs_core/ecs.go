package ecscore

import "github.com/fishykins/aero/pkg/slotmap"

type Entity slotmap.SlotMapIndex

type Component interface {
	Type() string
}

type System func(*Manager, ...QueryResult)

type Query struct {
	Components []string
}

type QueryResult struct {
	ID     uint32
	Result map[Entity][]Component
}
