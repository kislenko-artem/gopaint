package line

import "github.com/gotk3/gotk3/cairo"

type Line struct {
	XStart float64
	YStart float64
	XEnd   float64
	YEnd   float64

	start bool
}

func New() *Line {
	return &Line{}
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

func (l *Line) SetStop(x, y float64) {
	l.XEnd = x
	l.YEnd = y
}

func (l *Line) SetColor(cr *cairo.Context) {
	cr.SetSourceRGB(0, 0, 0)
}

func (l *Line) Draw(cr *cairo.Context) {
	cr.MoveTo(l.XStart, l.YStart)
	cr.LineTo(l.XEnd, l.YEnd)
	cr.Stroke()
}
