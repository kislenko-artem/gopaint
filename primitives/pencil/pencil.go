package pencil

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/kislenko-artem/gopaint/property/color"
)

type Pencil struct {
	color.Color
	XStart float64
	YStart float64
	XEnd   float64
	YEnd   float64

	start  bool
	points [][]float64
}

func New(color color.Color) *Pencil {
	return &Pencil{Color: color}
}

func (p *Pencil) IsWait() bool {
	return p.start
}

func (p *Pencil) Release() {
	p.start = false
}

func (p *Pencil) SetStart(x, y float64) {
	p.start = true
	p.points = append(p.points, []float64{x, y})
}

func (p *Pencil) SetStop(x, y float64) {
	p.points = append(p.points, []float64{x, y})
}

func (p *Pencil) SetColor(cr *cairo.Context) {
	cr.SetSourceRGB(p.RGB.R, p.RGB.G, p.RGB.B)
}

func (p *Pencil) Draw(cr *cairo.Context) {
	cr.SetLineWidth(0.5)
	cr.Scale(1, 1)
	if len(p.points) == 0 {
		return
	}
	//cr.MoveTo(p.points[0][0], p.points[0][1]);
	if len(p.points) < 3 {
		return
	}
	for i := 1; i < len(p.points); i += 3 {
		if i == 0 {
			continue
		}
		if len(p.points) < i-1 {
			break
		}
		if len(p.points) < i+1 {
			break
		}
		if len(p.points) <= i+2 {
			break
		}
		cr.CurveTo(p.points[i][0], p.points[i][1], p.points[i+1][0], p.points[i+1][1], p.points[i+2][0], p.points[i+2][1])
	}

	cr.Stroke()
}
