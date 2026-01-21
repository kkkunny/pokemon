package battle

import (
	"fmt"
	"image/color"
	"strings"

	stlval "github.com/kkkunny/stl/value"
	"golang.org/x/image/font"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/i18n"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/pokemon"
	"github.com/kkkunny/pokemon/src/system"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/animation"
	"github.com/kkkunny/pokemon/src/util/draw"
	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type BattleSystem struct {
	active    bool
	siteImage imgutil.Image // 战斗场地

	// 宝可梦
	pmRace       *pokemon.Race
	pms          [2]*pokemon.Pokemon
	pmAnimations [2]*animation.Player

	// 动作选项
	actionSelect int8
}

func NewBattleSystem() (*BattleSystem, error) {
	pok, err := pokemon.LoadPokemonRace(1)
	if err != nil {
		return nil, err
	}
	pms := [2]*pokemon.Pokemon{pok.RandomPokemon(), pok.RandomPokemon()}
	pmAnimations := [2]*animation.Player{
		pok.Back.NewPlayer(config.DefaultFPS),
		pok.Front.NewPlayer(config.DefaultFPS),
	}
	return &BattleSystem{
		pmRace:       pok,
		pms:          pms,
		pmAnimations: pmAnimations,
	}, nil
}

func (s *BattleSystem) StartOneBattle(site string) error {
	siteImage, err := util.FindFileAndThenParse(consts.GFXBattleSitesPath, site, imgutil.NewImageFromFile)
	if err != nil {
		return err
	}
	s.siteImage = siteImage.Scale(config.Scale, config.Scale)
	s.active = true
	return nil
}

func (s *BattleSystem) OnAction(system system.SystemManager, action input.KeyInputAction) error {
	if !s.active {
		return system.Next().OnAction(system, action)
	}

	switch {
	case action == input.KeyInputActionEnum.MoveUp.Released():
		s.actionSelect = (s.actionSelect + 3) % 4
	case action == input.KeyInputActionEnum.MoveDown.Released():
		s.actionSelect = (s.actionSelect + 1) % 4
	case action == input.KeyInputActionEnum.MoveLeft.Released():
		s.actionSelect = (s.actionSelect + 2) % 4
	case action == input.KeyInputActionEnum.MoveRight.Released():
		s.actionSelect = (s.actionSelect + 2) % 4
	}

	return nil
}

func (s *BattleSystem) OnUpdate(system system.SystemManager) error {
	if !s.active {
		return system.Next().OnUpdate(system)
	}
	return nil
}

func (s *BattleSystem) OnDraw(system system.SystemManager, drawer draw.OptionDrawer) error {
	if !s.active {
		return system.Next().OnDraw(system, drawer)
	}

	draw.OverlayColor(drawer, color.White)

	screenWidth, screenHeight := drawer.Bounds().Dx(), drawer.Bounds().Dy()

	// 敌方
	opponentSiteX, opponentSiteY := screenWidth-s.siteImage.Bounds().Dx(), screenHeight/2-s.siteImage.Bounds().Dy()
	draw.PrepareDrawImage(drawer, s.siteImage).Move(opponentSiteX, opponentSiteY).Draw()
	s.pmAnimations[1].Update()
	pokemonImage := s.pmAnimations[1].GetCurrentFrame()
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(opponentSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, opponentSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()
	s.drawPokemonStatusCard(drawer.Move(80, 50), s.pms[0])

	// 我方
	fontW, fontH := s.frontSize()
	_, bgH := fontW*(19+2), fontH*(2+2)

	selfSiteX, selfSiteY := 0, screenHeight-bgH-10-s.siteImage.Bounds().Dy()/3*2
	draw.PrepareDrawImage(drawer, s.siteImage).Move(selfSiteX, selfSiteY).Draw()
	s.pmAnimations[0].Update()
	pokemonImage = s.pmAnimations[0].GetCurrentFrame()
	draw.PrepareDrawImage(drawer, pokemonImage).Scale(config.Scale, config.Scale).Move(selfSiteX+s.siteImage.Bounds().Dx()/2-pokemonImage.Bounds().Dx()/2*config.Scale, selfSiteY+s.siteImage.Bounds().Dy()/4*3-pokemonImage.Bounds().Dy()*config.Scale).Draw()
	s.drawPokemonStatusCard(drawer.Move(340, 250), s.pms[0])

	// 对话栏
	draw.PrepareDrawRect(drawer, screenWidth, 126, color.Black).Move(0, screenHeight-126).Draw()

	// 行为框背景
	draw.PrepareDrawRect(drawer, screenWidth/2, bgH+10, color.Black).Move(screenWidth/2, screenHeight-bgH-10).Draw()
	draw.PrepareDrawRect(drawer, screenWidth/2-10, bgH, util.NewNRGBColor(132, 131, 188)).Move(screenWidth/2+5, screenHeight-bgH-5).SetRadius(4).Draw()
	draw.PrepareDrawRect(drawer, screenWidth/2-14, bgH-4, util.NewNRGBColor(112, 104, 128)).Move(screenWidth/2+7, screenHeight-bgH-3).Draw()
	draw.PrepareDrawRect(drawer, screenWidth/2-24, bgH-14, util.NewNRGBColor(248, 248, 248)).Move(screenWidth/2+12, screenHeight-bgH+2).SetRadius(6).Draw()

	draw.PrepareDrawText(drawer, i18n.Get("battle"), util.GetFont(util.FontTypeEnum.Normal, 32), color.Black).Move(screenWidth/2+50, screenHeight-bgH+14).Draw()
	draw.PrepareDrawText(drawer, i18n.Get("backpack"), util.GetFont(util.FontTypeEnum.Normal, 32), color.Black).Move(screenWidth/2+240, screenHeight-bgH+14).Draw()
	draw.PrepareDrawText(drawer, i18n.Get("pokemons"), util.GetFont(util.FontTypeEnum.Normal, 32), color.Black).Move(screenWidth/2+50, screenHeight-bgH+60).Draw()
	draw.PrepareDrawText(drawer, i18n.Get("escape"), util.GetFont(util.FontTypeEnum.Normal, 32), color.Black).Move(screenWidth/2+240, screenHeight-bgH+60).Draw()
	actionSelectDrawer := draw.PrepareDrawText(drawer, "🔻", util.GetFont(util.FontTypeEnum.Emoji, 32), color.Black)
	switch s.actionSelect {
	case 0:
		actionSelectDrawer = actionSelectDrawer.Move(screenWidth/2+10, screenHeight-bgH+14)
	case 1:
		actionSelectDrawer = actionSelectDrawer.Move(screenWidth/2+10, screenHeight-bgH+60)
	case 2:
		actionSelectDrawer = actionSelectDrawer.Move(screenWidth/2+200, screenHeight-bgH+14)
	case 3:
		actionSelectDrawer = actionSelectDrawer.Move(screenWidth/2+200, screenHeight-bgH+60)
	}
	actionSelectDrawer.Draw()
	return nil
}

func (s *BattleSystem) Drop() bool {
	return false
}

func (s *BattleSystem) Active() bool {
	return s.active
}

func (s *BattleSystem) frontSize() (int, int) {
	displayText := i18n.Get("game_name")
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
	typeName := i18n.Get(fmt.Sprintf("pokemon_type.%s", typ))
	bounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 11).UnsafeInternal(), typeName)
	x, y := (16+52)/2-bounds.Max.X.Round()/2, 3
	draw.PrepareDrawText(drawer, typeName, util.GetFont(util.FontTypeEnum.Normal, 11), color.White).Move(x, y).Draw()
}

