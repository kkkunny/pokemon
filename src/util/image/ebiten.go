package imgutil

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ebitenImage struct {
	Image *ebiten.Image
}

func WrapImage(img image.Image) Image {
	if imgInst, ok := img.(*ebitenImage); ok {
		return imgInst
	} else if imgInst, ok := img.(*ebiten.Image); ok {
		return &ebitenImage{Image: imgInst}
	} else {
		return &ebitenImage{Image: ebiten.NewImageFromImage(img)}
	}
}

func NewImageFromFile(path string) (Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	return WrapImage(img), nil
}

func NewImage(w, h int) Image {
	return WrapImage(ebiten.NewImage(w, h))
}

func (i *ebitenImage) EbitenImage() *ebiten.Image {
	return i.Image
}

func (i *ebitenImage) ColorModel() color.Model {
	return i.Image.ColorModel()
}

func (i *ebitenImage) Bounds() image.Rectangle {
	return i.Image.Bounds()
}

func (i *ebitenImage) At(x int, y int) color.Color {
	return i.Image.At(x, y)
}

func (i *ebitenImage) RGBA64At(x, y int) color.RGBA64 {
	return i.Image.RGBA64At(x, y)
}

func (i *ebitenImage) Set(x int, y int, c color.Color) {
	i.Image.Set(x, y, c)
}

func (i *ebitenImage) SetRGBA64(x, y int, c color.RGBA64) {
	i.Image.Set(x, y, c)
}

func (i *ebitenImage) SubImage(r image.Rectangle) Image {
	return WrapImage(i.Image.SubImage(r))
}

func (i *ebitenImage) Fill(c color.Color) {
	i.Image.Fill(c)
}
