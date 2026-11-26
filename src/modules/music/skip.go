package music

import "github.com/loonbac/Ressy-Bot/src/modules/core"

// NewSkipCommand advances the guild queue and restarts playback on the next track.
func NewSkipCommand(node AudioNode, queues *guildQueues, lookup VoiceLookup) core.Command {
	return core.Command{
		Name:        "skip",
		Description: "Skip the current track",
		Module:      "music",
		Execute: func(context core.CommandContext) error {
			track, found := queues.Next(context.GuildID)
			voice, ok := lookup(context.GuildID)
			if !found || !ok {
				return context.Reply("The queue is empty or you are not in a voice channel.")
			}
			if err := node.Stop(context.GuildID); err != nil {
				return err
			}
			if err := node.Play(context.GuildID, track, voice); err != nil {
				return err
			}
			return context.Reply("Now playing: " + track.Title)
		},
	}
}
