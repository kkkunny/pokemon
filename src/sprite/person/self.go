package person

import (
	"errors"
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

type Self struct {
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

func NewSelf(name string) (*Self, error) {
	dirpath := filepath.Join("./resource/map_item/people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}

	directionImages, behaviorAnimations, err := sprite.LoadPersonAnimations(name, sprite.BehaviorEnum.Walk, sprite.BehaviorEnum.Run)
	if err != nil {
		return nil, err
	}

	return &Self{
		directionImages:    directionImages,
		behaviorAnimations: behaviorAnimations,
		direction:          consts.DirectionEnum.Down,
		moveStartingFoot:   sprite.FootEnum.Left,
		speed:              1,
		pos:                [2]int{6, 8},
	}, nil
}

func (s *Self) Move() bool {
	return s.pos != s.expectPos
}

func (s *Self) SetPosition(x, y int) {
	s.pos = [2]int{x, y}
	s.expectPos = [2]int{x, y}
}

func (s *Self) Position() (int, int) {
	return s.pos[0], s.pos[1]
}

func (s *Self) OnAction(cfg *config.Config, action input.Action, info sprite.UpdateInfo) {
	if info == nil {
		return
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return
	}

	if s.Move() {
		return
	}
	preDirection := s.direction
	switch action {
	case input.ActionEnum.MoveUp:
		s.direction = consts.DirectionEnum.Up
	case input.ActionEnum.MoveDown:
		s.direction = consts.DirectionEnum.Down
	case input.ActionEnum.MoveLeft:
		s.direction = consts.DirectionEnum.Left
	case input.ActionEnum.MoveRight:
		s.direction = consts.DirectionEnum.Right
	}
	if preDirection == s.direction {
		expectPos := s.pos
		switch action {
		case input.ActionEnum.MoveUp:
			expectPos = [2]int{s.pos[0], s.pos[1] + int(consts.DirectionEnum.Up)%2}
		case input.ActionEnum.MoveDown:
			expectPos = [2]int{s.pos[0], s.pos[1] + int(consts.DirectionEnum.Down)%2}
		case input.ActionEnum.MoveLeft:
			expectPos = [2]int{s.pos[0] + int(consts.DirectionEnum.Left)%2, s.pos[1]}
		case input.ActionEnum.MoveRight:
			expectPos = [2]int{s.pos[0] + int(consts.DirectionEnum.Right)%2, s.pos[1]}
		}
		if !updateInfo.World.CheckCollision(expectPos[0], expectPos[1]) {
			s.expectPos = expectPos
		}
	}
	return
}

func (s *Self) PixelPosition(cfg *config.Config) (x, y int) {
	img := s.directionImages[s.direction]
	return cfg.ScreenWidth/2 - img.Bounds().Dx()/2, cfg.ScreenHeight/2 - img.Bounds().Dy()/2
}

func (s *Self) Update(cfg *config.Config, info sprite.UpdateInfo) error {
	if info == nil {
		return errors.New("expect UpdateInfo")
	}
	updateInfo, ok := info.(*UpdateInfo)
	if !ok {
		return errors.New("expect UpdateInfo")
	}

	// 移动
	if s.Move() {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
		a.SetFrameTime(cfg.TileSize / s.speed / a.FrameCount())
		a.Update()

		diff := cfg.TileSize - s.moveCounter
		if diff > s.speed {
			s.moveCounter += s.speed
		} else {
			s.moveCounter = 0
			targetMap, targetX, targetY, _ := updateInfo.World.GetActualPosition(s.expectPos[0], s.expectPos[1])
			updateInfo.World.MoveTo(targetMap)
			s.expectPos = [2]int{targetX, targetY}
			s.pos = s.expectPos
			s.moveStartingFoot = -s.moveStartingFoot
			a.Reset()
		}
	}

	// 更新地图位置
	img := s.directionImages[s.direction]
	x, y := s.pos[0]*cfg.TileSize, s.pos[1]*cfg.TileSize+cfg.TileSize-img.Bounds().Dy()
	switch s.direction {
	case consts.DirectionEnum.Up:
		y -= s.moveCounter
	case consts.DirectionEnum.Down:
		y += s.moveCounter
	case consts.DirectionEnum.Left:
		x -= s.moveCounter
	case consts.DirectionEnum.Right:
		x += s.moveCounter
	}
	selfX, selfY := s.PixelPosition(cfg)
	x = -x + selfX
	y = -y + selfY
	updateInfo.World.MovePixelPosTo(x, y)
	return nil
}

func (s *Self) Draw(cfg *config.Config, screen *ebiten.Image, _ *ebiten.DrawImageOptions) error {
	img := s.directionImages[s.direction]

	x, y := s.PixelPosition(cfg)
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(float64(x), float64(y))

	if s.Move() {
		a := s.behaviorAnimations[sprite.BehaviorEnum.Walk][s.direction][s.moveStartingFoot]
		a.Draw(screen, ops)
	} else {
		screen.DrawImage(img, ops)
	}
	return nil
}
