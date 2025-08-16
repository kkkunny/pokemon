package sub_system

import (
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/pokemon/src/output/voice"
)

var PrevPlayingVoicePlayers []*RealTimeVoicePlayer
var PlayingVoicePlayers []*RealTimeVoicePlayer

type RealTimeVoicePlayer struct {
	*voice.Player
}

func NewRealTimeVoicePlayer(player *voice.Player) *RealTimeVoicePlayer {
	return &RealTimeVoicePlayer{Player: player}
}

func (p *RealTimeVoicePlayer) Play() error {
	PlayingVoicePlayers = append(PlayingVoicePlayers, p)
	return nil
}

func (p *RealTimeVoicePlayer) Pause() {
	if stlslices.Contain(PlayingVoicePlayers, p) {
		PlayingVoicePlayers = stlslices.DiffTo(PlayingVoicePlayers, []*RealTimeVoicePlayer{p})
	}
	p.Pause()
}
