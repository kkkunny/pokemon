package config

type Config struct {
	ScreenWidth, ScreenHeight int
	Scale                     int // 放大倍数
	TileSize                  int // 地块像素点大小
	MaterFontName             string
}

func NewConfig() *Config {
	return &Config{
		ScreenWidth:   720,
		ScreenHeight:  480,
		Scale:         1,
		TileSize:      16,
		MaterFontName: "fusion-pixel-12px-monospaced-zh_hans",
	}
}
