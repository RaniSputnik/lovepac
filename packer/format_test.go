package packer_test

import (
	"testing"

	"github.com/RaniSputnik/lovepac/packer"
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
