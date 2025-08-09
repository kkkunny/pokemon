package world

import (
	"image"
	"image/color"
	"time"

	stlmaps "github.com/kkkunny/stl/container/maps"
	"github.com/kkkunny/stl/container/pqueue"
	"golang.org/x/image/font"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/system/world/render"
	"github.com/kkkunny/pokemon/src/system/world/sprite"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type World struct {
	ctx             context.Context
	tileCache       *render.TileCache
	mapCache        map[string]*Map
	currentMap      *Map
	pixPos          [2]int
	firstRenderTime time.Time

	// 地图名
	nameMoveSpeed   int // 地图名移动速度
	nameMoveCounter int // 地图名移动计数器

	// 缓存地图碰撞
	selfPos [2]int // 主角所在当前地图位置
}

func NewWorld(ctx context.Context, initMapName string) (*World, error) {
	tileCache := render.NewTileCache()
	w := &World{
		ctx:           ctx,
		tileCache:     tileCache,
		mapCache:      make(map[string]*Map),
		nameMoveSpeed: 1,
	}
	return w, w.MoveTo(initMapName)
}

func (w *World) Update(ctx context.Context, sprites []sprite.Sprite, info sprite.UpdateInfo) error {
	// 全局精灵
	var selfX, selfY int
	for _, s := range sprites {
		if !s.Collision() {
			continue
		}
		selfX, selfY = s.Position()
		x, y := s.CollisionPosition()
		w.selfPos = [2]int{x, y}
	}
	// 地图精灵
	for _, s := range w.CurrentMap().Sprites() {
		err := s.Update(ctx, info)
		if err != nil {
			return err
		}
	}

	// 洞
	for _, hole := range w.currentMap.GetHoles() {
		x, y := int(hole.X+hole.Width/2)/config.TileSize, int(hole.Y+hole.Height/2)/config.TileSize
		if selfX != x || selfY != y {
			continue
		}
		toMap, toX, toY := hole.Properties.GetString("to_map"), hole.Properties.GetInt("to_x"), hole.Properties.GetInt("to_y")
		err := w.MoveTo(toMap)
		if err != nil {
			return err
		}
		for _, s := range sprites {
			s.SetPosition(toX, toY)
		}
		break
	}
	return nil
}

// 获取需要绘制的地图信息（参数、范围）
func (w *World) getNeedDrawMap() (map[*Map]image.Point, map[*Map]image.Rectangle, error) {
	map2Pos := make(map[*Map]image.Point, 5)
	map2Rect := make(map[*Map]image.Rectangle, 5)

	var loopFn func(m *Map, pixX, pixY int) error
	loopFn = func(m *Map, pixX, pixY int) error {
		map2Pos[m] = image.Pt(pixX, pixY)
		x0, y0 := max(0-pixX, 0)/config.TileSize, max(0-pixY, 0)/config.TileSize
		mapPixWidth, mapPixHeight := m.PixelSize()
		mapWidth, mapHeight := m.Size()
		x1, y1 := mapWidth-max((pixX+mapPixWidth)*config.Scale-w.ctx.Config().ScreenWidth, 0)/(config.TileSize*config.Scale), mapHeight-max((pixY+mapPixHeight)*config.Scale-w.ctx.Config().ScreenHeight, 0)/(config.TileSize*config.Scale)
		map2Rect[m] = image.Rect(x0, y0, x1, y1)

		needDrawAdjacentMaps := stlmaps.Filter(m.AdjacentMaps(), func(d util.Direction, id string) bool {
			switch d {
			case util.DirectionEnum.Up:
				return pixY > 0
			case util.DirectionEnum.Down:
				return (pixY+mapPixHeight)*config.Scale < w.ctx.Config().ScreenHeight
			case util.DirectionEnum.Left:
				return pixX > 0
			case util.DirectionEnum.Right:
				return (pixX+mapPixWidth)*config.Scale < w.ctx.Config().ScreenWidth
			default:
				return false
			}
		})

		for d, adjacentMapID := range needDrawAdjacentMaps {
			adjacentMap, err := w.loadMap(adjacentMapID)
			if err != nil {
				return err
			} else if stlmaps.ContainKey(map2Pos, adjacentMap) {
				continue
			}
			adjacentPixX, adjacentPixY := pixX, pixY
			adjacentPixWidth, adjacentPixHeight := adjacentMap.PixelSize()
			switch d {
			case util.DirectionEnum.Up:
				adjacentPixY -= adjacentPixHeight
			case util.DirectionEnum.Down:
				adjacentPixY += mapPixHeight
			case util.DirectionEnum.Left:
				adjacentPixX -= adjacentPixWidth
			case util.DirectionEnum.Right:
				adjacentPixX += mapPixWidth
			}
			loopFn(adjacentMap, adjacentPixX, adjacentPixY)
		}
		return nil
	}
	err := loopFn(w.currentMap, w.pixPos[0], w.pixPos[1])
	if err != nil {
		return nil, nil, err
	}

	return map2Pos, map2Rect, nil
}

func (w *World) OnDraw(drawer draw.OptionDrawer, sprites []sprite.Sprite) error {
	now := time.Now()
	var defaultTime time.Time
	if w.firstRenderTime == defaultTime {
		w.firstRenderTime = now
	}

	map2Pos, map2Rect, err := w.getNeedDrawMap()
	if err != nil {
		return err
	}

	// 背景
	for drawMap, pos := range map2Pos {
		err = drawMap.DrawBackground(drawer.Move(pos.X, pos.Y), map2Rect[drawMap], now.Sub(w.firstRenderTime))
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
	currentMapPos := map2Pos[w.currentMap]
	for i := len(spritePairs) - 1; i >= 0; i-- {
		err = spritePairs[i].E2().Draw(w.ctx, drawer.Move(currentMapPos.X, currentMapPos.Y))
		if err != nil {
			return err
		}
	}
	// 前景
	for drawMap, pos := range map2Pos {
		err = drawMap.DrawForeground(drawer.Move(pos.X, pos.Y), map2Rect[drawMap], now.Sub(w.firstRenderTime))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) MovePixelPosTo(x, y int) {
	w.pixPos = [2]int{x, y}
}

func (w *World) loadMap(id string) (*Map, error) {
	if m := w.mapCache[id]; m != nil {
		return m, nil
	}
	targetMap, err := NewMap(w.ctx, w.tileCache, id)
	if err != nil {
		return nil, err
	}
	w.mapCache[id] = targetMap
	return targetMap, nil
}

func (w *World) MoveTo(id string) error {
	if w.currentMap != nil && w.currentMap.id == id {
		return nil
	}
	targetMap, err := w.loadMap(id)
	if err != nil {
		return err
	}
	w.currentMap = targetMap
	w.nameMoveCounter = 0
	return nil
}

func (w *World) GetActualPosition(x, y int) (*Map, int, int, bool) {
	curMap := w.currentMap
	for {
		adjacentMaps := curMap.AdjacentMaps()
		if y < 0 {
			upMapID, ok := adjacentMaps[util.DirectionEnum.Up]
			if !ok || !stlmaps.ContainKey(w.mapCache, upMapID) {
				return nil, 0, 0, false
			}
			upMap := w.mapCache[upMapID]
			x, y = x, y+upMap.define.Height
			curMap = upMap
			continue
		} else if y >= curMap.define.Height {
			downMapID, ok := adjacentMaps[util.DirectionEnum.Down]
			if !ok || !stlmaps.ContainKey(w.mapCache, downMapID) {
				return nil, 0, 0, false
			}
			downMap := w.mapCache[downMapID]
			x, y = x, y-curMap.define.Height
			curMap = downMap
			continue
		} else if x < 0 {
			leftMapID, ok := adjacentMaps[util.DirectionEnum.Left]
			if !ok || !stlmaps.ContainKey(w.mapCache, leftMapID) {
				return nil, 0, 0, false
			}
			leftMap := w.mapCache[leftMapID]
			x, y = x+leftMap.define.Width, y
			curMap = leftMap
			continue
		} else if x >= curMap.define.Width {
			rightMapID, ok := adjacentMaps[util.DirectionEnum.Right]
			if !ok || !stlmaps.ContainKey(w.mapCache, rightMapID) {
				return nil, 0, 0, false
			}
			rightMap := w.mapCache[rightMapID]
			x, y = x-curMap.define.Width, y
			curMap = rightMap
			continue
		}
		return curMap, x, y, true
	}
}

func (w *World) CurrentMap() *Map {
	return w.currentMap
}

func (w *World) CheckCollision(d util.Direction, x, y int) bool {
	if [2]int{x, y} == w.selfPos {
		return true
	}
	targetMap, x, y, ok := w.GetActualPosition(x, y)
	if !ok {
		return true
	}
	return targetMap.CheckCollision(d, x, y)
}

// DrawMapName 绘制地图名
func (w *World) DrawMapName(drawer draw.OptionDrawer) error {
	height := w.ctx.Config().ScreenHeight / 7
	if w.nameMoveCounter < 0 || w.nameMoveCounter >= height*4 {
		return nil
	}

	img, ok := w.getMapNameDisplayImage()
	if !ok {
		return nil
	}

	if w.nameMoveCounter < height {
		drawer = drawer.Move(10, w.nameMoveCounter-height)
	} else if w.nameMoveCounter < height*3 {
		drawer = drawer.Move(10, 0)
	} else {
		drawer = drawer.Move(10, -w.nameMoveCounter%height)
	}
	draw.PrepareDrawImage(drawer, img).Draw()
	w.nameMoveCounter += w.nameMoveSpeed
	return nil
}

func (w *World) getMapNameDisplayImage() (imgutil.Image, bool) {
	mapName := w.currentMap.Name()
	if mapName == "" {
		return nil, false
	}

	width, height := w.ctx.Config().ScreenWidth/3, w.ctx.Config().ScreenHeight/7
	img := imgutil.NewImage(width, height)
	draw.PrepareDrawRect(img, width, height, util.NewNRGBColor(248, 248, 255)).Draw()
	draw.PrepareDrawRect(img, width, height, nil).SetBorderWidth(12).SetBorderColor(util.NewNRGBColor(176, 196, 222)).Draw()
	draw.PrepareDrawRect(img, width, height, nil).SetBorderWidth(8).SetBorderColor(util.NewNRGBColor(119, 136, 153)).Draw()
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 32).UnsafeInternal(), mapName)
	draw.PrepareDrawText(img, mapName, util.GetFont(util.FontTypeEnum.Normal, 32), color.Black).Move((width+10)/2-(bounds.Max.X.Floor()-bounds.Min.X.Floor())/2, (height-6)/2-(bounds.Max.Y.Floor()-bounds.Min.Y.Floor())/2).Draw()
	return img, true
}
