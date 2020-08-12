package inspect

import (
	"fmt"
	"go/ast"

	"github.com/Hu13er/broodmother"
)

var (
	Name                       = "inspector"
	inspectTag broodmother.Tag = "inspect"
)

type Inspector struct {
	counter int
}

var _ broodmother.Generator = (*Inspector)(nil)

func (*Inspector) Name() string {
	return Name
}

func (*Inspector) Filter() broodmother.Filterer {
	return broodmother.FilterNothing
}

func (hg *Inspector) Visit(ctx broodmother.Context, node ast.Node) (bool, error) {
	if node == nil {
		hg.counter--
		return false, nil
	}
	hg.counter++
	value, _ := ctx.Get(inspectTag)
	if s, _ := value.(string); s == "on" {
		fmt.Printf("%d]\t(%T)\t%+v\n", hg.counter, node, node)
	}
	return true, nil
}
