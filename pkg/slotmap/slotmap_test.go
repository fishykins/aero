package slotmap

import "testing"

type owl struct {
	Name  string
	sound string
}

func TestSlotMap(t *testing.T) {
	var slotmap = New[owl]("owls")
	slotmap.Add(owl{Name: "Hillan", sound: "Hoot"})
	fishyIndex := slotmap.Add(owl{Name: "Fishy", sound: "Twit Twoo"})
	slotmap.Add(owl{Name: "Bucky", sound: "Mwaaarkk"})
	slotmap.Add(owl{Name: "Bilbo", sound: "?????"})

	if slotmap.Len() != 4 {
		t.Errorf("Expected 4, got %d", slotmap.Len())
	}

	fishy, err := slotmap.Get(fishyIndex)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if fishy.Name != "Fishy" {
		t.Errorf("Expected Fishy, got %s", fishy.Name)
	}

	slotmap.Remove(fishyIndex)

	if slotmap.Len() != 3 {
		t.Errorf("Expected 3, got %d", slotmap.Len())
	}

	slotmap.Add(owl{Name: "Doddy", sound: "Twoo"})

	if slotmap.Len() != 4 {
		t.Errorf("Expected 4, got %d", slotmap.Len())
	}

	_, err = slotmap.Get(fishyIndex)
	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}
}
