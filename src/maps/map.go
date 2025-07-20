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
	mapTMX, err := tiled.LoadFile("maps/hoenn/littleroot-town.tmx")
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

func (m *Map) TileSize() int {
	return m.define.TileWidth
}

func (m *Map) Size() (w int, h int) {
	return m.define.Width, m.define.Height
}

func (m *Map) Rect() (w int, h int) {
	return m.define.Width / m.define.TileWidth, m.define.Height / m.define.TileHeight
}
