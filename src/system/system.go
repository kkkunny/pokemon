package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/dialogue"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/sprite/person"
	"github.com/kkkunny/pokemon/src/voice"
)

type System struct {
	ctx            context.Context
	input          *input.System
	world          *maps.World
	self           person.Self
	mapVoicePlayer *voice.Player
	dialogue       *dialogue.DialogueSystem
}

func NewSystem(ctx context.Context) (*System, error) {
	// 地图
	world, err := maps.NewWorld(ctx.Config(), "Pallet_Town")
	if err != nil {
		return nil, err
	}
	// 对话系统
	ds, err := dialogue.NewDialogueSystem(ctx.Config())
	if err != nil {
		return nil, err
	}
	ds.ResetText("欢迎来到口袋妖怪世界！开始属于你的冒险吧！")
	ds.SetDisplay(true)
	// 主角
	self, err := person.NewSelf("master")
	if err != nil {
		return nil, err
	}
	self.SetPosition(6, 8)
	return &System{
		ctx:            ctx,
		input:          input.NewSystem(),
		world:          world,
		self:           self,
		mapVoicePlayer: voice.NewPlayer(),
		dialogue:       ds,
	}, nil
}

func (s *System) Update() error {
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

	drawInfo := &person.UpdateInfo{World: s.world}
	// 处理输入
	action, err := s.input.Action()
	if err != nil {
		return err
	}
	if action != nil {
		if !s.dialogue.Display() {
			err = s.self.OnAction(s.ctx, *action, drawInfo)
			if err != nil {
				return err
			}
			for _, sp := range s.world.CurrentMap().Sprites() {
				err = sp.OnAction(s.ctx, *action, drawInfo)
				if err != nil {
					return err
				}
			}
		} else if s.dialogue.StreamDone() && *action == input.ActionEnum.A {
			s.dialogue.SetDisplay(false)
		}
	}
	// 更新
	err = s.self.Update(s.ctx, drawInfo)
	if err != nil {
		return err
	}
	err = s.world.Update(s.ctx, []sprite.Sprite{s.self}, drawInfo)
	if err != nil {
		return err
	}
	return nil
}

func (s *System) Draw(screen *ebiten.Image) error {
	// 地图
	originSizeScreen := ebiten.NewImage(screen.Bounds().Dx()*s.ctx.Config().Scale, screen.Bounds().Dy()*s.ctx.Config().Scale)
	err := s.world.Draw(s.ctx, originSizeScreen, []sprite.Sprite{s.self})
	if err != nil {
		return err
	}

	var ops ebiten.DrawImageOptions
	ops.GeoM.Scale(float64(s.ctx.Config().Scale), float64(s.ctx.Config().Scale))
	ops.GeoM.Translate(float64(screen.Bounds().Dx()/2*(1-s.ctx.Config().Scale)), float64(screen.Bounds().Dy()/2*(1-s.ctx.Config().Scale)))
	screen.DrawImage(originSizeScreen, &ops)

	// 地图名
	err = s.world.DrawMapName(s.ctx, screen)
	if err != nil {
		return err
	}

	// 对话
	err = s.dialogue.Draw(screen)
	if err != nil {
		return err
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
	return nil
}
