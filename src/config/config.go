package config

const (
	Scale    = 2  // 放大倍数
	TileSize = 16 // 地图原大小
)

var (
	// 屏幕宽
	ScreenHeight int
	// 屏幕高
	ScreenWidth int
)

func Init() {
	ScreenWidth, ScreenHeight = 720, 480
}
