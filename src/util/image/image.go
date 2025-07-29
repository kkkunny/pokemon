package image

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Image struct {
	*ebiten.Image
}

func WrapImage(img *ebiten.Image) *Image {
	return &Image{Image: img}
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

func (i *Image) Width() int {
	return i.Image.Bounds().Dx()
}

func (i *Image) Height() int {
	return i.Image.Bounds().Dy()
}

func (i *Image) SubImage(x, y, w, h int) *Image {
	return i.SubImageByRect(image.Rect(x, y, x+w, y+h))
}

func (i *Image) SubImageByRect(r image.Rectangle) *Image {
	return WrapImage(i.Image.SubImage(r).(*ebiten.Image))
}

func (i *Image) DrawImage(img *Image, options *ebiten.DrawImageOptions) {
	i.Image.DrawImage(img.Image, options)
}

func (i *Image) Overlay(c color.Color) {
	mask := NewImage(i.Width(), i.Height())
	mask.Fill(c)
	i.DrawImage(mask, nil)
}
