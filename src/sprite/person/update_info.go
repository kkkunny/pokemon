package person

import "github.com/kkkunny/pokemon/src/maps"

type UpdateInfo struct {
	World *maps.World
}

func (i *UpdateInfo) UpdateInfo() {}
