package consts

import (
	"image/color"
)

const (
	MaleText   = "♂"
	FemaleText = "♀"
)

var (
	MaleColor   = color.NRGBA{R: 51, G: 85, B: 255, A: 0xff}   // 雄性颜色
	FemaleColor = color.NRGBA{R: 255, G: 119, B: 221, A: 0xff} // 雌性颜色
)
