package util

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/context"
)

type Drawer interface {
	Draw(ctx context.Context, screen *ebiten.Image, options ebiten.DrawImageOptions) error
}
