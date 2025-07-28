package input

import (
	stlval "github.com/kkkunny/stl/value"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/tnnmigga/enum"
)

type KeyInputAction input.Action

func (a KeyInputAction) action() input.Action {
	return input.Action(a)
}

func (a KeyInputAction) Pressed() KeyInputAction {
	return a + 100
}

func (a KeyInputAction) Released() KeyInputAction {
	return a + 200
}

var KeyInputActionEnum = enum.New[struct {
	MoveUp    KeyInputAction
	MoveDown  KeyInputAction
	MoveLeft  KeyInputAction
	MoveRight KeyInputAction

	A KeyInputAction
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
		KeyInputActionEnum.MoveUp.action():    {input.KeyGamepadUp, input.KeyW},
		KeyInputActionEnum.MoveDown.action():  {input.KeyGamepadDown, input.KeyS},
		KeyInputActionEnum.MoveLeft.action():  {input.KeyGamepadLeft, input.KeyA},
		KeyInputActionEnum.MoveRight.action(): {input.KeyGamepadRight, input.KeyD},
		KeyInputActionEnum.A.action():         {input.KeyGamepadA, input.KeyJ},
	}
	s.actionHandler = s.inputSystem.NewHandler(0, keymap)
	return s
}

func (s *System) KeyInputAction() (*KeyInputAction, error) {
	for _, a := range enum.Values[KeyInputAction](KeyInputActionEnum) {
		if s.actionHandler.ActionIsJustPressed(a.action()) {
			return stlval.Ptr(a.Pressed()), nil
		} else if s.actionHandler.ActionIsPressed(a.action()) {
			return &a, nil
		}
	}
	return nil, nil
}
