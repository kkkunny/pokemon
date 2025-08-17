package system

import (
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/util/draw"
)

type System interface {
	// 只会在action被触发时调用，且只会触发最后一个子系统的OnAction
	OnAction(system SystemManager, action input.KeyInputAction) error
	OnUpdate(system SystemManager) error
	OnDraw(system SystemManager, drawer draw.OptionDrawer) error
	Drop() bool
}

type emptySystem struct{}

func NewEmptySystem() System {
	return emptySystem{}
}

func (emptySystem) OnAction(_ SystemManager, _ input.KeyInputAction) error {
	return nil
}

func (emptySystem) OnUpdate(system SystemManager) error {
	return nil
}

func (emptySystem) OnDraw(system SystemManager, drawer draw.OptionDrawer) error {
	return nil
}

func (emptySystem) Drop() bool {
	return false
}
