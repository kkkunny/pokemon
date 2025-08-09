package draw

import (
	"image"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/kkkunny/pokemon/src/util/draw/option"
)

type ebitenImageGetter interface {
	EbitenImage() *ebiten.Image
}

func getEbitenImage(img image.Image) (*ebiten.Image, bool) {
	if drawer, ok := img.(_optionDrawer); ok {
		return getEbitenImage(drawer.Image)
	}

	getter, ok := img.(ebitenImageGetter)
	if ok {
		return getter.EbitenImage(), true
	}
	ebitenImg, ok := img.(*ebiten.Image)
	return ebitenImg, ok
}

func drawEbitenImage(drawer draw.Image, opts option.DrawImageOptions) bool {
	bgImg, ok := getEbitenImage(drawer)
	if !ok {
		return false
	}

	img, ok := getEbitenImage(opts.Image)
	if !ok {
		img = ebiten.NewImageFromImage(opts.Image)
	}

	globalOpts := getDrawOptions(drawer)
	opts.ScaleX *= globalOpts.scaleX
	opts.ScaleY *= globalOpts.scaleY
	opts.X = int(float64(opts.X)*globalOpts.scaleX + float64(globalOpts.x))
	opts.Y = int(float64(opts.Y)*globalOpts.scaleY + float64(globalOpts.y))

	var imgOps ebiten.DrawImageOptions
	imgOps.GeoM.Scale(opts.ScaleX, opts.ScaleY)
	imgOps.GeoM.Translate(float64(opts.X), float64(opts.Y))
	bgImg.DrawImage(img, &imgOps)
	return true
}

func drawEbitenText(drawer draw.Image, opts option.DrawTextOptions) bool {
	bgImg, ok := getEbitenImage(drawer)
	if !ok {
		return false
	}

	globalOpts := getDrawOptions(drawer)
	opts.ScaleX *= globalOpts.scaleX
	opts.ScaleY *= globalOpts.scaleY
	opts.X = int(float64(opts.X)*globalOpts.scaleX + float64(globalOpts.x))
	opts.Y = int(float64(opts.Y)*globalOpts.scaleY + float64(globalOpts.y))

	var textOps text.DrawOptions
	textOps.GeoM.Scale(opts.ScaleX, opts.ScaleY)
	textOps.GeoM.Translate(float64(opts.X), float64(opts.Y))
	textOps.ColorScale.ScaleWithColor(opts.Color)
	text.Draw(bgImg, opts.Text, opts.Font, &textOps)
	return true
}

