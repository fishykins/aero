package main

import (
	app "github.com/fishykins/aero/pkg/ecs_app"
	ecs "github.com/fishykins/aero/pkg/ecs_core"
	log "github.com/fishykins/aero/pkg/logging"
)

type age struct {
	Age int
}

type ageResource struct {
	speed int
}

func (a *age) Type() string {
	return "age"
}

func ageSystem(m *ecs.WorldManager, resources ecs.RMap, queries ...ecs.QueryResult) {
	ageResource := resources["ageResource"].(ageResource)
	ageQuery := queries[0]
	for _, components := range ageQuery.Result {
		age := components[0].(*age)
		age.Age += ageResource.speed
	}
}

func printAgeSystem(m *ecs.WorldManager, resources ecs.RMap, queries ...ecs.QueryResult) {
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
	app.AddResource(ageResource{speed: 2})
	app.AddSystem(ageSystem).WithQuery("age").WithResource("ageResource")
	app.AddSystem(printAgeSystem).After("ageSystem").WithQuery("age")
	app.Run()
}
