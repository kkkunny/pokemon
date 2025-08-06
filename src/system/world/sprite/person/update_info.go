package person

import (
	"github.com/kkkunny/pokemon/src/system/world"
)

type UpdateInfo struct {
	World *world.World
}

func (i *UpdateInfo) UpdateInfo() {}
