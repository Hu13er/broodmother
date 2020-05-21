package broodmother

import (
	"fmt"
	"go/build"
	"path"
	"strings"
)

func GetPackage(p string) (pkgName, importPath string, err error) {
	if path.Ext(p) != "" {
		p, _ = path.Split(p)
	}
	p = strings.TrimSuffix(p, "/")
	p = strings.TrimPrefix(p, build.Default.GOPATH+"/src/")
	pkg, err := build.Default.Import(p, "", build.ImportComment)
	if err != nil {
		return "", "", err
	}
	pkgName = pkg.Name
	importPath = pkg.ImportPath
	return pkgName, importPath, nil
}

func SnakeCase(camel string) string {
	isUpper := func(c byte) bool {
		return c >= 'A' && c <= 'Z'
	}
	outp := ""
	camel += "!"
	for i := range camel {
		if i > 0 && i < len(camel)-1 {
			if isUpper(camel[i]) && !isUpper(camel[i+1]) {
				outp += "-"
			}
		}
		outp += string(camel[i])
	}
	outp = outp[0 : len(outp)-1]
	return strings.ToLower(outp)
}

func JoinCamel(name ...string) string {
	outp := ""
	for _, n := range name {
		outp += CamelCase(n, true)
	}
	return outp
}

func CamelCase(snake string, first bool) string {
	outp := ""
	nextCapital := first
	for _, c := range snake {
		if c == '-' {
			nextCapital = true
			continue
		}
		if nextCapital {
			outp += strings.ToUpper(string(c))
			nextCapital = false
		} else {
			outp += string(c)
		}
	}
	return outp
}

func NthString(n int) string {
	lst := []string{"first", "second", "third", "forth"}
	if n-1 < 4 {
		return lst[n-1]
	}
	return fmt.Sprintf("%dth", n)
}
