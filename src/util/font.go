package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kkkunny/stl/container/tuple"
	"github.com/tnnmigga/enum"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
)

type FontType uint8

var FontTypeEnum = enum.New[struct {
	Normal FontType
	Emoji  FontType
}]()

var innerFontCache = make(map[FontType]*opentype.Font)
var fontFaceCache = make(map[tuple.Tuple2[FontType, int]]*text.GoXFace)

func init() {
	fontNames := enum.Keys(FontTypeEnum)
	fontTypeEnums := enum.Values[FontType](FontTypeEnum)
	for i, fontName := range fontNames {
		fontData, err := os.ReadFile(filepath.Join(config.FontsPath, strings.ToLower(fontName)+".ttf"))
		if err != nil {
			panic(err)
		}
		fontInst, err := opentype.Parse(fontData)
		if err != nil {
			panic(err)
		}
		innerFontCache[fontTypeEnums[i]] = fontInst
	}
}

func GetFont(fontType FontType, size int) *text.GoXFace {
	fontFace, ok := fontFaceCache[tuple.Pack2(fontType, size)]
	if ok {
		return fontFace
	}
	innerFont, ok := innerFontCache[fontType]
	if !ok {
		panic("unknown font type")
	}
	stdFontFace, err := opentype.NewFace(innerFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
	fontFace = text.NewGoXFace(stdFontFace)
	fontFaceCache[tuple.Pack2(fontType, size)] = fontFace
	return fontFace
}
