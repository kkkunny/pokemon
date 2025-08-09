package animation

import (
	"image"
	"image/gif"
	"time"

	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/util/image"
)

type Animation struct {
	frameImages   []imgutil.Image
	frameTime     int
	curFrameIndex int

	counter int
}

func NewAnimation(frameImages []imgutil.Image, frameTime int) *Animation {
	return &Animation{
		frameImages:   frameImages,
		frameTime:     frameTime,
		curFrameIndex: 0,
		counter:       0,
	}
}

func NewAnimationFromGIF(g *gif.GIF) *Animation {
	var frameTime int
	if len(g.Delay) != 0 {
		frameTime = int((time.Millisecond * 10 * time.Duration(g.Delay[0])) / (time.Second / 60))
	}
	return &Animation{
		frameImages:   stlslices.Map(g.Image, func(_ int, img *image.Paletted) imgutil.Image { return imgutil.WrapImage(img) }),
		frameTime:     frameTime,
		curFrameIndex: 0,
		counter:       0,
	}
}

func (a *Animation) Frames() []imgutil.Image {
	return a.frameImages
}

func (a *Animation) AddFrame(frame imgutil.Image) {
	a.frameImages = append(a.frameImages, frame)
}

func (a *Animation) SetFrameTime(frameTime int) {
	a.frameTime = frameTime
}

func (a *Animation) FrameTime() int {
	return a.frameTime
}

func (a *Animation) FrameCount() int {
	return len(a.frameImages)
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

func (a *Animation) GetFrameImage(i int) imgutil.Image {
	return a.frameImages[i]
}

func (a *Animation) GetCurrentFrameImage() imgutil.Image {
	return a.frameImages[a.curFrameIndex]
}
