package nesting

import (
	"errors"

	"github.com/fishykins/aero/pkg/slotmap"
)

type Nest struct {
	owls slotmap.SlotMap[Owl]
}

//NOTE: All functions here should only accept reference parameters of type 'SlotMapIndex' and return 'NestResult's where possible.

func NewNest(initialCapacity int) *Nest {
	return &Nest{
		owls: slotmap.WithCapacity[Owl]("nest", initialCapacity),
	}
}

// Adds a new data structure to the nest, returning the index of the subsequent Owl.
func (n *Nest) Add() NestResult {
	var owl = Owl{
		Parent:   nil,
		Children: make([]*slotmap.SlotMapIndex, 0),
	}

	index := n.owls.Add(owl)
	return NestResult{
		nest:   n,
		result: &index,
	}
}

// Removes the data structure at the given index, or returns an error if the index is invalid.
func (n *Nest) Remove(index slotmap.SlotMapIndex) error {
	return n.owls.Remove(index)
}

func (n *Nest) Adopt(parent slotmap.SlotMapIndex, child slotmap.SlotMapIndex) error {
	parentOwl, err := n.owls.Get(parent)
	if err != nil {
		return errors.New("parent index error: " + err.Error())
	}
	childOwl, err := n.owls.Get(child)
	if err != nil {
		return errors.New("child index error: " + err.Error())
	}
	childOwl.Parent = &parent
	parentOwl.Children = append(parentOwl.Children, &child)
	return nil
}

func (n *Nest) Orphan(index slotmap.SlotMapIndex) error {
	owl, err := n.owls.Get(index)
	if err != nil {
		return err
	}
	parentIndex := owl.Parent
	owl.Parent = nil
	if parentIndex != nil {
		parent, err := n.owls.Get(*parentIndex)
		if err != nil {
			return errors.New("parent index not found, despite being set")
		}
		for i, childIndex := range parent.Children {
			if *childIndex == index {
				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
				return nil
			}
		}

	}
	return nil
}

func (n *Nest) GetParent(index slotmap.SlotMapIndex) (NestResult, error) {
	owl, err := n.owls.Get(index)
	if err != nil {
		return NestResult{nest: n, result: nil}, err
	}
	return NestResult{nest: n, result: owl.Parent}, nil
}

func (n *Nest) GetChildren(index slotmap.SlotMapIndex) ([]NestResult, error) {
	owl, err := n.owls.Get(index)
	if err != nil {
		return nil, err
	}
	children := owl.Children
	var childrenResults []NestResult = make([]NestResult, len(children))
	for i, child := range children {
		childrenResults[i] = NestResult{nest: n, result: child}
	}
	return childrenResults, nil
}

func (n *Nest) GetSiblings(index slotmap.SlotMapIndex) ([]NestResult, error) {
	owl, err := n.owls.Get(index)
	if err != nil {
		return make([]NestResult, 0), err
	}
	if owl.Parent == nil {
		// This is a very sad result :(
		return make([]NestResult, 0), nil
	}
	parent, err := n.owls.Get(*owl.Parent)
	if err != nil {
		return make([]NestResult, 0), errors.New("Parent not found: " + err.Error())
	}
	children := parent.Children
	var siblings []NestResult = make([]NestResult, len(children)-1)
	i := 0
	for _, child := range children {
		if *child != index {
			siblings[i] = NestResult{nest: n, result: child}
			i++
		}
	}
	return siblings, nil
}
