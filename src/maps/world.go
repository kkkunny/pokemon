package maps

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kkkunny/stl/container/pqueue"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/maps/render"
	"github.com/kkkunny/pokemon/src/sprite"
)

type World struct {
	tileCache       *render.TileCache
	currentMap      *Map
	pos             [2]int
	firstRenderTime time.Time
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

func (w *World) Draw(cfg *config.Config, screen *ebiten.Image, sprites []sprite.Sprite) error {
	now := time.Now()
	var defaultTime time.Time
	if w.firstRenderTime == defaultTime {
		w.firstRenderTime = now
	}

	drawMaps := make(map[*Map]ebiten.DrawImageOptions, len(w.currentMap.adjacentMaps)+1)

	x, y := float64(w.pos[0]), float64(w.pos[1])
	var ops ebiten.DrawImageOptions
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
		var adjacentMapOps ebiten.DrawImageOptions
		adjacentMapOps.GeoM.Translate(adjacentMapX, adjacentMapY)
		drawMaps[adjacentMap] = adjacentMapOps
	}

	// 背景
	for drawMap, ops := range drawMaps {
		err := drawMap.DrawBackground(screen, ops, now.Sub(w.firstRenderTime))
		if err != nil {
			return err
		}
	}
	// 精灵
	drawSprites := pqueue.AnyWith[int, sprite.Sprite]()
	// 全局精灵
	for _, s := range sprites {
		_, y := s.Position()
		drawSprites.Push(y, s)
	}
	// 地图精灵
	for _, s := range w.currentMap.sprites {
		_, y := s.Position()
		drawSprites.Push(y, s)
	}
	spritePairs := drawSprites.ToSlice()
	for i := len(spritePairs) - 1; i >= 0; i-- {
		err := spritePairs[i].E2().Draw(cfg, screen, ops)
		if err != nil {
			return err
		}
	}
	// 前景
	for drawMap, ops := range drawMaps {
		err := drawMap.DrawForeground(screen, ops, now.Sub(w.firstRenderTime))
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

func (w *World) CurrentMap() *Map {
	return w.currentMap
}

func (w *World) CheckCollision(x, y int) bool {
	targetMap, x, y, ok := w.GetActualPosition(x, y)
	if !ok {
		return true
	}
	return targetMap.CheckCollision(x, y)
}
