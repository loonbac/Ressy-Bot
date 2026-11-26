package music

import (
	"strings"
	"testing"
)

func TestNowPlayingReportsCurrentTrack(t *testing.T) {
	queues := newGuildQueues()
	queues.Enqueue("g1", Track{Title: "stale"})
	queues.Next("g1")
	node := &fakeAudioNode{now: Track{Title: "current"}}

	reply := runCommand(NewNowPlayingCommand(node, queues), "g1")

	if !strings.Contains(reply, "current") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestNowPlayingReportsIdleState(t *testing.T) {
	reply := runCommand(NewNowPlayingCommand(&fakeAudioNode{}, newGuildQueues()), "g1")
	if !strings.Contains(reply, "Nothing is playing") {
		t.Fatalf("reply = %q", reply)
	}
}
