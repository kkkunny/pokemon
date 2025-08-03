package person

import "github.com/kkkunny/pokemon/src/world"

type UpdateInfo struct {
	World *world.World
}

func (i *UpdateInfo) UpdateInfo() {}
