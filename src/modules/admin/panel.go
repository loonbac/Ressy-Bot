package admin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type panel struct {
	config   Config
	catalog  core.CommandCatalog
	settings core.ModuleSettings
	server   *http.Server
}

func (p *panel) Name() string { return "admin" }

func (p *panel) Register(registration core.Registration) error {
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
