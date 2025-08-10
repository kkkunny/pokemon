package pokemon

import (
	"encoding/json"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strconv"

	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/animation"
)

// Race 宝可梦种族
type Race struct {
	ID              uint16          // 图鉴编号
	Type            Type            // 属性
	Category        Category        // 类别
	Height          float64         // 身高（m）
	Weight          float64         // 体重（kg）
	PokedexColor    color.Color     // 图鉴颜色
	CatchRate       float64         // 捕获概率(%)
	MaleRate        float64         // 男性概率(%)
	SpeciesStrength SpeciesStrength // 种族值
	BaseExperience  int             // 基础经验值

	Front *animation.Animation // 战斗正面图
	Back  *animation.Animation // 战斗背面图
}

// define结构体
type pokemonRaceDefine struct {
	Type            []string        `json:"type"`
	Category        string          `json:"category"`
	Height          float64         `json:"height"`
	Weight          float64         `json:"weight"`
	Shape           string          `json:"shape"`
	PokeDexColor    [3]uint8        `json:"pokedex_color"`
	CatchRate       float64         `json:"catch_rate"`
	MaleRate        float64         `json:"male_rate"`
	SpeciesStrength SpeciesStrength `json:"species_strength"`
	BaseExperience  uint            `json:"base_experience"`
}

func LoadPokemonRace(id uint16) (*Race, error) {
	dirpath := filepath.Join(config.PokemonDefinePath, strconv.FormatUint(uint64(id), 10))
	definePath := filepath.Join(dirpath, "define.json")
	defineData, err := os.ReadFile(definePath)
	if err != nil {
		return nil, err
	}
	var define pokemonRaceDefine
	err = json.Unmarshal(defineData, &define)
	if err != nil {
		return nil, err
	}

	types := stlslices.Map(define.Type, func(_ int, s string) Type {
		return ParseType(s)
	})
	var typ Type
	if stlslices.Contain(types, TypeEnum.Unknown) {
		typ = TypeEnum.Unknown
	} else {
		for _, t := range types {
			typ |= t
		}
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

	backFile, err := os.Open(filepath.Join(dirpath, "back.gif"))
	if err != nil {
		return nil, err
	}
	defer backFile.Close()
	backGif, err := gif.DecodeAll(backFile)
	if err != nil {
		return nil, err
	}
	return &Race{
		ID:              id,
		Type:            typ,
		Category:        ParseCategory(define.Category),
		Height:          define.Height,
		Weight:          define.Weight,
		PokedexColor:    util.NewNRGBColor(define.PokeDexColor[0], define.PokeDexColor[1], define.PokeDexColor[2]),
		CatchRate:       define.CatchRate,
		MaleRate:        define.MaleRate,
		SpeciesStrength: define.SpeciesStrength,
		BaseExperience:  int(define.BaseExperience),

		Front: animation.NewAnimationFromGIF(frontGif),
		Back:  animation.NewAnimationFromGIF(backGif),
	}, nil
}
