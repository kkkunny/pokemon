package maps

import (
	"errors"
	"fmt"
	"image"
	"path/filepath"
	"time"

	stlslices "github.com/kkkunny/stl/container/slices"
	"github.com/lafriks/go-tiled"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/maps/render"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/util/draw"
)

type ObjectLayerType = string

var ObjectLayerTypeEnum = enum.New[struct {
	Sprite ObjectLayerType `enum:"sprite"`
	Split  ObjectLayerType `enum:"split"`
}]()

type Map struct {
	ctx          context.Context
	id           string
	define       *tiled.Map
	tileCache    *render.TileCache
	adjacentMaps map[consts.Direction]*Map
	songFilepath string
	sprites      []sprite.Sprite
}

func NewMap(ctx context.Context, tileCache *render.TileCache, id string) (*Map, error) {
	return newMapWithAdjacent(ctx, tileCache, id, make(map[string]*Map))
}

func newMapWithAdjacent(ctx context.Context, tileCache *render.TileCache, id string, existMap map[string]*Map) (*Map, error) {
	// 缓存
	curMap := existMap[id]
	if curMap != nil {
		return curMap, nil
	}

	// 地图
	mapTMX, err := tiled.LoadFile(fmt.Sprintf("map/maps/%s.tmx", id))
	if err != nil {
		return nil, err
	}
	if mapTMX.TileWidth != mapTMX.TileHeight || mapTMX.TileWidth != ctx.Config().TileSize {
		return nil, errors.New("map tile is not valid")
	}

	curMap = &Map{
		ctx:       ctx,
		id:        id,
		define:    mapTMX,
		tileCache: tileCache,
	}
	existMap[id] = curMap

	// 音频
	if mapTMX.Properties != nil {
		songFileName := mapTMX.Properties.GetString("song")
		if songFileName != "" {
			curMap.songFilepath = filepath.Join(config.VoicePath, "map", songFileName)
		}
	}

	// 精灵
	for _, objectGroup := range mapTMX.ObjectGroups {
		if objectGroup.Class != ObjectLayerTypeEnum.Sprite {
			continue
		}
		for _, object := range objectGroup.Objects {
			x, y := int(object.X)/ctx.Config().TileSize, int(object.Y)/ctx.Config().TileSize
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
		directionMap, err := newMapWithAdjacent(ctx, tileCache, mapName, existMap)
		if err != nil {
			return nil, err
		}
		curMap.adjacentMaps[direction] = directionMap
	}
	return curMap, nil
}

func (m *Map) getSpriteLayerName() string {
	var layerName string
	for _, layer := range m.define.Layers {
		layerName = max(layerName, layer.Name)
	}
	spritesLayers := stlslices.Filter(m.define.ObjectGroups, func(_ int, og *tiled.ObjectGroup) bool {
		return og.Class == ObjectLayerTypeEnum.Sprite
	})
	if len(spritesLayers) > 0 {
		layerName = spritesLayers[0].Name
	}
	return layerName
}

func (m *Map) DrawBackground(drawer draw.Drawer, rect image.Rectangle, dur time.Duration) error {
	renderer := render.NewRenderer(m.define, m.tileCache, dur)

	objectLayerName := m.getSpriteLayerName()
	for i, layer := range m.define.Layers {
		if layer.Name > objectLayerName {
			continue
		}
		err := renderer.RenderRectLayer(drawer, i, rect)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) DrawForeground(drawer draw.Drawer, rect image.Rectangle, dur time.Duration) error {
	renderer := render.NewRenderer(m.define, m.tileCache, dur)

	objectLayerName := m.getSpriteLayerName()
	for i, layer := range m.define.Layers {
		if layer.Name <= objectLayerName {
			continue
		}
		err := renderer.RenderRectLayer(drawer, i, rect)
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

func (m *Map) ID() string {
	return m.id
}

func (m *Map) Name() string {
	return m.ctx.Localisation().Get(m.define.Properties.GetString("name"))
}

func (m *Map) SongFilepath() (string, bool) {
	return m.songFilepath, m.songFilepath != ""
}

func (m *Map) CheckCollision(d consts.Direction, x, y int) bool {
	for _, s := range m.sprites {
		if !s.Collision() {
			continue
		}
		sx, sy := s.CollisionPosition()
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
		if tileDef.Properties.GetBool("collision") {
			return true
		} else if allowDirectionStr := tileDef.Properties.GetString("allow_direction"); allowDirectionStr != "" {
			return d != -consts.ParseDirection(allowDirectionStr)
		}
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

func (m *Map) GetHoles() []*tiled.Object {
	splitLayers := stlslices.Filter(m.define.ObjectGroups, func(_ int, og *tiled.ObjectGroup) bool {
		return og.Class == ObjectLayerTypeEnum.Split
	})
	return stlslices.FlatMap(splitLayers, func(_ int, ob *tiled.ObjectGroup) []*tiled.Object {
		return stlslices.Filter(ob.Objects, func(_ int, o *tiled.Object) bool {
			return o.Type == "hole"
		})
	})
}

func (m *Map) Indoor() bool {
	return m.define.Properties.GetBool("indoor")
}
