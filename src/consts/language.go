package consts

import "github.com/tnnmigga/enum"

type Language string

var LanguageEnum = enum.New[struct {
	ZhCN Language `enum:"zh_cn"`
}]()
