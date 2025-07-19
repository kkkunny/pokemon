package sprite

import "image"

type Sprite interface {
	Update() error
	Image() (image.Image, error)
}
