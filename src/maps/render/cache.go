package render

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type TileCache struct {
	tilesets map[string]*ebiten.Image
	tiles    map[string]map[int]*ebiten.Image
}

func NewTileCache() *TileCache {
	return &TileCache{
		tilesets: make(map[string]*ebiten.Image),
		tiles:    make(map[string]map[int]*ebiten.Image),
	}
}

func (c *TileCache) AddTilesetImage(path string, img *ebiten.Image) {
	c.tilesets[path] = img
}

func (c *TileCache) GetTilesetImage(path string) *ebiten.Image {
	return c.tilesets[path]
}

func (c *TileCache) AddTileImage(tilesetPath string, index int, img *ebiten.Image) {
	if c.tiles[tilesetPath] == nil {
		c.tiles[tilesetPath] = make(map[int]*ebiten.Image)
	}
	c.tiles[tilesetPath][index] = img
}

func (c *TileCache) GetTileImage(tilesetPath string, index int) *ebiten.Image {
	return c.tiles[tilesetPath][index]
}
