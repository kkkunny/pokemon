package sub_system

import (
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/util/draw"
)

type SubSystemType uint8

var SubSystemTypeEnum = enum.New[struct {
	Empty    SubSystemType
	Battle   SubSystemType
	Dialogue SubSystemType
	World    SubSystemType
}]()

type SubSystem interface {
	Type() SubSystemType
	// 只会在action被触发时调用，且只会触发最后一个子系统的OnAction
	OnAction(system SubSystemManager, action input.KeyInputAction) error
	OnUpdate(system SubSystemManager) error
	OnDraw(system SubSystemManager, drawer draw.OptionDrawer) error
}

type emptySubSystem struct{}

func NewEmptySubSystem() SubSystem {
	return emptySubSystem{}
}

func (emptySubSystem) Type() SubSystemType {
	return SubSystemTypeEnum.Empty
}

func (emptySubSystem) OnAction(_ SubSystemManager, _ input.KeyInputAction) error {
	return nil
}

func (emptySubSystem) OnUpdate(system SubSystemManager) error {
	return nil
}

func (emptySubSystem) OnDraw(system SubSystemManager, drawer draw.OptionDrawer) error {
	return nil
}
