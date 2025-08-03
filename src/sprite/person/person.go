package person

import (
	"errors"
	"math/rand/v2"

	stlmaps "github.com/kkkunny/stl/container/maps"
	stlval "github.com/kkkunny/stl/value"
	"github.com/lafriks/go-tiled"

	"github.com/kkkunny/pokemon/src/animation"
	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/context"
	"github.com/kkkunny/pokemon/src/input"
	"github.com/kkkunny/pokemon/src/sprite"
	"github.com/kkkunny/pokemon/src/sprite/item"
	"github.com/kkkunny/pokemon/src/util/draw"
)

func init() {
	sprite.RegisterCreateFunc([]string{"person"}, func(object *tiled.Object) (sprite.Sprite, error) {
		person, err := NewPersonByTile(object)
		if err != nil {
			return nil, err
		}
		return person, nil
	})
}

type Person interface {
	sprite.MovableSprite
}

type _Person struct {
	item.Item
	// 静态资源
	behaviorAnimations map[sprite.Behavior]map[consts.Direction]map[Foot]*animation.Animation // 行为动画
	// 属性
	direction consts.Direction // 当前所处方向
	// 移动
	movable           bool             // 是否可移动
	speed             int              // 移动速度
	moveStartingFoot  Foot             // 移动时的起始脚
	nextStepDirection consts.Direction // 下一步预期所处方向
	pos               [2]int           // 当前地块位置
	nextStepPos       [2]int           // 下一步预期所处的地块位置，用于移动
	moveCounter       int              // 移动时的计数器，用于显示动画
}

func NewPerson(name string) (Person, error) {
	behaviorAnimations, err := loadPersonAnimations(name, sprite.BehaviorEnum.Walk)
	if err != nil {
		return nil, err
	}

	itemSprite, err := item.NewItem()
	if err != nil {
		return nil, err
	}

	return &_Person{
		Item:               itemSprite,
		behaviorAnimations: behaviorAnimations,
		direction:          consts.DirectionEnum.Down,
		movable:            true,
		nextStepDirection:  consts.DirectionEnum.Down,
		moveStartingFoot:   FootEnum.Left,
		speed:              1,
	}, nil
}

func NewPersonByTile(object *tiled.Object) (Person, error) {
	imgName := object.Properties.GetString("image")
	behaviorAnimations, err := loadPersonAnimations(imgName, sprite.BehaviorEnum.Walk)
	if err != nil {
		return nil, err
	}

	itemSprite, err := item.NewItemByTile(object)
	if err != nil {
		return nil, err
	}

	return &_Person{
		Item:               itemSprite,
		behaviorAnimations: behaviorAnimations,
		direction:          consts.DirectionEnum.Down,
		movable:            true,
		nextStepDirection:  consts.DirectionEnum.Down,
		moveStartingFoot:   FootEnum.Left,
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

func (p *_Person) Direction() consts.Direction {
	return p.direction
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
	if !p.Turn(d) {
		return false
	}
	p.nextStepPos[0], p.nextStepPos[1] = GetNextPositionByDirection(d, p.pos[0], p.pos[1])
	return true
}

func (p *_Person) OnAction(_ context.Context, _ input.KeyInputAction, _ sprite.UpdateInfo) error {
	return nil
}

func (p *_Person) PixelPosition() (x, y float64) {
	width := stlmaps.First(stlmaps.First(p.behaviorAnimations[sprite.BehaviorEnum.Walk]).E2()).E2().GetFrameImage(0).Height()
	x, y = float64(p.pos[0]*config.TileSize), float64((p.pos[1]+1)*config.TileSize-width)

	if p.Moving() && !p.Turning() {
		switch p.nextStepDirection {
		case consts.DirectionEnum.Up:
			y -= float64(p.moveCounter)
		case consts.DirectionEnum.Down:
			y += float64(p.moveCounter)
		case consts.DirectionEnum.Left:
			x -= float64(p.moveCounter)
		case consts.DirectionEnum.Right:
			x += float64(p.moveCounter)
		}
	}

	return x, y
}

func (p *_Person) Update(ctx context.Context, info sprite.UpdateInfo) error {
	if info == nil {
		return errors.New("expect UpdateInfo")
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return errors.New("expect UpdateInfo")
	}

	if p.Turning() {
		if p.moveCounter < config.TileSize {
			p.moveCounter += 2
		} else {
			p.moveCounter = 0
			p.direction = p.nextStepDirection
			p.moveStartingFoot = -p.moveStartingFoot
		}
	} else if p.Moving() {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.nextStepDirection][p.moveStartingFoot]
		a.SetFrameTime(config.TileSize / p.speed / a.FrameCount())
		a.Update()

		diff := config.TileSize - p.moveCounter
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
	} else if p.movable {
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
		} else if x, y := GetNextPositionByDirection(nextStepDirection, p.pos[0], p.pos[1]); !updateInfo.World.CheckCollision(p.direction, x, y) {
			p.SetNextStepDirection(nextStepDirection)
		}
	}
	return nil
}

func (p *_Person) Draw(ctx context.Context, drawer draw.Drawer) error {
	x, y := p.PixelPosition()
	drawer = drawer.Move(x, y)

	if p.Turning() {
		if p.direction == -p.nextStepDirection {
			p.moveStartingFoot = FootEnum.Right
		} else if p.direction == consts.DirectionEnum.Up {
			p.moveStartingFoot = stlval.Ternary(p.nextStepDirection == consts.DirectionEnum.Left, FootEnum.Left, FootEnum.Right)
		} else if p.direction == consts.DirectionEnum.Down {
			p.moveStartingFoot = stlval.Ternary(p.nextStepDirection == consts.DirectionEnum.Right, FootEnum.Left, FootEnum.Right)
		} else if p.direction == consts.DirectionEnum.Left {
			p.moveStartingFoot = stlval.Ternary(p.nextStepDirection == consts.DirectionEnum.Down, FootEnum.Left, FootEnum.Right)
		} else if p.direction == consts.DirectionEnum.Right {
			p.moveStartingFoot = stlval.Ternary(p.nextStepDirection == consts.DirectionEnum.Up, FootEnum.Left, FootEnum.Right)
		}
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.nextStepDirection][p.moveStartingFoot]
		return drawer.DrawImage(a.GetFrameImage(1))
	} else {
		a := p.behaviorAnimations[sprite.BehaviorEnum.Walk][p.nextStepDirection][p.moveStartingFoot]
		return a.Draw(drawer)
	}
}

func (p *_Person) Turn(d consts.Direction) bool {
	if p.Busying() {
		return false
	}
	p.nextStepDirection = d
	return true
}

func (p *_Person) SetMovable(movable bool) {
	p.movable = movable
}

func (p *_Person) Movable() bool {
	return p.movable
}

func (p *_Person) Collision() bool {
	return true
}

func (p *_Person) CollisionPosition() (int, int) {
	return p.nextStepPos[0], p.nextStepPos[1]
}
