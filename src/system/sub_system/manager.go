package sub_system

type SubSystemManager interface {
	Pop()
	Next() SubSystem
	DisplayLabel(text string) error
	DisplayDialogue(text string) error
}
