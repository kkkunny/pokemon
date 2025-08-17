package system

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/kkkunny/pokemon/src/i18n"
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

func (s *cursor) DisplayText(boxStyle sub_system.BoxStyle, text string, fontColor color.Color) error {
	ss, err := dialogue.NewText(boxStyle, text, fontColor)
	if err != nil {
		return err
	}
	s.sys.subSystems = append(s.sys.subSystems, ss)
	return nil
}

func (s *cursor) StartOneBattle(site string) error {
	ss := sub_system.ExtractSubSystem[*battle.BattleSystem](s)
	err := ss.StartOneBattle(site)
	if err != nil {
		return err
	}

	pokemonName := i18n.Get(fmt.Sprintf("pokemon.%d", ss.GetSelfPokemon().ID))
	text := strings.ReplaceAll(i18n.Get("default_battle_action_question"), "%POKEMON_NAME%", pokemonName)
	return s.DisplayText(sub_system.BoxStyleEnum.Battle, text, color.White)
}