func drawEbitenRect(drawer draw.Image, opts option.DrawRectOptions) bool {
	bgImg, ok := getEbitenImage(drawer)
	if !ok {
		return false
	}

	globalOpts := getDrawOptions(drawer)
	opts.Width = int(float64(opts.Width) * globalOpts.scaleX)
	opts.Height = int(float64(opts.Height) * globalOpts.scaleY)
	opts.X = int(float64(opts.X)*globalOpts.scaleX + float64(globalOpts.x))
	opts.Y = int(float64(opts.Y)*globalOpts.scaleY + float64(globalOpts.y))
	opts.BorderWidth = int(float64(opts.BorderWidth) * globalOpts.scaleX)
	opts.Radius = int(float64(opts.Radius) * globalOpts.scaleX)

	radius := min(min(opts.Width/2, opts.Height/2), opts.Radius)
	if radius == 0 {
		if opts.Color != nil {
			vector.DrawFilledRect(bgImg, float32(opts.X), float32(opts.Y), float32(opts.Width), float32(opts.Height), opts.Color, false)
		}
		if opts.BorderWidth > 0 && opts.BorderColor != nil {
			vector.StrokeRect(bgImg, float32(opts.X), float32(opts.Y), float32(opts.Width), float32(opts.Height), float32(opts.BorderWidth), opts.BorderColor, false)
		}
	} else if float64(radius) == min(float64(opts.Width)/2, float64(opts.Height)/2) {
		if opts.Color != nil {
			vector.DrawFilledCircle(bgImg, float32(opts.X)+float32(opts.Width)/2, float32(opts.Y)+float32(opts.Height)/2, float32(radius), opts.Color, false)
		}
		if opts.BorderWidth > 0 && opts.BorderColor != nil {
			vector.StrokeCircle(bgImg, float32(opts.X)+float32(opts.Width)/2, float32(opts.Y)+float32(opts.Height)/2, float32(radius), float32(opts.BorderWidth), opts.BorderColor, false)
		}
	} else {
		if opts.Color != nil {
			// 左上
			vector.DrawFilledCircle(bgImg, float32(opts.X+radius), float32(opts.Y+radius), float32(opts.Radius), opts.Color, false)
			// 右上
			vector.DrawFilledCircle(bgImg, float32(opts.X+opts.Width-radius), float32(opts.Y+radius), float32(opts.Radius), opts.Color, false)
			// 左下
			vector.DrawFilledCircle(bgImg, float32(opts.X+radius), float32(opts.Y+opts.Height-radius), float32(opts.Radius), opts.Color, false)
			// 右下
			vector.DrawFilledCircle(bgImg, float32(opts.X+opts.Width-radius), float32(opts.Y+opts.Height-radius), float32(opts.Radius), opts.Color, false)
		}
		if opts.BorderWidth > 0 && opts.BorderColor != nil {
			// 左上
			vector.StrokeCircle(bgImg, float32(opts.X+radius), float32(opts.Y+radius), float32(radius), float32(opts.BorderWidth)*2, opts.BorderColor, false)
			// 右上
			vector.StrokeCircle(bgImg, float32(opts.X+opts.Width-radius), float32(opts.Y+radius), float32(radius), float32(opts.BorderWidth)*2, opts.BorderColor, false)
			// 左下
			vector.StrokeCircle(bgImg, float32(opts.X+radius), float32(opts.Y+opts.Height-radius), float32(radius), float32(opts.BorderWidth)*2, opts.BorderColor, false)
			// 右下
			vector.StrokeCircle(bgImg, float32(opts.X+opts.Width-radius), float32(opts.Y+opts.Height-radius), float32(radius), float32(opts.BorderWidth)*2, opts.BorderColor, false)
		}
		if opts.Color != nil {
			// 中
			vector.DrawFilledRect(bgImg, float32(opts.X+radius), float32(opts.Y+radius), float32(opts.Width-radius*2), float32(opts.Height-radius*2), opts.Color, false)
			// 上
			vector.DrawFilledRect(bgImg, float32(opts.X+radius), float32(opts.Y), float32(opts.Width-radius*2), float32(radius), opts.Color, false)
			// 下
			vector.DrawFilledRect(bgImg, float32(opts.X+radius), float32(opts.Y+opts.Height-radius), float32(opts.Width-radius*2), float32(radius), opts.Color, false)
			// 左
			vector.DrawFilledRect(bgImg, float32(opts.X), float32(opts.Y+radius), float32(radius), float32(opts.Height-radius*2), opts.Color, false)
			// 右
			vector.DrawFilledRect(bgImg, float32(opts.X+opts.Width-radius), float32(opts.Y+radius), float32(radius), float32(opts.Height-radius*2), opts.Color, false)
		}
		if opts.BorderWidth > 0 && opts.BorderColor != nil {
			// 上
			vector.DrawFilledRect(bgImg, float32(opts.X+radius), float32(opts.Y), float32(opts.Width-radius*2), float32(opts.BorderWidth), opts.BorderColor, false)
			// 下
			vector.DrawFilledRect(bgImg, float32(opts.X+radius), float32(opts.Y+opts.Height-opts.BorderWidth), float32(opts.Width-radius*2), float32(opts.BorderWidth), opts.BorderColor, false)
			// 左
			vector.DrawFilledRect(bgImg, float32(opts.X), float32(opts.Y+radius), float32(opts.BorderWidth), float32(opts.Height-radius*2), opts.BorderColor, false)
			// 右
			vector.DrawFilledRect(bgImg, float32(opts.X+opts.Width-opts.BorderWidth), float32(opts.Y+radius), float32(opts.BorderWidth), float32(opts.Height-radius*2), opts.BorderColor, false)
		}
	}
	return true
}

func init() {
	drawImageSpecialFuncList = append(drawImageSpecialFuncList, drawEbitenImage)
	drawTextSpecialFuncList = append(drawTextSpecialFuncList, drawEbitenText)
	drawRectSpecialFuncList = append(drawRectSpecialFuncList, drawEbitenRect)
}
