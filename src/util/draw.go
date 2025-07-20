package util

import "github.com/hajimehoshi/ebiten/v2"

type Drawer interface {
	Draw(screen *ebiten.Image) error
}