func (s *BattleSystem) drawPokemonStatusCard(drawer draw.OptionDrawer, pm *pokemon.Pokemon) {
	draw.PrepareDrawRect(drawer, 300, 80, util.NewNRGBColor(248, 248, 216)).SetBorderWidth(5).SetBorderColor(color.Black).Draw()
	pokemonName := i18n.Get(fmt.Sprintf("pokemon.%d", pm.ID))
	pokemonNameBounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Normal, 26).UnsafeInternal(), pokemonName)
	types := pm.Type.Flatten()
	if len(types) == 1 {
		s.drawPokemonType(drawer.Move(5, 16), types[0])
	} else {
		s.drawPokemonType(drawer.Move(5, 7), types[0])
		s.drawPokemonType(drawer.Move(5, 25), types[1])
	}
	draw.PrepareDrawText(drawer, pokemonName, util.GetFont(util.FontTypeEnum.Normal, 26), color.Black).Move(65, 10).Draw()
	genderText := stlval.Ternary(pm.Gender, consts.MaleText, consts.FemaleText)
	genderTextColor := stlval.Ternary(pm.Gender, consts.MaleColor, consts.FemaleColor)
	genderBounds, _ := font.BoundString(util.GetFont(util.FontTypeEnum.Emoji, 16).UnsafeInternal(), pokemonName)
	draw.PrepareDrawText(drawer, genderText, util.GetFont(util.FontTypeEnum.Emoji, 16), genderTextColor).Move(65+pokemonNameBounds.Max.X.Round(), 10+pokemonNameBounds.Max.Y.Round()-genderBounds.Max.Y.Round()).Draw()
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

func (s *BattleSystem) GetSelfPokemon() *pokemon.Pokemon {
	return s.pms[1]
}
