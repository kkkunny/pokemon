package imgutil

import (
	"image"
	"image/color"
	"image/draw"
)

type Image interface {
	draw.RGBA64Image
	SubImage(r image.Rectangle) Image
	Fill(c color.Color)
	Scale(x, y float64) Image
}
