package ecs

import (
	"log"
	"testing"

	core "github.com/fishykins/aero/pkg/ecs_core"
)

type age struct {
	Age int
}

func (a *age) Type() string {
	return "age"
}

func ageSystem(m *core.Manager, queries ...core.QueryResult) {
	ageQuery := queries[0]
	for _, components := range ageQuery.Result {
		age := components[0].(*age)
		age.Age += 1
	}
}

func printAgeSystem(m *core.Manager, queries ...core.QueryResult) {
	ageQuery := queries[0]
	for _, components := range ageQuery.Result {
		age := components[0].(*age)
		log.Println("Age:", age.Age)
	}
}

func TestApp(t *testing.T) {
	app := New()
	app.AddEntity("Fishy").With(&age{Age: 29})
	app.AddSystem(ageSystem, core.NewQuery("age"))
	app.AddSystem(printAgeSystem, core.NewQuery("age")).After("ageSystem")
	for i := 0; i < 4; i++ {
		app.Update()
	}
}
