// Package app wires the bot process from its configured modules.
package app

import "os"

// Config holds the process settings read from the environment.
type Config struct {
	BotToken        string
	PluginDirectory string
	PanelAddress    string
}

// LoadConfig reads settings from environment variables.
func LoadConfig() Config {
	return Config{
		BotToken:        os.Getenv("DISCORD_BOT_TOKEN"),
		PluginDirectory: envOrDefault("RESSY_PLUGIN_DIR", "./plugins"),
		PanelAddress:    envOrDefault("RESSY_PANEL_ADDR", ":8080"),
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
