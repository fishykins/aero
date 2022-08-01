package main

import (
	log "github.com/fishykins/aero/pkg/logging"

	app "github.com/fishykins/aero/pkg/ecs_app"
	core "github.com/fishykins/aero/pkg/ecs_core"
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
	for _, components := range ageQuery.Result {
		age := components[0].(*age)
		log.Info("Age:", age.Age)
	}
}

func main() {
	log.SetLevel(5)
	app := app.New()
	app.AddEntity().WithComponent(&age{Age: 29}).Named("Fishy")
	app.AddSystem(ageSystem).WithQuery("age")
	app.AddSystem(printAgeSystem).After("ageSystem").WithQuery("age")
	app.Run()
}
