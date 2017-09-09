package target_test

import (
	"testing"

	"github.com/RaniSputnik/lovepac/target"
)

func TestIsValid(t *testing.T) {
	formats := map[target.Format]bool{
		target.Unknown:                                            false,
		target.Love:                                               true,
		target.Starling:                                           true,
		target.Format{Ext: "lua"}:                                 false,
		target.Format{Template: target.Love.Template}:             false,
		target.Format{Template: target.Love.Template, Ext: "lua"}: true,
	}

	for test, expect := range formats {
		got := test.IsValid()
		if got != expect {
			t.Errorf("Expected 'IsValid' for format '%v' to return '%t' but got '%t'", test, expect, got)
		}
	}
}
