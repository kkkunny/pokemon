package sprite

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/util/image"
)

type Behavior string

var BehaviorEnum = enum.New[struct {
	Walk Behavior `enum:"walk"`
	Run  Behavior `enum:"run"`
	// 无动画
	Talk   Behavior `enum:"talk"`
	Script Behavior `enum:"script"`
}]()

type ActionType string

var ActionTypeEnum = enum.New[struct {
	None     ActionType `enum:""`
	Script   ActionType `enum:"script"`
	Label    ActionType `enum:"label"`
	Dialogue ActionType `enum:"dialogue"`
}]()

type UpdateInfo interface {
	UpdateInfo()
}

type Sprite interface {
	ActionType() ActionType
	GetScript() string
	GetText() string

	SetPosition(x, y int)
	Position() (int, int)

	OnAction(ctx context.Context, action input.KeyInputAction, info UpdateInfo) error
	Update(ctx context.Context, info UpdateInfo) error
	Draw(ctx context.Context, screen *image.Image, options ebiten.DrawImageOptions) error
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
