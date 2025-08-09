package animation

import (
	"image"

	"github.com/kkkunny/pokemon/src/util/image"
)

type Animation struct {
	frameSheet              *imgutil.Image
	frameWidth, frameHeight int
	frameTime               int
	curFrameIndex           int

	counter int
}

func NewAnimation(frameSheet *imgutil.Image, frameWidth, frameHeight, frameTime int) *Animation {
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

func (a *Animation) GetFrameImage(i int) *imgutil.Image {
	x := (i % a.FrameCount()) * a.frameWidth
	return a.frameSheet.SubImage(image.Rect(x, 0, x+a.frameWidth, a.frameHeight))
}

func (a *Animation) GetCurrentFrameImage() *imgutil.Image {
	return a.GetFrameImage(a.curFrameIndex)
}
