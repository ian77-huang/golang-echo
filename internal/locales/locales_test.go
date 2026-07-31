package locales

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestEmbeddedTomlFilesExist(t *testing.T) {
	entries, err := FS.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected embedded locale files")
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".toml" {
			t.Fatalf("unexpected embedded file %q", entry.Name())
		}
		data, err := FS.ReadFile(entry.Name())
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		if strings.TrimSpace(string(data)) == "" {
			t.Fatalf("locale file %s is empty", entry.Name())
		}
	}
}

func TestEmbeddedLocaleFilesIncludeExpectedLocales(t *testing.T) {
	data, err := FS.ReadFile("active.en.toml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "welcome") {
		t.Fatalf("active.en.toml missing expected keys: %s", data)
	}
}
