// Package core defines the contract shared by the bot host and every plugin.
package core

import "errors"

// ModuleSymbol is the literal exported variable every plugin must declare.
const ModuleSymbol = "Module"

// GuildID identifies a Discord guild.
type GuildID string

// ModuleName identifies an installed module.
type ModuleName string

// VoiceState is the Discord data an audio node needs to join a channel.
type VoiceState struct {
	Token     string
	Endpoint  string
	SessionID string
	ChannelID string
}

// VoiceLookup finds the bot voice state for a guild.
type VoiceLookup func(guildID string) (VoiceState, bool)

var (
	// ErrPluginPathEmpty means no plugin path was supplied.
	ErrPluginPathEmpty = errors.New("plugin path is empty")
	// ErrPluginContractInvalid means a plugin did not export the expected module variable.
	ErrPluginContractInvalid = errors.New("plugin does not export core.Module")
	// ErrCommandNameEmpty means a command has no literal name.
	ErrCommandNameEmpty = errors.New("command name is empty")
	// ErrCommandAlreadyRegistered means two modules claimed the same command name.
	ErrCommandAlreadyRegistered = errors.New("command already registered")
)

// Module is the single capability a plugin grants to the bot.
type Module interface {
	Name() string
	Register(registration Registration) error
	Close() error
}

// Registration carries the host services a module may use.
type Registration struct {
	Commands    CommandCatalog
	Settings    ModuleSettings
	VoiceLookup VoiceLookup
}

// Command describes one slash command contributed by a module.
type Command struct {
	Name        string
	Description string
	Module      string
	Execute     func(context CommandContext) error
}

// CommandContext carries the Discord facts a command needs to reply.
type CommandContext struct {
	GuildID   string
	ChannelID string
	UserID    string
	Args      []string
	Reply     func(string) error
}

// CommandCatalog holds commands indexed by their literal name.
type CommandCatalog interface {
	Add(command Command) error
	List() []Command
}

// PluginSymbols is the lookup surface returned by a plugin opener.
type PluginSymbols interface {
	Lookup(name string) (any, error)
}

// Opener opens plugin symbols at a path.
type Opener func(path string) (PluginSymbols, error)
