package voice

import (
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Player struct {
	ctx     *audio.Context
	path    string
	file    *os.File
	decoder io.Reader
	player  *audio.Player
}

func NewPlayer() *Player {
	return &Player{
		ctx: audio.NewContext(44100),
	}
}

func (p *Player) Close() error {
	if p.path == "" {
		return nil
	}
	err := p.player.Close()
	if err != nil {
		return err
	}
	err = p.file.Close()
	if err != nil {
		return err
	}
	p.path = ""
	return nil
}

func (p *Player) LoadFile(path string) error {
	if p.path == path {
		return nil
	}
	err := p.Close()
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	decoder, err := vorbis.DecodeWithoutResampling(file)
	if err != nil {
		return err
	}
	player, err := p.ctx.NewPlayer(decoder)
	if err != nil {
		return err
	}
	p.file, p.decoder, p.player, p.path = file, decoder, player, path
	return nil
}

func (p *Player) IsPlaying() bool {
	return p.player != nil && p.player.IsPlaying()
}

func (p *Player) Play() error {
	if p.IsPlaying() {
		return nil
	}
	p.ContinuePlay()
	if p.IsPlaying() {
		return nil
	}
	err := p.player.Rewind()
	if err != nil {
		return err
	}
	p.ContinuePlay()
	return nil
}

func (p *Player) ContinuePlay() {
	if p.path == "" {
		return
	}
	p.player.Play()
	return
}
