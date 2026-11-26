package main

import (
	"errors"
	"os"
	"testing"

	"github.com/loonbac/Ressy-Bot/src/app"
)

func TestApplyBotIdentityUsesResolvedIDEverywhere(t *testing.T) {
	t.Setenv("RESSY_BOT_USER_ID", "")
	config := &app.Config{AppID: "configured-elsewhere"}
	fetch := func() (string, error) { return "bot-123", nil }

	if err := applyBotIdentity(config, fetch); err != nil {
		t.Fatalf("applyBotIdentity returned error: %v", err)
	}
	if config.AppID != "bot-123" {
		t.Fatalf("config.AppID = %q, want %q", config.AppID, "bot-123")
	}
	if got := os.Getenv("RESSY_BOT_USER_ID"); got != "bot-123" {
		t.Fatalf("RESSY_BOT_USER_ID = %q, want %q", got, "bot-123")
	}
}

func TestApplyBotIdentityReportsLookupFailure(t *testing.T) {
	t.Setenv("RESSY_BOT_USER_ID", "unchanged")
	config := &app.Config{AppID: "unchanged"}
	cause := errors.New("discord unavailable")
	fetch := func() (string, error) { return "", cause }

	err := applyBotIdentity(config, fetch)
	if !errors.Is(err, cause) {
		t.Fatalf("error = %v, want wrapped %v", err, cause)
	}
	if config.AppID != "unchanged" {
		t.Fatalf("config.AppID = %q, want %q after failure", config.AppID, "unchanged")
	}
	if got := os.Getenv("RESSY_BOT_USER_ID"); got != "unchanged" {
		t.Fatalf("RESSY_BOT_USER_ID = %q, want %q after failure", got, "unchanged")
	}
}

func TestApplyBotIdentityRejectsEmptyID(t *testing.T) {
	t.Setenv("RESSY_BOT_USER_ID", "unchanged")
	config := &app.Config{AppID: "unchanged"}
	fetch := func() (string, error) { return "", nil }

	if err := applyBotIdentity(config, fetch); err == nil {
		t.Fatal("applyBotIdentity returned nil error for an empty bot user ID")
	}
	if config.AppID != "unchanged" {
		t.Fatalf("config.AppID = %q, want %q after empty ID", config.AppID, "unchanged")
	}
	if got := os.Getenv("RESSY_BOT_USER_ID"); got != "unchanged" {
		t.Fatalf("RESSY_BOT_USER_ID = %q, want %q after empty ID", got, "unchanged")
	}
}
