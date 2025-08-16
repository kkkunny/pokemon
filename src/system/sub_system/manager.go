package sub_system

import stlval "github.com/kkkunny/stl/value"

type SubSystemManager interface {
	SubSystems() []SubSystem
	Next() SubSystem
	DisplayLabel(text string) error
	DisplayDialogue(text string) error
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
