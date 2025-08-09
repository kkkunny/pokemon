package option

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type DrawImageOptions struct {
	Do             func(opts DrawImageOptions)
	Image          image.Image
	X, Y           int
	ScaleX, ScaleY float64
}

func NewDrawImageOptions(img image.Image, do func(opts DrawImageOptions)) DrawImageOptions {
	return DrawImageOptions{
		Do:     do,
		Image:  img,
		ScaleX: 1.0,
		ScaleY: 1.0,
	}
}
func (opts DrawImageOptions) Draw() { opts.Do(opts) }
func (opts DrawImageOptions) SetImage(img image.Image) DrawImageOptions {
	opts.Image = img
	return opts
}
func (opts DrawImageOptions) Move(x, y int) DrawImageOptions { opts.X += x; opts.Y += y; return opts }
func (opts DrawImageOptions) Scale(x, y float64) DrawImageOptions {
	opts.ScaleX *= x
	opts.ScaleY *= y
	return opts
}

type DrawTextOptions struct {
	Do             func(opts DrawTextOptions)
	Text           string
	Font           text.Face
	Color          color.Color
	X, Y           int
	ScaleX, ScaleY float64
}

func NewDrawTextOptions(s string, font text.Face, c color.Color, do func(opts DrawTextOptions)) DrawTextOptions {
	return DrawTextOptions{
		Do:     do,
		Text:   s,
		Color:  c,
		Font:   font,
		ScaleX: 1.0,
		ScaleY: 1.0,
	}
}
func (opts DrawTextOptions) Draw()                                  { opts.Do(opts) }
func (opts DrawTextOptions) SetText(text string) DrawTextOptions    { opts.Text = text; return opts }
func (opts DrawTextOptions) SetFont(font text.Face) DrawTextOptions { opts.Font = font; return opts }
func (opts DrawTextOptions) SetColor(c color.Color) DrawTextOptions { opts.Color = c; return opts }
func (opts DrawTextOptions) Move(x, y int) DrawTextOptions          { opts.X += x; opts.Y += y; return opts }
func (opts DrawTextOptions) Scale(x, y float64) DrawTextOptions {
	opts.ScaleX *= x
	opts.ScaleY *= y
	return opts
}

type DrawRectOptions struct {
	Do            func(opts DrawRectOptions)
	Width, Height int
	Color         color.Color
	X, Y          int
	BorderWidth   int
	BorderColor   color.Color
	Radius        int
}

func NewDrawRectOptions(w, h int, c color.Color, do func(opts DrawRectOptions)) DrawRectOptions {
	return DrawRectOptions{
		Do:     do,
		Width:  w,
		Height: h,
		Color:  c,
	}
}
func (opts DrawRectOptions) Draw() { opts.Do(opts) }
func (opts DrawRectOptions) SetSize(w, h int) DrawRectOptions {
	opts.Width = w
	opts.Height = h
	return opts
}
func (opts DrawRectOptions) Move(x, y int) DrawRectOptions { opts.X += x; opts.Y += y; return opts }
func (opts DrawRectOptions) SetBorderWidth(w int) DrawRectOptions {
	opts.BorderWidth = w
	return opts
}
func (opts DrawRectOptions) SetBorderColor(c color.Color) DrawRectOptions {
	opts.BorderColor = c
	return opts
}
func (opts DrawRectOptions) SetRadius(r int) DrawRectOptions {
	opts.Radius = r
	return opts
}
