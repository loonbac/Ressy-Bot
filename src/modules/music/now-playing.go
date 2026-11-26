package music

import "github.com/loonbac/Ressy-Bot/src/modules/core"

// NewNowPlayingCommand reports the current track of a guild.
func NewNowPlayingCommand(queues *guildQueues) core.Command {
	return core.Command{
		Name:        "now-playing",
		Description: "Show the current track",
		Module:      "music",
		Execute: func(context core.CommandContext) error {
			track, found := queues.Current(context.GuildID)
			if !found {
				return context.Reply("Nothing is playing here.")
			}
			return context.Reply("Now playing: " + track.Title)
		},
	}
}
