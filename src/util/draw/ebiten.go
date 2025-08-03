package draw

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	image2 "github.com/kkkunny/pokemon/src/util/image"
)

type ebitenImage struct {
	_options
	img *ebiten.Image
}

func NewDrawerFromEbiten(img *ebiten.Image) Drawer {
	drawer := &ebitenImage{img: img}
	drawer._options = newOptions(drawer)
	return drawer
}

func (d *ebitenImage) copyWithOptions(opts _options) Drawer {
	return &ebitenImage{
		_options: opts,
		img:      d.img,
	}
}

func (d *ebitenImage) Size() (width int, height int) {
	bounds := d.img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func (d *ebitenImage) Set(x int, y int, clr color.Color) {
	d.img.Set(x, y, clr)
}

func (d *ebitenImage) DrawImage(dst image.Image) error {
	var ebitenImg *ebiten.Image
	if i, ok := dst.(*ebiten.Image); ok {
		ebitenImg = i
	} else if i, ok := dst.(*image2.Image); ok {
		ebitenImg = i.Image
	} else {
		ebitenImg = ebiten.NewImageFromImage(dst)
	}

	var imageOpts ebiten.DrawImageOptions
	imageOpts.GeoM.Scale(d.scaleX, d.scaleX)
	imageOpts.GeoM.Translate(float64(d.x), float64(d.y))
	if d.scaleWithColor != nil {
		imageOpts.ColorScale.ScaleWithColor(d.scaleWithColor)
	}
	d.img.DrawImage(ebitenImg, &imageOpts)
	return nil
}

func (d *ebitenImage) DrawText(renderText string, fontFace text.Face) error {
	var textOpts text.DrawOptions
	textOpts.GeoM.Scale(d.scaleX, d.scaleX)
	textOpts.GeoM.Translate(float64(d.x), float64(d.y))
	if d.scaleWithColor != nil {
		textOpts.ColorScale.ScaleWithColor(d.scaleWithColor)
	}
	text.Draw(d.img, renderText, fontFace, &textOpts)
	return nil
}
