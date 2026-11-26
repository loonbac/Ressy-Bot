package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// Start loads installed modules and registers their commands.
func Start(config Config) (core.CommandCatalog, []core.Module, error) {
	catalog := core.NewCommandCatalog()
	modules := make([]core.Module, 0)

	entries, err := os.ReadDir(config.PluginDirectory)
	if err != nil {
		return nil, nil, fmt.Errorf("read plugin directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".so" {
			continue
		}
		path := filepath.Join(config.PluginDirectory, entry.Name())
		module, err := core.OpenPlugin(path)
		if err != nil {
			return nil, nil, fmt.Errorf("load %s: %w", entry.Name(), err)
		}
		if err := module.Register(core.Registration{Commands: catalog}); err != nil {
			return nil, nil, fmt.Errorf("register %s: %w", module.Name(), err)
		}
		modules = append(modules, module)
	}
	return catalog, modules, nil
}
