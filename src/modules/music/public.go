// Package music provides YouTube playback as an installable bot module.
package music

import "github.com/loonbac/Ressy-Bot/src/modules/core"

// LavalinkConfig locates the audio node embedded in the bot container.
type LavalinkConfig struct {
	BaseURL  string
	Password string
	UserID   string
}

// Config carries the music module dependencies.
type Config struct {
	Node        AudioNode
	VoiceLookup VoiceLookup
	Lavalink    LavalinkConfig
}

// New creates a music module. It uses Config.Node when supplied, otherwise it creates a Lavalink node.
func New(config Config) core.Module {
	node := config.Node
	if node == nil {
		node = NewLavalinkNode(config.Lavalink)
	}
	return &musicModule{config: config, node: node, queues: newGuildQueues()}
}
