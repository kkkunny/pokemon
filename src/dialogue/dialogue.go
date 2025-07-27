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
	"github.com/kkkunny/pokemon/src/util"
)

type DialogueSystem struct {
	fontFace        *text.GoXFace
	displayInterval time.Duration

	display        bool
	text           []rune
	index          int
	lastUpdateTime time.Time
	lines          [][]rune
}

func NewDialogueSystem(cfg *config.Config) (*DialogueSystem, error) {
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
	return &DialogueSystem{
		displayInterval: time.Millisecond * 150,
		fontFace:        text.NewGoXFace(fontFace),
	}, nil
}

func (d *DialogueSystem) SetDisplay(v bool) {
	d.display = v
}

func (d *DialogueSystem) Display() bool {
	return d.display
}

func (d *DialogueSystem) ResetText(s string) {
	d.text = []rune(s)
	d.index = 0
	d.lines = nil
	d.lastUpdateTime = time.Time{}
}

func (d *DialogueSystem) DisplayText(s string) {
	d.ResetText(s)
	d.SetDisplay(true)
}

func (d *DialogueSystem) StreamDone() bool {
	return d.index > len(d.text)-1
}

func (d *DialogueSystem) Draw(screen *ebiten.Image) error {
	if !d.display {
		return nil
	}

	const bgHInterval, bgVInterval = 40, 30
	const outerLayerWidth, innerLayerWidth = 6, 18
	const fontVSpacing = 10

	screenW, screenH := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	bounds, _ := font.BoundString(d.fontFace.UnsafeInternal(), "好")
	_, fontH := float64((bounds.Max.X - bounds.Min.X).Floor()), float64((bounds.Max.Y - bounds.Min.Y).Floor())
	bgW, bgH := float64(screenW)-bgHInterval*2, fontH*2+fontVSpacing*2+innerLayerWidth*2
	fgW, _ := bgW-innerLayerWidth*2, bgH-innerLayerWidth*2

	// 背景
	x, y := float64(bgHInterval), screenH-bgH-bgVInterval
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(bgW), float32(bgH), util.NewRGBColor(248, 248, 255), false)
	vector.StrokeRect(screen, float32(x+outerLayerWidth), float32(y+outerLayerWidth), float32(bgW-innerLayerWidth+outerLayerWidth), float32(bgH-innerLayerWidth+outerLayerWidth), innerLayerWidth, util.NewRGBColor(176, 196, 222), false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(bgW), float32(bgH), outerLayerWidth, util.NewRGBColor(119, 136, 153), false)

	// 文字
	fontColor := util.NewRGBColor(119, 136, 153)
	x, y = x+innerLayerWidth, y+innerLayerWidth+fontVSpacing/2

	if len(d.lines) != 0 {
		// 存量行（第一行）
		var textOps text.DrawOptions
		textOps.ColorScale.ScaleWithColor(fontColor)
		textOps.GeoM.Translate(x, y)
		text.Draw(screen, string(d.lines[len(d.lines)-1]), d.fontFace, &textOps)

		y += fontH + fontVSpacing
	}

	// 输出行（第二行或第一行）
	var renderedIndex int
	for _, line := range d.lines {
		renderedIndex += len(line)
	}
	renderText := d.text[renderedIndex:d.index]

	var textOps text.DrawOptions
	textOps.ColorScale.ScaleWithColor(fontColor)
	textOps.GeoM.Translate(x, y)
	text.Draw(screen, string(renderText), d.fontFace, &textOps)

	if d.StreamDone() || (d.lastUpdateTime != stlval.Default[time.Time]() && time.Now().Sub(d.lastUpdateTime) < d.displayInterval) {
		return nil
	}

	// 更新
	d.lastUpdateTime = time.Now()
	d.index++

	// 换行
	nextBounds, _ := font.BoundString(d.fontFace.UnsafeInternal(), string(d.text[renderedIndex:d.index]))
	nextRenderW := float64((nextBounds.Max.X - nextBounds.Min.X).Floor())
	changeLine := nextRenderW > fgW

	if d.index <= len(d.text)-1 && d.text[d.index-1] == '\n' {
		changeLine = true
		renderText = d.text[renderedIndex:d.index]
		d.index++
	}
	if changeLine {
		d.lines = append(d.lines, renderText)
	}
	return nil
}
