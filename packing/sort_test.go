package packing_test

import (
	"sort"
	"testing"

	. "github.com/RaniSputnik/lovepac/packing"
)

var blocks = []Block{
	&TestBlock{id: "1", w: 200, h: 200},
	&TestBlock{id: "2", w: 100, h: 100},
	&TestBlock{id: "3", w: 100, h: 50},
	&TestBlock{id: "4", w: 20, h: 600},
	&TestBlock{id: "5", w: 512, h: 200},
}

func TestSortByArea(t *testing.T) {
	expected := []string{"5", "1", "4", "2", "3"}

	sort.Sort(ByArea(blocks))

	for i := range blocks {
		got := blocks[i].(*TestBlock)
		if got.id != expected[i] {
			t.Errorf("Expected '%s' at index %d, got '%s'", expected[i], i, got.id)
		} else if testing.Verbose() {
			t.Logf("Found '%s' at index %d - this is correct", got.id, i)
		}
	}
}

func TestSortByMaxSide(t *testing.T) {
	expected := []string{"4", "5", "1", "2", "3"}

	sort.Sort(ByMaxSide(blocks))

	for i := range blocks {
		got := blocks[i].(*TestBlock)
		if got.id != expected[i] {
			t.Errorf("Expected '%s' at index %d, got '%s'", expected[i], i, got.id)
		} else if testing.Verbose() {
			t.Logf("Found '%s' at index %d - this is correct", got.id, i)
		}
	}
}
