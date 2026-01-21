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
	frames []Frame
}

func NewAnimation(frame ...Frame) *Animation {
	return &Animation{
		frames: frame,
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

		if len(g.Delay) != len(g.Image) {
			return nil, fmt.Errorf("gif frame count is not equal to delay count")
		}
		frames := stlslices.Map(g.Image, func(i int, img *image.Paletted) Frame {
			return Frame{
				Image: imgutil.WrapImage(img),
				Time:  time.Second / 100 * time.Duration(g.Delay[i]),
			}
		})

		return NewAnimation(frames...), nil
	case "png":
		a, err := apng.DecodeAll(file)
		if err != nil {
			return nil, err
		}

		frames := stlslices.Map(a.Frames, func(i int, frame apng.Frame) Frame {
			return Frame{
				Image: imgutil.WrapImage(frame.Image),
				Time:  time.Duration(frame.GetDelay() * float64(time.Second)),
			}
		})

		return NewAnimation(frames...), nil
	default:
		return nil, fmt.Errorf("%s is not a valid animation", path)
	}
}

func (a *Animation) Frames() []Frame {
	return a.frames
}

func (a *Animation) AddFrame(frame Frame) Frame {
	a.frames = append(a.frames, frame)
	return frame
}

func (a *Animation) NewPlayer(fps int) *Player {
	return &Player{
		animation: a,
		frameCounters: stlslices.Map(a.frames, func(_ int, f Frame) int {
			return int(f.Time / (time.Second / time.Duration(fps)))
		}),
	}
}
