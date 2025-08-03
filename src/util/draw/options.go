package draw

import "image/color"

type options interface {
	Move(x, y int) Drawer
	Scale(x, y float64) Drawer
	ScaleWithColor(c color.Color) Drawer
}

type _options struct {
	drawer Drawer

	x, y           int
	scaleX, scaleY float64
	scaleWithColor color.Color
}

func newOptions(drawer Drawer) _options {
	return _options{
		drawer: drawer,
		scaleX: 1.0,
		scaleY: 1.0,
	}
}

func (d _options) Move(x, y int) Drawer {
	d.x += x
	d.y += y
	return d.drawer.copyWithOptions(d)
}

func (d _options) Scale(x, y float64) Drawer {
	d.scaleX *= x
	d.scaleY *= y
	return d.drawer.copyWithOptions(d)
}

func (d _options) ScaleWithColor(c color.Color) Drawer {
	d.scaleWithColor = c
	return d.drawer.copyWithOptions(d)
}
