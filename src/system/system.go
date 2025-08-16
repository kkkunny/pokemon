package system

import (
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system/context"
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
		subSystems: []sub_system.SubSystem{sub_system.NewEmptySubSystem(), ws},
	}, nil
}

func (s *System) OnAction(action input.KeyInputAction) error {
	return stlslices.Last(s.subSystems).OnAction(newCursor(s), action)
}

func (s *System) OnUpdate() error {
	return stlslices.Last(s.subSystems).OnUpdate(newCursor(s))
}

func (s *System) OnDraw(drawer draw.OptionDrawer) error {
	return stlslices.Last(s.subSystems).OnDraw(newCursor(s), drawer)
}
