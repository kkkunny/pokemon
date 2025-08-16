package system

import (
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/system/dialogue"
	"github.com/kkkunny/pokemon/src/system/sub_system"
	worldsubsystem "github.com/kkkunny/pokemon/src/system/world/sub_system"
	"github.com/kkkunny/pokemon/src/util/draw"
)

type System struct {
	ctx        context.Context
	subSystems []sub_system.SubSystem
}

func NewSystem(ctx context.Context) (*System, error) {
	// 世界
	ws, err := worldsubsystem.NewWorldSystem(ctx)
	if err != nil {
		return nil, err
	}

	// // 战斗
	// bs, err := battle.NewBattleSystem(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return &System{
		ctx:        ctx,
		subSystems: []sub_system.SubSystem{ws},
	}, nil
}

func (s *System) OnAction(action input.KeyInputAction) error {
	if len(s.subSystems) == 0 {
		return nil
	}
	return stlslices.Last(s.subSystems).OnAction(s, action)
}

func (s *System) OnUpdate() error {
	for _, subSystem := range s.subSystems {
		err := subSystem.OnUpdate(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *System) OnDraw(drawer draw.OptionDrawer) error {
	for _, subSystem := range s.subSystems {
		err := subSystem.OnDraw(drawer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *System) Pop() {
	if len(s.subSystems) == 0 {
		return
	}
	s.subSystems = s.subSystems[:len(s.subSystems)-1]
}

func (s *System) DisplayLabel(text string) error {
	// 对话
	ds, err := dialogue.NewDialogueSystem(s.ctx)
	if err != nil {
		return err
	}
	s.subSystems = append(s.subSystems, ds)
	ds.SetLabel(text)
	return nil
}

func (s *System) DisplayDialogue(text string) error {
	// 对话
	ds, err := dialogue.NewDialogueSystem(s.ctx)
	if err != nil {
		return err
	}
	s.subSystems = append(s.subSystems, ds)
	ds.SetDialogue(text)
	return nil
}
