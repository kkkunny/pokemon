package sub_system

type SubSystemManager interface {
	Pop()
	DisplayLabel(text string) error
	DisplayDialogue(text string) error
}
