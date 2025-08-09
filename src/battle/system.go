package battle

import (
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	ctx      context.Context
	fontFace *text.GoXFace

	active    bool
	siteImage imgutil.Image // 战斗场地
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
	siteImage, err := imgutil.NewImageFromFile(filepath.Join(config.GFXBattleSitesPath, site+".png"))
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

func (s *System) OnDraw(drawer draw.OptionDrawer) error {
	draw.OverlayColor(drawer, color.White)

	fontW, fontH := s.frontSize()
	bgW, bgH := fontW*(19+2), fontH*(2+2)

	img := imgutil.NewImage(bgW, bgH)
	draw.PrepareDrawRect(img, bgW, bgH, util.NewNRGBColor(104, 112, 120)).Draw()
	draw.PrepareDrawRect(img, fontW/4, fontH/4, util.NewNRGBColor(200, 200, 216)).Move(bgW-fontW/2, bgH-fontH/2).Draw()
	draw.PrepareDrawRect(img, fontW/2, fontH/2, util.NewNRGBColor(248, 248, 248)).Move(fontW/2, fontH/2).Draw()

	screenWidth, screenHeight := drawer.Bounds().Dx(), drawer.Bounds().Dy()
	draw.PrepareDrawImage(drawer, s.siteImage).Move(screenWidth-s.siteImage.Bounds().Dx(), screenHeight/2-s.siteImage.Bounds().Dy()).Draw()
	return nil
}
