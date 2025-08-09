package draw

type options interface {
	Options() _options
	Reset() OptionDrawer
	Move(x, y int) OptionDrawer
	MoveTo(x, y int) OptionDrawer
	Scale(x, y float64) OptionDrawer
	SetScale(x, y float64) OptionDrawer
}

type _options struct {
	drawer _optionDrawer

	x, y           int
	scaleX, scaleY float64
}

func newOptions() _options {
	return _options{
		scaleX: 1.0,
		scaleY: 1.0,
	}
}

func (d _options) withOptions(opts options) OptionDrawer {
	newOpts := opts.(_options)
	newOpts.drawer = d.drawer
	newOpts.drawer.options = newOpts
	return newOpts.drawer
}

func (d _options) Options() _options {
	return d
}

func (d _options) Reset() OptionDrawer {
	return d.withOptions(newOptions())
}

func (d _options) Move(x, y int) OptionDrawer {
	d.x += int(float64(x) * d.scaleX)
	d.y += int(float64(y) * d.scaleX)
	return d.withOptions(d)
}

func (d _options) MoveTo(x, y int) OptionDrawer {
	d.x = int(float64(x) * d.scaleX)
	d.y = int(float64(y) * d.scaleX)
	return d.withOptions(d)
}

func (d _options) Scale(x, y float64) OptionDrawer {
	d.scaleX *= x
	d.scaleY *= y
	return d.withOptions(d)
}

func (d _options) SetScale(x, y float64) OptionDrawer {
	d.scaleX = x
	d.scaleY = y
	return d.withOptions(d)
}
