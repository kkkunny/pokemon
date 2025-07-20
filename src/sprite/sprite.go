package sprite

import (
	input "github.com/quasilyte/ebitengine-input"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/maps"
	"github.com/kkkunny/pokemon/src/util"
)

type UpdateInfo struct {
	Person *PersonUpdateInfo
}

type PersonUpdateInfo struct {
	World *maps.World
}

type Sprite interface {
	util.Drawer
	OnAction(cfg *config.Config, action input.Action, info *UpdateInfo)
	Update(cfg *config.Config, info *UpdateInfo) error
}
