package pokemon

import (
	"fmt"
	"image/gif"
	"os"
	"path/filepath"

	"github.com/kkkunny/pokemon/src/config"
)

// PokemonRace 宝可梦种族
type PokemonRace struct {
	ID    int16    // 图鉴编号
	Type  Type     // 属性
	Front *gif.GIF // 战斗正面图
}

func NewPokemonRace(id int16) (*PokemonRace, error) {
	dirpath := filepath.Join(config.PokemonDefinePath, fmt.Sprintf("%d", id))
	dirinfo, err := os.Stat(dirpath)
	if (err != nil && os.IsNotExist(err)) || (err == nil && !dirinfo.IsDir()) {
		return nil, fmt.Errorf("not exist pokemon, id=%d", id)
	} else if err != nil {
		return nil, err
	}

	frontFile, err := os.Open(filepath.Join(dirpath, "front.gif"))
	if err != nil {
		return nil, err
	}
	defer frontFile.Close()
	frontGif, err := gif.DecodeAll(frontFile)
	if err != nil {
		return nil, err
	}
	return &PokemonRace{
		ID:    id,
		Front: frontGif,
	}, nil
}
