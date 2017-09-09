package target

import "text/template"

//go:generate go run gen.go

// Format represents a target atlas format.
//
// Name is simply the name of the format.
//
// Template is the text template that will
// be used to render the atlas descriptor file.
//
// Ext is the file extension that should be
// used when the descriptor file is written to
// the file system.
type Format struct {
	Name     string
	Template *template.Template
	Ext      string
}

var (
	// Unknown format, should used for error responses
	Unknown = Format{"unknown", nil, ""}
	// Love format for the love2d game engine
	Love = Format{"love", loveTemplate, "lua"}
	// Starling format for the Starling game engine
	Starling = Format{"starling", starlingTemplate, "xml"}
)

var allFormats = []Format{Love, Starling}

// FormatNamed returns the format with the given name
func FormatNamed(name string) Format {
	for _, format := range allFormats {
		if format.Name == name {
			return format
		}
	}
	return Unknown
}
