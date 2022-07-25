package ecs2

import (
	"log"
	"testing"
)

type soundComponent struct {
	string
}

type ageComponent struct {
	int
}

func (c *soundComponent) Type() string {
	return "sound"
}

func (a *ageComponent) Type() string {
	return "age"
}

var basicQuery = NewQuery([]string{"sound", "age"}, []string{})
var ageQuery = NewQuery([]string{"age"}, []string{})

func TestApp(t *testing.T) {
	app := New()
	app.AddSystem(basicSystem, ageQuery).Named("basic").Before("second")
	app.AddSystem(secondSystem, basicQuery).Named("second")

	app.AddEntity().Named("Fishy").WithComponent(&ageComponent{29}, &soundComponent{"meeeep"})
	app.AddEntity().Named("Hillan").WithComponent(&ageComponent{30}).WithComponent(&soundComponent{"eeeeghhh"})

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
		for entity, components := range query.Entities {
			sound := components[0].(*soundComponent)
			age := components[1].(*ageComponent)
			log.Printf("%s is %d years old and says %s\n", entity.id, age.int, sound.string)
		}
	}
}
