package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/loonbac/Ressy-Bot/src/app"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// applyBotIdentity stores the bot user ID derived from the token in config.AppID,
// which slash-command registration needs, and exports it as RESSY_BOT_USER_ID,
// which the music plugin reads when the host loads it.
func applyBotIdentity(config *app.Config, fetch func() (string, error)) error {
	id, err := fetch()
	if err != nil {
		return fmt.Errorf("resolve bot user ID from DISCORD_BOT_TOKEN: %w", err)
	}
	if id == "" {
		return fmt.Errorf("resolve bot user ID from DISCORD_BOT_TOKEN: Discord returned an empty user ID")
	}
	if err := os.Setenv("RESSY_BOT_USER_ID", id); err != nil {
		return fmt.Errorf("export RESSY_BOT_USER_ID: %w", err)
	}
	config.AppID = id
	return nil
}

func main() {
	config := app.LoadConfig()
	if err := os.MkdirAll(filepath.Dir(config.SettingsPath), 0o755); err != nil {
		log.Fatal(err)
	}
	settings, err := core.NewFileModuleSettings(config.SettingsPath)
	if err != nil {
		log.Fatal(err)
	}

	// Resolve the bot identity before plugins load: the music plugin reads
	// RESSY_BOT_USER_ID at load time, and slash-command registration uses the
	// same ID. DISCORD_APP_ID is no longer needed.
	var session *discordgo.Session
	if config.BotToken != "" {
		session, err = discordgo.New("Bot " + config.BotToken)
		if err != nil {
			log.Fatalf("create Discord session from DISCORD_BOT_TOKEN: %v", err)
		}
		fetchIdentity := func() (string, error) {
			user, err := session.User("@me")
			if err != nil {
				return "", err
			}
			return user.ID, nil
		}
		if err := applyBotIdentity(&config, fetchIdentity); err != nil {
			log.Fatal(err)
		}
		log.Printf("resolved bot user ID %s from DISCORD_BOT_TOKEN", config.AppID)
	}

	voices := app.NewVoiceStateStore()
	config.VoiceLookup = voices.Lookup
	catalog, modules, err := app.Start(config, settings)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		for _, module := range modules {
			if err := module.Close(); err != nil {
				log.Printf("close module %s: %v", module.Name(), err)
			}
		}
	}()

	log.Printf("loaded %d modules and %d commands", len(modules), len(catalog.List()))

	if config.BotToken == "" {
		log.Print("DISCORD_BOT_TOKEN is not set; Discord stays offline, panel only")
		<-make(chan struct{})
	}
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates
	session.AddHandler(app.ServeInteractions(catalog, settings))
	voices.Attach(session)
	if err := session.Open(); err != nil {
		log.Fatalf("open Discord session: %v", err)
	}
	defer session.Close()

	if config.AppID != "" {
		publisher := app.NewDiscordCommandPublisher(session)
		publish := func(loaded core.CommandCatalog) error { return app.RegisterCommands(publisher, config.AppID, loaded) }
		if err := publish(catalog); err != nil {
			log.Printf("register commands: %v", err)
		}
		for _, module := range modules {
			if setter, ok := module.(interface {
				SetCommandCatalogPublisher(func(core.CommandCatalog) error)
			}); ok {
				setter.SetCommandCatalogPublisher(publish)
			}
		}
	}
	log.Print("bot connected; press CTRL-C to stop")
	<-make(chan struct{})
}
