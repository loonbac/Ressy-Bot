package music

import (
	"strings"
	"testing"
)

func TestQueueListsWaitingTracks(t *testing.T) {
	queues := newGuildQueues()
	queues.Enqueue("g1", Track{Title: "first"})
	queues.Enqueue("g1", Track{Title: "second"})

	reply := runCommand(NewQueueCommand(queues), "g1")

	if !strings.Contains(reply, "1. first") || !strings.Contains(reply, "2. second") {
		t.Fatalf("reply = %q", reply)
	}
}

func TestQueueReportsEmptyState(t *testing.T) {
	reply := runCommand(NewQueueCommand(newGuildQueues()), "g1")
	if !strings.Contains(reply, "empty") {
		t.Fatalf("reply = %q", reply)
	}
}
