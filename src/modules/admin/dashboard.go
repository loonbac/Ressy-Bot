package admin

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

var dashboardTemplate = template.Must(template.New("dashboard").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>Ressy control panel</title></head>
<body>
<h1>Ressy control panel</h1>
<h2>Installed commands</h2>
<ul>{{range .Commands}}<li>{{.}}</li>{{else}}<li>None</li>{{end}}</ul>
<h2>Compiled plugins</h2>
<ul>{{range .Plugins}}<li>{{.}}</li>{{else}}<li>None</li>{{end}}</ul>
<form method="POST" action="/modules/load">
<label>Plugin file <input name="name" required></label><button type="submit">Load</button>
</form>
<form method="POST" action="/modules/toggle">
<label>Guild <input name="guild" required></label>
<label>Module <input name="module" required></label>
<label>Enabled <input type="checkbox" name="enabled"></label>
<button type="submit">Save</button>
</form>
</body></html>`))

func (p *panel) serveDashboard(w http.ResponseWriter, r *http.Request) {
	if !p.hasSession(r) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	view := struct {
		Commands []string
		Plugins  []string
	}{}
	if p.catalog != nil {
		for _, command := range p.catalog.List() {
			view.Commands = append(view.Commands, command.Name)
		}
	}
	entries, err := os.ReadDir(p.config.PluginDirectory)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".so") {
				view.Plugins = append(view.Plugins, entry.Name())
			}
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := dashboardTemplate.Execute(w, view); err != nil {
		http.Error(w, "dashboard template failed", http.StatusInternalServerError)
	}
}

func (p *panel) serveToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !p.hasSession(r) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	if p.settings == nil {
		http.Error(w, "module settings unavailable", http.StatusServiceUnavailable)
		return
	}
	guild := r.FormValue("guild")
	module := r.FormValue("module")
	if guild == "" || module == "" {
		http.Error(w, "guild and module are required", http.StatusBadRequest)
		return
	}
	enabled := r.FormValue("enabled") == "on"
	if err := p.settings.SetEnabled(core.GuildID(guild), core.ModuleName(module), enabled); err != nil {
		http.Error(w, "cannot save module setting", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (p *panel) serveLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !p.hasSession(r) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	name := r.FormValue("name")
	if name == "" || filepath.Base(name) != name || !strings.HasSuffix(name, ".so") || strings.Contains(name, string(filepath.Separator)) {
		http.Error(w, "plugin name must be a literal .so basename", http.StatusBadRequest)
		return
	}
	path := filepath.Join(p.config.PluginDirectory, name)
	if _, err := os.Stat(path); err != nil {
		http.Error(w, "plugin file not found", http.StatusNotFound)
		return
	}
	module, err := p.open(path)
	if err != nil {
		http.Error(w, "plugin could not be loaded", http.StatusBadGateway)
		return
	}
	if err := module.Register(core.Registration{Commands: p.catalog, Settings: p.settings}); err != nil {
		http.Error(w, "plugin could not be registered", http.StatusConflict)
		return
	}
	if p.publish != nil {
		if err := p.publish(p.catalog); err != nil {
			http.Error(w, "commands could not be published to Discord", http.StatusBadGateway)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
