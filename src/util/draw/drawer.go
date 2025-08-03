package draw

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Drawer interface {
	options
	copyWithOptions(opts _options) Drawer
	Size() (width int, height int)
	Set(x int, y int, clr color.Color)
	DrawImage(image image.Image) error
	DrawText(renderText string, fontFace text.Face) error
	OverlayColor(c color.Color) error
}
