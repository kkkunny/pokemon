package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type Game struct {
	gameMap *image.NRGBA
}

func NewGame() (*Game, error) {
	g := &Game{}
	err := g.Init()
	return g, err
}

func (g *Game) Init() error {
	mapTMX, err := tiled.LoadFile("maps/hoenn/littleroot-town.tmx")
	if err != nil {
		return err
	}
	renderer, err := render.NewRenderer(mapTMX)
	if err != nil {
		return err
	}
	err = renderer.RenderVisibleLayers()
	if err != nil {
		return err
	}
	// defer renderer.Clear()
	g.gameMap = renderer.Result
	return nil
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(ebiten.NewImageFromImage(g.gameMap), nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	game, err := NewGame()
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pokemon")
	if err = ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
