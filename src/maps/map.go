package maps

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type Map struct {
	define *tiled.Map
	render *render.Renderer
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
	err = renderer.RenderVisibleLayers()
	if err != nil {
		return nil, err
	}
	return &Map{
		define: mapTMX,
		render: renderer,
	}, nil
}

func (m *Map) Image() *ebiten.Image {
	return ebiten.NewImageFromImage(m.render.Result)
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
