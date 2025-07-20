package consts

import "github.com/tnnmigga/enum"

type Direction int8

var DirectionEnum = enum.New[struct {
	Up    Direction `enum:"-1"`
	Down  Direction `enum:"1"`
	Left  Direction `enum:"-3"`
	Right Direction `enum:"3"`
}]()
