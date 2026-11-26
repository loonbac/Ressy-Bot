package app

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type stubCommandPublisher struct {
	names []string
	err   error
}

func (p *stubCommandPublisher) Create(_, _ string, command *discordgo.ApplicationCommand) error {
	p.names = append(p.names, command.Name)
	return p.err
}

func TestRegisterCommandsPublishesCatalog(t *testing.T) {
	catalog := core.NewCommandCatalog()
	for _, name := range []string{"play", "skip"} {
		if err := catalog.Add(core.Command{Name: name, Description: name + " audio"}); err != nil {
			t.Fatalf("Add(%q) returned error: %v", name, err)
		}
	}
	publisher := &stubCommandPublisher{}

	if err := RegisterCommands(publisher, "app-1", catalog); err != nil {
		t.Fatalf("RegisterCommands returned error: %v", err)
	}
	if len(publisher.names) != 2 || publisher.names[0] != "play" || publisher.names[1] != "skip" {
		t.Fatalf("published %v, want [play skip]", publisher.names)
	}
}

func TestRegisterCommandsReportsPublishFailure(t *testing.T) {
	catalog := core.NewCommandCatalog()
	if err := catalog.Add(core.Command{Name: "play"}); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	cause := errors.New("discord unavailable")

	err := RegisterCommands(&stubCommandPublisher{err: cause}, "app-1", catalog)
	if !errors.Is(err, cause) {
		t.Fatalf("error = %v, want wrapped %v", err, cause)
	}
}
