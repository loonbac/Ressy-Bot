package music

import (
	"strings"
	"testing"
)

func TestStopClearsQueueAndStopsNode(t *testing.T) {
	node := &fakeAudioNode{}
	queues := newGuildQueues()
	queues.Enqueue("g1", Track{Title: "first"})

	reply := runCommand(NewStopCommand(node, queues), "g1")

	if !strings.Contains(reply, "stopped") {
		t.Fatalf("reply = %q", reply)
	}
	if len(queues.List("g1")) != 0 {
		t.Fatalf("queue was not cleared: %+v", queues.List("g1"))
	}
	if len(node.stopped) != 1 || node.stopped[0] != "g1" {
		t.Fatalf("stopped = %v", node.stopped)
	}
}
