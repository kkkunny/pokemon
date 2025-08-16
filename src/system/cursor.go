package system

import (
	"github.com/kkkunny/pokemon/src/system/battle"
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

func (s *cursor) SubSystems() []sub_system.SubSystem {
	return s.sys.subSystems
}

func (s *cursor) Next() sub_system.SubSystem {
	s.index--
	if s.index < 0 || s.index >= len(s.sys.subSystems) {
		return sub_system.NewEmptySubSystem()
	}
	return s.sys.subSystems[s.index]
}

func (s *cursor) DisplayLabel(text string) error {
	ss := sub_system.ExtractSubSystem[*dialogue.DialogueSystem](s)
	ss.SetDisplay(true)
	ss.SetLabel(text)
	return nil
}

func (s *cursor) DisplayDialogue(text string) error {
	ss := sub_system.ExtractSubSystem[*dialogue.DialogueSystem](s)
	ss.SetDisplay(true)
	ss.SetDialogue(text)
	return nil
}

func (s *cursor) StartOneBattle(site string) error {
	ss := sub_system.ExtractSubSystem[*battle.BattleSystem](s)
	return ss.StartOneBattle(site)
}
