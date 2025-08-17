package src

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/system"
	"github.com/kkkunny/pokemon/src/system/battle"
	"github.com/kkkunny/pokemon/src/system/dialogue"
)

type cursor struct {
	index int
}

func newCursor() *cursor {
	return &cursor{
		index: len(subSystems) - 1,
	}
}

func (s *cursor) SubSystems() []system.System {
	return subSystems
}

func (s *cursor) Next() system.System {
	s.index--
	if s.index < 0 || s.index >= len(subSystems) {
		return system.NewEmptySystem()
	}
	return subSystems[s.index]
}

func (s *cursor) DisplayText(boxStyle system.BoxStyle, text string, fontColor color.Color) error {
	ss, err := dialogue.NewText(boxStyle, text, fontColor)
	if err != nil {
		return err
	}
	subSystems = append(subSystems, ss)
	return nil
}

func (s *cursor) StartOneBattle(site string) error {
	ss := system.ExtractSystem[*battle.BattleSystem](s)
	err := ss.StartOneBattle(site)
	if err != nil {
		return err
	}

	pokemonName := i18n.Get(fmt.Sprintf("pokemon.%d", ss.GetSelfPokemon().ID))
	text := strings.ReplaceAll(i18n.Get("default_battle_action_question"), "%POKEMON_NAME%", pokemonName)
	return s.DisplayText(system.BoxStyleEnum.Battle, text, color.White)
}
