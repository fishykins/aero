package nesting

import (
	"testing"
)

type animalSound struct {
	string
}

func TestNest(t *testing.T) {
	var nest = NewNest(8)
	fishy := nest.Add(animalSound{"Quack"})
	hillan := nest.Add(animalSound{"BUNNIESSS"}, "Hillan")

	noise, err := fishy.Data()
	if err != nil {
		t.Error(err)
	}
	t.Log(*noise)

	fishy.Remove()

	_, err = nest.Data(fishy.Index())
	if err == nil {
		t.Error("Expected error")
	}

	frogbert := hillan.Birth(animalSound{"Ribbit"}, "Frogbert").Index()

	var parent NestResult
	parent, err = nest.GetParent(frogbert)
	if err != nil {
		t.Error(err)
	}
	if parent.Index() != hillan.Index() {
		t.Error("Expected parent to be Hillan")
	}

	var majasSiblings []NestResult
	majasSiblings, err = hillan.Birth(animalSound{"nus nus nus"}, "Maja").Siblings()
	if err != nil {
		t.Error(err)
	}
	if len(majasSiblings) != 1 {
		for _, sibling := range majasSiblings {
			t.Log("sibling: ", sibling)
		}
		t.Error("Expected Maja to only have one sibling, but she has", len(majasSiblings))
	}
	if majasSiblings[0].Index() != frogbert {
		t.Error("Expected Maja to be a sibling of Frogbert")
	}
}
