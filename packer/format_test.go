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

func TestGetFormatNamed(t *testing.T) {
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
		gotFormat, gotErr := packer.GetFormatNamed(test)
		if expect {
			if gotErr != nil {
				t.Errorf("Expected format '%s' returned without error, but got '%s'", test, gotErr)
			}
			if gotFormat == nil {
				t.Errorf("Expected format '%s' returned but got 'nil'", test)
				continue
			}
			if gotFormat.Template == nil {
				t.Errorf("Expected format '%s' returned with valid template but got 'nil'", test)
			}
			if gotFormat.Ext == "" {
				t.Errorf("Expected format '%s' returned with valid file extension but got empty string", test)
			}
		} else {
			if gotFormat != nil {
				t.Errorf("Expected invalid format '%s' returned but got '%v'", test, gotFormat)
			}
			expectedErr := packer.ErrFormatIsInvalid
			if gotErr != expectedErr {
				t.Errorf("Expected invalid format '%s' returned with '%s', but got '%s'", test, expectedErr, gotErr)
			}
		}
	}
}
