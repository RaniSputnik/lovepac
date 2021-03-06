// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at 2017-09-08 19:24:49.178047 +0100 BST
// TODO add the commit hash in here too

package target

import (
	"text/template"
)

var loveTemplate = template.Must(template.New("love").Parse(`local quads = {}

{{range .Sprites -}}
quads['{{.Name}}'] = love.graphics.newQuad({{.Left}},{{.Top}},{{.Width}},{{.Height}},{{$.Width}},{{$.Height}})
{{end}}
return quads
`))

var starlingTemplate = template.Must(template.New("starling").Parse(`<TextureAtlas imagePath="{{.ImageFilename}}">
{{- range .Sprites}}
    <SubTexture name="{{.Name}}" x="{{.Left}}" y="{{.Top}}" width="{{.Width}}" height="{{.Height}}"/>
{{- end}}
</TextureAtlas>
`))
