package sprite

import "github.com/kkkunny/pokemon/src/consts"

type MovableSprite interface {
	Sprite
	Direction() consts.Direction
	Turn(d consts.Direction) bool
	SetMovable(movable bool)
	Movable() bool
	Moving() bool
}
