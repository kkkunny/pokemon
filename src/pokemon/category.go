package pokemon

import (
	"strings"

	"github.com/tnnmigga/enum"
)

// Category 分类
type Category uint8

var CategoryEnum = enum.New[struct {
	Unknown Category // ???
	Seed    Category // 种子宝可梦
}]()

// 将string转成分类
func ParseCategory(s string) Category {
	for i, name := range enum.Keys(CategoryEnum) {
		if strings.ToLower(name) == s {
			return enum.Values[Category](CategoryEnum)[i]
		}
	}
	return CategoryEnum.Unknown
}
