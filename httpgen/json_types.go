package httpgen

import (
	"bytes"
	"fmt"
	"go/format"
	"path"
	"strings"
	"text/template"

	"gitlab.com/pirates1/broodmother"
)

var ioJSONType = `
package {{.Package}}

import (
	{{range .Imports}}{{.}}
	{{end}}
)

{{range .Requests}}
type {{.Name}}JSONRequest struct {
    {{range .Fields}}{{.Name}}    {{.Type}} ` + "`{{.Tag}}`" + `
    {{end}}
}
{{end}}

{{range .Responses}}
type {{.Name}}JSONResponses struct {
    {{range .Fields}}{{.Name}}    {{.Type}} ` + "`{{.Tag}}`" + `
    {{end}}
}
{{end}}
`

func (g *HttpGen) genJSONTypes(ctx broodmother.Context) string {
	t := template.Must(template.New("httpgen").Parse(ioJSONType))
	type field struct {
		Name string
		Type string
		Tag  string
	}
	type strct struct {
		Name   string
		Fields []field
	}
	requests := make([]strct, 0)
	responses := make([]strct, 0)
	imports := make(map[string]string)

	// TODO: we should use go/build package for
	// checking imports.
	// Now, this code assumes that folder name and
	// package name are same thing.
	checkForImport := func(inp string) {
		splited := strings.SplitN(inp, ".", 2)
		if len(splited) <= 1 {
			return
		}
		pkg := splited[0]
		for p, n := range ctx.Imports() {
			if n == "" {
				n = path.Base(p)
			}
			if n == pkg {
				imports[p] = n
			}
		}
	}
	stringImports := func() []string {
		outp := make([]string, 0)
		for p, n := range imports {
			if path.Base(p) == n {
				outp = append(outp, fmt.Sprintf("%q", p))
			} else {
				outp = append(outp, fmt.Sprintf("%s %q", n, p))
			}
		}
		return outp
	}

	for _, f := range g.funcs {
		if len(f.params) > 0 {
			fields := make([]field, 0)
			for _, p := range f.params {
				checkForImport(p.typ)
				fields = append(fields, field{
					Name: broodmother.CamelCase(p.name, true),
					Type: p.typ,
					Tag:  `json:"` + broodmother.SnakeCase(p.name) + `"`,
				})
			}
			requests = append(requests, strct{
				Name:   f.name,
				Fields: fields,
			})
		}
		if len(f.result) > 0 {
			fields := make([]field, 0)
			for _, r := range f.result {
				checkForImport(r.typ)
				fields = append(fields, field{
					Name: broodmother.CamelCase(r.name, true),
					Type: r.typ,
					Tag:  `json:"` + broodmother.SnakeCase(r.name) + `"`,
				})
			}
			responses = append(responses, strct{
				Name:   f.name,
				Fields: fields,
			})
		}
	}
	buf := &bytes.Buffer{}
	t.Execute(buf, struct {
		Package   string
		Imports   []string
		Requests  []strct
		Responses []strct
	}{
		Package:   broodmother.CamelCase(g.name, false),
		Imports:   stringImports(),
		Requests:  requests,
		Responses: responses,
	})
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formatted)
}
