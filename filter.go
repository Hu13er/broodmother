package broodmother

import (
	"go/ast"
)

type (
	Filterer interface {
		Allowed(ctx Context, node ast.Node) bool
	}

	FilterFunc func(ctx Context, node ast.Node) bool
	FilterList []Filterer
	FilterAny  []Filterer
	FilterTags []string
)

var (
	// Ensure implementing Filterer
	_ Filterer = (FilterFunc)(nil)
	_ Filterer = (FilterList)(nil)

	// Some Singletons
	FilterAll     = FilterFunc(func(Context, ast.Node) bool { return false })
	FilterNothing = FilterFunc(func(Context, ast.Node) bool { return true })
)

func (f FilterFunc) Allowed(ctx Context, node ast.Node) bool {
	return f(ctx, node)
}

func (fl FilterList) Allowed(ctx Context, node ast.Node) bool {
	for _, f := range fl {
		if !f.Allowed(ctx, node) {
			return false
		}
	}
	return true
}

func (fa FilterAny) Allowed(ctx Context, node ast.Node) bool {
	for _, f := range fa {
		if f.Allowed(ctx, node) {
			return true
		}
	}
	return false
}

func (tf FilterTags) Allowed(ctx Context, node ast.Node) bool {
	for _, t := range tf {
		if _, exists := ctx.Get(Tag(t)); !exists {
			return false
		}
	}
	return true
}
