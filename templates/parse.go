package templates

import (
	"text/template"
)

var (
	Starling *template.Template
	Love     *template.Template
)

func init() {
	Starling = template.Must(template.New("starling").Parse(starlingStr))
	Love = template.Must(template.New("love").Parse(loveStr))
}
