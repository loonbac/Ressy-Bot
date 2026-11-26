package music

import (
	"errors"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

var (
	// ErrQueueEmpty means a guild has no track to skip.
	ErrQueueEmpty = errors.New("queue is empty")
	// ErrVoiceStateMissing means the bot is not in a voice channel.
	ErrVoiceStateMissing = errors.New("voice state is missing")
	// ErrLavalinkSessionMissing means the audio node has not sent its ready event.
	ErrLavalinkSessionMissing = errors.New("lavalink session is missing")
)

// Track is one playable audio item.
type Track struct {
	Encoded string
	Title   string
	URI     string
}

// VoiceState is the Discord data Lavalink needs to join a voice channel.
type VoiceState = core.VoiceState

// VoiceLookup finds the bot voice state for a guild.
type VoiceLookup = core.VoiceLookup

// AudioNode resolves and plays tracks without exposing its transport.
type AudioNode interface {
	Resolve(query string) ([]Track, error)
	Play(guildID string, track Track, voice VoiceState) error
	Stop(guildID string) error
	NowPlaying(guildID string) (Track, bool)
	Close() error
}
