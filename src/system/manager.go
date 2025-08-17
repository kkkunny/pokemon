package system

import (
	"image/color"

	stlval "github.com/kkkunny/stl/value"
)

type SystemManager interface {
	SubSystems() []System
	Next() System
	DisplayText(boxStyle BoxStyle, text string, fontColor color.Color) error
	StartOneBattle(site string) error
}

func ExtractSystem[T System](ssm SystemManager) T {
	for _, s := range ssm.SubSystems() {
		ss, ok := s.(T)
		if ok {
			return ss
		}
	}
	return stlval.Default[T]()
}
