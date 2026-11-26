package admin

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

type stubLoadedPlugin struct{ core.Module }

func (stubLoadedPlugin) Register(r core.Registration) error {
	return r.Commands.Add(core.Command{Name: "loaded"})
}

type stubAdminSettings struct {
	guild   core.GuildID
	module  core.ModuleName
	enabled bool
}

func (s *stubAdminSettings) IsEnabled(core.GuildID, core.ModuleName) bool { return true }
func (s *stubAdminSettings) SetEnabled(guild core.GuildID, module core.ModuleName, enabled bool) error {
	s.guild, s.module, s.enabled = guild, module, enabled
	return nil
}
func (s *stubAdminSettings) List(core.GuildID) []core.ModuleName { return nil }

func newTestPanel() (*panel, http.Handler, *stubAdminSettings) {
	settings := &stubAdminSettings{}
	p := &panel{config: Config{
		ClientID: "client", ClientSecret: "secret", SessionSecret: "session-secret",
		RedirectURL: "https://example.test/callback", PanelAddress: "127.0.0.1:0",
		PluginDirectory: "plugins", AdminUserIDs: []string{"admin-1"},
	}}
	return p, p.handler(), settings
}

func sessionRequest(method, target string) *http.Request {
	request := httptest.NewRequest(method, target, nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: signValue("session-secret", "admin-1")})
	return request
}

func TestServeLoginSetsSignedStateCookie(t *testing.T) {
	_, handler, _ := newTestPanel()
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/login", nil))

	if recorder.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", recorder.Code)
	}
	if !strings.Contains(recorder.Header().Get("Location"), "discord.com/oauth2/authorize") {
		t.Fatalf("location = %q", recorder.Header().Get("Location"))
	}
	var stateCookie *http.Cookie
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == oauthStateCookie {
			stateCookie = cookie
		}
	}
	if stateCookie == nil {
		t.Fatal("no OAuth state cookie was set")
	}
	if _, valid := verifyValue("session-secret", stateCookie.Value); !valid {
		t.Fatal("state cookie is not signed")
	}
}

func TestCallbackRejectsMismatchedState(t *testing.T) {
	_, handler, _ := newTestPanel()
	request := httptest.NewRequest(http.MethodGet, "/callback?code=c&state=other", nil)
	request.AddCookie(&http.Cookie{Name: oauthStateCookie, Value: signValue("session-secret", "expected")})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
}

func TestCallbackRejectsUserOutsideAllowlist(t *testing.T) {
	p, handler, _ := newTestPanel()
	p.config.ExchangeCode = func(context.Context, string) (string, error) { return "token", nil }
	p.config.Identify = func(context.Context, string) (string, error) { return "intruder", nil }
	request := httptest.NewRequest(http.MethodGet, "/callback?code=c&state=s", nil)
	request.AddCookie(&http.Cookie{Name: oauthStateCookie, Value: signValue("session-secret", "s")})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", recorder.Code)
	}
}

func TestDashboardRequiresSession(t *testing.T) {
	_, handler, _ := newTestPanel()
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
}

func TestToggleUpdatesModuleSettings(t *testing.T) {
	p, handler, settings := newTestPanel()
	p.settings = settings
	form := url.Values{"guild": {"g1"}, "module": {"music"}, "enabled": {"on"}}
	request := sessionRequest(http.MethodPost, "/modules/toggle")
	request.Body = io.NopCloser(strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", recorder.Code)
	}
	if settings.guild != "g1" || settings.module != "music" || !settings.enabled {
		t.Fatalf("settings = %+v", settings)
	}
}

func TestRegisterRejectsUnsafeConfig(t *testing.T) {
	p := &panel{config: Config{AdminUserIDs: []string{"admin-1"}}}
	if err := p.Register(core.Registration{}); err == nil || len(p.config.AdminUserIDs) != 0 {
		t.Fatalf("insecure config accepted: %v", err)
	}
}

func TestMutatingHandlersRejectGET(t *testing.T) {
	p, handler, settings := newTestPanel()
	p.settings = settings
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, sessionRequest(http.MethodGet, "/modules/toggle?guild=g1&module=music&enabled=on"))
	if recorder.Code != http.StatusMethodNotAllowed || settings.guild != "" {
		t.Fatalf("status = %d, settings = %+v", recorder.Code, settings)
	}
}

func TestLoadPublishesLoadedCommands(t *testing.T) {
	p, handler, _ := newTestPanel()
	p.catalog, p.config.PluginDirectory = core.NewCommandCatalog(), t.TempDir()
	p.open = func(string) (core.Module, error) { return stubLoadedPlugin{}, nil }
	_ = os.WriteFile(p.config.PluginDirectory+"/x.so", nil, 0o644)
	var published string
	p.publish = func(c core.CommandCatalog) error { published = c.List()[0].Name; return nil }
	handler.ServeHTTP(httptest.NewRecorder(), sessionRequest(http.MethodPost, "/modules/load?name=x.so"))
	if published != "loaded" {
		t.Fatalf("published = %q, want loaded", published)
	}
}

func TestLoadRejectsPathTraversal(t *testing.T) {
	_, handler, _ := newTestPanel()
	form := url.Values{"name": {"../evil.so"}}
	request := sessionRequest(http.MethodPost, "/modules/load")
	request.Body = io.NopCloser(strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}
