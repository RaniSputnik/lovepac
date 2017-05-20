package packer

import (
	"text/template"

	"github.com/RaniSputnik/lovepac/templates"
)

type Format struct {
	Template *template.Template
	Ext      string
}

var FormatLookup = map[string]Format{
	"starling": Format{Template: templates.Starling, Ext: "xml"},
	"love":     Format{Template: templates.Love, Ext: "lua"},
}

func FormatIsValid(format string) bool {
	_, ok := FormatLookup[format]
	return ok
}
