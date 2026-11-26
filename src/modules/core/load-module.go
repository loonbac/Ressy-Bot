package core

import "fmt"

// LoadModule reads the exported Module variable from a plugin.
func LoadModule(open Opener, path string) (Module, error) {
	if path == "" {
		return nil, ErrPluginPathEmpty
	}
	symbols, err := open(path)
	if err != nil {
		return nil, fmt.Errorf("open plugin %q: %w", path, err)
	}
	symbol, err := symbols.Lookup(ModuleSymbol)
	if err != nil {
		return nil, fmt.Errorf("look up %s in %q: %w", ModuleSymbol, path, err)
	}
	switch value := symbol.(type) {
	case *Module:
		if value == nil || *value == nil {
			return nil, fmt.Errorf("%w: %q exports an empty Module variable", ErrPluginContractInvalid, path)
		}
		return *value, nil
	case Module:
		if value == nil {
			return nil, fmt.Errorf("%w: %q exports an empty Module", ErrPluginContractInvalid, path)
		}
		return value, nil
	default:
		return nil, fmt.Errorf("%w: %q exports %T", ErrPluginContractInvalid, path, symbol)
	}
}
