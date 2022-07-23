package slotmap

import "errors"

type SlotMap[T interface{}] struct {
	key   string
	slots []T
	gen   []int
	empty []int
}

type SlotMapIndex struct {
	key   string
	index int
	gen   int
}

func (sm *SlotMap[T]) Len() int {
	return len(sm.slots) - len(sm.empty)
}

func New[T interface{}](key string) SlotMap[T] {
	return SlotMap[T]{
		key:   key,
		slots: make([]T, 0),
		gen:   make([]int, 0),
		empty: make([]int, 0),
	}
}

func WithCapacity[T interface{}](key string, capacity int) SlotMap[T] {
	return SlotMap[T]{
		key:   key,
		slots: make([]T, 0, capacity),
		gen:   make([]int, 0, capacity),
		empty: make([]int, 0),
	}
}

func (sm *SlotMap[T]) Add(item T) SlotMapIndex {

	var index int
	var gen int

	if len(sm.empty) > 0 {
		index = sm.empty[0]
		sm.empty = sm.empty[1:]
		sm.slots[index] = item
		gen = sm.gen[index]
		key := sm.key
		return SlotMapIndex{key, index, gen}
	}

	sm.slots = append(sm.slots, item)
	sm.gen = append(sm.gen, 0)
	return SlotMapIndex{
		key:   sm.key,
		index: len(sm.slots) - 1,
		gen:   0,
	}
}

func (sm *SlotMap[T]) Remove(i SlotMapIndex) error {
	index, err := sm.validateIndex(i)
	if err != nil {
		return err
	}
	sm.empty = append(sm.empty, index)
	sm.gen[index]++
	return nil
}

func (sm *SlotMap[T]) Get(i SlotMapIndex) (*T, error) {
	index, err := sm.validateIndex(i)
	if err != nil {
		return nil, err
	}
	return &sm.slots[index], nil
}

func (sm *SlotMap[t]) validateIndex(i SlotMapIndex) (int, error) {
	if i.key != sm.key {
		return -1, errors.New("invalid map key")
	}
	index := i.index
	if index < 0 || index >= len(sm.slots) {
		return -1, errors.New("index out of range")
	}
	if sm.gen[index] != i.gen {
		return -1, errors.New("index is stale")
	}
	for _, i := range sm.empty {
		if i == index {
			return -1, errors.New("index is empty")
		}
	}
	return index, nil
}
