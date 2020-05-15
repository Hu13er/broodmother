// +build ignore

package broodmother

import (
	"shit/foo"
)

//+bm:inspect=on
//+bm:httpgen;httpgen.name=helllllooo
type Shit interface {
	Hello(a string, b string) (greeting string)
	Hello2(a, b string, c string) (greeting error)
	Hello3(a, b string, c string) (greeting foo.Tag)
}
