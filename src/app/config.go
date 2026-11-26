// Package app wires the bot process from its configured modules.
package app

import (
	"os"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// Config holds the process settings read from the environment.
type Config struct {
	BotToken        string
	AppID           string // bot user ID, set at startup from DISCORD_BOT_TOKEN
	PluginDirectory string
	SettingsPath    string
	PanelAddress    string
	VoiceLookup     core.VoiceLookup
}

// LoadConfig reads settings from environment variables. The bot user ID is not
// read here: the host derives it from DISCORD_BOT_TOKEN via GET /users/@me.
func LoadConfig() Config {
	return Config{
		BotToken:        os.Getenv("DISCORD_BOT_TOKEN"),
		PluginDirectory: envOrDefault("RESSY_PLUGIN_DIR", "./plugins"),
		SettingsPath:    envOrDefault("RESSY_SETTINGS_PATH", "data/modules.json"),
		PanelAddress:    envOrDefault("RESSY_PANEL_ADDR", ":8080"),
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
