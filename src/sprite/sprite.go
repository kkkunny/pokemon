package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"

	"github.com/kkkunny/pokemon/src/maps"
)

type DrawInfo struct {
	Person *PersonDrawInfo
}

type PersonDrawInfo struct {
	Map *maps.Map
}

type Sprite interface {
	OnAction(action input.Action)
	Update() error
	Draw(screen *ebiten.Image, info *DrawInfo)
}
