package pokemon

import (
	randv1 "math/rand"
	"math/rand/v2"
	"time"
)

type Pokemon struct {
	Race       Race  // 种族
	Level      uint8 // 级别
	Gender     bool  // 性别（true=雄，false=雌）
	Experience int   // 经验值
}

func (r *Race) RandomPokemon() *Pokemon {
	rander := rand.New(randv1.New(randv1.NewSource(time.Now().UnixNano())))
	gender := float64(rander.IntN(100)) < r.MaleRate
	return &Pokemon{
		Race:   *r,
		Level:  uint8(rander.IntN(100)) + 1,
		Gender: gender,
	}
}
