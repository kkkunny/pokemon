package draw

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/kkkunny/pokemon/src/util/draw/option"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

func getDrawOptions(drawer draw.Image) _options {
	if drawerWithOpts, ok := drawer.(OptionDrawer); ok {
		return drawerWithOpts.Options()
	}
	return newOptions()
}

var drawImageSpecialFuncList []func(drawer draw.Image, opts option.DrawImageOptions) bool

// PrepareDrawImage 绘制图像
func PrepareDrawImage(drawer draw.Image, img image.Image) option.DrawImageOptions {
	return option.NewDrawImageOptions(img, func(opts option.DrawImageOptions) {
		for _, fn := range drawImageSpecialFuncList {
			if fn(drawer, opts) {
				return
			}
		}
		// TODO: 其他类型的实现
		panic("todo")
	})
}

var drawTextSpecialFuncList []func(drawer draw.Image, opts option.DrawTextOptions) bool

// PrepareDrawText 绘制文字
func PrepareDrawText(drawer draw.Image, text string, font text.Face, c color.Color) option.DrawTextOptions {
	return option.NewDrawTextOptions(text, font, c, func(opts option.DrawTextOptions) {
		for _, fn := range drawTextSpecialFuncList {
			if fn(drawer, opts) {
				return
			}
		}
		// TODO: 其他类型的实现
		panic("todo")
	})
}

var drawRectSpecialFuncList []func(drawer draw.Image, opts option.DrawRectOptions) bool

// PrepareDrawRect 绘制矩形
func PrepareDrawRect(drawer draw.Image, w, h int, c color.Color) option.DrawRectOptions {
	return option.NewDrawRectOptions(w, h, c, func(opts option.DrawRectOptions) {
		for _, fn := range drawRectSpecialFuncList {
			if fn(drawer, opts) {
				return
			}
		}
		// TODO: 其他类型的实现
		panic("todo")
	})
}

func OverlayColor(drawer draw.Image, c color.Color) {
	mask := imgutil.NewImage(drawer.Bounds().Dx(), drawer.Bounds().Dy())
	mask.Fill(c)
	PrepareDrawImage(drawer, mask).Draw()
}
