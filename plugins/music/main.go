package main

import (
	"os"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
	"github.com/loonbac/Ressy-Bot/src/modules/music"
)

func pluginEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

// Module is the literal symbol the host loads from this plugin.
var Module core.Module = music.New(music.Config{
	Lavalink: music.LavalinkConfig{
		BaseURL:  pluginEnv("RESSY_LAVALINK_URL", "http://127.0.0.1:2333"),
		Password: os.Getenv("RESSY_LAVALINK_PASSWORD"),
		UserID:   os.Getenv("RESSY_BOT_USER_ID"),
	},
})
