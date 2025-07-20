package sprite

import (
	input "github.com/quasilyte/ebitengine-input"

	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/util"
)

type UpdateInfo struct {
	Person *PersonUpdateInfo
}

type PersonUpdateInfo struct {
	Map *maps.Map
}

type Sprite interface {
	util.Drawer
	OnAction(action input.Action, info *UpdateInfo)
	Update() error
}
