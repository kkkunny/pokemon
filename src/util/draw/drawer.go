package draw

import "image/draw"

type OptionDrawer interface {
	draw.Image
	options
}

type _optionDrawer struct {
	draw.Image
	options
}

func NewDrawerFromImage(img draw.Image) OptionDrawer {
	drawer := _optionDrawer{Image: img}
	opts := newOptions()
	opts.drawer = drawer
	drawer.options = opts
	return drawer
}
