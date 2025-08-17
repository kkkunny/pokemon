package config

import (
	"github.com/kkkunny/pokemon/src/consts"
)

const (
	Scale    = 2  // 放大倍数
	TileSize = 16 // 地图原大小
)

var (
	// 屏幕宽
	ScreenHeight int
	// 屏幕高
	ScreenWidth int
	// 语言
	Language consts.Language
)

func Init() {
	ScreenWidth, ScreenHeight = 720, 480
	Language = consts.LanguageEnum.ZhCN
}
