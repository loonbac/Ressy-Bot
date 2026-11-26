package core

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileModuleSettingsDefaultsToEnabled(t *testing.T) {
	settings, err := NewFileModuleSettings(filepath.Join(t.TempDir(), "settings.json"))
	if err != nil {
		t.Fatalf("NewFileModuleSettings returned error: %v", err)
	}
	if !settings.IsEnabled("guild-1", "music") {
		t.Fatal("new module was disabled by default")
	}
}

func TestFileModuleSettingsPersistsDisabledModule(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	settings, err := NewFileModuleSettings(path)
	if err != nil {
		t.Fatalf("NewFileModuleSettings returned error: %v", err)
	}
	if err := settings.SetEnabled("guild-1", "music", false); err != nil {
		t.Fatalf("SetEnabled returned error: %v", err)
	}

	reloaded, err := NewFileModuleSettings(path)
	if err != nil {
		t.Fatalf("reload returned error: %v", err)
	}
	if reloaded.IsEnabled("guild-1", "music") {
		t.Fatal("disabled module became enabled after reload")
	}
	if !reloaded.IsEnabled("guild-2", "music") {
		t.Fatal("another guild inherited the disabled setting")
	}
}

func TestFileModuleSettingsListsNamesInOrder(t *testing.T) {
	settings, err := NewFileModuleSettings(filepath.Join(t.TempDir(), "settings.json"))
	if err != nil {
		t.Fatalf("NewFileModuleSettings returned error: %v", err)
	}
	for _, name := range []string{"music", "admin", "music"} {
		if err := settings.SetEnabled("guild-1", ModuleName(name), false); err != nil {
			t.Fatalf("SetEnabled(%q) returned error: %v", name, err)
		}
	}

	got := settings.List("guild-1")
	want := []ModuleName{"admin", "music"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List = %v, want %v", got, want)
	}
}

func TestFileModuleSettingsReportsUnwritablePath(t *testing.T) {
	_, err := NewFileModuleSettings(filepath.Join(t.TempDir(), "absent", "settings.json"))
	if err == nil {
		t.Fatal("NewFileModuleSettings returned no error for a missing parent directory")
	}
}
