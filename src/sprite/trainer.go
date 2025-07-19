package sprite

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tnnmigga/enum"
)

type Behavior string

var BehaviorEnum = enum.New[struct {
	Walk Behavior `enum:"walk"` // 行走
	Run  Behavior `enum:"run"`  // 奔跑
}]()

type Direction int8

var DirectionEnum = enum.New[struct {
	Down  Direction `enum:"1"`  // 下
	Up    Direction `enum:"-1"` // 上
	Left  Direction `enum:"2"`  // 左
	Right Direction `enum:"-2"` // 右
}]()

var trainerBehaviors = []Behavior{BehaviorEnum.Walk, BehaviorEnum.Run}

type Trainer struct {
	behaviorImages map[Behavior]image.Image
	direction      Direction
}

func NewTrainer(name string) (*Trainer, error) {
	dirpath := filepath.Join("./resource/map_item/people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}
	behaviorImages := make(map[Behavior]image.Image, len(trainerBehaviors))
	for _, behavior := range trainerBehaviors {
		_, behaviorImg, err := ebitenutil.NewImageFromFile(filepath.Join(dirpath, string(behavior)+".png"))
		if err != nil {
			return nil, err
		}
		behaviorImages[behavior] = behaviorImg
	}
	return &Trainer{
		behaviorImages: behaviorImages,
		direction:      DirectionEnum.Down,
	}, nil
}

func (t *Trainer) Update() error {
	return nil
}

func (t *Trainer) Image() (image.Image, error) {
	img := t.behaviorImages[BehaviorEnum.Walk]
	size := img.Bounds().Size()
	frameW, frameH := size.X/3, size.Y/4

	var frameLine int
	switch t.direction {
	case DirectionEnum.Down:
		frameLine = 0
	case DirectionEnum.Up:
		frameLine = 1
	case DirectionEnum.Left:
		frameLine = 2
	case DirectionEnum.Right:
		frameLine = 3
	}
	beginX, beginY := 0, frameLine*frameH

	img = imaging.Crop(img, image.Rect(beginX, beginY, beginX+frameW, beginY+frameH))
	return img, nil
}
