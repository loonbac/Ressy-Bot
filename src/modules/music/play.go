package music

import (
	"strings"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// NewPlayCommand queues a YouTube track and starts playback.
func NewPlayCommand(node AudioNode, queues *guildQueues, lookup VoiceLookup) core.Command {
	return core.Command{
		Name:        "play",
		Description: "Queue a YouTube track",
		Module:      "music",
		Execute: func(context core.CommandContext) error {
			query := strings.Join(context.Args, " ")
			if query == "" {
				return context.Reply("Give me a YouTube URL or a search term.")
			}
			tracks, err := node.Resolve("ytsearch:" + query)
			if err != nil || len(tracks) == 0 {
				return context.Reply("No track found for that query.")
			}
			voice, found := lookup(context.GuildID)
			if !found {
				return context.Reply(ErrVoiceStateMissing.Error() + " Join a voice channel first.")
			}
			track := tracks[0]
			queues.Enqueue(context.GuildID, track)
			if err := node.Play(context.GuildID, track, voice); err != nil {
				return err
			}
			queues.SetCurrent(context.GuildID, track)
			return context.Reply("Queued: " + track.Title)
		},
	}
}
