package music

import (
	"fmt"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type musicModule struct {
	config Config
	node   AudioNode
	queues *guildQueues
}

func (m *musicModule) Name() string { return "music" }

func (m *musicModule) Register(registration core.Registration) error {
	lookup := m.config.VoiceLookup
	if lookup == nil {
		lookup = registration.VoiceLookup
	}
	commands := []core.Command{
		NewPlayCommand(m.node, m.queues, lookup),
		NewSkipCommand(m.node, m.queues, lookup),
		NewQueueCommand(m.queues),
		NewStopCommand(m.node, m.queues),
		NewNowPlayingCommand(m.node, m.queues),
	}
	for _, command := range commands {
		if err := registration.Commands.Add(command); err != nil {
			return fmt.Errorf("register music command %s: %w", command.Name, err)
		}
	}
	return nil
}

func (m *musicModule) Close() error {
	return m.node.Close()
}
