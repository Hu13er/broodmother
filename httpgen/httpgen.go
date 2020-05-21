package httpgen

import (
	"errors"
	"fmt"
	"go/ast"
	"path"

	"gitlab.com/pirates1/broodmother"
)

type funcDef struct {
	name   string
	params []varDef
	result []varDef
}

type varDef struct {
	name string
	typ  string
}

type HttpGen struct {
	name         string
	serverPath   string
	structName   string
	corePkg      string
	coreImport   string
	coreTypeName string
	funcs        []funcDef
}

var (
	_ broodmother.Generator = (*HttpGen)(nil)
)

func (g *HttpGen) Name() string {
	return "httpgen"
}

func (g *HttpGen) Filter() broodmother.Filterer {
	return broodmother.FilterList{
		broodmother.FilterTags{"httpgen"},
		broodmother.FilterFunc(func(ctx broodmother.Context, node ast.Node) bool {
			if ts, ok := node.(*ast.TypeSpec); ok {
				if _, ok := ts.Type.(*ast.InterfaceType); ok {
					return true
				}
			}
			return false
		}),
	}
}

func (g *HttpGen) Visit(ctx broodmother.Context, node ast.Node) (bool, error) {
	if node == nil {
		return false, nil
	}

	var err error
	g.corePkg, g.coreImport, err =
		broodmother.GetPackage(ctx.Path())
	if err != nil {
		return false, err
	}
	g.coreTypeName = node.(*ast.TypeSpec).Name.String()
	g.name = node.(*ast.TypeSpec).Name.String()
	if name, exists := ctx.
		Get(broodmother.Tag("httpgen.name")); exists {
		g.name = broodmother.CamelCase(name.(string), true)
	}
	g.serverPath = path.Join(path.Dir(ctx.Path()), "http")
	if sp, exists := ctx.
		Get(broodmother.Tag("httpgen.server-path")); exists {
		g.serverPath = path.Join(ctx.Path(), sp.(string))
	}
	g.structName = "HttpServer"
	if sn, exists := ctx.
		Get(broodmother.Tag("httpgen.struct-name")); exists {
		g.structName = sn.(string)
	}

	iface := node.(*ast.TypeSpec).Type.(*ast.InterfaceType)
	for _, method := range iface.Methods.List {
		f := funcDef{}
		f.name = method.Names[0].String()
		ftype := method.Type.(*ast.FuncType)
		for _, p := range ftype.Params.List {
			f.params = append(f.params, parseVarList(p)...)
		}
		for _, p := range f.params {
			if p.name == "" {
				// TODO: Better errors messages.
				return false, errors.New("empty name for params")
			}
		}
		for _, r := range ftype.Results.List {
			f.result = append(f.result, parseVarList(r)...)
		}
		for _, r := range f.result {
			if r.name == "" {
				// TODO: Better errors messages.
				return false, errors.New("empty name for results")
			}
		}
		g.funcs = append(g.funcs, f)
	}
	return true, nil
}

func (g *HttpGen) Finalize(ctx broodmother.Context) ([]broodmother.File, error) {
	jsons := g.genJSONTypes(ctx)
	fmt.Println(jsons)
	httpserver := g.genHttpServer(ctx)
	fmt.Println(httpserver)
	client := g.genClient(ctx)
	fmt.Println(client)
	return []broodmother.File{
		{
			Path:    path.Join(g.serverPath, "httpserver.httpgen.go"),
			Content: httpserver,
		},
		{
			Path:    path.Join(g.serverPath, "jsons.httpgen.go"),
			Content: jsons,
		},
		{
			Path:    path.Join(g.serverPath, "/client/client.httpgen.go"),
			Content: client,
		},
	}, nil
}

func parseVarList(lst *ast.Field) []varDef {
	outp := make([]varDef, 0)
	typ := ""
	switch typed := lst.Type.(type) {
	case *ast.Ident:
		typ = typed.String()
	case *ast.SelectorExpr:
		typ = typed.X.(*ast.Ident).String() + "." + typed.Sel.Name
	}
	for _, n := range lst.Names {
		outp = append(outp, varDef{n.String(), typ})
	}
	if len(lst.Names) <= 0 {
		outp = append(outp, varDef{typ: typ})
	}
	return outp
}
