package render

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
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

func (r *Renderer) getTileImage(tile *tiled.LayerTile) (*ebiten.Image, error) {
	tilesetPath := tile.Tileset.GetFileFullPath(tile.Tileset.Image.Source)
	if img := r.cache.GetTileImage(tilesetPath, int(tile.ID)); img != nil {
		return img, nil
	}

	tilesetImg := r.cache.GetTilesetImage(tilesetPath)
	if tilesetImg == nil {
		var err error
		tilesetImg, _, err = ebitenutil.NewImageFromFile(tilesetPath)
		if err != nil {
			return nil, err
		}
		r.cache.AddTilesetImage(tilesetPath, tilesetImg)
	}

	img := tilesetImg.SubImage(tile.Tileset.GetTileRect(tile.ID)).(*ebiten.Image)
	r.cache.AddTileImage(tilesetPath, int(tile.ID), img)

	return img, nil
}

func (r *Renderer) renderLayer(target *ebiten.Image, layer *tiled.Layer, options *ebiten.DrawImageOptions) error {
	var retErr error
	r.foreachLayerTile(layer, func(x, y int, layerTile *tiled.LayerTile) bool {
		err := func() error {
			if layerTile == nil || layerTile.IsNil() {
				return nil
			}

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
				ops := *options
				ops.GeoM.Translate(float64(x*r.m.TileWidth), float64(y*r.m.TileHeight))
				target.DrawImage(img, &ops)
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

func (r *Renderer) RenderLayer(target *ebiten.Image, id int, options *ebiten.DrawImageOptions) error {
	if id >= len(r.m.Layers) {
		return render.ErrOutOfBounds
	}
	return r.renderLayer(target, r.m.Layers[id], options)
}
