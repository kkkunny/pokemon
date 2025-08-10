package pokemon

import (
	"strings"

	"github.com/tnnmigga/enum"
)

// Shape 体型
type Shape uint8

var ShapeEnum = enum.New[struct {
	Unknown         Shape // ???
	Spherical       Shape // 球形
	Cylindrical     Shape // 柱形
	Snake           Shape // 蛇形
	Bipedal         Shape // 双腿形
	FourLeggedBeast Shape // 四足兽形
	Tentacle        Shape // 触手形
	Insect          Shape // 虫形
	TwoHanded       Shape // 双手形
	DoubleWinged    Shape // 双翅形
	MultiWinged     Shape // 多翅形
	Fish            Shape // 鱼形
	Humanoid        Shape // 人形
	BipedalBeast    Shape // 双足兽形
	Combined        Shape // 组合形
}]()

// 将string转成分类
func ParseShape(s string) Shape {
	for i, name := range enum.Keys(ShapeEnum) {
		if strings.ToLower(name) == s {
			return enum.Values[Shape](ShapeEnum)[i]
		}
	}
	return ShapeEnum.Unknown
}
