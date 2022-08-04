package ecs

import "time"

type Name string

func (n Name) Type() string {
	return "Name"
}

type UpdateFrequency uint32

func (UpdateFrequency) Type() string {
	return "UpdateFrequency"
}

func (u UpdateFrequency) FPS() time.Duration {
	return time.Second / time.Duration(u)
}
