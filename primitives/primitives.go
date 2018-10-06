package primitives

type Line interface {
	IsWait() bool
	SetStart(x, y float64)
	SetEnd(x, y float64)
	Release()
	GetStart() (x float64, y float64)
	GetEnd() (x float64, y float64)
}
