package util

import "image/color"

func NewRGBColor(r, g, b uint8) *color.RGBA {
	return &color.RGBA{R: r, G: g, B: b, A: 0xff}
}

func NewRGBAColor(r, g, b, a uint8) *color.RGBA {
	return &color.RGBA{R: r, G: g, B: b, A: a}
}
