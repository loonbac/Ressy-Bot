package admin

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	oauthStateCookie = "ressy_oauth_state"
	sessionCookieName = "ressy_session"
)

func signValue(secret, value string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(value))
	return value + "." + hex.EncodeToString(mac.Sum(nil))
}

func verifyValue(secret, signed string) (string, bool) {
	value, signature, found := strings.Cut(signed, ".")
	if !found {
		return "", false
	}
	return value, hmac.Equal([]byte(signature), []byte(strings.TrimPrefix(signValue(secret, value), value+".")))
}

func randomToken() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return ""
	}
	return hex.EncodeToString(buffer)
}

func (p *panel) serveLogin(w http.ResponseWriter, r *http.Request) {
	state := randomToken()
	if state == "" {
		http.Error(w, "cannot create OAuth state", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: oauthStateCookie, Value: signValue(p.config.SessionSecret, state),
		Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	target := "https://discord.com/oauth2/authorize?client_id=" + url.QueryEscape(p.config.ClientID) +
		"&response_type=code&scope=identify&state=" + url.QueryEscape(state) +
		"&redirect_uri=" + url.QueryEscape(p.config.RedirectURL)
	http.Redirect(w, r, target, http.StatusFound)
}

func (p *panel) serveCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie(oauthStateCookie)
	if err != nil {
		http.Error(w, "missing OAuth state", http.StatusUnauthorized)
		return
	}
	storedState, valid := verifyValue(p.config.SessionSecret, cookie.Value)
	if !valid || storedState != state || state == "" {
		http.Error(w, "invalid OAuth state", http.StatusUnauthorized)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	exchange := p.config.ExchangeCode
	if exchange == nil {
		exchange = exchangeDiscordCode
	}
	token, err := exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusBadGateway)
		return
	}
	identify := p.config.Identify
	if identify == nil {
		identify = identifyDiscordUser
	}
	userID, err := identify(r.Context(), token)
	if err != nil {
		http.Error(w, "identify failed", http.StatusBadGateway)
		return
	}
	if !p.isAdmin(userID) {
		http.Error(w, "user is not an administrator", http.StatusForbidden)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookieName, Value: signValue(p.config.SessionSecret, userID),
		Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (p *panel) isAdmin(userID string) bool {
	for _, allowed := range p.config.AdminUserIDs {
		if allowed == userID {
			return true
		}
	}
	return false
}

func (p *panel) hasSession(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return false
	}
	userID, valid := verifyValue(p.config.SessionSecret, cookie.Value)
	return valid && p.isAdmin(userID)
}

func exchangeDiscordCode(ctx context.Context, code string) (string, error) {
	values := url.Values{
		"client_id": {os.Getenv("DISCORD_CLIENT_ID")}, "client_secret": {os.Getenv("DISCORD_CLIENT_SECRET")},
		"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {os.Getenv("RESSY_PANEL_REDIRECT")},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://discord.com/api/oauth2/token", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange returned %s", response.Status)
	}
	var payload struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	return payload.AccessToken, nil
}

func identifyDiscordUser(ctx context.Context, token string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://discord.com/api/users/@me", nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("identify returned %s", response.Status)
	}
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	return payload.ID, nil
}
