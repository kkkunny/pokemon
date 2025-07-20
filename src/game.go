package src

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/util"
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
	masterCharacter.SetPosition(6, 8)
	g.sprites = append(g.sprites, masterCharacter)
	return nil
}

func (g *Game) Update() error {
	drawInfo := &sprite.UpdateInfo{Person: &sprite.PersonUpdateInfo{Map: g.gameMap}}
	// 处理输入
	action, err := g.input.Action()
	if err != nil {
		return err
	}
	if action != nil {
		for _, s := range g.sprites {
			s.OnAction(*action, drawInfo)
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
	drawer := make([]util.Drawer, len(g.sprites))
	for i, s := range g.sprites {
		drawer[i] = s
	}
	err := g.gameMap.Draw(screen, drawer)
	if err != nil {
		panic(err)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
