package animation

import (
	"fmt"
	"image"
	"image/gif"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kettek/apng"
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

func NewAnimationFromFile(path string) (*Animation, error) {
	filename := filepath.Base(path)
	fileExt := strings.TrimPrefix(filepath.Ext(filename), ".")

	allowExts := []string{"png", "gif"}
	if !stlslices.Contain(allowExts, fileExt) {
		return nil, fmt.Errorf("%s is not a valid animation", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	switch fileExt {
	case "gif":
		g, err := gif.DecodeAll(file)
		if err != nil {
			return nil, err
		}

		var frameTime int
		if len(g.Delay) != 0 {
			frameTime = int((time.Second / 100 * time.Duration(g.Delay[0])) / (time.Second / 60))
		}
		return &Animation{
			frameImages:   stlslices.Map(g.Image, func(_ int, img *image.Paletted) imgutil.Image { return imgutil.WrapImage(img) }),
			frameTime:     frameTime,
			curFrameIndex: 0,
			counter:       0,
		}, nil
	case "png":
		a, err := apng.DecodeAll(file)
		if err != nil {
			return nil, err
		}

		var frameTime int
		if len(a.Frames) != 0 {
			frameTime = int(time.Duration(a.Frames[0].GetDelay()*float64(time.Second)) / (time.Second / 60))
		}
		return &Animation{
			frameImages:   stlslices.Map(a.Frames, func(_ int, frame apng.Frame) imgutil.Image { return imgutil.WrapImage(frame.Image) }),
			frameTime:     frameTime,
			curFrameIndex: 0,
			counter:       0,
		}, nil
	default:
		return nil, fmt.Errorf("%s is not a valid animation", path)
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
