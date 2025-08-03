package battle

import (
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	"github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	ctx      context.Context
	fontFace *text.GoXFace

	active    bool
	siteImage *image.Image // 战斗场地
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
		ctx:      ctx,
		fontFace: text.NewGoXFace(fontFace),
	}, nil
}

func (s *System) Active() bool {
	return s.active
}

func (s *System) StartOneBattle(site string) error {
	siteImage, err := image.NewImageFromFile(filepath.Join(config.GFXBattleSitesPath, site+".png"))
	if err != nil {
		return err
	}
	s.siteImage = siteImage.Scale(config.Scale, config.Scale)
	s.active = true
	return nil
}

func (s *System) OnUpdate() error {
	return nil
}

func (s *System) frontSize() (int, int) {
	displayText := s.ctx.Localisation().Get("game_name")
	bounds, _ := font.BoundString(s.fontFace.UnsafeInternal(), displayText)
	return (bounds.Max.X - bounds.Min.X).Round() / len([]rune(displayText)), (bounds.Max.Y - bounds.Min.Y).Round()
}

func (s *System) OnDraw(drawer draw.Drawer) error {
	err := drawer.OverlayColor(color.White)
	if err != nil {
		return err
	}

	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(w+2), fontH*(h+2)

	img := image.NewImage(bgW, bgH)
	vector.DrawFilledRect(img.Image, 0, 0, float32(bgW), float32(bgH), util.NewRGBColor(104, 112, 120), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/4, float32(fontH)/4, float32(bgW)-float32(fontW)/2, float32(bgH)-float32(fontH)/2, util.NewRGBColor(200, 200, 216), false)
	vector.DrawFilledRect(img.Image, float32(fontW)/2, float32(fontH)/2, float32(bgW)-float32(fontW), float32(bgH)-float32(fontH), util.NewRGBColor(248, 248, 248), false)

	screenWidth, screenHeight := drawer.Size()
	err = drawer.Move(float64(screenWidth-s.siteImage.Width()), float64(screenHeight/2-s.siteImage.Height())).DrawImage(s.siteImage)
	if err != nil {
		return err
	}
	return nil
}
