package internal

import (
	app "github.com/fishykins/aero/pkg/ecs_app"
	ecs "github.com/fishykins/aero/pkg/ecs_core"
	log "github.com/fishykins/aero/pkg/logging"
)

func planetSystem(m *ecs.WorldManager, resources ecs.RMap, queries ...ecs.QueryResult) {
	for _, r := range queries[0].Result {
		body := r[0].(CelestialBody)
		sphere := r[1].(*Sphere)
		name := r[2].(ecs.Name)
		sphere.Radius += 1
		log.InfoWith("Found planet", map[string]interface{}{"body": body, "sphere": sphere, "name": name})
	}
}

func AddPlanetSystem(a *app.App) {
	a.AddSystem(planetSystem).WithQuery("CelestialBody", "Sphere", "Name")
}
