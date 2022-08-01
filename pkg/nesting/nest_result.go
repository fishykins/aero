package nesting

import "github.com/fishykins/aero/pkg/slotmap"

// A useful struct to help with the chainging of nest functions. Since we never really want to expose
// the raw index value or Owl struct, we can use this as a handle for the end user.
type NestResult struct {
	nest   *Nest
	result *slotmap.SlotMapIndex
}

// Returns the index value of the result.
func (n NestResult) Unwrap() slotmap.SlotMapIndex {
	return *n.result
}

// Returns the data structure at the given index, or returns an error if the index is invalid.
func (n NestResult) Remove() error {
	return n.nest.Remove(*n.result)
}

// Removes parent data from the given owl.
func (n NestResult) Orphan() error {
	return n.nest.Orphan(*n.result)
}

// Adopts the given child data structure to the given parent data structure.
func (n NestResult) Adopt(child *slotmap.SlotMapIndex) error {
	return n.nest.Adopt(*n.result, *child)
}

// Instantiates a new Owl and adds it to the nest as a child of the given parent.
func (n NestResult) Birth() NestResult {
	child := n.nest.Add()
	n.nest.Adopt(*n.result, child.Unwrap())
	return child
}

// Returns the parent of the given owl.
func (n NestResult) Parent() (NestResult, error) {
	parent, err := n.nest.GetParent(*n.result)
	if err != nil {
		return NestResult{nest: n.nest}, err
	}
	return parent, nil
}

func (n NestResult) Children() ([]NestResult, error) {
	children, err := n.nest.GetChildren(*n.result)
	if err != nil {
		return nil, err
	}
	return children, nil
}

func (n NestResult) Siblings() ([]NestResult, error) {
	siblings, err := n.nest.GetSiblings(*n.result)
	if err != nil {
		return nil, err
	}
	return siblings, nil
}
