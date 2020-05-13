package inspect

import (
	"context"
	"fmt"
	"go/ast"

	"gitlab.com/pirates1/broodmother"
)

var (
	Name       = "inspector"
	inspectTag = "inspect"
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

type inspectPowerKeyType string

var inspectPowerKey inspectPowerKeyType = "inspect-power"

func (hg *Inspector) Visit(ctx context.Context, node ast.Node) (bool, context.Context) {
	if node == nil {
		hg.counter--
		return false, ctx
	}
	tags := broodmother.ParseDocumentTags(node)
	if tags.Has(inspectTag) {
		ctx = context.WithValue(ctx,
			inspectPowerKey, tags.GetBool(inspectTag))
	}

	hg.counter++
	if ok, _ := ctx.Value(inspectPowerKey).(bool); ok {
		fmt.Printf("%d]\t(%T)\t%s\t(%d, %d)\n",
			hg.counter, node, node, node.Pos(), node.End())
	}
	return true, ctx
}
