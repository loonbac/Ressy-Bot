package music

import (
	"strings"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// NewQueueCommand lists the tracks waiting in a guild.
func NewQueueCommand(queues *guildQueues) core.Command {
	return core.Command{
		Name:        "queue",
		Description: "Show the guild queue",
		Module:      "music",
		Execute: func(context core.CommandContext) error {
			tracks := queues.List(context.GuildID)
			if len(tracks) == 0 {
				return context.Reply("The queue is empty.")
			}
			titles := make([]string, 0, len(tracks))
			for index, track := range tracks {
				titles = append(titles, itoa(index+1)+". "+track.Title)
			}
			return context.Reply(strings.Join(titles, "\n"))
		},
	}
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 20)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
