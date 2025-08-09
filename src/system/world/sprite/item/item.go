package item

import (
	"github.com/lafriks/go-tiled"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/system/world/sprite"
	"github.com/kkkunny/pokemon/src/util/draw"
)

func init() {
	sprite.RegisterCreateFunc([]string{"label"}, func(object *tiled.Object) (sprite.Sprite, error) {
		item, err := NewItemByTile(object)
		if err != nil {
			return nil, err
		}
		return item, nil
	})
}

type Item interface {
	sprite.Sprite
}

type _Item struct {
	pos        [2]int            // 位置
	actionType sprite.ActionType // 交互类型
	script     string            // 脚本id
	text       string            // 对话文本
}

func NewItem() (Item, error) {
	return &_Item{}, nil
}

func NewItemByTile(object *tiled.Object) (Item, error) {
	return &_Item{
		actionType: sprite.ActionType(object.Properties.GetString("action_type")),
		script:     object.Properties.GetString("script"),
		text:       object.Properties.GetString("text"),
	}, nil
}

func (i *_Item) ActionType() sprite.ActionType {
	return i.actionType
}

func (i *_Item) Position() (int, int) {
	return i.pos[0], i.pos[1]
}

func (i *_Item) SetPosition(x int, y int) {
	i.pos = [2]int{x, y}
}

func (i *_Item) Collision() bool {
	return true
}

func (i *_Item) CollisionPosition() (int, int) {
	return i.Position()
}

func (i *_Item) GetScript() string {
	return i.script
}

func (i *_Item) OnAction(_ context.Context, _ input.KeyInputAction, _ sprite.UpdateInfo) error {
	return nil
}

func (i *_Item) Update(_ context.Context, _ sprite.UpdateInfo) error {
	return nil
}

func (i *_Item) Draw(_ context.Context, _ draw.OptionDrawer) error {
	return nil
}

func (i *_Item) GetText() string {
	return i.text
}
