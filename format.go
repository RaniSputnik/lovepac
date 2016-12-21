package main

type Format struct {
	Template string
	Ext      string
}

var FormatLookup = map[string]Format{
	"starling": Format{Template: "starling.template", Ext: "xml"},
}

func FormatIsValid(format string) bool {
	_, ok := FormatLookup[format]
	return ok
}
