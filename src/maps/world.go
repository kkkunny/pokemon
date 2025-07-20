package maps

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/util"
)

type World struct {
	currentMap *Map
	pos        [2]int
}

func NewWorld(cfg *config.Config, initMapName string) (*World, error) {
	enterMap, err := NewMap(cfg, initMapName)
	if err != nil {
		return nil, err
	}
	return &World{currentMap: enterMap}, nil
}

func (w *World) Draw(cfg *config.Config, screen *ebiten.Image, sprites []util.Drawer) error {
	x, y := float64(w.pos[0]), float64(w.pos[1])
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(x, y)
	err := w.currentMap.Draw(cfg, screen, sprites, ops)
	if err != nil {
		return err
	}
	currentMapW, currentMapH := w.currentMap.PixelSize()
	for direction, adjacentMap := range w.currentMap.adjacentMaps {
		width, height := adjacentMap.PixelSize()
		adjacentMapX, adjacentMapY := x, y
		switch direction {
		case consts.DirectionEnum.Up:
			adjacentMapY -= float64(height)
		case consts.DirectionEnum.Down:
			adjacentMapY += float64(currentMapH)
		case consts.DirectionEnum.Left:
			adjacentMapX -= float64(width)
		case consts.DirectionEnum.Right:
			adjacentMapX += float64(currentMapW)
		}
		directionOps := &ebiten.DrawImageOptions{}
		directionOps.GeoM.Translate(adjacentMapX, adjacentMapY)
		err = adjacentMap.Draw(cfg, screen, nil, directionOps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) MovePixelPosTo(x, y int) {
	w.pos = [2]int{x, y}
}

func (w *World) MoveTo(targetMap *Map) {
	w.currentMap = targetMap
}

func (w *World) GetActualPosition(x, y int) (*Map, int, int, bool) {
	return w.currentMap.GetActualPosition(x, y)
}

func (w *World) CheckCollision(x, y int) bool {
	targetMap, x, y, ok := w.GetActualPosition(x, y)
	if !ok {
		return true
	}
	return targetMap.CheckCollision(x, y)
}
