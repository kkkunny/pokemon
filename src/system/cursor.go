package system

import (
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/system/dialogue"
	"github.com/kkkunny/pokemon/src/system/sub_system"
)

type cursor struct {
	sys   *System
	index int
}

func newCursor(sys *System) *cursor {
	return &cursor{
		sys:   sys,
		index: len(sys.subSystems) - 1,
	}
}

func (s *cursor) Pop() {
	if stlslices.Last(s.sys.subSystems).Type() == sub_system.SubSystemTypeEnum.Empty {
		return
	}
	s.sys.subSystems = s.sys.subSystems[:len(s.sys.subSystems)-1]
}

func (s *cursor) Next() sub_system.SubSystem {
	s.index--
	if s.index < 0 || s.index >= len(s.sys.subSystems) {
		return sub_system.NewEmptySubSystem()
	}
	return s.sys.subSystems[s.index]
}

func (s *cursor) DisplayLabel(text string) error {
	ds, err := dialogue.NewDialogueSystem(s.sys.ctx)
	if err != nil {
		return err
	}
	s.sys.subSystems = append(s.sys.subSystems, ds)
	ds.SetLabel(text)
	return nil
}

func (s *cursor) DisplayDialogue(text string) error {
	ds, err := dialogue.NewDialogueSystem(s.sys.ctx)
	if err != nil {
		return err
	}
	s.sys.subSystems = append(s.sys.subSystems, ds)
	ds.SetDialogue(text)
	return nil
}
