package broodmother

import (
	"go/ast"
	"strings"
)

var (
	TagsPrefix = "+bm:"
	trues      = []string{"t", "true", "on", "yes"}
)

type Tags map[string]string

func ParseDocumentTags(node ast.Node) Tags {
	if cg := Document(node); cg != nil {
		return ParseCommentGroupTags(cg)
	}
	return make(Tags)
}

func ParseCommentTags(node ast.Node) Tags {
	if cg := Comment(node); cg != nil {
		return ParseCommentGroupTags(cg)
	}
	return make(Tags)
}

func ParseCommentGroupTags(cg *ast.CommentGroup) Tags {
	if cg == nil {
		return make(Tags)
	}
	tagss := make(Tags)
	for _, cm := range cg.List {
		txt := trimComments(cm.Text)
		if !strings.HasPrefix(txt, TagsPrefix) {
			continue
		}
		txt = strings.TrimPrefix(txt, TagsPrefix)

		tags := make(Tags)
		for _, expr := range strings.Split(txt, ";") {
			var (
				sliced = strings.SplitN(expr, "=", 2)
				key    = sliced[0]
				value  = ""
			)
			if len(sliced) >= 2 {
				value = sliced[1]
			}
			tags.Set(key, value)
		}
		tagss = tagss.Join(tags)
	}
	return tagss
}

func (t Tags) Copy() Tags {
	cpy := make(Tags)
	for k, v := range t {
		cpy[k] = v
	}
	return cpy
}

func (t Tags) Get(key string) string {
	return t[key]
}

func (t Tags) Set(key, value string) {
	t[key] = value
}

func (t Tags) GetBool(key string) bool {
	value := t.Get(key)
	for _, tr := range trues {
		if value == tr {
			return true
		}
	}
	return false
}

func (t Tags) Has(key string) bool {
	_, exists := t[key]
	return exists
}

func (t Tags) Add(key, value string) Tags {
	cpy := t.Copy()
	cpy[key] = value
	return cpy
}

func (t Tags) Del(key string) Tags {
	cpy := t.Copy()
	delete(cpy, key)
	return cpy
}

func (t Tags) Join(t2 Tags) Tags {
	cpy := t.Copy()
	for k, v := range t2 {
		cpy[k] = v
	}
	return cpy
}

func trimComments(cmnt string) string {
	cmnt = strings.TrimPrefix(cmnt, "//")
	cmnt = strings.TrimPrefix(cmnt, "/*")
	cmnt = strings.TrimSuffix(cmnt, "*/")
	return cmnt
}
