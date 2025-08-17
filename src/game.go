package src

import (
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system/battle"
	"github.com/kkkunny/pokemon/src/system/sub_system"
	worldsubsystem "github.com/kkkunny/pokemon/src/system/world/sub_system"
	"github.com/kkkunny/pokemon/src/util/draw"
)

var subSystems []sub_system.SubSystem

func Init() error {
	// 世界
	ws, err := worldsubsystem.NewWorldSystem()
	if err != nil {
		return err
	}

	// 战斗
	bs, err := battle.NewBattleSystem()
	if err != nil {
		return err
	}

	subSystems = []sub_system.SubSystem{
		sub_system.NewEmptySubSystem(),
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
	for _, p := range sub_system.PlayingVoicePlayers {
		err = p.Player.Play()
		if err != nil {
			return err
		}
	}
	stopPlayer := stlslices.DiffTo(sub_system.PrevPlayingVoicePlayers, sub_system.PlayingVoicePlayers)
	sub_system.PrevPlayingVoicePlayers = sub_system.PlayingVoicePlayers
	sub_system.PlayingVoicePlayers = nil
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
	subSystems = stlslices.Filter(subSystems, func(_ int, ss sub_system.SubSystem) bool {
		return !ss.Drop()
	})
}
