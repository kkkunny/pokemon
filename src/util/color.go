package util

import "image/color"

func NewRGBColor(r, g, b uint8) *color.NRGBA {
	return NewRGBAColor(r, g, b, 0xff)
}

func NewRGBAColor(r, g, b, a uint8) *color.NRGBA {
	return &color.NRGBA{R: r, G: g, B: b, A: a}
}
