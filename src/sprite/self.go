package sprite

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/input"
)

type Behavior string

var BehaviorEnum = enum.New[struct {
	Walk Behavior `enum:"walk"` // 行走
	Run  Behavior `enum:"run"`  // 奔跑
}]()

type Direction int8

var DirectionEnum = enum.New[struct {
	Up    Direction `enum:"1"`  // 上
	Down  Direction `enum:"-1"` // 下
	Left  Direction `enum:"2"`  // 左
	Right Direction `enum:"-2"` // 右
}]()

var trainerBehaviors = []Behavior{BehaviorEnum.Walk, BehaviorEnum.Run}

type Self struct {
	behaviorImages map[Behavior]*ebiten.Image
	direction      Direction // 当前所处方向
	move           bool      // 是否在移动
	x, y           int
}

func NewSelf(name string) (*Self, error) {
	dirpath := filepath.Join("./resource/map_item/people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}
	behaviorImages := make(map[Behavior]*ebiten.Image, len(trainerBehaviors))
	for _, behavior := range trainerBehaviors {
		behaviorImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(dirpath, string(behavior)+".png"))
		if err != nil {
			return nil, err
		}
		behaviorImages[behavior] = behaviorImg
	}
	return &Self{
		behaviorImages: behaviorImages,
		direction:      DirectionEnum.Down,
	}, nil
}

func (s *Self) OnAction(action input.Action) error {
	preDirection := s.direction
	switch action {
	case input.ActionEnum.MoveUp:
		s.direction = DirectionEnum.Up
	case input.ActionEnum.MoveDown:
		s.direction = DirectionEnum.Down
	case input.ActionEnum.MoveLeft:
		s.direction = DirectionEnum.Left
	case input.ActionEnum.MoveRight:
		s.direction = DirectionEnum.Right
	}
	if preDirection == s.direction {
		s.move = true
	}
	return nil
}

func (s *Self) Update() error {
	defer func() {
		s.move = false
	}()
	if s.move {
		switch s.direction {
		case DirectionEnum.Up:
			s.y--
		case DirectionEnum.Down:
			s.y++
		case DirectionEnum.Left:
			s.x--
		case DirectionEnum.Right:
			s.x++
		}
	}
	return nil
}

func (s *Self) Image() (*ebiten.Image, error) {
	img := s.behaviorImages[BehaviorEnum.Walk]
	size := img.Bounds().Size()
	frameW, frameH := size.X/3, size.Y/4

	var frameLine int
	switch s.direction {
	case DirectionEnum.Up:
		frameLine = 1
	case DirectionEnum.Down:
		frameLine = 0
	case DirectionEnum.Left:
		frameLine = 2
	case DirectionEnum.Right:
		frameLine = 3
	}
	beginX, beginY := 0, frameLine*frameH
	img = img.SubImage(image.Rect(beginX, beginY, beginX+frameW, beginY+frameH)).(*ebiten.Image)
	return img, nil
}

func (s *Self) Position() (x, y int, display bool) {
	return s.x, s.y, true
}
