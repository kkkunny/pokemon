package input

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/tnnmigga/enum"
)

type Action = input.Action

var ActionEnum = enum.New[struct {
	MoveUp    Action
	MoveDown  Action
	MoveLeft  Action
	MoveRight Action
}]()

type System struct {
	inputSystem   input.System
	actionHandler *input.Handler
}

func NewSystem() *System {
	s := &System{}
	s.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		ActionEnum.MoveUp:    {input.KeyGamepadUp, input.KeyW},
		ActionEnum.MoveDown:  {input.KeyGamepadDown, input.KeyS},
		ActionEnum.MoveLeft:  {input.KeyGamepadLeft, input.KeyA},
		ActionEnum.MoveRight: {input.KeyGamepadRight, input.KeyD},
	}
	s.actionHandler = s.inputSystem.NewHandler(0, keymap)
	return s
}

func (s *System) Action() (*Action, error) {
	var action *Action
	switch {
	case s.actionHandler.ActionIsPressed(ActionEnum.MoveUp):
		action = &ActionEnum.MoveUp
	case s.actionHandler.ActionIsPressed(ActionEnum.MoveDown):
		action = &ActionEnum.MoveDown
	case s.actionHandler.ActionIsPressed(ActionEnum.MoveLeft):
		action = &ActionEnum.MoveLeft
	case s.actionHandler.ActionIsPressed(ActionEnum.MoveRight):
		action = &ActionEnum.MoveRight
	}
	return action, nil
}
