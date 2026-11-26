package admin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type panel struct {
	config   Config
	catalog  core.CommandCatalog
	settings core.ModuleSettings
	server   *http.Server
	publish  func(core.CommandCatalog) error
	open     func(string) (core.Module, error)
}

func (p *panel) Name() string { return "admin" }

// SetCommandCatalogPublisher installs the callback that publishes catalog commands to Discord.
func (p *panel) SetCommandCatalogPublisher(publish func(core.CommandCatalog) error) {
	p.publish = publish
}

func (p *panel) Register(registration core.Registration) error {
	admins := make([]string, 0, len(p.config.AdminUserIDs))
	for _, id := range p.config.AdminUserIDs {
		if id = strings.TrimSpace(id); id != "" {
			admins = append(admins, id)
		}
	}
	if strings.TrimSpace(p.config.SessionSecret) == "" || len(admins) == 0 {
		p.config.AdminUserIDs = nil
		return errors.New("admin panel: a session secret and at least one admin user id are required")
	}
	p.config.AdminUserIDs, p.open = admins, core.OpenPlugin
	p.catalog = registration.Commands
	p.settings = registration.Settings
	p.server = &http.Server{Addr: p.config.PanelAddress, Handler: p.handler()}
	go func() {
		if err := p.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("admin panel: %v", err)
		}
	}()
	return nil
}

func (p *panel) Close() error {
	if p.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.server.Shutdown(ctx)
}

func (p *panel) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", p.serveLogin)
	mux.HandleFunc("/callback", p.serveCallback)
	mux.HandleFunc("/modules/toggle", p.serveToggle)
	mux.HandleFunc("/modules/load", p.serveLoad)
	mux.HandleFunc("/", p.serveDashboard)
	return mux
}
