package battle

import (
	"image/color"
	"path/filepath"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/util/draw"
	"github.com/kkkunny/pokemon/src/util/image"
)

type System struct {
	ctx       context.Context
	active    bool
	siteImage *image.Image // 战斗场地
}

func NewSystem(ctx context.Context) (*System, error) {
	return &System{ctx: ctx}, nil
}

func (s *System) Active() bool {
	return s.active
}

func (s *System) StartOneBattle(site string) error {
	var err error
	s.siteImage, err = image.NewImageFromFile(filepath.Join(config.GFXBattleSitesPath, site+".png"))
	if err != nil {
		return err
	}
	s.active = true
	return nil
}

func (s *System) OnUpdate() error {
	return nil
}

func (s *System) OnDraw(drawer draw.Drawer) error {
	err := drawer.OverlayColor(color.White)
	if err != nil {
		return err
	}
	screenWidth, screenHeight := drawer.Size()
	return drawer.Move(float64(screenWidth-s.siteImage.Width()*s.ctx.Config().Scale), float64(screenHeight/2-s.siteImage.Height()*s.ctx.Config().Scale)).DrawImage(s.siteImage)
}
