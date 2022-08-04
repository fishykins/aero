package ecs

import (
	"errors"

	"github.com/fishykins/aero/pkg/nesting"
)

type World struct {
	Entities   nesting.Nest
	Components map[string]map[Entity]interface{}
	Systems    map[string]System
	Queries    map[uint32]Query
	Resources  map[string]interface{}
}

type Entity nesting.NestResult

type Labled interface {
	Type() string
}

type System struct {
	run       SystemFunc
	runsAfter []string
	queries   []uint32
	resources []string
}

type SystemFunc func(manager *WorldManager, resources RMap, queries ...QueryResult)

type Query []string

type QueryResult struct {
	ID     uint32
	Result map[Entity][]interface{}
}

type RMap map[string]interface{}

func (w *World) GetResource(name string) (interface{}, error) {
	if res, ok := w.Resources[name]; ok {
		return res, nil
	}
	return nil, errors.New("resource not found")
}

func (s System) GetQueries() []uint32 {
	return s.queries
}

func (s System) GetResources() []string {
	return s.resources
}

func (s System) SystemsBefore() []string {
	return s.runsAfter
}
