package person

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/system/world/sprite"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/animation"
	"github.com/kkkunny/pokemon/src/util/image"

	"github.com/kkkunny/pokemon/src/config"
)

type Foot int8

var FootEnum = enum.New[struct {
	Left  Foot `enum:"1"`
	Right Foot `enum:"-1"`
}]()

// 载入人类动画
func loadPersonAnimations(name string, behaviors ...sprite.Behavior) (map[sprite.Behavior]map[util.Direction]map[Foot]*animation.Animation, error) {
	dirpath := filepath.Join(config.GFXMapPath, "people", name)
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("can not found trainer `%s`", name)
	}

	behaviorAnimations := make(map[sprite.Behavior]map[util.Direction]map[Foot]*animation.Animation, len(behaviors))
	for _, behavior := range behaviors {
		behaviorImgSheetRect, err := imgutil.NewImageFromFile(filepath.Join(dirpath, string(behavior)+".png"))
		if err != nil {
			return nil, err
		}
		var behaviorDirectionAnimations map[util.Direction]map[Foot]*animation.Animation
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
func loadSimplePersonDirectionAnimations(imgSheet *imgutil.Image) (map[util.Direction]map[Foot]*animation.Animation, error) {
	directions := enum.Values[util.Direction](util.DirectionEnum)
	directionAnimations := make(map[util.Direction]map[Foot]*animation.Animation, len(directions))
	frameW, frameH := imgSheet.Width()/3, imgSheet.Height()/3
	for i, direction := range []util.Direction{util.DirectionEnum.Down, util.DirectionEnum.Up, util.DirectionEnum.Left} {
		y := i * frameH
		leftFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
		rightFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
		for j := range 3 {
			x := j * frameW
			img := imgSheet.SubImage(image.Rect(x, y, x+frameW, y+frameH))
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

	left := directionAnimations[util.DirectionEnum.Left][FootEnum.Left].GetFrameImage(0)
	right := imgutil.NewImage(frameW, frameH)
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Scale(-1, 1)
	ops.GeoM.Translate(float64(frameW), 0)
	right.DrawImage(left, ops)

	leftFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
	ops = &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(0, 0)
	leftFootAnimationFrameSheet.DrawImage(right, ops)
	ops.GeoM.Scale(-1, 1)
	ops.GeoM.Translate(float64(frameW)*2, 0)
	leftFootAnimationFrameSheet.DrawImage(directionAnimations[util.DirectionEnum.Left][FootEnum.Right].GetFrameImage(1), ops)

	rightFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
	ops = &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(0, 0)
	rightFootAnimationFrameSheet.DrawImage(right, ops)
	ops.GeoM.Scale(-1, 1)
	ops.GeoM.Translate(float64(frameW)*2, 0)
	rightFootAnimationFrameSheet.DrawImage(directionAnimations[util.DirectionEnum.Left][FootEnum.Left].GetFrameImage(1), ops)

	directionAnimations[util.DirectionEnum.Right] = map[Foot]*animation.Animation{
		FootEnum.Left:  animation.NewAnimation(leftFootAnimationFrameSheet, frameW, frameH, 0),
		FootEnum.Right: animation.NewAnimation(rightFootAnimationFrameSheet, frameW, frameH, 0),
	}
	return directionAnimations, nil
}
func loadCompletePersonDirectionAnimations(imgSheet *imgutil.Image) (map[util.Direction]map[Foot]*animation.Animation, error) {
	directions := enum.Values[util.Direction](util.DirectionEnum)
	directionAnimations := make(map[util.Direction]map[Foot]*animation.Animation, len(directions))
	frameW, frameH := imgSheet.Width()/3, imgSheet.Height()/3
	for i, direction := range []util.Direction{util.DirectionEnum.Down, util.DirectionEnum.Up, util.DirectionEnum.Left, util.DirectionEnum.Right} {
		y := i * frameH
		leftFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
		rightFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
		for j := range 3 {
			x := j * frameW
			img := imgSheet.SubImage(image.Rect(x, y, x+frameW, y+frameH))
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
func GetNextPositionByDirection(d util.Direction, x, y int) (int, int) {
	switch d {
	case util.DirectionEnum.Up:
		return x, y + int(util.DirectionEnum.Up)%2
	case util.DirectionEnum.Down:
		return x, y + int(util.DirectionEnum.Down)%2
	case util.DirectionEnum.Left:
		return x + int(util.DirectionEnum.Left)%2, y
	case util.DirectionEnum.Right:
		return x + int(util.DirectionEnum.Right)%2, y
	default:
		return x, y
	}
}
