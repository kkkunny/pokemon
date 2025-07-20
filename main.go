package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src"
)

func main() {
	game, err := src.NewGame()
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(240, 160)
	ebiten.SetWindowTitle("Pokemon")
	if err = ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
