package app

import "testing"

func TestLoadConfigDoesNotReadAppIDFromEnvironment(t *testing.T) {
	t.Setenv("DISCORD_APP_ID", "app-1")
	t.Setenv("DISCORD_BOT_TOKEN", "token")

	config := LoadConfig()

	if config.AppID != "" {
		t.Fatalf("LoadConfig() AppID = %q, want empty because the bot ID is derived from DISCORD_BOT_TOKEN", config.AppID)
	}
	if config.BotToken != "token" {
		t.Fatalf("LoadConfig() BotToken = %q, want %q", config.BotToken, "token")
	}
}
