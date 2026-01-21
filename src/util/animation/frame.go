package animation

import (
	"time"

	imgutil "github.com/kkkunny/pokemon/src/util/image"
)

type Frame struct {
	Image imgutil.Image
	Time  time.Duration
}
