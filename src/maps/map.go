package maps

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"

	"github.com/kkkunny/pokemon/src/maps/render"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/util/image"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
)

type Map struct {
	name         string
	define       *tiled.Map
	tileCache    *render.TileCache
	adjacentMaps map[consts.Direction]*Map
	songFilepath string
	sprites      []sprite.Sprite
}

func NewMap(cfg *config.Config, tileCache *render.TileCache, name string) (*Map, error) {
	return newMapWithAdjacent(cfg, tileCache, name, make(map[string]*Map))
}

func newMapWithAdjacent(cfg *config.Config, tileCache *render.TileCache, name string, existMap map[string]*Map) (*Map, error) {
	// 缓存
	curMap := existMap[name]
	if curMap != nil {
		return curMap, nil
	}

	// 地图
	mapTMX, err := tiled.LoadFile(fmt.Sprintf("map/maps/%s.tmx", name))
	if err != nil {
		return nil, err
	}
	if mapTMX.TileWidth != mapTMX.TileHeight || mapTMX.TileWidth != cfg.TileSize {
		return nil, errors.New("map tile is not valid")
	}

	curMap = &Map{
		name:      name,
		define:    mapTMX,
		tileCache: tileCache,
	}
	existMap[name] = curMap

	// 音频
	songFileName := mapTMX.Properties.GetString("song")
	if songFileName != "" {
		curMap.songFilepath = filepath.Join(config.VoicePath, "map", songFileName)
	}

	// 精灵
	for _, objectGroup := range mapTMX.ObjectGroups {
		for _, object := range objectGroup.Objects {
			x, y := int(object.X+object.Width/2)/cfg.TileSize, int(object.Y+object.Height/2)/cfg.TileSize
			spriteObj, err := sprite.NewSprite(object)
			if err != nil {
				return nil, err
			}
			spriteObj.SetPosition(x, y)
			curMap.sprites = append(curMap.sprites, spriteObj)
		}
	}

	adjacentMaps := curMap.AdjacentMaps()
	curMap.adjacentMaps = make(map[consts.Direction]*Map, len(adjacentMaps))
	for direction, mapName := range adjacentMaps {
		directionMap, err := newMapWithAdjacent(cfg, tileCache, mapName, existMap)
		if err != nil {
			return nil, err
		}
		curMap.adjacentMaps[direction] = directionMap
	}
	return curMap, nil
}

func (m *Map) DrawBackground(screen *image.Image, options ebiten.DrawImageOptions, dur time.Duration) error {
	renderer := render.NewRenderer(m.define, m.tileCache, dur)

	// 找到对象层级
	var objectLayerName string
	for _, layer := range m.define.Layers {
		objectLayerName = max(objectLayerName, layer.Name)
	}
	if len(m.define.ObjectGroups) > 0 {
		objectLayerName = m.define.ObjectGroups[0].Name
	}

	for i, layer := range m.define.Layers {
		if layer.Name > objectLayerName {
			continue
		}
		err := renderer.RenderLayer(screen, i, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) DrawForeground(screen *image.Image, options ebiten.DrawImageOptions, dur time.Duration) error {
	renderer := render.NewRenderer(m.define, m.tileCache, dur)

	// 找到对象层级
	var objectLayerName string
	for _, layer := range m.define.Layers {
		objectLayerName = max(objectLayerName, layer.Name)
	}
	if len(m.define.ObjectGroups) > 0 {
		objectLayerName = m.define.ObjectGroups[0].Name
	}

	for i, layer := range m.define.Layers {
		if layer.Name <= objectLayerName {
			continue
		}
		err := renderer.RenderLayer(screen, i, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) PixelSize() (w int, h int) {
	return m.define.Width * m.define.TileWidth, m.define.Height * m.define.TileHeight
}

func (m *Map) Size() (w int, h int) {
	return m.define.Width, m.define.Height
}

func (m *Map) NameLocKey() string {
	return m.define.Properties.GetString("name")
}

func (m *Map) SongFilepath() (string, bool) {
	return m.songFilepath, m.songFilepath != ""
}

func (m *Map) CheckCollision(x, y int) bool {
	for _, s := range m.sprites {
		sx, sy := s.Position()
		if movable, ok := s.(sprite.MovableSprite); ok {
			sx, sy = movable.NextStepPosition()
		}
		if sx == x && sy == y {
			return true
		}
	}
	for _, layer := range m.define.Layers {
		index := y*m.define.Width + x
		tile := layer.Tiles[index]
		if tile.Tileset == nil {
			continue
		}
		tileDef, err := tile.Tileset.GetTilesetTile(tile.ID)
		if err != nil {
			continue
		}
		return tileDef.Properties.GetBool("collision")
	}
	return false
}

func (m *Map) AdjacentMaps() map[consts.Direction]string {
	directions := map[consts.Direction]string{
		consts.DirectionEnum.Up:    "up",
		consts.DirectionEnum.Down:  "down",
		consts.DirectionEnum.Left:  "left",
		consts.DirectionEnum.Right: "right",
	}
	maps := make(map[consts.Direction]string, len(directions))
	for direction, attr := range directions {
		mapName := m.define.Properties.GetString(attr)
		if mapName == "" {
			continue
		}
		maps[direction] = mapName
	}
	return maps
}

func (m *Map) Sprites() []sprite.Sprite {
	return m.sprites
}

func (m *Map) GetActualPosition(x, y int) (*Map, int, int, bool) {
	if y < 0 {
		upMap := m.adjacentMaps[consts.DirectionEnum.Up]
		if upMap == nil {
			return nil, 0, 0, false
		}
		return upMap.GetActualPosition(x, y+upMap.define.Height)
	} else if y >= m.define.Height {
		downMap := m.adjacentMaps[consts.DirectionEnum.Down]
		if downMap == nil {
			return nil, 0, 0, false
		}
		return downMap.GetActualPosition(x, y-m.define.Height)
	} else if x < 0 {
		leftMap := m.adjacentMaps[consts.DirectionEnum.Left]
		if leftMap == nil {
			return nil, 0, 0, false
		}
		return leftMap.GetActualPosition(x+leftMap.define.Width, y)
	} else if x >= m.define.Width {
		rightMap := m.adjacentMaps[consts.DirectionEnum.Right]
		if rightMap == nil {
			return nil, 0, 0, false
		}
		return rightMap.GetActualPosition(x-m.define.Width, y)
	}
	return m, x, y, true
}

func (m *Map) GetSpriteByPosition(x, y int) (sprite.Sprite, bool) {
	for _, s := range m.sprites {
		sx, sy := s.Position()
		if sx == x && sy == y {
			return s, true
		}
	}
	return nil, false
}
