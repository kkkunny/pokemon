package sprite

import "github.com/kkkunny/pokemon/src/consts"

type Talker interface {
	Sprite
	Direction() consts.Direction
	Turn(d consts.Direction) bool
	TalkTo(talker Talker, passive bool) (bool, error)
	EndTalk() error
}
