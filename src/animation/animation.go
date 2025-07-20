package animation

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	frameSheet              *ebiten.Image
	frameWidth, frameHeight int
	frameTime               int
	curFrameIndex           int

	counter int
}

func NewAnimation(frameSheet *ebiten.Image, frameWidth, frameHeight, frameTime int) *Animation {
	return &Animation{
		frameSheet:    frameSheet,
		frameWidth:    frameWidth,
		frameHeight:   frameHeight,
		frameTime:     frameTime,
		curFrameIndex: 0,
		counter:       0,
	}
}

func (a *Animation) SetFrameTime(frameTime int) {
	a.frameTime = frameTime
}

func (a *Animation) FrameCount() int {
	return a.frameSheet.Bounds().Dx() / a.frameWidth
}

func (a *Animation) Reset() {
	a.curFrameIndex = 0
	a.counter = 0
}

// Update @return: 此轮动画是否结束
func (a *Animation) Update() bool {
	a.counter++
	if a.counter >= a.frameTime {
		a.counter = 0
		a.curFrameIndex = (a.curFrameIndex + 1) % a.FrameCount()
	}
	return a.counter == 0 && a.curFrameIndex == 0
}

func (a *Animation) Draw(screen *ebiten.Image, x, y float64) {
	sx := (a.curFrameIndex % (a.frameSheet.Bounds().Dx() / a.frameWidth)) * a.frameWidth
	sy := (a.curFrameIndex / (a.frameSheet.Bounds().Dx() / a.frameWidth)) * a.frameHeight

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	screen.DrawImage(
		a.frameSheet.SubImage(image.Rect(sx, sy, sx+a.frameWidth, sy+a.frameHeight)).(*ebiten.Image),
		op,
	)
}
