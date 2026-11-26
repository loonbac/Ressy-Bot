package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// ModuleSettings decides which modules are active in each guild.
type ModuleSettings interface {
	IsEnabled(guild GuildID, module ModuleName) bool
	SetEnabled(guild GuildID, module ModuleName, enabled bool) error
	List(guild GuildID) []ModuleName
}

type fileModuleSettings struct {
	mu     sync.Mutex
	path   string
	guilds map[GuildID]map[ModuleName]bool
}

// NewFileModuleSettings opens a JSON-backed store of per-guild module settings.
func NewFileModuleSettings(path string) (ModuleSettings, error) {
	store := &fileModuleSettings{
		path:   path,
		guilds: make(map[GuildID]map[ModuleName]bool),
	}
	parent := filepath.Dir(path)
	info, err := os.Stat(parent)
	if err != nil {
		return nil, fmt.Errorf("module settings directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("module settings parent %q is not a directory", parent)
	}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &store.guilds); err != nil {
			return nil, fmt.Errorf("parse module settings: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read module settings: %w", err)
	}
	return store, nil
}

func (s *fileModuleSettings) IsEnabled(guild GuildID, module ModuleName) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if modules, exists := s.guilds[guild]; exists {
		if enabled, configured := modules[module]; configured {
			return enabled
		}
	}
	return true
}

func (s *fileModuleSettings) SetEnabled(guild GuildID, module ModuleName, enabled bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := make(map[GuildID]map[ModuleName]bool, len(s.guilds))
	for existingGuild, modules := range s.guilds {
		copied := make(map[ModuleName]bool, len(modules))
		for name, value := range modules {
			copied[name] = value
		}
		next[existingGuild] = copied
	}
	if next[guild] == nil {
		next[guild] = make(map[ModuleName]bool)
	}
	next[guild][module] = enabled

	data, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return fmt.Errorf("encode module settings: %w", err)
	}
	temporary := s.path + ".tmp"
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return fmt.Errorf("write module settings: %w", err)
	}
	if err := os.Rename(temporary, s.path); err != nil {
		return fmt.Errorf("replace module settings: %w", err)
	}
	s.guilds = next
	return nil
}

func (s *fileModuleSettings) List(guild GuildID) []ModuleName {
	s.mu.Lock()
	defer s.mu.Unlock()
	modules := s.guilds[guild]
	names := make([]ModuleName, 0, len(modules))
	for name := range modules {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })
	return names
}
