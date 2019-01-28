package color

type rgb struct {
	R float64
	G float64
	B float64
}

type Color struct {
	RGB rgb
}

func (c *Color) PickColor(newColor *Color) {
	c = newColor
}
