package src

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/sprite"
)

type Game struct {
	input *input.System

	gameMap *image.NRGBA
	sprites []sprite.Sprite
}

func NewGame() (*Game, error) {
	g := &Game{
		input: input.NewSystem(),
	}
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
	masterCharacter, err := sprite.NewSelf("master")
	if err != nil {
		return err
	}
	g.sprites = append(g.sprites, masterCharacter)
	return nil
}

func (g *Game) Update() error {
	// 处理输入
	action, err := g.input.Action()
	if err != nil {
		return err
	}
	if action != nil {
		for _, s := range g.sprites {
			err := s.OnAction(*action)
			if err != nil {
				return err
			}
		}
	}
	// 更新帧
	for _, s := range g.sprites {
		err = s.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(ebiten.NewImageFromImage(g.gameMap), nil)
	for _, s := range g.sprites {
		x, y, display := s.Position()
		if !display {
			continue
		}
		spriteImg, err := s.Image()
		if err != nil {
			fmt.Println(err)
			continue
		}
		ops := &ebiten.DrawImageOptions{}
		ops.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(spriteImg, ops)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
