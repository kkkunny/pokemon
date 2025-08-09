package battle

import (
	"image/color"
	"path/filepath"
	"time"

	"golang.org/x/image/font"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/pokemon"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	ctx context.Context

	active    bool
	siteImage imgutil.Image // 战斗场地

	pok                *pokemon.PokemonRace
	pokLastUpdateTime  time.Time
	pokFrontFrameIndex int
	pokBackFrameIndex  int
}

func NewSystem(ctx context.Context) (*System, error) {
	pok, err := pokemon.NewPokemonRace(1)
	if err != nil {
		return nil, err
	}
	return &System{
		ctx: ctx,
		pok: pok,
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
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 32).UnsafeInternal(), displayText)
	return (bounds.Max.X - bounds.Min.X).Round() / len([]rune(displayText)), (bounds.Max.Y - bounds.Min.Y).Round()
}

func (s *System) drawPokemonStatusCard(drawer draw.OptionDrawer) {
	draw.PrepareDrawRect(drawer, 300, 80, util.NewNRGBColor(248, 248, 216)).SetBorderWidth(5).SetBorderColor(color.Black).Draw()
	opponentName := s.ctx.Localisation().Get("pokemon.1")
	opponentNameBounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 26).UnsafeInternal(), opponentName)
	draw.PrepareDrawText(drawer, opponentName, util.GetFont(util.FontTypeEnum.Normal, 26), color.Black).Move(20, 10).Draw()
	genderText := "♂"
	genderBounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Emoji, 16).UnsafeInternal(), opponentName)
	draw.PrepareDrawText(drawer, genderText, util.GetFont(util.FontTypeEnum.Emoji, 16), util.NewNRGBColor(65, 200, 248)).Move(20+opponentNameBounds.Max.X.Round(), 10+opponentNameBounds.Max.Y.Round()-genderBounds.Max.Y.Round()).Draw()
	draw.PrepareDrawText(drawer, "Lv  5", util.GetFont(util.FontTypeEnum.Normal, 26), color.Black).Move(220, 10).Draw()
	draw.PrepareDrawRect(drawer, 220, 20, util.NewNRGBColor(80, 104, 88)).Move(70, 50).SetRadius(7).Draw()
	draw.PrepareDrawText(drawer, "HP", util.GetFont(util.FontTypeEnum.Normal, 20), util.NewNRGBColor(248, 178, 65)).Move(76, 50).Draw()
	draw.PrepareDrawRect(drawer, 192, 16, color.White).Move(96, 52).SetRadius(5).Draw()
	draw.PrepareDrawRect(drawer, 188, 12, util.NewNRGBColor(80, 104, 88)).Move(98, 54).SetRadius(3).Draw()
	hp := 100
	draw.PrepareDrawRect(drawer, int(float64(188)/100*float64(hp)), 12, util.NewNRGBColor(110, 245, 165)).Move(98, 54).SetRadius(3).Draw()
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
	s.drawPokemonStatusCard(drawer.Move(80, 50))

	// 我方
	fontW, fontH := s.frontSize()
	_, bgH := fontW*(19+2), fontH*(2+2)

	selfSiteX, selfSiteY := 0, screenHeight-bgH-10-s.siteImage.Bounds().Dy()/3*2
	draw.PrepareDrawImage(drawer, s.siteImage).Move(selfSiteX, selfSiteY).Draw()
	pokemonImage = s.pok.Back.Image[s.pokBackFrameIndex]
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(selfSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, selfSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()
	s.drawPokemonStatusCard(drawer.Move(340, 250))

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
