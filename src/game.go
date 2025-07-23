package src

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/voice"
)

type Game struct {
	cfg   *config.Config
	input *input.System

	world   *maps.World
	sprites []sprite.Sprite

	mapVoicePlayer *voice.Player
}

func NewGame(cfg *config.Config) (*Game, error) {
	g := &Game{
		cfg:            cfg,
		input:          input.NewSystem(),
		mapVoicePlayer: voice.NewPlayer(),
	}
	err := g.Init()
	return g, err
}

func (g *Game) Init() (err error) {
	// 地图
	g.world, err = maps.NewWorld(g.cfg, "Pallet_Town")
	if err != nil {
		return err
	}
	// 主角
	masterCharacter, err := sprite.NewSelf("master")
	if err != nil {
		return err
	}
	masterCharacter.SetPosition(6, 8)
	g.sprites = append(g.sprites, masterCharacter)
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

	drawInfo := &sprite.UpdateInfo{Person: &sprite.PersonUpdateInfo{World: g.world}}
	// 处理输入
	action, err := g.input.Action()
	if err != nil {
		return err
	}
	if action != nil {
		for _, s := range g.sprites {
			s.OnAction(g.cfg, *action, drawInfo)
		}
	}
	// 更新帧
	for _, s := range g.sprites {
		err = s.Update(g.cfg, drawInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	originSizeScreen := ebiten.NewImage(screen.Bounds().Dx()*g.cfg.Scale, screen.Bounds().Dy()*g.cfg.Scale)
	drawer := make([]util.Drawer, len(g.sprites))
	for i, s := range g.sprites {
		drawer[i] = s
	}
	err := g.world.Draw(g.cfg, originSizeScreen, drawer)
	if err != nil {
		panic(err)
	}

	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Scale(float64(g.cfg.Scale), float64(g.cfg.Scale))
	ops.GeoM.Translate(float64(screen.Bounds().Dx()/2*(1-g.cfg.Scale)), float64(screen.Bounds().Dy()/2*(1-g.cfg.Scale)))
	screen.DrawImage(originSizeScreen, ops)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
