package imgutil

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

func (i *Image) SubImage(r image.Rectangle) *Image {
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

type DrawRectOptions struct {
	drawFn        func(opts *DrawRectOptions)
	x, y          int
	width, height int
	color         color.Color
	borderWidth   int
	borderColor   color.Color
	radius        int
}

func (opts *DrawRectOptions) Draw()                          { opts.drawFn(opts) }
func (opts *DrawRectOptions) Move(x, y int) *DrawRectOptions { opts.x += x; opts.y += y; return opts }
func (opts *DrawRectOptions) SetBorderWidth(w int) *DrawRectOptions {
	opts.borderWidth = w
	return opts
}
func (opts *DrawRectOptions) SetBorderColor(c color.Color) *DrawRectOptions {
	opts.borderColor = c
	return opts
}
func (opts *DrawRectOptions) SetRadius(r int) *DrawRectOptions {
	opts.radius = r
	return opts
}

func (i *Image) DrawRect(w, h int, c color.Color) *DrawRectOptions {
	return &DrawRectOptions{
		drawFn: func(opts *DrawRectOptions) {
			radius := min(min(opts.width/2, opts.height/2), opts.radius)
			if radius == 0 {
				if opts.color != nil {
					vector.DrawFilledRect(i.Image, float32(opts.x), float32(opts.y), float32(opts.width), float32(opts.height), opts.color, false)
				}
				if opts.borderWidth > 0 && opts.borderColor != nil {
					vector.StrokeRect(i.Image, float32(opts.x), float32(opts.y), float32(opts.width), float32(opts.height), float32(opts.borderWidth), opts.borderColor, false)
				}
			} else if float64(radius) == min(float64(opts.width)/2, float64(opts.height)/2) {
				if opts.color != nil {
					vector.DrawFilledCircle(i.Image, float32(opts.x)+float32(opts.width)/2, float32(opts.y)+float32(opts.height)/2, float32(radius), opts.color, false)
				}
				if opts.borderWidth > 0 && opts.borderColor != nil {
					vector.StrokeCircle(i.Image, float32(opts.x)+float32(opts.width)/2, float32(opts.y)+float32(opts.height)/2, float32(radius), float32(opts.borderWidth), opts.borderColor, false)
				}
			} else {
				if opts.color != nil {
					// 左上
					vector.DrawFilledCircle(i.Image, float32(opts.x+radius), float32(opts.y+radius), float32(opts.radius), opts.color, false)
					// 右上
					vector.DrawFilledCircle(i.Image, float32(opts.x+opts.width-radius), float32(opts.y+radius), float32(opts.radius), opts.color, false)
					// 左下
					vector.DrawFilledCircle(i.Image, float32(opts.x+radius), float32(opts.y+opts.height-radius), float32(opts.radius), opts.color, false)
					// 右下
					vector.DrawFilledCircle(i.Image, float32(opts.x+opts.width-radius), float32(opts.y+opts.height-radius), float32(opts.radius), opts.color, false)
				}
				if opts.borderWidth > 0 && opts.borderColor != nil {
					// 左上
					vector.StrokeCircle(i.Image, float32(opts.x+radius), float32(opts.y+radius), float32(radius), float32(opts.borderWidth)*2, opts.borderColor, false)
					// 右上
					vector.StrokeCircle(i.Image, float32(opts.x+opts.width-radius), float32(opts.y+radius), float32(radius), float32(opts.borderWidth)*2, opts.borderColor, false)
					// 左下
					vector.StrokeCircle(i.Image, float32(opts.x+radius), float32(opts.y+opts.height-radius), float32(radius), float32(opts.borderWidth)*2, opts.borderColor, false)
					// 右下
					vector.StrokeCircle(i.Image, float32(opts.x+opts.width-radius), float32(opts.y+opts.height-radius), float32(radius), float32(opts.borderWidth)*2, opts.borderColor, false)
				}
				if opts.color != nil {
					// 中
					vector.DrawFilledRect(i.Image, float32(opts.x+radius), float32(opts.y+radius), float32(opts.width-radius*2), float32(opts.height-radius*2), opts.color, false)
					// 上
					vector.DrawFilledRect(i.Image, float32(opts.x+radius), float32(opts.y), float32(opts.width-radius*2), float32(radius), opts.color, false)
					// 下
					vector.DrawFilledRect(i.Image, float32(opts.x+radius), float32(opts.y+opts.height-radius), float32(opts.width-radius*2), float32(radius), opts.color, false)
					// 左
					vector.DrawFilledRect(i.Image, float32(opts.x), float32(opts.y+radius), float32(radius), float32(opts.height-radius*2), opts.color, false)
					// 右
					vector.DrawFilledRect(i.Image, float32(opts.x+opts.width-radius), float32(opts.y+radius), float32(radius), float32(opts.height-radius*2), opts.color, false)
				}
				if opts.borderWidth > 0 && opts.borderColor != nil {
					// 上
					vector.DrawFilledRect(i.Image, float32(opts.x+radius), float32(opts.y), float32(opts.width-radius*2), float32(opts.borderWidth), opts.borderColor, false)
					// 下
					vector.DrawFilledRect(i.Image, float32(opts.x+radius), float32(opts.y+opts.height-opts.borderWidth), float32(opts.width-radius*2), float32(opts.borderWidth), opts.borderColor, false)
					// 左
					vector.DrawFilledRect(i.Image, float32(opts.x), float32(opts.y+radius), float32(opts.borderWidth), float32(opts.height-radius*2), opts.borderColor, false)
					// 右
					vector.DrawFilledRect(i.Image, float32(opts.x+opts.width-opts.borderWidth), float32(opts.y+radius), float32(opts.borderWidth), float32(opts.height-radius*2), opts.borderColor, false)
				}
			}
		},
		width:  w,
		height: h,
		color:  c,
	}
}

type DrawTextOptions struct {
	drawFn func(opts *DrawTextOptions)
	text   string
	x, y   int
	color  color.Color
	font   text.Face
}

func (opts *DrawTextOptions) Draw()                          { opts.drawFn(opts) }
func (opts *DrawTextOptions) Move(x, y int) *DrawTextOptions { opts.x += x; opts.y += y; return opts }

func (i *Image) DrawText(s string, font text.Face, c color.Color) *DrawTextOptions {
	return &DrawTextOptions{
		drawFn: func(opts *DrawTextOptions) {
			var textOps text.DrawOptions
			textOps.ColorScale.ScaleWithColor(color.Black)
			textOps.GeoM.Translate(float64(opts.x), float64(opts.y))
			text.Draw(i.Image, opts.text, opts.font, &textOps)
		},
		text:  s,
		color: c,
		font:  font,
	}
}
