package main

import (
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src"
	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/i18n"
)

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
	game, err := src.NewGame()
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle(game.Name())
	if err = ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
