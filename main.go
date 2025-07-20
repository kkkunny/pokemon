package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src"
	"github.com/kkkunny/pokemon/src/config"
)

func main() {
	cfg := config.NewConfig()
	cfg.Scale = 2
	game, err := src.NewGame(cfg)
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle("Pokemon")
	if err = ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
