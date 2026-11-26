package main

import (
	"log"

	"github.com/loonbac/Ressy-Bot/src/app"
)

func main() {
	config := app.LoadConfig()
	catalog, modules, err := app.Start(config)
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
}
