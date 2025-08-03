package dialogue

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	"github.com/kkkunny/pokemon/src/util/image"
)

const (
	normalDisplayInterval   = time.Millisecond * 150
	fastModeDisplayInterval = time.Millisecond * 30
)

type System struct {
	ctx context.Context

	fontFace        *text.GoXFace
	emojiFontFace   *text.GoXFace
	displayInterval time.Duration

	// 显示文字的必备属性
	display        bool
	isDialogue     bool
	text           []rune
	index          int
	lastUpdateTime time.Time
	waitFrame      int
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
	// emoji字体
	fontBytes, err = os.ReadFile(filepath.Join(config.FontsPath, "NotoEmoji-VariableFont_wght") + ".ttf")
	if err != nil {
		return nil, err
	}
	fontInst, err = opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	emojiFontFace, err := opentype.NewFace(fontInst, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}
	return &System{
		ctx:             ctx,
		displayInterval: normalDisplayInterval,
		fontFace:        text.NewGoXFace(fontFace),
		emojiFontFace:   text.NewGoXFace(emojiFontFace),
	}, nil
}

func (s *System) SetDisplay(v bool) {
	s.display = v
}

func (s *System) Display() bool {
	return s.display
}

func (s *System) SetLabel(text string) {
	s.isDialogue = false
	s.text = []rune(text)
	s.index = 0
	s.lastUpdateTime = time.Time{}
}

func (s *System) DisplayLabel(text string) {
	s.SetLabel(text)
	s.SetDisplay(true)
}

func (s *System) SetDialogue(text string) {
	s.isDialogue = true
	s.text = []rune(text)
	s.index = 0
	s.lastUpdateTime = time.Time{}
}

func (s *System) DisplayDialogue(text string) {
	s.SetDialogue(text)
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

func (s *System) getLabelBackground(w, h int) *image.Image {
	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(w+2), fontH*(h+2)

	img := image.NewImage(bgW, bgH)
	vector.DrawFilledRect(img.Image, 0, 0, float32(bgW), float32(bgH), util.NewRGBColor(104, 112, 120), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/4, float32(fontH)/4, float32(bgW)-float32(fontW)/2, float32(bgH)-float32(fontH)/2, util.NewRGBColor(200, 200, 216), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/2, float32(fontH)/2, float32(bgW)-float32(fontW), float32(bgH)-float32(fontH), util.NewRGBColor(248, 248, 248), false)
	return img
}

func (s *System) getDialogueBackground(w, h int) *image.Image {
	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(w+2), fontH*(h+2)

	img := image.NewImage(bgW, bgH)
	vector.DrawFilledRect(img.Image, 0, 0, float32(bgW), float32(bgH), util.NewRGBColor(160, 208, 224), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/4, float32(fontH)/4, float32(bgW)-float32(fontW)/2, float32(bgH)-float32(fontH)/2, util.NewRGBColor(224, 240, 248), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/2, float32(fontH)/2, float32(bgW)-float32(fontW), float32(bgH)-float32(fontH), util.NewRGBColor(248, 248, 248), false)
	return img
}

func (s *System) splitDoneLines(text []rune, maxLineCount int) (lines [][]rune) {
	var beginIndex, curIndex int
	for _, ch := range text {
		if ch == '\n' {
			lines = append(lines, text[beginIndex:curIndex])
			curIndex++
			beginIndex = curIndex
		} else if curIndex-beginIndex >= maxLineCount {
			lines = append(lines, text[beginIndex:curIndex])
			beginIndex = curIndex
			curIndex++
		} else {
			curIndex++
		}
	}
	if curIndex > beginIndex {
		lines = append(lines, text[beginIndex:curIndex])
	}
	return lines
}

func (s *System) OnDraw(drawer draw.Drawer) error {
	if !s.display {
		return nil
	}

	fontW, fontH := s.frontSize()
	screenW, screenH := drawer.Size()
	hFrontMaxCount, vFrontMaxCount := int(float64(screenW)/float64(fontW))-4, int(float64(screenH)/float64(fontH))-4
	if hFrontMaxCount < 2 || vFrontMaxCount < 3 {
		return nil
	}

	// 背景
	bgImg := stlval.Ternary(s.isDialogue, s.getDialogueBackground, s.getLabelBackground)(hFrontMaxCount, 2)
	x, y := (screenW-bgImg.Width())/2, screenH-bgImg.Height()-fontH
	err := drawer.Move(x, y).DrawImage(bgImg)
	if err != nil {
		return err
	}

	// 文字
	fontColor := util.NewRGBColor(100, 100, 100)

	x, y = x+fontW/2+fontW/4, y+fontH/2+fontH/3

	lines := s.splitDoneLines(s.text[:stlval.Ternary(s.index < len(s.text), s.index+1, s.index)], hFrontMaxCount)
	if len(lines) > 1 {
		// 存量行（第一行）
		renderStr := strings.Replace(string(lines[len(lines)-2]), string([]rune{consts.WaitForContinueChar}), "", -1)
		err = drawer.Move(x, y).ScaleWithColor(fontColor).DrawText(renderStr, s.fontFace)
		if err != nil {
			return err
		}

		y += fontH + fontH/3
	}

	// 输出行（第二行或第一行）
	renderStr := strings.Replace(string(lines[len(lines)-1]), string([]rune{consts.WaitForContinueChar}), "", -1)
	err = drawer.Move(x, y).ScaleWithColor(fontColor).DrawText(renderStr, s.fontFace)
	if err != nil {
		return err
	}

	if s.WaitForContinue() {
		bounds, _ := font.BoundString(s.fontFace.UnsafeInternal(), renderStr)
		x += (bounds.Max.X - bounds.Min.X).Round()
		y += (fontH/5)*2 + s.waitFrame
		waitString := string([]rune{consts.WaitForContinueChar})
		bounds, _ = font.BoundString(s.emojiFontFace.UnsafeInternal(), renderStr)
		waitFontHeight := (bounds.Max.Y - bounds.Min.Y).Round()
		y -= waitFontHeight / 2
		err = drawer.Move(x, y).ScaleWithColor(util.NewRGBColor(224, 8, 8)).DrawText(waitString, s.emojiFontFace)
		if time.Since(s.lastUpdateTime) > s.displayInterval*2 {
			s.waitFrame = (s.waitFrame + 1) % 3
			s.lastUpdateTime = time.Now()
		}
		return nil
	} else if s.StreamDone() || (s.lastUpdateTime != stlval.Default[time.Time]() && time.Since(s.lastUpdateTime) < s.displayInterval) {
		return nil
	}

	s.lastUpdateTime = time.Now()
	s.index++
	return nil
}

func (s *System) SetFastMode(v bool) {
	s.displayInterval = stlval.Ternary(v, fastModeDisplayInterval, normalDisplayInterval)
}

func (s *System) FastMode() bool {
	return s.displayInterval != normalDisplayInterval
}

func (s *System) WaitForContinue() bool {
	if !s.Display() || s.index >= len(s.text) {
		return false
	}
	return s.text[s.index] == consts.WaitForContinueChar
}

func (s *System) Continue() {
	if !s.WaitForContinue() {
		return
	}
	s.index++
}
