package sprite

import (
	"github.com/kkkunny/pokemon/src/util"
)

type MovableSprite interface {
	Sprite
	Direction() util.Direction
	Turn(d util.Direction) bool
	SetMovable(movable bool)
	Movable() bool
	Moving() bool
}
