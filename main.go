package main

import (
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src"
	"github.com/kkkunny/pokemon/src/config"
)

func main() {
	icon, err := imaging.Open(filepath.Join(config.GFXPath, "icon.png"))
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowIcon([]image.Image{icon})
	cfg := config.NewConfig()
	game, err := src.NewGame(cfg)
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle(game.Name())
	if err = ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
