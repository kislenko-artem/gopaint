package pencil

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/kislenko-artem/gopaint/property/color"
	"log"
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
	for _, point := range p.points {
		log.Println(point[0], point[1])
		cr.Rectangle(point[0], point[1], 1, 1)
		cr.Fill()
	}
}
