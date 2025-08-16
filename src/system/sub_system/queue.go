package sub_system

type SubSystemQueue interface {
	Pop()
	DisplayLabel(text string) error
	DisplayDialogue(text string) error
}
