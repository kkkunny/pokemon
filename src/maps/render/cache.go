package render

import (
	"github.com/kkkunny/pokemon/src/util/image"
)

type TileCache struct {
	tilesets map[string]*image.Image
	tiles    map[string]map[int]*image.Image
}

func NewTileCache() *TileCache {
	return &TileCache{
		tilesets: make(map[string]*image.Image),
		tiles:    make(map[string]map[int]*image.Image),
	}
}

func (c *TileCache) AddTilesetImage(path string, img *image.Image) {
	c.tilesets[path] = img
}

func (c *TileCache) GetTilesetImage(path string) *image.Image {
	return c.tilesets[path]
}

func (c *TileCache) AddTileImage(tilesetPath string, index int, img *image.Image) {
	if c.tiles[tilesetPath] == nil {
		c.tiles[tilesetPath] = make(map[int]*image.Image)
	}
	c.tiles[tilesetPath][index] = img
}

func (c *TileCache) GetTileImage(tilesetPath string, index int) *image.Image {
	return c.tiles[tilesetPath][index]
}
