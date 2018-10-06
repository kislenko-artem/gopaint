package line

import "../../primitives"

type Line struct {
	XStart float64
	YStart float64
	XEnd   float64
	YEnd   float64

	start bool
}

func New() primitives.Line {
	var line primitives.Line = &Line{}
	return line
}

func (l *Line) IsWait() bool {
	return l.start
}

func (l *Line) Release() {
	l.start = false
}

func (l *Line) SetStart(x, y float64) {
	l.XStart = x
	l.YStart = y
	l.start = true
}

func (l *Line) SetEnd(x, y float64) {
	l.XEnd = x
	l.YEnd = y
}

func (l *Line) GetStart() (x float64, y float64) {
	return l.XStart, l.YStart
}

func (l *Line) GetEnd() (x float64, y float64) {
	return l.XEnd, l.YEnd
}
