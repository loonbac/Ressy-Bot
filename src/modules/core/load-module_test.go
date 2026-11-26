package core

import (
	"errors"
	"testing"
)

type stubModule struct {
	name string
}

func (m stubModule) Name() string                { return m.name }
func (m stubModule) Register(Registration) error { return nil }
func (m stubModule) Close() error                { return nil }

type stubSymbols struct {
	value any
	err   error
	asked string
}

func (s *stubSymbols) Lookup(name string) (any, error) {
	s.asked = name
	return s.value, s.err
}

func TestLoadModuleReturnsExportedModule(t *testing.T) {
	module := stubModule{name: "music"}
	symbols := &stubSymbols{value: &module}
	open := func(string) (PluginSymbols, error) { return symbols, nil }

	got, err := LoadModule(open, "/plugins/music.so")
	if err != nil {
		t.Fatalf("LoadModule returned error: %v", err)
	}
	if got.Name() != "music" {
		t.Fatalf("module name = %q, want music", got.Name())
	}
	if symbols.asked != ModuleSymbol {
		t.Fatalf("looked up %q, want %q", symbols.asked, ModuleSymbol)
	}
}

func TestLoadModuleAcceptsModuleValue(t *testing.T) {
	module := stubModule{name: "music"}
	open := func(string) (PluginSymbols, error) {
		return &stubSymbols{value: module}, nil
	}

	got, err := LoadModule(open, "/plugins/music.so")
	if err != nil {
		t.Fatalf("LoadModule returned error: %v", err)
	}
	if got.Name() != "music" {
		t.Fatalf("module name = %q, want music", got.Name())
	}
}

func TestLoadModuleReportsLookupFailure(t *testing.T) {
	cause := errors.New("symbol unavailable")
	open := func(string) (PluginSymbols, error) {
		return &stubSymbols{err: cause}, nil
	}

	_, err := LoadModule(open, "/plugins/music.so")
	if !errors.Is(err, cause) {
		t.Fatalf("error = %v, want wrapped %v", err, cause)
	}
}

func TestLoadModuleRejectsWrongSymbolType(t *testing.T) {
	open := func(string) (PluginSymbols, error) {
		return &stubSymbols{value: "not a module"}, nil
	}

	_, err := LoadModule(open, "/plugins/music.so")
	if !errors.Is(err, ErrPluginContractInvalid) {
		t.Fatalf("error = %v, want ErrPluginContractInvalid", err)
	}
}

func TestLoadModuleRejectsEmptyPath(t *testing.T) {
	open := func(string) (PluginSymbols, error) {
		t.Fatal("open must not be called for an empty path")
		return nil, nil
	}

	_, err := LoadModule(open, "")
	if !errors.Is(err, ErrPluginPathEmpty) {
		t.Fatalf("error = %v, want ErrPluginPathEmpty", err)
	}
}
