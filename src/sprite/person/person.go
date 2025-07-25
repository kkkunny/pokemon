package person

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/kkkunny/pokemon/src/animation"
	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/sprite"
)

func init() {
	sprite.RegisterCreateFunc([]string{"person"}, func(_ string, imageName string) (sprite.Sprite, error) {
		person, err := NewPerson(imageName)
		if err != nil {
			return nil, err
		}
		return person, nil
	})
}

type Person struct {
	// 静态资源
	directionImages    map[consts.Direction]*ebiten.Image                                            // 方向图片
	behaviorAnimations map[sprite.Behavior]map[consts.Direction]map[sprite.Foot]*animation.Animation // 行为动画
	// 属性
	direction        consts.Direction // 当前所处方向
	moveStartingFoot sprite.Foot      // 移动时的起始脚
	speed            int              // 移动速度
	pos              [2]int           // 当前地块位置
	expectPos        [2]int           // 预期所处的地块位置，用于移动

	moveCounter int // 移动时的计数器，用于显示动画
}

func NewPerson(name string) (*Person, error) {
	dirpath := filepath.Join("./resource/map_item/people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}

	directionImages, behaviorAnimations, err := sprite.LoadPersonAnimations(name, sprite.BehaviorEnum.Walk)
	if err != nil {
		return nil, err
	}

	return &Person{
		directionImages:    directionImages,
		behaviorAnimations: behaviorAnimations,
		direction:          consts.DirectionEnum.Down,
		moveStartingFoot:   sprite.FootEnum.Left,
		speed:              1,
		pos:                [2]int{6, 8},
	}, nil
}

func (s *Person) Move() bool {
	return s.pos != s.expectPos
}

func (s *Person) SetPosition(x, y int) {
	s.pos = [2]int{x, y}
	s.expectPos = [2]int{x, y}
}

func (s *Person) OnAction(_ *config.Config, _ input.Action, _ sprite.UpdateInfo) {
	return
}

func (s *Person) PixelPosition(cfg *config.Config) (x, y int) {
	img := s.directionImages[s.direction]
	return s.pos[0] * cfg.TileSize, (s.pos[1]+1)*cfg.TileSize - img.Bounds().Dy()
}

func (s *Person) Update(cfg *config.Config, _ sprite.UpdateInfo) error {
	return nil
}

func (s *Person) Draw(cfg *config.Config, screen *ebiten.Image, ops *ebiten.DrawImageOptions) error {
	img := s.directionImages[s.direction]

	x, y := s.PixelPosition(cfg)
	if ops == nil {
		ops = &ebiten.DrawImageOptions{}
	} else {
		copyOps := *ops
		ops = &copyOps
	}
	ops.GeoM.Translate(float64(x), float64(y))

	if s.Move() {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
		a.Draw(screen, ops)
	} else {
		screen.DrawImage(img, ops)
	}
	return nil
}
