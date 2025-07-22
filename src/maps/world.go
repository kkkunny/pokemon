package maps

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/maps/render"
	"github.com/kkkunny/pokemon/src/util"
)

type World struct {
	tileCache  *render.TileCache
	currentMap *Map
	pos        [2]int
}

func NewWorld(cfg *config.Config, initMapName string) (*World, error) {
	tileCache := render.NewTileCache()
	enterMap, err := NewMap(cfg, tileCache, initMapName)
	if err != nil {
		return nil, err
	}
	return &World{
		tileCache:  tileCache,
		currentMap: enterMap,
	}, nil
}

func (w *World) Draw(cfg *config.Config, screen *ebiten.Image, sprites []util.Drawer) error {
	drawMaps := make(map[*Map]*ebiten.DrawImageOptions, len(w.currentMap.adjacentMaps)+1)

	x, y := float64(w.pos[0]), float64(w.pos[1])
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(x, y)
	drawMaps[w.currentMap] = ops

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
		ops = &ebiten.DrawImageOptions{}
		ops.GeoM.Translate(adjacentMapX, adjacentMapY)
		drawMaps[adjacentMap] = ops
	}

	for drawMap, ops := range drawMaps {
		err := drawMap.DrawBackground(screen, ops)
		if err != nil {
			return err
		}
	}
	for _, s := range sprites {
		err := s.Draw(cfg, screen, ops)
		if err != nil {
			return err
		}
	}
	for drawMap, ops := range drawMaps {
		err := drawMap.DrawForeground(screen, ops)
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
