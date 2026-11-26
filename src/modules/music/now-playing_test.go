package music

import (
	"strings"
	"testing"
)

func TestNowPlayingReportsCurrentTrack(t *testing.T) {
	queues := newGuildQueues()
	queues.Enqueue("g1", Track{Title: "current"})
	queues.Next("g1")

	reply := runCommand(NewNowPlayingCommand(queues), "g1")

	if !strings.Contains(reply, "current") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestNowPlayingReportsIdleState(t *testing.T) {
	reply := runCommand(NewNowPlayingCommand(newGuildQueues()), "g1")
	if !strings.Contains(reply, "Nothing is playing") {
		t.Fatalf("reply = %q", reply)
	}
}
