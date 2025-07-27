package person

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	stlmaps "github.com/kkkunny/stl/container/maps"
	stlval "github.com/kkkunny/stl/value"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/sprite"
)

var actionToDirection = map[input.Action]consts.Direction{
	input.ActionEnum.MoveUp:    consts.DirectionEnum.Up,
	input.ActionEnum.MoveDown:  consts.DirectionEnum.Down,
	input.ActionEnum.MoveLeft:  consts.DirectionEnum.Left,
	input.ActionEnum.MoveRight: consts.DirectionEnum.Right,
}

type Self interface {
	Person
	self()
}

type _Self struct {
	_Person
}

func NewSelf(name string) (Self, error) {
	personObj, err := NewPerson(name)
	if err != nil {
		return nil, err
	}
	person := personObj.(*_Person)

	behaviorAnimations, err := sprite.LoadPersonAnimations(name, sprite.BehaviorEnum.Run)
	if err != nil {
		return nil, err
	}
	person.behaviorAnimations = stlmaps.Union(person.behaviorAnimations, behaviorAnimations)

	person.SetPosition(6, 8)
	return &_Self{_Person: *person}, nil
}

func (s *_Self) self() {}

func (s *_Self) OnAction(_ context.Context, action input.Action, info sprite.UpdateInfo) {
	if info == nil {
		return
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return
	}

	if s.Busying() {
		return
	}

	nextStepDirection := actionToDirection[action]
	if s.direction != nextStepDirection {
		s.nextStepDirection = nextStepDirection
	} else if s.SetNextStepDirection(nextStepDirection) && updateInfo.World.CheckCollision(s.nextStepPos[0], s.nextStepPos[1]) {
		s.nextStepPos = s.pos
	}
}

func (s *_Self) PixelPosition(cfg *config.Config) (x, y int) {
	bounds := stlmaps.First(stlmaps.First(s.behaviorAnimations[sprite.BehaviorEnum.Walk]).E2()).E2().GetFrameImage(0).Bounds()
	return cfg.ScreenWidth/2 - bounds.Dx()/2, cfg.ScreenHeight/2 - bounds.Dy()/2
}

func (s *_Self) Update(ctx context.Context, info sprite.UpdateInfo) error {
	if info == nil {
		return errors.New("expect UpdateInfo")
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return errors.New("expect UpdateInfo")
	}

	if s.Turning() {
		if s.moveCounter < ctx.Config().TileSize {
			s.moveCounter += 2
		} else {
			s.moveCounter = 0
			s.direction = s.nextStepDirection
			s.moveStartingFoot = -s.moveStartingFoot
		}
	} else if s.Moving() {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.nextStepDirection][s.moveStartingFoot]
		a.SetFrameTime(ctx.Config().TileSize / s.speed / a.FrameCount())
		a.Update()

		diff := ctx.Config().TileSize - s.moveCounter
		if diff > s.speed {
			s.moveCounter += s.speed
		} else {
			s.moveCounter = 0
			targetMap, targetX, targetY, _ := updateInfo.World.GetActualPosition(s.nextStepPos[0], s.nextStepPos[1])
			updateInfo.World.MoveTo(targetMap)
			s.nextStepPos = [2]int{targetX, targetY}
			s.pos = s.nextStepPos
			s.moveStartingFoot = -s.moveStartingFoot
			a.Reset()
		}
	}

	// 更新地图位置
	x, y := s._Person.PixelPosition(ctx.Config())
	selfX, selfY := s.PixelPosition(ctx.Config())
	x = -x + selfX
	y = -y + selfY
	updateInfo.World.MovePixelPosTo(x, y)
	return nil
}

func (s *_Self) Draw(ctx context.Context, screen *ebiten.Image, _ ebiten.DrawImageOptions) error {
	x, y := s.PixelPosition(ctx.Config())
	var ops ebiten.DrawImageOptions
	ops.GeoM.Translate(float64(x), float64(y))

	if s.Turning() {
		if s.direction == -s.nextStepDirection {
			s.moveStartingFoot = sprite.FootEnum.Right
		} else if s.direction == consts.DirectionEnum.Up {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Left, sprite.FootEnum.Left, sprite.FootEnum.Right)
		} else if s.direction == consts.DirectionEnum.Down {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Right, sprite.FootEnum.Left, sprite.FootEnum.Right)
		} else if s.direction == consts.DirectionEnum.Left {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Down, sprite.FootEnum.Left, sprite.FootEnum.Right)
		} else if s.direction == consts.DirectionEnum.Right {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Up, sprite.FootEnum.Left, sprite.FootEnum.Right)
		}
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.nextStepDirection][s.moveStartingFoot]
		screen.DrawImage(a.GetFrameImage(1), &ops)
	} else {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.nextStepDirection][s.moveStartingFoot]
		a.Draw(screen, ops)
	}
	return nil
}
