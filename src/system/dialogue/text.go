package dialogue

import (
	"image"
	"image/color"
	"path/filepath"
	"strings"
	"time"

	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"

	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/system/sub_system"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

const (
	waitForContinueChar     = '🔻'
	normalDisplayInterval   = time.Millisecond * 150
	fastModeDisplayInterval = time.Millisecond * 30
)

var waitIcon imgutil.Image

func init() {
	icon, err := imgutil.NewImageFromFile(filepath.Join(consts.GFXInterfacePath, "red_wedge_1.png"))
	if err != nil {
		panic(err)
	}
	waitIcon = icon
}

type Text struct {
	boxStyle  sub_system.BoxStyle
	needDrop  bool
	fontColor color.Color

	displayInterval time.Duration

	text           []rune
	index          int
	lastUpdateTime time.Time
	waitFrame      int
}

func NewText(botStyle sub_system.BoxStyle, text string, fontColor color.Color) (*Text, error) {
	return &Text{
		boxStyle:        botStyle,
		fontColor:       fontColor,
		displayInterval: normalDisplayInterval,
		text:            []rune(text),
	}, nil
}

func (s *Text) OnAction(system sub_system.SubSystemManager, action input.KeyInputAction) error {
	switch {
	case s.waitForContinue() && action == input.KeyInputActionEnum.A.Pressed():
		s.continueNext()
	case s.streamDone() && action == input.KeyInputActionEnum.A.Pressed():
		s.needDrop = true
	case s.fastMode() && action == input.KeyInputActionEnum.A.Released():
		s.setFastMode(false)
	case !s.fastMode() && action == input.KeyInputActionEnum.A.Pressed():
		s.setFastMode(true)
	}
	return nil
}

func (s *Text) OnUpdate(system sub_system.SubSystemManager) error {
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

func (s *Text) OnDraw(system sub_system.SubSystemManager, drawer draw.OptionDrawer) error {
	err := system.Next().OnDraw(system, drawer)
	if err != nil {
		return err
	}

	// 背景
	_, innerRect, err := s.boxStyle.OnUpdate(drawer, image.Rect(-1, -1, -1, -1))
	if err != nil {
		return err
	}

	// 文字
	displayText := i18n.Get("game_name")
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 36).UnsafeInternal(), displayText)
	fontW, fontH := (bounds.Max.X-bounds.Min.X).Round()/len([]rune(displayText)), (bounds.Max.Y - bounds.Min.Y).Round()
	hSep := 10
	hFrontMaxCount := (innerRect.Dx() - hSep*2) / fontW

	x, y := innerRect.Min.X+hSep, (innerRect.Min.Y+(innerRect.Min.Y+innerRect.Max.Y)/2)/2-fontH/2

	lines := s.splitDoneLines(s.text[:stlval.Ternary(s.index < len(s.text), s.index+1, s.index)], hFrontMaxCount)
	if len(lines) > 1 {
		// 存量行（第一行）
		renderStr := strings.Replace(string(lines[len(lines)-2]), string([]rune{waitForContinueChar}), "", -1)
		draw.PrepareDrawText(drawer, renderStr, util.GetFont(util.FontTypeEnum.Normal, 36), s.fontColor).Move(x, y).Draw()

		y = ((innerRect.Min.Y+innerRect.Max.Y)/2+innerRect.Max.Y)/2 - fontH/2
	}

	// 输出行（第二行或第一行）
	renderStr := strings.Replace(string(lines[len(lines)-1]), string([]rune{waitForContinueChar}), "", -1)
	draw.PrepareDrawText(drawer, renderStr, util.GetFont(util.FontTypeEnum.Normal, 36), s.fontColor).Move(x, y).Draw()

	// 等待符
	if s.waitForContinue() {
		x += fontW*len([]rune(renderStr)) + fontW/4
		y += (fontH/5)*3 - stlval.Ternary(s.waitFrame < 4, s.waitFrame, 6-s.waitFrame)
		ratio := 20 / float64(waitIcon.Bounds().Dx())
		draw.PrepareDrawImage(drawer, waitIcon).Scale(ratio, ratio).Move(x, y).Draw()
		return nil
	}
	return nil
}

func (s *Text) Drop() bool {
	return s.needDrop
}

func (s *Text) streamDone() bool {
	return s.index > len(s.text)-1
}

func (s *Text) splitDoneLines(text []rune, maxLineCount int) (lines [][]rune) {
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

func (s *Text) setFastMode(v bool) {
	s.displayInterval = stlval.Ternary(v, fastModeDisplayInterval, normalDisplayInterval)
}

func (s *Text) fastMode() bool {
	return s.displayInterval != normalDisplayInterval
}

func (s *Text) waitForContinue() bool {
	if s.index >= len(s.text) {
		return false
	}
	return s.text[s.index] == waitForContinueChar
}

func (s *Text) continueNext() {
	if !s.waitForContinue() {
		return
	}
	s.index++
}
