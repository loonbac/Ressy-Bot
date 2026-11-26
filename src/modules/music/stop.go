package music

import "github.com/loonbac/Ressy-Bot/src/modules/core"

// NewStopCommand clears a guild queue and stops its player.
func NewStopCommand(node AudioNode, queues *guildQueues) core.Command {
	return core.Command{
		Name:        "stop",
		Description: "Stop playback and clear the queue",
		Module:      "music",
		Execute: func(context core.CommandContext) error {
			queues.Clear(context.GuildID)
			if err := node.Stop(context.GuildID); err != nil {
				return err
			}
			return context.Reply("Playback stopped.")
		},
	}
}
