package system

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/dialogue"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/sprite/person"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/image"
	"github.com/kkkunny/pokemon/src/voice"
)

type System struct {
	ctx            context.Context
	world          *maps.World
	self           person.Self
	mapVoicePlayer *voice.Player
	dialogue       *dialogue.System

	time time.Time // 游戏世界时间
}

func NewSystem(ctx context.Context) (*System, error) {
	// 地图
	world, err := maps.NewWorld(ctx.Config(), "pallet_town")
	if err != nil {
		return nil, err
	}
	// 对话系统
	ds, err := dialogue.NewSystem(ctx)
	if err != nil {
		return nil, err
	}
	ds.DisplayLabel("欢迎来到口袋妖怪世界！开始属于你的冒险吧！")
	// 主角
	self, err := person.NewSelf("master")
	if err != nil {
		return nil, err
	}
	self.SetPosition(6, 8)
	return &System{
		ctx:            ctx,
		world:          world,
		self:           self,
		mapVoicePlayer: voice.NewPlayer(),
		dialogue:       ds,
		time:           time.Now(),
	}, nil
}

func (s *System) OnAction(action input.KeyInputAction) error {
	if !s.dialogue.Display() {
		drawInfo := &person.UpdateInfo{World: s.world}
		err := s.self.OnAction(s.ctx, action, drawInfo)
		if err != nil {
			return err
		}
		for _, sp := range s.world.CurrentMap().Sprites() {
			err = sp.OnAction(s.ctx, action, drawInfo)
			if err != nil {
				return err
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
					text := s.ctx.Localisation().Get(targetSprite.GetText())
					s.dialogue.DisplayLabel(text)
				case sprite.ActionTypeEnum.Dialogue:
					movableSprite, ok := targetSprite.(sprite.MovableSprite)
					if ok {
						movableSprite.SetMovable(false)
					}
					text := s.ctx.Localisation().Get(targetSprite.GetText())
					s.dialogue.DisplayDialogue(text)
				}
			}
		}
	} else if s.dialogue.StreamDone() && action == input.KeyInputActionEnum.A.Pressed() {
		actionSprite := s.self.ActionSprite()
		if actionSprite != nil {
			s.self.SetActionSprite(nil)
			movableSprite, ok := actionSprite.(sprite.MovableSprite)
			if ok {
				movableSprite.SetMovable(true)
			}
		}
		s.dialogue.SetDisplay(false)
	}
	return nil
}

func (s *System) OnUpdate() error {
	// 地图音乐
	songFilepath, ok := s.world.CurrentMap().SongFilepath()
	if ok {
		err := s.mapVoicePlayer.LoadFile(songFilepath)
		if err != nil {
			return err
		}
		err = s.mapVoicePlayer.Play()
		if err != nil {
			return err
		}
	}

	// 时间
	s.time = s.time.Add(time.Minute)

	// 主角
	drawInfo := &person.UpdateInfo{World: s.world}
	err := s.self.Update(s.ctx, drawInfo)
	if err != nil {
		return err
	}
	// 世界
	err = s.world.Update(s.ctx, []sprite.Sprite{s.self}, drawInfo)
	if err != nil {
		return err
	}
	return nil
}

func (s *System) getSkyMaskColor() color.Color {
	hour, minute := float64(s.time.Hour()), float64(s.time.Minute())
	hour += minute / 60

	switch {
	case hour < 4:
		return util.NewRGBAColor(0, 0, 0, 180)
	case 4 <= hour && hour < 10:
		return util.GradientColor(util.NewRGBAColor(0, 0, 0, 180), util.NewRGBAColor(255, 255, 255, 0), (hour-4)/6)
	case 10 <= hour && hour < 15:
		return util.NewRGBAColor(255, 255, 255, 0)
	case 15 <= hour && hour < 17:
		return util.GradientColor(util.NewRGBAColor(255, 255, 255, 0), util.NewRGBAColor(255, 128, 64, 80), (hour-15)/2)
	case 17 <= hour && hour < 18:
		return util.GradientColor(util.NewRGBAColor(255, 128, 64, 80), util.NewRGBAColor(0, 0, 0, 180), (hour-17)/1)
	case 18 <= hour:
		return util.NewRGBAColor(0, 0, 0, 180)
	default:
		return util.NewRGBAColor(255, 255, 255, 0)
	}
}

func (s *System) OnDraw(screen *image.Image) error {
	// 地图
	originSizeScreen := image.NewImage(screen.Width()*s.ctx.Config().Scale, screen.Height()*s.ctx.Config().Scale)
	err := s.world.Draw(s.ctx, originSizeScreen, []sprite.Sprite{s.self})
	if err != nil {
		return err
	}

	var ops ebiten.DrawImageOptions
	ops.GeoM.Scale(float64(s.ctx.Config().Scale), float64(s.ctx.Config().Scale))
	ops.GeoM.Translate(float64(screen.Width()/2*(1-s.ctx.Config().Scale)), float64(screen.Height()/2*(1-s.ctx.Config().Scale)))
	screen.DrawImage(originSizeScreen, &ops)

	// 地图名
	err = s.world.DrawMapName(s.ctx, screen)
	if err != nil {
		return err
	}

	// 天色
	screen.Overlay(s.getSkyMaskColor())

	// 对话
	err = s.dialogue.Draw(screen)
	if err != nil {
		return err
	}

	return nil
}
