package util

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/config"
)

type Drawer interface {
	Draw(cfg *config.Config, screen *ebiten.Image, options ebiten.DrawImageOptions) error
}
