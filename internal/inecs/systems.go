package inecs

import core "github.com/fishykins/aero/pkg/ecs_core"

type SystemData struct {
	System    core.System
	Queries   []uint32
	RunsAfter []string
}
