package src

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/util/draw"
	i18n2 "github.com/kkkunny/pokemon/src/util/i18n"
)

type Game struct {
	cfg   *config.Config
	loc   *i18n2.Localisation
	input *input.System
	sys   *system.System
}

func NewGame(cfg *config.Config) (*Game, error) {
	// 翻译
	loc, err := i18n2.LoadLocalisation(i18n2.LanguageEnum.ZH_CN)
	if err != nil {
		return nil, err
	}
	sys, err := system.NewSystem(context.NewContext(cfg, loc))
	if err != nil {
		return nil, err
	}
	return &Game{
		cfg:   cfg,
		loc:   loc,
		input: input.NewSystem(),
		sys:   sys,
	}, err
}

func (g *Game) Name() string {
	return g.loc.Get("game_name")
}

func (g *Game) Update() error {
	action, err := g.input.KeyInputAction()
	if err != nil {
		return err
	}
	if action != nil {
		err = g.sys.OnAction(*action)
		if err != nil {
			return err
		}
	}
	return g.sys.OnUpdate()
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := g.sys.OnDraw(draw.NewDrawerFromEbiten(screen))
	if err != nil {
		panic(err)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
