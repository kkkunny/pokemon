package src

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/sprite"
)

type Game struct {
	input *input.System

	gameMap *maps.Map
	sprites []sprite.Sprite
}

func NewGame() (*Game, error) {
	g := &Game{
		input: input.NewSystem(),
	}
	err := g.Init()
	return g, err
}

func (g *Game) Init() (err error) {
	// 地图
	g.gameMap, err = maps.NewMap()
	if err != nil {
		return err
	}
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
			s.OnAction(*action)
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
	screen.DrawImage(g.gameMap.Image(), nil)
	drawInfo := &sprite.DrawInfo{Person: &sprite.PersonDrawInfo{Map: g.gameMap}}
	for _, s := range g.sprites {
		s.Draw(screen, drawInfo)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
