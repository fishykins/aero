package main

import (
	internal "github.com/fishykins/aero/internal"
	app "github.com/fishykins/aero/pkg/ecs_app"
	log "github.com/fishykins/aero/pkg/logging"
)

func main() {
	log.SetLevel(5)
	app := app.New()
	app.AddEntity().Named("earth").WithComponent(&internal.Sphere{Radius: 6371}).WithComponent(internal.CelestialBody{})
	app.AddEntity().Named("moon").WithComponent(&internal.Sphere{Radius: 1737}).WithComponent(internal.CelestialBody{})
	internal.AddPlanetSystem(app)
	app.Run()
}
