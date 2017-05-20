package templates

const loveStr = `local quads = {}

{{range .Sprites -}}
quads['{{.Name}}'] = love.graphics.newQuad({{.Left}},{{.Top}},{{.Width}},{{.Height}},{{$.Width}},{{$.Height}})
{{end}}
return quads`
