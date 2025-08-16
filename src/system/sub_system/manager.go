package sub_system

import (
	"image/color"

	stlval "github.com/kkkunny/stl/value"
)

type SubSystemManager interface {
	SubSystems() []SubSystem
	Next() SubSystem
	DisplayText(boxStyle BoxStyle, text string, fontColor color.Color) error
	StartOneBattle(site string) error
}

func ExtractSubSystem[T SubSystem](ssm SubSystemManager) T {
	for _, s := range ssm.SubSystems() {
		ss, ok := s.(T)
		if ok {
			return ss
		}
	}
	return stlval.Default[T]()
}
