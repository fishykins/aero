package main

import (
	app "github.com/fishykins/aero/pkg/ecs_app"
	core "github.com/fishykins/aero/pkg/ecs_core"
	log "github.com/fishykins/aero/pkg/logging"
)

type age struct {
	Age int
}

func (a *age) Type() string {
	return "age"
}

func ageSystem(m *core.WorldManager, queries ...core.QueryResult) {
	ageQuery := queries[0]
	for _, components := range ageQuery.Result {
		age := components[0].(*age)
		age.Age += 1
	}
}

func printAgeSystem(m *core.WorldManager, queries ...core.QueryResult) {
	ageQuery := queries[0]
	for id, components := range ageQuery.Result {
		age := components[0].(*age)
		log.TraceWith("Age", map[string]interface{}{"age": age.Age, "name": id})
	}
}

func main() {
	log.Default().WithLevel(6)
	app := app.New()
	app.AddEntity("Fishy").With(&age{Age: 29})
	app.AddSystem(ageSystem, core.NewQuery("age"))
	app.AddSystem(printAgeSystem, core.NewQuery("age")).After("ageSystem")
	app.AddResource(core.UpdateFrequency(30))
	app.Run()
}
