package consts

import "github.com/tnnmigga/enum"

type Direction int8

var DirectionEnum = enum.New[struct {
	Up    Direction `enum:"-1"`
	Down  Direction `enum:"1"`
	Left  Direction `enum:"-3"`
	Right Direction `enum:"3"`
}]()

func ParseDirection(s string) Direction {
	switch s {
	case "up":
		return DirectionEnum.Up
	case "down":
		return DirectionEnum.Down
	case "left":
		return DirectionEnum.Left
	case "right":
		return DirectionEnum.Right
	default:
		return DirectionEnum.Up
	}
}

func (d Direction) String() string {
	switch d {
	case DirectionEnum.Up:
		return "up"
	case DirectionEnum.Down:
		return "down"
	case DirectionEnum.Left:
		return "left"
	case DirectionEnum.Right:
		return "right"
	default:
		return ""
	}
}
