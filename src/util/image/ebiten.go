package imgutil

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Image struct {
	Image *ebiten.Image
}

func WrapImage(img image.Image) *Image {
	if imgInst, ok := img.(*Image); ok {
		return imgInst
	} else if imgInst, ok := img.(*ebiten.Image); ok {
		return &Image{Image: imgInst}
	} else {
		return &Image{Image: ebiten.NewImageFromImage(img)}
	}
}

func NewImageFromFile(path string) (*Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	return WrapImage(img), nil
}

func NewImage(w, h int) *Image {
	return WrapImage(ebiten.NewImage(w, h))
}

func (i *Image) EbitenImage() *ebiten.Image {
	return i.Image
}

func (i *Image) ColorModel() color.Model {
	return i.Image.ColorModel()
}

func (i *Image) Bounds() image.Rectangle {
	return i.Image.Bounds()
}

func (i *Image) At(x int, y int) color.Color {
	return i.Image.At(x, y)
}

func (i *Image) RGBA64At(x, y int) color.RGBA64 {
	return i.Image.RGBA64At(x, y)
}

func (i *Image) Set(x int, y int, c color.Color) {
	i.Image.Set(x, y, c)
}

func (i *Image) SetRGBA64(x, y int, c color.RGBA64) {
	i.Image.Set(x, y, c)
}

func (i *Image) SubImage(r image.Rectangle) *Image {
	return WrapImage(i.Image.SubImage(r))
}
