package nesting

import "github.com/fishykins/aero/pkg/slotmap"

// The Object-World-Level struct is a data structure used to represent the objects position in the object-world.
type Owl struct {
	Name     string
	Parent   *slotmap.SlotMapIndex
	Children []*slotmap.SlotMapIndex
	Data     interface{}
}
