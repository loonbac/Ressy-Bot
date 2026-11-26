package main

import (
	"os"
	"strings"

	"github.com/loonbac/Ressy-Bot/src/modules/admin"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

func pluginEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

// Module is the literal symbol the host loads from this plugin.
var Module core.Module = admin.New(admin.Config{
	ClientID:        os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret:    os.Getenv("DISCORD_CLIENT_SECRET"),
	SessionSecret:   os.Getenv("RESSY_SESSION_SECRET"),
	RedirectURL:     os.Getenv("RESSY_PANEL_REDIRECT"),
	PanelAddress:    pluginEnv("RESSY_PANEL_ADDR", ":8080"),
	PluginDirectory: pluginEnv("RESSY_PLUGIN_DIR", "./plugins"),
	AdminUserIDs:    strings.Split(os.Getenv("RESSY_ADMIN_IDS"), ","),
})
