package person

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/tnnmigga/enum"

	"github.com/kkkunny/pokemon/src/system/world/sprite"
	"github.com/kkkunny/pokemon/src/util"
	"github.com/kkkunny/pokemon/src/util/animation"
	"github.com/kkkunny/pokemon/src/util/draw"
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
		if behaviorImgSheetRect.Bounds().Dy() == 60 {
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
func loadSimplePersonDirectionAnimations(imgSheet imgutil.Image) (map[util.Direction]map[Foot]*animation.Animation, error) {
	directions := enum.Values[util.Direction](util.DirectionEnum)
	directionAnimations := make(map[util.Direction]map[Foot]*animation.Animation, len(directions))
	frameW, frameH := imgSheet.Bounds().Dx()/3, imgSheet.Bounds().Dy()/3
	for i, direction := range []util.Direction{util.DirectionEnum.Down, util.DirectionEnum.Up, util.DirectionEnum.Left} {
		y := i * frameH
		leftFootAnimation := animation.NewAnimation(nil, 0)
		rightFootAnimation := animation.NewAnimation(nil, 0)
		for j := range 3 {
			x := j * frameW
			img := imgSheet.SubImage(image.Rect(x, y, x+frameW, y+frameH))
			switch j {
			case 0:
				leftFootAnimation.AddFrame(img)
				rightFootAnimation.AddFrame(img)
			case 1:
				leftFootAnimation.AddFrame(img)
			case 2:
				rightFootAnimation.AddFrame(img)
			}
		}
		directionAnimations[direction] = map[Foot]*animation.Animation{
			FootEnum.Left:  leftFootAnimation,
			FootEnum.Right: rightFootAnimation,
		}
	}

	right := directionAnimations[util.DirectionEnum.Left][FootEnum.Left].GetFrameImage(0).Scale(-1, 1)

	leftFootAnimation := animation.NewAnimation(nil, 0)
	leftFootAnimation.AddFrame(right)
	leftFootAnimation.AddFrame(directionAnimations[util.DirectionEnum.Left][FootEnum.Right].GetFrameImage(1).Scale(-1, 1))

	rightFootAnimation := animation.NewAnimation(nil, 0)
	rightFootAnimation.AddFrame(right)
	rightFootAnimation.AddFrame(directionAnimations[util.DirectionEnum.Left][FootEnum.Left].GetFrameImage(1).Scale(-1, 1))

	directionAnimations[util.DirectionEnum.Right] = map[Foot]*animation.Animation{
		FootEnum.Left:  leftFootAnimation,
		FootEnum.Right: rightFootAnimation,
	}
	return directionAnimations, nil
}
func loadCompletePersonDirectionAnimations(imgSheet imgutil.Image) (map[util.Direction]map[Foot]*animation.Animation, error) {
	directions := enum.Values[util.Direction](util.DirectionEnum)
	directionAnimations := make(map[util.Direction]map[Foot]*animation.Animation, len(directions))
	frameW, frameH := imgSheet.Bounds().Dx()/3, imgSheet.Bounds().Dy()/3
	for i, _ := range []util.Direction{util.DirectionEnum.Down, util.DirectionEnum.Up, util.DirectionEnum.Left, util.DirectionEnum.Right} {
		y := i * frameH
		leftFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
		rightFootAnimationFrameSheet := imgutil.NewImage(2*frameW, frameH)
		for j := range 3 {
			x := j * frameW
			img := imgSheet.SubImage(image.Rect(x, y, x+frameW, y+frameH))
			switch j {
			case 0:
				draw.PrepareDrawImage(leftFootAnimationFrameSheet, img).Draw()
				draw.PrepareDrawImage(rightFootAnimationFrameSheet, img).Draw()
			case 1:
				draw.PrepareDrawImage(leftFootAnimationFrameSheet, img).Move(frameW, 0).Draw()
			case 2:
				draw.PrepareDrawImage(rightFootAnimationFrameSheet, img).Move(frameW, 0).Draw()
			}
		}
		// directionAnimations[direction] = map[Foot]*animation.Animation{
		// 	FootEnum.Left:  animation.NewAnimation(leftFootAnimationFrameSheet, frameW, frameH, 0),
		// 	FootEnum.Right: animation.NewAnimation(rightFootAnimationFrameSheet, frameW, frameH, 0),
		// }
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
