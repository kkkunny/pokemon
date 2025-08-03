package config

const (
	Scale    = 2
	TileSize = 16
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
