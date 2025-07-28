package dialogue

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	fontFace        *text.GoXFace
	displayInterval time.Duration

	display        bool
	text           []rune
	index          int
	lastUpdateTime time.Time
	lines          [][]rune
}

func NewSystem(cfg *config.Config) (*System, error) {
	// 字体
	fontBytes, err := os.ReadFile(filepath.Join(config.FontsPath, cfg.MaterFontName) + ".ttf")
	if err != nil {
		return nil, err
	}
	fontInst, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	fontFace, err := opentype.NewFace(fontInst, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}
	return &System{
		displayInterval: time.Millisecond * 150,
		fontFace:        text.NewGoXFace(fontFace),
	}, nil
}

func (s *System) SetDisplay(v bool) {
	s.display = v
}

func (s *System) Display() bool {
	return s.display
}

func (s *System) ResetText(text string) {
	s.text = []rune(text)
	s.index = 0
	s.lines = nil
	s.lastUpdateTime = time.Time{}
}

func (s *System) DisplayText(text string) {
	s.ResetText(text)
	s.SetDisplay(true)
}

func (s *System) StreamDone() bool {
	return s.index > len(s.text)-1
}

func (s *System) Draw(screen *image.Image) error {
	if !s.display {
		return nil
	}

	const bgHInterval, bgVInterval = 40, 30
	const outerLayerWidth, innerLayerWidth = 6, 18
	const fontVSpacing = 10

	screenW, screenH := float64(screen.Width()), float64(screen.Height())
	bounds, _ := font.BoundString(s.fontFace.UnsafeInternal(), "好")
	_, fontH := float64((bounds.Max.X - bounds.Min.X).Floor()), float64((bounds.Max.Y - bounds.Min.Y).Floor())
	bgW, bgH := float64(screenW)-bgHInterval*2, fontH*2+fontVSpacing*2+innerLayerWidth*2
	fgW, _ := bgW-innerLayerWidth*2, bgH-innerLayerWidth*2

	// 背景
	x, y := float64(bgHInterval), screenH-bgH-bgVInterval
	vector.DrawFilledRect(screen.Image, float32(x), float32(y), float32(bgW), float32(bgH), util.NewRGBColor(248, 248, 255), false)
	vector.StrokeRect(screen.Image, float32(x+outerLayerWidth), float32(y+outerLayerWidth), float32(bgW-innerLayerWidth+outerLayerWidth), float32(bgH-innerLayerWidth+outerLayerWidth), innerLayerWidth, util.NewRGBColor(176, 196, 222), false)
	vector.StrokeRect(screen.Image, float32(x), float32(y), float32(bgW), float32(bgH), outerLayerWidth, util.NewRGBColor(119, 136, 153), false)

	// 文字
	fontColor := util.NewRGBColor(119, 136, 153)
	x, y = x+innerLayerWidth, y+innerLayerWidth+fontVSpacing/2

	if len(s.lines) != 0 {
		// 存量行（第一行）
		var textOps text.DrawOptions
		textOps.ColorScale.ScaleWithColor(fontColor)
		textOps.GeoM.Translate(x, y)
		text.Draw(screen.Image, string(s.lines[len(s.lines)-1]), s.fontFace, &textOps)

		y += fontH + fontVSpacing
	}

	// 输出行（第二行或第一行）
	var renderedIndex int
	for _, line := range s.lines {
		renderedIndex += len(line)
	}
	renderText := s.text[renderedIndex:s.index]

	var textOps text.DrawOptions
	textOps.ColorScale.ScaleWithColor(fontColor)
	textOps.GeoM.Translate(x, y)
	text.Draw(screen.Image, string(renderText), s.fontFace, &textOps)

	if s.StreamDone() || (s.lastUpdateTime != stlval.Default[time.Time]() && time.Now().Sub(s.lastUpdateTime) < s.displayInterval) {
		return nil
	}

	// 更新
	s.lastUpdateTime = time.Now()
	s.index++

	// 换行
	nextBounds, _ := font.BoundString(s.fontFace.UnsafeInternal(), string(s.text[renderedIndex:s.index]))
	nextRenderW := float64((nextBounds.Max.X - nextBounds.Min.X).Floor())
	changeLine := nextRenderW > fgW

	if s.index <= len(s.text)-1 && s.text[s.index-1] == '\n' {
		changeLine = true
		renderText = s.text[renderedIndex:s.index]
		s.index++
	}
	if changeLine {
		s.lines = append(s.lines, renderText)
	}
	return nil
}
