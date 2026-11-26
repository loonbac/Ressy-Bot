package music

import (
	"strings"
	"testing"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type fakeAudioNode struct {
	tracks  []Track
	played  map[string]Track
	stopped []string
}

func (n *fakeAudioNode) Resolve(string) ([]Track, error) { return n.tracks, nil }
func (n *fakeAudioNode) Play(guildID string, track Track, _ core.VoiceState) error {
	if n.played == nil {
		n.played = make(map[string]Track)
	}
	n.played[guildID] = track
	return nil
}
func (n *fakeAudioNode) Stop(guildID string) error {
	n.stopped = append(n.stopped, guildID)
	return nil
}
func (n *fakeAudioNode) NowPlaying(string) (Track, bool) { return Track{}, false }
func (n *fakeAudioNode) Close() error                    { return nil }

func runCommand(command core.Command, guildID string, args ...string) string {
	reply := ""
	_ = command.Execute(core.CommandContext{
		GuildID: guildID, Args: args,
		Reply: func(message string) error { reply = message; return nil },
	})
	return reply
}

func TestPlayQueuesAndPlaysResolvedTrack(t *testing.T) {
	node := &fakeAudioNode{tracks: []Track{{Encoded: "encoded", Title: "Song"}}}
	queues := newGuildQueues()
	lookup := func(string) (core.VoiceState, bool) { return core.VoiceState{ChannelID: "c1"}, true }

	reply := runCommand(NewPlayCommand(node, queues, lookup), "g1", "song")

	if !strings.Contains(reply, "Song") {
		t.Fatalf("reply = %q", reply)
	}
	if len(queues.List("g1")) != 1 || queues.List("g1")[0].Title != "Song" {
		t.Fatalf("queue = %+v", queues.List("g1"))
	}
	if played, ok := node.played["g1"]; !ok || played.Title != "Song" {
		t.Fatalf("played = %+v", node.played)
	}
}

func TestPlayRequiresVoiceState(t *testing.T) {
	node := &fakeAudioNode{tracks: []Track{{Title: "Song"}}}
	lookup := func(string) (core.VoiceState, bool) { return core.VoiceState{}, false }

	reply := runCommand(NewPlayCommand(node, newGuildQueues(), lookup), "g1", "song")

	if !strings.Contains(reply, "voice channel") {
		t.Fatalf("reply = %q", reply)
	}
}
