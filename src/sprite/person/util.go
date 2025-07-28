package person

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kkkunny/pokemon/src/util/image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/animation"
	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
	"github.com/kkkunny/pokemon/src/sprite"
)

type Foot int8

var FootEnum = enum.New[struct {
	Left  Foot `enum:"1"`
	Right Foot `enum:"-1"`
}]()

// 载入人类动画
func loadPersonAnimations(name string, behaviors ...sprite.Behavior) (map[sprite.Behavior]map[consts.Direction]map[Foot]*animation.Animation, error) {
	dirpath := filepath.Join(config.MapItemPath, "people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}

	behaviorAnimations := make(map[sprite.Behavior]map[consts.Direction]map[Foot]*animation.Animation, len(behaviors))
	for _, behavior := range behaviors {
		behaviorImgSheetRect, err := image.NewImageFromFile(filepath.Join(dirpath, string(behavior)+".png"))
		if err != nil {
			return nil, err
		}
		var behaviorDirectionAnimations map[consts.Direction]map[Foot]*animation.Animation
		if behaviorImgSheetRect.Height() == 60 {
			behaviorDirectionAnimations, err = loadSimplePersonDirectionAnimations(behaviorImgSheetRect)
		} else {
			behaviorDirectionAnimations, err = loadCompletePersonDirectionAnimations(behaviorImgSheetRect)
		}
		if err != nil {
			return nil, err
		}
		behaviorAnimations[behavior] = behaviorDirectionAnimations
	}
	return behaviorAnimations, nil
}
func loadSimplePersonDirectionAnimations(imgSheet *image.Image) (map[consts.Direction]map[Foot]*animation.Animation, error) {
	directions := enum.Values[consts.Direction](consts.DirectionEnum)
	directionAnimations := make(map[consts.Direction]map[Foot]*animation.Animation, len(directions))
	frameW, frameH := imgSheet.Width()/3, imgSheet.Height()/3
	for i, direction := range []consts.Direction{consts.DirectionEnum.Down, consts.DirectionEnum.Up, consts.DirectionEnum.Left} {
		y := i * frameH
		leftFootAnimationFrameSheet := image.NewImage(2*frameW, frameH)
		rightFootAnimationFrameSheet := image.NewImage(2*frameW, frameH)
		for j := range 3 {
			x := j * frameW
			img := imgSheet.SubImage(x, y, frameW, frameH)
			switch j {
			case 0:
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(0, 0)
				leftFootAnimationFrameSheet.DrawImage(img, ops)
				ops.GeoM.Translate(0, 0)
				rightFootAnimationFrameSheet.DrawImage(img, ops)
			case 1:
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(frameW), 0)
				leftFootAnimationFrameSheet.DrawImage(img, ops)
			case 2:
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(frameW), 0)
				rightFootAnimationFrameSheet.DrawImage(img, ops)
			}
		}
		directionAnimations[direction] = map[Foot]*animation.Animation{
			FootEnum.Left:  animation.NewAnimation(leftFootAnimationFrameSheet, frameW, frameH, 0),
			FootEnum.Right: animation.NewAnimation(rightFootAnimationFrameSheet, frameW, frameH, 0),
		}
	}

	left := directionAnimations[consts.DirectionEnum.Left][FootEnum.Left].GetFrameImage(0)
	right := image.NewImage(frameW, frameH)
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Scale(-1, 1)
	ops.GeoM.Translate(float64(frameW), 0)
	right.DrawImage(left, ops)

	leftFootAnimationFrameSheet := image.NewImage(2*frameW, frameH)
	ops = &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(0, 0)
	leftFootAnimationFrameSheet.DrawImage(right, ops)
	ops.GeoM.Scale(-1, 1)
	ops.GeoM.Translate(float64(frameW)*2, 0)
	leftFootAnimationFrameSheet.DrawImage(directionAnimations[consts.DirectionEnum.Left][FootEnum.Right].GetFrameImage(1), ops)

	rightFootAnimationFrameSheet := image.NewImage(2*frameW, frameH)
	ops = &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(0, 0)
	rightFootAnimationFrameSheet.DrawImage(right, ops)
	ops.GeoM.Scale(-1, 1)
	ops.GeoM.Translate(float64(frameW)*2, 0)
	rightFootAnimationFrameSheet.DrawImage(directionAnimations[consts.DirectionEnum.Left][FootEnum.Left].GetFrameImage(1), ops)

	directionAnimations[consts.DirectionEnum.Right] = map[Foot]*animation.Animation{
		FootEnum.Left:  animation.NewAnimation(leftFootAnimationFrameSheet, frameW, frameH, 0),
		FootEnum.Right: animation.NewAnimation(rightFootAnimationFrameSheet, frameW, frameH, 0),
	}
	return directionAnimations, nil
}
func loadCompletePersonDirectionAnimations(imgSheet *image.Image) (map[consts.Direction]map[Foot]*animation.Animation, error) {
	directions := enum.Values[consts.Direction](consts.DirectionEnum)
	directionAnimations := make(map[consts.Direction]map[Foot]*animation.Animation, len(directions))
	frameW, frameH := imgSheet.Width()/3, imgSheet.Height()/3
	for i, direction := range []consts.Direction{consts.DirectionEnum.Down, consts.DirectionEnum.Up, consts.DirectionEnum.Left, consts.DirectionEnum.Right} {
		y := i * frameH
		leftFootAnimationFrameSheet := image.NewImage(2*frameW, frameH)
		rightFootAnimationFrameSheet := image.NewImage(2*frameW, frameH)
		for j := range 3 {
			x := j * frameW
			img := imgSheet.SubImage(x, y, frameW, frameH)
			switch j {
			case 0:
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(0, 0)
				leftFootAnimationFrameSheet.DrawImage(img, ops)
				ops.GeoM.Translate(0, 0)
				rightFootAnimationFrameSheet.DrawImage(img, ops)
			case 1:
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(frameW), 0)
				leftFootAnimationFrameSheet.DrawImage(img, ops)
			case 2:
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(frameW), 0)
				rightFootAnimationFrameSheet.DrawImage(img, ops)
			}
		}
		directionAnimations[direction] = map[Foot]*animation.Animation{
			FootEnum.Left:  animation.NewAnimation(leftFootAnimationFrameSheet, frameW, frameH, 0),
			FootEnum.Right: animation.NewAnimation(rightFootAnimationFrameSheet, frameW, frameH, 0),
		}
	}
	return directionAnimations, nil
}

// GetNextPositionByDirection 获取该方向下一步位置
func GetNextPositionByDirection(d consts.Direction, x, y int) (int, int) {
	switch d {
	case consts.DirectionEnum.Up:
		return x, y + int(consts.DirectionEnum.Up)%2
	case consts.DirectionEnum.Down:
		return x, y + int(consts.DirectionEnum.Down)%2
	case consts.DirectionEnum.Left:
		return x + int(consts.DirectionEnum.Left)%2, y
	case consts.DirectionEnum.Right:
		return x + int(consts.DirectionEnum.Right)%2, y
	default:
		return x, y
	}
}
