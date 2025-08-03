package util

import "image/color"

func NewRGBColor(r, g, b uint8) color.NRGBA {
	return NewRGBAColor(r, g, b, 0xff)
}

func NewRGBAColor(r, g, b, a uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// 线性插值函数
func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a)*(1-t) + float64(b)*t)
}

// GradientColor 计算两个颜色之间的渐变色
func GradientColor(from, to color.Color, percent float64) color.Color {
	r1, g1, b1, a1 := from.RGBA()
	r2, g2, b2, a2 := to.RGBA()

	// RGBA() 返回的是 uint32，需要转换为 uint8
	r := lerp(uint8(r1>>8), uint8(r2>>8), percent)
	g := lerp(uint8(g1>>8), uint8(g2>>8), percent)
	b := lerp(uint8(b1>>8), uint8(b2>>8), percent)
	a := lerp(uint8(a1>>8), uint8(a2>>8), percent)

	return color.RGBA{R: r, G: g, B: b, A: a}
}

func GetNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return color.NRGBA{}
	}

	if nrgba, ok := c.(color.NRGBA); ok {
		return nrgba
	} else if rgba, ok := c.(color.RGBA); ok {
		return color.NRGBA{
			R: uint8(uint16(rgba.R) * 0xff / uint16(rgba.A)),
			G: uint8(uint16(rgba.G) * 0xff / uint16(rgba.A)),
			B: uint8(uint16(rgba.B) * 0xff / uint16(rgba.A)),
			A: rgba.A,
		}
	} else {
		return color.NRGBA{
			R: uint8((r * 0xff / a) >> 8),
			G: uint8((g * 0xff / a) >> 8),
			B: uint8((b * 0xff / a) >> 8),
			A: uint8(a >> 8),
		}
	}
}
