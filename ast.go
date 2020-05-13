package broodmother

import "go/ast"

func Document(node ast.Node) *ast.CommentGroup {
	switch typed := node.(type) {
	case *ast.Field:
		return typed.Doc
	case *ast.ImportSpec:
		return typed.Doc
	case *ast.ValueSpec:
		return typed.Doc
	case *ast.TypeSpec:
		return typed.Doc
	case *ast.GenDecl:
		return typed.Doc
	case *ast.File:
		return typed.Doc
	}
	return nil
}

func Comment(node ast.Node) *ast.CommentGroup {
	switch typed := node.(type) {
	case *ast.Field:
		return typed.Comment
	case *ast.ImportSpec:
		return typed.Comment
	case *ast.ValueSpec:
		return typed.Comment
	case *ast.TypeSpec:
		return typed.Comment
	}
	return nil
}
