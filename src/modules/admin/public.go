// Package admin serves the web control panel as a bot module.
package admin

import (
	"context"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// Config carries the panel dependencies supplied by the plugin.
type Config struct {
	ClientID       string
	ClientSecret   string
	SessionSecret  string
	RedirectURL    string
	PanelAddress   string
	PluginDirectory string
	AdminUserIDs   []string
	ExchangeCode   func(context.Context, string) (string, error)
	Identify       func(context.Context, string) (string, error)
}

// New creates an admin panel module.
func New(config Config) core.Module {
	return &panel{config: config}
}
