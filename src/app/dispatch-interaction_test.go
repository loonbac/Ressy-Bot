package app

import (
	"errors"
	"testing"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type stubModuleSettings struct {
	enabled bool
}

func (s stubModuleSettings) IsEnabled(core.GuildID, core.ModuleName) bool { return s.enabled }
func (s stubModuleSettings) SetEnabled(core.GuildID, core.ModuleName, bool) error {
	return nil
}
func (s stubModuleSettings) List(core.GuildID) []core.ModuleName { return nil }

func TestDispatchInteractionExecutesMatchingCommand(t *testing.T) {
	catalog := core.NewCommandCatalog()
	called := false
	if err := catalog.Add(core.Command{
		Name:   "play",
		Module: "music",
		Execute: func(context core.CommandContext) error {
			called = true
			return context.Reply("queued")
		},
	}); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	var reply string
	err := DispatchInteraction(catalog, stubModuleSettings{enabled: true}, InteractionRequest{
		Name: "play", GuildID: "g1", UserID: "u1",
	}, func(message string) error {
		reply = message
		return nil
	})
	if err != nil {
		t.Fatalf("DispatchInteraction returned error: %v", err)
	}
	if !called || reply != "queued" {
		t.Fatalf("called = %v, reply = %q", called, reply)
	}
}

func TestDispatchInteractionRejectsDisabledModule(t *testing.T) {
	catalog := core.NewCommandCatalog()
	called := false
	if err := catalog.Add(core.Command{
		Name: "play", Module: "music", Execute: func(core.CommandContext) error {
			called = true
			return nil
		},
	}); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	err := DispatchInteraction(catalog, stubModuleSettings{enabled: false}, InteractionRequest{
		Name: "play", GuildID: "g1",
	}, func(string) error { return nil })
	if !errors.Is(err, ErrModuleDisabled) {
		t.Fatalf("error = %v, want ErrModuleDisabled", err)
	}
	if called {
		t.Fatal("command ran while its module was disabled")
	}
}

func TestDispatchInteractionRejectsUnknownCommand(t *testing.T) {
	err := DispatchInteraction(core.NewCommandCatalog(), stubModuleSettings{enabled: true}, InteractionRequest{
		Name: "missing",
	}, func(string) error { return nil })
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("error = %v, want ErrCommandNotFound", err)
	}
}
