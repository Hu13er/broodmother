package broodmother

import (
	"go/ast"
	"strings"
)

var (
	TagsPrefix = "+bm:"
	trues      = []string{"t", "true", "on", "yes"}
)

type Tag string

func ParseDocumentTags(node ast.Node) map[string]string {
	if cg := Document(node); cg != nil {
		return ParseCommentGroupTags(cg)
	}
	return map[string]string{}
}

func ParseCommentTags(node ast.Node) map[string]string {
	if cg := Comment(node); cg != nil {
		return ParseCommentGroupTags(cg)
	}
	return map[string]string{}
}

func ParseCommentGroupTags(cg *ast.CommentGroup) map[string]string {
	if cg == nil {
		return map[string]string{}
	}
	tags := make(map[string]string)
	for _, cm := range cg.List {
		txt := trimComments(cm.Text)
		if !strings.HasPrefix(txt, TagsPrefix) {
			continue
		}
		txt = strings.TrimPrefix(txt, TagsPrefix)

		for _, expr := range strings.Split(txt, ";") {
			var (
				sliced = strings.SplitN(expr, "=", 2)
				key    = sliced[0]
				value  = ""
			)
			if len(sliced) >= 2 {
				value = sliced[1]
			}
			tags[key] = value
		}
	}
	return tags
}

func trimComments(cmnt string) string {
	cmnt = strings.TrimPrefix(cmnt, "//")
	cmnt = strings.TrimPrefix(cmnt, "/*")
	cmnt = strings.TrimSuffix(cmnt, "*/")
	return cmnt
}
