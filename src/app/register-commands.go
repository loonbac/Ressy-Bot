package app

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// CommandPublisher registers application commands with Discord.
type CommandPublisher interface {
	Create(appID, guildID string, command *discordgo.ApplicationCommand) error
}

type discordCommandPublisher struct {
	session *discordgo.Session
}

// NewDiscordCommandPublisher adapts a Discord session to CommandPublisher.
func NewDiscordCommandPublisher(session *discordgo.Session) CommandPublisher {
	return discordCommandPublisher{session: session}
}

func (p discordCommandPublisher) Create(appID, guildID string, command *discordgo.ApplicationCommand) error {
	_, err := p.session.ApplicationCommandCreate(appID, guildID, command)
	return err
}

// RegisterCommands publishes every installed command as a global slash command.
func RegisterCommands(publisher CommandPublisher, appID string, catalog core.CommandCatalog) error {
	for _, command := range catalog.List() {
		err := publisher.Create(appID, "", &discordgo.ApplicationCommand{
			Name:        command.Name,
			Description: command.Description,
		})
		if err != nil {
			return fmt.Errorf("publish command %s: %w", command.Name, err)
		}
	}
	return nil
}
