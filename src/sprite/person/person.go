package person

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	stlmaps "github.com/kkkunny/stl/container/maps"
	"github.com/tnnmigga/enum"

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

type Person interface {
	sprite.Sprite
	Move() bool
	PixelPosition(cfg *config.Config) (x, y int)
}

type _Person struct {
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

func NewPerson(name string) (Person, error) {
	dirpath := filepath.Join("./resource/map_item/people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}

	behaviorAnimations, err := sprite.LoadPersonAnimations(name, sprite.BehaviorEnum.Walk)
	if err != nil {
		return nil, err
	}
	directionImages := make(map[consts.Direction]*ebiten.Image, len(enum.Values[consts.Direction](consts.DirectionEnum)))
	for direction, directionAnimations := range behaviorAnimations[sprite.BehaviorEnum.Walk] {
		directionImages[direction] = stlmaps.First(directionAnimations).E2().GetFrameImage(0)
	}

	return &_Person{
		directionImages:    directionImages,
		behaviorAnimations: behaviorAnimations,
		direction:          consts.DirectionEnum.Down,
		moveStartingFoot:   sprite.FootEnum.Left,
		speed:              1,
	}, nil
}

func (p *_Person) Move() bool {
	return p.pos != p.expectPos
}

func (p *_Person) SetPosition(x, y int) {
	p.pos = [2]int{x, y}
	p.expectPos = [2]int{x, y}
}

func (p *_Person) Position() (int, int) {
	return p.pos[0], p.pos[1]
}

func (p *_Person) OnAction(_ *config.Config, _ input.Action, _ sprite.UpdateInfo) {
	return
}

func (p *_Person) PixelPosition(cfg *config.Config) (x, y int) {
	img := p.directionImages[p.direction]
	x, y = p.pos[0]*cfg.TileSize, (p.pos[1]+1)*cfg.TileSize-img.Bounds().Dy()

	if p.Move() {
		switch p.direction {
		case consts.DirectionEnum.Up:
			y -= p.moveCounter
		case consts.DirectionEnum.Down:
			y += p.moveCounter
		case consts.DirectionEnum.Left:
			x -= p.moveCounter
		case consts.DirectionEnum.Right:
			x += p.moveCounter
		}
	}

	return x, y
}

func (p *_Person) Update(cfg *config.Config, info sprite.UpdateInfo) error {
	if info == nil {
		return errors.New("expect UpdateInfo")
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return errors.New("expect UpdateInfo")
	}

	// 移动
	if p.Move() {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.direction][p.moveStartingFoot]
		a.SetFrameTime(cfg.TileSize / p.speed / a.FrameCount())
		a.Update()

		diff := cfg.TileSize - p.moveCounter
		if diff > p.speed {
			p.moveCounter += p.speed
		} else {
			p.moveCounter = 0
			_, targetX, targetY, _ := updateInfo.World.GetActualPosition(p.expectPos[0], p.expectPos[1])
			p.expectPos = [2]int{targetX, targetY}
			p.pos = p.expectPos
			p.moveStartingFoot = -p.moveStartingFoot
			a.Reset()
		}
	} else {
		preDirection := p.direction
		n := rand.IntN(500)
		if n >= 499 {
			p.direction = consts.DirectionEnum.Up
		} else if n >= 498 {
			p.direction = consts.DirectionEnum.Down
		} else if n >= 497 {
			p.direction = consts.DirectionEnum.Left
		} else if n >= 496 {
			p.direction = consts.DirectionEnum.Right
		} else {
			return nil
		}
		if preDirection == p.direction {
			expectPos := p.pos
			switch p.direction {
			case consts.DirectionEnum.Up:
				expectPos = [2]int{p.pos[0], p.pos[1] + int(consts.DirectionEnum.Up)%2}
			case consts.DirectionEnum.Down:
				expectPos = [2]int{p.pos[0], p.pos[1] + int(consts.DirectionEnum.Down)%2}
			case consts.DirectionEnum.Left:
				expectPos = [2]int{p.pos[0] + int(consts.DirectionEnum.Left)%2, p.pos[1]}
			case consts.DirectionEnum.Right:
				expectPos = [2]int{p.pos[0] + int(consts.DirectionEnum.Right)%2, p.pos[1]}
			}
			if !updateInfo.World.CheckCollision(expectPos[0], expectPos[1]) {
				p.expectPos = expectPos
			}
		}
	}
	return nil
}

func (p *_Person) Draw(cfg *config.Config, screen *ebiten.Image, ops ebiten.DrawImageOptions) error {
	img := p.directionImages[p.direction]

	x, y := p.PixelPosition(cfg)
	ops.GeoM.Translate(float64(x), float64(y))

	if p.Move() {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.direction][p.moveStartingFoot]
		a.Draw(screen, ops)
	} else {
		screen.DrawImage(img, &ops)
	}
	return nil
}
