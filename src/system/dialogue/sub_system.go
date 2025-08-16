package dialogue

import (
	"strings"
	"time"

	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"

	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/system/sub_system"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	"github.com/kkkunny/pokemon/src/util/image"
)

const (
	waitForContinueChar     = '🔻'
	normalDisplayInterval   = time.Millisecond * 150
	fastModeDisplayInterval = time.Millisecond * 30
)

type DialogueSystem struct {
	ctx             context.Context
	display         bool
	displayInterval time.Duration

	// 显示文字的必备属性
	isDialogue     bool
	text           []rune
	index          int
	lastUpdateTime time.Time
	waitFrame      int
}

func NewDialogueSystem(ctx context.Context) (*DialogueSystem, error) {
	return &DialogueSystem{
		ctx:             ctx,
		displayInterval: normalDisplayInterval,
	}, nil
}

func (s *DialogueSystem) Type() sub_system.SubSystemType {
	return sub_system.SubSystemTypeEnum.Dialogue
}

func (s *DialogueSystem) OnAction(system sub_system.SubSystemManager, action input.KeyInputAction) error {
	if !s.display {
		return system.Next().OnAction(system, action)
	}

	switch {
	case s.waitForContinue() && action == input.KeyInputActionEnum.A.Pressed():
		s.continueNext()
	case s.streamDone() && action == input.KeyInputActionEnum.A.Pressed():
		s.SetDisplay(false)
	case s.fastMode() && action == input.KeyInputActionEnum.A.Released():
		s.setFastMode(false)
	case !s.fastMode() && action == input.KeyInputActionEnum.A.Pressed():
		s.setFastMode(true)
	}
	return nil
}

func (s *DialogueSystem) OnUpdate(system sub_system.SubSystemManager) error {
	if !s.display {
		return system.Next().OnUpdate(system)
	}

	if s.waitForContinue() {
		if time.Since(s.lastUpdateTime) > s.displayInterval {
			s.waitFrame = (s.waitFrame + 1) % 6
			s.lastUpdateTime = time.Now()
		}
		return system.Next().OnUpdate(system)
	} else if s.streamDone() || (s.lastUpdateTime != stlval.Default[time.Time]() && time.Since(s.lastUpdateTime) < s.displayInterval) {
		return system.Next().OnUpdate(system)
	}

	s.lastUpdateTime = time.Now()
	s.index++
	return system.Next().OnUpdate(system)
}

func (s *DialogueSystem) OnDraw(system sub_system.SubSystemManager, drawer draw.OptionDrawer) error {
	err := system.Next().OnDraw(system, drawer)
	if err != nil {
		return err
	}

	if !s.display {
		return nil
	}

	_fontW, _fontH := s.frontSize()
	fontW, fontH := float64(_fontW), float64(_fontH)
	_screenW, _screenH := drawer.Bounds().Dx(), drawer.Bounds().Dy()
	screenW, screenH := float64(_screenW), float64(_screenH)
	hFrontMaxCount, vFrontMaxCount := int(screenW/fontW)-4, int(screenH/fontH)-4
	if hFrontMaxCount < 2 || vFrontMaxCount < 3 {
		return nil
	}

	// 背景
	bgImg := stlval.Ternary(s.isDialogue, s.getDialogueBackground, s.getLabelBackground)(hFrontMaxCount, 2)
	x, y := (screenW-float64(bgImg.Bounds().Dx()))/2, screenH-float64(bgImg.Bounds().Dy())-fontH
	draw.PrepareDrawImage(drawer, bgImg).Move(int(x), int(y)).Draw()

	// 文字
	fontColor := util.NewNRGBColor(100, 100, 100)

	x, y = x+fontW/2+fontW/4, y+fontH/2+fontH/3

	lines := s.splitDoneLines(s.text[:stlval.Ternary(s.index < len(s.text), s.index+1, s.index)], hFrontMaxCount)
	if len(lines) > 1 {
		// 存量行（第一行）
		renderStr := strings.Replace(string(lines[len(lines)-2]), string([]rune{waitForContinueChar}), "", -1)
		draw.PrepareDrawText(drawer, renderStr, util.GetFont(util.FontTypeEnum.Normal, 36), fontColor).Move(int(x), int(y)).Draw()

		y += fontH + fontH/3
	}

	// 输出行（第二行或第一行）
	renderStr := strings.Replace(string(lines[len(lines)-1]), string([]rune{waitForContinueChar}), "", -1)
	draw.PrepareDrawText(drawer, renderStr, util.GetFont(util.FontTypeEnum.Normal, 36), fontColor).Move(int(x), int(y)).Draw()

	if s.waitForContinue() {
		bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 36).UnsafeInternal(), renderStr)
		x += float64((bounds.Max.X - bounds.Min.X).Round())
		y += (fontH/5)*3 - stlval.Ternary(s.waitFrame < 4, float64(s.waitFrame), float64(6-s.waitFrame))
		waitString := string([]rune{waitForContinueChar})
		bounds, _ = font.BoundString(util.GetFont(util.FontTypeEnum.Emoji, 36).UnsafeInternal(), renderStr)
		y -= float64((bounds.Max.Y - bounds.Min.Y).Round()) / 2
		draw.PrepareDrawText(drawer, waitString, util.GetFont(util.FontTypeEnum.Emoji, 36), util.NewNRGBColor(224, 8, 8)).Move(int(x), int(y)).Draw()
		return nil
	}
	return nil
}

