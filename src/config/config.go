package config

type Config struct {
	ScreenWidth, ScreenHeight int
	Scale                     int // 放大倍数
	TileSize                  int // 地块像素点大小
}

func NewConfig() *Config {
	return &Config{
		ScreenWidth:  720,
		ScreenHeight: 480,
		Scale:        1,
		TileSize:     16,
	}
}
