package line

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/kislenko-artem/gopaint/property/color"
)

type Line struct {
	color.Color
	XStart float64
	YStart float64
	XEnd   float64
	YEnd   float64

	start bool
}

func New(color color.Color) *Line {
	return &Line{Color: color}
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
	cr.SetSourceRGB(l.RGB.R, l.RGB.G, l.RGB.B)
}

func (l *Line) Draw(cr *cairo.Context) {
	cr.MoveTo(l.XStart, l.YStart)
	cr.LineTo(l.XEnd, l.YEnd)
	cr.Stroke()
}
