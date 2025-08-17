package main

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src"
	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type Game struct {
	input *input.System
}

func NewGame() (*Game, error) {
	err := src.Init()
	if err != nil {
		return nil, err
	}
	return &Game{
		input: input.NewSystem(),
	}, err
}

func (g *Game) Name() string {
	return i18n.Get("game_name")
}

func (g *Game) Update() error {
	action, err := g.input.KeyInputAction()
	if err != nil {
		return err
	}
	if action != nil {
		err = src.OnAction(*action)
		if err != nil {
			return err
		}
	}
	return src.OnUpdate()
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := src.OnDraw(draw.NewDrawerFromImage(imgutil.WrapImage(screen)))
	if err != nil {
		panic(err)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	config.ScreenWidth, config.ScreenHeight = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	config.Init()
	if err := i18n.Init(); err != nil {
		panic(err)
	}

	icon, err := imaging.Open(filepath.Join(consts.GFXPath, "icon.png"))
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowIcon([]image.Image{icon})
	game, err := NewGame()
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle(game.Name())
	if err = ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
