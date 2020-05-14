package broodmother

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Generator interface {
	Name() string
	Filter() Filterer
	Visit(Context, ast.Node) bool
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
	// ctxStack := contextStack{context.Background()}
	ctx := newBackground(path)
	for _, g := range e.Generators {
		ast.Inspect(f, func(node ast.Node) bool {
			if node == nil {
				g.Visit(ctx, nil)
				ctx = ctx.Parent()
				// g.Visit(ctxStack.pop(), node)
				return false
			}
			if f := g.Filter(); f != nil {
				// if !f.Allowed(ctxStack.top(), node) {
				if !f.Allowed(ctx, node) {
					return false
				}
			}
			// cont, ctx := g.Visit(ctxStack.top(), node)
			ctx = newLayer(ctx)
			for k, v := range ParseDocumentTags(node) {
				ctx.Set(Tag(k), v)
			}
			return g.Visit(ctx, node)
			// cont := g.Visit(ctx, node)
			// ctxStack.push(ctx)
			// return cont
		})
	}
	return nil
}

// type contextStack []context.Context

// func (s *contextStack) push(ctx context.Context) {
// 	(*s) = append(*s, ctx)
// }

// func (s *contextStack) top() context.Context {
// 	return (*s)[len(*s)-1]
// }

// func (s *contextStack) pop() context.Context {
// 	top := s.top()
// 	(*s) = (*s)[0 : len(*s)-1]
// 	return top
// }

// func (e *Executor) ParseDir(path string) error {
// 	fs := token.NewFileSet()
// 	pkgs, err := parser.ParseDir(
// 		fs, path, nil, parser.AllErrors|parser.ParseComments)
// 	if err != nil {
// 		return err
// 	}
// 	for file, pkg := range pkgs {
// 		for _, g := range e.Generators {
// 			ast.Inspect(pkg, func(node ast.Node) bool {
// 				if node == nil {
// 					g.Visit(path, node)
// 					return false
// 				}
// 				if f := g.Filter(); f != nil {
// 					if !f.Allowed(node) {
// 						return false
// 					}
// 				}
// 				return g.Visit(file, node)
// 			})
// 		}
// 	}
// 	return nil
// }
