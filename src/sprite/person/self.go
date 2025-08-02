package person

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	stlmaps "github.com/kkkunny/stl/container/maps"
	stlval "github.com/kkkunny/stl/value"
	lua "github.com/yuin/gopher-lua"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/script"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/util/image"
)

var keyInputActionToDirection = map[input.KeyInputAction]consts.Direction{
	input.KeyInputActionEnum.MoveUp:    consts.DirectionEnum.Up,
	input.KeyInputActionEnum.MoveDown:  consts.DirectionEnum.Down,
	input.KeyInputActionEnum.MoveLeft:  consts.DirectionEnum.Left,
	input.KeyInputActionEnum.MoveRight: consts.DirectionEnum.Right,
}

type Self interface {
	Person
	ActionSprite() sprite.Sprite
	SetActionSprite(sp sprite.Sprite)
}

type _Self struct {
	_Person

	actionSprite sprite.Sprite
}

func NewSelf(name string) (Self, error) {
	personObj, err := NewPerson(name)
	if err != nil {
		return nil, err
	}
	person := personObj.(*_Person)

	behaviorAnimations, err := loadPersonAnimations(name, sprite.BehaviorEnum.Run)
	if err != nil {
		return nil, err
	}
	person.behaviorAnimations = stlmaps.Union(person.behaviorAnimations, behaviorAnimations)

	person.SetPosition(6, 8)
	return &_Self{_Person: *person}, nil
}

func (s *_Self) OnAction(_ context.Context, action input.KeyInputAction, info sprite.UpdateInfo) error {
	if info == nil {
		return nil
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return nil
	}

	if s.Busying() {
		return nil
	}

	if s.movable { // 移动
		nextStepDirection, ok := keyInputActionToDirection[action]
		if ok {
			if s.direction != nextStepDirection {
				s.nextStepDirection = nextStepDirection
			} else if x, y := GetNextPositionByDirection(nextStepDirection, s.pos[0], s.pos[1]); !updateInfo.World.CheckCollision(x, y) {
				s.SetNextStepDirection(nextStepDirection)
			}
		}
	}

	return nil
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
			err := updateInfo.World.MoveTo(targetMap.ID())
			if err != nil {
				return err
			}
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

func (s *_Self) Draw(ctx context.Context, screen *image.Image, _ ebiten.DrawImageOptions) error {
	x, y := s.PixelPosition(ctx.Config())
	var ops ebiten.DrawImageOptions
	ops.GeoM.Translate(float64(x), float64(y))

	if s.Turning() {
		if s.direction == -s.nextStepDirection {
			s.moveStartingFoot = FootEnum.Right
		} else if s.direction == consts.DirectionEnum.Up {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Left, FootEnum.Left, FootEnum.Right)
		} else if s.direction == consts.DirectionEnum.Down {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Right, FootEnum.Left, FootEnum.Right)
		} else if s.direction == consts.DirectionEnum.Left {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Down, FootEnum.Left, FootEnum.Right)
		} else if s.direction == consts.DirectionEnum.Right {
			s.moveStartingFoot = stlval.Ternary(s.nextStepDirection == consts.DirectionEnum.Up, FootEnum.Left, FootEnum.Right)
		}
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.nextStepDirection][s.moveStartingFoot]
		screen.DrawImage(a.GetFrameImage(1), &ops)
	} else {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.nextStepDirection][s.moveStartingFoot]
		a.Draw(screen, ops)
	}
	return nil
}

func (s *_Self) ActionSprite() sprite.Sprite {
	return s.actionSprite
}

func (s *_Self) SetActionSprite(sp sprite.Sprite) {
	s.actionSprite = sp
}

var luaModuleToGo = map[string]map[string]lua.LGFunction{
	"sprite": {
		"set_movable": func(rt *lua.LState) int {
			param1 := rt.CheckUserData(1)
			s, ok := param1.Value.(sprite.MovableSprite)
			if !ok {
				return 1
			}
			movable := rt.CheckBool(2)
			s.SetMovable(movable)
			return 0
		},
	},
	"game": {
		"display_message": func(rt *lua.LState) int {
			param1 := rt.CheckUserData(1)
			s, ok := param1.Value.(sprite.MovableSprite)
			if !ok {
				return 1
			}
			movable := rt.CheckBool(2)
			s.SetMovable(movable)
			return 0
		},
	},
}

func loadScriptFileWithSelf(w *maps.World, this sprite.Sprite, master *_Self, name string) (rt *lua.LState, err error) {
	rt, err = script.LoadScriptFile(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			rt.Close()
		}
	}()

	err = rt.PCall(0, lua.MultRet, nil)
	if err != nil {
		return nil, err
	}
	for moduleName, luaFuncToGo := range luaModuleToGo {
		rt.PreloadModule(moduleName, func(rt *lua.LState) int {
			rt.Push(rt.SetFuncs(rt.NewTable(), luaFuncToGo))
			return 1
		})
	}

	return rt, nil
}
