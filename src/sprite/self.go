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
	"github.com/kkkunny/pokemon/src/config"
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

type Foot int8

var FootEnum = enum.New[struct {
	Left  Foot `enum:"1"`
	Right Foot `enum:"-1"`
}]()

var trainerBehaviors = []Behavior{BehaviorEnum.Walk, BehaviorEnum.Run}

type Self struct {
	// 静态资源
	directionImages    map[Direction]*ebiten.Image                              // 方向图片
	behaviorAnimations map[Behavior]map[Direction]map[Foot]*animation.Animation // 行为动画
	// 属性
	direction        Direction // 当前所处方向
	moveStartingFoot Foot      // 移动时的起始脚
	speed            int       // 移动速度
	pos              [2]int    // 当前地块位置
	expectPos        [2]int    // 预期所处的地块位置，用于移动

	moveCounter int // 移动时的计数器，用于显示动画
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
	behaviorAnimations := make(map[Behavior]map[Direction]map[Foot]*animation.Animation, len(trainerBehaviors))
	for _, behavior := range trainerBehaviors {
		behaviorImgSheetRect, _, err := ebitenutil.NewImageFromFile(filepath.Join(dirpath, string(behavior)+".png"))
		if err != nil {
			return nil, err
		}
		frameW, frameH := behaviorImgSheetRect.Bounds().Dx()/3, behaviorImgSheetRect.Bounds().Dy()/4
		behaviorDirectionAnimations := make(map[Direction]map[Foot]*animation.Animation, len(directions))
		for i, direction := range []Direction{DirectionEnum.Down, DirectionEnum.Up, DirectionEnum.Left, DirectionEnum.Right} {
			y := i * frameH
			leftFootAnimationFrameSheet := ebiten.NewImage(2*frameW, frameH)
			rightFootAnimationFrameSheet := ebiten.NewImage(2*frameW, frameH)
			for j := range 3 {
				x := j * frameW
				img := behaviorImgSheetRect.SubImage(image.Rect(x, y, x+frameW, y+frameH)).(*ebiten.Image)
				switch j {
				case 0:
					if behavior == BehaviorEnum.Walk {
						directionImages[direction] = img
					}
					ops := &ebiten.DrawImageOptions{}
					ops.GeoM.Translate(0, 0)
					leftFootAnimationFrameSheet.DrawImage(img, ops)
					ops.GeoM.Translate(0, 0)
					rightFootAnimationFrameSheet.DrawImage(img, ops)
				case 1:
					ops := &ebiten.DrawImageOptions{}
					ops.GeoM.Translate(float64(frameW), 0)
					leftFootAnimationFrameSheet.DrawImage(img, ops)
				case 2:
					ops := &ebiten.DrawImageOptions{}
					ops.GeoM.Translate(float64(frameW), 0)
					rightFootAnimationFrameSheet.DrawImage(img, ops)
				}
			}
			behaviorDirectionAnimations[direction] = map[Foot]*animation.Animation{
				FootEnum.Left:  animation.NewAnimation(leftFootAnimationFrameSheet, frameW, frameH, 0),
				FootEnum.Right: animation.NewAnimation(rightFootAnimationFrameSheet, frameW, frameH, 0),
			}
		}
		behaviorAnimations[behavior] = behaviorDirectionAnimations
	}
	return &Self{
		directionImages:    directionImages,
		behaviorAnimations: behaviorAnimations,
		direction:          DirectionEnum.Down,
		moveStartingFoot:   FootEnum.Left,
		speed:              1,
		pos:                [2]int{6, 8},
	}, nil
}

func (s *Self) Move() bool {
	return s.pos != s.expectPos
}

func (s *Self) SetPosition(x, y int) {
	s.pos = [2]int{x, y}
	s.expectPos = [2]int{x, y}
}

func (s *Self) OnAction(action input.Action, info *UpdateInfo) {
	if info.Person == nil {
		return
	}

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
		expectPos := s.pos
		switch action {
		case input.ActionEnum.MoveUp:
			expectPos = [2]int{s.pos[0], s.pos[1] + int(DirectionEnum.Up)%2}
		case input.ActionEnum.MoveDown:
			expectPos = [2]int{s.pos[0], s.pos[1] + int(DirectionEnum.Down)%2}
		case input.ActionEnum.MoveLeft:
			expectPos = [2]int{s.pos[0] + int(DirectionEnum.Left)%2, s.pos[1]}
		case input.ActionEnum.MoveRight:
			expectPos = [2]int{s.pos[0] + int(DirectionEnum.Right)%2, s.pos[1]}
		}
		if !info.Person.Map.CheckCollision(expectPos[0], expectPos[1]) {
			s.expectPos = expectPos
		}
	}
	return
}

func (s *Self) Update() error {
	if !s.Move() {
		return nil
	}

	a := s.behaviorAnimations[BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
	a.SetFrameTime(config.TileSize / s.speed / a.FrameCount())
	a.Update()

	diff := config.TileSize - s.moveCounter
	if diff > s.speed {
		s.moveCounter += s.speed
	} else {
		s.moveCounter = 0
		s.pos = s.expectPos
		s.moveStartingFoot = -s.moveStartingFoot
		a.Reset()
	}
	return nil
}

func (s *Self) Draw(screen *ebiten.Image) error {
	img := s.directionImages[s.direction]

	x, y := s.pos[0]*config.TileSize, s.pos[1]*config.TileSize+config.TileSize-img.Bounds().Dy()
	switch s.direction {
	case DirectionEnum.Up:
		y -= s.moveCounter
	case DirectionEnum.Down:
		y += s.moveCounter
	case DirectionEnum.Left:
		x -= s.moveCounter
	case DirectionEnum.Right:
		x += s.moveCounter
	}

	if s.Move() {
		a := s.behaviorAnimations[BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
		a.Draw(screen, float64(x), float64(y))
	} else {
		ops := &ebiten.DrawImageOptions{}
		ops.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, ops)
	}
	return nil
}
