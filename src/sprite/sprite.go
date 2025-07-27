package sprite

import (
	"fmt"

	"github.com/lafriks/go-tiled"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/util"
)

type Behavior string

var BehaviorEnum = enum.New[struct {
	Walk Behavior `enum:"walk"`
	Run  Behavior `enum:"run"`
	// 无动画
	Talk   Behavior `enum:"talk"`
	Script Behavior `enum:"script"`
}]()

type UpdateInfo interface {
	UpdateInfo()
}

type Sprite interface {
	util.Drawer
	SetPosition(x, y int)
	Position() (int, int)
	NextStepPosition() (int, int)
	GetScript() string
	OnAction(ctx context.Context, action input.Action, info UpdateInfo) error
	Update(ctx context.Context, info UpdateInfo) error
}

var spriteCreateFuncMap = make(map[string]func(object *tiled.Object) (Sprite, error))

func RegisterCreateFunc(classes []string, fn func(object *tiled.Object) (Sprite, error)) {
	for _, class := range classes {
		spriteCreateFuncMap[class] = fn
	}
}

func NewSprite(object *tiled.Object) (Sprite, error) {
	fn, ok := spriteCreateFuncMap[object.Type]
	if !ok {
		return nil, fmt.Errorf("not found sprite class `%s`", object.Type)
	}
	return fn(object)
}
