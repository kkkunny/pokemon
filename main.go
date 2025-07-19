package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/kkkunny/pokemon/src/sprite"
)

type Game struct {
	gameMap *image.NRGBA

	sprites []sprite.Sprite
}

func NewGame() (*Game, error) {
	g := &Game{}
	err := g.Init()
	return g, err
}

func (g *Game) Init() error {
	// 地图
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
	// 主角
	masterCharacter, err := sprite.NewTrainer("master")
	if err != nil {
		return err
	}
	g.sprites = append(g.sprites, masterCharacter)
	return nil
}

func (g *Game) Update() error {
	for _, s := range g.sprites {
		err := s.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(ebiten.NewImageFromImage(g.gameMap), nil)
	for _, s := range g.sprites {
		spriteImg, err := s.Image()
		if err != nil {
			fmt.Println(err)
			continue
		}
		screen.DrawImage(ebiten.NewImageFromImage(spriteImg), nil)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
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
