package httpgen

import (
	"bytes"
	"strings"
	"text/template"

	"gitlab.com/pirates1/broodmother"
)

var httpServer = `
package {{.Package}}

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	{{range .Imports}}
	"{{.}}"{{end}}
)

type {{.StructName}} struct {
	Core {{.CoreType}}
}

func (h *{{$.StructName}}) Handler() http.Handler {
	r := mux.NewRouter()
	{{range .Methods}}r.HandleFunc("{{.Path}}", h.handle{{.Name}}).Methods("POST")
	{{end}}
	return r
}

{{range .Methods}}
func (h *{{$.StructName}}) handle{{.Name}}(w http.ResponseWriter, r *http.Request) { {{if .Params }}
	var request {{.Name}}JSONRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	{{end}} {{if .Results }}
	{{call $.Join ", " .Results}} := h.Core.{{.Name}}(
		{{call $.AllCamelCase .Params | call $.AllPrefix "request." | call $.Join ", "}})
	response := {{.Name}}JSONResponse { {{range .Results}}
		{{call $.CamelCase .}}: {{.}},{{end}}
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusServerInternalError)
		return
	}
	{{else}}
		h.Core.{{.Name}}(
			{{call $.Join .Params ","}})
	{{end}}
}
{{end}}
`

func (g *HttpGen) genHttpServer(ctx broodmother.Context) string {
	t := template.Must(template.New("httpgen-httpserver").Parse(httpServer))
	type method struct {
		Name    string
		Path    string
		Params  []string
		Results []string
	}
	methods := make([]method, 0)
	for _, f := range g.funcs {
		m := method{}
		m.Name = f.name
		m.Path = "/" + broodmother.SnakeCase(f.name)
		for _, p := range f.params {
			m.Params = append(m.Params, p.name)
		}
		for _, r := range f.result {
			m.Results = append(m.Results, r.name)
		}
		methods = append(methods, m)
	}
	buf := &bytes.Buffer{}
	t.Execute(buf, struct {
		Package    string
		Imports    []string
		StructName string
		CoreType   string
		Methods    []method

		CamelCase    func(string) string
		AllCamelCase func([]string) []string
		AllPrefix    func(string, []string) []string
		Join         func(string, []string) string
	}{
		Package:    broodmother.CamelCase(g.name, false),
		Imports:    []string{"holyshit"},
		StructName: "HttpServer",
		CoreType:   "slardar.Interface",
		Methods:    methods,

		CamelCase: func(s string) string {
			return broodmother.CamelCase(s, true)
		},
		AllCamelCase: func(ary []string) []string {
			outp := make([]string, len(ary))
			for i, s := range ary {
				outp[i] = broodmother.CamelCase(s, true)
			}
			return outp
		},
		AllPrefix: func(prefix string, ary []string) []string {
			outp := make([]string, len(ary))
			for i, s := range ary {
				outp[i] = prefix + s
			}
			return outp
		},
		Join: func(sep string, ary []string) string {
			return strings.Join(ary, sep)
		},
	})
	return buf.String()
}
