package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStartWithEmptyDirectory(t *testing.T) {
	catalog, modules, err := Start(Config{PluginDirectory: t.TempDir()})
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if len(modules) != 0 {
		t.Fatalf("loaded %d modules, want 0", len(modules))
	}
	if len(catalog.List()) != 0 {
		t.Fatalf("registered %d commands, want 0", len(catalog.List()))
	}
}

func TestStartIgnoresFilesThatAreNotPlugins(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "notes.txt"), []byte("ignored"), 0o600); err != nil {
		t.Fatalf("write note: %v", err)
	}

	catalog, modules, err := Start(Config{PluginDirectory: directory})
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if len(modules) != 0 || len(catalog.List()) != 0 {
		t.Fatalf("loaded %d modules and %d commands, want 0 and 0", len(modules), len(catalog.List()))
	}
}

func TestStartReportsMissingDirectory(t *testing.T) {
	_, _, err := Start(Config{PluginDirectory: filepath.Join(t.TempDir(), "absent")})
	if err == nil {
		t.Fatal("Start returned no error for a missing plugin directory")
	}
}
