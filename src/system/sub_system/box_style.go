package sub_system

import (
	"image"

	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type BoxStyle uint8

var BoxStyleEnum = enum.New[struct {
	Label    BoxStyle
	Dialogue BoxStyle
	Battle   BoxStyle
}]()

func (s BoxStyle) OnUpdate(drawer draw.OptionDrawer, rect image.Rectangle) (image.Rectangle, image.Rectangle, error) {
	x, y, w, h := rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy()
	switch s {
	case BoxStyleEnum.Label:
		sepX1, sepY1, sepX2, sepY2, sepX3, sepY3 := 10, 10, 6, 4, 20, 4
		if w <= 0 {
			w = drawer.Bounds().Dx() - sepX1*2 - sepX2*2 - sepX3*2
		}
		if h <= 0 {
			h = 90
		}
		if x == -1 {
			x = sepX1
		}
		if y == -1 {
			y = drawer.Bounds().Dy() - sepY1 - h - sepY2*2 - sepY3*2
		}
		draw.PrepareDrawRect(drawer, w+sepX2*2+sepX3*2, h+sepY2*2+sepY3*2, util.NewNRGBColor(104, 112, 120)).SetRadius(6).Move(x, y).Draw()
		draw.PrepareDrawRect(drawer, w+sepX3*2, h+sepY3*2, util.NewNRGBColor(200, 200, 216)).Move(x+sepX2, y+sepY2).SetRadius(4).Draw()
		draw.PrepareDrawRect(drawer, w, h, util.NewNRGBColor(248, 248, 248)).Move(x+sepX2+sepX3, y+sepY2+sepY3).Draw()
		return imgutil.Rect(x, y, w+sepX2*2+sepX3*2, h+sepY2*2+sepY3*2), imgutil.Rect(x+sepX2+sepX3, y+sepY2+sepY3, w, h), nil
	case BoxStyleEnum.Dialogue:
		sepX1, sepY1, sepX2, sepY2, sepX3, sepY3 := 10, 10, 6, 4, 20, 4
		if w <= 0 {
			w = drawer.Bounds().Dx() - sepX1*2 - sepX2*2 - sepX3*2
		}
		if h <= 0 {
			h = 90
		}
		if x == -1 {
			x = sepX1
		}
		if y == -1 {
			y = drawer.Bounds().Dy() - sepY1 - h - sepY2*2 - sepY3*2
		}
		draw.PrepareDrawRect(drawer, w+sepX2*2+sepX3*2, h+sepY2*2+sepY3*2, util.NewNRGBColor(160, 208, 224)).SetRadius(6).Move(x, y).Draw()
		draw.PrepareDrawRect(drawer, w+sepX3*2, h+sepY3*2, util.NewNRGBColor(224, 240, 248)).Move(x+sepX2, y+sepY2).SetRadius(4).Draw()
		draw.PrepareDrawRect(drawer, w, h, util.NewNRGBColor(248, 248, 248)).Move(x+sepX2+sepX3, y+sepY2+sepY3).Draw()
		return imgutil.Rect(x, y, w+sepX2*2+sepX3*2, h+sepY2*2+sepY3*2), imgutil.Rect(x+sepX2+sepX3, y+sepY2+sepY3, w, h), nil
	case BoxStyleEnum.Battle:
		sepX1, sepY1, sepX2, sepY2, sepX3, sepY3 := 5, 5, 10, 10, 4, 4
		if w <= 0 {
			w = drawer.Bounds().Dx() - sepX1*2 - sepX2*2 - sepX3*2
		}
		if h <= 0 {
			h = 90
		}
		if x == -1 {
			x = sepX1
		}
		if y == -1 {
			y = drawer.Bounds().Dy() - sepY1 - h - sepY2*2 - sepY3*2
		}
		draw.PrepareDrawRect(drawer, w+sepX2*2+sepX3*2, h+sepY2*2+sepY3*2, util.NewNRGBColor(200, 168, 72)).SetRadius(6).Move(x, y).Draw()
		draw.PrepareDrawRect(drawer, w+sepX3*2, h+sepY3*2, util.NewNRGBColor(224, 216, 224)).Move(x+sepX2, y+sepY2).SetRadius(4).Draw()
		draw.PrepareDrawRect(drawer, w, h, util.NewNRGBColor(40, 80, 104)).Move(x+sepX2+sepX3, y+sepY2+sepY3).Draw()
		return imgutil.Rect(x, y, w+sepX2*2+sepX3*2, h+sepY2*2+sepY3*2), imgutil.Rect(x+sepX2+sepX3, y+sepY2+sepY3, w, h), nil
	default:
		panic("unreachable")
	}
}
