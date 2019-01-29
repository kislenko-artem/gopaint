// list of primitives

package primitives

import "github.com/gotk3/gotk3/cairo"

type Primitive interface {
	SetColor(cr *cairo.Context)
	Draw(cr *cairo.Context)
	SetStart(x, y float64)
	SetStop(x, y float64)
	Release()
	IsWait() bool
}
