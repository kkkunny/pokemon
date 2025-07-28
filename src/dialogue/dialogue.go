package dialogue

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	ctx context.Context

	fontFace        *text.GoXFace
	displayInterval time.Duration

	display        bool
	text           []rune
	index          int
	lastUpdateTime time.Time
	lines          [][]rune
}

func NewSystem(ctx context.Context) (*System, error) {
	// 字体
	fontBytes, err := os.ReadFile(filepath.Join(config.FontsPath, ctx.Config().MaterFontName) + ".ttf")
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
		ctx:             ctx,
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

func (s *System) frontSize() (int, int) {
	displayText := s.ctx.Localisation().Get("game_name")
	bounds, _ := font.BoundString(s.fontFace.UnsafeInternal(), displayText)
	return (bounds.Max.X - bounds.Min.X).Round() / len([]rune(displayText)), (bounds.Max.Y - bounds.Min.Y).Round()
}

func (s *System) getDialogBackground(w, h int) *image.Image {
	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(w+2), fontH*(h+2)

	img := image.NewImage(bgW, bgH)
	vector.DrawFilledRect(img.Image, 0, 0, float32(bgW), float32(bgH), util.NewRGBColor(119, 136, 153), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/4, float32(fontH)/4, float32(bgW)-float32(fontW)/2, float32(bgH)-float32(fontH)/2, util.NewRGBColor(176, 196, 222), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/2, float32(fontH)/2, float32(bgW)-float32(fontW), float32(bgH)-float32(fontH), util.NewRGBColor(248, 248, 255), false)
	return img
}

func (s *System) Draw(screen *image.Image) error {
	if !s.display {
		return nil
	}

	fontW, fontH := s.frontSize()
	screenW, screenH := float64(screen.Width()), float64(screen.Height())
	hFrontMaxCount, vFrontMaxCount := int(screenW/float64(fontW))-4, int(screenH/float64(fontH))-4
	if hFrontMaxCount < 2 || vFrontMaxCount < 3 {
		return nil
	}

	// 背景
	bgImg := s.getDialogBackground(hFrontMaxCount, 2)
	x, y := (screenW-float64(bgImg.Width()))/2, screenH-float64(bgImg.Height())-float64(fontH)
	var bgOps ebiten.DrawImageOptions
	bgOps.GeoM.Translate(x, y)
	screen.DrawImage(bgImg, &bgOps)

	// 文字
	fontColor := util.NewRGBColor(119, 136, 153)

	x, y = x+float64(fontW)/2+float64(fontW)/4, y+float64(fontH)/2+float64(fontH)/3

	if len(s.lines) != 0 {
		// 存量行（第一行）
		var textOps text.DrawOptions
		textOps.ColorScale.ScaleWithColor(fontColor)
		textOps.GeoM.Translate(x, y)
		text.Draw(screen.Image, string(s.lines[len(s.lines)-1]), s.fontFace, &textOps)

		y += float64(fontH) + float64(fontH)/3
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
	changeLine := len(renderText) >= hFrontMaxCount

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
