package ecs

import (
	"log"
	"testing"
)

type nameComponent struct {
	string
}

type ageComponent struct {
	int
}

func (c *nameComponent) Type() string {
	return "name"
}

func (a *ageComponent) Type() string {
	return "age"
}

var basicQuery = NewQuery([]string{"name", "age"}, []string{})
var ageQuery = NewQuery([]string{"age"}, []string{})

func TestApp(t *testing.T) {
	app := New()
	app.AddSystem(basicSystem, ageQuery)
	app.AddSystem(secondSystem, basicQuery)

	fishy := app.AddEntity()
	app.AddComponent(fishy, &nameComponent{"fishy"})
	app.AddComponent(fishy, &ageComponent{29})

	hillan := app.AddEntity()
	app.AddComponent(hillan, &nameComponent{"hillan"})
	app.AddComponent(hillan, &ageComponent{30})

	// Loop 5 times
	for i := 0; i < 5; i++ {
		app.Update()
	}
}

func basicSystem(w *World, q ...QueryResult) {
	for _, query := range q {
		for _, components := range query.Entities {
			age := components[0].(*ageComponent)
			age.int += 1
		}
	}
}

func secondSystem(w *World, q ...QueryResult) {
	for _, query := range q {
		for _, components := range query.Entities {
			name := components[0].(*nameComponent)
			age := components[1].(*ageComponent)
			log.Printf("%s is %d years old\n", name.string, age.int)
		}
	}
}
