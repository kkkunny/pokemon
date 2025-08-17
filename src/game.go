package src

import (
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system"
	"github.com/kkkunny/pokemon/src/system/battle"
	worldsystem "github.com/kkkunny/pokemon/src/system/world/system"
	"github.com/kkkunny/pokemon/src/util/draw"
)

var subSystems []system.System

func Init() error {
	// 世界
	ws, err := worldsystem.NewWorldSystem()
	if err != nil {
		return err
	}

	// 战斗
	bs, err := battle.NewBattleSystem()
	if err != nil {
		return err
	}

	subSystems = []system.System{
		system.NewEmptySystem(),
		ws,
		bs,
	}
	return nil
}

func OnAction(action input.KeyInputAction) error {
	dropSubSystem()

	return stlslices.Last(subSystems).OnAction(newCursor(), action)
}

func OnUpdate() error {
	dropSubSystem()

	err := stlslices.Last(subSystems).OnUpdate(newCursor())
	if err != nil {
		return err
	}

	// 声音
	for _, p := range system.PlayingVoicePlayers {
		err = p.Player.Play()
		if err != nil {
			return err
		}
	}
	stopPlayer := stlslices.DiffTo(system.PrevPlayingVoicePlayers, system.PlayingVoicePlayers)
	system.PrevPlayingVoicePlayers = system.PlayingVoicePlayers
	system.PlayingVoicePlayers = nil
	for _, p := range stopPlayer {
		p.Pause()
	}

	return nil
}

func OnDraw(drawer draw.OptionDrawer) error {
	dropSubSystem()

	return stlslices.Last(subSystems).OnDraw(newCursor(), drawer)
}

func dropSubSystem() {
	subSystems = stlslices.Filter(subSystems, func(_ int, ss system.System) bool {
		return !ss.Drop()
	})
}
