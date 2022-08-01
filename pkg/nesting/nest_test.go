package nesting

import (
	"testing"
)

func TestNest(t *testing.T) {
	var nest = NewNest(8)
	fishy := nest.Add()
	hillan := nest.Add()

	fishy.Remove()

	frogbert := hillan.Birth().Unwrap()

	parent, err := nest.GetParent(frogbert)
	if err != nil {
		t.Error(err)
	}
	if parent.Unwrap() != hillan.Unwrap() {
		t.Error("Expected parent to be Hillan")
	}

	var majasSiblings []NestResult
	majasSiblings, err = hillan.Birth().Siblings()
	if err != nil {
		t.Error(err)
	}
	if len(majasSiblings) != 1 {
		for _, sibling := range majasSiblings {
			t.Log("sibling: ", sibling)
		}
		t.Error("Expected Maja to only have one sibling, but she has", len(majasSiblings))
	}
	if majasSiblings[0].Unwrap() != frogbert {
		t.Error("Expected Maja to be a sibling of Frogbert")
	}
}
