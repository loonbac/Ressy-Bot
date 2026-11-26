package music

import (
	"strings"
	"testing"
)

func TestSkipAdvancesQueue(t *testing.T) {
	node := &fakeAudioNode{}
	queues := newGuildQueues()
	queues.Enqueue("g1", Track{Title: "first"})
	queues.Enqueue("g1", Track{Title: "second"})

	reply := runCommand(NewSkipCommand(node, queues, foundVoice), "g1")

	if !strings.Contains(reply, "first") {
		t.Fatalf("reply = %q", reply)
	}
	if len(node.stopped) != 1 || node.played["g1"].Title != "first" {
		t.Fatalf("stopped=%v played=%+v", node.stopped, node.played)
	}
}

func TestSkipReportsEmptyQueue(t *testing.T) {
	reply := runCommand(NewSkipCommand(&fakeAudioNode{}, newGuildQueues(), foundVoice), "g1")
	if !strings.Contains(reply, "empty") {
		t.Fatalf("reply = %q", reply)
	}
}
