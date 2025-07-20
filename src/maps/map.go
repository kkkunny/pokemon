package maps

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/kkkunny/pokemon/src/util"
)

type Map struct {
	define *tiled.Map
	render *render.Renderer

	pos [2]int
}

func NewMap() (*Map, error) {
	mapTMX, err := tiled.LoadFile("map/maps/pallet_town.tmx")
	if err != nil {
		return nil, err
	}
	renderer, err := render.NewRenderer(mapTMX)
	if err != nil {
		return nil, err
	}
	return &Map{
		define: mapTMX,
		render: renderer,
	}, nil
}

func (m *Map) Draw(screen *ebiten.Image, sprites []util.Drawer) error {
	// 找到对象层级
	var objectLayerID uint32
	for _, layer := range m.define.Layers {
		objectLayerID = max(objectLayerID, layer.ID)
	}
	if len(m.define.ObjectGroups) > 0 {
		objectLayerID = m.define.ObjectGroups[0].ID
	}
	// 绘制背景
	m.render.Clear()
	for i, layer := range m.define.Layers {
		if layer.ID > objectLayerID {
			continue
		}
		err := m.render.RenderLayer(i)
		if err != nil {
			return err
		}
	}
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(float64(m.pos[0]), float64(m.pos[1]))
	screen.DrawImage(ebiten.NewImageFromImage(m.render.Result), ops)
	// 绘制对象
	for _, s := range sprites {
		err := s.Draw(screen)
		if err != nil {
			return err
		}
	}
	// 绘制前景
	m.render.Clear()
	for i, layer := range m.define.Layers {
		if layer.ID <= objectLayerID {
			continue
		}
		err := m.render.RenderLayer(i)
		if err != nil {
			return err
		}
	}
	screen.DrawImage(ebiten.NewImageFromImage(m.render.Result), ops)
	return nil
}

func (m *Map) TilePixelSize() int {
	return m.define.TileWidth
}

func (m *Map) PixelSize() (w int, h int) {
	return m.define.Width * m.define.TileWidth, m.define.Height * m.define.TileHeight
}

func (m *Map) Size() (w int, h int) {
	return m.define.Width, m.define.Height
}

func (m *Map) CheckCollision(x, y int) bool {
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

func (m *Map) MoveTo(x, y int) {
	m.pos = [2]int{x, y}
}
