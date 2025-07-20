package sprite

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/animation"
	"github.com/kkkunny/pokemon/src/input"
)

type Behavior string

var BehaviorEnum = enum.New[struct {
	Walk Behavior `enum:"walk"`
	Run  Behavior `enum:"run"`
}]()

type Direction int8

var DirectionEnum = enum.New[struct {
	Up    Direction `enum:"-1"`
	Down  Direction `enum:"1"`
	Left  Direction `enum:"-3"`
	Right Direction `enum:"3"`
}]()

var trainerBehaviors = []Behavior{BehaviorEnum.Walk, BehaviorEnum.Run}

type Self struct {
	// 静态资源
	directionImages    map[Direction]*ebiten.Image                     // 方向图片
	behaviorAnimations map[Behavior]map[Direction]*animation.Animation // 行为动画
	// 属性
	direction Direction // 当前所处方向
	pos       [2]int    // 当前位置
	expectPos [2]int    // 预期所处的位置，用于移动
}

func NewSelf(name string) (*Self, error) {
	dirpath := filepath.Join("./resource/map_item/people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}

	directions := enum.Values[Direction](DirectionEnum)
	directionImages := make(map[Direction]*ebiten.Image, len(directions))
	behaviorAnimations := make(map[Behavior]map[Direction]*animation.Animation, len(trainerBehaviors))
	for _, behavior := range trainerBehaviors {
		behaviorImgSheetRect, _, err := ebitenutil.NewImageFromFile(filepath.Join(dirpath, string(behavior)+".png"))
		if err != nil {
			return nil, err
		}
		frameW, frameH := behaviorImgSheetRect.Bounds().Dx()/3, behaviorImgSheetRect.Bounds().Dy()/4
		behaviorDirectionAnimations := make(map[Direction]*animation.Animation, len(directions))
		for i, direction := range []Direction{DirectionEnum.Down, DirectionEnum.Up, DirectionEnum.Left, DirectionEnum.Right} {
			y := i * frameH
			animationFrameSheet := ebiten.NewImage(4*frameW, frameH)
			for j := range 3 {
				x := j * frameW
				img := behaviorImgSheetRect.SubImage(image.Rect(x, y, x+frameW, y+frameH)).(*ebiten.Image)
				switch j {
				case 0:
					if behavior == BehaviorEnum.Walk {
						directionImages[direction] = img
					}
					ops := &ebiten.DrawImageOptions{}
					ops.GeoM.Translate(float64(frameW), 0)
					animationFrameSheet.DrawImage(img, ops)
					ops.GeoM.Translate(2*float64(frameW), 0)
					animationFrameSheet.DrawImage(img, ops)
				case 1:
					ops := &ebiten.DrawImageOptions{}
					ops.GeoM.Translate(0, 0)
					animationFrameSheet.DrawImage(img, ops)
				case 2:
					ops := &ebiten.DrawImageOptions{}
					ops.GeoM.Translate(2*float64(frameW), 0)
					animationFrameSheet.DrawImage(img, ops)
				}
			}
			behaviorDirectionAnimations[direction] = animation.NewAnimation(animationFrameSheet, frameW, frameH, 10)
		}
		behaviorAnimations[behavior] = behaviorDirectionAnimations
	}
	return &Self{
		directionImages:    directionImages,
		behaviorAnimations: behaviorAnimations,
		direction:          DirectionEnum.Down,
	}, nil
}

func (s *Self) Move() bool {
	return s.pos != s.expectPos
}

func (s *Self) OnAction(action input.Action) {
	if s.Move() {
		return
	}
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
		switch action {
		case input.ActionEnum.MoveUp:
			s.expectPos = [2]int{s.pos[0], s.pos[1] + int(DirectionEnum.Up)%2}
		case input.ActionEnum.MoveDown:
			s.expectPos = [2]int{s.pos[0], s.pos[1] + int(DirectionEnum.Down)%2}
		case input.ActionEnum.MoveLeft:
			s.expectPos = [2]int{s.pos[0] + int(DirectionEnum.Left)%2, s.pos[1]}
		case input.ActionEnum.MoveRight:
			s.expectPos = [2]int{s.pos[0] + int(DirectionEnum.Right)%2, s.pos[1]}
		}
	}
	return
}

func (s *Self) Update() error {
	if !s.Move() {
		return nil
	}
	// switch s.direction {
	// case DirectionEnum.Up:
	// 	s.y--
	// case DirectionEnum.Down:
	// 	s.y++
	// case DirectionEnum.Left:
	// 	s.x--
	// case DirectionEnum.Right:
	// 	s.x++
	// }
	a := s.behaviorAnimations[BehaviorEnum.Walk][s.direction]
	if a.Update() {
		s.pos = s.expectPos
	}
	return nil
}

func (s *Self) Draw(screen *ebiten.Image, info *DrawInfo) {
	if info.Person == nil {
		return
	}

	tileW, tileH := info.Person.Map.TileSize()
	x, y := tileW*s.pos[0], tileH*s.pos[1]
	img := s.directionImages[s.direction]
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, ops)
}
