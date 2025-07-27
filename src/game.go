package src

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/system"
)

type Game struct {
	cfg *config.Config
	loc *i18n.Localisation
	sys *system.System
}

func NewGame(cfg *config.Config) (*Game, error) {
	// 翻译
	loc, err := i18n.LoadLocalisation(i18n.LanguageEnum.ZH_CN)
	if err != nil {
		return nil, err
	}
	sys, err := system.NewSystem(context.NewContext(cfg, loc))
	if err != nil {
		return nil, err
	}
	return &Game{
		cfg: cfg,
		loc: loc,
		sys: sys,
	}, err
}

func (g *Game) Update() error {
	return g.sys.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := g.sys.Draw(screen)
	if err != nil {
		panic(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
