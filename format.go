package main

type Format struct {
	Template string
	Ext      string
}

var FormatLookup = map[string]Format{
	"starling": Format{Template: "starling.template", Ext: "xml"},
	"love":     Format{Template: "love.template", Ext: "lua"},
}

func FormatIsValid(format string) bool {
	_, ok := FormatLookup[format]
	return ok
}
