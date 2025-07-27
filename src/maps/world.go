package maps

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kkkunny/stl/container/pqueue"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/maps/render"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/util"
)

type World struct {
	tileCache       *render.TileCache
	currentMap      *Map
	pos             [2]int
	firstRenderTime time.Time
	// 切换地图
	nameMoveSpeed   int           // 地图名移动速度
	fontFace        *text.GoXFace // 显示地图名
	nameMoveCounter int           // 地图名移动计数器
}

func NewWorld(cfg *config.Config, initMapName string) (*World, error) {
	tileCache := render.NewTileCache()
	enterMap, err := NewMap(cfg, tileCache, initMapName)
	if err != nil {
		return nil, err
	}
	// 字体
	fontBytes, err := os.ReadFile(fmt.Sprintf("./resource/fonts/%s.ttf", cfg.MaterFontName))
	if err != nil {
		return nil, err
	}
	fontInst, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	fontFace, err := opentype.NewFace(fontInst, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}
	w := &World{
		tileCache:     tileCache,
		nameMoveSpeed: 1,
		fontFace:      text.NewGoXFace(fontFace),
	}
	w.MoveTo(enterMap)
	return w, nil
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
	if w.currentMap == targetMap {
		return
	}
	w.currentMap = targetMap
	w.nameMoveCounter = 0
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

// DrawMapName 绘制地图名
func (w *World) DrawMapName(cfg *config.Config, screen *ebiten.Image) error {
	height := cfg.ScreenHeight / 7
	if w.nameMoveCounter < 0 || w.nameMoveCounter >= height*4 {
		return nil
	}

	var ops ebiten.DrawImageOptions
	if w.nameMoveCounter < height {
		ops.GeoM.Translate(10, float64(w.nameMoveCounter-height))
	} else if w.nameMoveCounter < height*3 {
		ops.GeoM.Translate(10, 0)
	} else {
		ops.GeoM.Translate(10, -float64(w.nameMoveCounter%height))
	}
	screen.DrawImage(w.getMapNameDisplayImage(cfg), &ops)
	w.nameMoveCounter += w.nameMoveSpeed
	return nil
}

func (w *World) getMapNameDisplayImage(cfg *config.Config) *ebiten.Image {
	width, height := float32(cfg.ScreenWidth)/3, float32(cfg.ScreenHeight)/7
	img := ebiten.NewImage(int(width), int(height))

	vector.DrawFilledRect(img, 0, -6, width, height, util.NewRGBColor(248, 248, 255), false)
	vector.StrokeRect(img, 4, -4, width-8, height-6, 4, util.NewRGBColor(176, 196, 222), false)
	vector.StrokeRect(img, 0, -6, width, height, 6, util.NewRGBColor(119, 136, 153), false)

	displayText := w.currentMap.Name()
	bounds, _ := font.BoundString(w.fontFace.UnsafeInternal(), displayText)
	var textOps text.DrawOptions
	textOps.ColorScale.ScaleWithColor(color.Black)
	textOps.GeoM.Translate((float64(width)+10)/2-float64(bounds.Max.X.Floor()-bounds.Min.X.Floor())/2, (float64(height)-6)/2-float64(bounds.Max.Y.Floor()-bounds.Min.Y.Floor()))
	text.Draw(img, displayText, w.fontFace, &textOps)
	return img
}
