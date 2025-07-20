package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"

	"github.com/kkkunny/pokemon/src/maps"
)

type UpdateInfo struct {
	Person *PersonUpdateInfo
}

type PersonUpdateInfo struct {
	Map *maps.Map
}

type Sprite interface {
	OnAction(action input.Action)
	Update(info *UpdateInfo) error
	Draw(screen *ebiten.Image)
}
