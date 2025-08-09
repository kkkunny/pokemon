package config

const (
	Scale    = 2  // 放大倍数
	TileSize = 16 // 地图原大小
)

type Config struct {
	ScreenWidth, ScreenHeight int
	MaterFontName             string
}

func NewConfig() *Config {
	return &Config{
		ScreenWidth:   720,
		ScreenHeight:  480,
		MaterFontName: "fusion-pixel-12px-monospaced-zh_hans",
	}
}
