package core

import "plugin"

type nativeOpener struct{}

type nativeSymbols struct {
	loaded *plugin.Plugin
}

func (s nativeSymbols) Lookup(name string) (any, error) {
	return s.loaded.Lookup(name)
}

func (nativeOpener) Open(path string) (PluginSymbols, error) {
	loaded, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	return nativeSymbols{loaded: loaded}, nil
}

// OpenPlugin loads a compiled .so module from the local filesystem.
func OpenPlugin(path string) (Module, error) {
	return LoadModule(nativeOpener{}.Open, path)
}
