package person

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	stlmaps "github.com/kkkunny/stl/container/maps"

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
	Moving() bool
	PixelPosition(cfg *config.Config) (x, y int)
}

type _Person struct {
	// 静态资源
	behaviorAnimations map[sprite.Behavior]map[consts.Direction]map[sprite.Foot]*animation.Animation // 行为动画
	// 属性
	direction consts.Direction // 当前所处方向
	// 移动
	speed             int              // 移动速度
	moveStartingFoot  sprite.Foot      // 移动时的起始脚
	nextStepDirection consts.Direction // 下一步预期所处方向
	pos               [2]int           // 当前地块位置
	nextStepPos       [2]int           // 下一步预期所处的地块位置，用于移动
	moveCounter       int              // 移动时的计数器，用于显示动画
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

	return &_Person{
		behaviorAnimations: behaviorAnimations,
		direction:          consts.DirectionEnum.Down,
		nextStepDirection:  consts.DirectionEnum.Down,
		moveStartingFoot:   sprite.FootEnum.Left,
		speed:              1,
	}, nil
}

// Turning 是否正在转向
func (p *_Person) Turning() bool {
	return p.direction != p.nextStepDirection
}

// Moving 是否正在移动
func (p *_Person) Moving() bool {
	return p.pos != p.nextStepPos
}

// Busying 是否忙碌中
func (p *_Person) Busying() bool {
	return p.Turning() || p.Moving()
}

func (p *_Person) SetDirection(d consts.Direction) {
	if p.Busying() {
		return
	}
	p.direction = d
	p.nextStepDirection = d
}

func (p *_Person) SetPosition(x, y int) {
	if p.Busying() {
		return
	}
	p.pos = [2]int{x, y}
	p.nextStepPos = [2]int{x, y}
}

func (p *_Person) Position() (int, int) {
	return p.pos[0], p.pos[1]
}

// SetNextStepDirection 设置下一步方向，每次只可前进一格
// 设置时不会校验下一个是否可移动，会在Update时校验
func (p *_Person) SetNextStepDirection(d consts.Direction) bool {
	if p.Busying() {
		return false
	}
	p.nextStepDirection = d
	switch d {
	case consts.DirectionEnum.Up:
		p.nextStepPos = [2]int{p.pos[0], p.pos[1] + int(consts.DirectionEnum.Up)%2}
	case consts.DirectionEnum.Down:
		p.nextStepPos = [2]int{p.pos[0], p.pos[1] + int(consts.DirectionEnum.Down)%2}
	case consts.DirectionEnum.Left:
		p.nextStepPos = [2]int{p.pos[0] + int(consts.DirectionEnum.Left)%2, p.pos[1]}
	case consts.DirectionEnum.Right:
		p.nextStepPos = [2]int{p.pos[0] + int(consts.DirectionEnum.Right)%2, p.pos[1]}
	}
	return true
}

func (p *_Person) OnAction(_ *config.Config, _ input.Action, _ sprite.UpdateInfo) {
	return
}

func (p *_Person) PixelPosition(cfg *config.Config) (x, y int) {
	width := stlmaps.First(stlmaps.First(p.behaviorAnimations[sprite.BehaviorEnum.Walk]).E2()).E2().GetFrameImage(0).Bounds().Dy()
	x, y = p.pos[0]*cfg.TileSize, (p.pos[1]+1)*cfg.TileSize-width

	if p.Moving() && !p.Turning() {
		switch p.nextStepDirection {
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

	if p.Turning() {
		if p.moveCounter < cfg.TileSize {
			p.moveCounter += 2
		} else {
			p.moveCounter = 0
			p.direction = p.nextStepDirection
			p.moveStartingFoot = -p.moveStartingFoot
		}
	} else if p.Moving() {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.nextStepDirection][p.moveStartingFoot]
		a.SetFrameTime(cfg.TileSize / p.speed / a.FrameCount())
		a.Update()

		diff := cfg.TileSize - p.moveCounter
		if diff > p.speed {
			p.moveCounter += p.speed
		} else {
			p.moveCounter = 0
			_, targetX, targetY, _ := updateInfo.World.GetActualPosition(p.nextStepPos[0], p.nextStepPos[1])
			p.nextStepPos = [2]int{targetX, targetY}
			p.pos = p.nextStepPos
			p.moveStartingFoot = -p.moveStartingFoot
			a.Reset()
		}
	} else {
		var nextStepDirection consts.Direction
		n := rand.IntN(500)
		if n >= 499 {
			nextStepDirection = consts.DirectionEnum.Up
		} else if n >= 498 {
			nextStepDirection = consts.DirectionEnum.Down
		} else if n >= 497 {
			nextStepDirection = consts.DirectionEnum.Left
		} else if n >= 496 {
			nextStepDirection = consts.DirectionEnum.Right
		} else {
			return nil
		}
		if p.direction != nextStepDirection && rand.IntN(500) > 250 {
			p.nextStepDirection = nextStepDirection
		} else if p.SetNextStepDirection(nextStepDirection) && updateInfo.World.CheckCollision(p.nextStepPos[0], p.nextStepPos[1]) {
			p.nextStepPos = p.pos
		}
	}
	return nil
}

func (p *_Person) Draw(cfg *config.Config, screen *ebiten.Image, ops ebiten.DrawImageOptions) error {
	x, y := p.PixelPosition(cfg)
	ops.GeoM.Translate(float64(x), float64(y))

	if p.Turning() {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.nextStepDirection][p.moveStartingFoot]
		screen.DrawImage(a.GetFrameImage(1), &ops)
	} else {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.nextStepDirection][p.moveStartingFoot]
		a.Draw(screen, ops)
	}
	return nil
}
