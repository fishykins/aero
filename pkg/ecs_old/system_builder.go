package ecs2

import "log"

type SystemBuilder struct {
	ecs    *ECS
	system string
}

func (s *SystemBuilder) Named(name string) *SystemBuilder {
	s.ecs.World.Systems[name] = s.ecs.World.Systems[s.system]
	delete(s.ecs.World.Systems, s.system)
	for _, system := range s.ecs.World.Systems {
		for i, trigger := range system.after {
			if trigger == s.system {
				system.after[i] = name
			}
		}
		for i, trigger := range system.before {
			if trigger == s.system {
				system.before[i] = name
			}
		}
	}
	return s
}

func (s *SystemBuilder) After(other ...string) *SystemBuilder {
	data := s.ecs.World.Systems[s.system]
	data.after = append(data.after, other...)
	log.Println("Added system", s.system, "after", other, ":", data.after)

	s.ecs.World.Systems[s.system] = data
	return s
}

func (s *SystemBuilder) Before(other ...string) *SystemBuilder {
	data := s.ecs.World.Systems[s.system]
	data.before = append(data.before, other...)
	s.ecs.World.Systems[s.system] = data
	return s
}

func (s *SystemBuilder) Unwrap() string {
	return s.system
}
