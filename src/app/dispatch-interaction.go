package app

import (
	"errors"
	"fmt"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

var (
	// ErrCommandNotFound means no installed module owns an interaction name.
	ErrCommandNotFound = errors.New("command not found")
	// ErrModuleDisabled means the guild turned the owning module off.
	ErrModuleDisabled = errors.New("module disabled")
)

// InteractionRequest is the Discord interaction reduced to its command facts.
type InteractionRequest struct {
	Name      string
	GuildID   string
	ChannelID string
	UserID    string
	Options   []string
}

// DispatchInteraction runs one command when its module is enabled in the guild.
func DispatchInteraction(catalog core.CommandCatalog, settings core.ModuleSettings, request InteractionRequest, respond func(string) error) error {
	for _, command := range catalog.List() {
		if command.Name != request.Name {
			continue
		}
		if !settings.IsEnabled(core.GuildID(request.GuildID), core.ModuleName(command.Module)) {
			return fmt.Errorf("%w: %s", ErrModuleDisabled, command.Module)
		}
		return command.Execute(core.CommandContext{
			GuildID:   request.GuildID,
			ChannelID: request.ChannelID,
			UserID:    request.UserID,
			Args:      request.Options,
			Reply:     respond,
		})
	}
	return fmt.Errorf("%w: %s", ErrCommandNotFound, request.Name)
}
