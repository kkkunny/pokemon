package system

import (
	"image/color"
	"time"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/output/voice"
	"github.com/kkkunny/pokemon/src/system"
	"github.com/kkkunny/pokemon/src/system/world"
	"github.com/kkkunny/pokemon/src/system/world/sprite"
	"github.com/kkkunny/pokemon/src/system/world/sprite/person"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
)

type WorldSystem struct {
	world  *world.World // 世界
	self   person.Self  // 主角
	time   time.Time    // 游戏世界时间
	player *system.RealTimeVoicePlayer
}

func NewWorldSystem() (*WorldSystem, error) {
	// 地图
	w, err := world.NewWorld()
	if err != nil {
		return nil, err
	}
	err = w.MoveTo("pallet_town")
	if err != nil {
		return nil, err
	}

	// 主角
	self, err := person.NewSelf("master")
	if err != nil {
		return nil, err
	}
	self.SetPosition(6, 8)

	return &WorldSystem{
		world:  w,
		self:   self,
		time:   time.Now(),
		player: system.NewRealTimeVoicePlayer(voice.NewPlayer()),
	}, nil
}

func (s *WorldSystem) OnAction(sm system.SystemManager, action input.KeyInputAction) error {
	drawInfo := &person.UpdateInfo{World: s.world}
	err := s.self.OnAction(action, drawInfo)
	if err != nil {
		return err
	}
	for _, sp := range s.world.CurrentMap().Sprites() {
		err = sp.OnAction(action, drawInfo)
		if err != nil {
			return err
		}
	}

	actionSprite := s.self.ActionSprite()
	if actionSprite != nil {
		s.self.SetActionSprite(nil)
		movableSprite, ok := actionSprite.(sprite.MovableSprite)
		if ok {
			movableSprite.SetMovable(true)
		}
	}

	if action == input.KeyInputActionEnum.A.Pressed() {
		x, y := s.self.Position()
		targetX, targetY := person.GetNextPositionByDirection(s.self.Direction(), x, y)
		targetMap, targetX, targetY, _ := s.world.GetActualPosition(targetX, targetY)
		targetSprite, ok := targetMap.GetSpriteByPosition(targetX, targetY)
		if ok {
			s.self.SetActionSprite(targetSprite)
			switch targetSprite.ActionType() {
			case sprite.ActionTypeEnum.Script:
				// scriptName := targetSprite.GetScript()
				// rt, err := loadScriptFileWithSelf(updateInfo.World, targetSprite, s, scriptName)
				// if err != nil {
				// 	return err
				// }
				// defer rt.Close()
				//
				// param1 := rt.NewUserData()
				// param1.Value = targetSprite
				// err = rt.CallByParam(lua.P{
				// 	Fn:      rt.GetGlobal(scriptName),
				// 	NRet:    1,
				// 	Protect: true,
				// }, param1)
				// if err != nil {
				// 	return err
				// }
			case sprite.ActionTypeEnum.Label:
				text := i18n.Get(targetSprite.GetText())
				return sm.DisplayText(system.BoxStyleEnum.Label, text, util.NewNRGBColor(100, 100, 100))
			case sprite.ActionTypeEnum.Dialogue:
				movableSprite, ok := targetSprite.(sprite.MovableSprite)
				if ok {
					movableSprite.SetMovable(false)
				}
				text := i18n.Get(targetSprite.GetText())
				return sm.DisplayText(system.BoxStyleEnum.Label, text, util.NewNRGBColor(100, 100, 100))
			}
		}
	}
	return nil
}

func (s *WorldSystem) OnUpdate(system system.SystemManager) error {
	// TODO: MOCK
	return system.StartOneBattle("grassland")

	// 地图音乐
	songFilepath, ok := s.world.CurrentMap().SongFilepath()
	if ok {
		err := s.player.LoadFile(songFilepath)
		if err != nil {
			return err
		}
		err = s.player.Play()
		if err != nil {
			return err
		}
	}

	// 时间
	s.time = s.time.Add(time.Minute)

	// 主角
	drawInfo := &person.UpdateInfo{World: s.world}
	err := s.self.Update(drawInfo)
	if err != nil {
		return err
	}
	// 世界
	return s.world.Update([]sprite.Sprite{s.self}, drawInfo)
}

func (s *WorldSystem) Drop() bool {
	return false
}

func (s *WorldSystem) OnDraw(system system.SystemManager, drawer draw.OptionDrawer) error {
	// 地图
	err := s.world.OnDraw(drawer.Scale(config.Scale, config.Scale), []sprite.Sprite{s.self})
	if err != nil {
		return err
	}

	// 天色
	if !s.world.CurrentMap().Indoor() {
		draw.OverlayColor(drawer, s.getSkyMaskColor())
	}

	// 地图名
	err = s.world.DrawMapName(drawer)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorldSystem) getSkyMaskColor() color.Color {
	hour, minute := float64(s.time.Hour()), float64(s.time.Minute())
	hour += minute / 60

	switch {
	case hour < 4:
		return util.NewNRGBAColor(0, 0, 0, 180)
	case 4 <= hour && hour < 10:
		return util.GradientColor(util.NewNRGBAColor(0, 0, 0, 180), util.NewNRGBAColor(255, 255, 255, 0), (hour-4)/6)
	case 10 <= hour && hour < 15:
		return util.NewNRGBAColor(255, 255, 255, 0)
	case 15 <= hour && hour < 17:
		return util.GradientColor(util.NewNRGBAColor(255, 255, 255, 0), util.NewNRGBAColor(255, 128, 64, 80), (hour-15)/2)
	case 17 <= hour && hour < 18:
		return util.GradientColor(util.NewNRGBAColor(255, 128, 64, 80), util.NewNRGBAColor(0, 0, 0, 180), (hour-17)/1)
	case 18 <= hour:
		return util.NewNRGBAColor(0, 0, 0, 180)
	default:
		return util.NewNRGBAColor(255, 255, 255, 0)
	}
}
