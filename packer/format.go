package packer

import (
	"errors"
	"text/template"

	"github.com/RaniSputnik/lovepac/target"
)

type Format struct {
	Template *template.Template
	Ext      string
}

const (
	FormatStarling = "starling"
	FormatLove     = "love"
)

var ErrFormatIsInvalid = errors.New("Format is not valid")

var formatLookup = map[string]*Format{
	FormatStarling: &Format{Template: target.Starling, Ext: "xml"},
	FormatLove:     &Format{Template: target.Love, Ext: "lua"},
}

func FormatIsValid(format string) bool {
	_, ok := formatLookup[format]
	return ok
}

func GetFormatNamed(format string) (*Format, error) {
	if !FormatIsValid(format) {
		return nil, ErrFormatIsInvalid
	}
	return formatLookup[format], nil
}
