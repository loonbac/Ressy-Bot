package music

import (
	"strings"
	"testing"
)

func TestSkipAdvancesQueue(t *testing.T) {
	queues := newGuildQueues()
	queues.Enqueue("g1", Track{Title: "first"})
	queues.Enqueue("g1", Track{Title: "second"})

	reply := runCommand(NewSkipCommand(&fakeAudioNode{}, queues), "g1")

	if !strings.Contains(reply, "first") {
		t.Fatalf("reply = %q", reply)
	}
	waiting := queues.List("g1")
	if len(waiting) != 1 || waiting[0].Title != "second" {
		t.Fatalf("queue = %+v", waiting)
	}
}

func TestSkipReportsEmptyQueue(t *testing.T) {
	reply := runCommand(NewSkipCommand(&fakeAudioNode{}, newGuildQueues()), "g1")
	if !strings.Contains(reply, "empty") {
		t.Fatalf("reply = %q", reply)
	}
}
