package packer_test

import (
	"testing"

	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
)

func TestFormatIsValid(t *testing.T) {
	tests := map[string]bool{
		packer.FormatLove:     true,
		packer.FormatStarling: true,
		"foo":          false,
		"__bar__":      false,
		"nil":          false,
		"-1E+02":       false,
		"Œ„´‰ˇÁ¨ˆØ∏”’": false,
	}

	for test, expect := range tests {
		got := packer.FormatIsValid(test)
		if expect != got {
			t.Errorf("Expected 'FormatIsValid(\"%s\") == %t', but got '%t'", test, expect, got)
		}
	}
}

func TestGetFormatNamed(t *testing.T) {
	tests := map[string]target.Format{
		packer.FormatLove:     target.Love,
		packer.FormatStarling: target.Starling,
		"foo":          target.Unknown,
		"__bar__":      target.Unknown,
		"nil":          target.Unknown,
		"-1E+02":       target.Unknown,
		"Œ„´‰ˇÁ¨ˆØ∏”’": target.Unknown,
	}

	for test, expect := range tests {
		gotFormat := packer.GetFormatNamed(test)
		if gotFormat != expect {
			t.Errorf("Expected format '%s' returned but got '%s'", test, gotFormat.Name)
			continue
		}
		if expect == target.Unknown {
			continue // Don't check template and ext on an unknown format
		}
		if gotFormat.Template == nil {
			t.Errorf("Expected format '%s' returned with valid template but got 'nil'", test)
		}
		if gotFormat.Ext == "" {
			t.Errorf("Expected format '%s' returned with valid file extension but got empty string", test)
		}
	}
}
