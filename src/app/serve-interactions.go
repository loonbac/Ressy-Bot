package app

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// ServeInteractions builds the Discord handler that dispatches slash commands.
func ServeInteractions(catalog core.CommandCatalog, settings core.ModuleSettings) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		data := interaction.ApplicationCommandData()
		request := InteractionRequest{
			Name:      data.Name,
			GuildID:   interaction.GuildID,
			ChannelID: interaction.ChannelID,
		}
		if interaction.Member != nil && interaction.Member.User != nil {
			request.UserID = interaction.Member.User.ID
		} else if interaction.User != nil {
			request.UserID = interaction.User.ID
		}
		for _, option := range data.Options {
			if value, ok := option.Value.(string); ok {
				request.Options = append(request.Options, value)
			}
		}
		respond := func(message string) error {
			return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: message},
			})
		}
		if err := DispatchInteraction(catalog, settings, request, respond); err != nil {
			log.Printf("dispatch %s: %v", request.Name, err)
			_ = respond("That command is unavailable.")
		}
	}
}
