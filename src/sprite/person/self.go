package person

import (
	"errors"

	stlmaps "github.com/kkkunny/stl/container/maps"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/sprite"
)

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

	person.pos = [2]int{6, 8}
	return &_Self{_Person: *person}, nil
}

func (s *_Self) self() {}

func (s *_Self) Move() bool {
	return s.pos != s.expectPos
}

func (s *_Self) SetPosition(x, y int) {
	s.pos = [2]int{x, y}
	s.expectPos = [2]int{x, y}
}

func (s *_Self) Position() (int, int) {
	return s.pos[0], s.pos[1]
}

func (s *_Self) OnAction(cfg *config.Config, action input.Action, info sprite.UpdateInfo) {
	if info == nil {
		return
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return
	}

	if s.Move() {
		return
	}
	preDirection := s.direction
	switch action {
	case input.ActionEnum.MoveUp:
		s.direction = consts.DirectionEnum.Up
	case input.ActionEnum.MoveDown:
		s.direction = consts.DirectionEnum.Down
	case input.ActionEnum.MoveLeft:
		s.direction = consts.DirectionEnum.Left
	case input.ActionEnum.MoveRight:
		s.direction = consts.DirectionEnum.Right
	}
	if preDirection == s.direction {
		expectPos := s.pos
		switch s.direction {
		case consts.DirectionEnum.Up:
			expectPos = [2]int{s.pos[0], s.pos[1] + int(consts.DirectionEnum.Up)%2}
		case consts.DirectionEnum.Down:
			expectPos = [2]int{s.pos[0], s.pos[1] + int(consts.DirectionEnum.Down)%2}
		case consts.DirectionEnum.Left:
			expectPos = [2]int{s.pos[0] + int(consts.DirectionEnum.Left)%2, s.pos[1]}
		case consts.DirectionEnum.Right:
			expectPos = [2]int{s.pos[0] + int(consts.DirectionEnum.Right)%2, s.pos[1]}
		}
		if !updateInfo.World.CheckCollision(expectPos[0], expectPos[1]) {
			s.expectPos = expectPos
		}
	}
	return
}

func (s *_Self) PixelPosition(cfg *config.Config) (x, y int) {
	img := s.directionImages[s.direction]
	return cfg.ScreenWidth/2 - img.Bounds().Dx()/2, cfg.ScreenHeight/2 - img.Bounds().Dy()/2
}

func (s *_Self) Update(cfg *config.Config, info sprite.UpdateInfo) error {
	if info == nil {
		return errors.New("expect UpdateInfo")
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return errors.New("expect UpdateInfo")
	}

	// 移动
	if s.Move() {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
		a.SetFrameTime(cfg.TileSize / s.speed / a.FrameCount())
		a.Update()

		diff := cfg.TileSize - s.moveCounter
		if diff > s.speed {
			s.moveCounter += s.speed
		} else {
			s.moveCounter = 0
			targetMap, targetX, targetY, _ := updateInfo.World.GetActualPosition(s.expectPos[0], s.expectPos[1])
			updateInfo.World.MoveTo(targetMap)
			s.expectPos = [2]int{targetX, targetY}
			s.pos = s.expectPos
			s.moveStartingFoot = -s.moveStartingFoot
			a.Reset()
		}
	}

	// 更新地图位置
	img := s.directionImages[s.direction]
	x, y := s.pos[0]*cfg.TileSize, (s.pos[1]+1)*cfg.TileSize-img.Bounds().Dy()
	switch s.direction {
	case consts.DirectionEnum.Up:
		y -= s.moveCounter
	case consts.DirectionEnum.Down:
		y += s.moveCounter
	case consts.DirectionEnum.Left:
		x -= s.moveCounter
	case consts.DirectionEnum.Right:
		x += s.moveCounter
	}
	selfX, selfY := s.PixelPosition(cfg)
	x = -x + selfX
	y = -y + selfY
	updateInfo.World.MovePixelPosTo(x, y)
	return nil
}
