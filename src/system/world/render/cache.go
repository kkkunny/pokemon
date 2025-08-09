package render

import (
	"github.com/kkkunny/pokemon/src/util/image"
)

type TileCache struct {
	tilesets map[string]imgutil.Image
	tiles    map[string]map[int]imgutil.Image
}

func NewTileCache() *TileCache {
	return &TileCache{
		tilesets: make(map[string]imgutil.Image),
		tiles:    make(map[string]map[int]imgutil.Image),
	}
}

func (c *TileCache) AddTilesetImage(path string, img imgutil.Image) {
	c.tilesets[path] = img
}

func (c *TileCache) GetTilesetImage(path string) imgutil.Image {
	return c.tilesets[path]
}

func (c *TileCache) AddTileImage(tilesetPath string, index int, img imgutil.Image) {
	if c.tiles[tilesetPath] == nil {
		c.tiles[tilesetPath] = make(map[int]imgutil.Image)
	}
	c.tiles[tilesetPath][index] = img
}

func (c *TileCache) GetTileImage(tilesetPath string, index int) imgutil.Image {
	return c.tiles[tilesetPath][index]
}
