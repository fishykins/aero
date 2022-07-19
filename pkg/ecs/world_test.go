package ecs

import "testing"

type Name struct {
	value string
}

type Age struct {
	value int
}

func TestWorld(t *testing.T) {
	var world World = NewWorld()

	var entity1 = world.SpawnEntity()
	var name1 = Name{value: "John"}
	var age1 = Age{value: 30}
	world.AddComponent(entity1, name1)
	world.AddComponent(entity1, age1)

	var entity2 = world.SpawnEntity()
	var name2 = Name{value: "Jane"}
	var age2 = Age{value: 25}
	world.AddComponent(entity2, name2)
	world.AddComponent(entity2, age2)

	var entity3 = world.SpawnEntity()
	var name3 = Name{value: "Jack"}
	world.AddComponent(entity3, name3)

}
