package music

import "github.com/loonbac/Ressy-Bot/src/modules/core"

// NewSkipCommand advances the guild queue.
func NewSkipCommand(node AudioNode, queues *guildQueues) core.Command {
	return core.Command{
		Name:        "skip",
		Description: "Skip the current track",
		Module:      "music",
		Execute: func(context core.CommandContext) error {
			track, found := queues.Next(context.GuildID)
			if !found {
				return context.Reply(ErrQueueEmpty.Error())
			}
			return context.Reply("Now playing: " + track.Title)
		},
	}
}
