package sprite

import (
	"fmt"
	"image"
	"math"
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
	direction        Direction  // 当前所处方向
	moveStartingFoot Foot       // 移动时的起始脚
	speed            int        // 移动速度
	pos              [2]int     // 当前地块位置
	expectPos        [2]int     // 预期所处的地块位置，用于移动
	pixelPos         [2]float64 // 当前像素位置
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

func (s *Self) Update(info *UpdateInfo) error {
	if info.Person == nil {
		return nil
	}

	if !s.Move() {
		return nil
	}
	// 更新动画
	tileSize := info.Person.Map.TileSize()
	a := s.behaviorAnimations[BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
	a.SetFrameTime(tileSize / s.speed / a.FrameCount())
	a.Update()

	// 更新xy
	expectPixelX := float64(s.expectPos[0] * tileSize)
	expectPixelY := float64(s.expectPos[1] * tileSize)

	dx := expectPixelX - s.pixelPos[0]
	dy := expectPixelY - s.pixelPos[1]
	if math.Abs(dx) > float64(s.speed) {
		s.pixelPos[0] += float64(s.speed) * math.Copysign(1, dx)
	} else {
		s.pixelPos[0] = expectPixelX
	}
	if math.Abs(dy) > float64(s.speed) {
		s.pixelPos[1] += float64(s.speed) * math.Copysign(1, dy)
	} else {
		s.pixelPos[1] = expectPixelY
	}

	// 移动结束
	if s.pixelPos[0] == expectPixelX && s.pixelPos[1] == expectPixelY {
		s.pos = s.expectPos
		s.moveStartingFoot = -s.moveStartingFoot
		a.Reset()
	}
	return nil
}

func (s *Self) Draw(screen *ebiten.Image) {
	if s.Move() {
		a := s.behaviorAnimations[BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
		a.Draw(screen, s.pixelPos[0], s.pixelPos[1])
	} else {
		ops := &ebiten.DrawImageOptions{}
		ops.GeoM.Translate(s.pixelPos[0], s.pixelPos[1])
		img := s.directionImages[s.direction]
		screen.DrawImage(img, ops)
	}
}
