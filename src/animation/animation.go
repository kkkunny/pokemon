package animation

import (
	"github.com/kkkunny/pokemon/src/util/image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	frameSheet              *image.Image
	frameWidth, frameHeight int
	frameTime               int
	curFrameIndex           int

	counter int
}

func NewAnimation(frameSheet *image.Image, frameWidth, frameHeight, frameTime int) *Animation {
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

func (a *Animation) FrameTime() int {
	return a.frameTime
}

func (a *Animation) FrameCount() int {
	return a.frameSheet.Width() / a.frameWidth
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

func (a *Animation) GetFrameImage(i int) *image.Image {
	x := (i % a.FrameCount()) * a.frameWidth
	return a.frameSheet.SubImage(x, 0, a.frameWidth, a.frameHeight)
}

func (a *Animation) Draw(screen *image.Image, options ebiten.DrawImageOptions) {
	frameImg := a.GetFrameImage(a.curFrameIndex)
	screen.DrawImage(frameImg, &options)
}