func (s *DialogueSystem) SetLabel(text string) {
	s.isDialogue = false
	s.text = []rune(text)
	s.index = 0
	s.lastUpdateTime = time.Time{}
}

func (s *DialogueSystem) SetDialogue(text string) {
	s.isDialogue = true
	s.text = []rune(text)
	s.index = 0
	s.lastUpdateTime = time.Time{}
}

func (s *DialogueSystem) streamDone() bool {
	return s.index > len(s.text)-1
}

func (s *DialogueSystem) frontSize() (int, int) {
	displayText := s.ctx.Localisation().Get("game_name")
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 36).UnsafeInternal(), displayText)
	return (bounds.Max.X - bounds.Min.X).Round() / len([]rune(displayText)), (bounds.Max.Y - bounds.Min.Y).Round()
}

func (s *DialogueSystem) getLabelBackground(w, h int) imgutil.Image {
	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(w+2), fontH*(h+2)

	img := imgutil.NewImage(bgW, bgH)
	draw.PrepareDrawRect(img, bgW, bgH, util.NewNRGBColor(104, 112, 120)).Draw()
	draw.PrepareDrawRect(img, bgW-fontW/2, bgH-fontH/2, util.NewNRGBColor(200, 200, 216)).Move(fontW/4, fontH/4).Draw()
	draw.PrepareDrawRect(img, bgW-fontW, bgH-fontH, util.NewNRGBColor(248, 248, 248)).Move(fontW/2, fontH/2).Draw()
	return img
}

func (s *DialogueSystem) getDialogueBackground(w, h int) imgutil.Image {
	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(w+2), fontH*(h+2)

	img := imgutil.NewImage(bgW, bgH)
	draw.PrepareDrawRect(img, bgW, bgH, util.NewNRGBColor(160, 208, 224)).SetRadius(fontW / 2).Draw()
	draw.PrepareDrawRect(img, bgW-fontW/2, bgH-fontH/2, util.NewNRGBColor(224, 240, 248)).SetRadius(fontW/2).Move(fontW/4, fontH/4).Draw()
	draw.PrepareDrawRect(img, bgW-fontW, bgH-fontH, util.NewNRGBColor(248, 248, 248)).SetRadius(fontW/2).Move(fontW/2, fontH/2).Draw()
	return img
}

func (s *DialogueSystem) splitDoneLines(text []rune, maxLineCount int) (lines [][]rune) {
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

func (s *DialogueSystem) setFastMode(v bool) {
	s.displayInterval = stlval.Ternary(v, fastModeDisplayInterval, normalDisplayInterval)
}

func (s *DialogueSystem) fastMode() bool {
	return s.displayInterval != normalDisplayInterval
}

func (s *DialogueSystem) waitForContinue() bool {
	if s.index >= len(s.text) {
		return false
	}
	return s.text[s.index] == waitForContinueChar
}

func (s *DialogueSystem) continueNext() {
	if !s.waitForContinue() {
		return
	}
	s.index++
}

func (s *DialogueSystem) Display() bool {
	return s.display
}

func (s *DialogueSystem) SetDisplay(v bool) {
	s.display = v
}
