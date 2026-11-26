package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/loonbac/Ressy-Bot/src/app"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

func main() {
	config := app.LoadConfig()
	if err := os.MkdirAll(filepath.Dir(config.SettingsPath), 0o755); err != nil {
		log.Fatal(err)
	}
	settings, err := core.NewFileModuleSettings(config.SettingsPath)
	if err != nil {
		log.Fatal(err)
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
	session, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates
	session.AddHandler(app.ServeInteractions(catalog, settings))
	voices.Attach(session)
	if err := session.Open(); err != nil {
		log.Fatal(err)
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
