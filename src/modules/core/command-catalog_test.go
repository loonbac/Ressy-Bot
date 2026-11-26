package core

import (
	"errors"
	"testing"
)

func TestCommandCatalogListsCommandsInNameOrder(t *testing.T) {
	catalog := NewCommandCatalog()
	for _, name := range []string{"skip", "play", "queue"} {
		if err := catalog.Add(Command{Name: name}); err != nil {
			t.Fatalf("Add(%q) returned error: %v", name, err)
		}
	}

	got := catalog.List()
	want := []string{"play", "queue", "skip"}
	if len(got) != len(want) {
		t.Fatalf("got %d commands, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Fatalf("command %d = %q, want %q", i, got[i].Name, name)
		}
	}
}

func TestCommandCatalogRejectsDuplicateName(t *testing.T) {
	catalog := NewCommandCatalog()
	if err := catalog.Add(Command{Name: "play"}); err != nil {
		t.Fatalf("first Add returned error: %v", err)
	}

	err := catalog.Add(Command{Name: "play"})
	if !errors.Is(err, ErrCommandAlreadyRegistered) {
		t.Fatalf("error = %v, want ErrCommandAlreadyRegistered", err)
	}
}

func TestCommandCatalogIsSafeUnderConcurrentAddAndList(t *testing.T) {
	catalog := NewCommandCatalog()
	added := make(chan error, 1)
	go func() {
		added <- catalog.Add(Command{Name: "play"})
	}()
	catalog.List()
	<-added
}

func TestCommandCatalogRejectsEmptyName(t *testing.T) {
	err := NewCommandCatalog().Add(Command{})
	if !errors.Is(err, ErrCommandNameEmpty) {
		t.Fatalf("error = %v, want ErrCommandNameEmpty", err)
	}
}
