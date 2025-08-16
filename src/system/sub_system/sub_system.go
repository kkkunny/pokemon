package sub_system

import (
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/util/draw"
)

type SubSystemType uint8

var SubSystemTypeEnum = enum.New[struct {
	Battle   SubSystemType
	Dialogue SubSystemType
	World    SubSystemType
}]()

type SubSystem interface {
	Type() SubSystemType
	OnAction(system SubSystemQueue, action input.KeyInputAction) error
	OnUpdate(system SubSystemQueue) error
	OnDraw(drawer draw.OptionDrawer) error
}
