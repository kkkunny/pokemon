package render

import (
	"image"
	"time"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type Renderer struct {
	m     *tiled.Map
	cache *TileCache
	dur   time.Duration
}

func NewRenderer(m *tiled.Map, cache *TileCache, dur time.Duration) *Renderer {
	return &Renderer{
		m:     m,
		cache: cache,
		dur:   dur,
	}
}

func (r *Renderer) getRenderInfo() (xs int, xe int, xi int, ys int, ye int, yi int) {
	renderOrder := r.m.RenderOrder
	if renderOrder == "" {
		renderOrder = "right-down"
	}
	switch renderOrder {
	case "right-down":
		xs = 0
		xe = r.m.Width
		xi = 1
		ys = 0
		ye = r.m.Height
		yi = 1
	default:
		panic(render.ErrUnsupportedRenderOrder)
	}
	return
}

func (r *Renderer) foreachLayerTile(layer *tiled.Layer, fn func(x, y int, layerTile *tiled.LayerTile) bool) {
	xs, xe, xi, ys, ye, yi := r.getRenderInfo()
	i := 0
	for y := ys; y*yi < ye; y = y + yi {
		for x := xs; x*xi < xe; x = x + xi {
			layerTile := layer.Tiles[i]
			i++
			if !fn(x, y, layerTile) {
				return
			}
		}
	}
}

func (r *Renderer) getTileImage(tile *tiled.LayerTile) (*imgutil.Image, error) {
	tilesetPath := tile.Tileset.GetFileFullPath(tile.Tileset.Image.Source)
	if img := r.cache.GetTileImage(tilesetPath, int(tile.ID)); img != nil {
		return img, nil
	}

	tilesetImg := r.cache.GetTilesetImage(tilesetPath)
	if tilesetImg == nil {
		var err error
		tilesetImg, err = imgutil.NewImageFromFile(tilesetPath)
		if err != nil {
			return nil, err
		}
		r.cache.AddTilesetImage(tilesetPath, tilesetImg)
	}

	img := tilesetImg.SubImageByRect(tile.Tileset.GetTileRect(tile.ID))
	r.cache.AddTileImage(tilesetPath, int(tile.ID), img)

	return img, nil
}

func (r *Renderer) renderLayer(drawer draw.Drawer, layer *tiled.Layer, rect image.Rectangle) error {
	// TODO: 可优化，foreach函数里直接不遍历rect外的坐标
	var retErr error
	r.foreachLayerTile(layer, func(x, y int, layerTile *tiled.LayerTile) bool {
		if layerTile == nil || layerTile.IsNil() || x < rect.Min.X || x > rect.Max.X || y < rect.Min.Y || y > rect.Max.Y {
			return true
		}
		err := func() error {
			// 动画
			if layerTile.Tileset != nil {
				tileDef, err := layerTile.Tileset.GetTilesetTile(layerTile.ID)
				if err == nil && tileDef != nil && len(tileDef.Animation) > 0 {
					index := int(r.dur/(time.Millisecond*time.Duration(tileDef.Animation[0].Duration))) % len(tileDef.Animation)
					thisFrame := tileDef.Animation[index]
					newLayerTile := *layerTile
					newLayerTile.ID = thisFrame.TileID
					layerTile = &newLayerTile
				}
			}

			img, err := r.getTileImage(layerTile)
			if err != nil {
				return err
			}

			if layer.Opacity < 1 {
				// TODO
				panic("unexpected opacity")
			} else {
				err = drawer.Move(float64(x*r.m.TileWidth), float64(y*r.m.TileHeight)).DrawImage(img)
				if err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil {
			retErr = err
		}
		return err == nil
	})
	return retErr
}

func (r *Renderer) RenderLayer(drawer draw.Drawer, id int) error {
	return r.RenderRectLayer(drawer, id, image.Rect(0, 0, r.m.Width, r.m.Height))
}

func (r *Renderer) RenderRectLayer(drawer draw.Drawer, id int, rect image.Rectangle) error {
	if id >= len(r.m.Layers) {
		return render.ErrOutOfBounds
	}
	return r.renderLayer(drawer, r.m.Layers[id], rect)
}
