package src

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/sprite/person"
	"github.com/kkkunny/pokemon/src/voice"
)

type Game struct {
	ctx   context.Context
	input *input.System

	world *maps.World
	self  person.Self

	mapVoicePlayer *voice.Player
}

func NewGame(cfg *config.Config) (*Game, error) {
	g := &Game{
		input:          input.NewSystem(),
		mapVoicePlayer: voice.NewPlayer(),
	}
	err := g.Init(cfg)
	return g, err
}

func (g *Game) Init(cfg *config.Config) (err error) {
	// 翻译
	locs, err := i18n.LoadLocalisation(i18n.LanguageEnum.ZH_CN)
	if err != nil {
		return
	}
	g.ctx = context.NewContext(cfg, locs)
	// 地图
	g.world, err = maps.NewWorld(cfg, "Pallet_Town")
	if err != nil {
		return err
	}
	// 主角
	g.self, err = person.NewSelf("master")
	if err != nil {
		return err
	}
	g.self.SetPosition(6, 8)
	return nil
}

func (g *Game) Update() error {
	// 地图音乐
	songFilepath, ok := g.world.CurrentMap().SongFilepath()
	if ok {
		err := g.mapVoicePlayer.LoadFile(songFilepath)
		if err != nil {
			return err
		}
		err = g.mapVoicePlayer.Play()
		if err != nil {
			return err
		}
	}

	drawInfo := &person.UpdateInfo{World: g.world}
	// 处理输入
	action, err := g.input.Action()
	if err != nil {
		return err
	}
	if action != nil {
		err = g.self.OnAction(g.ctx, *action, drawInfo)
		if err != nil {
			return err
		}
		for _, s := range g.world.CurrentMap().Sprites() {
			err = s.OnAction(g.ctx, *action, drawInfo)
			if err != nil {
				return err
			}
		}
	}
	// 更新
	err = g.self.Update(g.ctx, drawInfo)
	if err != nil {
		return err
	}
	err = g.world.Update(g.ctx, []sprite.Sprite{g.self}, drawInfo)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	originSizeScreen := ebiten.NewImage(screen.Bounds().Dx()*g.ctx.Config().Scale, screen.Bounds().Dy()*g.ctx.Config().Scale)
	err := g.world.Draw(g.ctx, originSizeScreen, []sprite.Sprite{g.self})
	if err != nil {
		panic(err)
	}

	var ops ebiten.DrawImageOptions
	ops.GeoM.Scale(float64(g.ctx.Config().Scale), float64(g.ctx.Config().Scale))
	ops.GeoM.Translate(float64(screen.Bounds().Dx()/2*(1-g.ctx.Config().Scale)), float64(screen.Bounds().Dy()/2*(1-g.ctx.Config().Scale)))
	screen.DrawImage(originSizeScreen, &ops)

	// 地图名
	err = g.world.DrawMapName(g.ctx, screen)
	if err != nil {
		panic(err)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
