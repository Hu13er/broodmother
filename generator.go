package broodmother

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Generator interface {
	Name() string
	Filter() Filterer
	Visit(Context, ast.Node) (bool, error)
}

type Finalezer interface {
	Finalize(Context) ([]File, error)
}

type Executor struct {
	Generators []Generator
}

func (e *Executor) ParseFile(path string) error {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(
		fs, path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}
	imps := make(map[string]string)
	for _, imp := range f.Imports {
		p := strings.Trim(imp.Path.Value, "\"")
		n := ""
		if imp.Name != nil {
			n = imp.Name.String()
		}
		imps[p] = n
	}
	ctx := newBackground(f.Name.String(), path, imps)
	for _, g := range e.Generators {
		errs := make([]error, 0)
		ast.Inspect(f, func(node ast.Node) bool {
			if node == nil {
				g.Visit(ctx, nil)
				ctx = ctx.Parent()
				return false
			}
			ctx = newLayer(ctx)
			for k, v := range ParseDocumentTags(node) {
				ctx.Set(Tag(k), v)
			}
			if f := g.Filter(); f != nil {
				if !f.Allowed(ctx, node) {
					return true
				}
			}
			cont, err := g.Visit(ctx, node)
			if err != nil {
				errs = append(errs, err)
			}
			return cont
		})
		if len(errs) <= 0 {
			if f, ok := g.(Finalezer); ok {
				fs, err := f.Finalize(ctx)
				if err != nil {
					fmt.Println("[+] Generator", g.Name(), "error:", err)
				} else {
					for _, f := range fs {
						err := f.write()
						if err != nil {
							fmt.Println("[+] Generator", g.Name(),
								"writing file error:", err)
						}
					}
				}
			}
		} else {
			fmt.Println("[+] Generator", g.Name(), "errors:")
			for _, err := range errs {
				fmt.Println("\t[*]", err)
			}
		}
	}
	return nil
}
