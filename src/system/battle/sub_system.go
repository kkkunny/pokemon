package battle

import (
	"fmt"
	"image/color"
	"strings"

	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/pokemon"
	"github.com/kkkunny/pokemon/src/system/context"
	"github.com/kkkunny/pokemon/src/system/sub_system"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type BattleSystem struct {
	ctx context.Context

	active    bool
	siteImage imgutil.Image // 战斗场地

	pmRace *pokemon.Race
	pms    [2]*pokemon.Pokemon
}

func NewBattleSystem(ctx context.Context) (*BattleSystem, error) {
	pok, err := pokemon.LoadPokemonRace(1)
	if err != nil {
		return nil, err
	}
	pms := [2]*pokemon.Pokemon{pok.RandomPokemon(), pok.RandomPokemon()}
	return &BattleSystem{
		ctx:    ctx,
		pmRace: pok,
		pms:    pms,
	}, nil
}

func (s *BattleSystem) Active() bool {
	return s.active
}

func (s *BattleSystem) Type() sub_system.SubSystemType {
	return sub_system.SubSystemTypeEnum.Battle
}

func (s *BattleSystem) StartOneBattle(site string) error {
	siteImage, err := util.FindFileAndThenParse(config.GFXBattleSitesPath, site, imgutil.NewImageFromFile)
	if err != nil {
		return err
	}
	s.siteImage = siteImage.Scale(config.Scale, config.Scale)
	s.active = true
	return nil
}

func (s *BattleSystem) OnAction(action input.KeyInputAction) error {
	return nil
}

func (s *BattleSystem) OnUpdate() error {
	return nil
}

func (s *BattleSystem) frontSize() (int, int) {
	displayText := s.ctx.Localisation().Get("game_name")
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 32).UnsafeInternal(), displayText)
	return (bounds.Max.X - bounds.Min.X).Round() / len([]rune(displayText)), (bounds.Max.Y - bounds.Min.Y).Round()
}

func (s *BattleSystem) drawPokemonType(drawer draw.OptionDrawer, typ pokemon.Type) {
	typ = typ.Flatten()[0]
	draw.PrepareDrawRect(drawer, 55, 16, typ.Color()).SetRadius(7).Draw()
	icon, ok := pokemon.GetTypeIcon(typ)
	if ok {
		scale := 14 / float64(icon.Bounds().Dy())
		draw.PrepareDrawImage(drawer, icon).Scale(scale, scale).Move(2, 1).Draw()
	}
	typeName := s.ctx.Localisation().Get(fmt.Sprintf("pokemon_type.%s", typ))
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 11).UnsafeInternal(), typeName)
	x, y := (16+52)/2-bounds.Max.X.Round()/2, 3
	draw.PrepareDrawText(drawer, typeName, util.GetFont(util.FontTypeEnum.Normal, 11), color.White).Move(x, y).Draw()
}

func (s *BattleSystem) drawPokemonStatusCard(drawer draw.OptionDrawer, pm *pokemon.Pokemon) {
	draw.PrepareDrawRect(drawer, 300, 80, util.NewNRGBColor(248, 248, 216)).SetBorderWidth(5).SetBorderColor(color.Black).Draw()
	opponentName := s.ctx.Localisation().Get(fmt.Sprintf("pokemon.%d", pm.ID))
	opponentNameBounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 26).UnsafeInternal(), opponentName)
	types := pm.Type.Flatten()
	if len(types) == 1 {
		s.drawPokemonType(drawer.Move(5, 16), types[0])
	} else {
		s.drawPokemonType(drawer.Move(5, 7), types[0])
		s.drawPokemonType(drawer.Move(5, 25), types[1])
	}
	draw.PrepareDrawText(drawer, opponentName, util.GetFont(util.FontTypeEnum.Normal, 26), color.Black).Move(65, 10).Draw()
	genderText := stlval.Ternary(pm.Gender, consts.MaleText, consts.FemaleText)
	genderTextColor := stlval.Ternary(pm.Gender, consts.MaleColor, consts.FemaleColor)
	genderBounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Emoji, 16).UnsafeInternal(), opponentName)
	draw.PrepareDrawText(drawer, genderText, util.GetFont(util.FontTypeEnum.Emoji, 16), genderTextColor).Move(65+opponentNameBounds.Max.X.Round(), 10+opponentNameBounds.Max.Y.Round()-genderBounds.Max.Y.Round()).Draw()
	level := fmt.Sprintf("%03d", pm.Level)
	if simLevel := strings.TrimPrefix(level, "0"); len(simLevel) != len(level) {
		level = strings.Repeat(" ", len(level)-len(simLevel)) + simLevel
	}
	draw.PrepareDrawText(drawer, "Lv"+level, util.GetFont(util.FontTypeEnum.Normal, 26), color.Black).Move(224, 10).Draw()
	draw.PrepareDrawRect(drawer, 220, 20, util.NewNRGBColor(80, 104, 88)).Move(70, 50).SetRadius(7).Draw()
	draw.PrepareDrawText(drawer, "HP", util.GetFont(util.FontTypeEnum.Normal, 20), util.NewNRGBColor(248, 178, 65)).Move(76, 50).Draw()
	draw.PrepareDrawRect(drawer, 192, 16, color.White).Move(96, 52).SetRadius(5).Draw()
	draw.PrepareDrawRect(drawer, 188, 12, util.NewNRGBColor(80, 104, 88)).Move(98, 54).SetRadius(3).Draw()
	hpRatio := float64(pm.CurrentHP) / float64(pm.SpeciesStrength.HP()) * 100
	draw.PrepareDrawRect(drawer, int(float64(188)/100*hpRatio), 12, util.NewNRGBColor(110, 245, 165)).Move(98, 54).SetRadius(3).Draw()
}

func (s *BattleSystem) OnDraw(drawer draw.OptionDrawer) error {
	draw.OverlayColor(drawer, color.White)

	screenWidth, screenHeight := drawer.Bounds().Dx(), drawer.Bounds().Dy()

	// 敌方
	opponentSiteX, opponentSiteY := screenWidth-s.siteImage.Bounds().Dx(), screenHeight/2-s.siteImage.Bounds().Dy()
	draw.PrepareDrawImage(drawer, s.siteImage).Move(opponentSiteX, opponentSiteY).Draw()
	s.pmRace.Front.Update()
	pokemonImage := s.pmRace.Front.GetCurrentFrameImage()
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(opponentSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, opponentSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()
	s.drawPokemonStatusCard(drawer.Move(80, 50), s.pms[0])

	// 我方
	fontW, fontH := s.frontSize()
	_, bgH := fontW*(19+2), fontH*(2+2)

	selfSiteX, selfSiteY := 0, screenHeight-bgH-10-s.siteImage.Bounds().Dy()/3*2
	draw.PrepareDrawImage(drawer, s.siteImage).Move(selfSiteX, selfSiteY).Draw()
	s.pmRace.Back.Update()
	pokemonImage = s.pmRace.Back.GetCurrentFrameImage()
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(selfSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, selfSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()
	s.drawPokemonStatusCard(drawer.Move(340, 250), s.pms[0])

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
