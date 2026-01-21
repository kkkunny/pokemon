package animation

import imgutil "github.com/kkkunny/pokemon/src/util/image"

type Player struct {
	animation     *Animation
	frameCounters []int // 每帧的counter数
	counterTotal  int   // 总counter数

	counter int // 计数器
}

func (p *Player) Animation() *Animation {
	return p.animation
}

func (p *Player) Reset() {
	p.counter = 0
}

// Update @return: 此轮动画是否结束
func (p *Player) Update() bool {
	p.counter++
	if p.counterTotal == 0 {
		for _, fc := range p.frameCounters {
			p.counterTotal += fc
		}
	}
	if p.counter >= p.counterTotal {
		p.counter = 0
	}
	return p.counter == 0
}

// 获取当前帧下标
func (p *Player) GetCurrentFrameIndex() int {
	counter := p.counter
	for i, fc := range p.frameCounters {
		if counter < fc {
			return i
		}
		counter -= fc
	}
	return len(p.frameCounters) - 1
}

func (p *Player) GetCurrentFrame() imgutil.Image {
	return p.animation.frames[p.GetCurrentFrameIndex()].Image
}
