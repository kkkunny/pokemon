package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type Sprite interface {
	OnAction(action input.Action) error
	Update() error
	Image() (*ebiten.Image, error)
	Position() (x, y int, display bool)
}
