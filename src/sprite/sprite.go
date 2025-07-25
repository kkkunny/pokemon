package sprite

import (
	"fmt"

	input "github.com/quasilyte/ebitengine-input"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/util"
)

type UpdateInfo interface {
	UpdateInfo()
}

type Sprite interface {
	util.Drawer
	SetPosition(x, y int)
	OnAction(cfg *config.Config, action input.Action, info UpdateInfo)
	Update(cfg *config.Config, info UpdateInfo) error
}

var spriteCreateFuncMap = make(map[string]func(class string, imageName string) (Sprite, error))

func RegisterCreateFunc(classes []string, fn func(class string, imageName string) (Sprite, error)) {
	for _, class := range classes {
		spriteCreateFuncMap[class] = fn
	}
}

func NewSprite(class string, imageName string) (Sprite, error) {
	fn, ok := spriteCreateFuncMap[class]
	if !ok {
		return nil, fmt.Errorf("not found sprite class `%s`", class)
	}
	return fn(class, imageName)
}
