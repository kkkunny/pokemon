package battle

import (
	"image/color"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/pokemon"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	ctx           context.Context
	fontFace      *text.GoXFace
	emojiFontFace *text.GoXFace

	active    bool
	siteImage imgutil.Image // 战斗场地

	pok                *pokemon.PokemonRace
	pokLastUpdateTime  time.Time
	pokFrontFrameIndex int
	pokBackFrameIndex  int
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
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	pok, err := pokemon.NewPokemonRace(1)
	if err != nil {
		return nil, err
	}
	return &System{
		ctx:           ctx,
		fontFace:      text.NewGoXFace(fontFace),
		emojiFontFace: text.NewGoXFace(emojiFontFace),
		pok:           pok,
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

	screenWidth, screenHeight := drawer.Bounds().Dx(), drawer.Bounds().Dy()

	// 敌方
	opponentSiteX, opponentSiteY := screenWidth-s.siteImage.Bounds().Dx(), screenHeight/2-s.siteImage.Bounds().Dy()
	draw.PrepareDrawImage(drawer, s.siteImage).Move(opponentSiteX, opponentSiteY).Draw()
	if time.Now().Sub(s.pokLastUpdateTime) > time.Millisecond*70 {
		s.pokLastUpdateTime = time.Now()
		s.pokFrontFrameIndex = (s.pokFrontFrameIndex + 1) % len(s.pok.Front.Image)
		s.pokBackFrameIndex = (s.pokBackFrameIndex + 1) % len(s.pok.Back.Image)
	}
	pokemonImage := s.pok.Front.Image[s.pokFrontFrameIndex]
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(opponentSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, opponentSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()

	draw.PrepareDrawRect(drawer, 300, 100, util.NewNRGBColor(248, 248, 216)).SetBorderWidth(5).SetBorderColor(color.Black).Move(0, 0).Draw()
	draw.PrepareDrawText(drawer, s.ctx.Localisation().Get("pokemon.1"), s.fontFace, color.Black).Scale(0.8, 0.8).Move(0, 0).Draw()
	draw.PrepareDrawText(drawer, "♂", s.emojiFontFace, util.NewNRGBColor(65, 200, 248)).Move(100, 50).Draw()

	// 我方
	fontW, fontH := s.frontSize()
	_, bgH := fontW*(19+2), fontH*(2+2)

	selfSiteX, selfSiteY := 0, screenHeight-bgH-10-s.siteImage.Bounds().Dy()/3*2
	draw.PrepareDrawImage(drawer, s.siteImage).Move(selfSiteX, selfSiteY).Draw()
	pokemonImage = s.pok.Back.Image[s.pokBackFrameIndex]
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(selfSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, selfSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()

	// 对话栏

	// 对话栏总背景
	draw.PrepareDrawRect(drawer, screenWidth, bgH+10, color.Black).Move(0, screenHeight-bgH-10).Draw()

	// 对话栏背景
	draw.PrepareDrawRect(drawer, screenWidth-10, bgH, util.NewNRGBColor(200, 168, 72)).Move(5, screenHeight-bgH-5).SetRadius(10).Draw()
	draw.PrepareDrawRect(drawer, screenWidth-30, bgH-20, util.NewNRGBColor(224, 216, 224)).Move(15, screenHeight-bgH+5).SetRadius(4).Draw()
	draw.PrepareDrawRect(drawer, screenWidth-40, bgH-30, util.NewNRGBColor(40, 80, 104)).Move(20, screenHeight-bgH+10).Draw()

	// 行为框背景
	draw.PrepareDrawRect(drawer, screenWidth/2, bgH+10, color.Black).Move(screenWidth/2, screenHeight-bgH-10).Draw()
	draw.PrepareDrawRect(drawer, screenWidth/2-10, bgH, util.NewNRGBColor(132, 131, 188)).Move(screenWidth/2+5, screenHeight-bgH-5).SetRadius(4).Draw()
	draw.PrepareDrawRect(drawer, screenWidth/2-14, bgH-4, util.NewNRGBColor(112, 104, 128)).Move(screenWidth/2+7, screenHeight-bgH-3).Draw()
	draw.PrepareDrawRect(drawer, screenWidth/2-24, bgH-14, util.NewNRGBColor(248, 248, 248)).Move(screenWidth/2+12, screenHeight-bgH+2).SetRadius(6).Draw()
	return nil
}
