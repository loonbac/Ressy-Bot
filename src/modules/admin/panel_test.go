package admin

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

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
